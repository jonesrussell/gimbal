package ecs

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/input"
)

// We'll load assets from the game package for now
// TODO: Move assets to a shared location or copy them here

// ECSGame represents the main game state using ECS
type ECSGame struct {
	world        donburi.World
	config       *common.GameConfig
	inputHandler input.Interface
	logger       common.Logger
	isPaused     bool

	// Event system
	eventSystem *EventSystem

	// Resource management
	resourceManager *ResourceManager

	// System management
	systemManager *SystemManager

	// Entity references
	playerEntity donburi.Entity
	starEntities []donburi.Entity
}

// NewECSGame creates a new ECS-based game instance
func NewECSGame(config *common.GameConfig, logger common.Logger) (*ECSGame, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}
	if logger == nil {
		return nil, fmt.Errorf("logger cannot be nil")
	}

	logger.Debug("Creating new ECS game instance",
		"screen_size", config.ScreenSize,
		"player_size", config.PlayerSize,
		"num_stars", config.NumStars,
	)

	// Create ECS world
	world := donburi.NewWorld()

	// Create input handler
	inputHandler := input.New(logger)
	logger.Debug("Input handler created")

	// Create event system
	eventSystem := NewEventSystem(world)
	logger.Debug("Event system created")

	// Create resource manager
	resourceManager := NewResourceManager(logger)
	logger.Debug("Resource manager created")

	// Create system manager
	systemManager := NewSystemManager()
	logger.Debug("System manager created")

	// Create game instance
	game := &ECSGame{
		world:           world,
		config:          config,
		inputHandler:    inputHandler,
		logger:          logger,
		isPaused:        false,
		eventSystem:     eventSystem,
		resourceManager: resourceManager,
		systemManager:   systemManager,
	}

	// Load assets
	if err := game.loadAssets(); err != nil {
		return nil, fmt.Errorf("failed to load assets: %w", err)
	}

	// Create entities
	if err := game.createEntities(); err != nil {
		return nil, fmt.Errorf("failed to create entities: %w", err)
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
		return fmt.Errorf("failed to load sprites: %w", err)
	}

	g.logger.Debug("Assets loaded successfully", "resource_count", g.resourceManager.GetResourceCount())
	return nil
}

// createEntities creates all game entities
func (g *ECSGame) createEntities() error {
	// Get sprites from resource manager
	playerSprite, ok := g.resourceManager.GetSprite(SpritePlayer)
	if !ok {
		return fmt.Errorf("player sprite not found")
	}

	starSprite, ok := g.resourceManager.GetSprite(SpriteStar)
	if !ok {
		return fmt.Errorf("star sprite not found")
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
	if g.isPaused {
		return nil
	}

	// Handle input
	g.inputHandler.HandleInput()

	// Check for pause
	if g.inputHandler.IsPausePressed() {
		g.isPaused = !g.isPaused
		if g.isPaused {
			g.eventSystem.EmitGamePaused()
		} else {
			g.eventSystem.EmitGameResumed()
		}
		g.logger.Debug("Game paused", "is_paused", g.isPaused)
		return nil
	}

	// Get input angle for player movement
	inputAngle := g.inputHandler.GetMovementInput()

	// Run player input system (needs input angle)
	playerInputWrapper := NewPlayerInputSystemWrapper(inputAngle)
	if err := playerInputWrapper.Update(g.world); err != nil {
		g.logger.Error("Player input system failed", "error", err)
	}

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

	return nil
}

// Draw renders the game
func (g *ECSGame) Draw(screen *ebiten.Image) {
	// Clear screen
	screen.Fill(color.Black)

	// Run render system through wrapper
	renderWrapper := NewRenderSystemWrapper(screen)
	if err := renderWrapper.Update(g.world); err != nil {
		g.logger.Error("Render system failed", "error", err)
	}

	// Draw debug info if enabled
	if g.config.Debug {
		g.drawDebugInfo(screen)
	}
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
	return g.isPaused
}

// SetInputHandler sets the input handler (for testing)
func (g *ECSGame) SetInputHandler(handler input.Interface) {
	g.inputHandler = handler
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
