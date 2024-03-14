package game

import "github.com/hajimehoshi/ebiten/v2"

// InputHandlerInterface defines the methods for handling input.
type InputHandlerInterface interface {
	IsKeyPressed(key ebiten.Key) bool
}

// InputHandler implements HandlerInterface for the real game.
type InputHandler struct{}

func (rh *InputHandler) IsKeyPressed(key ebiten.Key) bool {
	return ebiten.IsKeyPressed(key)
}
