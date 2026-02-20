package resources

import (
	"context"
	"fmt"
	"image/color"
	"log"

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

	return sprite, nil
}

// loadSpriteWithFallback loads a sprite with fallback creation
func (rm *ResourceManager) loadSpriteWithFallback(ctx context.Context, config SpriteLoadConfig) error {
	if err := common.CheckContextCancellation(ctx); err != nil {
		return err
	}

	_, err := rm.LoadSprite(ctx, config.Name, config.Path)
	if err != nil {
		log.Printf("[WARN] Failed to load sprite, using placeholder: name=%s error=%v", config.Name, err)
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

// loadGameSprites loads all game sprite types using a configuration-driven approach.
//
//nolint:funlen // Config slice is inherently long; splitting would not reduce complexity.
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
			Name:           "ammo",
			Path:           "sprites/ammo.png",
			FallbackWidth:  16,
			FallbackHeight: 16,
			FallbackColor:  color.RGBA{255, 255, 0, 255}, // Yellow fallback
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

// loadUISprites loads UI-related sprites
func (rm *ResourceManager) loadUISprites(ctx context.Context) error {
	spriteConfigs := []SpriteLoadConfig{
		{
			Name:           "title_logo",
			Path:           "ui/title_logo.png",
			FallbackWidth:  320,
			FallbackHeight: 80,
			FallbackColor:  color.RGBA{0, 255, 255, 255}, // Cyan fallback
		},
		{
			Name:           "studio_screen",
			Path:           "ui/studio_screen.png",
			FallbackWidth:  640,
			FallbackHeight: 480,
			FallbackColor:  color.Black,
		},
		{
			Name:           "menu_frame",
			Path:           "ui/menu_frame.png",
			FallbackWidth:  640,
			FallbackHeight: 480,
			FallbackColor:  color.RGBA{128, 0, 255, 255}, // Purple fallback
		},
		{
			Name:           "button_highlight",
			Path:           "ui/button_highlight.png",
			FallbackWidth:  200,
			FallbackHeight: 40,
			FallbackColor:  color.RGBA{255, 0, 255, 255}, // Magenta fallback
		},
		{
			Name:           "warning_overlay",
			Path:           "ui/warning_overlay.png",
			FallbackWidth:  640,
			FallbackHeight: 480,
			FallbackColor:  color.RGBA{255, 0, 0, 128}, // Red semi-transparent fallback
		},
		{
			Name:           "scanline_overlay",
			Path:           "ui/scanline_overlay.png",
			FallbackWidth:  640,
			FallbackHeight: 480,
			FallbackColor:  color.RGBA{0, 0, 0, 128}, // Black semi-transparent fallback
		},
	}

	for _, cfg := range spriteConfigs {
		if err := rm.loadSpriteWithFallback(ctx, cfg); err != nil {
			log.Printf("[WARN] Failed to load UI sprite, using fallback: name=%s error=%v", cfg.Name, err)
			// Continue with fallback
		}
	}
	return nil
}

// loadCutsceneSprites loads cutscene-related sprites
func (rm *ResourceManager) loadCutsceneSprites(ctx context.Context) error {
	// Scanning grid
	spriteConfigs := []SpriteLoadConfig{
		{
			Name:           "scanning_grid",
			Path:           "cutscenes/scanning_grid.png",
			FallbackWidth:  640,
			FallbackHeight: 480,
			FallbackColor:  color.RGBA{0, 255, 255, 64}, // Cyan semi-transparent fallback
		},
	}

	// Warp tunnel frames
	for i := 1; i <= 8; i++ {
		cfg := SpriteLoadConfig{
			Name:           fmt.Sprintf("warp_tunnel_%02d", i),
			Path:           fmt.Sprintf("cutscenes/warp_tunnel_%02d.png", i),
			FallbackWidth:  640,
			FallbackHeight: 480,
			FallbackColor:  color.RGBA{128, 0, 255, 128}, // Purple semi-transparent fallback
		}
		spriteConfigs = append(spriteConfigs, cfg)
	}

	for _, cfg := range spriteConfigs {
		if err := rm.loadSpriteWithFallback(ctx, cfg); err != nil {
			log.Printf("[WARN] Failed to load cutscene sprite, using fallback: name=%s error=%v", cfg.Name, err)
			// Continue with fallback
		}
	}
	return nil
}

// loadPlanetSprites loads planet portrait sprites
func (rm *ResourceManager) loadPlanetSprites(ctx context.Context) error {
	planets := []string{"earth", "mars", "jupiter", "saturn", "uranus", "neptune", "pluto"}

	for _, planet := range planets {
		cfg := SpriteLoadConfig{
			Name:           fmt.Sprintf("planet_%s", planet),
			Path:           fmt.Sprintf("planets/%s.png", planet),
			FallbackWidth:  128,
			FallbackHeight: 128,
			FallbackColor:  color.RGBA{0, 128, 255, 255}, // Blue fallback
		}
		if err := rm.loadSpriteWithFallback(ctx, cfg); err != nil {
			log.Printf("[WARN] Failed to load planet sprite, using fallback: planet=%s error=%v", planet, err)
			// Continue with fallback
		}
	}
	return nil
}

// loadBossPortraitSprites loads boss portrait sprites
func (rm *ResourceManager) loadBossPortraitSprites(ctx context.Context) error {
	bosses := []string{"earth", "mars", "jupiter", "saturn", "uranus", "neptune"}

	for _, boss := range bosses {
		cfg := SpriteLoadConfig{
			Name:           fmt.Sprintf("boss_portrait_%s", boss),
			Path:           fmt.Sprintf("bosses/%s_boss_portrait.png", boss),
			FallbackWidth:  128,
			FallbackHeight: 128,
			FallbackColor:  color.RGBA{255, 0, 0, 255}, // Red fallback
		}
		if err := rm.loadSpriteWithFallback(ctx, cfg); err != nil {
			log.Printf("[WARN] Failed to load boss portrait sprite, using fallback: boss=%s error=%v", boss, err)
			// Continue with fallback
		}
	}
	return nil
}

// loadEndingSprites loads ending sequence sprites
func (rm *ResourceManager) loadEndingSprites(ctx context.Context) error {
	spriteConfigs := []SpriteLoadConfig{
		{
			Name:           "starfield_bg",
			Path:           "ending/starfield_bg.png",
			FallbackWidth:  640,
			FallbackHeight: 480,
			FallbackColor:  color.Black,
		},
		{
			Name:           "mission_complete",
			Path:           "ending/mission_complete.png",
			FallbackWidth:  400,
			FallbackHeight: 100,
			FallbackColor:  color.RGBA{0, 255, 255, 255}, // Cyan fallback
		},
	}

	for _, cfg := range spriteConfigs {
		if err := rm.loadSpriteWithFallback(ctx, cfg); err != nil {
			log.Printf("[WARN] Failed to load ending sprite, using fallback: name=%s error=%v", cfg.Name, err)
			// Continue with fallback
		}
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

	// Load UI sprites (title logo, menu frame, etc.)
	if err := rm.loadUISprites(ctx); err != nil {
		log.Printf("[WARN] Failed to load some UI sprites, continuing: error=%v", err)
	}

	// Load cutscene sprites
	if err := rm.loadCutsceneSprites(ctx); err != nil {
		log.Printf("[WARN] Failed to load some cutscene sprites, continuing: error=%v", err)
	}

	// Load planet sprites
	if err := rm.loadPlanetSprites(ctx); err != nil {
		log.Printf("[WARN] Failed to load some planet sprites, continuing: error=%v", err)
	}

	// Load boss portrait sprites
	if err := rm.loadBossPortraitSprites(ctx); err != nil {
		log.Printf("[WARN] Failed to load some boss portrait sprites, continuing: error=%v", err)
	}

	// Load ending sprites
	if err := rm.loadEndingSprites(ctx); err != nil {
		log.Printf("[WARN] Failed to load some ending sprites, continuing: error=%v", err)
	}

	return nil
}
