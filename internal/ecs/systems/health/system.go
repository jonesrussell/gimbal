package health

import (
	"time"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/ecs/core"
)

// HealthSystem manages player health, invincibility, and respawning
type HealthSystem struct {
	world            donburi.World
	config           *common.GameConfig
	eventSystem      interface{} // Using interface to avoid circular dependency
	gameStateManager interface{} // Using interface to avoid circular dependency
	logger           common.Logger
	lastUpdate       time.Time
}

// NewHealthSystem creates a new health system
func NewHealthSystem(
	world donburi.World,
	config *common.GameConfig,
	eventSystem interface{},
	gameStateManager interface{},
	logger common.Logger,
) *HealthSystem {
	return &HealthSystem{
		world:            world,
		config:           config,
		eventSystem:      eventSystem,
		gameStateManager: gameStateManager,
		logger:           logger,
		lastUpdate:       time.Now(),
	}
}

// Update updates the health system
func (hs *HealthSystem) Update() {
	now := time.Now()
	deltaTime := now.Sub(hs.lastUpdate).Seconds()
	hs.lastUpdate = now

	// Update invincibility timers for all entities with health
	query.NewQuery(
		filter.And(
			filter.Contains(core.Health),
		),
	).Each(hs.world, func(entry *donburi.Entry) {
		health := core.Health.Get(entry)
		if health.IsInvincible {
			health.InvincibilityTime -= deltaTime
			if health.InvincibilityTime <= 0 {
				health.IsInvincible = false
				health.InvincibilityTime = 0
			}
			core.Health.SetValue(entry, *health)
		}
	})

	// Check for game over condition
	hs.checkGameOverCondition()
}

// GetPlayerHealth returns the current health of the first player found
func (hs *HealthSystem) GetPlayerHealth() (current, maximum int) {
	players := make([]donburi.Entity, 0)
	query.NewQuery(
		filter.And(
			filter.Contains(core.PlayerTag),
			filter.Contains(core.Health),
		),
	).Each(hs.world, func(entry *donburi.Entry) {
		players = append(players, entry.Entity())
	})

	if len(players) > 0 {
		playerEntry := hs.world.Entry(players[0])
		if playerEntry.Valid() {
			health := core.Health.Get(playerEntry)
			current = health.Current
			maximum = health.Maximum
			return
		}
	}

	return
}

// IsPlayerInvincible returns whether the first player found is currently invincible
func (hs *HealthSystem) IsPlayerInvincible() bool {
	players := make([]donburi.Entity, 0)
	query.NewQuery(
		filter.And(
			filter.Contains(core.PlayerTag),
			filter.Contains(core.Health),
		),
	).Each(hs.world, func(entry *donburi.Entry) {
		players = append(players, entry.Entity())
	})

	if len(players) > 0 {
		playerEntry := hs.world.Entry(players[0])
		if playerEntry.Valid() {
			health := core.Health.Get(playerEntry)
			return health.IsInvincible
		}
	}

	return false
}
