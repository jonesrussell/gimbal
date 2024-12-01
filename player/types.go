package player

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// InputHandlerInterface defines the contract for handling player input
type InputHandlerInterface interface {
	HandleInput() (float64, float64)
	IsKeyPressed(ebiten.Key) bool
}
