package core

import (
	"math"
	"time"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/jonesrussell/gimbal/internal/common"
)

// MovementSystem handles entity movement
func MovementSystem(w donburi.World) {
	query.NewQuery(
		filter.And(
			filter.Contains(Position),
			filter.Contains(Movement),
		),
	).Each(w, func(entry *donburi.Entry) {
		pos := Position.Get(entry)
		mov := Movement.Get(entry)

		// Apply velocity
		pos.X += mov.Velocity.X
		pos.Y += mov.Velocity.Y

		// Clamp to max speed
		speed := math.Sqrt(mov.Velocity.X*mov.Velocity.X + mov.Velocity.Y*mov.Velocity.Y)
		if speed > mov.MaxSpeed {
			scale := mov.MaxSpeed / speed
			mov.Velocity.X *= scale
			mov.Velocity.Y *= scale
		}
	})
}

// OrbitalMovementSystem handles orbital movement for entities
func OrbitalMovementSystem(w donburi.World) {
	query.NewQuery(
		filter.And(
			filter.Contains(Position),
			filter.Contains(Orbital),
		),
	).Each(w, func(entry *donburi.Entry) {
		pos := Position.Get(entry)
		orb := Orbital.Get(entry)

		// Calculate position based on orbital angle
		angleRad := float64(orb.OrbitalAngle) * 0.017453292519943295 // DegreesToRadians
		pos.X = orb.Center.X + orb.Radius*math.Sin(angleRad)
		pos.Y = orb.Center.Y - orb.Radius*math.Cos(angleRad) // Subtract because Y increases downward
	})
}

// StarMovementSystem handles star movement in Gyruss-style pattern
func StarMovementSystem(ecsInstance *ecs.ECS, config *common.GameConfig) {
	// Create star field helper with configuration from game config (but without speed)
	starConfig := &StarFieldConfig{
		SpawnRadiusMin: config.StarSpawnRadiusMin,
		SpawnRadiusMax: config.StarSpawnRadiusMax,
		Speed:          0, // Will be overridden by individual star speed
		MinScale:       config.StarMinScale,
		MaxScale:       config.StarMaxScale,
		ScaleDistance:  config.StarScaleDistance,
		ResetMargin:    config.StarResetMargin,
		Seed:           time.Now().UnixNano(),
	}
	starHelper := NewStarFieldHelper(starConfig, config.ScreenSize)

	starCount := 0
	StarTag.Each(ecsInstance.World, func(entry *donburi.Entry) {
		pos := Position.Get(entry)
		scale := Scale.Get(entry)
		speed := Speed.Get(entry)

		// Update star using helper with individual speed
		starHelper.UpdateStarWithSpeed(pos, scale, *speed)
		starCount++
	})

	// Removed empty logging branch - not needed for functionality
}

// PlayerInputSystem handles player input
func PlayerInputSystem(w donburi.World, inputAngle common.Angle) {
	movement := &PlayerMovement{}

	query.NewQuery(
		filter.And(
			filter.Contains(PlayerTag),
			filter.Contains(Orbital),
		),
	).Each(w, func(entry *donburi.Entry) {
		orb := Orbital.Get(entry)

		// Update orbital angle
		movement.UpdateOrbitalAngle(orb, inputAngle)

		// Update facing angle
		movement.UpdateFacingAngle(orb)
	})
}
