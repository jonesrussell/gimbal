package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/jonesrussell/gimbal/internal/ui/core"
)

// SpriteManager handles sprite operations and scaling
type SpriteManager struct {
	heartSprite *ebiten.Image
	ammoSprite  *ebiten.Image
}

// NewSpriteManager creates a new sprite manager
func NewSpriteManager(heartSprite, ammoSprite *ebiten.Image) *SpriteManager {
	return &SpriteManager{
		heartSprite: scaleSprite(heartSprite, core.HeartIconSize, core.HeartIconSize),
		ammoSprite:  scaleSprite(ammoSprite, core.AmmoIconSize, core.AmmoIconSize),
	}
}

// GetHeartSprite returns the scaled heart sprite
func (sm *SpriteManager) GetHeartSprite() *ebiten.Image {
	return sm.heartSprite
}

// GetAmmoSprite returns the scaled ammo sprite or creates a fallback
func (sm *SpriteManager) GetAmmoSprite() *ebiten.Image {
	if sm.ammoSprite != nil {
		return sm.ammoSprite
	}
	return sm.createFallbackAmmoSprite()
}

// scaleSprite scales a sprite to the specified dimensions
func scaleSprite(sprite *ebiten.Image, width, height int) *ebiten.Image {
	if sprite == nil {
		return nil
	}

	bounds := sprite.Bounds()
	if bounds.Dx() == width && bounds.Dy() == height {
		return sprite
	}

	scaled := ebiten.NewImage(width, height)
	op := &ebiten.DrawImageOptions{}
	scaleX := float64(width) / float64(bounds.Dx())
	scaleY := float64(height) / float64(bounds.Dy())
	op.GeoM.Scale(scaleX, scaleY)
	scaled.DrawImage(sprite, op)
	return scaled
}

// createFallbackAmmoSprite creates a simple colored square as fallback
func (sm *SpriteManager) createFallbackAmmoSprite() *ebiten.Image {
	sprite := ebiten.NewImage(core.AmmoIconSize, core.AmmoIconSize)
	sprite.Fill(color.NRGBA{R: 255, G: 255, B: 0, A: 255}) // Yellow square
	return sprite
}
