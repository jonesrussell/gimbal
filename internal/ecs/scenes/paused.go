package ecs

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
	pulseSpeed     = 3.0
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
	selection         int
	options           []string
	animationTime     float64
	selectionChanged  bool
	lastSelectionTime time.Time
	fadeIn            float64
}

// NewPausedScene creates a new pause scene instance
func NewPausedScene(manager *SceneManager) *PausedScene {
	return &PausedScene{
		manager:           manager,
		selection:         0,
		options:           []string{"Resume", "Return to Menu", "Quit"},
		animationTime:     0,
		selectionChanged:  false,
		lastSelectionTime: time.Now(),
		fadeIn:            0,
	}
}

// Update handles input and animations for the pause scene
func (s *PausedScene) Update() error {
	dt := 1.0 / frameRate
	s.animationTime += dt

	s.updateFadeIn(dt)
	s.updateSelectionAnimation()
	s.handleInput()

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

// handleInput processes keyboard and mouse input
func (s *PausedScene) handleInput() {
	s.handleKeyboardInput()
	s.handleMouseInput()
}

// handleKeyboardInput processes keyboard navigation and actions
func (s *PausedScene) handleKeyboardInput() {
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		s.changeSelection((s.selection - 1 + len(s.options)) % len(s.options))
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		s.changeSelection((s.selection + 1) % len(s.options))
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		s.activateSelection()
	}

	// Quick resume with ESC
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		s.manager.SwitchScene(ScenePlaying)
	}
}

// handleMouseInput processes mouse hover and click events
func (s *PausedScene) handleMouseInput() {
	x, y := ebiten.CursorPosition()
	menuY := float64(s.manager.config.ScreenSize.Height) / 2
	hoveredItem := -1

	for i := range s.options {
		itemY := menuY + float64(i*menuSpacing)
		if s.isMouseOverItem(x, y, s.options[i], itemY) {
			hoveredItem = i
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				s.activateSelection()
			}
		}
	}

	if hoveredItem != -1 && hoveredItem != s.selection {
		s.changeSelection(hoveredItem)
	}
}

// isMouseOverItem checks if the mouse cursor is over a menu item
func (s *PausedScene) isMouseOverItem(x, y int, option string, itemY float64) bool {
	width, height := text.Measure(option, defaultFontFace, 0)
	w := int(width)
	h := int(height)

	// Create larger hit area for better UX
	itemRect := struct{ x0, y0, x1, y1 int }{
		int(float64(s.manager.config.ScreenSize.Width)/2) - w/2 - hitAreaPadding,
		int(itemY) - h/2 - 15,
		int(float64(s.manager.config.ScreenSize.Width)/2) + w/2 + hitAreaPadding,
		int(itemY) + h/2 + 15,
	}

	return x >= itemRect.x0 && x <= itemRect.x1 && y >= itemRect.y0 && y <= itemRect.y1
}

// changeSelection updates the selected menu item
func (s *PausedScene) changeSelection(newSelection int) {
	if newSelection != s.selection {
		s.selection = newSelection
		s.selectionChanged = true
		s.lastSelectionTime = time.Now()
	}
}

// activateSelection executes the action for the currently selected option
func (s *PausedScene) activateSelection() {
	switch PauseOption(s.selection) {
	case PauseOptionResume:
		s.manager.SwitchScene(ScenePlaying)
	case PauseOptionReturnToMenu:
		s.manager.SwitchScene(SceneMenu)
	case PauseOptionQuit:
		os.Exit(0)
	}
}

// Draw renders the pause menu
func (s *PausedScene) Draw(screen *ebiten.Image) {
	s.drawOverlay(screen)
	s.drawTitle(screen)
	s.drawMenuOptions(screen)
	s.drawHintText(screen)
}

