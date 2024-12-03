package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/jonesrussell/gimbal/internal/config"
	"github.com/jonesrussell/gimbal/internal/engine"
)

const (
	debugGridSpacing = 32
)

// DrawDebugGrid draws a grid overlay for debugging purposes.
func (g *GimlarGame) DrawDebugGrid(screen *ebiten.Image) error {
	if engine.Debug {
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
