package engine

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"go.uber.org/zap"

	"github.com/jonesrussell/gimbal/internal/config"
)

// Game represents the main game engine
type Game struct {
	logger    *zap.Logger
	config    *config.Config
	gameState GameEngine
	stars     []Star
	state     GameState
	player    *Player
}

// NewGame creates a new game instance with dependencies
func NewGame(logger *zap.Logger, config *config.Config, gameState GameEngine) (*Game, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger is required")
	}
	if config == nil {
		return nil, fmt.Errorf("config is required")
	}
	if gameState == nil {
		return nil, fmt.Errorf("gameState is required")
	}

	// Set global debug flag
	Debug = config.Game.Debug

	logger.Debug("Debug mode",
		zap.Bool("enabled", Debug))

	g := &Game{
		logger:    logger,
		config:    config,
		gameState: gameState,
		state:     StatePlaying,
	}

	// Initialize player
	playerImage := ebiten.NewImage(20, 20)
	playerImage.Fill(color.RGBA{R: 0, G: 255, B: 0, A: 255}) // Green player
	g.player = NewPlayer(playerImage, float64(g.config.Screen.Width/2), float64(g.config.Screen.Height/2))

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
	g.updateStars()
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
