package engine

import (
	"github.com/hajimehoshi/ebiten/v2"
	"go.uber.org/zap"

	"github.com/jonesrussell/gimbal/internal/config"
)

// Game represents the main game engine
type Game struct {
	logger *zap.Logger
	config *config.Config
	stars  []Star
	state  GameState
}

// NewGame creates a new game instance with dependencies
func NewGame(logger *zap.Logger, config *config.Config) (*Game, error) {
	g := &Game{
		logger: logger,
		config: config,
		state:  StateTitle,
	}

	return g, nil
}

// Update handles game logic per frame
func (g *Game) Update() error {
	g.updateStars()
	return nil
}

// Draw handles rendering
func (g *Game) Draw(screen *ebiten.Image) {
	g.drawStars(screen)
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
