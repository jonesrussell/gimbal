package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// Theme defines all visual properties for the UI system
type Theme struct {
	Colors struct {
		Text      color.RGBA
		TextLight color.RGBA
		Heart     color.RGBA
		Warning   color.RGBA
		Debug     color.RGBA
	}

	Fonts struct {
		UI    text.Face
		HUD   text.Face
		Title text.Face
	}

	Sizes struct {
		SmallIcon  float64 // 24px
		MediumIcon float64 // 32px
		LargeIcon  float64 // 48px
		Padding    float64 // 8px
		Margin     float64 // 16px
	}
}

// DefaultTheme provides sensible defaults for all UI elements
var DefaultTheme = &Theme{}

func init() {
	// Initialize colors
	DefaultTheme.Colors.Text = color.RGBA{255, 255, 255, 255}      // White
	DefaultTheme.Colors.TextLight = color.RGBA{200, 200, 200, 255} // Light gray
	DefaultTheme.Colors.Heart = color.RGBA{255, 100, 100, 255}     // Red
	DefaultTheme.Colors.Warning = color.RGBA{255, 200, 0, 255}     // Yellow
	DefaultTheme.Colors.Debug = color.RGBA{0, 255, 0, 128}         // Semi-transparent green

	// Initialize sizes (logical units, not pixels)
	DefaultTheme.Sizes.SmallIcon = 24
	DefaultTheme.Sizes.MediumIcon = 32
	DefaultTheme.Sizes.LargeIcon = 48
	DefaultTheme.Sizes.Padding = 8
	DefaultTheme.Sizes.Margin = 16
}

// SetFonts allows setting the font faces after they're loaded
func (t *Theme) SetFonts(ui, hud, title text.Face) {
	t.Fonts.UI = ui
	t.Fonts.HUD = hud
	t.Fonts.Title = title
}
