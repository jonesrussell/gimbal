package enemy

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/config"
	"github.com/jonesrussell/gimbal/internal/ecs/core"
)

// Enemy weapon constants
const (
	EnemyProjectileSize   = 6    // Size of enemy projectiles
	EnemyProjectileMargin = 50.0 // Screen margin for cleanup
)

// EnemyWeaponSystem manages enemy shooting and projectiles
type EnemyWeaponSystem struct {
	world            donburi.World
	config           *config.GameConfig
	logger           common.Logger
	projectileSprite *ebiten.Image
	enemySystem      *EnemySystem // Reference to enemy system for getting enemy type data

	// Track fire timers per enemy entity
	enemyFireTimers map[donburi.Entity]float64
}

// NewEnemyWeaponSystem creates a new enemy weapon system
func NewEnemyWeaponSystem(
	world donburi.World,
	gameConfig *config.GameConfig,
	logger common.Logger,
	enemySystem *EnemySystem,
) *EnemyWeaponSystem {
	ews := &EnemyWeaponSystem{
		world:           world,
		config:          gameConfig,
		logger:          logger,
		enemySystem:     enemySystem,
		enemyFireTimers: make(map[donburi.Entity]float64),
	}
	ews.initializeProjectileSprite()
	return ews
}

// initializeProjectileSprite creates the enemy projectile sprite
func (ews *EnemyWeaponSystem) initializeProjectileSprite() {
	// Red projectile for enemies
	ews.projectileSprite = ebiten.NewImage(EnemyProjectileSize, EnemyProjectileSize)
	ews.projectileSprite.Fill(color.RGBA{R: 255, G: 50, B: 50, A: 255})
}

// Update updates enemy weapons and their projectiles
func (ews *EnemyWeaponSystem) Update(deltaTime float64) {
	// Get player position for targeting
	playerPos := ews.getPlayerPosition()
	if playerPos == nil {
		// No player to shoot at
		return
	}

	// Update fire timers and shoot for each enemy
	ews.updateEnemyShooting(deltaTime, *playerPos)

	// Update existing enemy projectiles
	ews.updateProjectiles()

	// Clean up fire timers for removed enemies
	ews.cleanupFireTimers()
}

// getPlayerPosition returns the player's current position
func (ews *EnemyWeaponSystem) getPlayerPosition() *common.Point {
	var playerPos *common.Point
	query.NewQuery(
		filter.And(
			filter.Contains(core.PlayerTag),
			filter.Contains(core.Position),
		),
	).Each(ews.world, func(entry *donburi.Entry) {
		pos := core.Position.Get(entry)
		playerPos = pos
	})
	return playerPos
}

// updateEnemyShooting handles shooting logic for all enemies
func (ews *EnemyWeaponSystem) updateEnemyShooting(deltaTime float64, playerPos common.Point) {
	query.NewQuery(
		filter.And(
			filter.Contains(core.EnemyTag),
			filter.Contains(core.Position),
		),
	).Each(ews.world, func(entry *donburi.Entry) {
		entity := entry.Entity()
		enemyPos := core.Position.Get(entry)

		// Get enemy type from component (preferred) or fall back to health heuristic
		enemyData, err := ews.getEnemyTypeDataForWeapon(entry)
		if err != nil {
			ews.logger.Warn("Failed to get enemy type data for weapon", "error", err)
			return // Skip this enemy
		}

		// Skip if enemy can't shoot
		if !enemyData.CanShoot || enemyData.FireRate <= 0 {
			return
		}

		// Initialize fire timer if needed
		if _, exists := ews.enemyFireTimers[entity]; !exists {
			// Randomize initial timer so enemies don't all shoot at once
			ews.enemyFireTimers[entity] = float64(entity%100) / 100.0 / enemyData.FireRate
		}

		// Update fire timer
		ews.enemyFireTimers[entity] += deltaTime

		// Check if it's time to fire
		fireInterval := 1.0 / enemyData.FireRate
		if ews.enemyFireTimers[entity] >= fireInterval {
			ews.fireAtPlayer(*enemyPos, playerPos, enemyData.ProjectileSpeed)
			ews.enemyFireTimers[entity] = 0
		}
	})
}

