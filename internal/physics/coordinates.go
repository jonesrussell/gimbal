package physics

import (
	"math"

	"github.com/jonesrussell/gimbal/internal/common"
)

const (
	// MinAngle represents the minimum allowed angle in radians
	MinAngle = -math.Pi
	// MaxAngle represents the maximum allowed angle in radians
	MaxAngle = 3 * math.Pi / 2
	// RotationOffset represents the rotation offset in radians
	RotationOffset = math.Pi / 2
)

// CoordinateSystem handles coordinate transformations and calculations
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

// CalculateCircularPosition calculates a position on a circle given an angle
func (cs *CoordinateSystem) CalculateCircularPosition(angle common.Angle, heightOffset float64) common.Point {
	angleRad := angle.ToRadians()

	// Calculate raw coordinates
	rawX := cs.center.X + cs.radius*math.Cos(angleRad)
	// Subtract sin for screen coordinates (Y increases downward)
	rawY := cs.center.Y - cs.radius*math.Sin(angleRad) - heightOffset

	return common.Point{
		X: math.Round(rawX),
		Y: math.Round(rawY),
	}
}

// CalculateAngle calculates the angle between a point and the center
func (cs *CoordinateSystem) CalculateAngle(pos common.Point) common.Angle {
	dx := cs.center.X - pos.X
	dy := cs.center.Y - pos.Y
	angleRad := math.Atan2(dy, dx) + RotationOffset
	return common.FromRadians(angleRad)
}

// ValidateAngle ensures an angle is within the valid range
func ValidateAngle(angle common.Angle) common.Angle {
	rad := angle.ToRadians()
	if rad < MinAngle {
		return common.FromRadians(MinAngle)
	}
	if rad > MaxAngle {
		return common.FromRadians(MaxAngle)
	}
	return angle
}
