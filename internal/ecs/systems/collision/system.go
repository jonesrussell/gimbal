package collision

import (
	"context"
	"time"

	"github.com/yohamta/donburi"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/config"
	"github.com/jonesrussell/gimbal/internal/ecs/contracts"
)

// CollisionSystem manages collision detection and response with proper type safety
type CollisionSystem struct {
	world    donburi.World
	config   *config.GameConfig
	registry contracts.SystemRegistry
	logger   common.Logger
}

// NewCollisionSystem creates a new collision system with proper dependency injection
func NewCollisionSystem(
	world donburi.World,
	config *config.GameConfig,
	registry contracts.SystemRegistry,
	logger common.Logger,
) *CollisionSystem {
	return &CollisionSystem{
		world:    world,
		config:   config,
		registry: registry,
		logger:   logger,
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
