package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/jonesrussell/gimbal/internal/entity/stars"
)

const (
	debugGridSpacing = 32
)

// DebugPrintStar prints the debug information for a star.
func (g *GimlarGame) DebugPrintStar(screen *ebiten.Image, star *stars.Star) {
	if g.config.Debug {
		pos := star.GetPosition()
		size := star.GetSize()
		ebitenutil.DebugPrint(screen, fmt.Sprintf("Star: X=%.2f, Y=%.2f, Size=%.2f", pos.X, pos.Y, size))
	}
}

// DrawDebugGridOverlay draws a grid overlay for debugging purposes.
func (g *GimlarGame) DrawDebugGridOverlay(screen *ebiten.Image) {
	if g.config.Debug {
		for i := 0; i < g.config.ScreenSize.Width; i += debugGridSpacing {
			vector.StrokeLine(screen, float32(i), 0, float32(i), float32(g.config.ScreenSize.Height), 1, color.White, false)
		}
		for i := 0; i < g.config.ScreenSize.Height; i += debugGridSpacing {
			vector.StrokeLine(screen, 0, float32(i), float32(g.config.ScreenSize.Width), float32(i), 1, color.White, false)
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
	for i := 0; i < g.config.ScreenSize.Width; i += debugGridSpacing {
		vector.StrokeLine(screen, float32(i), 0, float32(i), float32(g.config.ScreenSize.Height), 1, color.White, false)
	}
	for i := 0; i < g.config.ScreenSize.Height; i += debugGridSpacing {
		vector.StrokeLine(screen, 0, float32(i), float32(g.config.ScreenSize.Width), float32(i), 1, color.White, false)
	}
}
