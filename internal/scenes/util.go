package scenes

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// TextDrawOptions holds options for drawing centered text
type TextDrawOptions struct {
	Text  string
	X, Y  float64
	Alpha float64
	Font  text.Face
}

// DrawCenteredTextWithOptions draws text with options struct
func DrawCenteredTextWithOptions(screen *ebiten.Image, opts TextDrawOptions) {
	width, height := text.Measure(opts.Text, opts.Font, 0)
	drawOp := &text.DrawOptions{}
	drawOp.GeoM.Translate(float64(int(opts.X)-int(width)/2), float64(int(opts.Y)+int(height)/2))
	drawOp.ColorScale.SetR(1)
	drawOp.ColorScale.SetG(1)
	drawOp.ColorScale.SetB(1)
	drawOp.ColorScale.SetA(float32(opts.Alpha))
	text.Draw(screen, opts.Text, opts.Font, drawOp)
}
