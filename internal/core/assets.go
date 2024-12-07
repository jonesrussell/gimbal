// Package core provides the core game engine functionality including asset management,
// rendering, input handling, and game loop management.
package core

import (
	"context"
	"fmt"
	"image"
	"os"
	"path/filepath"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"go.uber.org/zap"
)

// Implementation-specific types
type assetManagerImpl struct {
	// 8-byte aligned fields
	images  map[string]*ebiten.Image // 8 bytes
	sounds  map[string][]byte        // 8 bytes
	logger  *zap.Logger              // 8 bytes
	config  *AssetManagerConfig      // 8 bytes
	baseDir string                   // 16 bytes
	mu      sync.RWMutex             // 8 bytes
}

// Memory layout visualization:
// |-----------------------------------------------|
// | mu (8 bytes)                                  |
// |-----------------------------------------------|
// | logger (8) | config (8)                       |
// |-----------------------------------------------|
// | images (8) | sounds (8)                       |
// |-----------------------------------------------|
// | baseDir (16 bytes)                            |
// |-----------------------------------------------|

// Verify that assetManagerImpl implements AssetManager interface
var _ NewAssetManager = NewAssetManagerImpl

func NewAssetManagerImpl(logger *zap.Logger, opts ...AssetOption) (AssetManager, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger is required")
	}

	config := &AssetManagerConfig{
		BaseDir:     "assets",
		CacheSize:   1000,
		EnableSound: true,
		Debug:       false,
	}

	for _, opt := range opts {
		opt(config)
	}

	am := &assetManagerImpl{
		logger:  logger,
		baseDir: config.BaseDir,
		config:  config,
		images:  make(map[string]*ebiten.Image, config.CacheSize),
		sounds:  make(map[string][]byte, config.CacheSize),
	}

	return am, nil
}

func (am *assetManagerImpl) LoadImage(ctx context.Context, path string) (*ebiten.Image, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		am.mu.RLock()
		if img, ok := am.images[path]; ok {
			am.mu.RUnlock()
			return img, nil
		}
		am.mu.RUnlock()

		// Load new image
		am.mu.Lock()
		defer am.mu.Unlock()

		// Double check after acquiring write lock
		if img, ok := am.images[path]; ok {
			return img, nil
		}

		fullPath := filepath.Join(am.baseDir, path)
		file, err := os.Open(fullPath)
		if err != nil {
			return nil, fmt.Errorf("failed to open image %s: %w", path, err)
		}
		defer file.Close()

		img, _, err := image.Decode(file)
		if err != nil {
			return nil, fmt.Errorf("failed to decode image %s: %w", path, err)
		}

		ebitenImg := ebiten.NewImageFromImage(img)
		am.images[path] = ebitenImg

		am.logger.Debug("loaded image",
			zap.String("path", path),
			zap.Int("width", ebitenImg.Bounds().Dx()),
			zap.Int("height", ebitenImg.Bounds().Dy()))

		return ebitenImg, nil
	}
}

func (am *assetManagerImpl) LoadSound(ctx context.Context, path string) ([]byte, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		am.mu.RLock()
		if sound, ok := am.sounds[path]; ok {
			am.mu.RUnlock()
			return sound, nil
		}
		am.mu.RUnlock()

		// Load new sound
		am.mu.Lock()
		defer am.mu.Unlock()

		// Double check after acquiring write lock
		if sound, ok := am.sounds[path]; ok {
			return sound, nil
		}

		fullPath := filepath.Join(am.baseDir, path)
		sound, err := os.ReadFile(fullPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load sound %s: %w", path, err)
		}

		am.sounds[path] = sound
		am.logger.Debug("loaded sound",
			zap.String("path", path),
			zap.Int("size", len(sound)))

		return sound, nil
	}
}

func (am *assetManagerImpl) Preload(ctx context.Context, dir string) error {
	fullDir := filepath.Join(am.baseDir, dir)
	return filepath.Walk(fullDir, am.handlePreloadFile(ctx))
}

func (am *assetManagerImpl) handlePreloadFile(ctx context.Context) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed to access path %s: %w", path, err)
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if info.IsDir() {
				return nil
			}
			return am.preloadAsset(ctx, path)
		}
	}
}

func (am *assetManagerImpl) preloadAsset(ctx context.Context, path string) error {
	relPath, err := filepath.Rel(am.baseDir, path)
	if err != nil {
		return fmt.Errorf("failed to get relative path for %s: %w", path, err)
	}

	ext := filepath.Ext(path)
	var loadErr error

	switch ext {
	case ".png", ".jpg", ".jpeg":
		_, loadErr = am.LoadImage(ctx, relPath)
	case ".wav", ".mp3":
		_, loadErr = am.LoadSound(ctx, relPath)
	default:
		return nil // Skip unsupported file types
	}

	if loadErr != nil {
		am.logger.Error("failed to preload asset",
			zap.String("path", relPath),
			zap.Error(loadErr))
		return fmt.Errorf("failed to preload %s: %w", relPath, loadErr)
	}

	return nil
}

func (am *assetManagerImpl) Cleanup(ctx context.Context) error {
	am.mu.Lock()
	defer am.mu.Unlock()

	for _, img := range am.images {
		if img != nil {
			img.Deallocate()
		}
	}
	am.images = nil
	am.sounds = nil
	return nil
}

// Asset manager options
func WithBaseDir(dir string) AssetOption {
	return func(c *AssetManagerConfig) {
		c.BaseDir = dir
	}
}

func WithCacheSize(size int) AssetOption {
	return func(c *AssetManagerConfig) {
		c.CacheSize = size
	}
}

func WithSound(enable bool) AssetOption {
	return func(c *AssetManagerConfig) {
		c.EnableSound = enable
	}
}

func WithAssetDebug(debug bool) AssetOption {
	return func(c *AssetManagerConfig) {
		c.Debug = debug
	}
}
