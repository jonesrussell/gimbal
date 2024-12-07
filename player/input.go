package player

import (
	"github.com/hajimehoshi/ebiten/v2"
	"go.uber.org/zap"
)

// InputHandler implements the InputHandlerInterface
type InputHandler struct {
	logger *zap.Logger
}

// NewInputHandler creates a new input handler instance
func NewInputHandler(logger *zap.Logger) *InputHandler {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &InputHandler{
		logger: logger,
	}
}

// HandleInput processes keyboard input and returns the horizontal and vertical movement values
func (i *InputHandler) HandleInput() (horizontalMove, verticalMove float64) {
	// For now returning 0,0 as noted in the comment about circular movement
	return 0, 0
}

// IsKeyPressed checks if a specific key is currently pressed
func (i *InputHandler) IsKeyPressed(key ebiten.Key) bool {
	return ebiten.IsKeyPressed(key)
}
