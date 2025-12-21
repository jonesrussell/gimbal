package health

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/jonesrussell/gimbal/internal/ecs/core"
)

// checkGameOverCondition checks if the game should end
func (hs *HealthSystem) checkGameOverCondition() {
	players := hs.getPlayerEntities()
	if len(players) == 0 {
		hs.triggerGameOver("no player entity found")
		return
	}

	if hs.areAllPlayersDead(players) {
		hs.triggerGameOver("all players dead")
	}
}

// getPlayerEntities returns all player entities with health components
func (hs *HealthSystem) getPlayerEntities() []donburi.Entity {
	players := make([]donburi.Entity, 0)
	query.NewQuery(
		filter.And(
			filter.Contains(core.PlayerTag),
			filter.Contains(core.Health),
		),
	).Each(hs.world, func(entry *donburi.Entry) {
		players = append(players, entry.Entity())
	})
	return players
}

// areAllPlayersDead checks if all players have zero health
func (hs *HealthSystem) areAllPlayersDead(players []donburi.Entity) bool {
	for _, playerEntity := range players {
		playerEntry := hs.world.Entry(playerEntity)
		if playerEntry.Valid() {
			health := core.Health.Get(playerEntry)
			if health.Current > 0 {
				return false
			}
		}
	}
	return true
}

// triggerGameOver triggers game over state and events
func (hs *HealthSystem) triggerGameOver(reason string) {
	if hs.gameStateManager != nil {
		hs.gameStateManager.SetGameOver(true)
	}
	if hs.eventSystem != nil {
		hs.eventSystem.EmitGameOver()
	}
	hs.logger.Debug("Game over", "reason", reason)
}
