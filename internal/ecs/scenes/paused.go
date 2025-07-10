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

type PauseOption int

const (
	PauseOptionResume PauseOption = iota
	PauseOptionReturnToMenu
	PauseOptionQuit
)

type PausedScene struct {
	manager           *SceneManager
	selection         int
	options           []string
	animationTime     float64
	selectionChanged  bool
	lastSelectionTime time.Time
	fadeIn            float64

	// Animation constants
	fadeInDuration  float64
	pulseSpeed      float64
	chevronOffset   float64
	hoverTransition float64
}

func NewPausedScene(manager *SceneManager) *PausedScene {
	return &PausedScene{
		manager:           manager,
		selection:         0,
		options:           []string{"Resume", "Return to Menu", "Quit"},
		animationTime:     0,
		selectionChanged:  false,
		lastSelectionTime: time.Now(),
		fadeIn:            0,
		fadeInDuration:    0.3,
		pulseSpeed:        3.0,
		chevronOffset:     0,
		hoverTransition:   0,
	}
}

func (s *PausedScene) Update() error {
	dt := 1.0 / 60.0 // Assume 60 FPS
	s.animationTime += dt

	// Handle fade-in animation
	if s.fadeIn < 1.0 {
		s.fadeIn = math.Min(1.0, s.fadeIn+dt/s.fadeInDuration)
	}

	// Handle selection change animation
	if s.selectionChanged && time.Since(s.lastSelectionTime).Seconds() > 0.1 {
		s.selectionChanged = false
	}

	// Keyboard navigation
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		s.changeSelection((s.selection - 1 + len(s.options)) % len(s.options))
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		s.changeSelection((s.selection + 1) % len(s.options))
	}

	// Mouse interaction
	s.handleMouseInput()

	// Keyboard activation
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		s.activateSelection()
	}

	// Quick resume with ESC
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		s.manager.SwitchScene(ScenePlaying)
	}

	return nil
}

func (s *PausedScene) changeSelection(newSelection int) {
	if newSelection != s.selection {
		s.selection = newSelection
		s.selectionChanged = true
		s.lastSelectionTime = time.Now()
	}
}

func (s *PausedScene) handleMouseInput() {
	x, y := ebiten.CursorPosition()
	menuY := float64(s.manager.config.ScreenSize.Height) / 2
	hoveredItem := -1

	for i := range s.options {
		itemY := menuY + float64(i*50) // Increased spacing
		width, height := text.Measure(s.options[i], defaultFontFace, 0)
		w := int(width)
		h := int(height)

		// Create larger hit area for better UX
		padding := 50
		itemRect := struct{ x0, y0, x1, y1 int }{
			int(float64(s.manager.config.ScreenSize.Width)/2) - w/2 - padding,
			int(itemY) - h/2 - 15,
			int(float64(s.manager.config.ScreenSize.Width)/2) + w/2 + padding,
			int(itemY) + h/2 + 15,
		}

		if x >= itemRect.x0 && x <= itemRect.x1 && y >= itemRect.y0 && y <= itemRect.y1 {
			hoveredItem = i
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				s.activateSelection()
			}
		}
	}

	// Smooth hover transition
	if hoveredItem != -1 && hoveredItem != s.selection {
		s.changeSelection(hoveredItem)
	}
}

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

func (s *PausedScene) Draw(screen *ebiten.Image) {
	// Draw animated overlay with fade-in
	overlayAlpha := uint8(128 * s.fadeIn)
	s.drawOverlay(screen, overlayAlpha)

	// Draw title with fade-in
	s.drawTitle(screen)

	// Draw menu options
	s.drawMenuOptions(screen)

	// Draw hint text
	s.drawHintText(screen)
}

func (s *PausedScene) drawOverlay(screen *ebiten.Image, alpha uint8) {
	overlay := ebiten.NewImage(s.manager.config.ScreenSize.Width, s.manager.config.ScreenSize.Height)
	overlay.Fill(color.RGBA{0, 0, 0, alpha})
	screen.DrawImage(overlay, &ebiten.DrawImageOptions{})
}

