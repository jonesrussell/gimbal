package enemy

import (
	"context"
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
	bossSprite := es.getBossSprite(ctx)

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

	// Set size (boss is larger)
	bossData := GetEnemyTypeData(EnemyTypeBoss)
	core.Size.SetValue(entry, config.Size{Width: bossData.Size, Height: bossData.Size})

	// Set health
	core.Health.SetValue(entry, core.NewHealthData(bossData.Health, bossData.Health))

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

	es.logger.Debug("Enemy spawned",
		"type", EnemyTypeBoss.String(),
		"sprite", "enemy_boss",
		"health", bossData.Health,
		"position", common.Point{X: spawnX, Y: spawnY},
		"angle", initialAngle)

	return entity
}

// getBossSprite loads or creates the boss sprite (with caching)
func (es *EnemySystem) getBossSprite(ctx context.Context) *ebiten.Image {
	// Check cache first
	if sprite, ok := es.enemySprites[EnemyTypeBoss]; ok {
		return sprite
	}

	// Try to load boss sprite from resource manager
	bossSprite, exists := es.resourceMgr.GetSprite(ctx, "enemy_boss")
	if !exists {
		// Try alternative name
		bossSprite, exists = es.resourceMgr.GetSprite(ctx, "boss")
		if !exists {
			es.logger.Warn("Boss sprite not found, using placeholder")
			// Create a larger placeholder sprite (purple to distinguish from regular enemies)
			bossData := GetEnemyTypeData(EnemyTypeBoss)
			bossSprite = ebiten.NewImage(bossData.Size, bossData.Size)
			bossSprite.Fill(color.RGBA{128, 0, 128, 255}) // Purple
		}
	}

	// Cache the sprite
	es.enemySprites[EnemyTypeBoss] = bossSprite
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
