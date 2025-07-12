package menu

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

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
	width, height := text.Measure(option, m.font, 0)
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
