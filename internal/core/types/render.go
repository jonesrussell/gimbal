package types

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"go.uber.org/zap"
)

// Renderer handles all drawing operations
type Renderer interface {
	Draw(screen *ebiten.Image) error
	AddSprite(id string, sprite *ebiten.Image, pos image.Point, z int) error
	Clear()
}

// RenderConfig holds renderer-specific configuration
type RenderConfig struct {
	MaxSprites   int
	VSync        bool
	DoubleBuffer bool
	Debug        bool
}

// RenderOption configures a Renderer instance
type RenderOption func(*RenderConfig)

// NewRenderer creates a new renderer
type NewRenderer func(logger *zap.Logger, opts ...RenderOption) (Renderer, error)
