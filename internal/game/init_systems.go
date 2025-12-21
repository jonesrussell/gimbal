package game

import (
	"context"
	"fmt"

	"github.com/jonesrussell/gimbal/internal/ecs/debug"
	"github.com/jonesrussell/gimbal/internal/ecs/events"
	"github.com/jonesrussell/gimbal/internal/ecs/managers"
	resources "github.com/jonesrussell/gimbal/internal/ecs/managers/resource"
	"github.com/jonesrussell/gimbal/internal/ecs/systems/collision"
	"github.com/jonesrussell/gimbal/internal/ecs/systems/enemy"
	"github.com/jonesrussell/gimbal/internal/ecs/systems/health"
	"github.com/jonesrussell/gimbal/internal/ecs/systems/movement"
	"github.com/jonesrussell/gimbal/internal/ecs/systems/weapon"
	"github.com/jonesrussell/gimbal/internal/ecs/core"
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

	// Load level definitions
	levelDefs := managers.GetDefaultLevelDefinitions()
	if err := g.levelManager.LoadLevels(levelDefs); err != nil {
		return fmt.Errorf("failed to load levels: %w", err)
	}
	g.logger.Debug("Level definitions loaded", "count", len(levelDefs))

	return nil
}

// createGameplaySystems creates gameplay ECS systems
func (g *ECSGame) createGameplaySystems(ctx context.Context) error {
	// Load entity configurations from JSON (no fallback - errors if missing)
	playerConfig, err := managers.LoadPlayerConfig(ctx, g.logger)
	if err != nil {
		return fmt.Errorf("failed to load player config: %w", err)
	}
	g.playerConfig = playerConfig
	g.logger.Debug("Player config loaded", "health", playerConfig.Health, "size", playerConfig.Size)

	enemyConfigs, err := managers.LoadEnemyConfigs(ctx, g.logger)
	if err != nil {
		return fmt.Errorf("failed to load enemy configs: %w", err)
	}

	// Convert enemy configs to map[EnemyType]EnemyTypeData
	enemyConfigMap := make(map[enemy.EnemyType]enemy.EnemyTypeData)
	for _, enemyConfig := range enemyConfigs.EnemyTypes {
		enemyType, typeErr := enemy.GetEnemyTypeFromString(enemyConfig.Type)
		if typeErr != nil {
			return fmt.Errorf("invalid enemy type '%s': %w", enemyConfig.Type, typeErr)
		}
		enemyData, convertErr := enemy.ConvertEnemyTypeConfig(&enemyConfig, enemyType)
		if convertErr != nil {
			return fmt.Errorf("failed to convert enemy config for type '%s': %w", enemyConfig.Type, err)
		}
		enemyConfigMap[enemyType] = enemyData
	}
	g.logger.Debug("Enemy configs loaded and converted", "count", len(enemyConfigMap))

	g.healthSystem = health.NewHealthSystem(g.world, g.config, g.eventSystem, g.stateManager, g.logger)
	g.logger.Debug("Health system created")
	g.movementSystem = movement.NewMovementSystem(g.world, g.config, g.logger, g.inputHandler)
	g.logger.Debug("Movement system created")
	g.enemySystem = enemy.NewEnemySystem(g.world, g.config, g.resourceManager, g.logger)

	// Load enemy configs into enemy system
	g.enemySystem.LoadEnemyConfigs(enemyConfigMap)
	g.logger.Debug("Enemy configs loaded into enemy system")

	// Load initial level configuration into enemy system
	levelConfig := g.levelManager.GetCurrentLevelConfig()
	if levelConfig != nil {
		enemyWaves := convertWaveConfigs(levelConfig.Waves)
		g.enemySystem.LoadLevelConfig(enemyWaves, &levelConfig.Boss)
		g.logger.Debug("Initial level config loaded into enemy system",
			"level", levelConfig.LevelNumber,
			"waves", len(levelConfig.Waves))
	}

	g.enemyWeaponSystem = enemy.NewEnemyWeaponSystem(g.world, g.config, g.logger, g.enemySystem)
	g.logger.Debug("Enemy weapon system created")
	g.weaponSystem = weapon.NewWeaponSystem(g.world, g.config)
	g.collisionSystem = collision.NewCollisionSystem(&collision.CollisionSystemConfig{
		World:        g.world,
		Config:       g.config,
		HealthSystem: g.healthSystem,
		EventSystem:  g.eventSystem,
		ScoreManager: g.scoreManager,
		EnemySystem:  g.enemySystem,
		Logger:       g.logger,
	})
	return nil
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

