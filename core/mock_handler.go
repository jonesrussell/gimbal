package game

import "github.com/hajimehoshi/ebiten/v2"

// MockHandler is a mock implementation of the input.Handler interface
// for use in unit tests.
type MockHandler struct {
	pressedKeys map[ebiten.Key]bool
}

func NewMockHandler() *MockHandler {
	return &MockHandler{pressedKeys: make(map[ebiten.Key]bool)}
}

func (mh *MockHandler) PressKey(key ebiten.Key) {
	mh.pressedKeys[key] = true
}

func (mh *MockHandler) ReleaseKey(key ebiten.Key) {
	mh.pressedKeys[key] = false
}

func (mh *MockHandler) IsKeyPressed(key ebiten.Key) bool {
	return mh.pressedKeys[key]
}
