package ecs

import (
	"image/color"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

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
