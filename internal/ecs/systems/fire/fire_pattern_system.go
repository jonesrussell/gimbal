package fire

import (
	"context"
	"image/color"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/config"
	"github.com/jonesrussell/gimbal/internal/ecs/core"
)

const (
	DefaultProjectileSpeed = 5.0
	DefaultProjectileSize  = 6
)

// FirePatternSystem handles enemy firing patterns
type FirePatternSystem struct {
	world            donburi.World
	config           *config.GameConfig
	screenCenter     common.Point
	projectileSprite *ebiten.Image
}

// NewFirePatternSystem creates a new fire pattern system
func NewFirePatternSystem(
	world donburi.World,
	cfg *config.GameConfig,
) *FirePatternSystem {
	fps := &FirePatternSystem{
		world:  world,
		config: cfg,
		screenCenter: common.Point{
			X: float64(cfg.ScreenSize.Width) / 2,
			Y: float64(cfg.ScreenSize.Height) / 2,
		},
	}

	// Create projectile sprite
	fps.createProjectileSprite()

	return fps
}

// createProjectileSprite creates a default enemy projectile sprite
func (fps *FirePatternSystem) createProjectileSprite() {
	fps.projectileSprite = ebiten.NewImage(DefaultProjectileSize, DefaultProjectileSize)
	fps.projectileSprite.Fill(color.RGBA{R: 255, G: 100, B: 100, A: 255}) // Red-ish
}

// Update processes all entities with fire patterns
func (fps *FirePatternSystem) Update(ctx context.Context, deltaTime float64) error {
	// Check for cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Get player position for targeting
	playerPos := fps.getPlayerPosition()

	// Query entities with FirePattern component
	query.NewQuery(
		filter.And(
			filter.Contains(core.FirePattern),
			filter.Contains(core.Position),
			filter.Contains(core.BehaviorState),
		),
	).Each(fps.world, func(entry *donburi.Entry) {
		fps.updateFirePattern(entry, playerPos, deltaTime)
	})

	// Update existing enemy projectiles
	fps.updateProjectiles(deltaTime)

	return nil
}

// updateFirePattern processes an entity's fire pattern
func (fps *FirePatternSystem) updateFirePattern(entry *donburi.Entry, playerPos common.Point, deltaTime float64) {
	fireData := core.FirePattern.Get(entry)
	behaviorData := core.BehaviorState.Get(entry)

	// Check if firing is allowed in current state
	canFire := fps.canFireInState(fireData, behaviorData)
	if !canFire {
		return
	}

	// Update fire cooldown
	fireData.LastFireTime += time.Duration(deltaTime * float64(time.Second))

	// Check if ready to fire
	if fireData.FireRate <= 0 {
		return
	}

	fireCooldown := time.Duration(float64(time.Second) / fireData.FireRate)
	if fireData.LastFireTime < fireCooldown {
		core.FirePattern.SetValue(entry, *fireData)
		return
	}

	// Execute fire pattern
	pos := core.Position.Get(entry)
	fps.executeFirePattern(pos, playerPos, fireData)

	// Reset fire timer
	fireData.LastFireTime = 0
	core.FirePattern.SetValue(entry, *fireData)
}

// canFireInState checks if the entity can fire in its current behavior state
func (fps *FirePatternSystem) canFireInState(fireData *core.FirePatternData, behaviorData *core.BehaviorStateData) bool {
	switch behaviorData.CurrentState {
	case core.StateOrbiting:
		return fireData.CanFireWhileOrbit
	case core.StateAttacking:
		return fireData.CanFireWhileAttack
	case core.StateHovering:
		return true // Can fire while hovering
	default:
		return false
	}
}

// executeFirePattern executes the appropriate fire pattern
func (fps *FirePatternSystem) executeFirePattern(enemyPos *common.Point, playerPos common.Point, fireData *core.FirePatternData) {
	switch fireData.PatternType {
	case core.FireSingleShot:
		fps.fireSingleShot(enemyPos, playerPos)
	case core.FireBurst:
		fps.fireBurst(enemyPos, playerPos, fireData)
	case core.FireSpray:
		fps.fireSpray(enemyPos, playerPos, fireData)
	}
}

// fireSingleShot fires a single aimed shot at the player
func (fps *FirePatternSystem) fireSingleShot(enemyPos *common.Point, playerPos common.Point) {
	// Calculate direction to player
	dx := playerPos.X - enemyPos.X
	dy := playerPos.Y - enemyPos.Y
	distance := math.Sqrt(dx*dx + dy*dy)

	if distance < 1 {
		return
	}

	// Normalize
	dx /= distance
	dy /= distance

	// Create projectile
	fps.createProjectile(*enemyPos, dx, dy, DefaultProjectileSpeed)
}

