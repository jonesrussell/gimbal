package engine

import (
	"context"
	"fmt"
	"image"
	"image/color"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"go.uber.org/zap"

	"github.com/jonesrussell/gimbal/internal/config"
	"github.com/jonesrussell/gimbal/internal/core"
	"github.com/jonesrussell/gimbal/player"
)

// Game represents the main game engine
type Game struct {
	// 8-byte aligned fields (pointers and slices)
	player    *player.Player     // 8 bytes
	logger    *zap.Logger        // 8 bytes
	config    *config.Config     // 8 bytes
	assets    core.AssetManager  // 8 bytes (interface)
	gameState GameEngine         // 8 bytes (interface)
	cancel    context.CancelFunc // 8 bytes (function pointer)
	stars     []Star             // 24 bytes (slice header)

	// 4-byte aligned fields
	state GameState // 4 bytes
}

// NewGame creates a new game instance with dependencies
func NewGame(logger *zap.Logger, cfg *config.Config, gameState GameEngine, assets core.AssetManager, cancel context.CancelFunc) (*Game, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger is required")
	}
	if cfg == nil {
		return nil, fmt.Errorf("config is required")
	}
	if gameState == nil {
		return nil, fmt.Errorf("gameState is required")
	}
	if assets == nil {
		return nil, fmt.Errorf("assets is required")
	}
	if cancel == nil {
		return nil, fmt.Errorf("cancel function is required")
	}

	// Set global debug flag
	Debug = cfg.Game.Debug

	logger.Debug("Debug mode",
		zap.Bool("enabled", Debug))

	g := &Game{
		logger:    logger,
		config:    cfg,
		gameState: gameState,
		assets:    assets,
		state:     StatePlaying,
		cancel:    cancel,
	}

	// Initialize player with proper circular movement
	center := image.Point{X: g.config.Screen.Width / 2, Y: g.config.Screen.Height / 2}
	inputHandler := player.NewInputHandler(g.logger)
	playerSprite, err := g.assets.LoadImage(context.Background(), "images/player.png")
	if err != nil {
		return nil, fmt.Errorf("failed to load player sprite: %w", err)
	}

	p, err := player.NewPlayer(inputHandler, g.config.Game.Speed, playerSprite, center)
	if err != nil {
		return nil, fmt.Errorf("failed to create player: %w", err)
	}
	g.player = p

	// Initialize stars
	starImage := ebiten.NewImage(1, 1)
	starImage.Fill(color.White)
	stars, err := initializeStars(g.config.Game.NumStars, starImage)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize stars: %w", err)
	}
	g.stars = stars

	logger.Debug("Initialized stars",
		zap.Int("numStars", g.config.Game.NumStars))

	return g, nil
}

// Update handles game logic per frame
func (g *Game) Update() error {
	if ebiten.IsWindowBeingClosed() {
		g.logger.Info("window close requested")
		g.cancel()
		return fmt.Errorf("window closed")
	}

	if err := g.updateStars(); err != nil {
		return fmt.Errorf("failed to update stars: %w", err)
	}
	g.player.Update()
	return g.gameState.Update()
}

// Draw handles rendering
func (g *Game) Draw(screen *ebiten.Image) {
	g.drawStars(screen)
	g.player.Draw(screen)
	g.gameState.Draw(screen)
}

// Layout returns the game's screen dimensions
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

// Run starts the game loop
func (g *Game) Run() error {
	ebiten.SetWindowSize(g.config.Screen.Width, g.config.Screen.Height)
	ebiten.SetWindowTitle(g.config.Screen.Title)

	return ebiten.RunGame(g)
}
