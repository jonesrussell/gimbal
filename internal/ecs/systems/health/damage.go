package health

import (
	"github.com/yohamta/donburi"

	"github.com/jonesrussell/gimbal/internal/ecs/core"
)

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
