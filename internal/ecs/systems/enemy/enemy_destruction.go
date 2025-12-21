package enemy

import (
	"github.com/yohamta/donburi"

	"github.com/jonesrussell/gimbal/internal/ecs/core"
)

// DestroyEnemy destroys an enemy entity and returns points based on type
func (es *EnemySystem) DestroyEnemy(entity donburi.Entity) int {
	entry := es.world.Entry(entity)
	if !entry.Valid() {
		return 0
	}

	// Get enemy type from component (preferred) or fall back to health heuristic
	var enemyType EnemyType
	if entry.HasComponent(core.EnemyTypeID) {
		typeID := core.EnemyTypeID.Get(entry)
		enemyType = EnemyType(*typeID)
	} else if entry.HasComponent(core.Health) {
		// Fallback for legacy entities without EnemyTypeID
		health := core.Health.Get(entry)
		if health.Maximum >= 10 {
			enemyType = EnemyTypeBoss
		} else if health.Maximum >= 2 {
			enemyType = EnemyTypeHeavy
		} else {
			enemyType = EnemyTypeBasic
		}
	} else {
		enemyType = EnemyTypeBasic
	}

	enemyData, err := es.GetEnemyTypeData(enemyType)
	if err != nil {
		es.logger.Error("Failed to get enemy type data for points", "type", enemyType, "error", err)
		return 0 // Skip scoring
	}
	points := enemyData.Points

	// Mark enemy killed in wave manager
	es.waveManager.MarkEnemyKilled()

	// Remove the entity from the world
	es.world.Remove(entity)

	return points
}
