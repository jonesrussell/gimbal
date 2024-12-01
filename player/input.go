package player

import "github.com/hajimehoshi/ebiten/v2"

// InputHandler implements InputHandlerInterface
type InputHandler struct{}

// HandleInput processes player input and returns movement values
func (ih *InputHandler) HandleInput() (float64, float64) {
	// Implementation will go here
	return 0, 0
}

// IsKeyPressed checks if a key is currently pressed
func (ih *InputHandler) IsKeyPressed(key ebiten.Key) bool {
	return ebiten.IsKeyPressed(key)
}

// NewInputHandler creates and returns a new input handler
func NewInputHandler() InputHandlerInterface {
	return &InputHandler{}
}
