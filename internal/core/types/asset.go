package types

import (
	"context"

	"github.com/hajimehoshi/ebiten/v2"
	"go.uber.org/zap"
)

// AssetManager handles loading and managing game assets
type AssetManager interface {
	LoadImage(ctx context.Context, path string) (*ebiten.Image, error)
	LoadSound(ctx context.Context, path string) ([]byte, error)
	Preload(ctx context.Context, dir string) error
	Cleanup(ctx context.Context) error
}

// AssetManagerConfig holds asset manager configuration
type AssetManagerConfig struct {
	BaseDir     string
	CacheSize   int
	EnableSound bool
	Debug       bool
}

// AssetOption configures the asset manager
type AssetOption func(*AssetManagerConfig)

// NewAssetManager creates a new asset manager
type NewAssetManager func(logger *zap.Logger, opts ...AssetOption) (AssetManager, error)
