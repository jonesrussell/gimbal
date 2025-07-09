package ecs

import (
	"bytes"
	"fmt"
	"image/color"
	"image/png"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/jonesrussell/gimbal/assets"
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
	Type     ResourceType
	Name     string
	Data     interface{}
	RefCount int
}

// ResourceManager manages all game resources
type ResourceManager struct {
	resources map[string]*Resource
	mutex     sync.RWMutex
	logger    common.Logger
}

// NewResourceManager creates a new resource manager
func NewResourceManager(logger common.Logger) *ResourceManager {
	return &ResourceManager{
		resources: make(map[string]*Resource),
		logger:    logger,
	}
}

// LoadSprite loads a sprite from the embedded assets
func (rm *ResourceManager) LoadSprite(name, path string) (*ebiten.Image, error) {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	// Check if already loaded
	if resource, exists := rm.resources[name]; exists {
		if sprite, ok := resource.Data.(*ebiten.Image); ok {
			resource.RefCount++
			rm.logger.Debug("Sprite reused", "name", name, "ref_count", resource.RefCount)
			return sprite, nil
		}
	}

	// Load from embedded assets
	rm.logger.Debug("Attempting to load sprite from embed", "name", name, "path", path)

	// List all embedded files for debugging
	files, listErr := assets.Assets.ReadDir("sprites")
	if listErr != nil {
		rm.logger.Error("Failed to list embedded files", "error", listErr)
	} else {
		rm.logger.Debug("Embedded files found", "files", files)
		for _, f := range files {
			rm.logger.Debug("Embedded file", "name", f.Name(), "is_dir", f.IsDir())
		}
	}

	imageData, err := assets.Assets.ReadFile(path)
	if err != nil {
		rm.logger.Error("Failed to read sprite file from embed", "name", name, "path", path, "error", err)
		return nil, fmt.Errorf("failed to read sprite file %s: %w", path, err)
	}

	rm.logger.Debug("Sprite file read successfully", "name", name, "size", len(imageData), "path", path)

	// Use PNG decoder specifically for PNG files
	img, err := png.Decode(bytes.NewReader(imageData))
	if err != nil {
		rm.logger.Error("Failed to decode PNG sprite", "name", name, "path", path, "error", err)
		return nil, fmt.Errorf("failed to decode sprite %s: %w", path, err)
	}

	rm.logger.Debug("Sprite decoded successfully", "name", name, "bounds", img.Bounds())

	sprite := ebiten.NewImageFromImage(img)

	// Store in resource manager
	rm.resources[name] = &Resource{
		Type:     ResourceSprite,
		Name:     name,
		Data:     sprite,
		RefCount: 1,
	}

	rm.logger.Debug("Sprite loaded", "name", name, "path", path, "bounds", img.Bounds())
	return sprite, nil
}

// CreateSprite creates a simple colored sprite
func (rm *ResourceManager) CreateSprite(name string, width, height int, color color.Color) (*ebiten.Image, error) {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	// Check if already created
	if resource, exists := rm.resources[name]; exists {
		if sprite, ok := resource.Data.(*ebiten.Image); ok {
			resource.RefCount++
			rm.logger.Debug("Sprite reused", "name", name, "ref_count", resource.RefCount)
			return sprite, nil
		}
	}

	// Create new sprite
	sprite := ebiten.NewImage(width, height)
	sprite.Fill(color)

	// Store in resource manager
	rm.resources[name] = &Resource{
		Type:     ResourceSprite,
		Name:     name,
		Data:     sprite,
		RefCount: 1,
	}

	rm.logger.Debug("Sprite created", "name", name, "size", fmt.Sprintf("%dx%d", width, height))
	return sprite, nil
}

// GetSprite retrieves a loaded sprite
func (rm *ResourceManager) GetSprite(name string) (*ebiten.Image, bool) {
	rm.mutex.RLock()
	defer rm.mutex.RUnlock()

	if resource, exists := rm.resources[name]; exists {
		if sprite, ok := resource.Data.(*ebiten.Image); ok {
			return sprite, true
		}
	}
	return nil, false
}

// ReleaseSprite releases a sprite reference
func (rm *ResourceManager) ReleaseSprite(name string) {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	if resource, exists := rm.resources[name]; exists {
		resource.RefCount--
		rm.logger.Debug("Sprite reference released", "name", name, "ref_count", resource.RefCount)

		// Remove if no more references
		if resource.RefCount <= 0 {
			delete(rm.resources, name)
			rm.logger.Debug("Sprite removed", "name", name)
		}
	}
}

// LoadAllSprites loads all required sprites for the game
func (rm *ResourceManager) LoadAllSprites() error {
	// Load player sprite from file
	_, err := rm.LoadSprite("player", "sprites/player.png")
	if err != nil {
		rm.logger.Warn("Failed to load player sprite, using placeholder", "error", err)
		_, err = rm.CreateSprite("player", 32, 32, color.RGBA{0, 255, 0, 255})
		if err != nil {
			return fmt.Errorf("failed to create player placeholder: %w", err)
		}
	}

	// Create star sprite
	_, err = rm.CreateSprite("star", 10, 10, color.White)
	if err != nil {
		return fmt.Errorf("failed to create star sprite: %w", err)
	}

	// Create UI sprites
	_, err = rm.CreateSprite("button", 100, 30, color.RGBA{100, 100, 100, 255})
	if err != nil {
		return fmt.Errorf("failed to create button sprite: %w", err)
	}

	_, err = rm.CreateSprite("background", 1, 1, color.Black)
	if err != nil {
		return fmt.Errorf("failed to create background sprite: %w", err)
	}

	rm.logger.Info("All sprites loaded successfully")
	return nil
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
			"type":      resource.Type,
			"ref_count": resource.RefCount,
		}
	}
	return info
}

// Cleanup releases all resources
func (rm *ResourceManager) Cleanup() {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	rm.logger.Info("Cleaning up resources", "count", len(rm.resources))
	rm.resources = make(map[string]*Resource)
}

// Predefined sprite names for easy access
const (
	SpritePlayer     = "player"
	SpriteStar       = "star"
	SpriteButton     = "button"
	SpriteBackground = "background"
)