// getEnemyTypeDataForWeapon gets enemy type data from entry
func (ews *EnemyWeaponSystem) getEnemyTypeDataForWeapon(entry *donburi.Entry) (EnemyTypeData, error) {
	if entry.HasComponent(core.EnemyTypeID) {
		typeID := core.EnemyTypeID.Get(entry)
		return ews.enemySystem.GetEnemyTypeData(EnemyType(*typeID))
	}
	if entry.HasComponent(core.Health) {
		// Fallback for legacy entities without EnemyTypeID
		health := core.Health.Get(entry)
		if health.Maximum >= 10 {
			return ews.enemySystem.GetEnemyTypeData(EnemyTypeBoss)
		}
		if health.Maximum >= 2 {
			return ews.enemySystem.GetEnemyTypeData(EnemyTypeHeavy)
		}
		return ews.enemySystem.GetEnemyTypeData(EnemyTypeBasic)
	}
	return ews.enemySystem.GetEnemyTypeData(EnemyTypeBasic)
}

// fireAtPlayer creates a projectile aimed at the player
func (ews *EnemyWeaponSystem) fireAtPlayer(enemyPos, playerPos common.Point, speed float64) {
	// Calculate direction from enemy to player
	dirX := playerPos.X - enemyPos.X
	dirY := playerPos.Y - enemyPos.Y

	// Normalize direction
	distance := math.Sqrt(dirX*dirX + dirY*dirY)
	if distance < 1 {
		// Enemy is basically on top of player, skip
		return
	}
	dirX /= distance
	dirY /= distance

	// Create projectile entity
	entity := ews.world.Create(
		core.EnemyProjectileTag,
		core.Position,
		core.Sprite,
		core.Movement,
		core.Size,
	)
	entry := ews.world.Entry(entity)

	// Set position (at enemy location)
	core.Position.SetValue(entry, enemyPos)

	// Set size
	core.Size.SetValue(entry, config.Size{Width: EnemyProjectileSize, Height: EnemyProjectileSize})

	// Set velocity toward player
	velocity := common.Point{
		X: dirX * speed,
		Y: dirY * speed,
	}
	core.Movement.SetValue(entry, core.MovementData{
		Velocity: velocity,
		MaxSpeed: speed,
	})

	// Set sprite
	core.Sprite.SetValue(entry, ews.projectileSprite)

	ews.logger.Debug("Enemy fired projectile", "from", enemyPos, "toward", playerPos)
}

// updateProjectiles updates all enemy projectile positions
func (ews *EnemyWeaponSystem) updateProjectiles() {
	query.NewQuery(
		filter.And(
			filter.Contains(core.EnemyProjectileTag),
			filter.Contains(core.Position),
			filter.Contains(core.Movement),
		),
	).Each(ews.world, func(entry *donburi.Entry) {
		pos := core.Position.Get(entry)
		mov := core.Movement.Get(entry)

		// Update position based on velocity
		pos.X += mov.Velocity.X
		pos.Y += mov.Velocity.Y

		// Check if projectile is off screen
		if ews.isOffScreen(*pos) {
			ews.world.Remove(entry.Entity())
		}
	})
}

// isOffScreen checks if a position is off screen
func (ews *EnemyWeaponSystem) isOffScreen(pos common.Point) bool {
	margin := EnemyProjectileMargin
	return pos.X < -margin || pos.X > float64(ews.config.ScreenSize.Width)+margin ||
		pos.Y < -margin || pos.Y > float64(ews.config.ScreenSize.Height)+margin
}

// cleanupFireTimers removes fire timers for entities that no longer exist
func (ews *EnemyWeaponSystem) cleanupFireTimers() {
	for entity := range ews.enemyFireTimers {
		entry := ews.world.Entry(entity)
		if !entry.Valid() {
			delete(ews.enemyFireTimers, entity)
		}
	}
}