// drawOverlay renders the semi-transparent background overlay
func (s *PausedScene) drawOverlay(screen *ebiten.Image) {
	alpha := uint8(overlayAlpha * s.fadeIn)
	overlay := ebiten.NewImage(s.manager.config.ScreenSize.Width, s.manager.config.ScreenSize.Height)
	overlay.Fill(color.RGBA{0, 0, 0, alpha})
	screen.DrawImage(overlay, &ebiten.DrawImageOptions{})
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

// drawMenuOptions renders all menu options with animations
func (s *PausedScene) drawMenuOptions(screen *ebiten.Image) {
	menuY := float64(s.manager.config.ScreenSize.Height) / 2

	for i, option := range s.options {
		y := menuY + float64(i*menuSpacing)
		isSelected := i == s.selection

		s.drawMenuOption(screen, option, y, isSelected)
	}
}

// drawMenuOption renders a single menu option
func (s *PausedScene) drawMenuOption(screen *ebiten.Image, option string, y float64, isSelected bool) {
	_, height := text.Measure(option, defaultFontFace, 0)
	textTopY := y - height/2

	alpha := s.fadeIn
	scale := 1.0

	if isSelected {
		alpha *= 0.8 + 0.2*math.Sin(s.animationTime*pulseSpeed)
		scale = 1.0 + 0.05*math.Sin(s.animationTime*pulseSpeed)

		s.drawSelectionBackground(screen, option, textTopY, 0.6*s.fadeIn)
		s.drawChevron(screen, textTopY)
	} else {
		alpha *= dimmedAlpha
	}

	s.drawOptionText(screen, option, textTopY, alpha, scale)
}

// drawChevron renders the animated selection chevron
func (s *PausedScene) drawChevron(screen *ebiten.Image, y float64) {
	chevronPulse := 0.7 + 0.3*math.Sin(s.animationTime*4.0)
	chevronX := float64(s.manager.config.ScreenSize.Width)/2 - 140 + 10*math.Sin(s.animationTime*2.0)

	op := &text.DrawOptions{}
	op.GeoM.Translate(chevronX, y)
	op.ColorScale.SetR(0)
	op.ColorScale.SetG(1)
	op.ColorScale.SetB(1)
	op.ColorScale.SetA(float32(chevronPulse * s.fadeIn))

	text.Draw(screen, ">", defaultFontFace, op)
}

// drawSelectionBackground renders the background highlight for selected items
func (s *PausedScene) drawSelectionBackground(screen *ebiten.Image, option string, y, alpha float64) {
	width, height := text.Measure(option, defaultFontFace, 0)
	w := int(width)
	h := int(height)

	// Create gradient effect
	rectCol := color.RGBA{0, 255, 255, uint8(128 * alpha)}
	rect := ebiten.NewImage(w+paddingX*2, h+paddingY*2)
	rect.Fill(rectCol)

	// Add subtle border glow
	borderCol := color.RGBA{255, 255, 255, uint8(64 * alpha)}
	border := ebiten.NewImage(w+paddingX*2+4, h+paddingY*2+4)
	border.Fill(borderCol)

	// Draw border first
	borderOp := &ebiten.DrawImageOptions{}
	borderOp.GeoM.Translate(
		float64(s.manager.config.ScreenSize.Width)/2-float64(w+paddingX*2+4)/2,
		y-float64(h+paddingY*2+4)/2,
	)
	screen.DrawImage(border, borderOp)

	// Draw main background
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(
		float64(s.manager.config.ScreenSize.Width)/2-float64(w+paddingX*2)/2,
		y-float64(h+paddingY*2)/2,
	)
	screen.DrawImage(rect, op)
}

// drawOptionText renders the text for a menu option
func (s *PausedScene) drawOptionText(screen *ebiten.Image, option string, y, alpha, scale float64) {
	op := &text.DrawOptions{}

	// Apply scaling from center
	width, height := text.Measure(option, defaultFontFace, 0)
	centerX := float64(s.manager.config.ScreenSize.Width) / 2
	centerY := y

	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(
		centerX-width*scale/2,
		centerY-height*scale/2,
	)

	// Color with fade-in
	op.ColorScale.SetR(1.0)
	op.ColorScale.SetG(1.0)
	op.ColorScale.SetB(1.0)
	op.ColorScale.SetA(float32(alpha))

	text.Draw(screen, option, defaultFontFace, op)
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
}

// Exit is called when the scene becomes inactive
func (s *PausedScene) Exit() {
	s.manager.logger.Debug("Exiting paused scene")
}

// GetType returns the scene type identifier
func (s *PausedScene) GetType() SceneType {
	return ScenePaused
}
