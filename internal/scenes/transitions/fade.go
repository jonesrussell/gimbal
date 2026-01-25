package transitions

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// FadeTransition implements a fade in/out transition
type FadeTransition struct {
	duration     float64
	elapsed      float64
	fromAlpha    float64
	toAlpha      float64
	currentAlpha float64
	complete     bool
}

// NewFadeTransition creates a new fade transition
func NewFadeTransition(duration float64) *FadeTransition {
	return &FadeTransition{
		duration:     duration,
		fromAlpha:    1.0,
		toAlpha:      0.0,
		currentAlpha: 1.0,
	}
}

// Update updates the fade transition
func (f *FadeTransition) Update(deltaTime float64) bool {
	if f.complete {
		return true
	}

	f.elapsed += deltaTime
	progress := f.elapsed / f.duration

	if progress >= 1.0 {
		progress = 1.0
		f.complete = true
	}

	// Fade from fromAlpha to toAlpha
	f.currentAlpha = f.fromAlpha + (f.toAlpha-f.fromAlpha)*progress

	return f.complete
}

// Draw draws the fade transition
func (f *FadeTransition) Draw(screen, from, to *ebiten.Image) {
	// Draw the "to" scene with current alpha
	if to != nil {
		op := &ebiten.DrawImageOptions{}
		op.ColorScale.SetA(float32(f.currentAlpha))
		screen.DrawImage(to, op)
	}

	// Draw the "from" scene with inverse alpha
	if from != nil {
		op := &ebiten.DrawImageOptions{}
		op.ColorScale.SetA(float32(1.0 - f.currentAlpha))
		screen.DrawImage(from, op)
	}

	// Draw fade overlay
	if f.currentAlpha < 1.0 {
		overlay := ebiten.NewImage(screen.Bounds().Dx(), screen.Bounds().Dy())
		overlayAlpha := uint8((1.0 - f.currentAlpha) * 255)
		overlay.Fill(color.RGBA{0, 0, 0, overlayAlpha})
		screen.DrawImage(overlay, nil)
	}
}

// Reset resets the fade transition
func (f *FadeTransition) Reset() {
	f.elapsed = 0
	f.complete = false
	f.currentAlpha = f.fromAlpha
}

// GetProgress returns the transition progress (0.0 to 1.0)
func (f *FadeTransition) GetProgress() float64 {
	if f.duration == 0 {
		return 1.0
	}
	progress := f.elapsed / f.duration
	if progress > 1.0 {
		return 1.0
	}
	return progress
}
