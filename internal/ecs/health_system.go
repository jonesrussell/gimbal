package ecs

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
	eventSystem      *EventSystem
	gameStateManager *GameStateManager
	logger           common.Logger
	lastUpdate       time.Time
}

// NewHealthSystem creates a new health system
func NewHealthSystem(
	world donburi.World,
	config *common.GameConfig,
	eventSystem *EventSystem,
	gameStateManager *GameStateManager,
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

// DamagePlayer damages the player and handles invincibility
func (hs *HealthSystem) DamagePlayer(playerEntity donburi.Entity, damage int) {
	playerEntry := hs.world.Entry(playerEntity)
	if !playerEntry.Valid() {
		return
	}

	health := core.Health.Get(playerEntry)
	if health.IsInvincible {
		return // Player is invincible, no damage taken
	}

	// Apply damage
	health.Current -= damage
	if health.Current < 0 {
		health.Current = 0
	}

	// Set invincibility
	health.IsInvincible = true
	health.InvincibilityTime = health.InvincibilityDuration

	// Update health component
	core.Health.SetValue(playerEntry, *health)

	// Emit player damaged event
	hs.eventSystem.EmitPlayerDamaged(playerEntity, damage, health.Current)

	hs.logger.Debug("Player damaged", "damage", damage, "remaining_lives", health.Current)

	// Check if player should respawn or game over
	if health.Current > 0 {
		hs.respawnPlayer(playerEntity)
	} else {
		hs.gameStateManager.SetGameOver(true)
		hs.eventSystem.EmitGameOver()
		hs.logger.Debug("Game over - no lives remaining")
	}
}

// respawnPlayer respawns the player at the center bottom of the screen
func (hs *HealthSystem) respawnPlayer(playerEntity donburi.Entity) {
	playerEntry := hs.world.Entry(playerEntity)
	if !playerEntry.Valid() {
		return
	}

	// Reset position to center bottom
	center := common.Point{
		X: float64(hs.config.ScreenSize.Width) / 2,
		Y: float64(hs.config.ScreenSize.Height) / 2,
	}

	// Update position
	core.Position.SetValue(playerEntry, center)

	// Reset orbital data to bottom position (180 degrees)
	orbitalData := core.Orbital.Get(playerEntry)
	orbitalData.Center = center
	orbitalData.OrbitalAngle = common.HalfCircleDegrees // 180 degrees
	core.Orbital.SetValue(playerEntry, *orbitalData)

	// Reset angle
	core.Angle.SetValue(playerEntry, common.Angle(0))

	hs.logger.Debug("Player respawned at center bottom")
}

// checkGameOverCondition checks if the game should end
func (hs *HealthSystem) checkGameOverCondition() {
	// Check if any player exists and has health
	players := make([]donburi.Entity, 0)
	query.NewQuery(
		filter.And(
			filter.Contains(core.PlayerTag),
			filter.Contains(core.Health),
		),
	).Each(hs.world, func(entry *donburi.Entry) {
		players = append(players, entry.Entity())
	})

	if len(players) == 0 {
		// No player exists, trigger game over
		hs.gameStateManager.SetGameOver(true)
		hs.eventSystem.EmitGameOver()
		hs.logger.Debug("Game over - no player entity found")
		return
	}

	// Check if all players are dead
	allDead := true
	for _, playerEntity := range players {
		playerEntry := hs.world.Entry(playerEntity)
		if playerEntry.Valid() {
			health := core.Health.Get(playerEntry)
			if health.Current > 0 {
				allDead = false
				break
			}
		}
	}

	if allDead {
		hs.gameStateManager.SetGameOver(true)
		hs.eventSystem.EmitGameOver()
		hs.logger.Debug("Game over - all players dead")
	}
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

// IsPlayerInvincible returns whether the player is currently invincible
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

// AddLife adds a life to the player (for bonus lives)
func (hs *HealthSystem) AddLife(playerEntity donburi.Entity) {
	playerEntry := hs.world.Entry(playerEntity)
	if !playerEntry.Valid() {
		return
	}

	health := core.Health.Get(playerEntry)
	if health.Current < health.Maximum {
		health.Current++
		core.Health.SetValue(playerEntry, *health)
		hs.eventSystem.EmitLifeAdded(playerEntity, health.Current)
		hs.logger.Debug("Life added", "new_lives", health.Current)
	}
}
