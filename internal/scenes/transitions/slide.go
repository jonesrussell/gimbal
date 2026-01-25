package transitions

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// SlideDirection represents the direction of a slide transition
type SlideDirection int

const (
	SlideLeft SlideDirection = iota
	SlideRight
	SlideUp
	SlideDown
)

// SlideTransition implements a slide transition
type SlideTransition struct {
	duration     float64
	elapsed      float64
	direction    SlideDirection
	screenWidth  int
	screenHeight int
	complete     bool
}

// NewSlideTransition creates a new slide transition
func NewSlideTransition(duration float64, direction SlideDirection, screenWidth, screenHeight int) *SlideTransition {
	return &SlideTransition{
		duration:     duration,
		direction:    direction,
		screenWidth:  screenWidth,
		screenHeight: screenHeight,
	}
}

// Update updates the slide transition
func (s *SlideTransition) Update(deltaTime float64) bool {
	if s.complete {
		return true
	}

	s.elapsed += deltaTime
	if s.elapsed >= s.duration {
		s.elapsed = s.duration
		s.complete = true
	}

	return s.complete
}

// Draw draws the slide transition
func (s *SlideTransition) Draw(screen, from, to *ebiten.Image) {
	progress := s.GetProgress()

	var fromX, fromY, toX, toY float64

	switch s.direction {
	case SlideLeft:
		fromX = -progress * float64(s.screenWidth)
		toX = (1.0 - progress) * float64(s.screenWidth)
	case SlideRight:
		fromX = progress * float64(s.screenWidth)
		toX = -(1.0 - progress) * float64(s.screenWidth)
	case SlideUp:
		fromY = -progress * float64(s.screenHeight)
		toY = (1.0 - progress) * float64(s.screenHeight)
	case SlideDown:
		fromY = progress * float64(s.screenHeight)
		toY = -(1.0 - progress) * float64(s.screenHeight)
	}

	// Draw "to" scene
	if to != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(toX, toY)
		screen.DrawImage(to, op)
	}

	// Draw "from" scene
	if from != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(fromX, fromY)
		screen.DrawImage(from, op)
	}
}

// Reset resets the slide transition
func (s *SlideTransition) Reset() {
	s.elapsed = 0
	s.complete = false
}

// GetProgress returns the transition progress (0.0 to 1.0)
func (s *SlideTransition) GetProgress() float64 {
	if s.duration == 0 {
		return 1.0
	}
	progress := s.elapsed / s.duration
	if progress > 1.0 {
		return 1.0
	}
	return progress
}
