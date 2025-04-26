package player

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// TestSprite implements Drawable for testing
type TestSprite struct {
	bounds image.Rectangle
	img    *ebiten.Image
}

// NewTestSprite creates a new test sprite
func NewTestSprite(width, height int) *TestSprite {
	return &TestSprite{
		bounds: image.Rect(0, 0, width, height),
		img:    ebiten.NewImage(width, height),
	}
}

// Bounds returns the bounds of the sprite
func (m *TestSprite) Bounds() image.Rectangle {
	return m.bounds
}

// Draw implements the Drawable interface
func (m *TestSprite) Draw(screen, op any) {
	// No-op for testing
}

// Image returns the underlying ebiten.Image
func (m *TestSprite) Image() *ebiten.Image {
	return m.img
}
