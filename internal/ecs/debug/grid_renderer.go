package debug

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// drawGrid draws a debug grid overlay
func (dr *DebugRenderer) drawGrid(screen *ebiten.Image) {
	bounds := screen.Bounds()
	gridSize := 50

	// Draw vertical lines - barely visible guide lines
	for x := 0; x < bounds.Dx(); x += gridSize {
		vector.StrokeLine(screen, float32(x), 0, float32(x), float32(bounds.Dy()),
			1, color.RGBA{255, 255, 255, 20}, false)
	}

	// Draw horizontal lines - barely visible guide lines
	for y := 0; y < bounds.Dy(); y += gridSize {
		vector.StrokeLine(screen, 0, float32(y), float32(bounds.Dx()), float32(y),
			1, color.RGBA{255, 255, 255, 20}, false)
	}
}
