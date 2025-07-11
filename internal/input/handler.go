package input

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/jonesrussell/gimbal/internal/common"
)

// Handler handles input for the game
type Handler struct {
	logger        common.Logger
	touchState    *common.TouchState
	lastEvent     common.InputEvent
	mousePos      common.Point
	simulatedKeys map[ebiten.Key]bool // Track simulated key presses
}

// New creates a new input handler
func New(logger common.Logger) common.GameInputHandler {
	return &Handler{
		logger:        logger,
		touchState:    nil,
		lastEvent:     common.InputEventNone,
		mousePos:      common.Point{},
		simulatedKeys: make(map[ebiten.Key]bool),
	}
}

// HandleInput processes input events
func (h *Handler) HandleInput() {
	h.lastEvent = common.InputEventNone

	// Handle keyboard input
	h.handleKeyboardInput()

	// Handle touch input
	h.handleTouchInput()

	// Handle mouse input
	h.handleMouseInput()

	// Input state logging removed for cleaner output
}

func (h *Handler) handleKeyboardInput() {
	if h.IsKeyPressed(ebiten.KeyLeft) || h.IsKeyPressed(ebiten.KeyRight) {
		h.lastEvent = common.InputEventMove
	} else if h.IsKeyPressed(ebiten.KeyEscape) {
		h.lastEvent = common.InputEventPause
	} else if h.IsKeyPressed(ebiten.KeySpace) {
		h.lastEvent = common.InputEventQuit
	}
}

func (h *Handler) handleTouchInput() {
	// Handle touch input
	var touchIDs []ebiten.TouchID
	touchIDs = ebiten.AppendTouchIDs(touchIDs)
	if len(touchIDs) > 0 {
		if h.touchState == nil {
			// New touch
			x, y := ebiten.TouchPosition(touchIDs[0])
			h.touchState = &common.TouchState{
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
			h.lastEvent = common.InputEventTouch
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
		h.lastEvent = common.InputEventMouseMove
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
		return -common.Angle(common.PlayerMovementSpeed) // Clockwise (left)
	}
	if h.IsKeyPressed(ebiten.KeyRight) {
		return common.Angle(common.PlayerMovementSpeed) // Counterclockwise (right)
	}

	// Handle touch/mouse movement if needed
	if h.touchState != nil && h.touchState.Duration > common.MinTouchDuration {
		// Calculate movement based on touch position relative to screen center
		// This is just an example - adjust the calculation based on your needs
		dx := h.touchState.LastPos.X - h.touchState.StartPos.X
		if math.Abs(dx) > common.TouchThreshold {
			return common.Angle(math.Copysign(float64(common.PlayerMovementSpeed), dx))
		}
	}

	return 0
}

// IsPausePressed checks if the pause key is pressed
func (h *Handler) IsPausePressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyEscape)
}

// IsQuitPressed checks if the quit key is pressed
func (h *Handler) IsQuitPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeySpace)
}

// GetTouchState returns the current touch state
func (h *Handler) GetTouchState() *common.TouchState {
	var touchIDs []ebiten.TouchID
	touchIDs = ebiten.AppendTouchIDs(touchIDs)

	if len(touchIDs) == 0 {
		return nil
	}

	// Get the first touch point
	x, y := ebiten.TouchPosition(touchIDs[0])
	return &common.TouchState{
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

// GetLastEvent returns the last input event
func (h *Handler) GetLastEvent() common.InputEvent {
	return h.lastEvent
}

// SimulateKeyPress simulates a key press for testing
func (h *Handler) SimulateKeyPress(key ebiten.Key) {
	h.simulatedKeys[key] = true
}

// SimulateKeyRelease simulates a key release for testing
func (h *Handler) SimulateKeyRelease(key ebiten.Key) {
	h.simulatedKeys[key] = false
}
