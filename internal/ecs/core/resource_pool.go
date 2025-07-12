package core

import (
	"sync"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/jonesrussell/gimbal/internal/common"
)

// ImagePool manages reusable ebiten.Image instances
type ImagePool struct {
	pool   map[string][]*ebiten.Image
	mutex  sync.RWMutex
	logger common.Logger
}

// NewImagePool creates a new image pool
func NewImagePool(logger common.Logger) *ImagePool {
	return &ImagePool{
		pool:   make(map[string][]*ebiten.Image),
		logger: logger,
	}
}

// GetImage retrieves an image from the pool or creates a new one
func (ip *ImagePool) GetImage(width, height int) *ebiten.Image {
	key := ip.createKey(width, height)

	ip.mutex.Lock()
	defer ip.mutex.Unlock()

	// Check if we have a cached image
	if images, exists := ip.pool[key]; exists && len(images) > 0 {
		// Pop the last image from the pool
		image := images[len(images)-1]
		ip.pool[key] = images[:len(images)-1]
		ip.logger.Debug("Image reused from pool", "size", key)
		return image
	}

	// Create new image
	image := ebiten.NewImage(width, height)
	ip.logger.Debug("New image created", "size", key)
	return image
}

// ReturnImage returns an image to the pool for reuse
func (ip *ImagePool) ReturnImage(image *ebiten.Image) {
	if image == nil {
		return
	}

	bounds := image.Bounds()
	key := ip.createKey(bounds.Dx(), bounds.Dy())

	ip.mutex.Lock()
	defer ip.mutex.Unlock()

	// Clear the image before returning to pool
	image.Clear()

	// Add to pool
	ip.pool[key] = append(ip.pool[key], image)
	ip.logger.Debug("Image returned to pool", "size", key, "pool_size", len(ip.pool[key]))
}

// createKey creates a string key for the image dimensions
func (ip *ImagePool) createKey(width, height int) string {
	return string(rune(width)) + "x" + string(rune(height))
}

// Cleanup releases all pooled images
func (ip *ImagePool) Cleanup() {
	ip.mutex.Lock()
	defer ip.mutex.Unlock()

	totalImages := 0
	for key, images := range ip.pool {
		totalImages += len(images)
		ip.pool[key] = nil
	}

	ip.logger.Info("Image pool cleaned up", "total_images", totalImages)
	ip.pool = make(map[string][]*ebiten.Image)
}

// GetPoolStats returns statistics about the image pool
func (ip *ImagePool) GetPoolStats() map[string]interface{} {
	ip.mutex.RLock()
	defer ip.mutex.RUnlock()

	stats := make(map[string]interface{})
	totalImages := 0

	for key, images := range ip.pool {
		stats[key] = len(images)
		totalImages += len(images)
	}

	stats["total_pooled"] = totalImages
	stats["pool_count"] = len(ip.pool)

	return stats
}
