package menu

import (
	"image/color"
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
