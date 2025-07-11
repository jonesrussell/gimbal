package collision

import (
	"github.com/yohamta/donburi"

	"github.com/jonesrussell/gimbal/internal/common"
)

// CollisionSystem manages collision detection and response
type CollisionSystem struct {
	world        donburi.World
	config       *common.GameConfig
	healthSystem interface{} // Using interface to avoid circular dependency
	eventSystem  interface{} // Using interface to avoid circular dependency
	logger       common.Logger
}

// NewCollisionSystem creates a new collision system
func NewCollisionSystem(
	world donburi.World,
	config *common.GameConfig,
	healthSystem interface{},
	eventSystem interface{},
	logger common.Logger,
) *CollisionSystem {
	return &CollisionSystem{
		world:        world,
		config:       config,
		healthSystem: healthSystem,
		eventSystem:  eventSystem,
		logger:       logger,
	}
}

// Update updates the collision system
func (cs *CollisionSystem) Update() {
	// Check projectile-enemy collisions
	cs.checkProjectileEnemyCollisions()

	// Check player-enemy collisions
	cs.checkPlayerEnemyCollisions()
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
