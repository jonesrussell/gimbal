package resources

import (
	"context"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

// GetScaledSprite returns a sprite scaled to the given width and height, with caching
func (rm *ResourceManager) GetScaledSprite(ctx context.Context, name string, width, height int) (*ebiten.Image, error) {
	cacheKey := fmt.Sprintf("%s_%dx%d", name, width, height)

	rm.mutex.RLock()
	if img, ok := rm.scaledCache[cacheKey]; ok {
		rm.mutex.RUnlock()
		return img, nil
	}
	rm.mutex.RUnlock()

	// Get the base sprite
	sprite, ok := rm.GetSprite(ctx, name)
	if !ok || sprite == nil {
		return nil, fmt.Errorf("sprite '%s' not found", name)
	}

	// If already correct size, return as is
	bounds := sprite.Bounds()
	if bounds.Dx() == width && bounds.Dy() == height {
		return sprite, nil
	}

	// Scale the sprite
	scaled := ebiten.NewImage(width, height)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(width)/float64(bounds.Dx()), float64(height)/float64(bounds.Dy()))
	scaled.DrawImage(sprite, op)

	rm.mutex.Lock()
	rm.scaledCache[cacheKey] = scaled
	rm.mutex.Unlock()

	return scaled, nil
}

// GetUISprite is a convenience method for square UI icons
func (rm *ResourceManager) GetUISprite(ctx context.Context, name string, size int) (*ebiten.Image, error) {
	return rm.GetScaledSprite(ctx, name, size, size)
}
