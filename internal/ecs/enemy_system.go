package ecs

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/ecs/core"
	"github.com/jonesrussell/gimbal/internal/ecs/resources"
)

// EnemySystem manages enemy spawning, movement, and behavior
// Simplified: only one enemy type, basic downward movement
type EnemySystem struct {
	world         donburi.World
	config        *common.GameConfig
	spawnTimer    float64
	spawnInterval float64
	resourceMgr   *resources.ResourceManager
}

func NewEnemySystem(
	world donburi.World,
	config *common.GameConfig,
	resourceMgr *resources.ResourceManager,
) *EnemySystem {
	es := &EnemySystem{
		world:         world,
		config:        config,
		spawnTimer:    0,
		spawnInterval: 60, // Spawn every 60 frames (1 second at 60fps)
		resourceMgr:   resourceMgr,
	}

	// Global RNG is automatically seeded in Go 1.20+
	// No need to call rand.Seed() anymore

	return es
}

func (es *EnemySystem) Update(deltaTime float64) {
	es.spawnTimer += deltaTime
	if es.spawnTimer >= es.spawnInterval {
		es.spawnEnemy()
		es.spawnTimer = 0
	}
	es.updateEnemies()
}

func (es *EnemySystem) spawnEnemy() {
	// Spawn at screen center (Gyruss-style)
	centerX := float64(es.config.ScreenSize.Width) / 2
	centerY := float64(es.config.ScreenSize.Height) / 2
	spawnPos := common.Point{X: centerX, Y: centerY}

	// Get enemy sprite from resource manager
	enemySprite, exists := es.resourceMgr.GetSprite("enemy")
	if !exists {
		// Fallback to red square if sprite not found
		enemySprite = ebiten.NewImage(16, 16)
		enemySprite.Fill(color.RGBA{255, 0, 0, 255})
	}

	entity := es.world.Create(core.EnemyTag, core.Position, core.Sprite, core.Movement, core.Size, core.Health)
	entry := es.world.Entry(entity)
	core.Position.SetValue(entry, spawnPos)
	core.Sprite.SetValue(entry, enemySprite)
	core.Size.SetValue(entry, common.Size{Width: 32, Height: 32})
	core.Health.SetValue(entry, core.HealthData{Current: 1, Maximum: 1, InvincibilityDuration: 0})

	// Calculate random angle for outward movement
	angle := rand.Float64() * 2 * math.Pi //nolint:gosec // Game logic randomness is acceptable
	speed := 2.0

	// Move outward from center toward player orbital ring
	velocity := common.Point{
		X: math.Cos(angle) * speed,
		Y: math.Sin(angle) * speed,
	}

	core.Movement.SetValue(entry, core.MovementData{
		Velocity: velocity,
		MaxSpeed: speed,
	})
}

func (es *EnemySystem) updateEnemies() {
	query.NewQuery(
		filter.And(
			filter.Contains(core.EnemyTag),
			filter.Contains(core.Position),
			filter.Contains(core.Movement),
		),
	).Each(es.world, func(entry *donburi.Entry) {
		pos := core.Position.Get(entry)
		mov := core.Movement.Get(entry)
		pos.X += mov.Velocity.X
		pos.Y += mov.Velocity.Y

		// Remove enemies when they move too far from center (Gyruss-style)
		centerX := float64(es.config.ScreenSize.Width) / 2
		centerY := float64(es.config.ScreenSize.Height) / 2
		distanceFromCenter := math.Sqrt((pos.X-centerX)*(pos.X-centerX) + (pos.Y-centerY)*(pos.Y-centerY))
		maxDistance := math.Max(float64(es.config.ScreenSize.Width), float64(es.config.ScreenSize.Height)) * 0.8

		if distanceFromCenter > maxDistance {
			es.world.Remove(entry.Entity())
		}
	})
}
