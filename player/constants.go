package player

import "math"

const (
	// Player dimensions
	playerWidth  = 32
	playerHeight = 32

	// Player positioning
	radius = 8.0

	// Movement constants
	MinAngle       = -math.Pi
	MaxAngle       = 3 * math.Pi / 2
	AngleStep      = 0.1
	RotationOffset = math.Pi / 2
)

var Debug bool = false
