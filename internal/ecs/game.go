package ecs

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"

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

	// Entity references
	playerEntity donburi.Entity
	starEntities []donburi.Entity

	// Assets
	playerSprite *ebiten.Image
	starSprite   *ebiten.Image
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

	// Create game instance
	game := &ECSGame{
		world:        world,
		config:       config,
		inputHandler: inputHandler,
		logger:       logger,
		isPaused:     false,
	}

	// Load assets
	if err := game.loadAssets(); err != nil {
		return nil, fmt.Errorf("failed to load assets: %w", err)
	}

	// Create entities
	if err := game.createEntities(); err != nil {
		return nil, fmt.Errorf("failed to create entities: %w", err)
	}

	return game, nil
}

// loadAssets loads and prepares game assets
func (g *ECSGame) loadAssets() error {
	// For now, create a simple placeholder sprite
	// TODO: Load actual assets from a shared location
	g.playerSprite = ebiten.NewImage(32, 32)
	g.playerSprite.Fill(color.RGBA{0, 255, 0, 255}) // Green square
	g.logger.Debug("Player sprite created (placeholder)")

	// Create a simple star sprite (white square for now)
	g.starSprite = ebiten.NewImage(10, 10)
	g.starSprite.Fill(color.White)
	g.logger.Debug("Star sprite created")

	return nil
}

// createEntities creates all game entities
func (g *ECSGame) createEntities() error {
	// Create player
	g.playerEntity = CreatePlayer(g.world, g.playerSprite, g.config)
	g.logger.Debug("Player entity created", "entity_id", g.playerEntity)

	// Create star field
	g.starEntities = CreateStarField(g.world, g.starSprite, g.config)
	g.logger.Debug("Star entities created", "count", len(g.starEntities))

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
		g.logger.Debug("Game paused", "is_paused", g.isPaused)
		return nil
	}

	// Get input angle for player movement
	inputAngle := g.inputHandler.GetMovementInput()

	// Run ECS systems
	PlayerInputSystem(g.world, inputAngle)
	OrbitalMovementSystem(g.world)
	StarMovementSystem(g.world, g.config.ScreenSize.Height)

	return nil
}

// Draw renders the game
func (g *ECSGame) Draw(screen *ebiten.Image) {
	// Clear screen
	screen.Fill(color.Black)

	// Run render system
	RenderSystem(g.world, screen)

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

		// Draw debug text (simplified for now)
		// In a real implementation, you'd use a proper text rendering system
		_ = fmt.Sprintf("Player: Pos(%.1f, %.1f) Angle: %.1f°",
			pos.X, pos.Y, orb.OrbitalAngle)
	}
}

// Layout implements ebiten.Game interface
func (g *ECSGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.config.ScreenSize.Width, g.config.ScreenSize.Height
}

// Cleanup cleans up resources
func (g *ECSGame) Cleanup() {
	g.logger.Debug("Cleaning up ECS game")
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
