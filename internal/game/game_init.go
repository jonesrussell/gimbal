package game

import (
	"context"

	"github.com/yohamta/donburi"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/config"
	"github.com/jonesrussell/gimbal/internal/ecs/core"
	"github.com/jonesrussell/gimbal/internal/ecs/events"
	"github.com/jonesrussell/gimbal/internal/ecs/managers"
	resources "github.com/jonesrussell/gimbal/internal/ecs/managers/resource"
	"github.com/jonesrussell/gimbal/internal/ecs/systems/collision"
	"github.com/jonesrussell/gimbal/internal/ecs/systems/enemy"
	"github.com/jonesrussell/gimbal/internal/ecs/systems/health"
	"github.com/jonesrussell/gimbal/internal/ecs/systems/weapon"
	"github.com/jonesrussell/gimbal/internal/errors"
	"github.com/jonesrussell/gimbal/internal/scenes"
	"github.com/jonesrussell/gimbal/internal/ui"
)

// NewECSGame creates a new ECS-based game instance with dependency-injected UI
func NewECSGame(
	gameConfig *config.GameConfig,
	logger common.Logger,
	inputHandler common.GameInputHandler,
	uiFactory ui.UIFactory,
) (*ECSGame, error) {
	if gameConfig == nil {
		return nil, errors.NewGameError(errors.ErrorCodeConfigMissing, "config cannot be nil")
	}
	if logger == nil {
		return nil, errors.NewGameError(errors.ErrorCodeConfigMissing, "logger cannot be nil")
	}
	if inputHandler == nil {
		return nil, errors.NewGameError(errors.ErrorCodeConfigMissing, "inputHandler cannot be nil")
	}

	logger.Debug("Creating new ECS game instance",
		"screen_size", gameConfig.ScreenSize,
		"player_size", gameConfig.PlayerSize,
		"num_stars", gameConfig.NumStars,
	)

	// Create ECS world
	world := donburi.NewWorld()

	// Create game instance
	game := &ECSGame{
		world:        world,
		config:       gameConfig,
		inputHandler: inputHandler,
		logger:       logger,
	}

	// Initialize systems and managers
	if err := game.initializeSystems(); err != nil {
		return nil, err
	}

	// Load assets
	if err := game.loadAssets(); err != nil {
		return nil, errors.NewGameErrorWithCause(errors.ErrorCodeAssetLoadFailed, "failed to load assets", err)
	}

	// Create UI through factory (dependency injection)
	font, err := game.resourceManager.GetDefaultFont(context.Background())
	if err != nil {
		return nil, errors.NewGameErrorWithCause(errors.ErrorCodeAssetLoadFailed, "failed to get default font", err)
	}
	heartSprite, _ := game.resourceManager.GetUISprite(context.Background(), "heart", ui.HeartIconSize)
	uiConfig := ui.UIConfig{
		Font:  font,
		Theme: heartSprite,
	}
	game.ui = uiFactory.CreateGameUI(uiConfig)

	// Create entities
	if err := game.createEntities(); err != nil {
		return nil, errors.NewGameErrorWithCause(errors.ErrorCodeEntityCreationFailed, "failed to create entities", err)
	}

	// Set up event subscriptions
	game.setupEventSubscriptions()

	// Set up systems
	game.setupSystems()

	return game, nil
}

// initializeSystems creates all the systems and managers
func (g *ECSGame) initializeSystems() error {
	// Create event system
	g.eventSystem = events.NewEventSystem(g.world)
	g.logger.Debug("Event system created")

	// Create resource manager
	g.resourceManager = resources.NewResourceManager(g.logger)
	g.logger.Debug("Resource manager created")

	// Create game state managers
	g.stateManager = NewGameStateManager(g.eventSystem, g.logger)
	g.scoreManager = managers.NewScoreManager(10000) // Bonus life every 10,000 points
	g.levelManager = managers.NewLevelManager(g.logger)

	// Create health system (after state manager)
	g.healthSystem = health.NewHealthSystem(g.world, g.config, g.eventSystem, g.stateManager, g.logger)
	g.logger.Debug("Health system created")

	// Get font from resource manager
	font, err := g.resourceManager.GetDefaultFont(context.Background())
	if err != nil {
		return errors.NewGameErrorWithCause(errors.ErrorCodeAssetLoadFailed, "failed to get default font", err)
	}

	// Create scene manager
	g.sceneManager = scenes.NewSceneManager(&scenes.SceneManagerConfig{
		World:        g.world,
		Config:       g.config,
		Logger:       g.logger,
		InputHandler: g.inputHandler,
		Font:         font,
		ScoreManager: g.scoreManager,
		ResourceMgr:  g.resourceManager,
	})

	// Set resume callback to unpause game state
	g.sceneManager.SetResumeCallback(func() {
		g.stateManager.SetPaused(false)
	})

	// Set health system for scenes to access
	g.sceneManager.SetHealthSystem(g.healthSystem)

	// Set initial scene
	if err := g.sceneManager.SetInitialScene(scenes.SceneStudioIntro); err != nil {
		return errors.NewGameErrorWithCause(errors.ErrorCodeSystemFailed, "failed to set initial scene", err)
	}

	// Create combat systems
	g.enemySystem = enemy.NewEnemySystem(g.world, g.config, g.resourceManager, g.logger)
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

// loadAssets loads and prepares game assets
func (g *ECSGame) loadAssets() error {
	// Load all sprites through resource manager
	if err := g.resourceManager.LoadAllSprites(context.Background()); err != nil {
		return errors.NewGameErrorWithCause(errors.ErrorCodeAssetLoadFailed, "failed to load sprites", err)
	}

	g.logger.Debug("Assets loaded successfully", "resource_count", g.resourceManager.GetResourceCount())
	return nil
}

// createEntities creates all game entities
func (g *ECSGame) createEntities() error {
	// Get sprites from resource manager
	playerSprite, ok := g.resourceManager.GetSprite(context.Background(), resources.SpritePlayer)
	if !ok {
		return errors.NewGameError(errors.ErrorCodeSpriteNotFound, "player sprite not found")
	}

	starSprite, ok := g.resourceManager.GetSprite(context.Background(), resources.SpriteStar)
	if !ok {
		return errors.NewGameError(errors.ErrorCodeSpriteNotFound, "star sprite not found")
	}

	// Create player
	g.playerEntity = core.CreatePlayer(g.world, playerSprite, g.config)
	g.logger.Debug("Player entity created", "entity_id", g.playerEntity)

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

// setupSystems is no longer needed - systems are called directly in the Update loop
// This method is kept for compatibility but does nothing
func (g *ECSGame) setupSystems() {
	// Systems are now called directly in the Update loop
}
