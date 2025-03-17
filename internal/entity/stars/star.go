package stars

import (
	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/entity/player"
)

// Star represents a star in the game
type Star struct {
	position common.Point
	speed    float64
	size     float64
	sprite   player.Drawable
	bounds   common.Size
	angle    float64
}

// New creates a new star instance
func New(pos common.Point, speed, size float64, sprite player.Drawable) *Star {
	return &Star{
		position: pos,
		speed:    speed,
		size:     size,
		sprite:   sprite,
		bounds:   common.Size{Width: int(size), Height: int(size)},
		angle:    0,
	}
}

// Update implements Entity interface
func (s *Star) Update() {
	// Move star towards center
	s.position.Y += s.speed

	// Reset star if it goes off screen
	if s.position.Y > float64(s.bounds.Height) {
		s.position.Y = 0
	}
}

// Draw implements the Drawable interface
func (s *Star) Draw(screen any) {
	if s.sprite != nil {
		s.sprite.Draw(screen, nil)
	}
}

// GetPosition implements Entity interface
func (s *Star) GetPosition() common.Point {
	return s.position
}

// SetPosition implements Movable interface
func (s *Star) SetPosition(pos common.Point) {
	s.position = pos
}

// GetSpeed implements Movable interface
func (s *Star) GetSpeed() float64 {
	return s.speed
}

// SetSpeed sets the star's speed
func (s *Star) SetSpeed(speed float64) {
	s.speed = speed
}

// GetSize returns the star's size
func (s *Star) GetSize() float64 {
	return s.size
}

// SetSize sets the star's size
func (s *Star) SetSize(size float64) {
	s.size = size
	s.bounds = common.Size{Width: int(size), Height: int(size)}
}

// GetAngle returns the star's angle
func (s *Star) GetAngle() float64 {
	return s.angle
}

// SetAngle sets the star's angle
func (s *Star) SetAngle(angle float64) {
	s.angle = angle
}

// GetSprite returns the star's sprite
func (s *Star) GetSprite() player.Drawable {
	return s.sprite
}

// GetBounds implements Collidable interface
func (s *Star) GetBounds() common.Size {
	return s.bounds
}

// SetBounds sets the screen bounds for the star
func (s *Star) SetBounds(bounds common.Size) {
	s.bounds = bounds
}
