package transitions

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// FlashTransition implements a screen flash transition
type FlashTransition struct {
	duration     float64
	elapsed      float64
	intensity    float64
	complete     bool
	screenWidth  int
	screenHeight int
}

// NewFlashTransition creates a new flash transition
func NewFlashTransition(duration float64, intensity float64, screenWidth, screenHeight int) *FlashTransition {
	return &FlashTransition{
		duration:     duration,
		intensity:    intensity,
		screenWidth:  screenWidth,
		screenHeight: screenHeight,
	}
}

// Update updates the flash transition
func (f *FlashTransition) Update(deltaTime float64) bool {
	if f.complete {
		return true
	}

	f.elapsed += deltaTime
	if f.elapsed >= f.duration {
		f.elapsed = f.duration
		f.complete = true
	}

	return f.complete
}

// Draw draws the flash transition
func (f *FlashTransition) Draw(screen *ebiten.Image, from, to *ebiten.Image) {
	// Draw the base scene
	if to != nil {
		screen.DrawImage(to, nil)
	} else if from != nil {
		screen.DrawImage(from, nil)
	}

	// Calculate flash intensity using a pulse curve
	progress := f.GetProgress()
	// Use a pulse curve: quick flash in, slow fade out
	flashAlpha := 0.0
	if progress < 0.3 {
		// Quick flash in
		flashAlpha = progress / 0.3
	} else {
		// Slow fade out
		flashAlpha = 1.0 - ((progress - 0.3) / 0.7)
	}
	flashAlpha = math.Max(0.0, math.Min(1.0, flashAlpha)) * f.intensity

	// Draw white flash overlay
	if flashAlpha > 0 {
		overlay := ebiten.NewImage(f.screenWidth, f.screenHeight)
		alpha := uint8(flashAlpha * 255)
		overlay.Fill(color.RGBA{255, 255, 255, alpha})
		screen.DrawImage(overlay, nil)
	}
}

// Reset resets the flash transition
func (f *FlashTransition) Reset() {
	f.elapsed = 0
	f.complete = false
}

// GetProgress returns the transition progress (0.0 to 1.0)
func (f *FlashTransition) GetProgress() float64 {
	if f.duration == 0 {
		return 1.0
	}
	progress := f.elapsed / f.duration
	if progress > 1.0 {
		return 1.0
	}
	return progress
}
