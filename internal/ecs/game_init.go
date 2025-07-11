package ecs

import (
	"github.com/yohamta/donburi"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/ecs/core"
	scenes "github.com/jonesrussell/gimbal/internal/ecs/scenes"
)

// NewECSGame creates a new ECS-based game instance
func NewECSGame(
	config *common.GameConfig,
	logger common.Logger,
	inputHandler common.GameInputHandler,
) (*ECSGame, error) {
	if config == nil {
		return nil, common.NewGameError(common.ErrorCodeConfigMissing, "config cannot be nil")
	}
	if logger == nil {
		return nil, common.NewGameError(common.ErrorCodeConfigMissing, "logger cannot be nil")
	}
	if inputHandler == nil {
		return nil, common.NewGameError(common.ErrorCodeConfigMissing, "inputHandler cannot be nil")
	}

	logger.Debug("Creating new ECS game instance",
		"screen_size", config.ScreenSize,
		"player_size", config.PlayerSize,
		"num_stars", config.NumStars,
	)

	// Create ECS world
	world := donburi.NewWorld()

	// Create game instance
	game := &ECSGame{
		world:        world,
		config:       config,
		inputHandler: inputHandler,
		logger:       logger,
	}

	// Initialize systems and managers
	if err := game.initializeSystems(); err != nil {
		return nil, err
	}

	// Load assets
	if err := game.loadAssets(); err != nil {
		return nil, common.NewGameErrorWithCause(common.ErrorCodeAssetLoadFailed, "failed to load assets", err)
	}

	// Create entities
	if err := game.createEntities(); err != nil {
		return nil, common.NewGameErrorWithCause(common.ErrorCodeEntityCreationFailed, "failed to create entities", err)
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
	g.eventSystem = NewEventSystem(g.world)
	g.logger.Debug("Event system created")

	// Create health system
	g.healthSystem = NewHealthSystem(g.world, g.config, g.eventSystem, g.stateManager, g.logger)
	g.logger.Debug("Health system created")

	// Create resource manager
	g.resourceManager = NewResourceManager(g.logger)
	g.logger.Debug("Resource manager created")

	// Create game state managers
	g.stateManager = NewGameStateManager(g.eventSystem, g.logger)
	g.scoreManager = NewScoreManager(g.eventSystem, g.logger)
	g.levelManager = NewLevelManager(g.logger)

	// Create scene manager
	g.sceneManager = scenes.NewSceneManager(g.world, g.config, g.logger, g.inputHandler)

	// Set resume callback to unpause game state
	g.sceneManager.SetResumeCallback(func() {
		g.stateManager.SetPaused(false)
	})

	// Set health system for scenes to access
	g.sceneManager.SetHealthSystem(g.healthSystem)

	// Set initial scene
	if err := g.sceneManager.SetInitialScene(scenes.SceneStudioIntro); err != nil {
		return common.NewGameErrorWithCause(common.ErrorCodeSystemFailed, "failed to set initial scene", err)
	}

	// Create combat systems
	g.enemySystem = NewEnemySystem(g.world, g.config)
	g.weaponSystem = NewWeaponSystem(g.world, g.config)
	g.collisionSystem = NewCollisionSystem(g.world, g.config, g.healthSystem, g.eventSystem, g.logger)

	return nil
}

// loadAssets loads and prepares game assets
func (g *ECSGame) loadAssets() error {
	// Load all sprites through resource manager
	if err := g.resourceManager.LoadAllSprites(); err != nil {
		return common.NewGameErrorWithCause(common.ErrorCodeAssetLoadFailed, "failed to load sprites", err)
	}

	g.logger.Debug("Assets loaded successfully", "resource_count", g.resourceManager.GetResourceCount())
	return nil
}

// createEntities creates all game entities
func (g *ECSGame) createEntities() error {
	// Get sprites from resource manager
	playerSprite, ok := g.resourceManager.GetSprite(SpritePlayer)
	if !ok {
		return common.NewGameError(common.ErrorCodeSpriteNotFound, "player sprite not found")
	}

	starSprite, ok := g.resourceManager.GetSprite(SpriteStar)
	if !ok {
		return common.NewGameError(common.ErrorCodeSpriteNotFound, "star sprite not found")
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
	// No need for system manager with wrappers
}
