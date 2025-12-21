package game

import (
	"fmt"

	enemysys "github.com/jonesrussell/gimbal/internal/ecs/systems/enemy"
)

// formatFormationType formats a formation type as a string
func (g *ECSGame) formatFormationType(ft enemysys.FormationType) string {
	switch ft {
	case enemysys.FormationLine:
		return "Line"
	case enemysys.FormationCircle:
		return "Circle"
	case enemysys.FormationV:
		return "V"
	case enemysys.FormationDiamond:
		return "Diamond"
	case enemysys.FormationDiagonal:
		return "Diagonal"
	case enemysys.FormationSpiral:
		return "Spiral"
	case enemysys.FormationRandom:
		return "Random"
	default:
		return "Unknown"
	}
}

// formatEnemyTypes formats enemy types as a string
func (g *ECSGame) formatEnemyTypes(types []enemysys.EnemyType) string {
	if len(types) == 0 {
		return "None"
	}

	typeCounts := make(map[string]int)
	for _, t := range types {
		typeCounts[t.String()]++
	}

	result := ""
	first := true
	for name, count := range typeCounts {
		if !first {
			result += ", "
		}
		if count > 1 {
			result += fmt.Sprintf("%s x%d", name, count)
		} else {
			result += name
		}
		first = false
	}
	return result
}

// formatMovementPattern formats a movement pattern as a string
func (g *ECSGame) formatMovementPattern(mp enemysys.MovementPattern) string {
	switch mp {
	case enemysys.MovementPatternNormal:
		return "Normal"
	case enemysys.MovementPatternZigzag:
		return "Zigzag"
	case enemysys.MovementPatternAccelerating:
		return "Accelerating"
	case enemysys.MovementPatternPulsing:
		return "Pulsing"
	default:
		return "Unknown"
	}
}

