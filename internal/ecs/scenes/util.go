package scenes

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// drawCenteredText draws text centered on screen (helper method for scenes)
func drawCenteredText(screen *ebiten.Image, textStr string, x, y, alpha float64, font text.Face) {
	width, height := text.Measure(textStr, font, 0)
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(int(x)-int(width)/2), float64(int(y)+int(height)/2))
	op.ColorScale.SetR(1)
	op.ColorScale.SetG(1)
	op.ColorScale.SetB(1)
	op.ColorScale.SetA(float32(alpha))
	text.Draw(screen, textStr, font, op)
}
