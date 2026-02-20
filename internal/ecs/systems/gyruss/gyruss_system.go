package gyruss

import (
	"context"
	"embed"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/dbg"
	"github.com/jonesrussell/gimbal/internal/config"
	"github.com/jonesrussell/gimbal/internal/ecs/core"
	"github.com/jonesrussell/gimbal/internal/ecs/events"
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

	// Events (for BossDefeated)
	eventSystem *events.EventSystem

	// State
	currentStage int
}

// GyrussSystemConfig holds configuration for creating a GyrussSystem
type GyrussSystemConfig struct {
	World       donburi.World
	GameConfig  *config.GameConfig
	ResourceMgr *resources.ResourceManager
	Logger      common.Logger
	AssetsFS    embed.FS
	EventSystem *events.EventSystem
}

// NewGyrussSystem creates a new Gyruss gameplay system
func NewGyrussSystem(cfg *GyrussSystemConfig) *GyrussSystem {
	gs := &GyrussSystem{
		world:        cfg.World,
		gameConfig:   cfg.GameConfig,
		resourceMgr:  cfg.ResourceMgr,
		logger:       cfg.Logger,
		eventSystem:  cfg.EventSystem,
		currentStage: 1,
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

	// Load into wave manager
	gs.waveManager.LoadStage(stageConfig)

	gs.logger.Info("Gyruss stage loaded",
		"stage", stageNumber,
		"name", stageConfig.Metadata.Name,
		"waves", len(stageConfig.Waves))

	return nil
}

// offScreenEnemyMargin is the margin past screen bounds beyond which enemies are removed (matches behavior retreating_state)
const offScreenEnemyMargin = 100.0

// removeOffScreenEnemies removes enemy entities that have moved off-screen (e.g. after retreat), so wave completion can trigger
func (gs *GyrussSystem) removeOffScreenEnemies() {
	w := float64(gs.gameConfig.ScreenSize.Width)
	h := float64(gs.gameConfig.ScreenSize.Height)
	var toRemove []donburi.Entity
	query.NewQuery(
		filter.And(
			filter.Contains(core.EnemyTag),
			filter.Contains(core.Position),
		),
	).Each(gs.world, func(entry *donburi.Entry) {
		if entry.HasComponent(core.EnemyTypeID) {
			typeID := core.EnemyTypeID.Get(entry)
			if enemy.EnemyType(*typeID) == enemy.EnemyTypeBoss {
				return
			}
		}
		pos := core.Position.Get(entry)
		if pos.X < -offScreenEnemyMargin || pos.X > w+offScreenEnemyMargin ||
			pos.Y < -offScreenEnemyMargin || pos.Y > h+offScreenEnemyMargin {
			toRemove = append(toRemove, entry.Entity())
		}
	})
	for _, e := range toRemove {
		gs.world.Remove(e)
	}
}

// Update updates all Gyruss systems
func (gs *GyrussSystem) Update(ctx context.Context, deltaTime float64) error {
	// Check for cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Remove enemies that retreated off-screen so wave completion can trigger without waiting for timeout
	gs.removeOffScreenEnemies()

	// Update wave manager
	gs.waveManager.Update(deltaTime)

	// Handle enemy spawning
	gs.handleSpawning(ctx)

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

	// Get current spawn index (before MarkEnemySpawned) for orbit angle distribution
	spawnIndex := gs.waveManager.GetCurrentGroupSpawnIndex()

	// Spawn enemy
	gs.spawner.SpawnEnemy(ctx, groupConfig, spawnIndex)

	// Mark as spawned
	gs.waveManager.MarkEnemySpawned()
}

// SpawnBoss spawns the boss immediately when called (by StageStateMachine after delay)
func (gs *GyrussSystem) SpawnBoss(ctx context.Context) {
	bossConfig := gs.waveManager.GetBossConfig()
	if bossConfig == nil || !bossConfig.Enabled {
		return
	}
	gs.spawner.SpawnBoss(ctx, bossConfig)
	dbg.Log(dbg.Spawn, "boss spawned")
	gs.logger.Info("Gyruss boss spawned",
		"stage", gs.currentStage,
		"boss_type", bossConfig.BossType)
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

// DestroyEnemy destroys an enemy and returns points - implements EnemySystemInterface
func (gs *GyrussSystem) DestroyEnemy(entity donburi.Entity) int {
	if !gs.world.Valid(entity) {
		return 0
	}

	entry := gs.world.Entry(entity)
	points := 100 // Default points

	// Check if it's a boss for more points and emit BossDefeated before removing
	isBoss := false
	if entry.HasComponent(core.EnemyTypeID) {
		typeID := core.EnemyTypeID.Get(entry)
		switch enemy.EnemyType(*typeID) {
		case enemy.EnemyTypeBoss:
			isBoss = true
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

	// Emit BossDefeated before removing so StageStateMachine can transition
	if isBoss && gs.eventSystem != nil {
		gs.eventSystem.EmitBossDefeated()
		dbg.Log(dbg.Event, "EmitBossDefeated (world=%p eventSystem=%p)", gs.world, gs.eventSystem)
	}

	// Remove the entity
	gs.world.Remove(entity)

	gs.logger.Debug("Enemy destroyed", "entity", entity, "points", points)
	return points
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
	gs.waveManager.Reset()
	gs.logger.Debug("Gyruss system reset")
}

// LoadNextStage advances to and loads the next stage
func (gs *GyrussSystem) LoadNextStage() error {
	nextStage := gs.currentStage + 1
	return gs.LoadStage(nextStage)
}
