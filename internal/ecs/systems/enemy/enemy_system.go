package enemy

import (
	"context"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/config"
	"github.com/jonesrussell/gimbal/internal/ecs/core"
	"github.com/jonesrussell/gimbal/internal/ecs/managers"
	resources "github.com/jonesrussell/gimbal/internal/ecs/managers/resource"
)

// EnemySystem manages enemy spawning, movement, and behavior
type EnemySystem struct {
	world         donburi.World
	gameConfig    *config.GameConfig
	spawnTimer    float64
	spawnInterval float64
	resourceMgr   *resources.ResourceManager
	logger        common.Logger

	// Wave management
	waveManager *WaveManager

	// Boss spawning
	bossSpawnTimer float64
	bossSpawned    bool
	bossConfig     *managers.BossConfig // Current level's boss configuration

	// Enemy sprites cache
	enemySprites map[EnemyType]*ebiten.Image

	// Enemy configurations loaded from JSON
	enemyConfigs map[EnemyType]EnemyTypeData
}

// NewEnemySystem creates a new enemy management system with the provided dependencies
func NewEnemySystem(
	world donburi.World,
	gameConfig *config.GameConfig,
	resourceMgr *resources.ResourceManager,
	logger common.Logger,
) *EnemySystem {
	es := &EnemySystem{
		world:         world,
		gameConfig:    gameConfig,
		spawnTimer:    0,
		spawnInterval: DefaultSpawnIntervalSeconds,
		resourceMgr:   resourceMgr,
		logger:        logger,
		enemySprites:  make(map[EnemyType]*ebiten.Image),
		enemyConfigs:  make(map[EnemyType]EnemyTypeData),
	}

	// Initialize wave manager
	es.waveManager = NewWaveManager(world, logger)

	return es
}

func (es *EnemySystem) Update(ctx context.Context, deltaTime float64) error {
	// Check for cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Update wave manager
	es.waveManager.Update(deltaTime)

	// Check if we need to start a new wave (but not if we're waiting for level start or inter-wave delay)
	if es.waveManager.GetCurrentWave() == nil &&
		es.waveManager.HasMoreWaves() &&
		!es.waveManager.IsWaiting() &&
		!es.waveManager.IsWaitingForLevelStart() {
		es.waveManager.StartNextWave()
	}

	// Spawn enemies from current wave
	if es.waveManager.ShouldSpawnEnemy(deltaTime) {
		wave := es.waveManager.GetCurrentWave()
		if wave != nil {
			es.spawnWaveEnemy(ctx, wave)
			es.waveManager.MarkEnemySpawned()
		}
	}

	// Check if wave is complete and start next
	es.handleWaveCompletion(ctx, deltaTime)

	// Handle boss spawning after all waves are complete
	es.handleBossSpawning(ctx, deltaTime)

	// Update enemy movement (including boss)
	es.updateEnemies(deltaTime)
	es.UpdateBossMovement(deltaTime)

	return nil
}

// handleWaveCompletion handles wave completion and advances to next wave
func (es *EnemySystem) handleWaveCompletion(ctx context.Context, deltaTime float64) {
	if !es.waveManager.IsWaveComplete() {
		return
	}

	es.waveManager.CompleteWave()
	if es.waveManager.HasMoreWaves() {
		es.waveManager.StartNextWave()
		return
	}

	// All waves complete - boss will be handled in handleBossSpawning
	if es.bossConfig != nil && es.bossConfig.Enabled && !es.bossSpawned {
		es.logger.Debug("All waves complete, boss will spawn soon", "spawn_delay", es.bossConfig.SpawnDelay)
	}
}

// handleBossSpawning handles boss spawn timer and spawning
func (es *EnemySystem) handleBossSpawning(ctx context.Context, deltaTime float64) {
	// Only spawn boss if all waves are complete and boss hasn't been spawned yet
	if es.waveManager.HasMoreWaves() {
		return // Still have waves to complete
	}

	if es.bossSpawned {
		return // Boss already spawned
	}

	if es.bossConfig == nil || !es.bossConfig.Enabled {
		return // No boss configured for this level
	}

	// Increment boss spawn timer
	es.bossSpawnTimer += deltaTime
	spawnDelay := es.bossConfig.SpawnDelay
	if spawnDelay <= 0 {
		spawnDelay = BossSpawnDelay // Fallback to default
	}

	if es.bossSpawnTimer >= spawnDelay {
		es.SpawnBoss(ctx)
		es.bossSpawned = true
		es.logger.Debug("Boss spawned", "delay", es.bossSpawnTimer)
	}
}

// IsBossActive checks if there's an active boss
func (es *EnemySystem) IsBossActive() bool {
	if es.bossSpawned {
		// Check if boss still exists using EnemyTypeID component
		count := 0
		query.NewQuery(
			filter.And(
				filter.Contains(core.EnemyTag),
				filter.Contains(core.EnemyTypeID),
			),
		).Each(es.world, func(entry *donburi.Entry) {
			typeID := core.EnemyTypeID.Get(entry)
			if EnemyType(*typeID) == EnemyTypeBoss {
				count++
			}
		})
		return count > 0
	}
	return false
}

// WasBossSpawned returns true if boss was spawned (even if now killed)
func (es *EnemySystem) WasBossSpawned() bool {
	return es.bossSpawned
}

// GetWaveManager returns the wave manager
func (es *EnemySystem) GetWaveManager() *WaveManager {
	return es.waveManager
}
