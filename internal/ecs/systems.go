package ecs

import (
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
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
		angleRad := float64(orb.OrbitalAngle) * common.DegreesToRadians
		pos.X = orb.Center.X + orb.Radius*math.Sin(angleRad)
		pos.Y = orb.Center.Y - orb.Radius*math.Cos(angleRad) // Subtract because Y increases downward
	})
}

// RenderSystem handles entity rendering
func RenderSystem(w donburi.World, screen *ebiten.Image) {
	query.NewQuery(
		filter.And(
			filter.Contains(Position),
			filter.Contains(Sprite),
		),
	).Each(w, func(entry *donburi.Entry) {
		pos := Position.Get(entry)
		sprite := Sprite.Get(entry)

		if sprite != nil {
			op := &ebiten.DrawImageOptions{}

			// Apply scaling if size component exists
			if entry.HasComponent(Size) {
				size := Size.Get(entry)
				bounds := (*sprite).Bounds()
				scaleX := float64(size.Width) / float64(bounds.Dx())
				scaleY := float64(size.Height) / float64(bounds.Dy())

				// Apply additional scale if Scale component exists
				if entry.HasComponent(Scale) {
					scale := Scale.Get(entry)
					scaleX *= *scale
					scaleY *= *scale
				}

				op.GeoM.Scale(scaleX, scaleY)
			}

			// Apply rotation if orbital component exists (use facing angle)
			if entry.HasComponent(Orbital) {
				orb := Orbital.Get(entry)
				// Get scaled sprite center for rotation
				var centerX, centerY float64
				if entry.HasComponent(Size) {
					size := Size.Get(entry)
					centerX = float64(size.Width) / 2
					centerY = float64(size.Height) / 2
				} else {
					bounds := (*sprite).Bounds()
					centerX = float64(bounds.Dx()) / 2
					centerY = float64(bounds.Dy()) / 2
				}

				op.GeoM.Translate(-centerX, -centerY)
				op.GeoM.Rotate(float64(orb.FacingAngle) * common.DegreesToRadians)
				op.GeoM.Translate(centerX, centerY)
			} else if entry.HasComponent(Angle) {
				// Fallback to angle component for non-orbital entities
				angle := Angle.Get(entry)
				// Get scaled sprite center for rotation
				var centerX, centerY float64
				if entry.HasComponent(Size) {
					size := Size.Get(entry)
					centerX = float64(size.Width) / 2
					centerY = float64(size.Height) / 2
				} else {
					bounds := (*sprite).Bounds()
					centerX = float64(bounds.Dx()) / 2
					centerY = float64(bounds.Dy()) / 2
				}

				op.GeoM.Translate(-centerX, -centerY)
				op.GeoM.Rotate(float64(*angle) * common.DegreesToRadians)
				op.GeoM.Translate(centerX, centerY)
			}

			// Apply position translation
			op.GeoM.Translate(pos.X, pos.Y)

			screen.DrawImage(*sprite, op)
		}
	})
}

// StarMovementSystem handles star movement in Gyruss-style pattern
func StarMovementSystem(ecs *ecs.ECS) {
	centerX := float64(640) / 2 // TODO: Get from config
	centerY := float64(480) / 2

	// Initialize random seed once
	rand.Seed(time.Now().UnixNano())

	StarTag.Each(ecs.World, func(entry *donburi.Entry) {
		pos := Position.Get(entry)
		speed := Speed.Get(entry)
		scale := Scale.Get(entry)

		// Calculate direction from center to star
		dx := pos.X - centerX
		dy := pos.Y - centerY
		distance := math.Sqrt(dx*dx + dy*dy)

		// Normalize direction vector
		if distance > 0 {
			dx /= distance
			dy /= distance
		}

		// Move star outward from center
		pos.X += dx * *speed
		pos.Y += dy * *speed

		// Scale star based on distance from center (farther = bigger)
		*scale = 0.3 + (distance/200.0)*0.7 // Scale from 0.3 to 1.0

		// Reset star if it goes off screen
		if pos.X < -50 || pos.X > 690 || pos.Y < -50 || pos.Y > 530 {
			// Reset to truly random position along small orbital path
			// Random angle around the circle (0 to 2π)
			angle := rand.Float64() * 2 * math.Pi

			// Random radius within the spawn range (30-80 pixels from center)
			spawnRadius := 30.0 + rand.Float64()*50.0

			pos.X = centerX + math.Cos(angle)*spawnRadius
			pos.Y = centerY + math.Sin(angle)*spawnRadius

			// Reset to random small scale
			*scale = 0.3 + rand.Float64()*0.5
		}
	})
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

// FacingAngleSystem calculates the facing angle for orbital entities
// NOTE: This is now handled in PlayerInputSystem for better performance
func FacingAngleSystem(w donburi.World) {
	// Deprecated: Facing angle calculation moved to PlayerInputSystem
	// This function is kept for backward compatibility but does nothing
}
