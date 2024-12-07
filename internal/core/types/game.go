package types

import (
	"context"

	"github.com/hajimehoshi/ebiten/v2"
	"go.uber.org/zap"
)

// Game represents the main game interface
type Game interface {
	Update(ctx context.Context) error
	Draw(screen *ebiten.Image) error
	Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int)
	Shutdown(ctx context.Context) error
}

// GameConfig holds game-specific configuration
type GameConfig struct {
	ScreenWidth  int  // 8 bytes
	ScreenHeight int  // 8 bytes
	Debug        bool // 1 byte, padded to 8 bytes
}

// GameOption configures a Game instance
type GameOption func(*GameConfig)

// NewGame creates a new game instance
type NewGame func(logger *zap.Logger) (Game, error)
