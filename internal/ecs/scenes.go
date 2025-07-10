package ecs

import (
	"fmt"
	"image/color"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/jonesrussell/gimbal/internal/common"
)

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

func NewMenuScene(manager *SceneManager) *MenuScene {
	return &MenuScene{
		manager:   manager,
		selection: 0,
		options:   []string{"Start Game", "Options", "Credits", "Quit"},
	}
}

func (s *MenuScene) Update() error {
	// Keyboard navigation - use JustPressed to prevent rapid scrolling
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		s.selection = (s.selection - 1 + len(s.options)) % len(s.options)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		s.selection = (s.selection + 1) % len(s.options)
	}

	// Mouse hover
	x, y := ebiten.CursorPosition()
	menuY := float64(s.manager.config.ScreenSize.Height) / 2
	for i := range s.options {
		itemY := menuY + float64(i*40)
		width, height := text.Measure(s.options[i], defaultFontFace, 0)
		w := int(width)
		h := int(height)
		itemRect := struct{ x0, y0, x1, y1 int }{
			int(float64(s.manager.config.ScreenSize.Width)/2) - w/2 - 40, // extra for chevron
			int(itemY) - h/2 - 8,
			int(float64(s.manager.config.ScreenSize.Width)/2) + w/2 + 40,
			int(itemY) + h/2 + 8,
		}
		if x >= itemRect.x0 && x <= itemRect.x1 && y >= itemRect.y0 && y <= itemRect.y1 {
			s.selection = i
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				s.activateSelection()
			}
		}
	}

	// Keyboard select - use JustPressed to prevent multiple activations
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		s.activateSelection()
	}
	return nil
}

func (s *MenuScene) activateSelection() {
	switch s.selection {
	case 0: // Start Game
		s.manager.SwitchScene(ScenePlaying)
	case 1: // Options
		s.manager.SwitchScene(SceneOptions)
	case 2: // Credits
		s.manager.SwitchScene(SceneCredits)
	case 3: // Quit
		os.Exit(0)
	}
}

func (s *MenuScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)
	drawCenteredText(screen, "GIMBAL",
		float64(s.manager.config.ScreenSize.Width)/2,
		100, 1.0)

	menuY := float64(s.manager.config.ScreenSize.Height) / 2
	for i, option := range s.options {
		y := menuY + float64(i*40)
		alpha := 1.0
		bgAlpha := 0.0
		if i == s.selection {
			alpha = 1.0
			bgAlpha = 0.5
			// Animated chevron
			pulse := 0.7 + 0.3*float64((time.Now().UnixNano()/1e7)%20)/20.0
			chevron := ">"
			chevronOp := &text.DrawOptions{}
			chevronOp.GeoM.Translate(float64(int(float64(s.manager.config.ScreenSize.Width)/2)-120), float64(int(y)+8))
			chevronOp.ColorScale.SetR(0)
			chevronOp.ColorScale.SetG(1)
			chevronOp.ColorScale.SetB(1)
			chevronOp.ColorScale.SetA(float32(pulse))
			text.Draw(screen, chevron, defaultFontFace, chevronOp)
		}
		// Neon blue background highlight
		if bgAlpha > 0 {
			width, height := text.Measure(option, defaultFontFace, 0)
			w := int(width)
			h := int(height)
			paddingX := 24 // horizontal padding
			paddingY := 6  // vertical padding
			rectCol := color.RGBA{0, 255, 255, uint8(128 * bgAlpha)}
			rect := ebiten.NewImage(w+paddingX*2, h+paddingY*2)
			rect.Fill(rectCol)
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(
				float64(s.manager.config.ScreenSize.Width)/2-float64(w+paddingX*2)/2,
				y-float64(h+paddingY*2)/2+2, // fine-tuned for pixel-perfect alignment
			)
			screen.DrawImage(rect, op)
		}
		drawCenteredText(screen, option,
			float64(s.manager.config.ScreenSize.Width)/2, y, alpha)
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

// Placeholder scenes for Options and Credits
type OptionsScene struct{ manager *SceneManager }

func NewOptionsScene(manager *SceneManager) *OptionsScene { return &OptionsScene{manager: manager} }
func (s *OptionsScene) Update() error {
	if s.manager.inputHandler.GetLastEvent() != common.InputEventNone {
		s.manager.SwitchScene(SceneMenu)
	}
	return nil
}

func (s *OptionsScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)
	drawCenteredText(
		screen,
		"OPTIONS\nComing Soon!",
		float64(s.manager.config.ScreenSize.Width)/2,
		float64(s.manager.config.ScreenSize.Height)/2,
		1.0,
	)
}
func (s *OptionsScene) Enter()             { s.manager.logger.Debug("Entering options scene") }
func (s *OptionsScene) Exit()              { s.manager.logger.Debug("Exiting options scene") }
func (s *OptionsScene) GetType() SceneType { return SceneOptions }

type CreditsScene struct{ manager *SceneManager }

func NewCreditsScene(manager *SceneManager) *CreditsScene { return &CreditsScene{manager: manager} }
func (s *CreditsScene) Update() error {
	if s.manager.inputHandler.GetLastEvent() != common.InputEventNone {
		s.manager.SwitchScene(SceneMenu)
	}
	return nil
}

func (s *CreditsScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)
	drawCenteredText(
		screen,
		"CREDITS\nGimbal Studios\n2025",
		float64(s.manager.config.ScreenSize.Width)/2,
		float64(s.manager.config.ScreenSize.Height)/2,
		1.0,
	)
}
func (s *CreditsScene) Enter()             { s.manager.logger.Debug("Entering credits scene") }
func (s *CreditsScene) Exit()              { s.manager.logger.Debug("Exiting credits scene") }
func (s *CreditsScene) GetType() SceneType { return SceneCredits }
