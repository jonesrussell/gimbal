package scenes

import (
	"image/color"
	"math"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// PauseOption represents the available pause menu options
type PauseOption int

const (
	PauseOptionResume PauseOption = iota
	PauseOptionReturnToMenu
	PauseOptionQuit
)

const (
	// Animation constants
	frameRate      = 60.0
	fadeInDuration = 0.3
	selectionDelay = 0.1

	// Layout constants
	titleScale     = 1.5
	titleY         = 80
	menuSpacing    = 50
	overlayAlpha   = 128
	paddingX       = 30
	paddingY       = 10
	hitAreaPadding = 50
	hintTextY      = 60

	// Color constants
	dimmedAlpha   = 0.7
	hintBaseAlpha = 0.6
)

// PausedScene manages the pause menu state and rendering
type PausedScene struct {
	manager           *SceneManager
	menu              *MenuSystem
	animationTime     float64
	selectionChanged  bool
	lastSelectionTime time.Time
	fadeIn            float64
	overlayImage      *ebiten.Image // Cached overlay image
	escWasPressed     bool          // Track if ESC was pressed when we entered
	canUnpause        bool          // Flag to allow unpausing
}

// NewPausedScene creates a new pause scene instance
func NewPausedScene(manager *SceneManager) *PausedScene {
	scene := &PausedScene{
		manager:           manager,
		animationTime:     0,
		selectionChanged:  false,
		lastSelectionTime: time.Now(),
		fadeIn:            0,
	}

	options := []MenuOption{
		{"Resume", func() { manager.SwitchScene(ScenePlaying) }},
		{"Return to Menu", func() { manager.SwitchScene(SceneMenu) }},
		{"Quit", func() { os.Exit(0) }},
	}

	config := PausedMenuConfig()
	config.MenuY = float64(manager.config.ScreenSize.Height) / 2

	scene.menu = NewMenuSystem(options, &config, manager.config.ScreenSize.Width, manager.config.ScreenSize.Height)

	// Create overlay image once (TODO: handle resizing if needed)
	scene.overlayImage = ebiten.NewImage(manager.config.ScreenSize.Width, manager.config.ScreenSize.Height)
	return scene
}

// Update handles input and animations for the pause scene
func (s *PausedScene) Update() error {
	dt := 1.0 / frameRate
	s.animationTime += dt
	s.updateFadeIn(dt)
	s.updateSelectionAnimation()
	s.handleInput()
	s.menu.Update()
	return nil
}

// updateFadeIn handles the fade-in animation
func (s *PausedScene) updateFadeIn(dt float64) {
	if s.fadeIn < 1.0 {
		s.fadeIn = math.Min(1.0, s.fadeIn+dt/fadeInDuration)
	}
}

// updateSelectionAnimation manages selection change animations
func (s *PausedScene) updateSelectionAnimation() {
	if s.selectionChanged && time.Since(s.lastSelectionTime).Seconds() > selectionDelay {
		s.selectionChanged = false
	}
}

// handleInput processes pause-specific input (ESC key)
func (s *PausedScene) handleInput() {
	currentEscPressed := ebiten.IsKeyPressed(ebiten.KeyEscape)
	escJustPressed := inpututil.IsKeyJustPressed(ebiten.KeyEscape)

	// Debug logging every few frames
	if int(s.animationTime*60)%30 == 0 { // Log every half second
		s.manager.logger.Debug("Pause input state",
			"escWasPressed", s.escWasPressed,
			"canUnpause", s.canUnpause,
			"currentEscPressed", currentEscPressed,
			"escJustPressed", escJustPressed)
	}

	// If ESC was pressed when we entered, wait for it to be released
	if s.escWasPressed && currentEscPressed {
		s.manager.logger.Debug("Waiting for ESC release")
		return // ESC is still held down from the pause action
	}

	// ESC has been released (or wasn't pressed when we entered)
	if s.escWasPressed && !currentEscPressed {
		s.escWasPressed = false
		s.canUnpause = true
		s.manager.logger.Debug("ESC released, can now unpause")
		return // Don't process input this frame, just mark as ready
	}

	// Now we can check for new ESC presses
	if s.canUnpause && escJustPressed {
		s.manager.logger.Debug("ESC pressed, unpausing game")
		s.manager.SwitchScene(ScenePlaying)
	}

	// If we entered without ESC pressed, we can unpause immediately
	if !s.escWasPressed {
		s.canUnpause = true
		if escJustPressed {
			s.manager.logger.Debug("ESC pressed (immediate), unpausing game")
			s.manager.SwitchScene(ScenePlaying)
		}
	}
}

// Draw renders the pause menu
func (s *PausedScene) Draw(screen *ebiten.Image) {
	s.drawOverlay(screen)
	s.drawTitle(screen)
	s.menu.Draw(screen, s.fadeIn)
	s.drawHintText(screen)
}

// drawOverlay renders the semi-transparent background overlay
func (s *PausedScene) drawOverlay(screen *ebiten.Image) {
	alpha := uint8(overlayAlpha * s.fadeIn)
	s.overlayImage.Fill(color.RGBA{0, 0, 0, alpha})
	screen.DrawImage(s.overlayImage, &ebiten.DrawImageOptions{})
}

// drawTitle renders the "PAUSED" title with pulsing animation
func (s *PausedScene) drawTitle(screen *ebiten.Image) {
	titleAlpha := s.fadeIn
	// Add subtle pulsing effect to title
	pulse := 0.9 + 0.1*math.Sin(s.animationTime*2.0)
	titleAlpha *= pulse

	op := &text.DrawOptions{}
	op.GeoM.Scale(titleScale, titleScale)
	op.GeoM.Translate(
		float64(s.manager.config.ScreenSize.Width)/2-75, // Adjust for scaling
		titleY,
	)
	op.ColorScale.SetR(0.2)
	op.ColorScale.SetG(0.8)
	op.ColorScale.SetB(1.0)
	op.ColorScale.SetA(float32(titleAlpha))

	text.Draw(screen, "PAUSED", defaultFontFace, op)
}

// drawHintText renders the hint text at the bottom of the screen
func (s *PausedScene) drawHintText(screen *ebiten.Image) {
	hintText := "Press ESC to resume quickly"
	hintAlpha := hintBaseAlpha * s.fadeIn * (0.8 + 0.2*math.Sin(s.animationTime*1.5))

	op := &text.DrawOptions{}
	width, _ := text.Measure(hintText, defaultFontFace, 0)
	op.GeoM.Translate(
		float64(s.manager.config.ScreenSize.Width)/2-width/2,
		float64(s.manager.config.ScreenSize.Height)-hintTextY,
	)
	op.ColorScale.SetR(0.8)
	op.ColorScale.SetG(0.8)
	op.ColorScale.SetB(0.8)
	op.ColorScale.SetA(float32(hintAlpha))

	text.Draw(screen, hintText, defaultFontFace, op)
}

// Enter is called when the scene becomes active
func (s *PausedScene) Enter() {
	s.manager.logger.Debug("Entering paused scene")
	s.fadeIn = 0
	s.animationTime = 0
	s.selectionChanged = false

	// Check and log ESC state when entering
	s.escWasPressed = ebiten.IsKeyPressed(ebiten.KeyEscape)
	s.canUnpause = false

	s.manager.logger.Debug("Pause scene ESC state",
		"escWasPressed", s.escWasPressed,
		"canUnpause", s.canUnpause)
}

// Exit is called when the scene becomes inactive
func (s *PausedScene) Exit() {
	s.manager.logger.Debug("Exiting paused scene")
}

// GetType returns the scene type identifier
func (s *PausedScene) GetType() SceneType {
	return ScenePaused
}
