package pause

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

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
		float64(s.manager.GetConfig().ScreenSize.Width)/2-75, // Adjust for scaling
		titleY,
	)
	op.ColorScale.SetR(0.2)
	op.ColorScale.SetG(0.8)
	op.ColorScale.SetB(1.0)
	op.ColorScale.SetA(float32(titleAlpha))

	text.Draw(screen, "PAUSED", s.font, op)
}

// drawHintText renders the hint text at the bottom of the screen
func (s *PausedScene) drawHintText(screen *ebiten.Image) {
	hintText := "Press ESC to resume quickly"
	hintAlpha := hintBaseAlpha * s.fadeIn * (0.8 + 0.2*math.Sin(s.animationTime*1.5))

	op := &text.DrawOptions{}
	width, _ := text.Measure(hintText, s.font, 0)
	op.GeoM.Translate(
		float64(s.manager.GetConfig().ScreenSize.Width)/2-width/2,
		float64(s.manager.GetConfig().ScreenSize.Height)-hintTextY,
	)
	op.ColorScale.SetR(0.8)
	op.ColorScale.SetG(0.8)
	op.ColorScale.SetB(0.8)
	op.ColorScale.SetA(float32(hintAlpha))

	text.Draw(screen, hintText, s.font, op)
}
