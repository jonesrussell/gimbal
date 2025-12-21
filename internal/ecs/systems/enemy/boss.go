package enemy

import (
	"context"
	"fmt"
	"image/color"
	stdmath "math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/config"
	"github.com/jonesrussell/gimbal/internal/ecs/core"
	"github.com/jonesrussell/gimbal/internal/math"
)

// BossSpawnRadius is the orbital radius for the boss
const BossSpawnRadius = 150.0

// BossOrbitalSpeed is how fast the boss orbits
const BossOrbitalSpeed = 0.5 // radians per second

// SpawnBoss spawns a boss enemy with orbital movement
func (es *EnemySystem) SpawnBoss(ctx context.Context) donburi.Entity {
	// Load boss sprite if not already loaded
	// IMPORTANT: Use getBossSprite, not getEnemySprite, to ensure correct sprite
	bossSprite := es.getBossSprite(ctx)

	// Final verification: ensure sprite is not nil
	if bossSprite == nil {
		es.logger.Error("[BOSS_SPRITE] CRITICAL: Boss sprite is nil in SpawnBoss, creating emergency fallback")
		bossData, err := es.GetEnemyTypeData(EnemyTypeBoss)
		if err != nil {
			es.logger.Error("[BOSS_SPRITE] Failed to get boss data for fallback", "error", err)
			bossSprite = ebiten.NewImage(64, 64) // Default size
		} else {
			bossSprite = ebiten.NewImage(bossData.Size, bossData.Size)
		}
		bossSprite.Fill(color.RGBA{128, 0, 128, 255}) // Purple
	}

	centerX := float64(es.gameConfig.ScreenSize.Width) / 2
	centerY := float64(es.gameConfig.ScreenSize.Height) / 2

	// Spawn boss at top of orbital path (270 degrees = top)
	initialAngle := 270.0 * stdmath.Pi / 180.0
	spawnX := centerX + stdmath.Cos(initialAngle)*BossSpawnRadius
	spawnY := centerY + stdmath.Sin(initialAngle)*BossSpawnRadius

	entity := es.world.Create(
		core.EnemyTag, core.Position, core.Sprite, core.Orbital,
		core.Size, core.Health, core.Angle, core.EnemyTypeID,
	)
	entry := es.world.Entry(entity)

	// Set position
	core.Position.SetValue(entry, common.Point{X: spawnX, Y: spawnY})

	// Set sprite
	core.Sprite.SetValue(entry, bossSprite)

	// Verify sprite was set correctly and log pointer address for debugging
	if setSprite := core.Sprite.Get(entry); setSprite != nil && *setSprite != nil {
		bounds := (*setSprite).Bounds()
		// Get player sprite for comparison
		playerSprite, _ := es.resourceMgr.GetSprite(ctx, "player")
		spriteMatch := "different"
		if playerSprite != nil && *setSprite == playerSprite {
			spriteMatch = "SAME AS PLAYER - ERROR!"
		}
		es.logger.Debug("[BOSS_SPRITE] Boss sprite set successfully",
			"sprite_size", fmt.Sprintf("%dx%d", bounds.Dx(), bounds.Dy()),
			"sprite_ptr", fmt.Sprintf("%p", *setSprite),
			"player_ptr", fmt.Sprintf("%p", playerSprite),
			"comparison", spriteMatch,
			"entity", entity)
	} else {
		es.logger.Error("[BOSS_SPRITE] Failed to set boss sprite - sprite is nil")
	}

	// Use boss config if available, otherwise use loaded boss data
	bossData, err := es.GetEnemyTypeData(EnemyTypeBoss)
	if err != nil {
		es.logger.Error("[BOSS_SPRITE] Failed to get boss type data", "error", err)
		return donburi.Null
	}
	bossSize := bossData.Size
	bossHealth := bossData.Health

	if es.bossConfig != nil {
		if es.bossConfig.Size > 0 {
			bossSize = es.bossConfig.Size
		}
		if es.bossConfig.Health > 0 {
			bossHealth = es.bossConfig.Health
		}
	}

	// Set size (boss is larger)
	core.Size.SetValue(entry, config.Size{Width: bossSize, Height: bossSize})

	// Set health
	core.Health.SetValue(entry, core.NewHealthData(bossHealth, bossHealth))

	// Set enemy type ID for proper identification
	core.EnemyTypeID.SetValue(entry, int(EnemyTypeBoss))

	// Set up orbital movement
	orbitalData := core.OrbitalData{
		Center:       common.Point{X: centerX, Y: centerY},
		Radius:       BossSpawnRadius,
		OrbitalAngle: 270, // Start at top
		FacingAngle:  0,
	}
	core.Orbital.SetValue(entry, orbitalData)

	// Set initial angle
	core.Angle.SetValue(entry, 0)

	// Verify final sprite
	finalSprite := core.Sprite.Get(entry)
	spriteInfo := "nil"
	if finalSprite != nil && *finalSprite != nil {
		bounds := (*finalSprite).Bounds()
		spriteInfo = fmt.Sprintf("%dx%d", bounds.Dx(), bounds.Dy())
	}

	es.logger.Debug("Boss spawned",
		"type", EnemyTypeBoss.String(),
		"sprite_name", "enemy_boss",
		"sprite_size", spriteInfo,
		"health", bossHealth,
		"position", common.Point{X: spawnX, Y: spawnY},
		"angle", initialAngle,
		"entity", entity)

	return entity
}

