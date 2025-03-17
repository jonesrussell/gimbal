package test

// Screen represents a test screen
type Screen struct {
	width  int
	height int
}

// NewScreen creates a new test screen
func NewScreen(width, height int) *Screen {
	return &Screen{
		width:  width,
		height: height,
	}
}

// Width returns the screen width
func (s *Screen) Width() int {
	return s.width
}

// Height returns the screen height
func (s *Screen) Height() int {
	return s.height
}

// Draw implements the Drawable interface
func (s *Screen) Draw() {
	// No-op for testing
}
