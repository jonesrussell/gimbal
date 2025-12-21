package game

import (
	"github.com/jonesrussell/gimbal/internal/ui/state"
)

// updateHUD updates the heads-up display
func (g *ECSGame) updateHUD() {
	current, maximum := g.healthSystem.GetPlayerHealth()
	healthPercent := 1.0
	if maximum > 0 {
		healthPercent = float64(current) / float64(maximum)
	}

	uiData := state.HUDData{
		Score:  g.scoreManager.GetScore(),
		Lives:  current,
		Level:  g.levelManager.GetLevel(),
		Health: healthPercent,
	}

	if hudUI, ok := g.ui.(interface{ UpdateHUD(state.HUDData) }); ok {
		hudUI.UpdateHUD(uiData)
	}
}

