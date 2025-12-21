package game

import (
	"context"
	"fmt"

	"github.com/yohamta/donburi"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/config"
	"github.com/jonesrussell/gimbal/internal/ecs/core"
	"github.com/jonesrussell/gimbal/internal/ecs/debug"
	"github.com/jonesrussell/gimbal/internal/ecs/events"
	"github.com/jonesrussell/gimbal/internal/ecs/managers"
	resources "github.com/jonesrussell/gimbal/internal/ecs/managers/resource"
	"github.com/jonesrussell/gimbal/internal/ecs/systems/collision"
	"github.com/jonesrussell/gimbal/internal/ecs/systems/enemy"
	"github.com/jonesrussell/gimbal/internal/ecs/systems/health"
	"github.com/jonesrussell/gimbal/internal/ecs/systems/movement"
	"github.com/jonesrussell/gimbal/internal/ecs/systems/weapon"
	"github.com/jonesrussell/gimbal/internal/errors"
	"github.com/jonesrussell/gimbal/internal/scenes"
	uicore "github.com/jonesrussell/gimbal/internal/ui/core"
	"github.com/jonesrussell/gimbal/internal/ui/responsive"
)

// validateGameConfig validates the provided game configuration
func (g *ECSGame) validateGameConfig(cfg *config.GameConfig) error {
	if cfg == nil {
		return errors.NewGameError(errors.ConfigMissing, "config cannot be nil")
	}
	return nil
}

// createGameEntities creates all game entities
func (g *ECSGame) createGameEntities(ctx context.Context) error {
	return g.createEntities(ctx)
}

// setupInitialScene sets up the initial scene
func (g *ECSGame) setupInitialScene(ctx context.Context) error {
	font, err := g.resourceManager.GetDefaultFont(ctx)
	if err != nil {
		return errors.NewGameErrorWithCause(errors.AssetLoadFailed, "failed to get default font", err)
	}
	g.sceneManager = scenes.NewSceneManager(&scenes.SceneManagerConfig{
		World:        g.world,
		Config:       g.config,
		Logger:       g.logger,
		InputHandler: g.inputHandler,
		Font:         font,
		ScoreManager: g.scoreManager,
		ResourceMgr:  g.resourceManager,
	})
	g.sceneManager.SetResumeCallback(func() {
		g.stateManager.SetPaused(false)
	})
	g.sceneManager.SetHealthSystem(g.healthSystem)
	g.sceneManager.SetLevelManager(g.levelManager)
	g.sceneManager.SetRenderOptimizer(g.renderOptimizer)
	g.sceneManager.SetImagePool(g.imagePool)
	if sceneErr := g.sceneManager.SetInitialScene(scenes.SceneStudioIntro); sceneErr != nil {
		return errors.NewGameErrorWithCause(errors.SystemInitFailed, "failed to set initial scene", sceneErr)
	}
	return nil
}

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
		enemyData, convertErr := enemy.ConvertEnemyTypeConfig(enemyConfig, enemyType)
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

// initializeSystems initializes all game systems
func (g *ECSGame) initializeSystems(ctx context.Context) error {
	if err := g.validateGameConfig(g.config); err != nil {
		return err
	}
	if err := g.createCoreSystems(ctx); err != nil {
		return err
	}
	if err := g.createGameplaySystems(ctx); err != nil {
		return err
	}
	if err := g.registerAllSystems(ctx); err != nil {
		return err
	}
	if err := g.loadAssets(ctx); err != nil {
		return errors.NewGameErrorWithCause(errors.AssetLoadFailed, "failed to load assets", err)
	}
	return nil
}

