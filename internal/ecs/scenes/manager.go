package ecs

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

	// Initialize scenes (to be set up in main or via factory)
	// Example scene registration (add this where scenes are registered):
	sceneMgr.scenes[SceneCredits] = NewSimpleTextScene(sceneMgr, "CREDITS\nGimbal Studios\n2025", SceneCredits)
	sceneMgr.scenes[SceneOptions] = NewSimpleTextScene(sceneMgr, "OPTIONS\nComing Soon!", SceneOptions)
	return sceneMgr
}

func (sceneMgr *SceneManager) Update() error {
	return sceneMgr.currentScene.Update()
}

func (sceneMgr *SceneManager) Draw(screen *ebiten.Image) {
	sceneMgr.currentScene.Draw(screen)
}

func (sceneMgr *SceneManager) SwitchScene(sceneType SceneType) {
	if scene, exists := sceneMgr.scenes[sceneType]; exists {
		sceneMgr.logger.Debug("Switching scene",
			"from", sceneMgr.currentScene.GetType(),
			"to", sceneType)

		sceneMgr.currentScene.Exit()
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
