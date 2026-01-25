package path

import (
	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/ecs/core"
)

// StraightPath implements straight line entry from center outward
type StraightPath struct{}

// PathType returns the path type
func (sp *StraightPath) PathType() core.PathType {
	return core.PathTypeStraightIn
}

// Calculate returns position along straight path
// Simple linear interpolation from start to end
func (sp *StraightPath) Calculate(progress float64, start, end common.Point, params core.PathParams) common.Point {
	// Apply easing for smoother movement (ease out)
	easedProgress := 1 - (1-progress)*(1-progress)

	return common.Point{
		X: start.X + (end.X-start.X)*easedProgress,
		Y: start.Y + (end.Y-start.Y)*easedProgress,
	}
}
