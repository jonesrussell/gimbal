package scenes

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"

	"github.com/jonesrussell/gimbal/internal/common"
)

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
	currentScene Scene
	scenes       map[SceneType]Scene
	world        donburi.World
	config       *common.GameConfig
	logger       common.Logger
	inputHandler common.GameInputHandler
}

func NewSceneManager(
	world donburi.World,
	config *common.GameConfig,
	logger common.Logger,
	inputHandler common.GameInputHandler,
) *SceneManager {
	sceneMgr := &SceneManager{
		scenes:       make(map[SceneType]Scene),
		world:        world,
		config:       config,
		logger:       logger,
		inputHandler: inputHandler,
	}

	// Register all scenes
	sceneMgr.scenes[SceneStudioIntro] = NewStudioIntroScene(sceneMgr)
	sceneMgr.scenes[SceneTitleScreen] = NewTitleScreenScene(sceneMgr)
	sceneMgr.scenes[SceneMenu] = NewMenuScene(sceneMgr)
	sceneMgr.scenes[ScenePlaying] = NewPlayingScene(sceneMgr)
	sceneMgr.scenes[ScenePaused] = NewPausedScene(sceneMgr)
	sceneMgr.scenes[SceneCredits] = NewSimpleTextScene(sceneMgr, "CREDITS\nGimbal Studios\n2025", SceneCredits)
	sceneMgr.scenes[SceneOptions] = NewSimpleTextScene(sceneMgr, "OPTIONS\nComing Soon!", SceneOptions)

	return sceneMgr
}

func (sceneMgr *SceneManager) Update() error {
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
}

func (sceneMgr *SceneManager) SetInitialScene(sceneType SceneType) error {
	if scene, exists := sceneMgr.scenes[sceneType]; exists {
		sceneMgr.currentScene = scene
		sceneMgr.currentScene.Enter()
		sceneMgr.logger.Debug("Initial scene set", "scene_type", sceneType)
		return nil
	} else {
		sceneMgr.logger.Error("Scene not found for initial scene", "scene_type", sceneType)
		return common.NewGameError(common.ErrorCodeSceneNotFound, "initial scene not found")
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

func (sceneMgr *SceneManager) GetConfig() *common.GameConfig {
	return sceneMgr.config
}

func (sceneMgr *SceneManager) GetLogger() common.Logger {
	return sceneMgr.logger
}

func (sceneMgr *SceneManager) GetInputHandler() common.GameInputHandler {
	return sceneMgr.inputHandler
}
