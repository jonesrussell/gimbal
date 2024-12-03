package player

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// InputHandler implements the InputHandlerInterface
type InputHandler struct{}

func NewInputHandler() InputHandlerInterface {
	return &InputHandler{}
}

// HandleInput implements InputHandlerInterface
func (i *InputHandler) HandleInput() (float64, float64) {
	return 0, 0 // This method might not be needed for circular movement
}

// IsKeyPressed implements InputHandlerInterface
func (i *InputHandler) IsKeyPressed(key ebiten.Key) bool {
	return ebiten.IsKeyPressed(key)
}
