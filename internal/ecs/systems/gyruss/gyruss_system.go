package gyruss

import (
	"context"
	"embed"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/config"
	"github.com/jonesrussell/gimbal/internal/ecs/core"
	"github.com/jonesrussell/gimbal/internal/ecs/managers"
	resources "github.com/jonesrussell/gimbal/internal/ecs/managers/resource"
	"github.com/jonesrussell/gimbal/internal/ecs/systems/animation"
	"github.com/jonesrussell/gimbal/internal/ecs/systems/attack"
	"github.com/jonesrussell/gimbal/internal/ecs/systems/behavior"
	"github.com/jonesrussell/gimbal/internal/ecs/systems/enemy"
	"github.com/jonesrussell/gimbal/internal/ecs/systems/fire"
	"github.com/jonesrussell/gimbal/internal/ecs/systems/path"
	"github.com/jonesrussell/gimbal/internal/ecs/systems/powerup"
)

// GyrussSystem is the main system that coordinates all Gyruss-style gameplay
type GyrussSystem struct {
	world       donburi.World
	gameConfig  *config.GameConfig
	resourceMgr *resources.ResourceManager
	logger      common.Logger

	// Stage management
	stageLoader *managers.StageLoader
	waveManager *enemy.GyrussWaveManager
	spawner     *enemy.GyrussSpawner

	// Subsystems
	pathSystem     *path.PathSystem
	scaleSystem    *animation.ScaleAnimationSystem
	behaviorSystem *behavior.BehaviorSystem
	attackSystem   *attack.AttackSystem
	fireSystem     *fire.FirePatternSystem
	powerUpSystem  *powerup.PowerUpSystem

	// State
	currentStage int
	bossSpawned  bool
	bossTimer    float64
}

// GyrussSystemConfig holds configuration for creating a GyrussSystem
type GyrussSystemConfig struct {
	World       donburi.World
	GameConfig  *config.GameConfig
	ResourceMgr *resources.ResourceManager
	Logger      common.Logger
	AssetsFS    embed.FS
}

// NewGyrussSystem creates a new Gyruss gameplay system
func NewGyrussSystem(cfg *GyrussSystemConfig) *GyrussSystem {
	gs := &GyrussSystem{
		world:        cfg.World,
		gameConfig:   cfg.GameConfig,
		resourceMgr:  cfg.ResourceMgr,
		logger:       cfg.Logger,
		currentStage: 1,
		bossSpawned:  false,
		bossTimer:    0,
	}

	// Create stage loader
	gs.stageLoader = managers.NewStageLoader(cfg.Logger, cfg.AssetsFS)

	// Create wave manager
	gs.waveManager = enemy.NewGyrussWaveManager(cfg.World, cfg.Logger)

	// Create spawner
	gs.spawner = enemy.NewGyrussSpawner(cfg.World, cfg.GameConfig, cfg.ResourceMgr, cfg.Logger)

	// Create subsystems
	gs.createSubsystems(cfg)

	return gs
}

// createSubsystems creates all the Gyruss subsystems
func (gs *GyrussSystem) createSubsystems(cfg *GyrussSystemConfig) {
	// Path system for entry animations
	gs.pathSystem = path.NewPathSystem(cfg.World, cfg.GameConfig, cfg.Logger)

	// Scale animation system
	gs.scaleSystem = animation.NewScaleAnimationSystem(cfg.World, cfg.Logger)

	// Behavior state machine
	gs.behaviorSystem = behavior.NewBehaviorSystem(cfg.World, cfg.GameConfig, cfg.Logger)

	// Attack pattern system
	gs.attackSystem = attack.NewAttackSystem(cfg.World, cfg.GameConfig, cfg.Logger)

	// Fire pattern system
	gs.fireSystem = fire.NewFirePatternSystem(cfg.World, cfg.GameConfig, cfg.Logger)

	// Power-up system
	gs.powerUpSystem = powerup.NewPowerUpSystem(cfg.World, cfg.GameConfig, cfg.Logger)

	gs.logger.Debug("Gyruss subsystems created")
}

// LoadStage loads a stage by number
func (gs *GyrussSystem) LoadStage(stageNumber int) error {
	stageConfig, err := gs.stageLoader.LoadStage(stageNumber)
	if err != nil {
		return err
	}

	gs.currentStage = stageNumber
	gs.bossSpawned = false
	gs.bossTimer = 0

	// Load into wave manager
	gs.waveManager.LoadStage(stageConfig)

	gs.logger.Info("Gyruss stage loaded",
		"stage", stageNumber,
		"name", stageConfig.Metadata.Name,
		"waves", len(stageConfig.Waves))

	return nil
}

// Update updates all Gyruss systems
func (gs *GyrussSystem) Update(ctx context.Context, deltaTime float64) error {
	// Check for cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Update wave manager
	gs.waveManager.Update(deltaTime)

	// Handle enemy spawning
	gs.handleSpawning(ctx)

	// Handle boss spawning
	gs.handleBossSpawning(ctx, deltaTime)

	// Update path system (entry animations)
	if err := gs.pathSystem.Update(ctx, deltaTime); err != nil {
		return err
	}

	// Update scale animations
	if err := gs.scaleSystem.Update(ctx, deltaTime); err != nil {
		return err
	}

	// Update behavior state machine
	if err := gs.behaviorSystem.Update(ctx, deltaTime); err != nil {
		return err
	}

	// Update attack patterns
	if err := gs.attackSystem.Update(ctx, deltaTime); err != nil {
		return err
	}

	// Update fire patterns
	if err := gs.fireSystem.Update(ctx, deltaTime); err != nil {
		return err
	}

	// Update power-ups
	if err := gs.powerUpSystem.Update(ctx, deltaTime); err != nil {
		return err
	}

	return nil
}

