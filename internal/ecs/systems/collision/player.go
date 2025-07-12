package collision

import (
	"context"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/config"
	"github.com/jonesrussell/gimbal/internal/ecs/core"
)

// getPlayerEntity returns the first valid player entity
func (cs *CollisionSystem) getPlayerEntity() (donburi.Entity, *donburi.Entry) {
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
		return 0, nil
	}
	playerEntity := players[0]
	playerEntry := cs.world.Entry(playerEntity)
	if !playerEntry.Valid() {
		return 0, nil
	}
	return playerEntity, playerEntry
}

// getEnemyEntities returns all valid enemy entities
func (cs *CollisionSystem) getEnemyEntities() []donburi.Entity {
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
	return enemies
}

// checkPlayerEnemyCollisions checks for collisions between player and enemies
func (cs *CollisionSystem) checkPlayerEnemyCollisions(ctx context.Context) error {
	// Check for cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	playerEntity, playerEntry := cs.getPlayerEntity()
	if playerEntry == nil {
		return nil
	}
	playerPos := core.Position.Get(playerEntry)
	playerSize := core.Size.Get(playerEntry)

	enemies := cs.getEnemyEntities()
	return cs.checkCollisionsWithEnemies(ctx, playerEntity, playerEntry, playerPos, playerSize, enemies)
}

// checkCollisionsWithEnemies checks collisions between the player and a list of enemies
func (cs *CollisionSystem) checkCollisionsWithEnemies(
	ctx context.Context,
	playerEntity donburi.Entity,
	playerEntry *donburi.Entry,
	playerPos *common.Point,
	playerSize *config.Size,
	enemies []donburi.Entity,
) error {
	for _, enemyEntity := range enemies {
		// Check for cancellation
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		enemyEntry := cs.world.Entry(enemyEntity)
		if !enemyEntry.Valid() {
			continue
		}

		enemyPos := core.Position.Get(enemyEntry)
		enemySize := core.Size.Get(enemyEntry)

		// Check collision
		if cs.checkCollision(*playerPos, *playerSize, *enemyPos, *enemySize) {
			// Handle collision
			if err := cs.handlePlayerEnemyCollision(ctx, playerEntity, enemyEntity, playerEntry, enemyEntry); err != nil {
				return err
			}
		}
	}
	return nil
}

// handlePlayerEnemyCollision handles collision between the player and an enemy
func (cs *CollisionSystem) handlePlayerEnemyCollision(
	ctx context.Context,
	playerEntity, enemyEntity donburi.Entity,
	playerEntry, enemyEntry *donburi.Entry,
) error {
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
		cs.healthSystem.DamagePlayer(playerEntity, 1)

		cs.logger.Debug("Player damaged by enemy collision")
	}

	return nil
}
