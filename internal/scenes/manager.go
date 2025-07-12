package scenes

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/yohamta/donburi"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/config"
	"github.com/jonesrussell/gimbal/internal/ecs/core"
	"github.com/jonesrussell/gimbal/internal/ecs/debug"
	"github.com/jonesrussell/gimbal/internal/ecs/managers"
	resources "github.com/jonesrussell/gimbal/internal/ecs/managers/resource"
	"github.com/jonesrussell/gimbal/internal/errors"
)

// SceneManagerConfig groups all dependencies for SceneManager
// to avoid argument limit lint violations
type SceneManagerConfig struct {
	World        donburi.World
	Config       *config.GameConfig
	Logger       common.Logger
	InputHandler common.GameInputHandler
	Font         text.Face
	ScoreManager *managers.ScoreManager
	ResourceMgr  *resources.ResourceManager
}

type SceneType int

const (
	SceneStudioIntro SceneType = iota
	SceneTitleScreen
	SceneMenu
	ScenePlaying
	ScenePaused
	SceneGameOver
	SceneVictory
	SceneOptions
	SceneCredits
)

type Scene interface {
	Update() error
	Draw(screen *ebiten.Image)
	Enter()
	Exit()
	GetType() SceneType
}

type SceneManager struct {
	currentScene    Scene
	scenes          map[SceneType]Scene
	world           donburi.World
	config          *config.GameConfig
	logger          common.Logger
	inputHandler    common.GameInputHandler
	onResume        func()      // Callback to unpause game state
	healthSystem    interface{} // Health system interface for scenes to access
	font            text.Face
	resourceMgr     *resources.ResourceManager
	debugRenderer   *debug.DebugRenderer
	renderOptimizer *core.RenderOptimizer
	imagePool       *core.ImagePool
}

func NewSceneManager(cfg *SceneManagerConfig) *SceneManager {
	sceneMgr := &SceneManager{
		scenes:       make(map[SceneType]Scene),
		world:        cfg.World,
		config:       cfg.Config,
		logger:       cfg.Logger,
		inputHandler: cfg.InputHandler,
		font:         cfg.Font,
		resourceMgr:  cfg.ResourceMgr,
	}

	// Initialize debug renderer
	sceneMgr.debugRenderer = debug.NewDebugRenderer(cfg.Config, cfg.Logger)
	sceneMgr.debugRenderer.SetFont(cfg.Font)

	// Register all scenes
	sceneMgr.scenes[SceneStudioIntro] = NewStudioIntroScene(sceneMgr, cfg.Font)
	sceneMgr.scenes[SceneTitleScreen] = NewTitleScreenScene(sceneMgr, cfg.Font)
	sceneMgr.scenes[SceneMenu] = NewMenuScene(sceneMgr, cfg.Font)
	sceneMgr.scenes[ScenePlaying] = NewPlayingScene(sceneMgr, cfg.Font, cfg.ScoreManager, cfg.ResourceMgr)
	sceneMgr.scenes[ScenePaused] = NewPausedScene(sceneMgr, cfg.Font)
	sceneMgr.scenes[SceneGameOver] = NewGameOverScene(sceneMgr, cfg.Font)
	sceneMgr.scenes[SceneCredits] = NewSimpleTextScene(
		sceneMgr,
		"CREDITS\nGimbal Studios\n2025",
		SceneCredits,
		cfg.Font,
	)
	sceneMgr.scenes[SceneOptions] = NewSimpleTextScene(sceneMgr, "OPTIONS\nComing Soon!", SceneOptions, cfg.Font)

	return sceneMgr
}

func (sceneMgr *SceneManager) Update() error {
	// Handle debug controls
	if sceneMgr.inputHandler != nil {
		if handler, ok := sceneMgr.inputHandler.(interface {
			IsDebugTogglePressed() bool
			IsDebugLevelCyclePressed() bool
		}); ok {
			if handler.IsDebugTogglePressed() {
				sceneMgr.debugRenderer.Toggle()
			}
			if handler.IsDebugLevelCyclePressed() {
				sceneMgr.debugRenderer.CycleLevel()
			}
		}
	}

	if sceneMgr.currentScene == nil {
		return nil // No scene set yet, nothing to update
	}
	return sceneMgr.currentScene.Update()
}

