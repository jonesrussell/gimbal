package ebiten

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
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
		if ebitenOp, okOp := op.(*ebiten.DrawImageOptions); okOp {
			ebitenScreen.DrawImage(e.img, ebitenOp)
		}
	}
}

// Image returns the underlying ebiten.Image
func (e *Sprite) Image() *ebiten.Image {
	return e.img
}
