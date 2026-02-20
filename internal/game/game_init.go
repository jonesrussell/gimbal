package game

import (
	"context"
	"log"

	"github.com/yohamta/donburi"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/config"
	"github.com/jonesrussell/gimbal/internal/errors"
	"github.com/jonesrussell/gimbal/internal/scenes"
)

// validateGameConfig validates the provided game configuration
func (g *ECSGame) validateGameConfig(cfg *config.GameConfig) error {
	if cfg == nil {
		return errors.NewGameError(errors.ConfigMissing, "config cannot be nil")
	}
	return nil
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
	inputHandler common.GameInputHandler,
) (*ECSGame, error) {
	// Create ECS world
	world := donburi.NewWorld()

	// Create game context for lifecycle management
	gameCtx, cancel := context.WithCancel(ctx)

	game := &ECSGame{
		world:        world,
		config:       gameConfig,
		inputHandler: inputHandler,
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

	// Load all audio through resource manager
	if err := g.resourceManager.LoadAllAudio(ctx); err != nil {
		// Audio is optional, log warning but don't fail
		log.Printf("[WARN] Failed to load audio, continuing without it: %v", err)
	}

	return nil
}
