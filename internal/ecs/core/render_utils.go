package core

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

// ApplySpriteTransform applies scaling and rotation to sprite
// This is a shared utility to eliminate code duplication
func ApplySpriteTransform(entry *donburi.Entry, sprite *ebiten.Image, op *ebiten.DrawImageOptions) {
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

	// Apply rotation if orbital component exists
	if entry.HasComponent(Orbital) {
		orb := Orbital.Get(entry)
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
