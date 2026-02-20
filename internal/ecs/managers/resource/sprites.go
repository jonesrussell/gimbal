package resources

import (
	"bytes"
	"context"
	"image"
	"image/draw"
	"image/png"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/jonesrussell/gimbal/assets"
	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/errors"
)

// LoadSprite loads a sprite from embedded assets with simplified logic
func (rm *ResourceManager) LoadSprite(ctx context.Context, name, path string) (*ebiten.Image, error) {
	// Check context cancellation
	if err := common.CheckContextCancellation(ctx); err != nil {
		return nil, err
	}

	// Try cache first
	if cached := rm.getCachedSprite(name); cached != nil {
		return cached, nil
	}

	// Load and decode sprite
	return rm.loadAndCacheSprite(ctx, name, path)
}

// getCachedSprite retrieves sprite from cache if exists
func (rm *ResourceManager) getCachedSprite(name string) *ebiten.Image {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	if resource, exists := rm.resources[name]; exists {
		if sprite, ok := resource.Data.(*ebiten.Image); ok {
			return sprite
		}
	}
	return nil
}

// loadAndCacheSprite loads sprite from assets and caches it
func (rm *ResourceManager) loadAndCacheSprite(ctx context.Context, name, path string) (*ebiten.Image, error) {
	// Load from embedded assets
	imageData, err := rm.loadSpriteFile(path)
	if err != nil {
		return nil, err
	}

	// Decode PNG data
	sprite, err := rm.decodePNGData(name, path, imageData)
	if err != nil {
		return nil, err
	}

	// Cache the sprite
	rm.cacheSprite(name, sprite)

	return sprite, nil
}

// loadSpriteFile loads sprite file from embedded assets
func (rm *ResourceManager) loadSpriteFile(path string) ([]byte, error) {
	imageData, err := assets.Assets.ReadFile(path)
	if err != nil {
		return nil, errors.NewGameErrorWithCause(errors.AssetLoadFailed, "failed to read sprite file", err)
	}

	return imageData, nil
}

// decodePNGData decodes PNG data into ebiten image
func (rm *ResourceManager) decodePNGData(name, path string, imageData []byte) (*ebiten.Image, error) {
	// Use DecodeConfig first to check image format, then decode
	// This ensures we handle transparency correctly
	img, err := png.Decode(bytes.NewReader(imageData))
	if err != nil {
		return nil, errors.NewGameErrorWithCause(errors.AssetInvalid, "failed to decode sprite", err)
	}

	// Convert to NRGBA if not already to ensure alpha channel is preserved
	finalImg := img
	if _, ok := img.(*image.NRGBA); !ok {
		// Convert to NRGBA to ensure transparency works correctly
		// Use draw.Draw to properly preserve alpha channel
		bounds := img.Bounds()
		nrgba := image.NewNRGBA(bounds)
		draw.Draw(nrgba, bounds, img, bounds.Min, draw.Src)
		finalImg = nrgba
	}

	sprite := ebiten.NewImageFromImage(finalImg)
	return sprite, nil
}

// cacheSprite stores sprite in the resource cache
func (rm *ResourceManager) cacheSprite(name string, sprite *ebiten.Image) {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	if rm.resources == nil {
		rm.resources = make(map[string]*Resource)
	}

	rm.resources[name] = &Resource{
		Name: name,
		Type: ResourceSprite,
		Data: sprite,
	}
}

// GetSprite retrieves a loaded sprite
func (rm *ResourceManager) GetSprite(ctx context.Context, name string) (*ebiten.Image, bool) {
	// Check for cancellation
	if err := common.CheckContextCancellation(ctx); err != nil {
		return nil, false
	}

	rm.mutex.RLock()
	defer rm.mutex.RUnlock()

	if resource, exists := rm.resources[name]; exists {
		if sprite, ok := resource.Data.(*ebiten.Image); ok {
			return sprite, true
		}
	}
	return nil, false
}
