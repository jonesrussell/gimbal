package resources

import (
	"bytes"
	"context"
	"fmt"
	"image/color"
	"image/png"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/jonesrussell/gimbal/assets"
	"github.com/jonesrussell/gimbal/internal/errors"
	"github.com/jonesrussell/gimbal/internal/ui/core"
)

// LoadSprite loads a sprite from embedded assets with simplified logic
func (rm *ResourceManager) LoadSprite(ctx context.Context, name, path string) (*ebiten.Image, error) {
	// Check context cancellation
	if err := rm.checkContext(ctx); err != nil {
		return nil, err
	}

	// Try cache first
	if cached := rm.getCachedSprite(name); cached != nil {
		rm.logger.Debug("[SPRITE_CACHE] Sprite reused from cache", "name", name)
		return cached, nil
	}

	// Load and decode sprite
	return rm.loadAndCacheSprite(ctx, name, path)
}

// checkContext verifies context is not cancelled
func (rm *ResourceManager) checkContext(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
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
	rm.logger.Debug("[SPRITE_LOAD] Loading sprite from embed", "name", name, "path", path)

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

	rm.logger.Debug("[SPRITE_LOAD] Sprite loaded successfully", "name", name)
	return sprite, nil
}

// loadSpriteFile loads sprite file from embedded assets
func (rm *ResourceManager) loadSpriteFile(path string) ([]byte, error) {
	imageData, err := assets.Assets.ReadFile(path)
	if err != nil {
		rm.logger.Error("[SPRITE_ERROR] Failed to read sprite file", "path", path, "error", err)
		return nil, errors.NewGameErrorWithCause(errors.AssetLoadFailed, "failed to read sprite file", err)
	}

	rm.logger.Debug("[SPRITE_LOAD] Sprite file read successfully", "path", path, "size", len(imageData))
	return imageData, nil
}

