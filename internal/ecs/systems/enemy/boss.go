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
	bossSprite = es.ensureBossSprite(bossSprite)

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
	es.verifyBossSprite(ctx, entry, entity)

	// Setup boss entity components
	if err := es.setupBossEntity(entry, centerX, centerY); err != nil {
		es.logger.Error("[BOSS_SPRITE] Failed to setup boss entity", "error", err)
		return donburi.Null
	}

	es.logBossSpawned(entry, spawnX, spawnY, initialAngle, entity)
	return entity
}

// ensureBossSprite ensures boss sprite is not nil
func (es *EnemySystem) ensureBossSprite(bossSprite *ebiten.Image) *ebiten.Image {
	if bossSprite != nil {
		return bossSprite
	}

	es.logger.Error("[BOSS_SPRITE] CRITICAL: Boss sprite is nil in SpawnBoss, creating emergency fallback")
	bossData, err := es.GetEnemyTypeData(EnemyTypeBoss)
	if err != nil {
		es.logger.Error("[BOSS_SPRITE] Failed to get boss data for fallback", "error", err)
		bossSprite = ebiten.NewImage(64, 64) // Default size
	} else {
		bossSprite = ebiten.NewImage(bossData.Size, bossData.Size)
	}
	bossSprite.Fill(color.RGBA{128, 0, 128, 255}) // Purple
	return bossSprite
}

// logBossSpawned logs boss spawn information
func (es *EnemySystem) logBossSpawned(
	entry *donburi.Entry,
	spawnX, spawnY, initialAngle float64,
	entity donburi.Entity,
) {
	finalSprite := core.Sprite.Get(entry)
	spriteInfo := "nil"
	if finalSprite != nil && *finalSprite != nil {
		bounds := (*finalSprite).Bounds()
		spriteInfo = fmt.Sprintf("%dx%d", bounds.Dx(), bounds.Dy())
	}

	bossData, err := es.GetEnemyTypeData(EnemyTypeBoss)
	if err != nil {
		es.logger.Error("[BOSS_SPRITE] Failed to get boss data for logging", "error", err)
		return
	}
	bossHealth := bossData.Health
	if es.bossConfig != nil && es.bossConfig.Health > 0 {
		bossHealth = es.bossConfig.Health
	}

	es.logger.Debug("Boss spawned",
		"type", EnemyTypeBoss.String(),
		"sprite_name", "enemy_boss",
		"sprite_size", spriteInfo,
		"health", bossHealth,
		"position", common.Point{X: spawnX, Y: spawnY},
		"angle", initialAngle,
		"entity", entity)
}

// verifyBossSprite verifies and logs sprite details
func (es *EnemySystem) verifyBossSprite(ctx context.Context, entry *donburi.Entry, entity donburi.Entity) {
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
}

// setupBossEntity sets up boss entity components
func (es *EnemySystem) setupBossEntity(entry *donburi.Entry, centerX, centerY float64) error {
	// Use boss config if available, otherwise use loaded boss data
	bossData, err := es.GetEnemyTypeData(EnemyTypeBoss)
	if err != nil {
		return fmt.Errorf("failed to get boss type data: %w", err)
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

	return nil
}

// getBossSprite loads or creates the boss sprite (with caching)
func (es *EnemySystem) getBossSprite(ctx context.Context) *ebiten.Image {
	// Check cache first (silent - this happens frequently)
	if sprite, ok := es.enemySprites[EnemyTypeBoss]; ok {
		return sprite
	}

	// Try to load boss sprite from resource manager
	bossSprite := es.tryLoadBossSprite(ctx)
	bossSprite = es.verifyBossSpriteNotPlayer(ctx, bossSprite)

	// Verify sprite is not nil before caching
	if bossSprite == nil {
		bossSprite = es.createBossPlaceholder()
	}

	// Cache the sprite (log only once during initialization)
	es.enemySprites[EnemyTypeBoss] = bossSprite
	es.logger.Info("[BOSS_SPRITE] Boss sprite loaded and cached", "size", bossSprite.Bounds())
	return bossSprite
}

// tryLoadBossSprite tries to load boss sprite from resource manager
func (es *EnemySystem) tryLoadBossSprite(ctx context.Context) *ebiten.Image {
	bossSprite, exists := es.resourceMgr.GetSprite(ctx, "enemy_boss")
	if !exists || bossSprite == nil {
		// Try alternative name
		bossSprite, exists = es.resourceMgr.GetSprite(ctx, "boss")
		if !exists || bossSprite == nil {
			es.logger.Warn("[BOSS_SPRITE] Boss sprite not found in resource manager")
			return nil
		}
		return bossSprite
	}
	return bossSprite
}

// verifyBossSpriteNotPlayer verifies that boss sprite is not the same as player sprite
func (es *EnemySystem) verifyBossSpriteNotPlayer(ctx context.Context, bossSprite *ebiten.Image) *ebiten.Image {
	if bossSprite == nil {
		return nil
	}

	playerSprite, playerExists := es.resourceMgr.GetSprite(ctx, "player")
	if !playerExists || playerSprite == nil {
		return bossSprite
	}

	// Check if they're the same pointer (same sprite object)
	if bossSprite == playerSprite {
		es.logger.Error(
			"[BOSS_SPRITE] CRITICAL ERROR: Boss sprite is the same as player sprite",
		)
		// Force create a proper boss sprite
		return es.createBossPlaceholder()
	}

	// Verification passed (no need to log success on every check)
	return bossSprite
}

// createBossPlaceholder creates a fallback placeholder sprite
func (es *EnemySystem) createBossPlaceholder() *ebiten.Image {
	es.logger.Error("[BOSS_SPRITE] Boss sprite is nil after all attempts, creating fallback")
	bossData, err := es.GetEnemyTypeData(EnemyTypeBoss)
	if err != nil {
		es.logger.Error("[BOSS_SPRITE] Failed to get boss data for fallback", "error", err)
		return es.createDefaultBossSprite(64)
	}
	return es.createDefaultBossSprite(bossData.Size)
}

// createDefaultBossSprite creates a default purple boss sprite
func (es *EnemySystem) createDefaultBossSprite(size int) *ebiten.Image {
	sprite := ebiten.NewImage(size, size)
	sprite.Fill(color.RGBA{128, 0, 128, 255}) // Purple
	// Placeholder creation only logged at Warn level in createBossPlaceholder
	return sprite
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
