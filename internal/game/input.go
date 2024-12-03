package game

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// InputHandlerInterface defines the interface for handling input
type InputHandlerInterface interface {
	IsKeyPressed(key ebiten.Key) bool
}

// InputHandler implements the InputHandlerInterface
type InputHandler struct{}

// IsKeyPressed checks if a key is pressed
func (i *InputHandler) IsKeyPressed(key ebiten.Key) bool {
	return ebiten.IsKeyPressed(key)
}
