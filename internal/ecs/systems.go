package ecs

import (
	"math"
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

// applySpriteTransform applies scaling and rotation to the DrawImageOptions for an entity
func applySpriteTransform(entry *donburi.Entry, sprite *ebiten.Image, op *ebiten.DrawImageOptions) {
	// Apply scaling if size component exists
	if entry.HasComponent(Size) {
		size := Size.Get(entry)
		bounds := sprite.Bounds()
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
			bounds := sprite.Bounds()
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
			bounds := sprite.Bounds()
			centerX = float64(bounds.Dx()) / 2
			centerY = float64(bounds.Dy()) / 2
		}

		op.GeoM.Translate(-centerX, -centerY)
		op.GeoM.Rotate(float64(*angle) * common.DegreesToRadians)
		op.GeoM.Translate(centerX, centerY)
	}
}

// RenderEntity handles rendering a single entity
func RenderEntity(entry *donburi.Entry, screen *ebiten.Image) {
	pos := Position.Get(entry)
	sprite := Sprite.Get(entry)

	if sprite == nil {
		return
	}

	op := &ebiten.DrawImageOptions{}
	applySpriteTransform(entry, *sprite, op)
	// Apply position translation
	op.GeoM.Translate(pos.X, pos.Y)
	screen.DrawImage(*sprite, op)
}

// RenderSystem handles entity rendering
func RenderSystem(w donburi.World, screen *ebiten.Image) {
	query.NewQuery(
		filter.And(
			filter.Contains(Position),
			filter.Contains(Sprite),
		),
	).Each(w, func(entry *donburi.Entry) {
		RenderEntity(entry, screen)
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