// decodePNGData decodes PNG data into ebiten image
func (rm *ResourceManager) decodePNGData(name, path string, imageData []byte) (*ebiten.Image, error) {
	img, err := png.Decode(bytes.NewReader(imageData))
	if err != nil {
		rm.logger.Error("[SPRITE_ERROR] Failed to decode PNG sprite", "name", name, "path", path, "error", err)
		return nil, errors.NewGameErrorWithCause(errors.AssetInvalid, "failed to decode sprite", err)
	}

	sprite := ebiten.NewImageFromImage(img)
	rm.logger.Debug("[SPRITE_DECODE] Sprite decoded successfully", "name", name, "bounds", img.Bounds())

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

// debugListEmbeddedFiles lists embedded files for debugging (optional helper)
func (rm *ResourceManager) debugListEmbeddedFiles() {
	files, err := assets.Assets.ReadDir("sprites")
	if err != nil {
		rm.logger.Error("[SPRITE_ERROR] Failed to list embedded files", "error", err)
		return
	}

	rm.logger.Debug("[SPRITE_FILES] Embedded files found", "count", len(files))
	for _, f := range files {
		rm.logger.Debug("[SPRITE_FILES] Embedded file", "name", f.Name(), "is_dir", f.IsDir())
	}
}

// CreateSprite creates a simple colored sprite
func (rm *ResourceManager) CreateSprite(
	name string, width, height int, spriteColor color.Color,
) (*ebiten.Image, error) {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	// Check if already created
	if resource, exists := rm.resources[name]; exists {
		if sprite, ok := resource.Data.(*ebiten.Image); ok {
			rm.logger.Debug("[SPRITE_CACHE] Sprite reused", "name", name)
			return sprite, nil
		}
	}

	// Create new sprite
	sprite := ebiten.NewImage(width, height)
	sprite.Fill(spriteColor)

	// Store in resource manager
	rm.resources[name] = &Resource{
		Type: ResourceSprite,
		Name: name,
		Data: sprite,
	}

	rm.logger.Debug("[SPRITE_CREATE] Sprite created", "name", name, "size", fmt.Sprintf("%dx%d", width, height))
	return sprite, nil
}

// GetSprite retrieves a loaded sprite
func (rm *ResourceManager) GetSprite(ctx context.Context, name string) (*ebiten.Image, bool) {
	// Check for cancellation
	select {
	case <-ctx.Done():
		return nil, false
	default:
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

type SpriteLoadConfig struct {
	Name           string
	Path           string
	FallbackWidth  int
	FallbackHeight int
	FallbackColor  color.Color
}

func (rm *ResourceManager) loadSpriteWithFallback(ctx context.Context, config SpriteLoadConfig) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	_, err := rm.LoadSprite(ctx, config.Name, config.Path)
	if err != nil {
		rm.logger.Warn("Failed to load sprite, using placeholder", "name", config.Name, "error", err)
		_, err = rm.CreateSprite(config.Name, config.FallbackWidth, config.FallbackHeight, config.FallbackColor)
		if err != nil {
			return errors.NewGameErrorWithCause(
				errors.AssetLoadFailed,
				fmt.Sprintf("failed to create %s placeholder", config.Name),
				err,
			)
		}
	}
	return nil
}

func (rm *ResourceManager) loadPlayerSprites(ctx context.Context) error {
	return rm.loadSpriteWithFallback(ctx, SpriteLoadConfig{
		Name:           "player",
		Path:           "sprites/player.png",
		FallbackWidth:  32,
		FallbackHeight: 32,
		FallbackColor:  color.RGBA{0, 255, 0, 255},
	})
}

func (rm *ResourceManager) loadHeartSprites(ctx context.Context) error {
	return rm.loadSpriteWithFallback(ctx, SpriteLoadConfig{
		Name:           "heart",
		Path:           "sprites/heart.png",
		FallbackWidth:  16,
		FallbackHeight: 16,
		FallbackColor:  color.RGBA{255, 0, 0, 255},
	})
}

func (rm *ResourceManager) loadEnemySprites(ctx context.Context) error {
	rm.logger.Debug("[SPRITE_LOAD] Attempting to load enemy sprite", "path", "sprites/enemy.png")
	err := rm.loadSpriteWithFallback(ctx, SpriteLoadConfig{
		Name:           "enemy",
		Path:           "sprites/enemy.png",
		FallbackWidth:  32,
		FallbackHeight: 32,
		FallbackColor:  color.RGBA{255, 0, 0, 255},
	})
	if err == nil {
		rm.logger.Debug("[SPRITE_LOAD] Enemy sprite loaded successfully")
	}
	return err
}

// createUISprites creates UI-related sprites
func (rm *ResourceManager) createUISprites(ctx context.Context) error {
	// Create star sprite
	if _, err := rm.CreateSprite("star", core.StarSpriteSize, core.StarSpriteSize, color.White); err != nil {
		return errors.NewGameErrorWithCause(errors.AssetLoadFailed, "failed to create star sprite", err)
	}

	// Create button sprite
	if _, err := rm.CreateSprite("button", core.ButtonSpriteWidth, core.ButtonSpriteHeight,
		color.RGBA{core.ButtonColorR, core.ButtonColorG, core.ButtonColorB, core.ButtonColorA}); err != nil {
		return errors.NewGameErrorWithCause(errors.AssetLoadFailed, "failed to create button sprite", err)
	}

	// Create background sprite
	if _, err := rm.CreateSprite("background", 1, 1, color.Black); err != nil {
		return errors.NewGameErrorWithCause(errors.AssetLoadFailed, "failed to create background sprite", err)
	}

	return nil
}

// LoadAllSprites loads all required sprites for the game
func (rm *ResourceManager) LoadAllSprites(ctx context.Context) error {
	// Check for cancellation at the start
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Load different sprite types
	if err := rm.loadPlayerSprites(ctx); err != nil {
		return err
	}
	if err := rm.loadHeartSprites(ctx); err != nil {
		return err
	}
	if err := rm.loadEnemySprites(ctx); err != nil {
		return err
	}
	if err := rm.createUISprites(ctx); err != nil {
		return err
	}

	rm.logger.Info("[SPRITE_LOAD] All sprites loaded successfully")
	return nil
}

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
