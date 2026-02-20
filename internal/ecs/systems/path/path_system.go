package path

import (
	"context"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/config"
	"github.com/jonesrussell/gimbal/internal/dbg"
	"github.com/jonesrussell/gimbal/internal/ecs/core"
)

// PathSystem handles parametric entry path execution for enemies
type PathSystem struct {
	world        donburi.World
	config       *config.GameConfig
	registry     *PathRegistry
	screenCenter common.Point
}

// NewPathSystem creates a new path system
func NewPathSystem(
	world donburi.World,
	cfg *config.GameConfig,
) *PathSystem {
	return &PathSystem{
		world:    world,
		config:   cfg,
		registry: NewPathRegistry(),
		screenCenter: common.Point{
			X: float64(cfg.ScreenSize.Width) / 2,
			Y: float64(cfg.ScreenSize.Height) / 2,
		},
	}
}

// Update processes all entities with entry paths
func (ps *PathSystem) Update(ctx context.Context, deltaTime float64) error {
	// Check for cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Query all entities with EntryPath component
	query.NewQuery(
		filter.Contains(core.EntryPath),
	).Each(ps.world, func(entry *donburi.Entry) {
		ps.updateEntryPath(entry, deltaTime)
	})

	return nil
}

// updateEntryPath updates a single entity's entry path
func (ps *PathSystem) updateEntryPath(entry *donburi.Entry, deltaTime float64) {
	pathData := core.EntryPath.Get(entry)

	// Skip if already complete
	if pathData.IsComplete {
		return
	}

	// Update elapsed time
	pathData.ElapsedTime += deltaTime

	// Calculate progress (0.0 to 1.0)
	if pathData.Duration > 0 {
		pathData.Progress = pathData.ElapsedTime / pathData.Duration
		if pathData.Progress > 1.0 {
			pathData.Progress = 1.0
		}
	} else {
		pathData.Progress = 1.0
	}

	// Get the appropriate path calculator
	calculator := ps.registry.Get(pathData.PathType)

	// Calculate new position along the path
	newPos := calculator.Calculate(
		pathData.Progress,
		pathData.StartPosition,
		pathData.EndPosition,
		pathData.Parameters,
	)

	// Update entity position
	if entry.HasComponent(core.Position) {
		pos := core.Position.Get(entry)
		pos.X = newPos.X
		pos.Y = newPos.Y
	}

	// Update scale animation if present
	ps.updateScaleAnimation(entry, pathData.Progress)

	// Check if path is complete
	if pathData.Progress >= 1.0 {
		pathData.IsComplete = true
		ps.onPathComplete(entry)
	}

	// Save updated path data
	core.EntryPath.SetValue(entry, *pathData)
}

// updateScaleAnimation syncs scale animation with path progress
func (ps *PathSystem) updateScaleAnimation(entry *donburi.Entry, pathProgress float64) {
	if !entry.HasComponent(core.ScaleAnimation) {
		return
	}

	scaleData := core.ScaleAnimation.Get(entry)
	if scaleData.IsComplete {
		return
	}

	// Sync scale animation progress with path progress
	scaleData.Progress = pathProgress

	// Calculate scale using easing
	easedProgress := applyEasing(scaleData.Progress, scaleData.Easing)
	currentScale := scaleData.StartScale + (scaleData.TargetScale-scaleData.StartScale)*easedProgress

	// Update the Scale component
	if entry.HasComponent(core.Scale) {
		scale := core.Scale.Get(entry)
		*scale = currentScale
	}

	// Mark complete if done
	if scaleData.Progress >= 1.0 {
		scaleData.IsComplete = true
	}

	core.ScaleAnimation.SetValue(entry, *scaleData)
}

// onPathComplete handles entry path completion
func (ps *PathSystem) onPathComplete(entry *donburi.Entry) {
	// Transition behavior state from Entering to next state
	if entry.HasComponent(core.BehaviorState) {
		behaviorData := core.BehaviorState.Get(entry)
		behaviorData.PreviousState = behaviorData.CurrentState
		behaviorData.CurrentState = core.StateOrbiting
		behaviorData.StateTime = 0
		core.BehaviorState.SetValue(entry, *behaviorData)

		dbg.Log(dbg.State, "Entry path complete, transitioning to orbiting")
	}

	// Ensure scale is at target
	if entry.HasComponent(core.Scale) && entry.HasComponent(core.ScaleAnimation) {
		scaleData := core.ScaleAnimation.Get(entry)
		scale := core.Scale.Get(entry)
		*scale = scaleData.TargetScale
	}

	// Remove the EntryPath component since it's no longer needed
	// (keeping it for debugging, but marking as complete)
}

// applyEasing applies the easing function to a progress value
func applyEasing(progress float64, easing core.EasingType) float64 {
	switch easing {
	case core.EasingEaseIn:
		return progress * progress
	case core.EasingEaseOut:
		return 1 - (1-progress)*(1-progress)
	case core.EasingEaseInOut:
		if progress < 0.5 {
			return 2 * progress * progress
		}
		return 1 - (-2*progress+2)*(-2*progress+2)/2
	default: // EasingLinear
		return progress
	}
}

// GetScreenCenter returns the screen center point
func (ps *PathSystem) GetScreenCenter() common.Point {
	return ps.screenCenter
}
