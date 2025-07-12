package resources

import (
	"bytes"
	"fmt"
	"image/color"
	"image/png"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/jonesrussell/gimbal/assets"
	"github.com/jonesrussell/gimbal/internal/common"
)

// LoadSprite loads a sprite from the embedded assets
func (rm *ResourceManager) LoadSprite(name, path string) (*ebiten.Image, error) {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	// Check if already loaded
	if resource, exists := rm.resources[name]; exists {
		if sprite, ok := resource.Data.(*ebiten.Image); ok {
			rm.logger.Debug("Sprite reused", "name", name)
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
		return nil, common.NewGameErrorWithCause(common.ErrorCodeAssetLoadFailed, "failed to read sprite file", err)
	}

	rm.logger.Debug("Sprite file read successfully", "name", name, "size", len(imageData), "path", path)

	// Use PNG decoder specifically for PNG files
	img, err := png.Decode(bytes.NewReader(imageData))
	if err != nil {
		rm.logger.Error("Failed to decode PNG sprite", "name", name, "path", path, "error", err)
		return nil, common.NewGameErrorWithCause(common.ErrorCodeAssetInvalid, "failed to decode sprite", err)
	}

	rm.logger.Debug("Sprite decoded successfully", "name", name, "bounds", img.Bounds())

	sprite := ebiten.NewImageFromImage(img)

	// Store in resource manager
	rm.resources[name] = &Resource{
		Type: ResourceSprite,
		Name: name,
		Data: sprite,
	}

	rm.logger.Debug("Sprite loaded", "name", name, "path", path, "bounds", img.Bounds())
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
			rm.logger.Debug("Sprite reused", "name", name)
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

// LoadAllSprites loads all required sprites for the game
func (rm *ResourceManager) LoadAllSprites() error {
	// Load player sprite from file
	_, err := rm.LoadSprite("player", "sprites/player.png")
	if err != nil {
		rm.logger.Warn("Failed to load player sprite, using placeholder", "error", err)
		_, err = rm.CreateSprite("player", 32, 32, color.RGBA{0, 255, 0, 255})
		if err != nil {
			return common.NewGameErrorWithCause(
				common.ErrorCodeAssetLoadFailed,
				"failed to create player placeholder",
				err,
			)
		}
	}

	// Load heart sprite for lives display
	_, err = rm.LoadSprite("heart", "sprites/heart.png")
	if err != nil {
		rm.logger.Warn("Failed to load heart sprite, using placeholder", "error", err)
		_, err = rm.CreateSprite("heart", 16, 16, color.RGBA{255, 0, 0, 255})
		if err != nil {
			return common.NewGameErrorWithCause(
				common.ErrorCodeAssetLoadFailed,
				"failed to create heart placeholder",
				err,
			)
		}
	}

	// Load enemy sprite sheet for animations
	_, err = rm.LoadSprite("enemy_sheet", "sprites/enemy_sheet.png")
	if err != nil {
		rm.logger.Warn("Failed to load enemy sprite sheet, using placeholder", "error", err)
		_, err = rm.CreateSprite("enemy_sheet", 512, 256, color.RGBA{255, 0, 0, 255})
		if err != nil {
			return common.NewGameErrorWithCause(
				common.ErrorCodeAssetLoadFailed,
				"failed to create enemy sheet placeholder",
				err,
			)
		}
	}

	// Create star sprite
	_, err = rm.CreateSprite("star", common.StarSpriteSize, common.StarSpriteSize, color.White)
	if err != nil {
		return common.NewGameErrorWithCause(common.ErrorCodeAssetLoadFailed, "failed to create star sprite", err)
	}

	// Create UI sprites
	_, err = rm.CreateSprite("button", common.ButtonSpriteWidth, common.ButtonSpriteHeight,
		color.RGBA{common.ButtonColorR, common.ButtonColorG, common.ButtonColorB, common.ButtonColorA})
	if err != nil {
		return common.NewGameErrorWithCause(common.ErrorCodeAssetLoadFailed, "failed to create button sprite", err)
	}

	_, err = rm.CreateSprite("background", 1, 1, color.Black)
	if err != nil {
		return common.NewGameErrorWithCause(common.ErrorCodeAssetLoadFailed, "failed to create background sprite", err)
	}

	rm.logger.Info("All sprites loaded successfully")
	return nil
}
