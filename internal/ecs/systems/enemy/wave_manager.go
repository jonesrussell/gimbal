package enemy

import (
	"math/rand"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/ecs/core"
)

// WaveConfig defines the configuration for a wave
type WaveConfig struct {
	FormationType FormationType
	EnemyCount    int
	EnemyTypes    []EnemyType // Types of enemies in this wave
	SpawnDelay    float64     // Delay between individual enemy spawns
	Timeout       float64     // Time before wave auto-completes (0 = no timeout)
}

// WaveState tracks the current state of a wave
type WaveState struct {
	WaveIndex      int
	Config         WaveConfig
	EnemiesSpawned int
	EnemiesKilled  int
	WaveTimer      float64
	IsComplete     bool
	IsSpawning     bool
	LastSpawnTime  float64
}

// WaveManager manages wave spawning and completion
type WaveManager struct {
	world       donburi.World
	currentWave *WaveState
	waves       []WaveConfig
	waveIndex   int
	logger      common.Logger
}

// NewWaveManager creates a new wave manager
func NewWaveManager(world donburi.World, logger common.Logger) *WaveManager {
	wm := &WaveManager{
		world:     world,
		waveIndex: 0,
		logger:    logger,
	}

	// Initialize default waves
	wm.waves = wm.generateDefaultWaves()

	return wm
}

// generateDefaultWaves creates default wave configurations
func (wm *WaveManager) generateDefaultWaves() []WaveConfig {
	waves := []WaveConfig{
		// Wave 1: Line formation, basic enemies
		{
			FormationType: FormationLine,
			EnemyCount:    5,
			EnemyTypes:    []EnemyType{EnemyTypeBasic},
			SpawnDelay:    0.2,
			Timeout:       30.0,
		},
		// Wave 2: Circle formation, basic enemies
		{
			FormationType: FormationCircle,
			EnemyCount:    8,
			EnemyTypes:    []EnemyType{EnemyTypeBasic},
			SpawnDelay:    0.15,
			Timeout:       30.0,
		},
		// Wave 3: V-formation, mix of basic and heavy
		{
			FormationType: FormationV,
			EnemyCount:    8,
			EnemyTypes:    []EnemyType{EnemyTypeBasic, EnemyTypeHeavy},
			SpawnDelay:    0.2,
			Timeout:       30.0,
		},
		// Wave 4: Circle with more enemies
		{
			FormationType: FormationCircle,
			EnemyCount:    10,
			EnemyTypes:    []EnemyType{EnemyTypeBasic, EnemyTypeHeavy},
			SpawnDelay:    0.15,
			Timeout:       30.0,
		},
		// Wave 5: Line with heavy enemies
		{
			FormationType: FormationLine,
			EnemyCount:    6,
			EnemyTypes:    []EnemyType{EnemyTypeHeavy},
			SpawnDelay:    0.25,
			Timeout:       30.0,
		},
	}

	return waves
}

// StartNextWave starts the next wave
func (wm *WaveManager) StartNextWave() *WaveConfig {
	if wm.waveIndex >= len(wm.waves) {
		return nil // All waves complete
	}

	config := wm.waves[wm.waveIndex]
	wm.currentWave = &WaveState{
		WaveIndex:      wm.waveIndex,
		Config:         config,
		EnemiesSpawned: 0,
		EnemiesKilled:  0,
		WaveTimer:      0,
		IsComplete:     false,
		IsSpawning:     true,
		LastSpawnTime:  -1, // Start spawning immediately
	}

	wm.logger.Debug("Wave started",
		"wave", wm.waveIndex+1,
		"formation", config.FormationType,
		"count", config.EnemyCount)

	return &config
}

