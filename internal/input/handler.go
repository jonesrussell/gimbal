package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jonesrussell/gimbal/internal/common"
)

const (
	// MovementSpeedDegreesPerFrame is the speed at which the player moves in degrees per frame
	MovementSpeedDegreesPerFrame = 5
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

// Handler handles input for the game
type Handler struct {
	logger common.Logger
}

// New creates a new input handler
func New(logger common.Logger) *Handler {
	return &Handler{
		logger: logger,
	}
}

// HandleInput processes input events
func (h *Handler) HandleInput() {
	// Log input state for debugging
	if h.logger != nil {
		h.logger.Debug("Input state",
			"left", ebiten.IsKeyPressed(ebiten.KeyLeft),
			"right", ebiten.IsKeyPressed(ebiten.KeyRight),
			"space", ebiten.IsKeyPressed(ebiten.KeySpace),
			"escape", ebiten.IsKeyPressed(ebiten.KeyEscape),
		)
	}
}

// IsKeyPressed checks if a key is pressed
func (h *Handler) IsKeyPressed(key ebiten.Key) bool {
	return ebiten.IsKeyPressed(key)
}

// IsPausePressed checks if the pause key is pressed
func (h *Handler) IsPausePressed() bool {
	return ebiten.IsKeyPressed(ebiten.KeySpace)
}

// IsQuitPressed checks if the quit key is pressed
func (h *Handler) IsQuitPressed() bool {
	return ebiten.IsKeyPressed(ebiten.KeyEscape)
}

// GetMovementInput returns the movement angle based on input
func (h *Handler) GetMovementInput() common.Angle {
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		return -common.Angle(MovementSpeedDegreesPerFrame)
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		return common.Angle(MovementSpeedDegreesPerFrame)
	}
	return 0
}

// SimulateKeyPress is a no-op for the real input handler
func (h *Handler) SimulateKeyPress(key ebiten.Key) {}

// SimulateKeyRelease is a no-op for the real input handler
func (h *Handler) SimulateKeyRelease(key ebiten.Key) {}
