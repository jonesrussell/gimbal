package physics

import (
	"math"

	"github.com/jonesrussell/gimbal/internal/common"
)

// CoordinateSystem handles coordinate transformations
type CoordinateSystem struct {
	center common.Point
	radius float64
}

// NewCoordinateSystem creates a new coordinate system
func NewCoordinateSystem(center common.Point, radius float64) *CoordinateSystem {
	return &CoordinateSystem{
		center: center,
		radius: radius,
	}
}

// CalculateCircularPosition calculates a point on a circle given an angle
func (cs *CoordinateSystem) CalculateCircularPosition(angle common.Angle) common.Point {
	rad := angle.ToRadians()
	// In our coordinate system:
	// - 0° points up
	// - 90° points right
	// - 180° points down
	// - 270° points left
	x := cs.center.X + cs.radius*math.Sin(rad) // Use sin for x to match our coordinate system
	y := cs.center.Y - cs.radius*math.Cos(rad) // Use negative cos for y to match our coordinate system

	return common.Point{
		X: x,
		Y: y,
	}
}

// SetPosition sets the center point of the coordinate system
func (cs *CoordinateSystem) SetPosition(pos common.Point) {
	cs.center = pos
}

// GetCenter returns the center point of the coordinate system
func (cs *CoordinateSystem) GetCenter() common.Point {
	return cs.center
}

// GetRadius returns the radius of the coordinate system
func (cs *CoordinateSystem) GetRadius() float64 {
	return cs.radius
}

// CalculateAngle calculates the angle between a point and the center
func (cs *CoordinateSystem) CalculateAngle(pos common.Point) common.Angle {
	dx := pos.X - cs.center.X
	dy := cs.center.Y - pos.Y
	angleRad := math.Atan2(dy, dx)
	return common.FromRadians(angleRad)
}
