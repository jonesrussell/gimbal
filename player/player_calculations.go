package player

import (
	"math"
)

// calculateCoordinates is a helper function for initial position calculation
func calculateCoordinates(angle float64) (float64, float64) {
	x := float64(center.X) + radius*math.Cos(angle)
	y := float64(center.Y) - radius*math.Sin(angle) - float64(playerHeight)/2
	return x, y
}
