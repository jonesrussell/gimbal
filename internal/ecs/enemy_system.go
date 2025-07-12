package ecs

import (
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
	"github.com/yohamta/ganim8/v2"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/ecs/core"
	"github.com/jonesrussell/gimbal/internal/ecs/resources"
)

// EnemySystem manages enemy spawning, movement, and behavior
// Enhanced: uses ganim8 for sprite animations
type EnemySystem struct {
	world         donburi.World
	config        *common.GameConfig
	spawnTimer    float64
	spawnInterval float64
	resourceMgr   *resources.ResourceManager

	// Animation assets
	enemySheet       *ebiten.Image
	idleAnimation    *ganim8.Animation
	explodeAnimation *ganim8.Animation
	grid             *ganim8.Grid
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

	// Initialize animations
	es.initializeAnimations()

	return es
}

// initializeAnimations sets up ganim8 animations for enemies
func (es *EnemySystem) initializeAnimations() {
	// Load enemy sprite sheet
	enemySheet, exists := es.resourceMgr.GetSprite("enemy_sheet")
	if !exists {
		// Create a placeholder sprite sheet for testing
		enemySheet = ebiten.NewImage(512, 256)
		enemySheet.Fill(color.RGBA{255, 0, 0, 255})
	}
	es.enemySheet = enemySheet

	// Create grid for 4x2 sprite sheet (128x128 per frame)
	es.grid = ganim8.NewGrid(128, 128, 512, 256)

	// Create idle animation (frames 1-2, row 1)
	es.idleAnimation = ganim8.New(es.enemySheet, es.grid.Frames("1-2", 1), 500*time.Millisecond)

	// Create explosion animation (frames 7-8, row 2)
	es.explodeAnimation = ganim8.New(es.enemySheet, es.grid.Frames("7-8", 2), 200*time.Millisecond)
}

func (es *EnemySystem) Update(deltaTime float64) {
	// Update animations
	es.idleAnimation.Update()
	es.explodeAnimation.Update()

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

	entity := es.world.Create(
		core.EnemyTag, core.Position, core.Sprite, core.Movement,
		core.Size, core.Health, core.Animation,
	)
	entry := es.world.Entry(entity)
	core.Position.SetValue(entry, spawnPos)

	// Set sprite to the sprite sheet
	core.Sprite.SetValue(entry, es.enemySheet)
	core.Size.SetValue(entry, common.Size{Width: 32, Height: 32})
	core.Health.SetValue(entry, core.HealthData{Current: 1, Maximum: 1, InvincibilityDuration: 0})

	// Set up animation data
	animationData := core.AnimationData{
		CurrentAnimation: es.idleAnimation,
		State:            core.EnemyStateIdle,
		IdleAnimation:    es.idleAnimation,
		ExplodeAnimation: es.explodeAnimation,
	}
	core.Animation.SetValue(entry, animationData)

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

		// Check if enemy is exploding and animation is complete
		if animation := core.Animation.Get(entry); animation != nil {
			if animation.State == core.EnemyStateExploding {
				// Check if explosion animation is finished
				// Note: ganim8 doesn't have a Finished() method, so we'll use a timer-based approach
				// For now, we'll remove the enemy after a fixed time (400ms for 2 frames at 200ms each)
				// In a more sophisticated implementation, you'd track the animation state
				es.world.Remove(entry.Entity())
				return
			}
		}

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

// TriggerExplosion starts the explosion animation for an enemy
func (es *EnemySystem) TriggerExplosion(entity donburi.Entity) {
	entry := es.world.Entry(entity)
	if !entry.Valid() {
		return
	}

	// Update animation state
	if animation := core.Animation.Get(entry); animation != nil {
		animation.State = core.EnemyStateExploding
		animation.CurrentAnimation = animation.ExplodeAnimation
		// Start explosion from beginning by creating a new animation instance
		animation.CurrentAnimation = ganim8.New(es.enemySheet, es.grid.Frames("7-8", 2), 200*time.Millisecond)
	}
}
