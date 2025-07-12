package ui

// Position represents a 2D coordinate in the UI system
type Position struct {
	X, Y float64
}

// Size represents the dimensions of a UI element
type Size struct {
	Width, Height float64
}

// Positioning helpers for common UI locations
func TopLeft(x, y float64) Position {
	return Position{X: x, Y: y}
}

// ScreenRelative positioning functions that need screen dimensions
func TopRightRelative(screenWidth, x, y float64) Position {
	return Position{X: screenWidth - x, Y: y}
}

// HorizontalLayout arranges elements horizontally with consistent spacing
type HorizontalLayout struct {
	elements []UIElement
	spacing  float64
}

// Ensure HorizontalLayout implements UIElement
var _ UIElement = (*HorizontalLayout)(nil)

// NewHorizontalLayout creates a new horizontal layout with the specified spacing
func NewHorizontalLayout(spacing float64) *HorizontalLayout {
	return &HorizontalLayout{
		elements: make([]UIElement, 0),
		spacing:  spacing,
	}
}

// Add adds an element to the layout
func (hl *HorizontalLayout) Add(element UIElement) {
	hl.elements = append(hl.elements, element)
}

// GetSize calculates the total size of the layout
func (hl *HorizontalLayout) GetSize() (width, height float64) {
	if len(hl.elements) == 0 {
		return 0, 0
	}

	totalWidth := 0.0
	maxHeight := 0.0

	for i, element := range hl.elements {
		w, h := element.GetSize()
		totalWidth += w
		if h > maxHeight {
			maxHeight = h
		}

		// Add spacing between elements (but not after the last one)
		if i < len(hl.elements)-1 {
			totalWidth += hl.spacing
		}
	}

	return totalWidth, maxHeight
}

// Draw renders all elements in the horizontal layout
func (hl *HorizontalLayout) Draw(renderer *UIRenderer, pos Position) {
	currentX := pos.X

	for _, element := range hl.elements {
		element.Draw(renderer, Position{X: currentX, Y: pos.Y})
		w, _ := element.GetSize()
		currentX += w + hl.spacing
	}
}
