package core

import (
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/jonesrussell/gimbal/internal/common"
)

func RenderSystem(w donburi.World, screen *ebiten.Image) {
	count := 0
	query.NewQuery(
		filter.And(
			filter.Contains(Position),
			filter.Contains(Sprite),
		),
	).Each(w, func(entry *donburi.Entry) {
		RenderEntity(entry, screen)
		count++
	})
	log.Printf("[RenderSystem] Entities rendered: %d, screen bounds: %+v", count, screen.Bounds())
}

func RenderEntity(entry *donburi.Entry, screen *ebiten.Image) {
	pos := Position.Get(entry)
	sprite := Sprite.Get(entry)

	if sprite == nil {
		return
	}

	renderStaticEntity(entry, screen, pos, *sprite)
}

func renderStaticEntity(entry *donburi.Entry, screen *ebiten.Image, pos *common.Point, sprite *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	ApplySpriteTransform(entry, sprite, op)

		// Apply invincibility flashing if entity has health and is invincible
		if entry.HasComponent(Health) {
			health := Health.Get(entry)
			if health.IsInvincible {
				// Flash every 0.2 seconds (5 times per second)
				flashRate := 0.2 * float64(time.Second) // Convert to time.Duration
				remainingTime := health.InvincibilityDuration - health.InvincibilityTime
				flashPhase := int(remainingTime / time.Duration(flashRate))
			if flashPhase%2 == 0 {
				// Make sprite semi-transparent during flash
				op.ColorScale.SetR(1)
				op.ColorScale.SetG(1)
				op.ColorScale.SetB(1)
				op.ColorScale.SetA(0.5)
			}
		}
	}

	// Apply position translation
	op.GeoM.Translate(pos.X, pos.Y)
	screen.DrawImage(sprite, op)
}