// getBossSprite loads or creates the boss sprite (with caching)
func (es *EnemySystem) getBossSprite(ctx context.Context) *ebiten.Image {
	// Check cache first
	if sprite, ok := es.enemySprites[EnemyTypeBoss]; ok {
		es.logger.Debug("[BOSS_SPRITE] Using cached boss sprite", "sprite_size", sprite.Bounds())
		return sprite
	}

	// Try to load boss sprite from resource manager
	es.logger.Debug("[BOSS_SPRITE] Loading boss sprite from resource manager", "sprite_name", "enemy_boss")
	bossSprite, exists := es.resourceMgr.GetSprite(ctx, "enemy_boss")
	if !exists || bossSprite == nil {
		es.logger.Warn("[BOSS_SPRITE] enemy_boss not found or nil", "exists", exists, "sprite_nil", bossSprite == nil)
		// Try alternative name
		bossSprite, exists = es.resourceMgr.GetSprite(ctx, "boss")
		if !exists || bossSprite == nil {
			es.logger.Warn("[BOSS_SPRITE] Boss sprite not found in resource manager, creating placeholder")
			// Create a larger placeholder sprite (purple to distinguish from regular enemies)
			bossData, err := es.GetEnemyTypeData(EnemyTypeBoss)
			if err != nil {
				es.logger.Error("[BOSS_SPRITE] Failed to get boss data for placeholder", "error", err)
				bossSprite = ebiten.NewImage(64, 64) // Default size
			} else {
				bossSprite = ebiten.NewImage(bossData.Size, bossData.Size)
			}
			bossSprite.Fill(color.RGBA{128, 0, 128, 255}) // Purple
			es.logger.Debug("[BOSS_SPRITE] Created purple placeholder", "size", bossData.Size)
		} else {
			bounds := bossSprite.Bounds()
			es.logger.Debug("[BOSS_SPRITE] Found boss sprite with name 'boss'",
				"sprite_size", fmt.Sprintf("%dx%d", bounds.Dx(), bounds.Dy()))
		}
	} else {
		bounds := bossSprite.Bounds()
		es.logger.Debug("[BOSS_SPRITE] Found boss sprite with name 'enemy_boss'",
			"sprite_size", fmt.Sprintf("%dx%d", bounds.Dx(), bounds.Dy()))

		// CRITICAL: Verify we didn't accidentally get the player sprite
		// Compare with player sprite to ensure they're different
		playerSprite, playerExists := es.resourceMgr.GetSprite(ctx, "player")
		if playerExists && playerSprite != nil {
			// Check if they're the same pointer (same sprite object)
			if bossSprite == playerSprite {
				es.logger.Error("[BOSS_SPRITE] CRITICAL ERROR: Boss sprite is the same as player sprite! This is wrong!")
				// Force create a proper boss sprite
				bossData, err := es.GetEnemyTypeData(EnemyTypeBoss)
				if err != nil {
					es.logger.Error("[BOSS_SPRITE] Failed to get boss data for replacement", "error", err)
					bossSprite = ebiten.NewImage(64, 64) // Default size
				} else {
					bossSprite = ebiten.NewImage(bossData.Size, bossData.Size)
				}
				bossSprite.Fill(color.RGBA{128, 0, 128, 255}) // Purple
				es.logger.Warn("[BOSS_SPRITE] Created replacement boss sprite to fix player sprite issue", "size", bossData.Size)
			} else {
				es.logger.Debug("[BOSS_SPRITE] Verified boss sprite is different from player sprite")
			}
		}
	}

	// Verify sprite is not nil before caching
	if bossSprite == nil {
		es.logger.Error("[BOSS_SPRITE] Boss sprite is nil after all attempts, creating fallback")
		bossData, err := es.GetEnemyTypeData(EnemyTypeBoss)
		if err != nil {
			es.logger.Error("[BOSS_SPRITE] Failed to get boss data for fallback", "error", err)
			bossSprite = ebiten.NewImage(64, 64) // Default size
		} else {
			bossSprite = ebiten.NewImage(bossData.Size, bossData.Size)
		}
		bossSprite.Fill(color.RGBA{128, 0, 128, 255}) // Purple
	}

	// Cache the sprite
	es.enemySprites[EnemyTypeBoss] = bossSprite
	es.logger.Debug("[BOSS_SPRITE] Boss sprite cached", "final_size", bossSprite.Bounds())
	return bossSprite
}

// UpdateBossMovement updates the boss's orbital movement
func (es *EnemySystem) UpdateBossMovement(deltaTime float64) {
	query.NewQuery(
		filter.And(
			filter.Contains(core.EnemyTag),
			filter.Contains(core.Orbital),
		),
	).Each(es.world, func(entry *donburi.Entry) {
		orbital := core.Orbital.Get(entry)
		health := core.Health.Get(entry)

		// Check if this is a boss (has orbital movement and high health)
		if health != nil && health.Maximum >= 10 {
			// Update orbital angle (convert speed from radians/sec to degrees/sec)
			angleDelta := math.Angle(BossOrbitalSpeed * deltaTime * 180.0 / stdmath.Pi)
			orbital.OrbitalAngle += angleDelta

			// Keep angle in 0-360 range
			orbital.OrbitalAngle = orbital.OrbitalAngle.Normalize()

			// Update position based on orbital angle
			radians := orbital.OrbitalAngle.ToRadians()
			pos := core.Position.Get(entry)
			pos.X = orbital.Center.X + stdmath.Cos(radians)*orbital.Radius
			pos.Y = orbital.Center.Y + stdmath.Sin(radians)*orbital.Radius

			// Update orbital component
			core.Orbital.SetValue(entry, *orbital)
		}
	})
}
