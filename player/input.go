package player

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// InputHandler implements the InputHandlerInterface
type InputHandler struct{}

// IsKeyPressed checks if a key is pressed
func (i *InputHandler) IsKeyPressed(key ebiten.Key) bool {
	return ebiten.IsKeyPressed(key)
}
