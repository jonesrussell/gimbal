package ecs

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// MenuOption represents a menu item with its action
type MenuOption struct {
	Text   string
	Action func()
}

// MenuConfig holds configuration for menu appearance and behavior
type MenuConfig struct {
	// Layout
	MenuY       float64
	ItemSpacing int

	// Visual styling
	PaddingX        int
	PaddingY        int
	HitAreaPadding  int
	HitAreaPaddingY int
	BackgroundAlpha float64
	DimmedAlpha     float64

	// Chevron positioning
	ChevronOffsetX float64
	ChevronOffsetY float64

	// Animation
	PulseSpeed     float64
	PulseAmplitude float64
	PulseBase      float64

	// Colors
	BackgroundColor color.RGBA
	ChevronColor    color.RGBA
	TextColor       color.RGBA
}

// DefaultMenuConfig returns a standard menu configuration
func DefaultMenuConfig() MenuConfig {
	return MenuConfig{
		ItemSpacing:     40,
		PaddingX:        24,
		PaddingY:        6,
		HitAreaPadding:  40,
		HitAreaPaddingY: 8,
		BackgroundAlpha: 0.5,
		DimmedAlpha:     1.0,
		ChevronOffsetX:  120,
		ChevronOffsetY:  8,
		PulseSpeed:      2.0,
		PulseAmplitude:  0.3,
		PulseBase:       0.7,
		BackgroundColor: color.RGBA{0, 255, 255, 128},
		ChevronColor:    color.RGBA{0, 255, 255, 255},
		TextColor:       color.RGBA{255, 255, 255, 255},
	}
}

// PausedMenuConfig returns configuration optimized for pause menus
func PausedMenuConfig() MenuConfig {
	config := DefaultMenuConfig()
	config.ItemSpacing = 50
	config.PaddingX = 30
	config.PaddingY = 10
	config.HitAreaPadding = 50
	config.HitAreaPaddingY = 15
	config.DimmedAlpha = 0.7
	config.PulseSpeed = 3.0
	config.ChevronOffsetX = 140
	return config
}

// MenuSystem handles common menu functionality
type MenuSystem struct {
	selection     int
	options       []MenuOption
	animationTime float64
	config        MenuConfig
	screenWidth   int
	screenHeight  int
}

// NewMenuSystem creates a new menu system with the given options and config
func NewMenuSystem(options []MenuOption, config MenuConfig, screenWidth, screenHeight int) *MenuSystem {
	menuConfig := config
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

// handleKeyboardInput processes keyboard navigation and actions
func (m *MenuSystem) handleKeyboardInput() {
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		m.changeSelection((m.selection - 1 + len(m.options)) % len(m.options))
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		m.changeSelection((m.selection + 1) % len(m.options))
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		m.activateSelection()
	}
}

// handleMouseInput processes mouse hover and click events
func (m *MenuSystem) handleMouseInput() {
	x, y := ebiten.CursorPosition()
	hoveredItem := -1

	for i := range m.options {
		itemY := m.config.MenuY + float64(i*m.config.ItemSpacing)
		if m.isMouseOverItem(x, y, m.options[i].Text, itemY) {
			hoveredItem = i
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				m.activateSelection()
			}
		}
	}

	if hoveredItem != -1 && hoveredItem != m.selection {
		m.changeSelection(hoveredItem)
	}
}

