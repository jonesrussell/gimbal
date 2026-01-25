package path

import (
	"math"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/ecs/core"
)

// SpiralPath implements spiral entry from center outward
type SpiralPath struct{}

// PathType returns the path type
func (sp *SpiralPath) PathType() core.PathType {
	return core.PathTypeSpiralIn
}

// Calculate returns position along spiral path
// The spiral starts from the start position (near center) and spirals outward to end position
func (sp *SpiralPath) Calculate(progress float64, start, end common.Point, params core.PathParams) common.Point {
	// Calculate center point (midpoint or screen center)
	centerX := (start.X + end.X) / 2
	centerY := (start.Y + end.Y) / 2

	// Calculate target radius (distance from center to end position)
	targetRadius := math.Sqrt(
		(end.X-centerX)*(end.X-centerX) +
			(end.Y-centerY)*(end.Y-centerY),
	)

	// Calculate start radius (distance from center to start position)
	startRadius := math.Sqrt(
		(start.X-centerX)*(start.X-centerX) +
			(start.Y-centerY)*(start.Y-centerY),
	)
	if startRadius < 10 {
		startRadius = params.StartRadius
		if startRadius < 10 {
			startRadius = 10 // Minimum start radius
		}
	}

	// Spiral parameters
	turns := params.SpiralTurns
	if turns == 0 {
		turns = 1.5 // Default 1.5 turns
	}
	direction := float64(params.RotationDirection)
	if direction == 0 {
		direction = 1 // Default clockwise
	}

	// Calculate base angle (angle from center to end position)
	baseAngle := math.Atan2(end.Y-centerY, end.X-centerX)

	// Spiral equation:
	// r(t) = startRadius + (targetRadius - startRadius) * t
	// theta(t) = baseAngle - direction * 2*pi * turns * (1-t)
	// (starts at baseAngle - turns*2*pi, ends at baseAngle)

	currentRadius := startRadius + (targetRadius-startRadius)*progress
	currentAngle := baseAngle - direction*2*math.Pi*turns*(1-progress)

	return common.Point{
		X: centerX + currentRadius*math.Cos(currentAngle),
		Y: centerY + currentRadius*math.Sin(currentAngle),
	}
}
