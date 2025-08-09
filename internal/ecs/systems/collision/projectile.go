package collision

import (
	"context"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

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
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
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