// handleSpawning handles enemy spawning from waves
func (gs *GyrussSystem) handleSpawning(ctx context.Context) {
	groupConfig, shouldSpawn := gs.waveManager.ShouldSpawnEnemy()
	if !shouldSpawn || groupConfig == nil {
		return
	}

	// Get current spawn index
	spawnIndex := gs.getSpawnIndexForGroup()

	// Spawn enemy
	gs.spawner.SpawnEnemy(ctx, groupConfig, spawnIndex)

	// Mark as spawned
	gs.waveManager.MarkEnemySpawned()
}

// handleBossSpawning handles boss spawning
func (gs *GyrussSystem) handleBossSpawning(ctx context.Context, deltaTime float64) {
	if gs.bossSpawned {
		return
	}

	if !gs.waveManager.IsBossTriggered() {
		return
	}

	bossConfig := gs.waveManager.GetBossConfig()
	if bossConfig == nil || !bossConfig.Enabled {
		return
	}

	// Wait for spawn delay
	gs.bossTimer += deltaTime
	if gs.bossTimer < bossConfig.SpawnDelay {
		return
	}

	// Spawn boss
	gs.spawner.SpawnBoss(ctx, bossConfig)
	gs.bossSpawned = true

	gs.logger.Info("Gyruss boss spawned",
		"stage", gs.currentStage,
		"boss_type", bossConfig.BossType)
}

// getSpawnIndexForGroup calculates the spawn index based on wave manager state
func (gs *GyrussSystem) getSpawnIndexForGroup() int {
	// This is a simplified calculation - the wave manager tracks internally
	return 0
}

// OnEnemyDestroyed is called when an enemy is destroyed
func (gs *GyrussSystem) OnEnemyDestroyed(position common.Point, isPowerUpTrigger bool) {
	if isPowerUpTrigger {
		gs.powerUpSystem.TrySpawnPowerUp(position)
	}
}

// HasDoubleShot returns whether player has double shot active
func (gs *GyrussSystem) HasDoubleShot() bool {
	return gs.powerUpSystem.HasDoubleShot()
}

// IsBossActive returns whether boss is currently active
func (gs *GyrussSystem) IsBossActive() bool {
	return gs.bossSpawned && !gs.IsBossDefeated()
}

// WasBossSpawned returns whether the boss has been spawned this stage
func (gs *GyrussSystem) WasBossSpawned() bool {
	return gs.bossSpawned
}

// IsBossDefeated returns whether boss has been defeated
func (gs *GyrussSystem) IsBossDefeated() bool {
	// Check if boss entity still exists
	bossExists := false
	query.NewQuery(
		filter.And(
			filter.Contains(core.EnemyTag),
			filter.Contains(core.EnemyTypeID),
		),
	).Each(gs.world, func(entry *donburi.Entry) {
		typeID := core.EnemyTypeID.Get(entry)
		if enemy.EnemyType(*typeID) == enemy.EnemyTypeBoss {
			bossExists = true
		}
	})
	return !bossExists
}

// DestroyEnemy destroys an enemy and returns points - implements EnemySystemInterface
func (gs *GyrussSystem) DestroyEnemy(entity donburi.Entity) int {
	if !gs.world.Valid(entity) {
		return 0
	}

	entry := gs.world.Entry(entity)
	points := 100 // Default points

	// Check if it's a boss for more points
	if entry.HasComponent(core.EnemyTypeID) {
		typeID := core.EnemyTypeID.Get(entry)
		switch enemy.EnemyType(*typeID) {
		case enemy.EnemyTypeBoss:
			// Get boss points from stage config
			bossConfig := gs.waveManager.GetBossConfig()
			if bossConfig != nil {
				points = bossConfig.Points
			} else {
				points = 1000
			}
		case enemy.EnemyTypeHeavy:
			points = 200
		default:
			points = 100
		}
	}

	// Get position for power-up spawn
	var position common.Point
	if entry.HasComponent(core.Position) {
		pos := core.Position.Get(entry)
		position = *pos
	}

	// Try to spawn power-up
	gs.powerUpSystem.TrySpawnPowerUp(position)

	// Remove the entity
	gs.world.Remove(entity)

	gs.logger.Debug("Enemy destroyed", "entity", entity, "points", points)
	return points
}

// IsStageComplete returns whether the current stage is complete
func (gs *GyrussSystem) IsStageComplete() bool {
	return gs.bossSpawned && gs.IsBossDefeated()
}

// GetCurrentStage returns the current stage number
func (gs *GyrussSystem) GetCurrentStage() int {
	return gs.currentStage
}

// GetWaveManager returns the wave manager
func (gs *GyrussSystem) GetWaveManager() *enemy.GyrussWaveManager {
	return gs.waveManager
}

// GetPowerUpSystem returns the power-up system
func (gs *GyrussSystem) GetPowerUpSystem() *powerup.PowerUpSystem {
	return gs.powerUpSystem
}

// Reset resets the Gyruss system state for a new game
func (gs *GyrussSystem) Reset() {
	gs.bossSpawned = false
	gs.bossTimer = 0
	gs.waveManager.Reset()
	gs.logger.Debug("Gyruss system reset")
}

// LoadNextStage advances to and loads the next stage
func (gs *GyrussSystem) LoadNextStage() error {
	nextStage := gs.currentStage + 1
	return gs.LoadStage(nextStage)
}
