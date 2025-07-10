package ecs

import (
	"fmt"
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/jonesrussell/gimbal/assets"
	"github.com/jonesrussell/gimbal/internal/common"
)

var defaultFontFace font.Face

func init() {
	fontBytes, err := assets.Assets.ReadFile("fonts/PressStart2P.ttf")
	if err != nil {
		log.Fatalf("failed to read font: %v", err)
	}
	fontTTF, err := opentype.Parse(fontBytes)
	if err != nil {
		log.Fatalf("failed to parse font: %v", err)
	}
	defaultFontFace, err = opentype.NewFace(fontTTF, &opentype.FaceOptions{
		Size:    16,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatalf("failed to create font face: %v", err)
	}
}

// drawCenteredText draws text centered on screen (helper method for scenes)
func drawCenteredText(screen *ebiten.Image, textStr string, x, y, alpha float64) {
	bounds := text.BoundString(defaultFontFace, textStr)
	w := bounds.Max.X - bounds.Min.X
	h := bounds.Max.Y - bounds.Min.Y
	col := color.RGBA{255, 255, 255, uint8(255 * alpha)}
	text.Draw(screen, textStr, defaultFontFace, int(x)-w/2, int(y)+h/2, col)
}

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
	inputHandler common.GameInputHandler
}

// NewSceneManager creates a new scene manager
func NewSceneManager(world donburi.World, config *common.GameConfig, logger common.Logger, inputHandler common.GameInputHandler) *SceneManager {
	sm := &SceneManager{
		scenes:       make(map[SceneType]Scene),
		world:        world,
		config:       config,
		logger:       logger,
		inputHandler: inputHandler,
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

// GetInputHandler returns the input handler
func (sm *SceneManager) GetInputHandler() common.GameInputHandler {
	return sm.inputHandler
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
	minTime   float64
	maxTime   float64
	finished  bool
}

// NewStudioIntroScene creates a new studio intro scene
func NewStudioIntroScene(manager *SceneManager) *StudioIntroScene {
	return &StudioIntroScene{
		manager:   manager,
		startTime: time.Now(),
		minTime:   2.0, // Minimum 2 seconds
		maxTime:   4.0, // Maximum 4 seconds
		finished:  false,
	}
}

func (s *StudioIntroScene) Update() error {
	elapsed := time.Since(s.startTime).Seconds()
	if s.finished {
		return nil
	}
	// Allow skip after minTime with any key or mouse
	if elapsed >= s.minTime {
		input := s.manager.GetInputHandler()
		if input != nil && (input.GetLastEvent() != common.InputEventNone) {
			s.finished = true
			s.manager.SwitchScene(SceneTitleScreen)
			return nil
		}
	}
	// Auto-advance after maxTime
	if elapsed >= s.maxTime {
		s.finished = true
		s.manager.SwitchScene(SceneTitleScreen)
	}
	return nil
}

func (s *StudioIntroScene) Draw(screen *ebiten.Image) {
	// Clear screen with black background
	screen.Fill(color.Black)

	// Calculate fade-in effect
	elapsed := time.Since(s.startTime).Seconds()
	fadeProgress := elapsed / s.maxTime
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
