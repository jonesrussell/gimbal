package collision

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/jonesrussell/gimbal/internal/ecs/core"
)

// checkProjectileEnemyCollisions checks for collisions between projectiles and enemies
func (cs *CollisionSystem) checkProjectileEnemyCollisions() {
	// Get all projectiles
	projectiles := make([]donburi.Entity, 0)
	query.NewQuery(
		filter.And(
			filter.Contains(core.ProjectileTag),
			filter.Contains(core.Position),
			filter.Contains(core.Size),
		),
	).Each(cs.world, func(entry *donburi.Entry) {
		projectiles = append(projectiles, entry.Entity())
	})

	// Get all enemies
	enemies := make([]donburi.Entity, 0)
	query.NewQuery(
		filter.And(
			filter.Contains(core.EnemyTag),
			filter.Contains(core.Position),
			filter.Contains(core.Size),
			filter.Contains(core.Health),
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

		projectilePos := core.Position.Get(projectileEntry)
		projectileSize := core.Size.Get(projectileEntry)

		for _, enemyEntity := range enemies {
			enemyEntry := cs.world.Entry(enemyEntity)
			if !enemyEntry.Valid() {
				continue
			}

			enemyPos := core.Position.Get(enemyEntry)
			enemySize := core.Size.Get(enemyEntry)

			// Check collision
			if cs.checkCollision(*projectilePos, *projectileSize, *enemyPos, *enemySize) {
				// Handle collision
				cs.handleProjectileEnemyCollision(projectileEntity, enemyEntity, projectileEntry, enemyEntry)
				break // Projectile can only hit one enemy
			}
		}
	}
}

// handleProjectileEnemyCollision handles collision between a projectile and an enemy
func (cs *CollisionSystem) handleProjectileEnemyCollision(
	projectileEntity, enemyEntity donburi.Entity,
	projectileEntry, enemyEntry *donburi.Entry,
) {
	// Get projectile and enemy data
	projectilePos := core.Position.Get(projectileEntry)
	projectileSize := core.Size.Get(projectileEntry)
	enemyPos := core.Position.Get(enemyEntry)
	enemySize := core.Size.Get(enemyEntry)
	enemyHealth := core.Health.Get(enemyEntry)

	// Check collision
	if cs.checkCollision(*projectilePos, *projectileSize, *enemyPos, *enemySize) {
		// Reduce enemy health
		enemyHealthData := *enemyHealth
		enemyHealthData.Current -= 1
		if enemyHealthData.Current < 0 {
			enemyHealthData.Current = 0
		}
		core.Health.SetValue(enemyEntry, enemyHealthData)

		// Remove projectile
		cs.world.Remove(projectileEntity)

		// Remove enemy if health reaches 0
		if enemyHealthData.Current <= 0 {
			// Award points for destroying enemy
			cs.scoreManager.AddScore(10) // 10 points per enemy
			cs.logger.Debug("Enemy destroyed, points awarded", "points", 10, "total_score", cs.scoreManager.GetScore())

			// Remove enemy entity
			cs.world.Remove(enemyEntity)
		}
	}
}
