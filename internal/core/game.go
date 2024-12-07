package core

import (
	"context"
	"fmt"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"go.uber.org/zap"
)

type gameImpl struct {
	systems  []System     // slice: 24 bytes
	input    InputHandler // interface: 16 bytes
	renderer Renderer     // interface: 16 bytes
	assets   AssetManager // interface: 16 bytes
	mu       sync.RWMutex // mutex: 8 bytes
	logger   *zap.Logger  // ptr: 8 bytes
	config   *GameConfig  // ptr: 8 bytes
}

// NewGameImpl creates a new game instance with named return values
func NewGameImpl(logger *zap.Logger, opts ...GameOption) (game Game, err error) {
	if logger == nil {
		return nil, fmt.Errorf("logger is required")
	}

	config := &GameConfig{
		ScreenWidth:  800,
		ScreenHeight: 600,
		Debug:        false,
	}

	for _, opt := range opts {
		opt(config)
	}

	return &gameImpl{
		logger: logger,
		config: config,
	}, nil
}

func (g *gameImpl) Update(ctx context.Context) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		// Update input
		if err := g.input.Process(ctx); err != nil {
			return fmt.Errorf("input processing error: %w", err)
		}

		// Update systems
		for _, sys := range g.systems {
			if err := sys.Update(ctx); err != nil {
				return fmt.Errorf("system update error: %w", err)
			}
		}

		return nil
	}
}

func (g *gameImpl) Draw(screen *ebiten.Image) error {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return g.renderer.Draw(screen)
}

func (g *gameImpl) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.config.ScreenWidth, g.config.ScreenHeight
}

func (g *gameImpl) Shutdown(ctx context.Context) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.logger.Info("shutting down game")

	// Cleanup systems in reverse order
	for i := len(g.systems) - 1; i >= 0; i-- {
		if err := g.systems[i].Cleanup(ctx); err != nil {
			g.logger.Error("system cleanup error", zap.Error(err))
		}
	}

	// Cleanup assets
	if err := g.assets.Cleanup(ctx); err != nil {
		g.logger.Error("asset cleanup error", zap.Error(err))
	}

	return nil
}

// Game options
func WithScreenSize(width, height int) GameOption {
	return func(c *GameConfig) {
		c.ScreenWidth = width
		c.ScreenHeight = height
	}
}

func WithDebugMode(debug bool) GameOption {
	return func(c *GameConfig) {
		c.Debug = debug
	}
}