func (s *PausedScene) drawTitle(screen *ebiten.Image) {
	titleAlpha := s.fadeIn
	// Add subtle pulsing effect to title
	pulse := 0.9 + 0.1*math.Sin(s.animationTime*2.0)
	titleAlpha *= pulse

	op := &text.DrawOptions{}
	op.GeoM.Scale(1.5, 1.5) // Larger title
	op.GeoM.Translate(
		float64(s.manager.config.ScreenSize.Width)/2-75, // Adjust for scaling
		80,
	)
	op.ColorScale.SetR(0.2)
	op.ColorScale.SetG(0.8)
	op.ColorScale.SetB(1.0)
	op.ColorScale.SetA(float32(titleAlpha))

	text.Draw(screen, "PAUSED", defaultFontFace, op)
}

func (s *PausedScene) drawMenuOptions(screen *ebiten.Image) {
	menuY := float64(s.manager.config.ScreenSize.Height) / 2

	for i, option := range s.options {
		y := menuY + float64(i*50)
		isSelected := i == s.selection

		// Calculate animations
		alpha := s.fadeIn
		bgAlpha := 0.0
		scale := 1.0

		if isSelected {
			// Smooth pulsing animation
			pulse := 0.8 + 0.2*math.Sin(s.animationTime*s.pulseSpeed)
			alpha *= pulse
			bgAlpha = 0.6 * s.fadeIn

			// Subtle scale effect
			scale = 1.0 + 0.05*math.Sin(s.animationTime*s.pulseSpeed)

			// Animated chevron with smooth movement
			chevronPulse := 0.7 + 0.3*math.Sin(s.animationTime*4.0)
			chevronX := float64(s.manager.config.ScreenSize.Width)/2 - 140 + 10*math.Sin(s.animationTime*2.0)

			chevronOp := &text.DrawOptions{}
			chevronOp.GeoM.Translate(chevronX, y-5)
			chevronOp.ColorScale.SetR(0)
			chevronOp.ColorScale.SetG(1)
			chevronOp.ColorScale.SetB(1)
			chevronOp.ColorScale.SetA(float32(chevronPulse * s.fadeIn))
			text.Draw(screen, ">", defaultFontFace, chevronOp)
		} else {
			alpha *= 0.7 // Dim non-selected items
		}

		// Draw background highlight with rounded corners effect
		if bgAlpha > 0 {
			s.drawSelectionBackground(screen, option, y, bgAlpha)
		}

		// Draw text with scaling
		s.drawOptionText(screen, option, y, alpha, scale)
	}
}

func (s *PausedScene) drawSelectionBackground(screen *ebiten.Image, option string, y, alpha float64) {
	width, height := text.Measure(option, defaultFontFace, 0)
	w := int(width)
	h := int(height)

	paddingX := 30
	paddingY := 10

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

func (s *PausedScene) drawHintText(screen *ebiten.Image) {
	hintText := "Press ESC to resume quickly"
	hintAlpha := 0.6 * s.fadeIn * (0.8 + 0.2*math.Sin(s.animationTime*1.5))

	op := &text.DrawOptions{}
	width, _ := text.Measure(hintText, defaultFontFace, 0)
	op.GeoM.Translate(
		float64(s.manager.config.ScreenSize.Width)/2-width/2,
		float64(s.manager.config.ScreenSize.Height)-60,
	)
	op.ColorScale.SetR(0.8)
	op.ColorScale.SetG(0.8)
	op.ColorScale.SetB(0.8)
	op.ColorScale.SetA(float32(hintAlpha))

	text.Draw(screen, hintText, defaultFontFace, op)
}

func (s *PausedScene) Enter() {
	s.manager.logger.Debug("Entering paused scene")
	s.fadeIn = 0
	s.animationTime = 0
	s.selectionChanged = false
}

func (s *PausedScene) Exit() {
	s.manager.logger.Debug("Exiting paused scene")
}

func (s *PausedScene) GetType() SceneType {
	return ScenePaused
}
