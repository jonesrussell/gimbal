// internal/ecs/managers/resource/interfaces.go
// Interfaces for the existing resource manager structure

package resources

import (
	"context"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// SpriteManager handles sprite operations (matches existing methods)
type SpriteManager interface {
	LoadSprite(ctx context.Context, name, path string) (*ebiten.Image, error)
	GetSprite(ctx context.Context, name string) (*ebiten.Image, error)
	CreateSprite(ctx context.Context, name string, width, height int, color string) (*ebiten.Image, error)
	GetScaledSprite(ctx context.Context, name string, width, height int) (*ebiten.Image, error)
	GetUISprite(ctx context.Context, name string, size int) (*ebiten.Image, error)
}

// FontManager handles font operations (matches existing methods)
type FontManager interface {
	GetDefaultFont() (*text.GoTextFace, error)
}

// ResourceInfo provides resource metadata (matches existing methods)
type ResourceInfo interface {
	GetResourceCount() int
	GetResourceInfo() map[string]interface{}
}

// ResourceCleaner handles cleanup operations (matches existing methods)
type ResourceCleaner interface {
	Cleanup() error
}

// Manager combines all resource management capabilities
// This interface matches your existing ResourceManager struct
type Manager interface {
	SpriteManager
	FontManager
	ResourceInfo
	ResourceCleaner
}

// SpriteLoader handles batch sprite loading operations
// Based on methods I can see in your sprites.go
type SpriteLoader interface {
	LoadAllSprites(ctx context.Context) error
}
