package ecs

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"

	"github.com/jonesrussell/gimbal/internal/common"
)

// We'll load assets from the game package for now
// TODO: Move assets to a shared location or copy them here

// ECSGame represents the main game state using ECS
type ECSGame struct {
	world        donburi.World
	config       *common.GameConfig
	inputHandler common.GameInputHandler
	logger       common.Logger

	// Event system
	eventSystem *EventSystem

	// Resource management
	resourceManager *ResourceManager

	// System management
	systemManager *SystemManager

	// Game state management
	stateManager *GameStateManager

	// Scene management
	sceneManager *SceneManager

	// Combat systems
	enemySystem     *EnemySystem
	weaponSystem    *WeaponSystem
	collisionSystem *CollisionSystem

	// Entity references
	playerEntity donburi.Entity
	starEntities []donburi.Entity
}

// NewECSGame creates a new ECS-based game instance
func NewECSGame(config *common.GameConfig, logger common.Logger, inputHandler common.GameInputHandler) (*ECSGame, error) {
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

	// Create event system
	eventSystem := NewEventSystem(world)
	logger.Debug("Event system created")

	// Create resource manager
	resourceManager := NewResourceManager(logger)
	logger.Debug("Resource manager created")

	// Create system manager
	systemManager := NewSystemManager()
	logger.Debug("System manager created")

	// Create game state manager
	stateManager := NewGameStateManager(eventSystem, logger)

	// Create scene manager
	sceneManager := NewSceneManager(world, config, logger, inputHandler)

	// Create combat systems
	enemySystem := NewEnemySystem(world, config)
	weaponSystem := NewWeaponSystem(world, config)
	collisionSystem := NewCollisionSystem(world, config)

	// Create game instance
	game := &ECSGame{
		world:           world,
		config:          config,
		inputHandler:    inputHandler,
		logger:          logger,
		eventSystem:     eventSystem,
		resourceManager: resourceManager,
		systemManager:   systemManager,
		stateManager:    stateManager,
		sceneManager:    sceneManager,
		enemySystem:     enemySystem,
		weaponSystem:    weaponSystem,
		collisionSystem: collisionSystem,
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
	g.playerEntity = CreatePlayer(g.world, playerSprite, g.config)
	g.logger.Debug("Player entity created", "entity_id", g.playerEntity)

	// Create star field
	g.starEntities = CreateStarField(g.world, starSprite, g.config)
	g.logger.Debug("Star entities created", "count", len(g.starEntities))

	// Log star positions for debugging
	for i, entity := range g.starEntities {
		if i < 5 { // Only log first 5 stars
			entry := g.world.Entry(entity)
			if entry.Valid() {
				pos := Position.Get(entry)
				g.logger.Debug("Star position", "star_id", i, "pos", pos)
			}
		}
	}

	return nil
}

// Update updates the game state
func (g *ECSGame) Update() error {
	// Handle input
	g.inputHandler.HandleInput()

	// Update scene manager
	if err := g.sceneManager.Update(); err != nil {
		g.logger.Error("Scene update failed", "error", err)
		return err
	}

	// Check current scene
	currentScene := g.sceneManager.GetCurrentScene()

	switch currentScene.GetType() {
	case SceneStudioIntro, SceneTitleScreen, SceneMenu:
		// Handle input for these scenes
		g.handleMenuInput()
	case ScenePlaying:
		// Handle gameplay
		if g.stateManager.IsPaused() {
			return nil
		}

		// Check for pause
		if g.inputHandler.IsPausePressed() {
			g.stateManager.TogglePause()
			g.sceneManager.SwitchScene(ScenePaused)
			return nil
		}

		// Get input angle for player movement
		inputAngle := g.inputHandler.GetMovementInput()

		// Run player input system (needs input angle)
		playerInputWrapper := NewPlayerInputSystemWrapper(inputAngle)
		if err := playerInputWrapper.Update(g.world); err != nil {
			g.logger.Error("Player input system failed", "error", err)
		}

		// Handle weapon firing
		g.handleWeaponFiring()

		// Update combat systems
		g.enemySystem.Update(1.0) // Assuming 60fps, so deltaTime = 1.0
		g.weaponSystem.Update(1.0)
		g.collisionSystem.Update()

		// Run other ECS systems through system manager
		if err := g.systemManager.UpdateAll(g.world); err != nil {
			g.logger.Error("System update failed", "error", err)
			return err
		}

		// Emit player movement event if player moved
		if inputAngle != 0 {
			playerEntry := g.world.Entry(g.playerEntity)
			if playerEntry.Valid() {
				pos := Position.Get(playerEntry)
				orb := Orbital.Get(playerEntry)
				g.eventSystem.EmitPlayerMoved(*pos, orb.OrbitalAngle)
			}
		}

		// Process all events
		g.eventSystem.ProcessEvents()

	case ScenePaused:
		// Handle pause menu input
		g.handlePauseInput()
	}

	return nil
}

// Draw renders the game
func (g *ECSGame) Draw(screen *ebiten.Image) {
	// Use scene manager to draw the current scene
	g.sceneManager.Draw(screen)
}

// drawDebugInfo renders debug information
func (g *ECSGame) drawDebugInfo(screen *ebiten.Image) {
	// Get player info for debug display
	playerEntry := g.world.Entry(g.playerEntity)
	if playerEntry.Valid() {
		pos := Position.Get(playerEntry)
		orb := Orbital.Get(playerEntry)

		// Log debug info
		g.logger.Debug("Debug Info",
			"player_pos", fmt.Sprintf("(%.1f, %.1f)", pos.X, pos.Y),
			"player_angle", fmt.Sprintf("%.1f°", orb.OrbitalAngle),
			"resource_count", g.resourceManager.GetResourceCount(),
			"entity_count", g.world.Len(),
		)
	}
}

// Layout implements ebiten.Game interface
func (g *ECSGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.config.ScreenSize.Width, g.config.ScreenSize.Height
}

// Cleanup cleans up resources
func (g *ECSGame) Cleanup() {
	g.logger.Debug("Cleaning up ECS game")

	// Clean up resources
	if g.resourceManager != nil {
		g.resourceManager.Cleanup()
	}

	// Donburi handles entity cleanup automatically
}

// IsPaused returns the pause state
func (g *ECSGame) IsPaused() bool {
	return g.stateManager.IsPaused()
}

// SetInputHandler sets the input handler (for testing)
func (g *ECSGame) SetInputHandler(handler common.GameInputHandler) {
	g.inputHandler = handler
}

// GetInputHandler returns the current input handler
func (g *ECSGame) GetInputHandler() common.GameInputHandler {
	return g.inputHandler
}

// setupEventSubscriptions sets up event handlers
func (g *ECSGame) setupEventSubscriptions() {
	// Subscribe to player movement events
	g.eventSystem.SubscribeToPlayerMoved(func(w donburi.World, event PlayerMovedEvent) {
		g.logger.Debug("Player moved",
			"position", event.Position,
			"angle", event.Angle)
	})

	// Subscribe to game state events
	g.eventSystem.SubscribeToGameState(func(w donburi.World, event GameStateEvent) {
		g.logger.Debug("Game state changed", "is_paused", event.IsPaused)
	})

	// Subscribe to score changes
	g.eventSystem.SubscribeToScoreChanged(func(w donburi.World, event ScoreChangedEvent) {
		g.logger.Debug("Score changed",
			"old_score", event.OldScore,
			"new_score", event.NewScore,
			"delta", event.Delta)
	})
}

// setupSystems sets up the system manager with all required systems
func (g *ECSGame) setupSystems() {
	// Add update systems in execution order
	g.systemManager.AddSystem(&MovementSystemWrapper{})
	g.systemManager.AddSystem(&OrbitalMovementSystemWrapper{})
	g.systemManager.AddSystem(NewStarMovementSystemWrapper(&ecs.ECS{World: g.world}, g.config))

	g.logger.Debug("Systems set up", "system_count", g.systemManager.GetSystemCount(), "systems", g.systemManager.GetSystemNames())
}

// handleWeaponFiring handles weapon firing based on input
func (g *ECSGame) handleWeaponFiring() {
	// Get player position and angle
	playerEntry := g.world.Entry(g.playerEntity)
	if !playerEntry.Valid() {
		return
	}

	pos := Position.Get(playerEntry)
	orb := Orbital.Get(playerEntry)

	// Check for fire input (Space key)
	if g.inputHandler.IsKeyPressed(ebiten.KeySpace) {
		g.weaponSystem.FireWeapon(WeaponTypePrimary, *pos, orb.FacingAngle)
	}

	// Check for secondary weapon (Shift key)
	if g.inputHandler.IsKeyPressed(ebiten.KeyShift) {
		g.weaponSystem.FireWeapon(WeaponTypeSecondary, *pos, orb.FacingAngle)
	}
}

// handleMenuInput handles input for menu scenes
func (g *ECSGame) handleMenuInput() {
	currentScene := g.sceneManager.GetCurrentScene()

	switch currentScene.GetType() {
	case SceneTitleScreen:
		// Any key to continue to main menu
		if g.inputHandler.GetLastEvent() != common.InputEventNone {
			g.sceneManager.SwitchScene(SceneMenu)
		}
	case SceneMenu:
		// Handle menu navigation
		menuScene := currentScene.(*MenuScene)

		// Navigation
		if g.inputHandler.IsKeyPressed(ebiten.KeyUp) {
			menuScene.selection = (menuScene.selection - 1 + len(menuScene.options)) % len(menuScene.options)
		}
		if g.inputHandler.IsKeyPressed(ebiten.KeyDown) {
			menuScene.selection = (menuScene.selection + 1) % len(menuScene.options)
		}

		// Selection
		if g.inputHandler.IsKeyPressed(ebiten.KeyEnter) || g.inputHandler.IsKeyPressed(ebiten.KeySpace) {
			switch menuScene.selection {
			case 0: // Start Game
				g.sceneManager.SwitchScene(ScenePlaying)
			case 1: // Options
				// TODO: Implement options menu
			case 2: // Credits
				// TODO: Implement credits screen
			case 3: // Quit
				// TODO: Implement quit functionality
			}
		}
	}
}

// handlePauseInput handles input for pause menu
func (g *ECSGame) handlePauseInput() {
	// Check for resume (ESC or any key)
	if g.inputHandler.IsQuitPressed() || g.inputHandler.GetLastEvent() != common.InputEventNone {
		g.sceneManager.SwitchScene(ScenePlaying)
	}
}
