package game

import (
	"context"
	"fmt"

	"github.com/jonesrussell/gimbal/assets"
	"github.com/jonesrussell/gimbal/internal/ecs/core"
	"github.com/jonesrussell/gimbal/internal/ecs/debug"
	"github.com/jonesrussell/gimbal/internal/ecs/events"
	"github.com/jonesrussell/gimbal/internal/ecs/managers"
	resources "github.com/jonesrussell/gimbal/internal/ecs/managers/resource"
	"github.com/jonesrussell/gimbal/internal/ecs/systems/collision"
	"github.com/jonesrussell/gimbal/internal/ecs/systems/gyruss"
	"github.com/jonesrussell/gimbal/internal/ecs/systems/health"
	"github.com/jonesrussell/gimbal/internal/ecs/systems/movement"
	"github.com/jonesrussell/gimbal/internal/ecs/systems/weapon"
)

// createCoreSystems creates core ECS systems
func (g *ECSGame) createCoreSystems(ctx context.Context) error {
	g.eventSystem = events.NewEventSystem(g.world)
	g.logger.Debug("Event system created")
	g.resourceManager = resources.NewResourceManager(ctx, g.logger)
	g.logger.Debug("Resource manager created")
	g.stateManager = NewGameStateManager(g.eventSystem, g.logger)
	g.scoreManager = managers.NewScoreManager(10000)
	g.levelManager = managers.NewLevelManager(g.logger)
	g.logger.Debug("Level manager created")

	return nil
}

// createGameplaySystems creates gameplay ECS systems
func (g *ECSGame) createGameplaySystems(ctx context.Context) error {
	// Load player config
	playerConfig, err := managers.LoadPlayerConfig(ctx, g.logger)
	if err != nil {
		return fmt.Errorf("failed to load player config: %w", err)
	}
	g.playerConfig = playerConfig
	g.logger.Debug("Player config loaded", "health", playerConfig.Health, "size", playerConfig.Size)

	if err := g.createBasicSystems(); err != nil {
		return err
	}

	g.createRemainingSystems(ctx)

	// Load initial stage into Gyruss system
	if err := g.gyrussSystem.LoadStage(1); err != nil {
		g.logger.Warn("Failed to load initial stage", "error", err)
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
	})
	g.logger.Debug("Gyruss system created")

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
