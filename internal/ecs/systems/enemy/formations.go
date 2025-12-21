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

// calculateDiamondFormation calculates positions for a diamond formation
func calculateDiamondFormation(params FormationParams) []FormationData {
	positions := make([]FormationData, params.EnemyCount)
	diamondSpacing := 30.0

	// Calculate how many enemies per side of diamond
	// Diamond has 4 sides: top, right, bottom, left
	enemiesPerSide := params.EnemyCount / 4
	remainder := params.EnemyCount % 4

	idx := 0
	// Top side
	topCount := enemiesPerSide
	if remainder > 0 {
		topCount++
		remainder--
	}
	for i := 0; i < topCount && idx < params.EnemyCount; i++ {
		offsetX := float64(i-topCount/2) * diamondSpacing
		positions[idx] = FormationData{
			Position: common.Point{X: params.CenterX + offsetX, Y: params.CenterY - params.SpawnRadius},
			Angle:    params.BaseAngle + stdmath.Pi/2, // Move downward
		}
		idx++
	}

	// Right side
	rightCount := enemiesPerSide
	if remainder > 0 {
		rightCount++
		remainder--
	}
	for i := 0; i < rightCount && idx < params.EnemyCount; i++ {
		offsetY := float64(i-rightCount/2) * diamondSpacing
		positions[idx] = FormationData{
			Position: common.Point{X: params.CenterX + params.SpawnRadius, Y: params.CenterY + offsetY},
			Angle:    params.BaseAngle - stdmath.Pi, // Move leftward
		}
		idx++
	}

	// Bottom side
	bottomCount := enemiesPerSide
	if remainder > 0 {
		bottomCount++
		remainder--
	}
	for i := 0; i < bottomCount && idx < params.EnemyCount; i++ {
		offsetX := float64(i-bottomCount/2) * diamondSpacing
		positions[idx] = FormationData{
			Position: common.Point{X: params.CenterX + offsetX, Y: params.CenterY + params.SpawnRadius},
			Angle:    params.BaseAngle - stdmath.Pi/2, // Move upward
		}
		idx++
	}

	// Left side (remaining enemies)
	for i := 0; i < remainder && idx < params.EnemyCount; i++ {
		offsetY := float64(i-remainder/2) * diamondSpacing
		positions[idx] = FormationData{
			Position: common.Point{X: params.CenterX - params.SpawnRadius, Y: params.CenterY + offsetY},
			Angle:    params.BaseAngle, // Move rightward
		}
		idx++
	}

	return positions
}

// calculateDiagonalFormation calculates positions for a diagonal formation
func calculateDiagonalFormation(params FormationParams) []FormationData {
	positions := make([]FormationData, params.EnemyCount)
	diagonalSpacing := 25.0

	// Split enemies into two diagonal lines
	halfCount := params.EnemyCount / 2
	firstLineCount := halfCount + (params.EnemyCount % 2)
	secondLineCount := halfCount

	// First diagonal (top-left to bottom-right)
	for i := 0; i < firstLineCount; i++ {
		offset := float64(i-firstLineCount/2) * diagonalSpacing
		angle1 := params.BaseAngle + stdmath.Pi/4 // 45 degrees
		posX := params.CenterX + offset*stdmath.Cos(angle1)
		posY := params.CenterY + offset*stdmath.Sin(angle1)

		positions[i] = FormationData{
			Position: common.Point{X: posX, Y: posY},
			Angle:    angle1 + stdmath.Pi, // Move along diagonal
		}
	}

	// Second diagonal (top-right to bottom-left)
	for i := 0; i < secondLineCount; i++ {
		offset := float64(i-secondLineCount/2) * diagonalSpacing
		angle2 := params.BaseAngle - stdmath.Pi/4 // -45 degrees
		posX := params.CenterX + offset*stdmath.Cos(angle2)
		posY := params.CenterY + offset*stdmath.Sin(angle2)

		positions[firstLineCount+i] = FormationData{
			Position: common.Point{X: posX, Y: posY},
			Angle:    angle2 + stdmath.Pi, // Move along diagonal
		}
	}

	return positions
}

// calculateSpiralFormation calculates positions for a spiral formation
func calculateSpiralFormation(params FormationParams) []FormationData {
	positions := make([]FormationData, params.EnemyCount)
	spiralTurns := 1.5 // Number of full turns in spiral
	minRadius := 20.0
	maxRadius := params.SpawnRadius

	for i := 0; i < params.EnemyCount; i++ {
		// Calculate progress through spiral (0 to 1)
		progress := float64(i) / float64(params.EnemyCount-1)
		if params.EnemyCount == 1 {
			progress = 0
		}

		// Calculate radius (increases with progress)
		radius := minRadius + progress*(maxRadius-minRadius)

		// Calculate angle (spirals outward)
		angle := params.BaseAngle + progress*spiralTurns*2*stdmath.Pi

		posX := params.CenterX + stdmath.Cos(angle)*radius
		posY := params.CenterY + stdmath.Sin(angle)*radius

		// Move outward radially
		positions[i] = FormationData{
			Position: common.Point{X: posX, Y: posY},
			Angle:    angle,
		}
	}

	return positions
}

// calculateRandomFormation calculates positions for a random scattered formation
func calculateRandomFormation(params FormationParams) []FormationData {
	positions := make([]FormationData, params.EnemyCount)
	// Use a simple seed based on formation params for reproducibility
	// In practice, this will be called with different base angles for variety

	for i := 0; i < params.EnemyCount; i++ {
		// Generate pseudo-random angle and radius
		// Using sine/cosine of index to create pseudo-random but deterministic pattern
		angleOffset := float64(i) * 2.34567 // Arbitrary multiplier for variety
		randomAngle := params.BaseAngle + angleOffset
		randomRadius := params.SpawnRadius * (0.5 + 0.5*stdmath.Sin(float64(i)*1.234))

		posX := params.CenterX + stdmath.Cos(randomAngle)*randomRadius
		posY := params.CenterY + stdmath.Sin(randomAngle)*randomRadius

		// Random movement angle (outward from center)
		moveAngle := stdmath.Atan2(posY-params.CenterY, posX-params.CenterX)

		positions[i] = FormationData{
			Position: common.Point{X: posX, Y: posY},
			Angle:    moveAngle,
		}
	}

	return positions
}

// GetSpawnRadius returns appropriate spawn radius based on screen size
func GetSpawnRadius(screenConfig *config.GameConfig) float64 {
	// Spawn enemies close to center, but not exactly at center
	return 50.0
}
