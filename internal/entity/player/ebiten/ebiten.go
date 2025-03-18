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
	return &Sprite{img: img}
}

// Bounds returns the bounds of the sprite
func (e *Sprite) Bounds() image.Rectangle {
	return e.img.Bounds()
}

// Draw implements the DrawableSprite interface
func (e *Sprite) Draw(screen any, op any) {
	if ebitenScreen, ok := screen.(*ebiten.Image); ok {
		// Create default options if none provided
		drawOp := &ebiten.DrawImageOptions{}
		if ebitenOp, okOp := op.(*ebiten.DrawImageOptions); okOp {
			drawOp = ebitenOp
		}

		// Calculate scale to maintain 32x32 size
		bounds := e.img.Bounds()
		scaleX := float64(TargetSize) / float64(bounds.Dx())
		scaleY := float64(TargetSize) / float64(bounds.Dy())

		// Order of transformations:
		// 1. Center the sprite (move origin to center)
		drawOp.GeoM.Translate(-float64(bounds.Dx())/HalfDivisor, -float64(bounds.Dy())/HalfDivisor)

		// 2. Scale to target size
		drawOp.GeoM.Scale(scaleX, scaleY)

		// The rotation and final position are handled by the player's Draw method

		ebitenScreen.DrawImage(e.img, drawOp)
	}
}

// Image returns the underlying ebiten.Image
func (e *Sprite) Image() *ebiten.Image {
	return e.img
}
