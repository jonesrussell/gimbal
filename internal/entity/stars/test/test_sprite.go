package test

import (
	"image"
)

// Sprite implements SpriteImage for testing
type Sprite struct {
	bounds image.Rectangle
}

// NewSprite creates a new test sprite
func NewSprite(width, height int) *Sprite {
	return &Sprite{
		bounds: image.Rect(0, 0, width, height),
	}
}

// Bounds returns the bounds of the sprite
func (s *Sprite) Bounds() image.Rectangle {
	return s.bounds
}

// Draw implements the DrawableSprite interface
func (s *Sprite) Draw(screen any, op any) {
	// No-op for testing
}
