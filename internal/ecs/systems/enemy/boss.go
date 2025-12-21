package enemy

import (
	"context"
	"fmt"
	stdmath "math"

	"github.com/yohamta/donburi"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/config"
	"github.com/jonesrussell/gimbal/internal/ecs/core"
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
