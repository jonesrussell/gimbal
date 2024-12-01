package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jonesrussell/gimbal/player"
)

// InputHandler implements player.InputHandlerInterface
type InputHandler struct{}

// IsKeyPressed checks if a key is currently pressed
func (rh *InputHandler) IsKeyPressed(key ebiten.Key) bool {
	return ebiten.IsKeyPressed(key)
}

// HandleInput processes player input and returns movement values
func (rh *InputHandler) HandleInput() (float64, float64) {
	// Implementation will go here
	return 0, 0
}

// NewInputHandler creates and returns a new input handler
func NewInputHandler() player.InputHandlerInterface {
	return &InputHandler{}
}
