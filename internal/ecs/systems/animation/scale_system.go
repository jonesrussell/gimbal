package animation

import (
	"context"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/ecs/core"
)

// ScaleAnimationSystem handles visual scaling animations
type ScaleAnimationSystem struct {
	world  donburi.World
	logger common.Logger
}

// NewScaleAnimationSystem creates a new scale animation system
func NewScaleAnimationSystem(world donburi.World, logger common.Logger) *ScaleAnimationSystem {
	return &ScaleAnimationSystem{
		world:  world,
		logger: logger,
	}
}

// Update processes all scale animations
func (sas *ScaleAnimationSystem) Update(ctx context.Context, deltaTime float64) error {
	// Check for cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Query all entities with ScaleAnimation and Scale components
	query.NewQuery(
		filter.And(
			filter.Contains(core.ScaleAnimation),
			filter.Contains(core.Scale),
		),
	).Each(sas.world, func(entry *donburi.Entry) {
		sas.updateScaleAnimation(entry, deltaTime)
	})

	return nil
}

// updateScaleAnimation updates a single entity's scale animation
func (sas *ScaleAnimationSystem) updateScaleAnimation(entry *donburi.Entry, deltaTime float64) {
	scaleData := core.ScaleAnimation.Get(entry)

	// Skip if already complete
	if scaleData.IsComplete {
		return
	}

	// Update elapsed time
	scaleData.ElapsedTime += deltaTime

	// Calculate progress
	if scaleData.Duration > 0 {
		scaleData.Progress = scaleData.ElapsedTime / scaleData.Duration
		if scaleData.Progress > 1.0 {
			scaleData.Progress = 1.0
		}
	} else {
		scaleData.Progress = 1.0
	}

	// Calculate current scale using easing
	currentScale := LerpWithEasing(
		scaleData.StartScale,
		scaleData.TargetScale,
		scaleData.Progress,
		scaleData.Easing,
	)

	// Update the Scale component
	scale := core.Scale.Get(entry)
	*scale = currentScale

	// Mark complete if done
	if scaleData.Progress >= 1.0 {
		scaleData.IsComplete = true
		*scale = scaleData.TargetScale // Ensure exact final value
	}

	// Save updated animation data
	core.ScaleAnimation.SetValue(entry, *scaleData)
}

// StartAnimation starts or restarts a scale animation on an entity
func (sas *ScaleAnimationSystem) StartAnimation(
	entry *donburi.Entry,
	startScale, targetScale, duration float64,
	easing core.EasingType,
) {
	// Set initial scale
	if entry.HasComponent(core.Scale) {
		scale := core.Scale.Get(entry)
		*scale = startScale
	}

	// Set animation data
	animData := core.ScaleAnimationData{
		StartScale:  startScale,
		TargetScale: targetScale,
		Progress:    0.0,
		Duration:    duration,
		ElapsedTime: 0.0,
		Easing:      easing,
		IsComplete:  false,
	}

	if entry.HasComponent(core.ScaleAnimation) {
		core.ScaleAnimation.SetValue(entry, animData)
	}
}

// IsAnimationComplete checks if an entity's scale animation is complete
func (sas *ScaleAnimationSystem) IsAnimationComplete(entry *donburi.Entry) bool {
	if !entry.HasComponent(core.ScaleAnimation) {
		return true
	}
	scaleData := core.ScaleAnimation.Get(entry)
	return scaleData.IsComplete
}
