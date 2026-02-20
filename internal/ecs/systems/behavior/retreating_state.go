package behavior

import (
	"math"

	"github.com/yohamta/donburi"

	"github.com/jonesrussell/gimbal/internal/config"
	"github.com/jonesrussell/gimbal/internal/dbg"
	"github.com/jonesrussell/gimbal/internal/ecs/core"
)

// RetreatingState handles retreat movement back to orbit or off-screen
type RetreatingState struct {
	config *config.GameConfig
}

// NewRetreatingState creates a new retreating state handler
func NewRetreatingState(cfg *config.GameConfig) *RetreatingState {
	return &RetreatingState{config: cfg}
}

// StateType returns the state type
func (rs *RetreatingState) StateType() core.BehaviorStateType {
	return core.StateRetreating
}

// Enter is called when transitioning into this state
func (rs *RetreatingState) Enter(entry *donburi.Entry, data *core.BehaviorStateData) {
	dbg.Log(dbg.State, "Entering retreat state")

	// Set up retreat parameters
	if entry.HasComponent(core.RetreatTimer) {
		retreatData := core.RetreatTimer.Get(entry)
		retreatData.IsRetreating = true

		// Calculate retreat angle (outward from center)
		if entry.HasComponent(core.Position) {
			pos := core.Position.Get(entry)
			centerX := float64(rs.config.ScreenSize.Width) / 2
			centerY := float64(rs.config.ScreenSize.Height) / 2
			retreatData.RetreatAngle = math.Atan2(pos.Y-centerY, pos.X-centerX)
		}

		core.RetreatTimer.SetValue(entry, *retreatData)
	}
}

// Update is called every frame while in this state
func (rs *RetreatingState) Update(entry *donburi.Entry, data *core.BehaviorStateData, deltaTime float64) {
	// Determine retreat behavior based on previous state and attack pattern
	shouldReturnToOrbit := rs.shouldReturnToOrbit(entry, data)

	if shouldReturnToOrbit {
		rs.retreatToOrbit(entry, data, deltaTime)
	} else {
		rs.retreatOffScreen(entry, deltaTime)
	}
}

// shouldReturnToOrbit determines if the entity should return to orbit
func (rs *RetreatingState) shouldReturnToOrbit(entry *donburi.Entry, data *core.BehaviorStateData) bool {
	// Check if max attacks reached
	if data.MaxAttacks > 0 && data.AttackCount >= data.MaxAttacks {
		return false // Exit screen
	}

	// Check health - very low health means exit
	if entry.HasComponent(core.Health) {
		health := core.Health.Get(entry)
		if health.Current <= 0 {
			return false
		}
		healthPercent := float64(health.Current) / float64(health.Maximum)
		if healthPercent < 0.2 {
			return false // Exit when very low health
		}
	}

	// Check retreat timer - if forced retreat due to timeout, exit
	if entry.HasComponent(core.RetreatTimer) {
		retreatData := core.RetreatTimer.Get(entry)
		if retreatData.ElapsedTime >= retreatData.TimeoutDuration && retreatData.TimeoutDuration > 0 {
			return false // Forced retreat = exit
		}
	}

	// Otherwise return to orbit
	return true
}

