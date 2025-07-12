package menu

import (
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// MenuSystem handles common menu functionality
type MenuSystem struct {
	selection     int
	options       []MenuOption
	animationTime float64
	config        MenuConfig
	screenWidth   int
	screenHeight  int
	font          text.Face
}

// NewMenuSystem creates a new menu system with the given options and config
func NewMenuSystem(options []MenuOption, config *MenuConfig, screenWidth, screenHeight int,
	font text.Face,
) *MenuSystem {
	menuConfig := *config // Copy for local modification
	if menuConfig.MenuY == 0 {
		menuConfig.MenuY = float64(screenHeight) / 2
	}

	return &MenuSystem{
		selection:     0,
		options:       options,
		animationTime: 0,
		config:        menuConfig,
		screenWidth:   screenWidth,
		screenHeight:  screenHeight,
		font:          font,
	}
}

// Update handles input and animations
func (m *MenuSystem) Update() {
	m.updateAnimations()
	m.handleKeyboardInput()
	m.handleMouseInput()
}

// updateAnimations updates time-based animations
func (m *MenuSystem) updateAnimations() {
	m.animationTime += 1.0 / 60.0
}

// changeSelection updates the selected menu item
func (m *MenuSystem) changeSelection(newSelection int) {
	if newSelection >= 0 && newSelection < len(m.options) {
		m.selection = newSelection
	}
}

// activateSelection executes the action for the currently selected option
func (m *MenuSystem) activateSelection() {
	if m.selection >= 0 && m.selection < len(m.options) {
		m.options[m.selection].Action()
	}
}

// GetSelection returns the currently selected option index
func (m *MenuSystem) GetSelection() int {
	return m.selection
}

// Reset resets the menu to its initial state
func (m *MenuSystem) Reset() {
	m.selection = 0
	m.animationTime = 0
}
