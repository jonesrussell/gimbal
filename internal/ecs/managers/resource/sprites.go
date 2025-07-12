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
	"github.com/jonesrussell/gimbal/internal/ui"
)

// LoadSprite loads a sprite from the embedded assets
func (rm *ResourceManager) LoadSprite(ctx context.Context, name, path string) (*ebiten.Image, error) {
	// Check for cancellation
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	// Check if already loaded
	if resource, exists := rm.resources[name]; exists {
		if sprite, ok := resource.Data.(*ebiten.Image); ok {
			rm.logger.Debug("[SPRITE_CACHE] Sprite reused from cache",
				"name", name, "was_processed", name == "enemy_sheet")
			return sprite, nil
		}
	}

	// Load from embedded assets
	rm.logger.Debug("[SPRITE_LOAD] Attempting to load sprite from embed", "name", name, "path", path)

	// List all embedded files for debugging
	files, listErr := assets.Assets.ReadDir("sprites")
	if listErr != nil {
		rm.logger.Error("[SPRITE_ERROR] Failed to list embedded files", "error", listErr)
	} else {
		rm.logger.Debug("[SPRITE_FILES] Embedded files found", "files", files)
		for _, f := range files {
			rm.logger.Debug("[SPRITE_FILES] Embedded file", "name", f.Name(), "is_dir", f.IsDir())
		}
	}

	imageData, err := assets.Assets.ReadFile(path)
	if err != nil {
		rm.logger.Error("[SPRITE_ERROR] Failed to read sprite file from embed",
			"name", name, "path", path, "error", err)
		return nil, errors.NewGameErrorWithCause(errors.ErrorCodeAssetLoadFailed, "failed to read sprite file", err)
	}

	rm.logger.Debug("[SPRITE_LOAD] Sprite file read successfully", "name", name, "size", len(imageData), "path", path)

	// Use PNG decoder specifically for PNG files
	img, err := png.Decode(bytes.NewReader(imageData))
	if err != nil {
		rm.logger.Error("[SPRITE_ERROR] Failed to decode PNG sprite", "name", name, "path", path, "error", err)
		return nil, errors.NewGameErrorWithCause(errors.ErrorCodeAssetInvalid, "failed to decode sprite", err)
	}

	rm.logger.Debug("[SPRITE_DECODE] Sprite decoded successfully", "name", name, "bounds", img.Bounds())

	// Create ebiten image from decoded image
	sprite := ebiten.NewImageFromImage(img)
	rm.logger.Debug("[SPRITE_PROCESS] Created ebiten image from decoded sprite")

	// Store in resource manager
	rm.resources[name] = &Resource{
		Type: ResourceSprite,
		Name: name,
		Data: sprite,
	}

	rm.logger.Debug("[SPRITE_LOAD] Sprite loaded", "name", name, "path", path, "bounds", img.Bounds())
	return sprite, nil
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

// LoadAllSprites loads all required sprites for the game
func (rm *ResourceManager) LoadAllSprites(ctx context.Context) error {
	// Check for cancellation at the start
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Load player sprite from file
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	_, err := rm.LoadSprite(ctx, "player", "sprites/player.png")
	if err != nil {
		rm.logger.Warn("Failed to load player sprite, using placeholder", "error", err)
		_, err = rm.CreateSprite("player", 32, 32, color.RGBA{0, 255, 0, 255})
		if err != nil {
			return errors.NewGameErrorWithCause(
				errors.ErrorCodeAssetLoadFailed,
				"failed to create player placeholder",
				err,
			)
		}
	}

	// Load heart sprite for lives display
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	_, err = rm.LoadSprite(ctx, "heart", "sprites/heart.png")
	if err != nil {
		rm.logger.Warn("Failed to load heart sprite, using placeholder", "error", err)
		_, err = rm.CreateSprite("heart", 16, 16, color.RGBA{255, 0, 0, 255})
		if err != nil {
			return errors.NewGameErrorWithCause(
				errors.ErrorCodeAssetLoadFailed,
				"failed to create heart placeholder",
				err,
			)
		}
	}

	// Load enemy sprite
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	rm.logger.Debug("[SPRITE_LOAD] Attempting to load enemy sprite", "path", "sprites/enemy.png")
	_, err = rm.LoadSprite(ctx, "enemy", "sprites/enemy.png")
	if err != nil {
		rm.logger.Warn("[SPRITE_WARN] Failed to load enemy sprite, using placeholder", "error", err)
		_, err = rm.CreateSprite("enemy", 32, 32, color.RGBA{255, 0, 0, 255})
		if err != nil {
			return errors.NewGameErrorWithCause(
				errors.ErrorCodeAssetLoadFailed,
				"failed to create enemy placeholder",
				err,
			)
		}
	} else {
		rm.logger.Debug("[SPRITE_LOAD] Enemy sprite loaded successfully")
	}

	// Create star sprite
	_, err = rm.CreateSprite("star", ui.StarSpriteSize, ui.StarSpriteSize, color.White)
	if err != nil {
		return errors.NewGameErrorWithCause(errors.ErrorCodeAssetLoadFailed, "failed to create star sprite", err)
	}

	// Create UI sprites
	_, err = rm.CreateSprite("button", ui.ButtonSpriteWidth, ui.ButtonSpriteHeight,
		color.RGBA{ui.ButtonColorR, ui.ButtonColorG, ui.ButtonColorB, ui.ButtonColorA})
	if err != nil {
		return errors.NewGameErrorWithCause(errors.ErrorCodeAssetLoadFailed, "failed to create button sprite", err)
	}

	_, err = rm.CreateSprite("background", 1, 1, color.Black)
	if err != nil {
		return errors.NewGameErrorWithCause(errors.ErrorCodeAssetLoadFailed, "failed to create background sprite", err)
	}

	rm.logger.Info("[SPRITE_LOAD] All sprites loaded successfully")
	return nil
}
