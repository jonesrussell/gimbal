package player

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// InputHandlerInterface defines the contract for handling player input
type InputHandlerInterface interface {
	IsKeyPressed(ebiten.Key) bool
	HandleInput()
	// Add other required methods based on your implementation
}
