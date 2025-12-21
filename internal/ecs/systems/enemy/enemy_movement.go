package enemy

import (
	stdmath "math"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/ecs/core"
)

func (es *EnemySystem) updateEnemies(deltaTime float64) {
	query.NewQuery(
		filter.And(
			filter.Contains(core.EnemyTag),
			filter.Contains(core.Position),
			filter.Contains(core.Movement),
		),
	).Each(es.world, func(entry *donburi.Entry) {
		pos := core.Position.Get(entry)
		mov := core.Movement.Get(entry)

		// Update pattern time
		mov.PatternTime += deltaTime

		// Apply movement pattern
		velocity := es.applyMovementPattern(*mov)

		// Velocity is in pixels per frame (at 60fps), scale by deltaTime
		// deltaTime is typically 1/60 seconds, so multiply by 60 to get frame-equivalent
		frameScale := deltaTime * 60.0
		pos.X += velocity.X * frameScale
		pos.Y += velocity.Y * frameScale

		// Update movement component with new pattern time
		core.Movement.SetValue(entry, *mov)

		// Remove enemies when they move too far from center (Gyruss-style)
		centerX := float64(es.gameConfig.ScreenSize.Width) / 2
		centerY := float64(es.gameConfig.ScreenSize.Height) / 2
		distanceFromCenter := stdmath.Sqrt((pos.X-centerX)*(pos.X-centerX) + (pos.Y-centerY)*(pos.Y-centerY))
		screenWidth := float64(es.gameConfig.ScreenSize.Width)
		screenHeight := float64(es.gameConfig.ScreenSize.Height)
		maxDistance := stdmath.Max(screenWidth, screenHeight) * 0.8

		if distanceFromCenter > maxDistance {
			es.world.Remove(entry.Entity())
		}
	})
}

// applyMovementPattern applies the movement pattern to calculate velocity
func (es *EnemySystem) applyMovementPattern(mov core.MovementData) common.Point {
	pattern := MovementPattern(mov.Pattern)

	switch pattern {
	case MovementPatternZigzag:
		return es.calculateZigzagVelocity(mov)
	case MovementPatternAccelerating:
		return es.calculateAcceleratingVelocity(mov)
	case MovementPatternPulsing:
		return es.calculatePulsingVelocity(mov)
	default:
		// Normal movement
		return mov.Velocity
	}
}

// calculateZigzagVelocity calculates zigzag movement (oscillates side-to-side)
func (es *EnemySystem) calculateZigzagVelocity(mov core.MovementData) common.Point {
	// Zigzag frequency (how fast it oscillates)
	zigzagFreq := 3.0      // oscillations per second
	zigzagAmplitude := 0.3 // how much it deviates

	// Calculate perpendicular angle (90 degrees to base direction)
	perpendicularAngle := mov.BaseAngle + stdmath.Pi/2

	// Oscillate perpendicular to movement direction
	oscillation := stdmath.Sin(mov.PatternTime*zigzagFreq*2*stdmath.Pi) * zigzagAmplitude

	// Base velocity
	baseVelX := stdmath.Cos(mov.BaseAngle) * mov.BaseSpeed
	baseVelY := stdmath.Sin(mov.BaseAngle) * mov.BaseSpeed

	// Add perpendicular oscillation
	perpendicularX := stdmath.Cos(perpendicularAngle) * oscillation * mov.BaseSpeed
	perpendicularY := stdmath.Sin(perpendicularAngle) * oscillation * mov.BaseSpeed

	return common.Point{
		X: baseVelX + perpendicularX,
		Y: baseVelY + perpendicularY,
	}
}

// calculateAcceleratingVelocity calculates accelerating movement (starts slow, speeds up)
func (es *EnemySystem) calculateAcceleratingVelocity(mov core.MovementData) common.Point {
	// Acceleration factor (0 to 1, where 1 is max speed)
	accelTime := 2.0 // seconds to reach max speed
	accelFactor := stdmath.Min(1.0, mov.PatternTime/accelTime)

	// Start at 30% speed, accelerate to 100%
	speedMultiplier := 0.3 + (accelFactor * 0.7)
	currentSpeed := mov.BaseSpeed * speedMultiplier

	return common.Point{
		X: stdmath.Cos(mov.BaseAngle) * currentSpeed,
		Y: stdmath.Sin(mov.BaseAngle) * currentSpeed,
	}
}

// calculatePulsingVelocity calculates pulsing movement (fast-slow-fast bursts)
func (es *EnemySystem) calculatePulsingVelocity(mov core.MovementData) common.Point {
	// Pulse frequency (how often it pulses)
	pulseFreq := 2.0 // pulses per second
	pulsePhase := mov.PatternTime * pulseFreq * 2 * stdmath.Pi

	// Use sine wave to create smooth pulsing (0.5 to 1.0 speed multiplier)
	speedMultiplier := 0.5 + 0.5*stdmath.Sin(pulsePhase)
	currentSpeed := mov.BaseSpeed * speedMultiplier

	return common.Point{
		X: stdmath.Cos(mov.BaseAngle) * currentSpeed,
		Y: stdmath.Sin(mov.BaseAngle) * currentSpeed,
	}
}
