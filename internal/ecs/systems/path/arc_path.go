package path

import (
	"math"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/ecs/core"
)

// ArcPath implements arc sweep entry
type ArcPath struct{}

// PathType returns the path type
func (ap *ArcPath) PathType() core.PathType {
	return core.PathTypeArcSweep
}

// Calculate returns position along arc path
// The arc sweeps from start position around to end position
func (ap *ArcPath) Calculate(progress float64, start, end common.Point, params core.PathParams) common.Point {
	// Calculate center point for the arc
	centerX := (start.X + end.X) / 2
	centerY := (start.Y + end.Y) / 2

	// Calculate angles from center to start and end
	startAngle := math.Atan2(start.Y-centerY, start.X-centerX)
	endAngle := math.Atan2(end.Y-centerY, end.X-centerX)

	// Arc sweep angle (how much to curve)
	arcAngle := params.ArcAngle
	if arcAngle == 0 {
		arcAngle = 180 // Default 180 degree arc
	}
	arcRad := arcAngle * math.Pi / 180

	// Direction of arc
	direction := float64(params.RotationDirection)
	if direction == 0 {
		direction = 1 // Default clockwise
	}

	// Calculate radii
	startRadius := math.Sqrt(
		(start.X-centerX)*(start.X-centerX) +
			(start.Y-centerY)*(start.Y-centerY),
	)
	endRadius := math.Sqrt(
		(end.X-centerX)*(end.X-centerX) +
			(end.Y-centerY)*(end.Y-centerY),
	)
	if startRadius < 10 {
		startRadius = params.StartRadius
		if startRadius < 10 {
			startRadius = 10
		}
	}

	// Interpolate radius
	currentRadius := startRadius + (endRadius-startRadius)*progress

	// Calculate current angle along the arc
	// Start from startAngle, sweep by arcRad amount
	angleOffset := direction * arcRad * progress
	currentAngle := startAngle + angleOffset

	// Also interpolate toward end angle for smooth landing
	blendFactor := progress * progress // Ease into final position
	finalAngle := currentAngle*(1-blendFactor) + endAngle*blendFactor

	// When near the end, snap to end position
	if progress > 0.95 {
		blendToEnd := (progress - 0.95) / 0.05
		return common.Point{
			X: centerX + currentRadius*math.Cos(finalAngle)*(1-blendToEnd) + end.X*blendToEnd,
			Y: centerY + currentRadius*math.Sin(finalAngle)*(1-blendToEnd) + end.Y*blendToEnd,
		}
	}

	return common.Point{
		X: centerX + currentRadius*math.Cos(currentAngle),
		Y: centerY + currentRadius*math.Sin(currentAngle),
	}
}
