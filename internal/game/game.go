package game

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jonesrussell/gimbal/internal/config"
	"github.com/jonesrussell/gimbal/internal/engine"
	"github.com/jonesrussell/gimbal/player"
	"github.com/solarlune/resolv"
	"go.uber.org/zap"
)

type GimlarGame struct {
	logger *zap.Logger    // 8 bytes
	config *config.Config // 8 bytes
	player *player.Player // 8 bytes
	space  *resolv.Space  // 8 bytes
	speed  float64        // 8 bytes
	prevX  float64        // 8 bytes
	prevY  float64        // 8 bytes
}

// GimlarGame should implement engine.GameEngine
var _ engine.GameEngine = (*GimlarGame)(nil)

// NewGimlarGame creates a new game instance
func NewGimlarGame(logger *zap.Logger, cfg *config.Config) (*GimlarGame, error) {
	switch {
	case logger == nil:
		return nil, fmt.Errorf("logger is required")
	case cfg == nil:
		return nil, fmt.Errorf("config is required")
	}

	logger.Info("NewGimlarGame called with config",
		zap.Int("screen_width", cfg.Screen.Width),
		zap.Int("screen_height", cfg.Screen.Height))

	// Set Debug based on environment variable or config
	engine.Debug = cfg.Game.Debug

	g := &GimlarGame{
		player: &player.Player{},
		speed:  cfg.Game.Speed,
		space:  resolv.NewSpace(0, 0, cfg.Screen.Width, cfg.Screen.Height),
		prevX:  0,
		prevY:  0,
		logger: logger,
		config: cfg,
	}

	logger.Info("Game struct initialized successfully")
	return g, nil
}

// Update implements engine.GameEngine
func (g *GimlarGame) Update() error {
	g.logger.Debug("Update frame")
	return nil
}

// Draw implements engine.GameEngine
func (g *GimlarGame) Draw(screen *ebiten.Image) {
	if screen == nil {
		g.logger.Error("Draw called with nil screen")
		return
	}

	g.logger.Debug("Draw frame")

	// Draw the player
	if g.player != nil {
		g.player.Draw(screen)
	}
}

func (g *GimlarGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	g.logger.Debug("Layout called",
		zap.Int("outsideWidth", outsideWidth),
		zap.Int("outsideHeight", outsideHeight),
		zap.Int("returnWidth", g.config.Screen.Width),
		zap.Int("returnHeight", g.config.Screen.Height))

	return g.config.Screen.Width, g.config.Screen.Height
}
