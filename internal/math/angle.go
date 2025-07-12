package math

import "math"

const (
	// Angle constants
	DegreesInCircle  = 360.0
	DegreesInRadian  = 180.0 / math.Pi
	RadiansInDegree  = math.Pi / 180.0
	DegreesToRadians = math.Pi / 180

	// Standard angles
	AngleUp    = 0
	AngleRight = 90
	AngleDown  = 180
	AngleLeft  = 270

	// MinAngle represents the minimum allowed angle in radians
	MinAngle = -math.Pi
	// MaxAngle represents the maximum allowed angle in radians
	MaxAngle = 3 * math.Pi / 2
	// RotationOffset represents the rotation offset in radians
	RotationOffset = math.Pi / 2
)

// Angle represents an angle in degrees
type Angle float64

// ToRadians converts the angle from degrees to radians
func (a Angle) ToRadians() float64 {
	return float64(a) * RadiansInDegree
}

// FromRadians creates an Angle from radians
func FromRadians(rad float64) Angle {
	return Angle(rad * DegreesInRadian)
}

// Add returns the sum of two angles
func (a Angle) Add(b Angle) Angle {
	return a + b
}

// Sub returns the difference between two angles
func (a Angle) Sub(b Angle) Angle {
	return a - b
}

// Mul returns the product of an angle and a scalar
func (a Angle) Mul(scalar float64) Angle {
	return Angle(float64(a) * scalar)
}

// Div returns the quotient of an angle and a scalar
func (a Angle) Div(scalar float64) Angle {
	return Angle(float64(a) / scalar)
}

// Normalize returns the angle normalized to the range [0, 360)
func (a Angle) Normalize() Angle {
	angle := float64(a)
	for angle < 0 {
		angle += DegreesInCircle
	}
	for angle >= DegreesInCircle {
		angle -= DegreesInCircle
	}
	return Angle(angle)
}

// Validate ensures an angle is within the valid range
func (a Angle) Validate() Angle {
	rad := a.ToRadians()
	if rad < MinAngle {
		return FromRadians(MinAngle)
	}
	if rad > MaxAngle {
		return FromRadians(MaxAngle)
	}
	return a
}
