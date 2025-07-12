package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// UIRenderer is the main UI rendering system
type UIRenderer struct {
	screen      *ebiten.Image
	theme       *Theme
	heartSprite *ebiten.Image
	debug       bool
}

// NewUIRenderer creates a new UI renderer
func NewUIRenderer(screen *ebiten.Image, theme *Theme) *UIRenderer {
	if theme == nil {
		theme = DefaultTheme
	}

	return &UIRenderer{
		screen: screen,
		theme:  theme,
		debug:  false,
	}
}

// SetHeartSprite sets the heart sprite for lives display
func (ur *UIRenderer) SetHeartSprite(sprite *ebiten.Image) {
	ur.heartSprite = sprite
}

// SetScreen sets the screen for rendering
func (ur *UIRenderer) SetScreen(screen *ebiten.Image) {
	ur.screen = screen
}

// SetDebug enables or disables debug mode
func (ur *UIRenderer) SetDebug(debug bool) {
	ur.debug = debug
}

// Draw renders a UI element at the specified position
func (ur *UIRenderer) Draw(element UIElement, pos Position) {
	element.Draw(ur, pos)
}

// drawDebugBounds draws debug bounding boxes around UI elements
func (ur *UIRenderer) drawDebugBounds(pos Position, width, height float64) {
	// Create a simple rectangle for debug visualization
	// In a full implementation, you'd draw actual rectangles
	// For now, we'll just log the bounds
	// TODO: Implement actual debug rectangle drawing
}

// GetTheme returns the current theme
func (ur *UIRenderer) GetTheme() *Theme {
	return ur.theme
}

// GetScreenWidth returns the width of the current screen
func (ur *UIRenderer) GetScreenWidth() float64 {
	if ur.screen == nil {
		return 0
	}
	return float64(ur.screen.Bounds().Dx())
}
