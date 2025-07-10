package ecs

import (
	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/jonesrussell/gimbal/internal/common"
)

// SceneType represents different game scenes
type SceneType int

const (
	SceneStudioIntro SceneType = iota
	SceneTitleScreen
	SceneMenu
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
	sm.scenes[SceneStudioIntro] = NewStudioIntroScene(sm)
	sm.scenes[SceneTitleScreen] = NewTitleScreenScene(sm)
	sm.scenes[SceneMenu] = NewMenuScene(sm)
	sm.scenes[ScenePlaying] = NewPlayingScene(sm)
	sm.scenes[ScenePaused] = NewPausedScene(sm)
	sm.scenes[SceneGameOver] = NewGameOverScene(sm)

	// Set initial scene
	sm.currentScene = sm.scenes[SceneStudioIntro]
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
	// Clear screen
	screen.Fill(color.Black)

	// Run render system through wrapper
	renderWrapper := NewRenderSystemWrapper(screen)
	if err := renderWrapper.Update(s.manager.world); err != nil {
		s.manager.logger.Error("Render system failed", "error", err)
	}

	// Draw debug info if enabled
	if s.manager.config.Debug {
		s.drawDebugInfo(screen)
	}
}

// drawDebugInfo renders debug information
func (s *PlayingScene) drawDebugInfo(screen *ebiten.Image) {
	// Get player info for debug display
	players := make([]donburi.Entity, 0)
	query.NewQuery(
		filter.And(
			filter.Contains(PlayerTag),
			filter.Contains(Position),
			filter.Contains(Orbital),
		),
	).Each(s.manager.world, func(entry *donburi.Entry) {
		players = append(players, entry.Entity())
	})

	if len(players) > 0 {
		playerEntry := s.manager.world.Entry(players[0])
		if playerEntry.Valid() {
			pos := Position.Get(playerEntry)
			orb := Orbital.Get(playerEntry)

			// Log debug info
			s.manager.logger.Debug("Debug Info",
				"player_pos", fmt.Sprintf("(%.1f, %.1f)", pos.X, pos.Y),
				"player_angle", fmt.Sprintf("%.1f°", orb.OrbitalAngle),
				"entity_count", s.manager.world.Len(),
			)
		}
	}
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

// StudioIntroScene represents the studio intro scene
type StudioIntroScene struct {
	manager   *SceneManager
	startTime time.Time
	duration  float64
}

// NewStudioIntroScene creates a new studio intro scene
func NewStudioIntroScene(manager *SceneManager) *StudioIntroScene {
	return &StudioIntroScene{
		manager:   manager,
		startTime: time.Now(),
		duration:  3.0, // 3 seconds
	}
}

func (s *StudioIntroScene) Update() error {
	// Auto-transition after duration
	if time.Since(s.startTime).Seconds() >= s.duration {
		s.manager.SwitchScene(SceneTitleScreen)
	}
	return nil
}

func (s *StudioIntroScene) Draw(screen *ebiten.Image) {
	// Clear screen with black background
	screen.Fill(color.Black)

	// Calculate fade-in effect
	elapsed := time.Since(s.startTime).Seconds()
	fadeProgress := elapsed / s.duration
	if fadeProgress > 1.0 {
		fadeProgress = 1.0
	}

	// Draw "Gimbal Studios" text
	drawCenteredText(screen, "GIMBAL STUDIOS",
		float64(s.manager.config.ScreenSize.Width)/2,
		float64(s.manager.config.ScreenSize.Height)/2,
		fadeProgress)

	// Draw subtitle
	drawCenteredText(screen, "Presents",
		float64(s.manager.config.ScreenSize.Width)/2,
		float64(s.manager.config.ScreenSize.Height)/2+50,
		fadeProgress*0.8)
}

func (s *StudioIntroScene) Enter() {
	s.manager.logger.Debug("Entering studio intro scene")
	s.startTime = time.Now()
}

func (s *StudioIntroScene) Exit() {
	s.manager.logger.Debug("Exiting studio intro scene")
}

func (s *StudioIntroScene) GetType() SceneType {
	return SceneStudioIntro
}

// TitleScreenScene represents the title screen scene
type TitleScreenScene struct {
	manager   *SceneManager
	startTime time.Time
}

// NewTitleScreenScene creates a new title screen scene
func NewTitleScreenScene(manager *SceneManager) *TitleScreenScene {
	return &TitleScreenScene{
		manager:   manager,
		startTime: time.Now(),
	}
}

func (s *TitleScreenScene) Update() error {
	// Check for any key press to continue
	// This will be handled by input system
	return nil
}

func (s *TitleScreenScene) Draw(screen *ebiten.Image) {
	// Clear screen with space background
	screen.Fill(color.Black)

	// Draw game title
	drawCenteredText(screen, "GIMBAL",
		float64(s.manager.config.ScreenSize.Width)/2,
		float64(s.manager.config.ScreenSize.Height)/2-50,
		1.0)

	// Draw subtitle
	drawCenteredText(screen, "Exoplanetary Gyruss-Inspired Shooter",
		float64(s.manager.config.ScreenSize.Width)/2,
		float64(s.manager.config.ScreenSize.Height)/2,
		1.0)

	// Draw "Press any key to continue" with blinking effect
	elapsed := time.Since(s.startTime).Seconds()
	blink := (elapsed * 2) < 1.0 // Blink every 0.5 seconds
	if blink {
		drawCenteredText(screen, "Press any key to continue",
			float64(s.manager.config.ScreenSize.Width)/2,
			float64(s.manager.config.ScreenSize.Height)/2+100,
			1.0)
	}
}

func (s *TitleScreenScene) Enter() {
	s.manager.logger.Debug("Entering title screen scene")
	s.startTime = time.Now()
}

func (s *TitleScreenScene) Exit() {
	s.manager.logger.Debug("Exiting title screen scene")
}

func (s *TitleScreenScene) GetType() SceneType {
	return SceneTitleScreen
}

// drawCenteredText draws text centered on screen (helper method for scenes)
func drawCenteredText(screen *ebiten.Image, text string, x, y, alpha float64) {
	// Create a simple "text" representation (colored rectangle)
	textWidth := len(text) * 8 // Approximate character width
	textHeight := 16

	img := ebiten.NewImage(textWidth, textHeight)
	img.Fill(color.RGBA{255, 255, 255, uint8(255 * alpha)})

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x-float64(textWidth)/2, y-float64(textHeight)/2)
	screen.DrawImage(img, op)
}

// MenuScene represents the main menu scene
type MenuScene struct {
	manager   *SceneManager
	selection int
	options   []string
}

// NewMenuScene creates a new menu scene
func NewMenuScene(manager *SceneManager) *MenuScene {
	return &MenuScene{
		manager:   manager,
		selection: 0,
		options:   []string{"Start Game", "Options", "Credits", "Quit"},
	}
}

func (s *MenuScene) Update() error {
	// Handle menu navigation
	// This will be handled by input system
	return nil
}

func (s *MenuScene) Draw(screen *ebiten.Image) {
	// Draw menu background
	screen.Fill(color.Black)

	// Draw game title
	drawCenteredText(screen, "GIMBAL",
		float64(s.manager.config.ScreenSize.Width)/2,
		100,
		1.0)

	// Draw menu options
	menuY := float64(s.manager.config.ScreenSize.Height) / 2
	for i, option := range s.options {
		y := menuY + float64(i*40)
		alpha := 1.0
		if i == s.selection {
			alpha = 1.0 // Highlight selected option
		} else {
			alpha = 0.7 // Dim unselected options
		}

		drawCenteredText(screen, option,
			float64(s.manager.config.ScreenSize.Width)/2,
			y,
			alpha)
	}
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