// isMouseOverItem checks if the mouse cursor is over a menu item
func (m *MenuSystem) isMouseOverItem(x, y int, option string, itemY float64) bool {
	width, height := text.Measure(option, defaultFontFace, 0)
	w := int(width)
	h := int(height)

	itemRect := struct{ x0, y0, x1, y1 int }{
		int(float64(m.screenWidth)/2) - w/2 - m.config.HitAreaPadding,
		int(itemY) - h/2 - m.config.HitAreaPaddingY,
		int(float64(m.screenWidth)/2) + w/2 + m.config.HitAreaPadding,
		int(itemY) + h/2 + m.config.HitAreaPaddingY,
	}

	return x >= itemRect.x0 && x <= itemRect.x1 && y >= itemRect.y0 && y <= itemRect.y1
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

// Draw renders the menu options
func (m *MenuSystem) Draw(screen *ebiten.Image, fadeAlpha float64) {
	for i, option := range m.options {
		y := m.config.MenuY + float64(i*m.config.ItemSpacing)
		isSelected := i == m.selection
		m.drawMenuOption(screen, option.Text, y, isSelected, fadeAlpha)
	}
}

// drawMenuOption renders a single menu option
func (m *MenuSystem) drawMenuOption(
	screen *ebiten.Image,
	option string,
	y float64,
	isSelected bool,
	fadeAlpha float64,
) {
	alpha := fadeAlpha
	scale := 1.0

	if isSelected {
		pulse := m.config.PulseBase + m.config.PulseAmplitude*math.Sin(m.animationTime*m.config.PulseSpeed)
		alpha *= pulse
		scale = 1.0 + 0.05*math.Sin(m.animationTime*m.config.PulseSpeed)

		m.drawSelectionBackground(screen, option, y, m.config.BackgroundAlpha*fadeAlpha)
		m.drawChevron(screen, y, fadeAlpha)
	} else {
		alpha *= m.config.DimmedAlpha
	}

	m.drawOptionText(screen, option, y, alpha, scale)
}

// drawSelectionBackground renders the background highlight for selected items
func (m *MenuSystem) drawSelectionBackground(screen *ebiten.Image, option string, y, alpha float64) {
	width, height := text.Measure(option, defaultFontFace, 0)
	w := int(width)
	h := int(height)

	bgColor := m.config.BackgroundColor
	bgColor.A = uint8(float64(bgColor.A) * alpha)

	rect := ebiten.NewImage(w+m.config.PaddingX*2, h+m.config.PaddingY*2)
	rect.Fill(bgColor)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(
		float64(m.screenWidth)/2-float64(w+m.config.PaddingX*2)/2,
		y-float64(h+m.config.PaddingY*2)/2+2,
	)
	screen.DrawImage(rect, op)
}

// drawChevron renders the animated selection chevron
func (m *MenuSystem) drawChevron(screen *ebiten.Image, y, fadeAlpha float64) {
	pulse := m.config.PulseBase + m.config.PulseAmplitude*math.Sin(m.animationTime*4.0)
	chevronX := float64(m.screenWidth)/2 - m.config.ChevronOffsetX
	chevronY := y + m.config.ChevronOffsetY

	op := &text.DrawOptions{}
	op.GeoM.Translate(chevronX, chevronY)
	op.ColorScale.SetR(float32(m.config.ChevronColor.R) / 255.0)
	op.ColorScale.SetG(float32(m.config.ChevronColor.G) / 255.0)
	op.ColorScale.SetB(float32(m.config.ChevronColor.B) / 255.0)
	op.ColorScale.SetA(float32(pulse * fadeAlpha))

	text.Draw(screen, ">", defaultFontFace, op)
}

// drawOptionText renders the text for a menu option
func (m *MenuSystem) drawOptionText(screen *ebiten.Image, option string, y, alpha, scale float64) {
	op := &text.DrawOptions{}

	width, height := text.Measure(option, defaultFontFace, 0)
	centerX := float64(m.screenWidth) / 2
	centerY := y

	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(
		centerX-width*scale/2,
		centerY-height*scale/2,
	)

	textColor := m.config.TextColor
	op.ColorScale.SetR(float32(textColor.R) / 255.0)
	op.ColorScale.SetG(float32(textColor.G) / 255.0)
	op.ColorScale.SetB(float32(textColor.B) / 255.0)
	op.ColorScale.SetA(float32(alpha))

	text.Draw(screen, option, defaultFontFace, op)
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
