package player

import "math"

const (
	// Player dimensions
	playerWidth  = 32
	playerHeight = 32

	// Player positioning
	radius = 8.0

	// Game state
	gameStarted = false

	// Movement constants
	MinAngle       = -math.Pi
	MaxAngle       = 3 * math.Pi / 2
	AngleStep      = 0.1
	RotationOffset = math.Pi / 2
)

var (
	// Define center as a struct or use float64 values
	center = struct{ X, Y float64 }{X: 400, Y: 300} // adjust values as needed
)

var Debug bool = false
