package core

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

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

func RenderEntity(entry *donburi.Entry, screen *ebiten.Image) {
	pos := Position.Get(entry)
	sprite := Sprite.Get(entry)

	if sprite == nil {
		return
	}

	op := &ebiten.DrawImageOptions{}
	applySpriteTransform(entry, *sprite, op)

	// Apply invincibility flashing if entity has health and is invincible
	if entry.HasComponent(Health) {
		health := Health.Get(entry)
		if health.IsInvincible {
			// Flash every 0.2 seconds (5 times per second)
			flashRate := 0.2
			flashPhase := int((health.InvincibilityDuration - health.InvincibilityTime) / flashRate)
			if flashPhase%2 == 0 {
				// Make sprite semi-transparent during flash
				op.ColorM.Scale(1, 1, 1, 0.5)
			}
		}
	}

	// Apply position translation
	op.GeoM.Translate(pos.X, pos.Y)
	screen.DrawImage(*sprite, op)
}

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
		op.GeoM.Rotate(float64(orb.FacingAngle) * 0.017453292519943295) // DegreesToRadians
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
		op.GeoM.Rotate(float64(*angle) * 0.017453292519943295) // DegreesToRadians
		op.GeoM.Translate(centerX, centerY)
	}
}
