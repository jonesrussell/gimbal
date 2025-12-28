package resources

import (
	"context"
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/errors"
	uicore "github.com/jonesrussell/gimbal/internal/ui/core"
)

// SpriteLoadConfig holds configuration for loading sprites with fallback
type SpriteLoadConfig struct {
	Name           string
	Path           string
	FallbackWidth  int
	FallbackHeight int
	FallbackColor  color.Color
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

// loadSpriteWithFallback loads a sprite with fallback creation
func (rm *ResourceManager) loadSpriteWithFallback(ctx context.Context, config SpriteLoadConfig) error {
	if err := common.CheckContextCancellation(ctx); err != nil {
		return err
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

// loadGameSprites loads all game sprite types using a configuration-driven approach
func (rm *ResourceManager) loadGameSprites(ctx context.Context) error {
	spriteConfigs := []SpriteLoadConfig{
		{
			Name:           "player",
			Path:           "sprites/player.png",
			FallbackWidth:  32,
			FallbackHeight: 32,
			FallbackColor:  color.RGBA{0, 255, 0, 255}, // Green fallback
		},
		{
			Name:           "heart",
			Path:           "sprites/heart.png",
			FallbackWidth:  16,
			FallbackHeight: 16,
			FallbackColor:  color.RGBA{255, 0, 0, 255}, // Red fallback
		},
		{
			Name:           "enemy",
			Path:           "sprites/enemy.png",
			FallbackWidth:  32,
			FallbackHeight: 32,
			FallbackColor:  color.RGBA{255, 0, 0, 255}, // Red fallback
		},
		{
			Name:           "enemy_heavy",
			Path:           "sprites/enemy_heavy.png",
			FallbackWidth:  32,
			FallbackHeight: 32,
			FallbackColor:  color.RGBA{255, 165, 0, 255}, // Orange fallback
		},
		{
			Name:           "enemy_boss",
			Path:           "sprites/enemy_boss.png",
			FallbackWidth:  64,
			FallbackHeight: 64,
			FallbackColor:  color.RGBA{128, 0, 128, 255}, // Purple fallback
		},
		{
			Name:           "enemy_ammo",
			Path:           "sprites/enemy_ammo.png",
			FallbackWidth:  6,
			FallbackHeight: 6,
			FallbackColor:  color.RGBA{255, 50, 50, 255}, // Red fallback
		},
		{
			Name:           "enemy_heavy_ammo",
			Path:           "sprites/enemy_heavy_ammo.png",
			FallbackWidth:  6,
			FallbackHeight: 6,
			FallbackColor:  color.RGBA{255, 165, 0, 255}, // Orange fallback
		},
		{
			Name:           "star",
			Path:           "sprites/star.png",
			FallbackWidth:  uicore.StarSpriteSize,
			FallbackHeight: uicore.StarSpriteSize,
			FallbackColor:  color.White,
		},
	}

	for _, cfg := range spriteConfigs {
		if err := rm.loadSpriteWithFallback(ctx, cfg); err != nil {
			return err
		}
	}
	return nil
}

// createUISprites creates UI-related sprites
func (rm *ResourceManager) createUISprites(ctx context.Context) error {
	// Create button sprite
	if _, err := rm.CreateSprite("button", uicore.ButtonSpriteWidth, uicore.ButtonSpriteHeight,
		color.RGBA{uicore.ButtonColorR, uicore.ButtonColorG, uicore.ButtonColorB, uicore.ButtonColorA}); err != nil {
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
	if err := common.CheckContextCancellation(ctx); err != nil {
		return err
	}

	// Load game sprites (player, enemies, heart)
	if err := rm.loadGameSprites(ctx); err != nil {
		return err
	}

	// Create UI sprites (button, background)
	if err := rm.createUISprites(ctx); err != nil {
		return err
	}

	rm.logger.Info("[SPRITE_LOAD] All sprites loaded successfully")
	return nil
}
