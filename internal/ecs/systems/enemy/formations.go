package enemy

import (
	stdmath "math"

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
	default:
		// Default to circle if unknown
		params.FormationType = FormationCircle
		return calculateCircleFormation(params)
	}
}

// calculateLineFormation calculates positions for a line formation
func calculateLineFormation(params FormationParams) []FormationData {
	positions := make([]FormationData, params.EnemyCount)
	lineLength := float64(params.EnemyCount-1) * 20.0 // 20 pixels between enemies
	startX := params.CenterX - lineLength/2
	startY := params.CenterY

	// Perpendicular angle for line orientation
	lineAngle := params.BaseAngle + stdmath.Pi/2

	for i := 0; i < params.EnemyCount; i++ {
		offsetX := float64(i) * 20.0
		posX := startX + offsetX*stdmath.Cos(lineAngle)
		posY := startY + offsetX*stdmath.Sin(lineAngle)

		positions[i] = FormationData{
			Position: common.Point{X: posX, Y: posY},
			Angle:    params.BaseAngle, // All move in same direction
		}
	}

	return positions
}

// calculateCircleFormation calculates positions for a circle formation
func calculateCircleFormation(params FormationParams) []FormationData {
	positions := make([]FormationData, params.EnemyCount)
	angleStep := 2 * stdmath.Pi / float64(params.EnemyCount)

	for i := 0; i < params.EnemyCount; i++ {
		angle := float64(i)*angleStep + params.BaseAngle
		posX := params.CenterX + stdmath.Cos(angle)*params.SpawnRadius
		posY := params.CenterY + stdmath.Sin(angle)*params.SpawnRadius

		// Each enemy moves outward radially from center
		positions[i] = FormationData{
			Position: common.Point{X: posX, Y: posY},
			Angle:    angle, // Move outward in their spawn direction
		}
	}

	return positions
}

// calculateVFormation calculates positions for a V-formation
func calculateVFormation(params FormationParams) []FormationData {
	positions := make([]FormationData, params.EnemyCount)
	midPoint := params.EnemyCount / 2
	vAngle := stdmath.Pi / 6 // 30 degrees for V spread
	vSpacing := 25.0

	for i := 0; i < params.EnemyCount; i++ {
		var offsetX, offsetY float64

		if i < midPoint {
			// Left side of V
			sideAngle := params.BaseAngle - vAngle
			offsetX = float64(midPoint-i) * vSpacing * stdmath.Cos(sideAngle)
			offsetY = float64(midPoint-i) * vSpacing * stdmath.Sin(sideAngle)
		} else if i > midPoint {
			// Right side of V
			sideAngle := params.BaseAngle + vAngle
			offsetX = float64(i-midPoint) * vSpacing * stdmath.Cos(sideAngle)
			offsetY = float64(i-midPoint) * vSpacing * stdmath.Sin(sideAngle)
		}
		// else: center point (offsetX, offsetY remain 0)

		posX := params.CenterX + offsetX
		posY := params.CenterY + offsetY

		// All move in base direction, maintaining V shape
		positions[i] = FormationData{
			Position: common.Point{X: posX, Y: posY},
			Angle:    params.BaseAngle,
		}
	}

	return positions
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
