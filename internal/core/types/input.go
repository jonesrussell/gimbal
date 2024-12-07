package types

import (
	"context"

	"github.com/hajimehoshi/ebiten/v2"
	"go.uber.org/zap"
)

// InputHandler manages all input processing
type InputHandler interface {
	Process(ctx context.Context) error
	IsKeyPressed(key ebiten.Key) bool
	RegisterBinding(name string, key ebiten.Key, handler func(ctx context.Context) error) error
}

// InputConfig holds input-specific configuration
type InputConfig struct {
	EnableKeyboard bool
	EnableMouse    bool
	EnableGamepad  bool
}

// InputOption configures an InputHandler instance
type InputOption func(*InputConfig)

// NewInputHandler creates a new input handler
type NewInputHandler func(logger *zap.Logger) (InputHandler, error)
