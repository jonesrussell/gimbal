package game

import (
	"github.com/jonesrussell/gimbal/internal/ui/state"
)

// updateHUD updates the heads-up display using event-driven dirty flag pattern.
// Only pushes updates to UI when state has actually changed.
func (g *ECSGame) updateHUD() {
	// Health changes aren't emitted as events for continuous updates,
	// so we need to check and update the presenter directly
	current, maximum := g.healthSystem.GetPlayerHealth()
	g.hudPresenter.SetHealth(current, maximum)

	// Only update UI if HUD data has changed
	if !g.hudPresenter.IsDirty() {
		return
	}

	uiData := g.hudPresenter.GetData()

	if hudUI, ok := g.ui.(interface{ UpdateHUD(state.HUDData) }); ok {
		hudUI.UpdateHUD(uiData)
	}
}
