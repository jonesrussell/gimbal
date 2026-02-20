package enemy

import (
	"time"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/jonesrussell/gimbal/internal/dbg"
	"github.com/jonesrussell/gimbal/internal/ecs/core"
	"github.com/jonesrussell/gimbal/internal/ecs/managers"
)

// GyrussWaveManager manages Gyruss-style wave spawning from stage configuration
type GyrussWaveManager struct {
	world       donburi.World
	stageConfig *managers.StageConfig

	// Wave state
	currentWaveIndex  int
	currentGroupIndex int
	spawnIndex        int

	// Timing
	waveTimer       time.Duration
	groupSpawnTimer time.Duration

	// Flags
	isSpawning bool
}

// NewGyrussWaveManager creates a new Gyruss-style wave manager
func NewGyrussWaveManager(world donburi.World) *GyrussWaveManager {
	return &GyrussWaveManager{world: world}
}

// LoadStage loads a stage configuration
func (gwm *GyrussWaveManager) LoadStage(config *managers.StageConfig) {
	gwm.stageConfig = config
	gwm.Reset()

	dbg.Log(dbg.System, "Gyruss stage loaded (stage=%d waves=%d)", config.StageNumber, len(config.Waves))
}

// Reset resets the wave manager state
func (gwm *GyrussWaveManager) Reset() {
	gwm.currentWaveIndex = 0
	gwm.currentGroupIndex = 0
	gwm.spawnIndex = 0
	gwm.waveTimer = 0
	gwm.groupSpawnTimer = 0
	gwm.isSpawning = false
}

// Update updates the wave manager state (spawn timers only; wave transitions are driven by StageStateMachine)
func (gwm *GyrussWaveManager) Update(deltaTime float64) {
	deltaDuration := time.Duration(deltaTime * float64(time.Second))
	if gwm.isSpawning {
		gwm.waveTimer += deltaDuration
		gwm.groupSpawnTimer += deltaDuration
	}
}

// getCurrentWave returns the current wave config
func (gwm *GyrussWaveManager) getCurrentWave() *managers.GyrussWave {
	if gwm.stageConfig == nil || gwm.currentWaveIndex >= len(gwm.stageConfig.Waves) {
		return nil
	}
	return &gwm.stageConfig.Waves[gwm.currentWaveIndex]
}

// getCurrentGroup returns the current enemy group config
func (gwm *GyrussWaveManager) getCurrentGroup() *managers.EnemyGroupConfig {
	wave := gwm.getCurrentWave()
	if wave == nil || gwm.currentGroupIndex >= len(wave.SpawnSequence) {
		return nil
	}
	return &wave.SpawnSequence[gwm.currentGroupIndex]
}

// StartWave starts spawning the given wave index (called by StageStateMachine)
func (gwm *GyrussWaveManager) StartWave(waveIndex int) {
	if gwm.stageConfig == nil || waveIndex < 0 || waveIndex >= len(gwm.stageConfig.Waves) {
		return
	}
	gwm.currentWaveIndex = waveIndex
	gwm.currentGroupIndex = 0
	gwm.spawnIndex = 0
	gwm.waveTimer = 0
	gwm.groupSpawnTimer = 0
	gwm.isSpawning = true

	wave := &gwm.stageConfig.Waves[waveIndex]
	dbg.Log(dbg.System, "Gyruss wave started (wave_index=%d groups=%d)", waveIndex, len(wave.SpawnSequence))
}

