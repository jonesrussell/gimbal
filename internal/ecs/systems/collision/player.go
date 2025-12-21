package collision

import (
	"context"

	"github.com/yohamta/donburi"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/config"
	"github.com/jonesrussell/gimbal/internal/ecs/core"
)

// getPlayerEntity returns the first valid player entity
func (cs *CollisionSystem) getPlayerEntity() (donburi.Entity, *donburi.Entry) {
	entries := core.GetPlayerEntries(cs.world)
	if len(entries) == 0 {
		return 0, nil
	}
	playerEntry := entries[0]
	if !playerEntry.Valid() {
		return 0, nil
	}
	playerEntity := playerEntry.Entity()
	return playerEntity, playerEntry
}

// checkPlayerEnemyCollisions checks for collisions between player and enemies
func (cs *CollisionSystem) checkPlayerEnemyCollisions(ctx context.Context) error {
	// Check for cancellation
	if err := common.CheckContextCancellation(ctx); err != nil {
		return err
	}

	playerEntity, playerEntry := cs.getPlayerEntity()
	if playerEntry == nil {
		return nil
	}
	playerPos := core.Position.Get(playerEntry)
	playerSize := core.Size.Get(playerEntry)

	enemies, err := cs.getEnemyEntities(ctx)
	if err != nil {
		return err
	}
	return cs.checkCollisionsWithEnemies(ctx, PlayerCollisionData{
		Entity: playerEntity,
		Entry:  playerEntry,
		Pos:    playerPos,
		Size:   playerSize,
	}, enemies)
}

type PlayerCollisionData struct {
	Entity donburi.Entity
	Entry  *donburi.Entry
	Pos    *common.Point
	Size   *config.Size
}

func (cs *CollisionSystem) checkCollisionsWithEnemies(
	ctx context.Context,
	player PlayerCollisionData,
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
		if cs.checkCollision(*player.Pos, *player.Size, *enemyPos, *enemySize) {
			if err := cs.handlePlayerEnemyCollision(
				ctx, player.Entity, enemyEntity, player.Entry, enemyEntry,
			); err != nil {
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
