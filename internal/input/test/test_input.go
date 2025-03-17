package test

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/input"
)

// Handler implements input handling for tests
type Handler struct {
	pressedKeys map[ebiten.Key]bool
}

// NewHandler creates a new test input handler
func NewHandler() *Handler {
	return &Handler{
		pressedKeys: make(map[ebiten.Key]bool),
	}
}

// HandleInput implements input.Interface
func (h *Handler) HandleInput() {
	// No-op for testing
}

// IsKeyPressed returns whether a key is currently pressed
func (h *Handler) IsKeyPressed(key ebiten.Key) bool {
	return h.pressedKeys[key]
}

// GetMovementInput implements input.Interface
func (h *Handler) GetMovementInput() common.Angle {
	var angle common.Angle

	if h.pressedKeys[ebiten.KeyLeft] {
		angle = common.Angle(-input.MovementSpeedDegreesPerFrame)
	} else if h.pressedKeys[ebiten.KeyRight] {
		angle = common.Angle(input.MovementSpeedDegreesPerFrame)
	}

	return angle
}

// IsQuitPressed implements input.Interface
func (h *Handler) IsQuitPressed() bool {
	return false // Quit not implemented in tests
}

// IsPausePressed implements input.Interface
func (h *Handler) IsPausePressed() bool {
	return h.pressedKeys[ebiten.KeySpace]
}

// SimulateKeyPress simulates a key press
func (h *Handler) SimulateKeyPress(key ebiten.Key) {
	h.pressedKeys[key] = true
}

// SimulateKeyRelease simulates a key release
func (h *Handler) SimulateKeyRelease(key ebiten.Key) {
	h.pressedKeys[key] = false
}
