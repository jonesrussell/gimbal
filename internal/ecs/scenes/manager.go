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
	sm := &SceneManager{
		scenes:       make(map[SceneType]Scene),
		world:        world,
		config:       config,
		logger:       logger,
		inputHandler: inputHandler,
	}

	// Initialize scenes (to be set up in main or via factory)
	// Example scene registration (add this where scenes are registered):
	sm.scenes[SceneCredits] = NewSimpleTextScene(sm, "CREDITS\nGimbal Studios\n2025", SceneCredits)
	sm.scenes[SceneOptions] = NewSimpleTextScene(sm, "OPTIONS\nComing Soon!", SceneOptions)
	return sm
}

func (sm *SceneManager) Update() error {
	return sm.currentScene.Update()
}

func (sm *SceneManager) Draw(screen *ebiten.Image) {
	sm.currentScene.Draw(screen)
}

func (sm *SceneManager) SwitchScene(sceneType SceneType) {
	if scene, exists := sm.scenes[sceneType]; exists {
		sm.logger.Debug("Switching scene",
			"from", sm.currentScene.GetType(),
			"to", sceneType)

		sm.currentScene.Exit()
		sm.currentScene = scene
		sm.currentScene.Enter()
	} else {
		sm.logger.Error("Scene not found", "scene_type", sceneType)
	}
}

func (sm *SceneManager) GetCurrentScene() Scene {
	return sm.currentScene
}

func (sm *SceneManager) GetWorld() donburi.World {
	return sm.world
}

func (sm *SceneManager) GetConfig() *common.GameConfig {
	return sm.config
}

func (sm *SceneManager) GetLogger() common.Logger {
	return sm.logger
}

func (sm *SceneManager) GetInputHandler() common.GameInputHandler {
	return sm.inputHandler
}