// Update updates the wave state
func (wm *WaveManager) Update(deltaTime float64) {
	if wm.currentWave == nil {
		return
	}

	if wm.currentWave.IsComplete {
		return
	}

	wm.currentWave.WaveTimer += deltaTime

	// Check timeout
	if wm.currentWave.Config.Timeout > 0 && wm.currentWave.WaveTimer >= wm.currentWave.Config.Timeout {
		wm.currentWave.IsComplete = true
		wm.logger.Debug("Wave completed by timeout", "wave", wm.currentWave.WaveIndex+1)
		return
	}

	// Check if all enemies are killed
	activeEnemies := wm.countActiveEnemies()
	if activeEnemies == 0 && wm.currentWave.EnemiesSpawned >= wm.currentWave.Config.EnemyCount {
		wm.currentWave.IsComplete = true
		wm.logger.Debug("Wave completed - all enemies killed", "wave", wm.currentWave.WaveIndex+1)
	}
}

// countActiveEnemies counts how many enemies are currently active
func (wm *WaveManager) countActiveEnemies() int {
	count := 0
	query.NewQuery(
		filter.And(
			filter.Contains(core.EnemyTag),
		),
	).Each(wm.world, func(entry *donburi.Entry) {
		count++
	})
	return count
}

// ShouldSpawnEnemy checks if it's time to spawn the next enemy
func (wm *WaveManager) ShouldSpawnEnemy(deltaTime float64) bool {
	if wm.currentWave == nil || wm.currentWave.IsComplete || !wm.currentWave.IsSpawning {
		return false
	}

	if wm.currentWave.EnemiesSpawned >= wm.currentWave.Config.EnemyCount {
		wm.currentWave.IsSpawning = false
		return false
	}

	// Check spawn delay
	if wm.currentWave.LastSpawnTime < 0 {
		return true // First enemy spawns immediately
	}

	timeSinceLastSpawn := wm.currentWave.WaveTimer - wm.currentWave.LastSpawnTime
	return timeSinceLastSpawn >= wm.currentWave.Config.SpawnDelay
}

// GetNextEnemyType returns the type of enemy to spawn next
func (wm *WaveManager) GetNextEnemyType() EnemyType {
	if wm.currentWave == nil {
		return EnemyTypeBasic
	}

	// If multiple types, randomly select one
	if len(wm.currentWave.Config.EnemyTypes) > 1 {
		//nolint:gosec // Game logic randomness is acceptable
		idx := rand.Intn(len(wm.currentWave.Config.EnemyTypes))
		return wm.currentWave.Config.EnemyTypes[idx]
	}

	if len(wm.currentWave.Config.EnemyTypes) > 0 {
		return wm.currentWave.Config.EnemyTypes[0]
	}

	return EnemyTypeBasic
}

// MarkEnemySpawned marks that an enemy has been spawned
func (wm *WaveManager) MarkEnemySpawned() {
	if wm.currentWave != nil {
		wm.currentWave.EnemiesSpawned++
		wm.currentWave.LastSpawnTime = wm.currentWave.WaveTimer
	}
}

// MarkEnemyKilled marks that an enemy has been killed
func (wm *WaveManager) MarkEnemyKilled() {
	if wm.currentWave != nil {
		wm.currentWave.EnemiesKilled++
	}
}

// IsWaveComplete returns true if the current wave is complete
func (wm *WaveManager) IsWaveComplete() bool {
	return wm.currentWave != nil && wm.currentWave.IsComplete
}

// CompleteWave marks the current wave as complete and advances to next
func (wm *WaveManager) CompleteWave() {
	if wm.currentWave != nil {
		wm.currentWave.IsComplete = true
	}
	wm.waveIndex++
}

// GetCurrentWave returns the current wave state
func (wm *WaveManager) GetCurrentWave() *WaveState {
	return wm.currentWave
}

// HasMoreWaves returns true if there are more waves
func (wm *WaveManager) HasMoreWaves() bool {
	return wm.waveIndex < len(wm.waves)
}

// Reset resets the wave manager to the first wave
func (wm *WaveManager) Reset() {
	wm.waveIndex = 0
	wm.currentWave = nil
}

// GetWaveCount returns the total number of waves
func (wm *WaveManager) GetWaveCount() int {
	return len(wm.waves)
}

// GetCurrentWaveIndex returns the current wave index (0-based)
func (wm *WaveManager) GetCurrentWaveIndex() int {
	return wm.waveIndex
}
