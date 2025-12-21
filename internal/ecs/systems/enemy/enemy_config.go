package enemy

import (
	"fmt"

	"github.com/jonesrussell/gimbal/internal/ecs/managers"
)

// Enemy system constants
const (
	// DefaultSpawnIntervalSeconds is the time between enemy spawns (legacy, now used for wave delays)
	DefaultSpawnIntervalSeconds = 1.0
	// DefaultEnemySpeed is the movement speed of enemies
	DefaultEnemySpeed = 2.0
	// DefaultEnemySize is the size of enemy sprites
	DefaultEnemySize = 32
	// BossSpawnDelay is the delay after last wave before boss spawns
	BossSpawnDelay = 2.0
)

// LoadEnemyConfigs loads enemy type configurations into the system
// This must be called after NewEnemySystem and before using the system
func (es *EnemySystem) LoadEnemyConfigs(configs map[EnemyType]EnemyTypeData) {
	es.enemyConfigs = configs
	es.logger.Debug("Enemy configs loaded into EnemySystem", "count", len(configs))
}

// GetEnemyTypeData returns the configuration for an enemy type from loaded configs
// Returns an error if the enemy type is not found in loaded configs
func (es *EnemySystem) GetEnemyTypeData(enemyType EnemyType) (EnemyTypeData, error) {
	enemyConfig, ok := es.enemyConfigs[enemyType]
	if !ok {
		return EnemyTypeData{}, fmt.Errorf("enemy type %d not found in loaded configs", enemyType)
	}
	return enemyConfig, nil
}

// LoadLevelConfig loads the waves and boss configuration for a level
func (es *EnemySystem) LoadLevelConfig(waves []WaveConfig, bossConfig *managers.BossConfig) {
	es.waveManager.LoadWaves(waves)
	es.bossConfig = bossConfig
	es.bossSpawned = false
	es.bossSpawnTimer = 0

	bossHealth := 0
	if bossConfig != nil {
		bossHealth = bossConfig.Health
	}

	es.logger.Debug("Level config loaded",
		"waves", len(waves),
		"boss_enabled", bossConfig != nil && bossConfig.Enabled,
		"boss_health", bossHealth)
}

// Reset resets the enemy system for a new level
func (es *EnemySystem) Reset() {
	es.waveManager.Reset()
	es.bossSpawned = false
	es.bossSpawnTimer = 0
	es.spawnTimer = 0
	// Note: bossConfig is not reset here as it should be set via LoadLevelConfig
}
