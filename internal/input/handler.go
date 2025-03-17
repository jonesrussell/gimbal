package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/logger"
)

const (
	// MovementSpeedDegreesPerFrame is the speed at which the player moves in degrees per frame
	MovementSpeedDegreesPerFrame = 2
)

// Interface defines the input handler interface
type Interface interface {
	HandleInput()
	IsKeyPressed(key ebiten.Key) bool
	GetMovementInput() common.Angle
	IsQuitPressed() bool
	IsPausePressed() bool
	// Simulation methods for testing
	SimulateKeyPress(key ebiten.Key)
	SimulateKeyRelease(key ebiten.Key)
}

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
	for key := ebiten.Key(0); key <= ebiten.KeyMax; key++ {
		wasPressed := h.keyState[key]
		isPressed := ebiten.IsKeyPressed(key)
		h.keyState[key] = isPressed

		// Log key state changes for arrow keys
		if key == ebiten.KeyLeft || key == ebiten.KeyRight {
			logger.GlobalLogger.Debug("Arrow key state",
				"key", key,
				"was_pressed", wasPressed,
				"is_pressed", isPressed,
				"key_left", ebiten.KeyLeft,
				"key_right", ebiten.KeyRight,
			)
		}
	}
}

// IsKeyPressed implements InputHandler interface
func (h *Handler) IsKeyPressed(key ebiten.Key) bool {
	return ebiten.IsKeyPressed(key)
}

// GetMovementInput returns the movement direction based on key states
func (h *Handler) GetMovementInput() common.Angle {
	var angle common.Angle

	leftPressed := ebiten.IsKeyPressed(ebiten.KeyLeft)
	rightPressed := ebiten.IsKeyPressed(ebiten.KeyRight)

	logger.GlobalLogger.Debug("Movement check",
		"left_pressed", leftPressed,
		"right_pressed", rightPressed,
		"key_left", ebiten.KeyLeft,
		"key_right", ebiten.KeyRight,
	)

	switch {
	case leftPressed:
		angle = common.Angle(-MovementSpeedDegreesPerFrame) // Move counter-clockwise
		logger.GlobalLogger.Debug("Movement input", "direction", "left", "angle", angle)
	case rightPressed:
		angle = common.Angle(MovementSpeedDegreesPerFrame) // Move clockwise
		logger.GlobalLogger.Debug("Movement input", "direction", "right", "angle", angle)
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

// SimulateKeyPress is a no-op for the real input handler
func (h *Handler) SimulateKeyPress(key ebiten.Key) {}

// SimulateKeyRelease is a no-op for the real input handler
func (h *Handler) SimulateKeyRelease(key ebiten.Key) {}