// initializeUI sets up the game UI
func (g *ECSGame) initializeUI(ctx context.Context) error {
	font, err := g.resourceManager.GetDefaultFont(ctx)
	if err != nil {
		return errors.NewGameErrorWithCause(errors.AssetLoadFailed, "failed to get default font", err)
	}

	heartSprite, err := g.resourceManager.GetUISprite(ctx, "heart", uicore.HeartIconSize)
	if err != nil {
		return fmt.Errorf("failed to load heart sprite: %w", err)
	}

	ammoSprite, err := g.resourceManager.GetUISprite(ctx, "ammo", uicore.AmmoIconSize)
	if err != nil {
		g.logger.Warn("Failed to load ammo sprite, using fallback", "error", err)
		ammoSprite = nil // Will use fallback in UI
	}

	uiConfig := &responsive.Config{
		Font:        font,
		HeartSprite: heartSprite,
		AmmoSprite:  ammoSprite,
	}

	gameUI, err := responsive.NewResponsiveUI(uiConfig)
	if err != nil {
		return fmt.Errorf("failed to create game UI: %w", err)
	}
	g.ui = gameUI
	return nil
}

// finalizeInitialization completes the game initialization
func (g *ECSGame) finalizeInitialization(ctx context.Context) error {
	err := g.createGameEntities(ctx)
	if err != nil {
		return errors.NewGameErrorWithCause(errors.EntityInvalid, "failed to create entities", err)
	}

	err = g.setupInitialScene(ctx)
	if err != nil {
		return err
	}

	g.setupEventSubscriptions()
	return nil
}

// NewECSGame creates a new ECS-based game instance with dependency-injected UI
func NewECSGame(
	ctx context.Context,
	gameConfig *config.GameConfig,
	logger common.Logger,
	inputHandler common.GameInputHandler,
) (*ECSGame, error) {
	logger.Debug("Creating new ECS game instance",
		"screen_size", gameConfig.ScreenSize,
		"player_size", gameConfig.PlayerSize,
		"num_stars", gameConfig.NumStars,
	)

	// Create ECS world
	world := donburi.NewWorld()

	// Create game context for lifecycle management
	gameCtx, cancel := context.WithCancel(ctx)

	game := &ECSGame{
		world:        world,
		config:       gameConfig,
		inputHandler: inputHandler,
		logger:       logger,
		ctx:          gameCtx,
		cancel:       cancel,
	}

	if err := game.initializeSystems(ctx); err != nil {
		return nil, err
	}

	if err := game.initializeUI(ctx); err != nil {
		return nil, err
	}

	if err := game.finalizeInitialization(ctx); err != nil {
		return nil, err
	}

	return game, nil
}

// loadAssets loads and prepares game assets
func (g *ECSGame) loadAssets(ctx context.Context) error {
	// Load all sprites through resource manager
	if err := g.resourceManager.LoadAllSprites(ctx); err != nil {
		return errors.NewGameErrorWithCause(errors.AssetLoadFailed, "failed to load sprites", err)
	}

	g.logger.Debug("Assets loaded successfully", "resource_count", g.resourceManager.GetResourceCount())
	return nil
}

// createEntities creates all game entities
func (g *ECSGame) createEntities(ctx context.Context) error {
	// Get sprites from resource manager
	playerSprite, ok := g.resourceManager.GetSprite(ctx, resources.SpritePlayer)
	if !ok {
		return errors.NewGameError(errors.AssetNotFound, "player sprite not found")
	}

	starSprite, ok := g.resourceManager.GetSprite(ctx, resources.SpriteStar)
	if !ok {
		return errors.NewGameError(errors.AssetNotFound, "star sprite not found")
	}

	// Create player with config from JSON
	if g.playerConfig == nil {
		return errors.NewGameError(errors.ConfigMissing, "player config not loaded")
	}
	g.playerEntity = core.CreatePlayer(g.world, playerSprite, g.config, g.playerConfig)
	g.logger.Debug("Player entity created",
		"entity_id", g.playerEntity,
		"health", g.playerConfig.Health,
		"size", g.playerConfig.Size)

	// Create star field
	g.starEntities = core.CreateStarField(g.world, starSprite, g.config)
	g.logger.Debug("Star entities created", "count", len(g.starEntities))

	// Log star positions for debugging
	for i, entity := range g.starEntities {
		if i < 5 { // Only log first 5 stars
			entry := g.world.Entry(entity)
			if entry.Valid() {
				pos := core.Position.Get(entry)
				g.logger.Debug("Star position", "star_id", i, "pos", pos)
			}
		}
	}

	return nil
}
