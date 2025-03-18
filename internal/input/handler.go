package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/logger"
	"go.uber.org/zap"
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

// Handler handles input for the game
type Handler struct {
	leftPressed  bool
	rightPressed bool
	pausePressed bool
	quitPressed  bool
}

// New creates a new input handler
func New() *Handler {
	return &Handler{}
}

// HandleInput handles input for the current frame
func (h *Handler) HandleInput() {
	// Check arrow keys
	h.leftPressed = ebiten.IsKeyPressed(ebiten.KeyArrowLeft)
	h.rightPressed = ebiten.IsKeyPressed(ebiten.KeyArrowRight)
	h.pausePressed = ebiten.IsKeyPressed(ebiten.KeySpace)
	h.quitPressed = ebiten.IsKeyPressed(ebiten.KeyEscape)

	logger.GlobalLogger.Debug("Arrow key state",
		zap.Bool("left", h.leftPressed),
		zap.Bool("right", h.rightPressed),
		zap.Bool("pause", h.pausePressed),
		zap.Bool("quit", h.quitPressed),
	)
}

// IsKeyPressed implements InputHandler interface
func (h *Handler) IsKeyPressed(key ebiten.Key) bool {
	return ebiten.IsKeyPressed(key)
}

// GetMovementInput returns the movement input angle
func (h *Handler) GetMovementInput() common.Angle {
	logger.GlobalLogger.Debug("Movement check",
		zap.Bool("left", h.leftPressed),
		zap.Bool("right", h.rightPressed),
	)

	if h.leftPressed {
		angle := common.Angle(-1)
		logger.GlobalLogger.Debug("Movement input",
			zap.String("direction", "left"),
			zap.Any("angle", angle),
		)
		return angle
	}

	if h.rightPressed {
		angle := common.Angle(1)
		logger.GlobalLogger.Debug("Movement input",
			zap.String("direction", "right"),
			zap.Any("angle", angle),
		)
		return angle
	}

	return 0
}

// IsQuitPressed returns whether the quit key is pressed
func (h *Handler) IsQuitPressed() bool {
	return h.quitPressed
}

// IsPausePressed returns whether the pause key is pressed
func (h *Handler) IsPausePressed() bool {
	return h.pausePressed
}

// SimulateKeyPress is a no-op for the real input handler
func (h *Handler) SimulateKeyPress(key ebiten.Key) {}

// SimulateKeyRelease is a no-op for the real input handler
func (h *Handler) SimulateKeyRelease(key ebiten.Key) {}
