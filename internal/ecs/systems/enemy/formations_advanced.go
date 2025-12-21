package enemy

import (
	stdmath "math"

	"github.com/jonesrussell/gimbal/internal/common"
)

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
