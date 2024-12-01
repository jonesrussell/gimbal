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

// GameEngine represents the core game interface
type GameEngine interface {
	Update() error
	Draw(screen *ebiten.Image)
}
