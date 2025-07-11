package ecs

import (
	"image/color"
)

// AnimationConfig holds all animation-related constants
type AnimationConfig struct {
	FadeInDuration   float64
	PulseSpeed       float64
	ChevronSpeed     float64
	ScaleAmplitude   float64
	ChevronAmplitude float64
	TitlePulseSpeed  float64
	HintPulseSpeed   float64
}

// MenuColors defines the color scheme
type MenuColors struct {
	Overlay      color.RGBA
	Title        color.RGBA
	Text         color.RGBA
	SelectedText color.RGBA
	Chevron      color.RGBA
	Background   color.RGBA
	Border       color.RGBA
	HintText     color.RGBA
}

// MenuLayout defines positioning and sizing
type MenuLayout struct {
	TitleY          float64
	TitleScale      float64
	MenuStartY      float64
	MenuItemSpacing float64
	HintY           float64
	ChevronOffsetX  float64
	PaddingX        int
	PaddingY        int
	BorderWidth     int
	HitAreaPadding  int
}

// MenuItem represents a single menu option
type MenuItem struct {
	Text   string
	Option PauseOption
	Action func(*PausedScene)
}

// Configuration factory functions (DRY - single source of truth)
func getAnimationConfig() AnimationConfig {
	return AnimationConfig{
		FadeInDuration:   0.3,
		PulseSpeed:       3.0,
		ChevronSpeed:     2.0,
		ScaleAmplitude:   0.05,
		ChevronAmplitude: 10.0,
		TitlePulseSpeed:  2.0,
		HintPulseSpeed:   1.5,
	}
}

func getMenuColors() MenuColors {
	return MenuColors{
		Overlay:      color.RGBA{0, 0, 0, 128},
		Title:        color.RGBA{51, 204, 255, 255},
		Text:         color.RGBA{255, 255, 255, 255},
		SelectedText: color.RGBA{255, 255, 255, 255},
		Chevron:      color.RGBA{0, 255, 255, 255},
		Background:   color.RGBA{0, 255, 255, 128},
		Border:       color.RGBA{255, 255, 255, 64},
		HintText:     color.RGBA{204, 204, 204, 255},
	}
}

func getMenuLayout(screenWidth, screenHeight int) MenuLayout {
	return MenuLayout{
		TitleY:          80,
		TitleScale:      1.5,
		MenuStartY:      float64(screenHeight) / 2,
		MenuItemSpacing: 50,
		HintY:           float64(screenHeight) - 60,
		ChevronOffsetX:  -140,
		PaddingX:        30,
		PaddingY:        10,
		BorderWidth:     2,
		HitAreaPadding:  50,
	}
}
