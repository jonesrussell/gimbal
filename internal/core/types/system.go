package types

import (
	"context"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// System represents a game subsystem (physics, AI, etc.)
type System interface {
	Init(ctx context.Context) error
	Update(ctx context.Context) error
	Cleanup(ctx context.Context) error
}

// Entity represents a game object
type Entity interface {
	ID() string
	Update(ctx context.Context) error
	Draw(screen *ebiten.Image) error
	Bounds() image.Rectangle
}