// fireBurst fires a burst of shots
func (fps *FirePatternSystem) fireBurst(enemyPos *common.Point, playerPos common.Point, fireData *core.FirePatternData) {
	burstCount := fireData.BurstCount
	if burstCount <= 0 {
		burstCount = 3
	}

	// Calculate base direction to player
	dx := playerPos.X - enemyPos.X
	dy := playerPos.Y - enemyPos.Y
	distance := math.Sqrt(dx*dx + dy*dy)

	if distance < 1 {
		return
	}

	baseAngle := math.Atan2(dy, dx)

	// Fire burst with slight spread
	spreadAngle := 0.1 // radians

	for i := 0; i < burstCount; i++ {
		// Calculate offset from base angle
		offset := float64(i-burstCount/2) * spreadAngle
		angle := baseAngle + offset

		dirX := math.Cos(angle)
		dirY := math.Sin(angle)

		fps.createProjectile(*enemyPos, dirX, dirY, DefaultProjectileSpeed)
	}
}

// fireSpray fires multiple projectiles in a spread pattern
func (fps *FirePatternSystem) fireSpray(enemyPos *common.Point, playerPos common.Point, fireData *core.FirePatternData) {
	projectileCount := fireData.ProjectileCount
	if projectileCount <= 0 {
		projectileCount = 5
	}

	sprayAngle := fireData.SprayAngle
	if sprayAngle <= 0 {
		sprayAngle = 60.0 // degrees
	}
	sprayRad := sprayAngle * math.Pi / 180

	// Calculate base direction to player
	dx := playerPos.X - enemyPos.X
	dy := playerPos.Y - enemyPos.Y
	distance := math.Sqrt(dx*dx + dy*dy)

	if distance < 1 {
		return
	}

	baseAngle := math.Atan2(dy, dx)

	// Fire spray
	angleStep := sprayRad / float64(projectileCount-1)
	startAngle := baseAngle - sprayRad/2

	for i := 0; i < projectileCount; i++ {
		angle := startAngle + float64(i)*angleStep
		dirX := math.Cos(angle)
		dirY := math.Sin(angle)

		fps.createProjectile(*enemyPos, dirX, dirY, DefaultProjectileSpeed*0.8) // Slightly slower
	}
}

// createProjectile creates a new enemy projectile entity
func (fps *FirePatternSystem) createProjectile(start common.Point, dirX, dirY, speed float64) {
	entity := fps.world.Create(
		core.EnemyProjectileTag,
		core.Position,
		core.Sprite,
		core.Movement,
		core.Size,
	)
	entry := fps.world.Entry(entity)

	// Set position
	core.Position.SetValue(entry, start)

	// Set sprite
	core.Sprite.SetValue(entry, fps.projectileSprite)

	// Set size
	core.Size.SetValue(entry, config.Size{
		Width:  DefaultProjectileSize,
		Height: DefaultProjectileSize,
	})

	// Set velocity (direction * speed)
	velocity := common.Point{
		X: dirX * speed,
		Y: dirY * speed,
	}
	core.Movement.SetValue(entry, core.MovementData{
		Velocity: velocity,
		MaxSpeed: speed,
	})
}

// updateProjectiles updates all enemy projectiles
func (fps *FirePatternSystem) updateProjectiles(deltaTime float64) {
	query.NewQuery(
		filter.And(
			filter.Contains(core.EnemyProjectileTag),
			filter.Contains(core.Position),
			filter.Contains(core.Movement),
		),
	).Each(fps.world, func(entry *donburi.Entry) {
		pos := core.Position.Get(entry)
		mov := core.Movement.Get(entry)

		// Update position
		pos.X += mov.Velocity.X
		pos.Y += mov.Velocity.Y

		// Remove if off-screen
		margin := 50.0
		if pos.X < -margin || pos.X > float64(fps.config.ScreenSize.Width)+margin ||
			pos.Y < -margin || pos.Y > float64(fps.config.ScreenSize.Height)+margin {
			fps.world.Remove(entry.Entity())
		}
	})
}

// getPlayerPosition finds the player position
func (fps *FirePatternSystem) getPlayerPosition() common.Point {
	var playerPos common.Point

	query.NewQuery(
		filter.Contains(core.PlayerTag),
	).Each(fps.world, func(entry *donburi.Entry) {
		if entry.HasComponent(core.Position) {
			pos := core.Position.Get(entry)
			playerPos = *pos
		}
	})

	// Default to center if no player found
	if playerPos.X == 0 && playerPos.Y == 0 {
		playerPos = fps.screenCenter
	}

	return playerPos
}