func (sceneMgr *SceneManager) Draw(screen *ebiten.Image) {
	if sceneMgr.currentScene == nil {
		return // No scene set yet, nothing to draw
	}
	sceneMgr.currentScene.Draw(screen)

	// Draw debug overlay on top of everything
	sceneMgr.debugRenderer.Render(screen, sceneMgr.world)
}

func (sceneMgr *SceneManager) SetInitialScene(sceneType SceneType) error {
	if scene, exists := sceneMgr.scenes[sceneType]; exists {
		sceneMgr.currentScene = scene
		sceneMgr.currentScene.Enter()
		sceneMgr.logger.Debug("Initial scene set", "scene_type", sceneType)
		return nil
	} else {
		sceneMgr.logger.Error("Scene not found for initial scene", "scene_type", sceneType)
		return errors.NewGameError(errors.ErrorCodeSceneNotFound, "initial scene not found")
	}
}

func (sceneMgr *SceneManager) SwitchScene(sceneType SceneType) {
	if scene, exists := sceneMgr.scenes[sceneType]; exists {
		if sceneMgr.currentScene != nil {
			sceneMgr.logger.Debug("Switching scene",
				"from", sceneMgr.currentScene.GetType(),
				"to", sceneType)
			sceneMgr.currentScene.Exit()
		} else {
			sceneMgr.logger.Debug("Setting initial scene", "scene_type", sceneType)
		}
		sceneMgr.currentScene = scene
		if sceneType == ScenePlaying {
			sceneMgr.logger.Debug("Entering gameplay scene",
				"player_entity", sceneMgr.world, // Replace with actual player entity if accessible
				"health_system", sceneMgr.healthSystem != nil,
				"resource_mgr", sceneMgr.resourceMgr != nil)
		}
		sceneMgr.currentScene.Enter()
	} else {
		sceneMgr.logger.Error("Scene not found", "scene_type", sceneType)
	}
}

func (sceneMgr *SceneManager) GetCurrentScene() Scene {
	return sceneMgr.currentScene
}

func (sceneMgr *SceneManager) GetWorld() donburi.World {
	return sceneMgr.world
}

func (sceneMgr *SceneManager) GetConfig() *config.GameConfig {
	return sceneMgr.config
}

func (sceneMgr *SceneManager) GetLogger() common.Logger {
	return sceneMgr.logger
}

func (sceneMgr *SceneManager) GetInputHandler() common.GameInputHandler {
	return sceneMgr.inputHandler
}

// SetResumeCallback sets the callback function to unpause game state
func (sceneMgr *SceneManager) SetResumeCallback(callback func()) {
	sceneMgr.onResume = callback
}

// SetHealthSystem sets the health system for scenes to access
func (sceneMgr *SceneManager) SetHealthSystem(healthSystem interface{}) {
	sceneMgr.healthSystem = healthSystem
}

// GetHealthSystem returns the health system
func (sceneMgr *SceneManager) GetHealthSystem() interface{} {
	return sceneMgr.healthSystem
}

// GetResourceManager returns the resource manager
func (sceneMgr *SceneManager) GetResourceManager() *resources.ResourceManager {
	return sceneMgr.resourceMgr
}

// GetDebugRenderer returns the debug renderer
func (sceneMgr *SceneManager) GetDebugRenderer() *debug.DebugRenderer {
	return sceneMgr.debugRenderer
}

// GetRenderOptimizer returns the render optimizer
func (sceneMgr *SceneManager) GetRenderOptimizer() *core.RenderOptimizer {
	return sceneMgr.renderOptimizer
}

// SetRenderOptimizer sets the render optimizer
func (sceneMgr *SceneManager) SetRenderOptimizer(optimizer *core.RenderOptimizer) {
	sceneMgr.renderOptimizer = optimizer
}

// GetImagePool returns the image pool
func (sceneMgr *SceneManager) GetImagePool() *core.ImagePool {
	return sceneMgr.imagePool
}

// SetImagePool sets the image pool
func (sceneMgr *SceneManager) SetImagePool(pool *core.ImagePool) {
	sceneMgr.imagePool = pool
}
