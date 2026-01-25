// Package scenes provides constants for scene animations and layouts.
package scenes

import "time"

// Animation timing constants
const (
	// FadeInDuration is the duration for fade-in animations
	FadeInDuration = 300 * time.Millisecond

	// SelectionDelay is the delay before selection changes are registered
	SelectionDelay = 100 * time.Millisecond

	// ScreenShakeDecay is how quickly screen shake reduces per frame
	ScreenShakeDecay = 0.1

	// ScreenShakeIntensity is the initial shake intensity
	ScreenShakeIntensity = 1.0

	// ScreenShakeMultiplier converts intensity to pixel offset
	ScreenShakeMultiplier = 5.0
)

// Layout constants
const (
	// TitleScale is the scale factor for scene titles
	TitleScale = 1.5

	// TitleY is the Y position for titles
	TitleY = 80

	// MenuSpacing is the vertical spacing between menu items
	MenuSpacing = 50

	// OverlayAlpha is the transparency for overlay backgrounds (0-255)
	OverlayAlpha = 128

	// PaddingX is horizontal padding for UI elements
	PaddingX = 30

	// PaddingY is vertical padding for UI elements
	PaddingY = 10

	// HitAreaPadding is extra padding for clickable areas
	HitAreaPadding = 50

	// HintTextY is the Y offset from bottom for hint text
	HintTextY = 60
)

// Alpha/opacity constants (0.0 to 1.0)
const (
	// DimmedAlpha is the alpha for dimmed UI elements
	DimmedAlpha = 0.7

	// HintBaseAlpha is the base alpha for hint text
	HintBaseAlpha = 0.6

	// FullAlpha is full opacity
	FullAlpha = 1.0
)

// Timing constants
const (
	// LevelTitleDuration is how long to show level titles
	LevelTitleDuration = 3 * time.Second

	// LevelTitleFadeTime is the fade in/out duration for titles
	LevelTitleFadeTime = 500 * time.Millisecond
)

// Frame rate constant (for animation calculations)
const (
	// FrameRate is the target frames per second
	FrameRate = 60.0
)
