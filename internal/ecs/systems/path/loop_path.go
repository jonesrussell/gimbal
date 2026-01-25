package path

import (
	"math"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/ecs/core"
)

// LoopPath implements loop entry (figure-8 or circular loop)
type LoopPath struct{}

// PathType returns the path type
func (lp *LoopPath) PathType() core.PathType {
	return core.PathTypeLoopEntry
}

// Calculate returns position along loop path
// Creates a looping motion from start to end
func (lp *LoopPath) Calculate(progress float64, start, end common.Point, params core.PathParams) common.Point {
	// Calculate distance for loop size
	distance := math.Sqrt(
		(end.X-start.X)*(end.X-start.X) +
			(end.Y-start.Y)*(end.Y-start.Y),
	)
	loopRadius := distance / 4 // Loop radius is 1/4 of total distance

	// Direction
	direction := float64(params.RotationDirection)
	if direction == 0 {
		direction = 1
	}

	// Base angle from start to end
	baseAngle := math.Atan2(end.Y-start.Y, end.X-start.X)

	// Loop parameters
	curveIntensity := params.CurveIntensity
	if curveIntensity == 0 {
		curveIntensity = 1.0
	}

	// Create a loop by adding sinusoidal offset perpendicular to the path
	// One full loop during the journey
	loopProgress := progress * 2 * math.Pi
	perpAngle := baseAngle + math.Pi/2 // Perpendicular angle

	// Offset from straight line (creates the loop)
	loopOffset := math.Sin(loopProgress) * loopRadius * curveIntensity

	// Linear progress along the path
	linearX := start.X + (end.X-start.X)*progress
	linearY := start.Y + (end.Y-start.Y)*progress

	// Add perpendicular offset for loop effect
	return common.Point{
		X: linearX + math.Cos(perpAngle)*loopOffset*direction,
		Y: linearY + math.Sin(perpAngle)*loopOffset*direction,
	}
}
