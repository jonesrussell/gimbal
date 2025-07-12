package collision

import (
	"github.com/yohamta/donburi"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/ecs/managers"
)

// CollisionSystemConfig groups all dependencies for CollisionSystem
// to avoid argument limit lint violations
type CollisionSystemConfig struct {
	World        donburi.World
	Config       *common.GameConfig
	HealthSystem interface{} // Using interface to avoid circular dependency
	EventSystem  interface{} // Using interface to avoid circular dependency
	ScoreManager *managers.ScoreManager
	EnemySystem  interface{} // Using interface to avoid circular dependency
	Logger       common.Logger
}

// CollisionSystem manages collision detection and response
type CollisionSystem struct {
	world        donburi.World
	config       *common.GameConfig
	healthSystem interface{} // Using interface to avoid circular dependency
	eventSystem  interface{} // Using interface to avoid circular dependency
	scoreManager *managers.ScoreManager
	enemySystem  interface{} // Using interface to avoid circular dependency
	logger       common.Logger
}

// NewCollisionSystem creates a new collision system
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
