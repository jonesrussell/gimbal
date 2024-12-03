package engine

import "github.com/hajimehoshi/ebiten/v2"

// GameState represents the current state of the game
type GameState int

const (
	StateTitle GameState = iota
	StatePlaying
	StatePaused
	StateGameOver
)

// Debug controls global debug state
var Debug bool

// GameEngine represents the core game interface
type GameEngine interface {
	Update() error
	Draw(screen *ebiten.Image)
}
