package game

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/jonesrussell/gimbal/internal/config"
)

// Runner handles game execution with Ebiten
type Runner struct {
	config *config.AppGameConfig
	game   ebiten.Game
}

// NewRunner creates a new game runner
func NewRunner(cfg *config.AppGameConfig, game ebiten.Game) *Runner {
	return &Runner{
		config: cfg,
		game:   game,
	}
}

// Run configures Ebiten and starts the game
func (r *Runner) Run() error {
	// Configure Ebiten window
	ebiten.SetWindowSize(r.config.WindowWidth, r.config.WindowHeight)
	ebiten.SetWindowTitle(r.config.WindowTitle)
	ebiten.SetTPS(r.config.TPS)

	if r.config.Resizable {
		ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	}

	// Run the game
	return ebiten.RunGame(r.game)
}
