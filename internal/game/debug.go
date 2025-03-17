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
func (g *GimlarGame) DebugPrintStar(screen *ebiten.Image, star Star) {
	if g.config.Debug {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("Star: X=%.2f, Y=%.2f, Size=%.2f", star.X, star.Y, star.Size))
	}
}

// DrawDebugGridOverlay draws a grid overlay for debugging purposes.
func (g *GimlarGame) DrawDebugGridOverlay(screen *ebiten.Image) {
	if g.config.Debug {
		for i := 0; i < g.config.ScreenWidth; i += debugGridSpacing {
			vector.StrokeLine(screen, float32(i), 0, float32(i), float32(g.config.ScreenHeight), 1, color.White, false)
		}
		for i := 0; i < g.config.ScreenHeight; i += debugGridSpacing {
			vector.StrokeLine(screen, 0, float32(i), float32(g.config.ScreenWidth), float32(i), 1, color.White, false)
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
	for i := 0; i < g.config.ScreenWidth; i += debugGridSpacing {
		vector.StrokeLine(screen, float32(i), 0, float32(i), float32(g.config.ScreenHeight), 1, color.White, false)
	}
	for i := 0; i < g.config.ScreenHeight; i += debugGridSpacing {
		vector.StrokeLine(screen, 0, float32(i), float32(g.config.ScreenWidth), float32(i), 1, color.White, false)
	}
}
