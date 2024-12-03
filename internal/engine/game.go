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
	gameState GameEngine // Interface for game state
	stars     []Star
	state     GameState
}

// NewGame creates a new game instance with dependencies
func NewGame(logger *zap.Logger, config *config.Config, gameState GameEngine) (*Game, error) {
	g := &Game{
		logger:    logger,
		config:    config,
		gameState: gameState,
		state:     StatePlaying, // Or StateTitle if you want to start with a title screen
	}

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
	return g.gameState.Update()
}

// Draw handles rendering
func (g *Game) Draw(screen *ebiten.Image) {
	g.drawStars(screen)
	g.gameState.Draw(screen)
}

// Layout returns the game's screen dimensions
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.config.Screen.Width, g.config.Screen.Height
}

// Run starts the game loop
func (g *Game) Run() error {
	ebiten.SetWindowSize(g.config.Screen.Width, g.config.Screen.Height)
	ebiten.SetWindowTitle(g.config.Screen.Title)

	return ebiten.RunGame(g)
}
