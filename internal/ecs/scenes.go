package ecs

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"

	"github.com/jonesrussell/gimbal/internal/common"
)

// SceneType represents different game scenes
type SceneType int

const (
	SceneMenu SceneType = iota
	ScenePlaying
	ScenePaused
	SceneGameOver
	SceneVictory
)

// Scene represents a game scene with its own update and draw logic
type Scene interface {
	Update() error
	Draw(screen *ebiten.Image)
	Enter()
	Exit()
	GetType() SceneType
}

// SceneManager manages different game scenes
type SceneManager struct {
	currentScene Scene
	scenes       map[SceneType]Scene
	world        donburi.World
	config       *common.GameConfig
	logger       common.Logger
}

// NewSceneManager creates a new scene manager
func NewSceneManager(world donburi.World, config *common.GameConfig, logger common.Logger) *SceneManager {
	sm := &SceneManager{
		scenes: make(map[SceneType]Scene),
		world:  world,
		config: config,
		logger: logger,
	}

	// Initialize scenes
	sm.scenes[SceneMenu] = NewMenuScene(sm)
	sm.scenes[ScenePlaying] = NewPlayingScene(sm)
	sm.scenes[ScenePaused] = NewPausedScene(sm)
	sm.scenes[SceneGameOver] = NewGameOverScene(sm)

	// Set initial scene
	sm.currentScene = sm.scenes[SceneMenu]
	sm.currentScene.Enter()

	return sm
}

// Update updates the current scene
func (sm *SceneManager) Update() error {
	return sm.currentScene.Update()
}

// Draw draws the current scene
func (sm *SceneManager) Draw(screen *ebiten.Image) {
	sm.currentScene.Draw(screen)
}

// SwitchScene switches to a different scene
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

// GetCurrentScene returns the current scene
func (sm *SceneManager) GetCurrentScene() Scene {
	return sm.currentScene
}

// GetWorld returns the ECS world
func (sm *SceneManager) GetWorld() donburi.World {
	return sm.world
}

// GetConfig returns the game config
func (sm *SceneManager) GetConfig() *common.GameConfig {
	return sm.config
}

// GetLogger returns the logger
func (sm *SceneManager) GetLogger() common.Logger {
	return sm.logger
}

// MenuScene represents the main menu scene
type MenuScene struct {
	manager *SceneManager
}

// NewMenuScene creates a new menu scene
func NewMenuScene(manager *SceneManager) *MenuScene {
	return &MenuScene{manager: manager}
}

func (s *MenuScene) Update() error {
	// Handle menu input (simplified for now)
	// In a real implementation, you'd handle menu navigation
	return nil
}

func (s *MenuScene) Draw(screen *ebiten.Image) {
	// Draw menu background
	screen.Fill(color.Black)

	// Draw menu text (simplified)
	// In a real implementation, you'd use a proper text rendering system
}

func (s *MenuScene) Enter() {
	s.manager.logger.Debug("Entering menu scene")
}

func (s *MenuScene) Exit() {
	s.manager.logger.Debug("Exiting menu scene")
}

func (s *MenuScene) GetType() SceneType {
	return SceneMenu
}

// PlayingScene represents the main gameplay scene
type PlayingScene struct {
	manager *SceneManager
}

// NewPlayingScene creates a new playing scene
func NewPlayingScene(manager *SceneManager) *PlayingScene {
	return &PlayingScene{manager: manager}
}

func (s *PlayingScene) Update() error {
	// This will be handled by the main game loop
	return nil
}

func (s *PlayingScene) Draw(screen *ebiten.Image) {
	// This will be handled by the main game loop
}

func (s *PlayingScene) Enter() {
	s.manager.logger.Debug("Entering playing scene")
}

func (s *PlayingScene) Exit() {
	s.manager.logger.Debug("Exiting playing scene")
}

func (s *PlayingScene) GetType() SceneType {
	return ScenePlaying
}

// PausedScene represents the paused game scene
type PausedScene struct {
	manager *SceneManager
}

// NewPausedScene creates a new paused scene
func NewPausedScene(manager *SceneManager) *PausedScene {
	return &PausedScene{manager: manager}
}

func (s *PausedScene) Update() error {
	// Handle pause menu input
	return nil
}

func (s *PausedScene) Draw(screen *ebiten.Image) {
	// Draw pause overlay
	screen.Fill(color.Black)

	// Draw pause text (simplified)
}

func (s *PausedScene) Enter() {
	s.manager.logger.Debug("Entering paused scene")
}

func (s *PausedScene) Exit() {
	s.manager.logger.Debug("Exiting paused scene")
}

func (s *PausedScene) GetType() SceneType {
	return ScenePaused
}

// GameOverScene represents the game over scene
type GameOverScene struct {
	manager *SceneManager
}

// NewGameOverScene creates a new game over scene
func NewGameOverScene(manager *SceneManager) *GameOverScene {
	return &GameOverScene{manager: manager}
}

func (s *GameOverScene) Update() error {
	// Handle game over input
	return nil
}

func (s *GameOverScene) Draw(screen *ebiten.Image) {
	// Draw game over screen
	screen.Fill(color.Black)

	// Draw game over text (simplified)
}

func (s *GameOverScene) Enter() {
	s.manager.logger.Debug("Entering game over scene")
}

func (s *GameOverScene) Exit() {
	s.manager.logger.Debug("Exiting game over scene")
}

func (s *GameOverScene) GetType() SceneType {
	return SceneGameOver
}
