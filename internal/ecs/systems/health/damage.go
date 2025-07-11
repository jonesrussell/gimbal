package health

import (
	"github.com/yohamta/donburi"

	"github.com/jonesrussell/gimbal/internal/ecs/core"
)

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
	if eventSystem, ok := hs.eventSystem.(interface {
		EmitPlayerDamaged(donburi.Entity, int, int)
	}); ok {
		eventSystem.EmitPlayerDamaged(playerEntity, damage, health.Current)
	}

	hs.logger.Debug("Player damaged", "damage", damage, "remaining_lives", health.Current)

	// Check if player should respawn or game over
	if health.Current > 0 {
		hs.respawnPlayer(playerEntity)
	} else {
		if gameStateManager, ok := hs.gameStateManager.(interface {
			SetGameOver(bool)
		}); ok {
			gameStateManager.SetGameOver(true)
		}
		if eventSystem, ok := hs.eventSystem.(interface {
			EmitGameOver()
		}); ok {
			eventSystem.EmitGameOver()
		}
		hs.logger.Debug("Game over - no lives remaining")
	}
}

// AddLife adds a life to the player
func (hs *HealthSystem) AddLife(playerEntity donburi.Entity) {
	playerEntry := hs.world.Entry(playerEntity)
	if !playerEntry.Valid() {
		return
	}

	health := core.Health.Get(playerEntry)
	health.Current++
	if health.Current > health.Maximum {
		health.Current = health.Maximum
	}

	core.Health.SetValue(playerEntry, *health)

	hs.logger.Debug("Life added to player", "new_lives", health.Current)
}