// retreatToOrbit moves the entity back toward the orbital ring
func (rs *RetreatingState) retreatToOrbit(entry *donburi.Entry, data *core.BehaviorStateData, deltaTime float64) {
	if !entry.HasComponent(core.Position) {
		return
	}

	pos := core.Position.Get(entry)
	centerX := float64(rs.config.ScreenSize.Width) / 2
	centerY := float64(rs.config.ScreenSize.Height) / 2

	// Target orbit radius (same as player orbit)
	targetRadius := float64(rs.config.ScreenSize.Height) / 2 * 0.8

	// Current position relative to center
	dx := pos.X - centerX
	dy := pos.Y - centerY
	currentRadius := math.Sqrt(dx*dx + dy*dy)

	// If already at orbit, we're done
	if math.Abs(currentRadius-targetRadius) < 10 {
		return
	}

	// Get retreat speed
	retreatSpeed := 200.0 // pixels per second
	if entry.HasComponent(core.RetreatTimer) {
		retreatData := core.RetreatTimer.Get(entry)
		if retreatData.RetreatSpeed > 0 {
			retreatSpeed = retreatData.RetreatSpeed
		}
	}

	// Move outward toward orbit radius
	if currentRadius < targetRadius {
		// Move outward
		angle := math.Atan2(dy, dx)
		moveDistance := retreatSpeed * deltaTime
		pos.X += math.Cos(angle) * moveDistance
		pos.Y += math.Sin(angle) * moveDistance
	} else {
		// Move inward
		angle := math.Atan2(dy, dx)
		moveDistance := retreatSpeed * deltaTime
		pos.X -= math.Cos(angle) * moveDistance
		pos.Y -= math.Sin(angle) * moveDistance
	}
}

// retreatOffScreen moves the entity off the screen
func (rs *RetreatingState) retreatOffScreen(entry *donburi.Entry, deltaTime float64) {
	if !entry.HasComponent(core.Position) {
		return
	}

	pos := core.Position.Get(entry)
	centerX := float64(rs.config.ScreenSize.Width) / 2
	centerY := float64(rs.config.ScreenSize.Height) / 2

	// Move outward from center
	dx := pos.X - centerX
	dy := pos.Y - centerY
	distance := math.Sqrt(dx*dx + dy*dy)

	if distance < 1 {
		// At center, pick random direction
		dx = 1
		dy = 0
		distance = 1
	}

	// Normalize
	dx /= distance
	dy /= distance

	// Get retreat speed
	retreatSpeed := 300.0 // Faster for exit
	if entry.HasComponent(core.RetreatTimer) {
		retreatData := core.RetreatTimer.Get(entry)
		if retreatData.RetreatSpeed > 0 {
			retreatSpeed = retreatData.RetreatSpeed * 1.5 // Faster when exiting
		}
	}

	// Move outward
	moveDistance := retreatSpeed * deltaTime
	pos.X += dx * moveDistance
	pos.Y += dy * moveDistance
}

// Exit is called when transitioning out of this state
func (rs *RetreatingState) Exit(entry *donburi.Entry, data *core.BehaviorStateData) {
	dbg.Log(dbg.State, "Exiting retreat state")

	// Reset retreat flag
	if entry.HasComponent(core.RetreatTimer) {
		retreatData := core.RetreatTimer.Get(entry)
		retreatData.IsRetreating = false
		core.RetreatTimer.SetValue(entry, *retreatData)
	}
}

// NextState determines the next state
func (rs *RetreatingState) NextState(entry *donburi.Entry, data *core.BehaviorStateData) core.BehaviorStateType {
	if !entry.HasComponent(core.Position) {
		return core.StateRetreating
	}

	pos := core.Position.Get(entry)
	centerX := float64(rs.config.ScreenSize.Width) / 2
	centerY := float64(rs.config.ScreenSize.Height) / 2

	// Check if off-screen (entity should be removed)
	margin := 100.0
	if pos.X < -margin || pos.X > float64(rs.config.ScreenSize.Width)+margin ||
		pos.Y < -margin || pos.Y > float64(rs.config.ScreenSize.Height)+margin {
		// Entity is off-screen - it should be removed by another system
		return core.StateRetreating // Stay in retreating until removed
	}

	// Check if returning to orbit
	if rs.shouldReturnToOrbit(entry, data) {
		// Check if back at orbit radius
		targetRadius := float64(rs.config.ScreenSize.Height) / 2 * 0.8
		dx := pos.X - centerX
		dy := pos.Y - centerY
		currentRadius := math.Sqrt(dx*dx + dy*dy)

		if math.Abs(currentRadius-targetRadius) < 15 {
			return core.StateOrbiting
		}
	}

	return core.StateRetreating
}
