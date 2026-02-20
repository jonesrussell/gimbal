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

// OrbitingState handles orbital movement around the center
type OrbitingState struct {
	config *config.GameConfig
	logger common.Logger
}

// NewOrbitingState creates a new orbiting state handler
func NewOrbitingState(cfg *config.GameConfig, logger common.Logger) *OrbitingState {
	return &OrbitingState{
		config: cfg,
		logger: logger,
	}
}

// StateType returns the state type
func (os *OrbitingState) StateType() core.BehaviorStateType {
	return core.StateOrbiting
}

// Enter is called when transitioning into this state
func (os *OrbitingState) Enter(entry *donburi.Entry, data *core.BehaviorStateData) {
	dbg.Log(dbg.State, "Entering orbit state")

	// Initialize orbital data if not present
	if !entry.HasComponent(core.Orbital) {
		// Calculate current position to determine orbital angle
		if entry.HasComponent(core.Position) {
			pos := core.Position.Get(entry)
			centerX := float64(os.config.ScreenSize.Width) / 2
			centerY := float64(os.config.ScreenSize.Height) / 2

			// Calculate angle from center to current position
			angle := math.Atan2(pos.Y-centerY, pos.X-centerX)
			angleDegrees := angle * 180 / math.Pi

			// Set target orbit angle
			data.TargetOrbitAngle = angleDegrees
		}
	}
}

// Update is called every frame while in this state
func (os *OrbitingState) Update(entry *donburi.Entry, data *core.BehaviorStateData, deltaTime float64) {
	// Update orbital position
	os.updateOrbitalPosition(entry, data, deltaTime)
}

// updateOrbitalPosition moves the entity along its orbital path
func (os *OrbitingState) updateOrbitalPosition(entry *donburi.Entry, data *core.BehaviorStateData, deltaTime float64) {
	if !entry.HasComponent(core.Position) {
		return
	}

	pos := core.Position.Get(entry)
	centerX := float64(os.config.ScreenSize.Width) / 2
	centerY := float64(os.config.ScreenSize.Height) / 2

	// Calculate current angle and radius
	dx := pos.X - centerX
	dy := pos.Y - centerY
	currentRadius := math.Sqrt(dx*dx + dy*dy)
	currentAngle := math.Atan2(dy, dx)

	// Update angle based on orbit speed and direction
	orbitSpeed := data.OrbitSpeed
	if orbitSpeed == 0 {
		orbitSpeed = 45.0 // Default 45 degrees per second
	}
	direction := float64(data.OrbitDirection)
	if direction == 0 {
		direction = 1 // Default clockwise
	}

	// Calculate angle delta
	angleDelta := orbitSpeed * deltaTime * math.Pi / 180 * direction
	newAngle := currentAngle + angleDelta

	// Calculate new position
	pos.X = centerX + currentRadius*math.Cos(newAngle)
	pos.Y = centerY + currentRadius*math.Sin(newAngle)

	// Update target orbit angle for tracking
	data.TargetOrbitAngle = newAngle * 180 / math.Pi
}

// Exit is called when transitioning out of this state
func (os *OrbitingState) Exit(entry *donburi.Entry, data *core.BehaviorStateData) {
	dbg.Log(dbg.State, "Exiting orbit state")
}

// NextState determines the next state
func (os *OrbitingState) NextState(entry *donburi.Entry, data *core.BehaviorStateData) core.BehaviorStateType {
	// Check if behavior is orbit_only - never attack
	if data.PostEntryBehavior == core.BehaviorOrbitOnly {
		// Check for retreat timeout
		if os.shouldRetreat(entry, data) {
			return core.StateRetreating
		}
		return core.StateOrbiting
	}

	// Check if orbit duration has elapsed for attack transition
	if data.StateTime >= data.OrbitDuration && data.OrbitDuration > 0 {
		// Check if max attacks reached
		if data.MaxAttacks > 0 && data.AttackCount >= data.MaxAttacks {
			return core.StateRetreating
		}
		return core.StateAttacking
	}

	// Check for retreat timeout
	if os.shouldRetreat(entry, data) {
		return core.StateRetreating
	}

	return core.StateOrbiting
}

// shouldRetreat checks if the entity should start retreating
func (os *OrbitingState) shouldRetreat(entry *donburi.Entry, data *core.BehaviorStateData) bool {
	// Check retreat timer if present
	if entry.HasComponent(core.RetreatTimer) {
		retreatData := core.RetreatTimer.Get(entry)
		if retreatData.ElapsedTime >= retreatData.TimeoutDuration && retreatData.TimeoutDuration > 0 {
			return true
		}
	}

	// Check health threshold (retreat at low health)
	if entry.HasComponent(core.Health) {
		health := core.Health.Get(entry)
		healthPercent := float64(health.Current) / float64(health.Maximum)
		if healthPercent < 0.2 { // Retreat at 20% health
			return true
		}
	}

	return false
}

// GetDefaultOrbitDuration returns the default orbit duration
func GetDefaultOrbitDuration() time.Duration {
	return 3 * time.Second
}
