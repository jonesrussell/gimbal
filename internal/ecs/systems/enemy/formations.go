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
	formationType := params.FormationType
	enemyCount := params.EnemyCount
	centerX := params.CenterX
	centerY := params.CenterY
	baseAngle := params.BaseAngle
	spawnRadius := params.SpawnRadius
	positions := make([]FormationData, enemyCount)

	switch formationType {
	case FormationLine:
		// Line formation: enemies spawn in a line, all move in same direction
		lineLength := float64(enemyCount-1) * 20.0 // 20 pixels between enemies
		startX := centerX - lineLength/2
		startY := centerY

		// Perpendicular angle for line orientation
		lineAngle := baseAngle + stdmath.Pi/2

		for i := 0; i < enemyCount; i++ {
			offsetX := float64(i) * 20.0
			posX := startX + offsetX*stdmath.Cos(lineAngle)
			posY := startY + offsetX*stdmath.Sin(lineAngle)

			positions[i] = FormationData{
				Position: common.Point{X: posX, Y: posY},
				Angle:    baseAngle, // All move in same direction
			}
		}

	case FormationCircle:
		// Circle formation: enemies spawn evenly around a circle
		angleStep := 2 * stdmath.Pi / float64(enemyCount)

		for i := 0; i < enemyCount; i++ {
			angle := float64(i)*angleStep + baseAngle
			posX := centerX + stdmath.Cos(angle)*spawnRadius
			posY := centerY + stdmath.Sin(angle)*spawnRadius

			// Each enemy moves outward radially from center
			positions[i] = FormationData{
				Position: common.Point{X: posX, Y: posY},
				Angle:    angle, // Move outward in their spawn direction
			}
		}

	case FormationV:
		// V-formation: enemies spawn in a V shape
		midPoint := enemyCount / 2
		vAngle := stdmath.Pi / 6 // 30 degrees for V spread
		vSpacing := 25.0

		for i := 0; i < enemyCount; i++ {
			var offsetX, offsetY float64

			if i < midPoint {
				// Left side of V
				sideAngle := baseAngle - vAngle
				offsetX = float64(midPoint-i) * vSpacing * stdmath.Cos(sideAngle)
				offsetY = float64(midPoint-i) * vSpacing * stdmath.Sin(sideAngle)
			} else if i > midPoint {
				// Right side of V
				sideAngle := baseAngle + vAngle
				offsetX = float64(i-midPoint) * vSpacing * stdmath.Cos(sideAngle)
				offsetY = float64(i-midPoint) * vSpacing * stdmath.Sin(sideAngle)
			} else {
				// Center point
				offsetX = 0
				offsetY = 0
			}

			posX := centerX + offsetX
			posY := centerY + offsetY

			// All move in base direction, maintaining V shape
			positions[i] = FormationData{
				Position: common.Point{X: posX, Y: posY},
				Angle:    baseAngle,
			}
		}

	default:
		// Default to circle if unknown
		circleParams := FormationParams{
			FormationType: FormationCircle,
			EnemyCount:    enemyCount,
			CenterX:       centerX,
			CenterY:       centerY,
			BaseAngle:     baseAngle,
			SpawnRadius:   spawnRadius,
		}
		return CalculateFormation(circleParams)
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
