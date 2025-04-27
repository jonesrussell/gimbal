package input

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jonesrussell/gimbal/internal/common"
)

const (
	// MovementSpeedDegreesPerFrame is the speed at which the player moves in degrees per frame
	MovementSpeedDegreesPerFrame = 5

	// Touch input constants
	MinTouchDuration = 10 // frames
	TouchThreshold   = 10 // pixels
)

// InputEvent represents a game input event
type InputEvent int

const (
	InputEventNone InputEvent = iota
	InputEventMove
	InputEventPause
	InputEventQuit
	InputEventTouch
	InputEventMouseMove
)

// TouchState tracks touch input state
type TouchState struct {
	ID       ebiten.TouchID
	StartPos common.Point
	LastPos  common.Point
	Duration int
}

// Interface defines the input handler interface
type Interface interface {
	HandleInput()
	IsKeyPressed(key ebiten.Key) bool
	GetMovementInput() common.Angle
	IsQuitPressed() bool
	IsPausePressed() bool
	// New input methods
	GetTouchState() *TouchState
	GetMousePosition() common.Point
	IsMouseButtonPressed(button ebiten.MouseButton) bool
	GetLastEvent() InputEvent
	// Simulation methods for testing
	SimulateKeyPress(key ebiten.Key)
	SimulateKeyRelease(key ebiten.Key)
}

// Handler handles input for the game
type Handler struct {
	logger        common.Logger
	touchState    *TouchState
	lastEvent     InputEvent
	mousePos      common.Point
	simulatedKeys map[ebiten.Key]bool // Track simulated key presses
}

// New creates a new input handler
func New(logger common.Logger) *Handler {
	return &Handler{
		logger:        logger,
		touchState:    nil,
		lastEvent:     InputEventNone,
		mousePos:      common.Point{},
		simulatedKeys: make(map[ebiten.Key]bool),
	}
}

// HandleInput processes input events
func (h *Handler) HandleInput() {
	h.lastEvent = InputEventNone

	// Handle keyboard input
	h.handleKeyboardInput()

	// Handle touch input
	h.handleTouchInput()

	// Handle mouse input
	h.handleMouseInput()

	// Log input state for debugging
	if h.logger != nil {
		h.logger.Debug("Input state",
			"left", h.IsKeyPressed(ebiten.KeyLeft),
			"right", h.IsKeyPressed(ebiten.KeyRight),
			"space", h.IsKeyPressed(ebiten.KeySpace),
			"escape", h.IsKeyPressed(ebiten.KeyEscape),
			"touch", h.touchState != nil,
			"mouse_pos", h.mousePos,
			"last_event", h.lastEvent,
		)
	}
}

func (h *Handler) handleKeyboardInput() {
	if h.IsKeyPressed(ebiten.KeyLeft) || h.IsKeyPressed(ebiten.KeyRight) {
		h.lastEvent = InputEventMove
	} else if h.IsKeyPressed(ebiten.KeySpace) {
		h.lastEvent = InputEventPause
	} else if h.IsKeyPressed(ebiten.KeyEscape) {
		h.lastEvent = InputEventQuit
	}
}

func (h *Handler) handleTouchInput() {
	// Handle touch input
	touchIDs := ebiten.TouchIDs()
	if len(touchIDs) > 0 {
		if h.touchState == nil {
			// New touch
			x, y := ebiten.TouchPosition(touchIDs[0])
			h.touchState = &TouchState{
				ID: touchIDs[0],
				StartPos: common.Point{
					X: float64(x),
					Y: float64(y),
				},
				LastPos: common.Point{
					X: float64(x),
					Y: float64(y),
				},
				Duration: 0,
			}
		} else {
			// Update existing touch
			x, y := ebiten.TouchPosition(h.touchState.ID)
			h.touchState.LastPos = common.Point{
				X: float64(x),
				Y: float64(y),
			}
			h.touchState.Duration++
			h.lastEvent = InputEventTouch
		}
	} else {
		h.touchState = nil
	}
}

func (h *Handler) handleMouseInput() {
	x, y := ebiten.CursorPosition()
	newPos := common.Point{X: float64(x), Y: float64(y)}
	if newPos != h.mousePos {
		h.mousePos = newPos
		h.lastEvent = InputEventMouseMove
	}
}

// IsKeyPressed checks if a key is pressed
func (h *Handler) IsKeyPressed(key ebiten.Key) bool {
	// Check simulated keys first
	if pressed, ok := h.simulatedKeys[key]; ok {
		return pressed
	}
	return ebiten.IsKeyPressed(key)
}

// GetMovementInput returns the movement angle based on input
func (h *Handler) GetMovementInput() common.Angle {
	if h.IsKeyPressed(ebiten.KeyLeft) {
		return -common.Angle(MovementSpeedDegreesPerFrame)
	}
	if h.IsKeyPressed(ebiten.KeyRight) {
		return common.Angle(MovementSpeedDegreesPerFrame)
	}

	// Handle touch/mouse movement if needed
	if h.touchState != nil && h.touchState.Duration > MinTouchDuration {
		// Calculate movement based on touch position relative to screen center
		// This is just an example - adjust the calculation based on your needs
		dx := h.touchState.LastPos.X - h.touchState.StartPos.X
		if math.Abs(dx) > TouchThreshold {
			return common.Angle(math.Copysign(MovementSpeedDegreesPerFrame, dx))
		}
	}

	return 0
}

// IsPausePressed checks if the pause key is pressed
func (h *Handler) IsPausePressed() bool {
	return h.IsKeyPressed(ebiten.KeySpace)
}

// IsQuitPressed checks if the quit key is pressed
func (h *Handler) IsQuitPressed() bool {
	return h.IsKeyPressed(ebiten.KeyEscape)
}

// GetTouchState returns the current touch state
func (h *Handler) GetTouchState() *TouchState {
	var touchIDs []ebiten.TouchID
	touchIDs = ebiten.AppendTouchIDs(touchIDs)

	if len(touchIDs) == 0 {
		return nil
	}

	// Get the first touch point
	x, y := ebiten.TouchPosition(touchIDs[0])
	return &TouchState{
		ID: touchIDs[0],
		StartPos: common.Point{
			X: float64(x),
			Y: float64(y),
		},
		LastPos: common.Point{
			X: float64(x),
			Y: float64(y),
		},
		Duration: 0,
	}
}

// GetMousePosition returns the current mouse position
func (h *Handler) GetMousePosition() common.Point {
	return h.mousePos
}

// IsMouseButtonPressed checks if a mouse button is pressed
func (h *Handler) IsMouseButtonPressed(button ebiten.MouseButton) bool {
	return ebiten.IsMouseButtonPressed(button)
}

// GetLastEvent returns the last input event that occurred
func (h *Handler) GetLastEvent() InputEvent {
	return h.lastEvent
}

// SimulateKeyPress simulates a key press for testing
func (h *Handler) SimulateKeyPress(key ebiten.Key) {
	h.simulatedKeys[key] = true
	h.handleKeyboardInput() // Update lastEvent based on simulated key
}

// SimulateKeyRelease simulates a key release for testing
func (h *Handler) SimulateKeyRelease(key ebiten.Key) {
	delete(h.simulatedKeys, key)
	h.handleKeyboardInput() // Update lastEvent based on remaining keys
}
