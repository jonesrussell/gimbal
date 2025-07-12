package collision

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/jonesrussell/gimbal/internal/ecs/core"
)

// checkPlayerEnemyCollisions checks for collisions between player and enemies
func (cs *CollisionSystem) checkPlayerEnemyCollisions() {
	// Get player
	players := make([]donburi.Entity, 0)
	query.NewQuery(
		filter.And(
			filter.Contains(core.PlayerTag),
			filter.Contains(core.Position),
			filter.Contains(core.Size),
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

	playerPos := core.Position.Get(playerEntry)
	playerSize := core.Size.Get(playerEntry)

	// Get all enemies
	enemies := make([]donburi.Entity, 0)
	query.NewQuery(
		filter.And(
			filter.Contains(core.EnemyTag),
			filter.Contains(core.Position),
			filter.Contains(core.Size),
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

		enemyPos := core.Position.Get(enemyEntry)
		enemySize := core.Size.Get(enemyEntry)

		// Check collision
		if cs.checkCollision(*playerPos, *playerSize, *enemyPos, *enemySize) {
			// Handle collision
			cs.handlePlayerEnemyCollision(playerEntity, enemyEntity, playerEntry, enemyEntry)
		}
	}
}

// handlePlayerEnemyCollision handles collision between the player and an enemy
func (cs *CollisionSystem) handlePlayerEnemyCollision(
	playerEntity, enemyEntity donburi.Entity,
	playerEntry, enemyEntry *donburi.Entry,
) {
	// Get player and enemy data
	playerPos := core.Position.Get(playerEntry)
	playerSize := core.Size.Get(playerEntry)
	enemyPos := core.Position.Get(enemyEntry)
	enemySize := core.Size.Get(enemyEntry)

	// Check collision
	if cs.checkCollision(*playerPos, *playerSize, *enemyPos, *enemySize) {
		// Remove enemy immediately
		cs.world.Remove(enemyEntity)

		// Damage player (1 damage per enemy collision)
		// Note: Using interface to avoid circular dependency
		if healthSystem, ok := cs.healthSystem.(interface {
			DamagePlayer(donburi.Entity, int)
		}); ok {
			healthSystem.DamagePlayer(playerEntity, 1)
		}
	}
}
