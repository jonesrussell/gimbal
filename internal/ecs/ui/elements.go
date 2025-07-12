package ui

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// UIElement is the interface that all UI elements must implement
type UIElement interface {
	Draw(renderer *UIRenderer, pos Position)
	GetSize() (width, height float64)
}

// TextElement represents a text element with consistent styling
type TextElement struct {
	text  string
	style TextStyle
}

// TextAlignment defines text alignment options
type TextAlignment int

const (
	AlignLeft TextAlignment = iota
	AlignCenter
	AlignRight
)

// TextStyle defines the appearance of text elements
type TextStyle struct {
	Font      text.Face
	Color     color.RGBA
	Size      float64
	Alignment TextAlignment
}

// NewTextWithStyle creates a new text element with custom styling
func NewTextWithStyle(content string, style TextStyle) *TextElement {
	return &TextElement{
		text:  content,
		style: style,
	}
}

// Draw renders the text element
func (te *TextElement) Draw(renderer *UIRenderer, pos Position) {
	if te.style.Font == nil {
		te.style.Font = DefaultTheme.Fonts.UI
	}

	op := &text.DrawOptions{}

	// Handle text alignment
	textWidth, _ := text.Measure(te.text, te.style.Font, 0)
	switch te.style.Alignment {
	case AlignCenter:
		op.GeoM.Translate(pos.X-textWidth/2, pos.Y)
	case AlignRight:
		op.GeoM.Translate(pos.X-textWidth, pos.Y)
	default: // AlignLeft
		op.GeoM.Translate(pos.X, pos.Y)
	}

	text.Draw(renderer.screen, te.text, te.style.Font, op)
}

// GetSize calculates the size of the text element
func (te *TextElement) GetSize() (width, height float64) {
	if te.style.Font == nil {
		te.style.Font = DefaultTheme.Fonts.UI
	}

	// Use ebiten's text/v2.Measure for accurate sizing
	width, height = text.Measure(te.text, te.style.Font, 0)

	return width, height
}

// IconElement represents an icon element with proper sizing
type IconElement struct {
	sprite *ebiten.Image
	size   IconSize
}

// IconSize defines the size of an icon
type IconSize float64

const (
	SmallIcon  IconSize = 24 // Reduced for UI elements
	MediumIcon IconSize = 32 // Reduced for UI elements
	LargeIcon  IconSize = 48 // Reduced for UI elements
)

// NewIcon creates a new icon element
func NewIcon(sprite *ebiten.Image, size IconSize) *IconElement {
	return &IconElement{
		sprite: sprite,
		size:   size,
	}
}

// Draw renders the icon element
func (ie *IconElement) Draw(renderer *UIRenderer, pos Position) {
	if ie.sprite == nil {
		return
	}

	op := &ebiten.DrawImageOptions{}

	// Get original sprite dimensions
	spriteWidth := float64(ie.sprite.Bounds().Dx())
	spriteHeight := float64(ie.sprite.Bounds().Dy())

	// Calculate scale to fit the target size
	scaleX := float64(ie.size) / spriteWidth
	scaleY := float64(ie.size) / spriteHeight

	op.GeoM.Scale(scaleX, scaleY)
	op.GeoM.Translate(pos.X, pos.Y)

	renderer.screen.DrawImage(ie.sprite, op)

	// Draw debug bounds if enabled
	if renderer.debug {
		renderer.drawDebugBounds(pos, float64(ie.size), float64(ie.size))
	}
}

// GetSize returns the size of the icon element
func (ie *IconElement) GetSize() (width, height float64) {
	return float64(ie.size), float64(ie.size)
}

// LivesDisplay represents a complete lives display component
type LivesDisplay struct {
	lives    int
	maxLives int
	style    LivesStyle
}

// LivesStyle defines the appearance of the lives display
type LivesStyle struct {
	IconSize  IconSize
	Spacing   float64
	TextStyle TextStyle
}

// NewLivesDisplay creates a new lives display component
func NewLivesDisplay(lives, maxLives int) *LivesDisplay {
	return &LivesDisplay{
		lives:    lives,
		maxLives: maxLives,
		style: LivesStyle{
			IconSize: MediumIcon,
			Spacing:  DefaultTheme.Sizes.Padding,
			TextStyle: TextStyle{
				Font:  DefaultTheme.Fonts.HUD,
				Color: DefaultTheme.Colors.Text,
				Size:  16,
			},
		},
	}
}

// Draw renders the lives display
func (ld *LivesDisplay) Draw(renderer *UIRenderer, pos Position) {
	// Create horizontal layout for lives display
	layout := NewHorizontalLayout(ld.style.Spacing)

	// Add "Lives:" text
	livesText := NewTextWithStyle("Lives:", ld.style.TextStyle)
	layout.Add(livesText)

	// Add heart icons for each life
	for i := 0; i < ld.lives; i++ {
		heartIcon := NewIcon(renderer.heartSprite, ld.style.IconSize)
		layout.Add(heartIcon)
	}

	// Draw the layout
	layout.Draw(renderer, pos)
}

// GetSize calculates the size of the lives display
func (ld *LivesDisplay) GetSize() (width, height float64) {
	// Create temporary layout to calculate size
	layout := NewHorizontalLayout(ld.style.Spacing)

	// Add "Lives:" text
	livesText := NewTextWithStyle("Lives:", ld.style.TextStyle)
	layout.Add(livesText)

	// Add heart icons for each life
	for i := 0; i < ld.lives; i++ {
		heartIcon := NewIcon(nil, ld.style.IconSize) // nil sprite for size calculation
		layout.Add(heartIcon)
	}

	return layout.GetSize()
}

// ScoreDisplay represents a score display component
type ScoreDisplay struct {
	score int
	style TextStyle
}

// NewScoreDisplay creates a new score display component
func NewScoreDisplay(score int) *ScoreDisplay {
	return &ScoreDisplay{
		score: score,
		style: TextStyle{
			Font:  DefaultTheme.Fonts.HUD,
			Color: DefaultTheme.Colors.Text,
			Size:  16,
		},
	}
}

// Draw renders the score display
func (sd *ScoreDisplay) Draw(renderer *UIRenderer, pos Position) {
	scoreText := fmt.Sprintf("Score: %d", sd.score)
	textElement := NewTextWithStyle(scoreText, sd.style)
	textElement.Draw(renderer, pos)
}

// GetSize calculates the size of the score display
func (sd *ScoreDisplay) GetSize() (width, height float64) {
	scoreText := fmt.Sprintf("Score: %d", sd.score)
	textElement := NewTextWithStyle(scoreText, sd.style)
	return textElement.GetSize()
}
