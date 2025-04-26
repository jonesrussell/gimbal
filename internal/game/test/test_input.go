package test

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/input"
)

const (
	// MovementSpeedDegreesPerFrame is the speed at which the player moves in degrees per frame
	MovementSpeedDegreesPerFrame = 2
)

// TestHandler is a test-specific input handler that allows simulating key presses
type TestHandler struct {
	simulatedKeys map[ebiten.Key]bool
	keyState      map[ebiten.Key]bool
}

// NewTestHandler creates a new test input handler
func NewTestHandler() input.Interface {
	handler := &TestHandler{
		simulatedKeys: make(map[ebiten.Key]bool),
		keyState:      make(map[ebiten.Key]bool),
	}
	return handler
}

// SimulateKeyPress simulates a key press for testing
func (h *TestHandler) SimulateKeyPress(key ebiten.Key) {
	h.simulatedKeys[key] = true
}

// SimulateKeyRelease simulates a key release for testing
func (h *TestHandler) SimulateKeyRelease(key ebiten.Key) {
	h.simulatedKeys[key] = false
}

// HandleInput overrides the base handler to use simulated keys
func (h *TestHandler) HandleInput() {
	// Update key states
	for key := ebiten.Key(0); key <= ebiten.KeyMax; key++ {
		h.keyState[key] = h.IsKeyPressed(key)
	}
}

// IsKeyPressed overrides the base handler to use simulated keys
func (h *TestHandler) IsKeyPressed(key ebiten.Key) bool {
	if isPressed, ok := h.simulatedKeys[key]; ok {
		return isPressed
	}
	return ebiten.IsKeyPressed(key)
}

// GetMovementInput overrides the base handler to use simulated keys
func (h *TestHandler) GetMovementInput() common.Angle {
	var angle common.Angle

	leftPressed := h.IsKeyPressed(ebiten.KeyLeft)
	rightPressed := h.IsKeyPressed(ebiten.KeyRight)

	switch {
	case leftPressed:
		angle = common.Angle(-MovementSpeedDegreesPerFrame)
	case rightPressed:
		angle = common.Angle(MovementSpeedDegreesPerFrame)
	}

	return angle
}

// IsQuitPressed overrides the base handler to use simulated keys
func (h *TestHandler) IsQuitPressed() bool {
	return h.IsKeyPressed(ebiten.KeyEscape)
}

// IsPausePressed overrides the base handler to use simulated keys
func (h *TestHandler) IsPausePressed() bool {
	return h.IsKeyPressed(ebiten.KeySpace)
}

// GetLastEvent returns the last input event that occurred
func (h *TestHandler) GetLastEvent() input.InputEvent {
	// For testing purposes, we'll return None since we're only testing keyboard input
	return input.InputEventNone
}

// GetTouchState returns the current touch state
func (h *TestHandler) GetTouchState() *input.TouchState {
	return nil
}

// GetMousePosition returns the current mouse position
func (h *TestHandler) GetMousePosition() common.Point {
	return common.Point{}
}

// IsMouseButtonPressed checks if a mouse button is pressed
func (h *TestHandler) IsMouseButtonPressed(button ebiten.MouseButton) bool {
	return false
}
