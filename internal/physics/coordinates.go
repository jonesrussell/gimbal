package physics

import (
	"math"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/logger"
)

const (
	// MinAngle represents the minimum allowed angle in radians
	MinAngle = -math.Pi
	// MaxAngle represents the maximum allowed angle in radians
	MaxAngle = 3 * math.Pi / 2
	// RotationOffset represents the rotation offset in radians
	RotationOffset = math.Pi / 2
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
	// - 0° points right
	// - 90° points down
	// - 180° points left
	// - 270° points down
	x := cs.center.X + cs.radius*math.Cos(rad)
	y := cs.center.Y - cs.radius*math.Sin(rad) // Negative sin to match our coordinate system

	logger.GlobalLogger.Debug("Calculating circular position",
		"center", map[string]float64{
			"x": cs.center.X,
			"y": cs.center.Y,
		},
		"radius", cs.radius,
		"angle_rad", rad,
		"angle_deg", rad/common.DegreesToRadians,
		"cos", math.Cos(rad),
		"sin", math.Sin(rad),
		"result", map[string]float64{
			"x": x,
			"y": y,
		},
	)

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
