package orbital

import (
	"math"

	"github.com/jonesrussell/gimbal/internal/common"
)

// Constants for orbital calculations
const (
	DegreesToRadians = math.Pi / 180
	RadiansToDegrees = 180 / math.Pi
)

// Position represents a position in orbital space
type Position struct {
	Orbital common.Angle // The angle around the orbit
	Facing  common.Angle // The direction the entity is facing
	Point   common.Point // The actual position in 2D space
}

// Config holds orbital movement configuration
type Config struct {
	Center common.Point
	Radius float64
}

// Calculator handles orbital movement calculations
type Calculator struct {
	config Config
}

// NewCalculator creates a new orbital calculator
func NewCalculator(config Config) *Calculator {
	return &Calculator{
		config: config,
	}
}

// CalculatePosition calculates the 2D position given an orbital angle
func (c *Calculator) CalculatePosition(angle common.Angle) common.Point {
	angleRad := float64(angle) * DegreesToRadians
	return common.Point{
		X: c.config.Center.X + c.config.Radius*math.Sin(angleRad),
		Y: c.config.Center.Y - c.config.Radius*math.Cos(angleRad), // Subtract because Y increases downward
	}
}

// NewPosition creates a new orbital position
func (c *Calculator) NewPosition(orbitalAngle, facingAngle common.Angle) Position {
	return Position{
		Orbital: orbitalAngle,
		Facing:  facingAngle,
		Point:   c.CalculatePosition(orbitalAngle),
	}
}

// UpdatePosition updates an existing position with new angles
func (c *Calculator) UpdatePosition(pos *Position, orbitalAngle, facingAngle common.Angle) {
	pos.Orbital = orbitalAngle
	pos.Facing = facingAngle
	pos.Point = c.CalculatePosition(orbitalAngle)
}
