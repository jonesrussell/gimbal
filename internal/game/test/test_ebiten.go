package test

import (
	"image"
	"os"
	"testing"

	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// TestMain handles Ebiten initialization for tests
func TestMain(m *testing.M) {
	// Set up test environment
	os.Setenv("EBITEN_TEST", "1")

	// Run tests
	code := m.Run()

	// Clean up
	os.Unsetenv("EBITEN_TEST")
	os.Exit(code)
}

// MockImage implements a mock ebiten.Image for testing
type MockImage struct {
	width  int
	height int
}

// NewMockImage creates a new mock image
func NewMockImage(width, height int) *MockImage {
	return &MockImage{
		width:  width,
		height: height,
	}
}

// Bounds returns the bounds of the mock image
func (m *MockImage) Bounds() image.Rectangle {
	return image.Rect(0, 0, m.width, m.height)
}

// ColorModel returns the color model of the mock image
func (m *MockImage) ColorModel() color.Model {
	return color.RGBAModel
}

// DrawImage implements ebiten.Image interface
func (m *MockImage) DrawImage(image *ebiten.Image, options *ebiten.DrawImageOptions) {
	// No-op for testing
}

// Fill implements ebiten.Image interface
func (m *MockImage) Fill(color color.Color) {
	// No-op for testing
}

// Clear implements ebiten.Image interface
func (m *MockImage) Clear() {
	// No-op for testing
}

// ReplacePixels implements ebiten.Image interface
func (m *MockImage) ReplacePixels(pixels []byte) {
	// No-op for testing
}

// At implements ebiten.Image interface
func (m *MockImage) At(x, y int) color.Color {
	return color.White
}

// Set implements ebiten.Image interface
func (m *MockImage) Set(x, y int, clr color.Color) {
	// No-op for testing
}

// SubImage implements ebiten.Image interface
func (m *MockImage) SubImage(r image.Rectangle) image.Image {
	return m
}

// Dispose implements ebiten.Image interface
func (m *MockImage) Dispose() {
	// No-op for testing
}
