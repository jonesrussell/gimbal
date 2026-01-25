package path

import (
	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/ecs/core"
)

// PathCalculator defines interface for parametric path calculation
type PathCalculator interface {
	// Calculate returns the position along the path at the given progress (0.0 to 1.0)
	Calculate(progress float64, start, end common.Point, params core.PathParams) common.Point

	// PathType returns the type of path this calculator handles
	PathType() core.PathType
}

// PathRegistry holds all registered path calculators
type PathRegistry struct {
	calculators map[core.PathType]PathCalculator
}

// NewPathRegistry creates a new path registry with all path types registered
func NewPathRegistry() *PathRegistry {
	registry := &PathRegistry{
		calculators: make(map[core.PathType]PathCalculator),
	}

	// Register all path types
	registry.Register(&SpiralPath{})
	registry.Register(&ArcPath{})
	registry.Register(&StraightPath{})
	registry.Register(&LoopPath{})

	return registry
}

// Register adds a path calculator to the registry
func (r *PathRegistry) Register(calc PathCalculator) {
	r.calculators[calc.PathType()] = calc
}

// Get retrieves a path calculator by type
func (r *PathRegistry) Get(pathType core.PathType) PathCalculator {
	if calc, exists := r.calculators[pathType]; exists {
		return calc
	}
	// Default to straight path if type not found
	return r.calculators[core.PathTypeStraightIn]
}
