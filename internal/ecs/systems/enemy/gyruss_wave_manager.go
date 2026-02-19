package enemy

import (
	"time"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/ecs/core"
	"github.com/jonesrussell/gimbal/internal/ecs/managers"
)

// GyrussWaveManager manages Gyruss-style wave spawning from stage configuration
type GyrussWaveManager struct {
	world       donburi.World
	logger      common.Logger
	stageConfig *managers.StageConfig

	// Wave state
	currentWaveIndex  int
	currentGroupIndex int
	spawnIndex        int

	// Timing
	waveTimer       time.Duration
	groupSpawnTimer time.Duration
	interWaveTimer  time.Duration

	// Flags
	isSpawning             bool
	isWaitingInterWave     bool
	isWaitingForLevelStart bool
	levelStartTimer        time.Duration
	levelStartDelay        time.Duration

	// Boss state
	bossTriggered bool
}

// NewGyrussWaveManager creates a new Gyruss-style wave manager
func NewGyrussWaveManager(world donburi.World, logger common.Logger) *GyrussWaveManager {
	return &GyrussWaveManager{
		world:           world,
		logger:          logger,
		levelStartDelay: 3500 * time.Millisecond,
	}
}

// LoadStage loads a stage configuration
func (gwm *GyrussWaveManager) LoadStage(config *managers.StageConfig) {
	gwm.stageConfig = config
	gwm.Reset()
	gwm.isWaitingForLevelStart = true
	gwm.levelStartTimer = 0

	gwm.logger.Debug("Gyruss stage loaded",
		"stage", config.StageNumber,
		"name", config.Metadata.Name,
		"waves", len(config.Waves))
}

// Reset resets the wave manager state
func (gwm *GyrussWaveManager) Reset() {
	gwm.currentWaveIndex = 0
	gwm.currentGroupIndex = 0
	gwm.spawnIndex = 0
	gwm.waveTimer = 0
	gwm.groupSpawnTimer = 0
	gwm.interWaveTimer = 0
	gwm.isSpawning = false
	gwm.isWaitingInterWave = false
	gwm.isWaitingForLevelStart = false
	gwm.bossTriggered = false
}

// Update updates the wave manager state
func (gwm *GyrussWaveManager) Update(deltaTime float64) {
	deltaDuration := time.Duration(deltaTime * float64(time.Second))

	// Handle level start delay
	if gwm.isWaitingForLevelStart {
		gwm.levelStartTimer += deltaDuration
		if gwm.levelStartTimer >= gwm.levelStartDelay {
			gwm.isWaitingForLevelStart = false
			gwm.startNextWave()
		}
		return
	}

	// Handle inter-wave delay
	if gwm.isWaitingInterWave {
		gwm.interWaveTimer += deltaDuration
		wave := gwm.getCurrentWave()
		if wave != nil {
			delay := time.Duration(wave.Timing.InterWaveDelay * float64(time.Second))
			if gwm.interWaveTimer >= delay {
				gwm.isWaitingInterWave = false
				gwm.startNextWave()
			}
		}
		return
	}

	// Update wave timer
	if gwm.isSpawning {
		gwm.waveTimer += deltaDuration
		gwm.groupSpawnTimer += deltaDuration
	}

	// Check for wave completion
	gwm.checkWaveCompletion()
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

// startNextWave starts the next wave
func (gwm *GyrussWaveManager) startNextWave() {
	if gwm.stageConfig == nil || gwm.currentWaveIndex >= len(gwm.stageConfig.Waves) {
		return
	}

	wave := &gwm.stageConfig.Waves[gwm.currentWaveIndex]
	gwm.currentGroupIndex = 0
	gwm.spawnIndex = 0
	gwm.waveTimer = 0
	gwm.groupSpawnTimer = 0
	gwm.isSpawning = true

	gwm.logger.Debug("Gyruss wave started",
		"wave_id", wave.WaveID,
		"description", wave.Description,
		"groups", len(wave.SpawnSequence))
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

// checkWaveCompletion checks if the current wave is complete
func (gwm *GyrussWaveManager) checkWaveCompletion() {
	if !gwm.isSpawning {
		return
	}

	wave := gwm.getCurrentWave()
	if wave == nil {
		return
	}

	// Check timeout
	timeout := time.Duration(wave.Timing.Timeout * float64(time.Second))
	if timeout > 0 && gwm.waveTimer >= timeout {
		gwm.completeWave()
		return
	}

	// Check if all groups are done spawning
	allSpawned := gwm.currentGroupIndex >= len(wave.SpawnSequence)

	// Check if all enemies are killed
	activeEnemies := gwm.countActiveEnemies()

	if allSpawned && activeEnemies == 0 {
		gwm.completeWave()
	}
}

// completeWave completes the current wave and advances
func (gwm *GyrussWaveManager) completeWave() {
	wave := gwm.getCurrentWave()
	if wave == nil {
		return
	}

	gwm.logger.Debug("Gyruss wave completed",
		"wave_id", wave.WaveID,
		"on_clear", wave.OnClear)

	gwm.isSpawning = false

	// Check what happens on clear
	if wave.OnClear == "boss" {
		gwm.bossTriggered = true
		gwm.logger.Debug("Boss wave triggered")
		return
	}

	// Move to next wave
	gwm.currentWaveIndex++
	if gwm.currentWaveIndex < len(gwm.stageConfig.Waves) {
		nextWave := &gwm.stageConfig.Waves[gwm.currentWaveIndex]
		if nextWave.Timing.InterWaveDelay > 0 {
			gwm.isWaitingInterWave = true
			gwm.interWaveTimer = 0
		} else {
			gwm.startNextWave()
		}
	}
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

// IsBossTriggered returns true if boss should spawn
func (gwm *GyrussWaveManager) IsBossTriggered() bool {
	return gwm.bossTriggered
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

// IsWaitingForLevelStart returns true if waiting for level start
func (gwm *GyrussWaveManager) IsWaitingForLevelStart() bool {
	return gwm.isWaitingForLevelStart
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
