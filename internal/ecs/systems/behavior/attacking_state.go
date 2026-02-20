package behavior

import (
	"math"
	"time"

	"github.com/yohamta/donburi"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/config"
	"github.com/jonesrussell/gimbal/internal/dbg"
	"github.com/jonesrussell/gimbal/internal/ecs/core"
)

// AttackingState handles attack execution
type AttackingState struct {
	config *config.GameConfig
}

// NewAttackingState creates a new attacking state handler
func NewAttackingState(cfg *config.GameConfig) *AttackingState {
	return &AttackingState{config: cfg}
}

// StateType returns the state type
func (as *AttackingState) StateType() core.BehaviorStateType {
	return core.StateAttacking
}

// Enter is called when transitioning into this state
func (as *AttackingState) Enter(entry *donburi.Entry, data *core.BehaviorStateData) {
	dbg.Log(dbg.State, "Entering attack state")

	// Increment attack count
	data.AttackCount++

	// Initialize attack pattern if present
	if entry.HasComponent(core.AttackPattern) {
		attackData := core.AttackPattern.Get(entry)
		attackData.IsActive = true
		attackData.AttackTimer = 0

		// Calculate target position (toward center for rush attacks)
		as.initializeAttackTarget(entry, attackData)

		core.AttackPattern.SetValue(entry, *attackData)
	}
}

// initializeAttackTarget sets up the attack target position
func (as *AttackingState) initializeAttackTarget(entry *donburi.Entry, attackData *core.AttackPatternData) {
	if !entry.HasComponent(core.Position) {
		return
	}

	pos := core.Position.Get(entry)
	centerX := float64(as.config.ScreenSize.Width) / 2
	centerY := float64(as.config.ScreenSize.Height) / 2

	switch attackData.PatternType {
	case core.AttackSingleRush, core.AttackPairedRush:
		// Rush toward center (but stop before reaching it)
		attackData.TargetPosition = common.Point{
			X: centerX,
			Y: centerY,
		}
		// Save return position
		attackData.ReturnPosition = *pos

	case core.AttackLoopbackRush:
		// Rush through center and loop back
		attackData.TargetPosition = common.Point{
			X: centerX,
			Y: centerY,
		}
		attackData.ReturnPosition = *pos

	case core.AttackSuicideDive:
		// Dive straight at center (no return)
		attackData.TargetPosition = common.Point{
			X: centerX,
			Y: centerY,
		}

	default:
		// Default to center
		attackData.TargetPosition = common.Point{X: centerX, Y: centerY}
	}
}

// Update is called every frame while in this state
func (as *AttackingState) Update(entry *donburi.Entry, data *core.BehaviorStateData, deltaTime float64) {
	if !entry.HasComponent(core.AttackPattern) || !entry.HasComponent(core.Position) {
		return
	}

	attackData := core.AttackPattern.Get(entry)
	pos := core.Position.Get(entry)

	// Update attack timer
	attackData.AttackTimer += time.Duration(deltaTime * float64(time.Second))

	// Execute attack movement
	as.executeAttackMovement(pos, attackData, deltaTime)

	core.AttackPattern.SetValue(entry, *attackData)
}

// executeAttackMovement moves the entity during attack
func (as *AttackingState) executeAttackMovement(pos *common.Point, attackData *core.AttackPatternData, deltaTime float64) {
	// Calculate direction to target
	dx := attackData.TargetPosition.X - pos.X
	dy := attackData.TargetPosition.Y - pos.Y
	distance := math.Sqrt(dx*dx + dy*dy)

	if distance < 1 {
		return // Already at target
	}

	// Normalize direction
	dx /= distance
	dy /= distance

	// Get attack speed
	speed := attackData.RushSpeed
	if speed == 0 {
		speed = 300.0 // Default rush speed
	}

	// Move toward target
	moveDistance := speed * deltaTime
	if moveDistance > distance {
		moveDistance = distance
	}

	pos.X += dx * moveDistance
	pos.Y += dy * moveDistance
}

// Exit is called when transitioning out of this state
func (as *AttackingState) Exit(entry *donburi.Entry, data *core.BehaviorStateData) {
	dbg.Log(dbg.State, "Exiting attack state")

	// Mark attack as inactive
	if entry.HasComponent(core.AttackPattern) {
		attackData := core.AttackPattern.Get(entry)
		attackData.IsActive = false
		core.AttackPattern.SetValue(entry, *attackData)
	}
}

// NextState determines the next state
func (as *AttackingState) NextState(entry *donburi.Entry, data *core.BehaviorStateData) core.BehaviorStateType {
	if !entry.HasComponent(core.AttackPattern) {
		return core.StateOrbiting
	}

	attackData := core.AttackPattern.Get(entry)
	pos := core.Position.Get(entry)

	// Check if attack is complete based on pattern type
	switch attackData.PatternType {
	case core.AttackSuicideDive:
		// Suicide dive ends when reaching center (or destroyed)
		centerX := float64(as.config.ScreenSize.Width) / 2
		centerY := float64(as.config.ScreenSize.Height) / 2
		dist := math.Sqrt(
			(pos.X-centerX)*(pos.X-centerX) +
				(pos.Y-centerY)*(pos.Y-centerY),
		)
		if dist < 20 {
			// Reached center - could trigger explosion or damage
			return core.StateRetreating
		}

	case core.AttackSingleRush, core.AttackPairedRush:
		// Rush ends when reaching near center, then retreat to orbit
		dist := as.distanceToTarget(pos, &attackData.TargetPosition)
		if dist < 50 {
			return core.StateRetreating
		}

	case core.AttackLoopbackRush:
		// Loopback: rush to center, then continue past and loop back
		dist := as.distanceToTarget(pos, &attackData.TargetPosition)
		if dist < 30 {
			// Continue to retreating which will loop back
			return core.StateRetreating
		}

	default:
		// Default: end after attack duration
		if attackData.AttackDuration > 0 && attackData.AttackTimer >= attackData.AttackDuration {
			return core.StateRetreating
		}
	}

	// Check attack timeout (5 seconds max attack time)
	maxAttackTime := 5 * time.Second
	if attackData.AttackTimer >= maxAttackTime {
		return core.StateRetreating
	}

	return core.StateAttacking
}

// distanceToTarget calculates distance between two points
func (as *AttackingState) distanceToTarget(pos, target *common.Point) float64 {
	dx := target.X - pos.X
	dy := target.Y - pos.Y
	return math.Sqrt(dx*dx + dy*dy)
}
