package collision

import (
	"context"

	"github.com/yohamta/donburi"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/ecs/core"
)

// checkProjectileEnemyCollisions checks for collisions between projectiles and enemies
func (cs *CollisionSystem) checkProjectileEnemyCollisions(ctx context.Context) error {
	projectiles, err := cs.getProjectileEntities(ctx)
	if err != nil {
		return err
	}
	enemies, err := cs.getEnemyEntities(ctx)
	if err != nil {
		return err
	}
	return cs.processProjectileEnemyCollisions(ctx, projectiles, enemies)
}

func (cs *CollisionSystem) getProjectileEntities(ctx context.Context) ([]donburi.Entity, error) {
	entries := core.GetProjectileEntries(cs.world)
	projectiles := make([]donburi.Entity, 0, len(entries))
	for _, entry := range entries {
		projectiles = append(projectiles, entry.Entity())
	}
	return projectiles, nil
}

func (cs *CollisionSystem) processProjectileEnemyCollisions(
	ctx context.Context, projectiles, enemies []donburi.Entity,
) error {
	for _, projectileEntity := range projectiles {
		if err := cs.checkSingleProjectileCollisions(ctx, projectileEntity, enemies); err != nil {
			return err
		}
	}
	return nil
}

func (cs *CollisionSystem) checkSingleProjectileCollisions(
	ctx context.Context, projectileEntity donburi.Entity, enemies []donburi.Entity,
) error {
	// Check for cancellation
	if err := common.CheckContextCancellation(ctx); err != nil {
		return err
	}

	projectileEntry := cs.world.Entry(projectileEntity)
	if !projectileEntry.Valid() {
		return nil
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
			if err := cs.handleProjectileEnemyCollision(
				ctx, projectileEntity, enemyEntity, projectileEntry, enemyEntry,
			); err != nil {
				return err
			}
			return nil // Projectile can only hit one enemy
		}
	}
	return nil
}

// handleProjectileEnemyCollision handles collision between a projectile and an enemy
func (cs *CollisionSystem) handleProjectileEnemyCollision(
	ctx context.Context,
	projectileEntity, enemyEntity donburi.Entity,
	projectileEntry, enemyEntry *donburi.Entry,
) error {
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
			// Remove enemy entity and get points
			points := cs.enemySystem.DestroyEnemy(enemyEntity)

			// Emit enemy destroyed event (includes points)
			if err := cs.eventSystem.EmitEnemyDestroyed(ctx, enemyEntity, points); err != nil {
				cs.logger.Error("Failed to emit enemy destroyed event", "error", err)
			}

			cs.logger.Debug("Enemy destroyed", "points", points)
		}
	}

	return nil
}

// checkEnemyProjectilePlayerCollisions checks for collisions between enemy projectiles and player
func (cs *CollisionSystem) checkEnemyProjectilePlayerCollisions(ctx context.Context) error {
	// Check for cancellation
	if err := common.CheckContextCancellation(ctx); err != nil {
		return err
	}

	// Get player entity
	playerEntity, playerEntry := cs.getPlayerEntity()
	if playerEntry == nil {
		return nil
	}

	playerPos := core.Position.Get(playerEntry)
	playerSize := core.Size.Get(playerEntry)

	// Get all enemy projectiles
	enemyProjectiles := cs.getEnemyProjectileEntities()

	for _, projectileEntity := range enemyProjectiles {
		projectileEntry := cs.world.Entry(projectileEntity)
		if !projectileEntry.Valid() {
			continue
		}

		projectilePos := core.Position.Get(projectileEntry)
		projectileSize := core.Size.Get(projectileEntry)

		// Check collision
		if cs.checkCollision(*projectilePos, *projectileSize, *playerPos, *playerSize) {
			// Remove the projectile
			cs.world.Remove(projectileEntity)

			// Damage the player (1 damage per projectile hit) with proper context propagation
			cs.healthSystem.DamagePlayer(ctx, playerEntity, 1)

			cs.logger.Debug("Player hit by enemy projectile")
		}
	}

	return nil
}

// getEnemyProjectileEntities returns all enemy projectile entities
func (cs *CollisionSystem) getEnemyProjectileEntities() []donburi.Entity {
	entries := core.GetEnemyProjectileEntries(cs.world)
	projectiles := make([]donburi.Entity, 0, len(entries))
	for _, entry := range entries {
		projectiles = append(projectiles, entry.Entity())
	}
	return projectiles
}
