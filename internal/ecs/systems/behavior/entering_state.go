package behavior

import (
	"github.com/yohamta/donburi"

	"github.com/jonesrussell/gimbal/internal/dbg"
	"github.com/jonesrussell/gimbal/internal/ecs/core"
)

// EnteringState handles the entry path phase
type EnteringState struct{}

// NewEnteringState creates a new entering state handler
func NewEnteringState() *EnteringState {
	return &EnteringState{}
}

// StateType returns the state type
func (es *EnteringState) StateType() core.BehaviorStateType {
	return core.StateEntering
}

// Enter is called when transitioning into this state
func (es *EnteringState) Enter(entry *donburi.Entry, data *core.BehaviorStateData) {
	dbg.Log(dbg.State, "Entering entry state")
}

// Update is called every frame while in this state
func (es *EnteringState) Update(entry *donburi.Entry, data *core.BehaviorStateData, deltaTime float64) {
	// Entry path movement is handled by PathSystem
	// This state just waits for the entry path to complete
}

// Exit is called when transitioning out of this state
func (es *EnteringState) Exit(entry *donburi.Entry, data *core.BehaviorStateData) {
	dbg.Log(dbg.State, "Exiting entry state")
}

// NextState determines the next state
func (es *EnteringState) NextState(entry *donburi.Entry, data *core.BehaviorStateData) core.BehaviorStateType {
	// Check if entry path is complete
	if entry.HasComponent(core.EntryPath) {
		pathData := core.EntryPath.Get(entry)
		if pathData.IsComplete {
			// Determine next state based on post-entry behavior
			switch data.PostEntryBehavior {
			case core.BehaviorImmediateAttack:
				return core.StateAttacking
			case core.BehaviorHoverCenterThenOrbit:
				return core.StateHovering
			default:
				return core.StateOrbiting
			}
		}
	}

	return core.StateEntering
}
