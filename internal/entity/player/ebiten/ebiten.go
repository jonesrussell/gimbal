package ebiten

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	// HalfDivisor is used to calculate half of a dimension
	HalfDivisor = 2
	// TargetSize is the desired size of the sprite in pixels
	TargetSize = 32
)

// Sprite wraps an ebiten.Image to implement SpriteImage
type Sprite struct {
	img *ebiten.Image
}

// NewSprite creates a new Sprite from an ebiten.Image
func NewSprite(img *ebiten.Image) *Sprite {
	return &Sprite{
		img: img,
	}
}

// Bounds returns the bounds of the sprite
func (e *Sprite) Bounds() image.Rectangle {
	return e.img.Bounds()
}

// Draw implements the DrawableSprite interface
func (e *Sprite) Draw(screen any, op any) {
	if ebitenScreen, ok := screen.(*ebiten.Image); ok {
		// Use provided options or create new ones
		var drawOp *ebiten.DrawImageOptions
		if ebitenOp, okOp := op.(*ebiten.DrawImageOptions); okOp {
			drawOp = ebitenOp
		} else {
			drawOp = &ebiten.DrawImageOptions{}
		}

		// Calculate scale to maintain 32x32 size
		bounds := e.img.Bounds()
		scaleX := float64(TargetSize) / float64(bounds.Dx())
		scaleY := float64(TargetSize) / float64(bounds.Dy())

		// Apply scale transformation first
		drawOp.GeoM.Scale(scaleX, scaleY)

		// Add a color tint to make the sprite more visible for debugging
		drawOp.ColorM.Scale(2.0, 2.0, 2.0, 1.0) // Increased brightness for better visibility

		// Draw the sprite
		ebitenScreen.DrawImage(e.img, drawOp)
	}
}

// Image returns the underlying ebiten.Image
func (e *Sprite) Image() *ebiten.Image {
	return e.img
}
