package transitions

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// TransitionType represents different types of transitions
type TransitionType int

const (
	TransitionFade TransitionType = iota
	TransitionSlide
	TransitionFlash
	TransitionWarpTunnel
)

// Transition is an interface for scene transitions
type Transition interface {
	// Update updates the transition state, returns true when complete
	Update(deltaTime float64) bool
	// Draw draws the transition effect
	Draw(screen *ebiten.Image, from, to *ebiten.Image)
	// Reset resets the transition to initial state
	Reset()
	// GetProgress returns the transition progress (0.0 to 1.0)
	GetProgress() float64
}
