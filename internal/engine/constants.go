package engine

import "github.com/hajimehoshi/ebiten/v2"

// GameState represents the current state of the game
type GameState int

const (
	StatePlaying GameState = iota
)

// Debug controls global debug state
var Debug bool

// GameEngine represents the core game interface
type GameEngine interface {
	Update() error
	Draw(screen *ebiten.Image)
}
