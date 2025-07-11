package ecs

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// InputHandler handles all input processing
type InputHandler struct {
	scene *PausedScene
}

// NewInputHandler creates a new input handler
func NewInputHandler(scene *PausedScene) *InputHandler {
	return &InputHandler{scene: scene}
}

// HandleInput processes all input events
func (ih *InputHandler) HandleInput() error {
	ih.handleKeyboardNavigation()
	ih.handleMouseInput()
	ih.handleActivation()
	ih.handleQuickActions()
	return nil
}

func (ih *InputHandler) handleKeyboardNavigation() {
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		newSelection := (ih.scene.selection - 1 + len(ih.scene.menuItems)) % len(ih.scene.menuItems)
		ih.scene.changeSelection(newSelection)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		newSelection := (ih.scene.selection + 1) % len(ih.scene.menuItems)
		ih.scene.changeSelection(newSelection)
	}
}

func (ih *InputHandler) handleMouseInput() {
	x, y := ebiten.CursorPosition()
	hoveredItem := ih.getHoveredItem(x, y)

	if hoveredItem != -1 && hoveredItem != ih.scene.selection {
		ih.scene.changeSelection(hoveredItem)
	}

	if hoveredItem != -1 && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		ih.scene.activateSelection()
	}
}

func (ih *InputHandler) handleActivation() {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		ih.scene.activateSelection()
	}
}

func (ih *InputHandler) handleQuickActions() {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		ih.scene.resumeGame()
	}
}

func (ih *InputHandler) getHoveredItem(x, y int) int {
	layout := ih.scene.renderer.layout

	for i, item := range ih.scene.menuItems {
		itemY := layout.MenuStartY + float64(i)*layout.MenuItemSpacing
		width, height := text.Measure(item.Text, defaultFontFace, 0)

		hitArea := ih.calculateHitArea(width, height, itemY, layout)
		if x >= hitArea.x0 && x <= hitArea.x1 && y >= hitArea.y0 && y <= hitArea.y1 {
			return i
		}
	}
	return -1
}

func (ih *InputHandler) calculateHitArea(width, height float64, itemY float64, layout MenuLayout) struct{ x0, y0, x1, y1 int } {
	w := int(width)
	h := int(height)
	centerX := int(float64(ih.scene.manager.config.ScreenSize.Width) / 2)

	return struct{ x0, y0, x1, y1 int }{
		centerX - w/2 - layout.HitAreaPadding,
		int(itemY) - h/2 - 15,
		centerX + w/2 + layout.HitAreaPadding,
		int(itemY) + h/2 + 15,
	}
}
