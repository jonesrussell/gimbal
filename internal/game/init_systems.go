package game

import (
	"context"
	"fmt"

	"github.com/jonesrussell/gimbal/assets"
	"github.com/jonesrussell/gimbal/internal/dbg"
	"github.com/jonesrussell/gimbal/internal/ecs/core"
	"github.com/jonesrussell/gimbal/internal/ecs/debug"
	"github.com/jonesrussell/gimbal/internal/ecs/events"
	"github.com/jonesrussell/gimbal/internal/ecs/managers"
	resources "github.com/jonesrussell/gimbal/internal/ecs/managers/resource"
	"github.com/jonesrussell/gimbal/internal/ecs/systems/collision"
	"github.com/jonesrussell/gimbal/internal/ecs/systems/gyruss"
	"github.com/jonesrussell/gimbal/internal/ecs/systems/health"
	"github.com/jonesrussell/gimbal/internal/ecs/systems/movement"
	"github.com/jonesrussell/gimbal/internal/ecs/systems/stage"
	"github.com/jonesrussell/gimbal/internal/ecs/systems/weapon"
)

// createCoreSystems creates core ECS systems
func (g *ECSGame) createCoreSystems(ctx context.Context) error {
	g.eventSystem = events.NewEventSystem(g.world)
	dbg.Log(dbg.World, "EventSystem created at %p", g.eventSystem)
	g.logger.Debug("Event system created")
	g.resourceManager = resources.NewResourceManager(ctx, g.logger)
	g.logger.Debug("Resource manager created")
	g.stateManager = NewGameStateManager(g.eventSystem, g.logger)
	g.scoreManager = managers.NewScoreManager(10000)
	g.levelManager = managers.NewLevelManager(g.logger)

	// Wire up level manager to emit events
	g.levelManager.SetEventEmitter(g.eventSystem)

	g.logger.Debug("Level manager created")

	return nil
}

// createGameplaySystems creates gameplay ECS systems
func (g *ECSGame) createGameplaySystems(ctx context.Context) error {
	// Load player config
	playerConfig, loadErr := managers.LoadPlayerConfig(ctx, g.logger)
	if loadErr != nil {
		return fmt.Errorf("failed to load player config: %w", loadErr)
	}
	g.playerConfig = playerConfig
	g.logger.Debug("Player config loaded", "health", playerConfig.Health, "size", playerConfig.Size)

	if basicErr := g.createBasicSystems(); basicErr != nil {
		return basicErr
	}

	g.createRemainingSystems(ctx)

	// Load initial stage via stage state machine (delegates to GyrussSystem)
	if stageErr := g.stageStateMachine.LoadStage(1); stageErr != nil {
		g.logger.Warn("Failed to load initial stage", "error", stageErr)
	}

	return nil
}

// createBasicSystems creates health, movement, and Gyruss systems
func (g *ECSGame) createBasicSystems() error {
	g.healthSystem = health.NewHealthSystem(g.world, g.config, g.eventSystem, g.stateManager, g.logger)
	g.logger.Debug("Health system created")

	g.movementSystem = movement.NewMovementSystem(g.world, g.config, g.logger, g.inputHandler)
	g.logger.Debug("Movement system created")

	// Create Gyruss system (replaces EnemySystem)
	g.gyrussSystem = gyruss.NewGyrussSystem(&gyruss.GyrussSystemConfig{
		World:       g.world,
		GameConfig:  g.config,
		ResourceMgr: g.resourceManager,
		Logger:      g.logger,
		AssetsFS:    assets.Assets,
		EventSystem: g.eventSystem,
	})
	g.logger.Debug("Gyruss system created")

	// Create stage state machine (authoritative owner of wave/boss/stage progression)
	g.stageStateMachine = stage.NewStageStateMachine(&stage.Config{
		EventSystem:  g.eventSystem,
		WaveManager:  g.gyrussSystem.GetWaveManager(),
		GyrussSystem: g.gyrussSystem,
		Logger:       g.logger,
	})
	dbg.Log(dbg.World, "StageStateMachine using EventSystem %p", g.eventSystem)
	g.logger.Debug("Stage state machine created")

	return nil
}

// createRemainingSystems creates weapon and collision systems
func (g *ECSGame) createRemainingSystems(ctx context.Context) {
	g.weaponSystem = weapon.NewWeaponSystem(g.world, g.config)
	g.logger.Debug("Weapon system created")

	g.collisionSystem = collision.NewCollisionSystem(&collision.CollisionSystemConfig{
		World:        g.world,
		Config:       g.config,
		HealthSystem: g.healthSystem,
		EventSystem:  g.eventSystem,
		ScoreManager: g.scoreManager,
		EnemySystem:  g.gyrussSystem, // GyrussSystem implements EnemySystemInterface
		Logger:       g.logger,
	})
	g.logger.Debug("Collision system created")
}

// registerAllSystems registers and initializes all systems
func (g *ECSGame) registerAllSystems(ctx context.Context) error {
	g.renderOptimizer = core.NewRenderOptimizer(g.config)
	g.imagePool = core.NewImagePool(g.logger)
	g.perfMonitor = debug.NewPerformanceMonitor(g.logger)

	// Initialize rendering debugger with default font
	font, err := g.resourceManager.GetDefaultFont(ctx)
	if err != nil {
		g.logger.Warn("Failed to get font for debugger, debug overlay disabled", "error", err)
	} else {
		g.renderDebugger = debug.NewRenderingDebugger(font, g.logger)
		g.logger.Debug("Rendering debugger initialized")
	}

	g.logger.Debug("Performance optimizations initialized")
	return nil
}
