package player

import "github.com/hajimehoshi/ebiten/v2"

// MockHandler is a mock implementation of the InputHandlerInterface
type MockHandler struct {
	pressedKeys map[ebiten.Key]bool
}

func NewMockHandler() InputHandlerInterface {
	return &MockHandler{
		pressedKeys: make(map[ebiten.Key]bool),
	}
}

// HandleInput implements InputHandlerInterface
func (m *MockHandler) HandleInput() (float64, float64) {
	return 0, 0
}

// IsKeyPressed implements InputHandlerInterface
func (m *MockHandler) IsKeyPressed(key ebiten.Key) bool {
	return m.pressedKeys[key]
}

// PressKey is a test helper to simulate key presses
func (m *MockHandler) PressKey(key ebiten.Key) {
	m.pressedKeys[key] = true
}

// ReleaseKey is a test helper to simulate key releases
func (m *MockHandler) ReleaseKey(key ebiten.Key) {
	m.pressedKeys[key] = false
}
