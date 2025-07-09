package ecs

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
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
		angleRad := float64(orb.OrbitalAngle) * math.Pi / 180
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
				op.GeoM.Scale(scaleX, scaleY)
			}

			// Apply rotation if angle component exists
			if entry.HasComponent(Angle) {
				angle := Angle.Get(entry)
				// Get sprite center for rotation
				bounds := (*sprite).Bounds()
				centerX := float64(bounds.Dx()) / 2
				centerY := float64(bounds.Dy()) / 2

				op.GeoM.Translate(-centerX, -centerY)
				op.GeoM.Rotate(float64(*angle) * math.Pi / 180)
				op.GeoM.Translate(centerX, centerY)
			}

			// Apply position translation
			op.GeoM.Translate(pos.X, pos.Y)

			screen.DrawImage(*sprite, op)
		}
	})
}

// StarMovementSystem handles star-specific movement (falling down)
func StarMovementSystem(w donburi.World, screenHeight int) {
	query.NewQuery(
		filter.And(
			filter.Contains(StarTag),
			filter.Contains(Position),
			filter.Contains(Speed),
		),
	).Each(w, func(entry *donburi.Entry) {
		pos := Position.Get(entry)
		speed := Speed.Get(entry)

		// Move star downward
		pos.Y += *speed

		// Reset star if it goes off screen
		if pos.Y > float64(screenHeight) {
			pos.Y = 0
		}
	})
}

// PlayerInputSystem handles player input
func PlayerInputSystem(w donburi.World, inputAngle common.Angle) {
	query.NewQuery(
		filter.And(
			filter.Contains(PlayerTag),
			filter.Contains(Orbital),
		),
	).Each(w, func(entry *donburi.Entry) {
		orb := Orbital.Get(entry)

		// Update orbital angle based on input
		if inputAngle != 0 {
			orb.OrbitalAngle += inputAngle
			// Normalize angle to [0, 360)
			if orb.OrbitalAngle < 0 {
				orb.OrbitalAngle += 360
			} else if orb.OrbitalAngle >= 360 {
				orb.OrbitalAngle -= 360
			}
		}
	})
}
