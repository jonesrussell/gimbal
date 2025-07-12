package collision

import (
	"context"
	"time"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/config"
	"github.com/jonesrussell/gimbal/internal/ecs/core"
	"github.com/jonesrussell/gimbal/internal/ecs/managers"
)

// CollisionSystemConfig groups all dependencies for CollisionSystem
// to avoid argument limit lint violations
type CollisionSystemConfig struct {
	World        donburi.World
	Config       *config.GameConfig
	HealthSystem HealthSystemInterface
	EventSystem  EventSystemInterface
	ScoreManager *managers.ScoreManager
	EnemySystem  EnemySystemInterface
	Logger       common.Logger
}

// CollisionSystem manages collision detection and response with proper type safety
type CollisionSystem struct {
	world        donburi.World
	config       *config.GameConfig
	healthSystem HealthSystemInterface
	eventSystem  EventSystemInterface
	scoreManager *managers.ScoreManager
	enemySystem  EnemySystemInterface
	logger       common.Logger
}

// NewCollisionSystem creates a new collision system with proper dependency injection
func NewCollisionSystem(cfg *CollisionSystemConfig) *CollisionSystem {
	return &CollisionSystem{
		world:        cfg.World,
		config:       cfg.Config,
		healthSystem: cfg.HealthSystem,
		eventSystem:  cfg.EventSystem,
		scoreManager: cfg.ScoreManager,
		enemySystem:  cfg.EnemySystem,
		logger:       cfg.Logger,
	}
}

// Update updates the collision system with context support
func (cs *CollisionSystem) Update(ctx context.Context) error {
	// Add timeout for collision detection to prevent hanging
	ctx, cancel := context.WithTimeout(ctx, 16*time.Millisecond) // 60 FPS budget
	defer cancel()

	// Check projectile-enemy collisions
	if err := cs.checkProjectileEnemyCollisions(ctx); err != nil {
		return err
	}

	// Check player-enemy collisions
	if err := cs.checkPlayerEnemyCollisions(ctx); err != nil {
		return err
	}

	return nil
}

// GetCollisionDistance calculates the distance between two points
func (cs *CollisionSystem) GetCollisionDistance(pos1, pos2 common.Point) float64 {
	dx := pos1.X - pos2.X
	dy := pos1.Y - pos2.Y
	return dx*dx + dy*dy // Return squared distance for efficiency
}

// IsWithinRange checks if two points are within a specified distance
func (cs *CollisionSystem) IsWithinRange(pos1, pos2 common.Point, maxDistance float64) bool {
	distance := cs.GetCollisionDistance(pos1, pos2)
	return distance <= maxDistance*maxDistance // Compare squared distances
}

// getEnemyEntities returns all valid enemy entities with health
func (cs *CollisionSystem) getEnemyEntities(ctx context.Context) ([]donburi.Entity, error) {
	// Check for cancellation
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

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
	return enemies, nil
}
