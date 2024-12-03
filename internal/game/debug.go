package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/jonesrussell/gimbal/internal/config"
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
func DrawDebugGridOverlay(screen *ebiten.Image) error {
	if Debug {
		cfg, err := config.New()
		if err != nil {
			return fmt.Errorf("failed to get config for debug grid: %w", err)
		}

		for i := 0; i < cfg.Screen.Width; i += debugGridSpacing {
			vector.StrokeLine(screen, float32(i), 0, float32(i), float32(cfg.Screen.Height), 1, color.White, false)
		}
		for i := 0; i < cfg.Screen.Height; i += debugGridSpacing {
			vector.StrokeLine(screen, 0, float32(i), float32(cfg.Screen.Width), float32(i), 1, color.White, false)
		}
	}
	return nil
}

func (g *GimlarGame) DrawDebugInfo(screen *ebiten.Image) {
	// Print the current FPS
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", ebiten.ActualFPS()))

	// Draw grid overlay
	g.DrawDebugGrid(screen)
}

func (g *GimlarGame) DrawDebugGrid(screen *ebiten.Image) error {
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("failed to get config for debug grid: %w", err)
	}

	// Draw grid overlay
	for i := 0; i < cfg.Screen.Width; i += debugGridSpacing {
		vector.StrokeLine(screen, float32(i), 0, float32(i), float32(cfg.Screen.Height), 1, color.White, false)
	}
	for i := 0; i < cfg.Screen.Height; i += debugGridSpacing {
		vector.StrokeLine(screen, 0, float32(i), float32(cfg.Screen.Width), float32(i), 1, color.White, false)
	}
	return nil
}
