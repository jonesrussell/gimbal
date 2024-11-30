package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	debugGridSpacing = 32
)

// DebugPrintStar prints the debug information for a star.
func DebugPrintStar(star Star) {
	if Debug {
		fmt.Printf("Star: X=%.2f, Y=%.2f, Size=%.2f\n", star.X, star.Y, star.Size)
	}
}

// DrawDebugGridOverlay draws a grid overlay for debugging purposes.
func DrawDebugGridOverlay(screen *ebiten.Image) {
	if Debug {
		for i := 0; i < screenWidth; i += debugGridSpacing {
			vector.StrokeLine(screen, float32(i), 0, float32(i), float32(screenHeight), 1, color.White, false)
		}
		for i := 0; i < screenHeight; i += debugGridSpacing {
			vector.StrokeLine(screen, 0, float32(i), float32(screenWidth), float32(i), 1, color.White, false)
		}
	}
}

func (g *GimlarGame) DrawDebugInfo(screen *ebiten.Image) {
	// Print the current FPS
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", ebiten.ActualFPS()))

	// Draw grid overlay
	g.DrawDebugGrid(screen)
}

func (g *GimlarGame) DrawDebugGrid(screen *ebiten.Image) {
	// Draw grid overlay
	for i := 0; i < screenWidth; i += debugGridSpacing {
		vector.StrokeLine(screen, float32(i), 0, float32(i), float32(screenHeight), 1, color.White, false)
	}
	for i := 0; i < screenHeight; i += debugGridSpacing {
		vector.StrokeLine(screen, 0, float32(i), float32(screenWidth), float32(i), 1, color.White, false)
	}
}
