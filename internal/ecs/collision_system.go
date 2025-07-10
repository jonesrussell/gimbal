package ecs

import (
	"math"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/jonesrussell/gimbal/internal/common"
)

// CollisionSystem manages collision detection and response
type CollisionSystem struct {
	world  donburi.World
	config *common.GameConfig
}

// NewCollisionSystem creates a new collision system
func NewCollisionSystem(world donburi.World, config *common.GameConfig) *CollisionSystem {
	return &CollisionSystem{
		world:  world,
		config: config,
	}
}

// Update updates the collision system
func (cs *CollisionSystem) Update() {
	// Check projectile-enemy collisions
	cs.checkProjectileEnemyCollisions()

	// Check player-enemy collisions
	cs.checkPlayerEnemyCollisions()
}

// checkProjectileEnemyCollisions checks for collisions between projectiles and enemies
func (cs *CollisionSystem) checkProjectileEnemyCollisions() {
	// Get all projectiles
	projectiles := make([]donburi.Entity, 0)
	query.NewQuery(
		filter.And(
			filter.Contains(ProjectileTag),
			filter.Contains(Position),
			filter.Contains(Size),
		),
	).Each(cs.world, func(entry *donburi.Entry) {
		projectiles = append(projectiles, entry.Entity())
	})

	// Get all enemies
	enemies := make([]donburi.Entity, 0)
	query.NewQuery(
		filter.And(
			filter.Contains(EnemyTag),
			filter.Contains(Position),
			filter.Contains(Size),
			filter.Contains(Health),
		),
	).Each(cs.world, func(entry *donburi.Entry) {
		enemies = append(enemies, entry.Entity())
	})

	// Check each projectile against each enemy
	for _, projectileEntity := range projectiles {
		projectileEntry := cs.world.Entry(projectileEntity)
		if !projectileEntry.Valid() {
			continue
		}

		projectilePos := Position.Get(projectileEntry)
		projectileSize := Size.Get(projectileEntry)

		for _, enemyEntity := range enemies {
			enemyEntry := cs.world.Entry(enemyEntity)
			if !enemyEntry.Valid() {
				continue
			}

			enemyPos := Position.Get(enemyEntry)
			enemySize := Size.Get(enemyEntry)

			// Check collision
			if cs.checkCollision(*projectilePos, *projectileSize, *enemyPos, *enemySize) {
				// Handle collision
				cs.handleProjectileEnemyCollision(projectileEntity, enemyEntity, projectileEntry, enemyEntry)
				break // Projectile can only hit one enemy
			}
		}
	}
}

// checkPlayerEnemyCollisions checks for collisions between player and enemies
func (cs *CollisionSystem) checkPlayerEnemyCollisions() {
	// Get player
	players := make([]donburi.Entity, 0)
	query.NewQuery(
		filter.And(
			filter.Contains(PlayerTag),
			filter.Contains(Position),
			filter.Contains(Size),
		),
	).Each(cs.world, func(entry *donburi.Entry) {
		players = append(players, entry.Entity())
	})

	if len(players) == 0 {
		return
	}

	playerEntity := players[0]
	playerEntry := cs.world.Entry(playerEntity)
	if !playerEntry.Valid() {
		return
	}

	playerPos := Position.Get(playerEntry)
	playerSize := Size.Get(playerEntry)

	// Get all enemies
	enemies := make([]donburi.Entity, 0)
	query.NewQuery(
		filter.And(
			filter.Contains(EnemyTag),
			filter.Contains(Position),
			filter.Contains(Size),
		),
	).Each(cs.world, func(entry *donburi.Entry) {
		enemies = append(enemies, entry.Entity())
	})

	// Check player against each enemy
	for _, enemyEntity := range enemies {
		enemyEntry := cs.world.Entry(enemyEntity)
		if !enemyEntry.Valid() {
			continue
		}

		enemyPos := Position.Get(enemyEntry)
		enemySize := Size.Get(enemyEntry)

		// Check collision
		if cs.checkCollision(*playerPos, *playerSize, *enemyPos, *enemySize) {
			// Handle collision
			cs.handlePlayerEnemyCollision(playerEntity, enemyEntity, playerEntry, enemyEntry)
		}
	}
}

// checkCollision checks if two entities are colliding using AABB collision detection
func (cs *CollisionSystem) checkCollision(
	pos1 common.Point, size1 common.Size,
	pos2 common.Point, size2 common.Size,
) bool {
	// Calculate bounding boxes
	left1 := pos1.X
	right1 := pos1.X + float64(size1.Width)
	top1 := pos1.Y
	bottom1 := pos1.Y + float64(size1.Height)

	left2 := pos2.X
	right2 := pos2.X + float64(size2.Width)
	top2 := pos2.Y
	bottom2 := pos2.Y + float64(size2.Height)

	// Check for overlap
	return left1 < right2 && right1 > left2 && top1 < bottom2 && bottom1 > top2
}

// handleProjectileEnemyCollision handles collision between a projectile and an enemy
func (cs *CollisionSystem) handleProjectileEnemyCollision(
	projectileEntity, enemyEntity donburi.Entity,
	projectileEntry, enemyEntry *donburi.Entry,
) {
	// Get projectile and enemy data
	projectilePos := Position.Get(projectileEntry)
	projectileSize := Size.Get(projectileEntry)
	enemyPos := Position.Get(enemyEntry)
	enemySize := Size.Get(enemyEntry)
	enemyHealth := Health.Get(enemyEntry)

	// Check collision
	if cs.checkCollision(*projectilePos, *projectileSize, *enemyPos, *enemySize) {
		// Reduce enemy health
		newHealth := *enemyHealth - 1
		Health.SetValue(enemyEntry, newHealth)

		// Remove projectile
		cs.world.Remove(projectileEntity)

		// Remove enemy if health reaches 0
		if newHealth <= 0 {
			cs.world.Remove(enemyEntity)
		}
	}
}

// handlePlayerEnemyCollision handles collision between the player and an enemy
func (cs *CollisionSystem) handlePlayerEnemyCollision(
	playerEntity, enemyEntity donburi.Entity,
	playerEntry, enemyEntry *donburi.Entry,
) {
	// Get player and enemy data
	playerPos := Position.Get(playerEntry)
	playerSize := Size.Get(playerEntry)
	enemyPos := Position.Get(enemyEntry)
	enemySize := Size.Get(enemyEntry)

	// Check collision
	if cs.checkCollision(*playerPos, *playerSize, *enemyPos, *enemySize) {
		// Remove enemy
		cs.world.Remove(enemyEntity)

		// TODO: Handle player damage/lives
		// For now, just remove the enemy
	}
}

// GetCollisionDistance calculates the distance between two entities
func (cs *CollisionSystem) GetCollisionDistance(pos1, pos2 common.Point) float64 {
	dx := pos1.X - pos2.X
	dy := pos1.Y - pos2.Y
	return math.Sqrt(dx*dx + dy*dy)
}

// IsWithinRange checks if two entities are within a specified range
func (cs *CollisionSystem) IsWithinRange(pos1, pos2 common.Point, maxDistance float64) bool {
	distance := cs.GetCollisionDistance(pos1, pos2)
	return distance <= maxDistance
}
