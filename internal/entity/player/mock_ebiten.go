package player

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// MockImage implements ebiten.Image for testing
type MockImage struct {
	bounds image.Rectangle
}

// NewMockImage creates a new mock image with the given dimensions
func NewMockImage(width, height int) *ebiten.Image {
	img := &MockImage{
		bounds: image.Rect(0, 0, width, height),
	}
	return ebiten.NewImageFromImage(img)
}

// Bounds returns the bounds of the mock image
func (m *MockImage) Bounds() image.Rectangle {
	return m.bounds
}

// At returns a black pixel for testing
func (m *MockImage) At(x, y int) color.Color {
	return color.Black
}

// ColorModel returns the RGBA color model
func (m *MockImage) ColorModel() color.Model {
	return color.RGBA64Model
}

// SubImage returns the mock image itself
func (m *MockImage) SubImage(r image.Rectangle) image.Image {
	return m
}