// ShouldSpawnEnemy returns true if an enemy should be spawned and the spawn config
func (gwm *GyrussWaveManager) ShouldSpawnEnemy() (*managers.EnemyGroupConfig, bool) {
	if !gwm.isSpawning {
		return nil, false
	}

	group := gwm.getCurrentGroup()
	if group == nil {
		return nil, false
	}

	// Check if we've spawned all enemies in this group
	if gwm.spawnIndex >= group.Count {
		// Move to next group
		gwm.currentGroupIndex++
		gwm.spawnIndex = 0
		gwm.groupSpawnTimer = 0
		return nil, false
	}

	// Check spawn delay for first enemy in group
	spawnDelay := time.Duration(group.SpawnDelay * float64(time.Second))
	if gwm.spawnIndex == 0 && gwm.waveTimer < spawnDelay {
		return nil, false
	}

	// Check spawn interval for subsequent enemies
	if gwm.spawnIndex > 0 {
		interval := time.Duration(group.SpawnInterval * float64(time.Second))
		if gwm.groupSpawnTimer < interval {
			return nil, false
		}
	}

	return group, true
}

// MarkEnemySpawned marks that an enemy was spawned from current group
func (gwm *GyrussWaveManager) MarkEnemySpawned() {
	gwm.spawnIndex++
	gwm.groupSpawnTimer = 0
}

// AllSpawnedForCurrentWave returns true when all groups in the current wave have finished spawning
func (gwm *GyrussWaveManager) AllSpawnedForCurrentWave() bool {
	wave := gwm.getCurrentWave()
	if wave == nil {
		return false
	}
	return gwm.currentGroupIndex >= len(wave.SpawnSequence)
}

// ActiveEnemyCount returns the number of active enemies in the world (excluding boss)
func (gwm *GyrussWaveManager) ActiveEnemyCount() int {
	return gwm.countActiveEnemies()
}

// countActiveEnemies counts active enemies (excluding boss)
func (gwm *GyrussWaveManager) countActiveEnemies() int {
	count := 0
	query.NewQuery(
		filter.And(
			filter.Contains(core.EnemyTag),
		),
	).Each(gwm.world, func(entry *donburi.Entry) {
		// Skip boss
		if entry.HasComponent(core.EnemyTypeID) {
			typeID := core.EnemyTypeID.Get(entry)
			if EnemyType(*typeID) == EnemyTypeBoss {
				return
			}
		}
		count++
	})
	return count
}

// HasMoreWaves returns true if there are more waves
func (gwm *GyrussWaveManager) HasMoreWaves() bool {
	if gwm.stageConfig == nil {
		return false
	}
	return gwm.currentWaveIndex < len(gwm.stageConfig.Waves)
}

// GetBossConfig returns the boss configuration
func (gwm *GyrussWaveManager) GetBossConfig() *managers.StageBossConfig {
	if gwm.stageConfig == nil {
		return nil
	}
	return &gwm.stageConfig.Boss
}

// GetPowerUpConfig returns the power-up configuration
func (gwm *GyrussWaveManager) GetPowerUpConfig() *managers.PowerUpConfig {
	if gwm.stageConfig == nil {
		return nil
	}
	return &gwm.stageConfig.PowerUps
}

// GetDifficulty returns the difficulty settings
func (gwm *GyrussWaveManager) GetDifficulty() *managers.DifficultySettings {
	if gwm.stageConfig == nil {
		return nil
	}
	return &gwm.stageConfig.Difficulty
}

// GetCurrentWaveIndex returns the current wave index
func (gwm *GyrussWaveManager) GetCurrentWaveIndex() int {
	return gwm.currentWaveIndex
}

// GetCurrentGroupSpawnIndex returns the spawn index for the current group (next enemy to spawn).
// Use before MarkEnemySpawned so the spawner gets the correct index for orbit angle.
func (gwm *GyrussWaveManager) GetCurrentGroupSpawnIndex() int {
	return gwm.spawnIndex
}

// GetWaveCount returns total wave count
func (gwm *GyrussWaveManager) GetWaveCount() int {
	if gwm.stageConfig == nil {
		return 0
	}
	return len(gwm.stageConfig.Waves)
}

// GetStageConfig returns the current stage configuration
func (gwm *GyrussWaveManager) GetStageConfig() *managers.StageConfig {
	return gwm.stageConfig
}
