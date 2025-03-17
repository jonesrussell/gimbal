package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jonesrussell/gimbal/internal/common"
)

// Handler handles game input
type Handler struct {
	keyState map[ebiten.Key]bool
}

// New creates a new input handler
func New() *Handler {
	return &Handler{
		keyState: make(map[ebiten.Key]bool),
	}
}

// HandleInput implements InputHandler interface
func (h *Handler) HandleInput() {
	// Update key states
	for key := ebiten.Key0; key <= ebiten.KeyMax; key++ {
		h.keyState[key] = ebiten.IsKeyPressed(key)
	}
}

// IsKeyPressed implements InputHandler interface
func (h *Handler) IsKeyPressed(key ebiten.Key) bool {
	return h.keyState[key]
}

// GetMovementInput returns the movement direction based on key states
func (h *Handler) GetMovementInput() common.Angle {
	var angle common.Angle

	switch {
	case h.IsKeyPressed(ebiten.KeyLeft):
		angle -= 1
	case h.IsKeyPressed(ebiten.KeyRight):
		angle += 1
	}

	return angle
}

// IsQuitPressed returns true if the quit key is pressed
func (h *Handler) IsQuitPressed() bool {
	return h.IsKeyPressed(ebiten.KeyEscape)
}

// IsPausePressed returns true if the pause key is pressed
func (h *Handler) IsPausePressed() bool {
	return h.IsKeyPressed(ebiten.KeySpace)
}
