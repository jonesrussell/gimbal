package enemy

import (
	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/config"
)

// FormationType represents the type of formation
type FormationType int

const (
	// FormationLine spawns enemies in a straight line
	FormationLine FormationType = iota
	// FormationCircle spawns enemies in a circle around center
	FormationCircle
	// FormationV spawns enemies in a V-formation
	FormationV
	// FormationDiamond spawns enemies in a diamond pattern
	FormationDiamond
	// FormationDiagonal spawns enemies in two diagonal lines crossing center
	FormationDiagonal
	// FormationSpiral spawns enemies sequentially along a spiral path
	FormationSpiral
	// FormationRandom spawns enemies in random scattered positions
	FormationRandom
)

// FormationData contains spawn information for a formation
type FormationData struct {
	Position common.Point
	Angle    float64 // Direction of movement
}

// FormationParams contains parameters for calculating a formation
type FormationParams struct {
	FormationType FormationType
	EnemyCount    int
	CenterX       float64
	CenterY       float64
	BaseAngle     float64 // Base rotation angle for the formation
	SpawnRadius   float64 // Distance from center to spawn
}

// CalculateFormation calculates spawn positions and angles for a formation
func CalculateFormation(params FormationParams) []FormationData {
	switch params.FormationType {
	case FormationLine:
		return calculateLineFormation(params)
	case FormationCircle:
		return calculateCircleFormation(params)
	case FormationV:
		return calculateVFormation(params)
	case FormationDiamond:
		return calculateDiamondFormation(params)
	case FormationDiagonal:
		return calculateDiagonalFormation(params)
	case FormationSpiral:
		return calculateSpiralFormation(params)
	case FormationRandom:
		return calculateRandomFormation(params)
	default:
		// Default to circle if unknown
		params.FormationType = FormationCircle
		return calculateCircleFormation(params)
	}
}

// GetFormationTypeFromIndex returns a formation type based on wave index
func GetFormationTypeFromIndex(waveIndex int) FormationType {
	switch waveIndex % 3 {
	case 0:
		return FormationLine
	case 1:
		return FormationCircle
	case 2:
		return FormationV
	default:
		return FormationCircle
	}
}

// GetSpawnRadius returns appropriate spawn radius based on screen size
func GetSpawnRadius(screenConfig *config.GameConfig) float64 {
	// Spawn enemies close to center, but not exactly at center
	return 50.0
}
