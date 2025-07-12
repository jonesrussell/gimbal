package resources

import (
	"context"
	"sync"

	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/jonesrussell/gimbal/internal/common"
)

// ResourceType represents different types of resources
type ResourceType int

const (
	ResourceSprite ResourceType = iota
	ResourceSound
	ResourceFont
	ResourceData
)

// Resource represents a loaded game resource
type Resource struct {
	Type ResourceType
	Name string
	Data interface{}
}

// ResourceManager manages all game resources
type ResourceManager struct {
	resources   map[string]*Resource
	mutex       sync.RWMutex
	logger      common.Logger
	defaultFont text.Face
}

// NewResourceManager creates a new resource manager
func NewResourceManager(logger common.Logger) *ResourceManager {
	rm := &ResourceManager{
		resources: make(map[string]*Resource),
		logger:    logger,
	}
	if err := rm.loadDefaultFont(context.Background()); err != nil {
		logger.Error("failed to load default font", "error", err)
	}
	return rm
}

// GetResourceCount returns the number of loaded resources
func (rm *ResourceManager) GetResourceCount() int {
	rm.mutex.RLock()
	defer rm.mutex.RUnlock()
	return len(rm.resources)
}

// GetResourceInfo returns information about loaded resources
func (rm *ResourceManager) GetResourceInfo() map[string]interface{} {
	rm.mutex.RLock()
	defer rm.mutex.RUnlock()

	info := make(map[string]interface{})
	for name, resource := range rm.resources {
		info[name] = map[string]interface{}{
			"type": resource.Type,
		}
	}
	return info
}

// Cleanup releases all resources
func (rm *ResourceManager) Cleanup(ctx context.Context) error {
	// Check for cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	rm.logger.Info("Cleaning up resources", "count", len(rm.resources))
	rm.resources = make(map[string]*Resource)
	return nil
}

// Predefined sprite names for easy access
const (
	SpritePlayer     = "player"
	SpriteStar       = "star"
	SpriteButton     = "button"
	SpriteBackground = "background"
)
