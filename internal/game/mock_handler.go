package game

import (
	"fmt"
	"net/http"

	input "github.com/quasilyte/ebitengine-input"
)

// MockHandler is a mock implementation of the input.Handler interface
// for use in unit tests.
type MockHandler struct {
	GamepadDeadzone float64
	pressedAction   input.Action
}

func (m *MockHandler) SetPressedAction(action input.Action) {
	m.pressedAction = action
	fmt.Println("Action set: ", action) // Debugging print statement
}

func (m *MockHandler) ActionIsPressed(action input.Action) bool {
	result := action == m.pressedAction
	fmt.Println("Action checked: ", action, " Result: ", result) // Debugging print statement
	return result
}

// Add other methods of input.Handler that you need to mock...
// These methods are placeholders and may need to be implemented based on your usage.
func (m *MockHandler) Remap(keymap input.Keymap) {
	// Implement this method if necessary
}

func (m *MockHandler) GamepadConnected() bool {
	// Implement this method if necessary
	return false
}

func (m *MockHandler) TouchEventsEnabled() bool {
	// Implement this method if necessary
	return false
}

func (m *MockHandler) CursorPos() input.Vec {
	// Implement this method if necessary
	return input.Vec{}
}

func (m *MockHandler) DefaultInputMask() input.DeviceKind {
	// Implement this method if necessary
	return 0
}

func (m *MockHandler) EmitKeyEvent(e input.SimulatedKeyEvent) {
	// Implement this method if necessary
}

func (m *MockHandler) EmitEvent(e input.SimulatedAction) {
	// Implement this method if necessary
}

func (m *MockHandler) ActionKeyNames(action input.Action, mask input.DeviceKind) []string {
	// Implement this method if necessary
	return nil
}

func (m *MockHandler) JustPressedActionInfo(action input.Action) (input.EventInfo, bool) {
	// Implement this method if necessary
	return input.EventInfo{}, false
}

func (m *MockHandler) PressedActionInfo(action input.Action) (input.EventInfo, bool) {
	// Implement this method if necessary
	return input.EventInfo{}, false
}

func (m *MockHandler) ActionIsJustPressed(action input.Action) bool {
	// Implement this method if necessary
	return false
}

func (m *MockHandler) HandleRequest(req *http.Request) *http.Response {
	return &http.Response{}
}

type MockInput struct{}

func (m *MockInput) HandleRequest(request string) bool {
	// Implement your mock behavior here
	return true
}
