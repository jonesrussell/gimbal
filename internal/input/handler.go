package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/math"
)

const (
	PlayerMovementSpeed = 5 // degrees per frame
	MinTouchDuration    = 10
	TouchThreshold      = 50
)

// Handler implements the GameInputHandler interface
type Handler struct {
	lastEvent     common.InputEvent
	pressedKeys   map[ebiten.Key]bool
	touchState    *common.TouchState
	mousePos      common.Point
	simulatedKeys map[ebiten.Key]bool
}

// NewHandler creates a new input handler
func NewHandler() *Handler {
	return &Handler{
		lastEvent:     common.InputEventNone,
		pressedKeys:   make(map[ebiten.Key]bool),
		simulatedKeys: make(map[ebiten.Key]bool),
		mousePos:      common.Point{X: 0, Y: 0},
	}
}

// HandleInput processes all input events
func (h *Handler) HandleInput() {
	h.updateKeyState()
	h.updateTouchState()
	h.updateMouseState()
	h.updateLastEvent()
}

// updateKeyState updates the state of all keys
func (h *Handler) updateKeyState() {
	// Update pressed keys state
	for key := range h.pressedKeys {
		h.pressedKeys[key] = ebiten.IsKeyPressed(key)
	}

	// Check for new key presses
	allKeys := []ebiten.Key{
		ebiten.KeyA, ebiten.KeyD, ebiten.KeyLeft, ebiten.KeyRight,
		ebiten.KeyEscape, ebiten.KeyP, ebiten.KeySpace,
	}

	for _, key := range allKeys {
		if ebiten.IsKeyPressed(key) {
			h.pressedKeys[key] = true
		}
	}
}

// updateTouchState updates touch input state
func (h *Handler) updateTouchState() {
	touches := inpututil.AppendJustPressedTouchIDs(nil)
	if len(touches) > 0 {
		touchID := touches[0]
		x, y := ebiten.TouchPosition(touchID)
		h.touchState = &common.TouchState{
			ID:       touchID,
			StartPos: common.Point{X: float64(x), Y: float64(y)},
			LastPos:  common.Point{X: float64(x), Y: float64(y)},
			Duration: 0,
		}
	}

	if h.touchState != nil {
		h.touchState.Duration++
	}
}

// updateMouseState updates mouse input state
func (h *Handler) updateMouseState() {
	x, y := ebiten.CursorPosition()
	h.mousePos = common.Point{X: float64(x), Y: float64(y)}
}

// updateLastEvent determines the last input event that occurred
func (h *Handler) updateLastEvent() {
	if h.IsQuitPressed() {
		h.lastEvent = common.InputEventQuit
	} else if h.IsPausePressed() {
		h.lastEvent = common.InputEventPause
	} else if h.GetMovementInput() != 0 {
		h.lastEvent = common.InputEventMove
	} else if h.touchState != nil {
		h.lastEvent = common.InputEventTouch
	} else if h.isAnyInputPressed() {
		h.lastEvent = common.InputEventAny
	} else {
		h.lastEvent = common.InputEventNone
	}
}

// isAnyInputPressed checks if any key or mouse button is pressed
func (h *Handler) isAnyInputPressed() bool {
	// Check for any key press
	for key := ebiten.Key(0); key <= ebiten.KeyMax; key++ {
		if inpututil.IsKeyJustPressed(key) {
			return true
		}
	}

	// Check for any mouse button press
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) ||
		inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) ||
		inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonMiddle) {
		return true
	}

	return false
}

// IsKeyPressed checks if a specific key is currently pressed
func (h *Handler) IsKeyPressed(key ebiten.Key) bool {
	return h.pressedKeys[key] || h.simulatedKeys[key]
}

// GetMovementInput returns the current movement input as an angle
func (h *Handler) GetMovementInput() math.Angle {
	if h.IsKeyPressed(ebiten.KeyA) || h.IsKeyPressed(ebiten.KeyLeft) {
		return math.Angle(-PlayerMovementSpeed)
	}
	if h.IsKeyPressed(ebiten.KeyD) || h.IsKeyPressed(ebiten.KeyRight) {
		return math.Angle(PlayerMovementSpeed)
	}

	// Handle touch input
	if h.touchState != nil && h.touchState.Duration > MinTouchDuration {
		deltaX := h.touchState.LastPos.X - h.touchState.StartPos.X
		if deltaX > TouchThreshold {
			return math.Angle(PlayerMovementSpeed)
		} else if deltaX < -TouchThreshold {
			return math.Angle(-PlayerMovementSpeed)
		}
	}

	return 0
}

// IsQuitPressed checks if the quit key is pressed
func (h *Handler) IsQuitPressed() bool {
	return h.IsKeyPressed(ebiten.KeyEscape)
}

// IsPausePressed checks if the pause key is pressed
func (h *Handler) IsPausePressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyP)
}

// IsShootPressed checks if the shoot key is pressed
func (h *Handler) IsShootPressed() bool {
	return h.IsKeyPressed(ebiten.KeySpace)
}

// GetTouchState returns the current touch state
func (h *Handler) GetTouchState() *common.TouchState {
	return h.touchState
}

// GetMousePosition returns the current mouse position
func (h *Handler) GetMousePosition() common.Point {
	return h.mousePos
}

// IsMouseButtonPressed checks if a specific mouse button is pressed
func (h *Handler) IsMouseButtonPressed(button ebiten.MouseButton) bool {
	return ebiten.IsMouseButtonPressed(button)
}

// GetLastEvent returns the last input event that occurred
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
