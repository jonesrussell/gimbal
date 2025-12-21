package debug

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/jonesrussell/gimbal/internal/ecs/core"
)

// drawEntityDebug draws entity positions, bounds, and states (only for nearby entities)
func (dr *DebugRenderer) drawEntityDebug(screen *ebiten.Image, world donburi.World) {
	// Draw all entities with position and sprite components
	query.NewQuery(
		filter.And(
			filter.Contains(core.Position),
			filter.Contains(core.Sprite),
		),
	).Each(world, func(entry *donburi.Entry) {
		pos := core.Position.Get(entry)
		sprite := core.Sprite.Get(entry)

		if pos == nil || sprite == nil {
			return
		}

		// Determine entity type and color
		entityColor := dr.getEntityColor(entry)

		// Draw entity center point - tiny colored dot
		vector.DrawFilledCircle(screen, float32(pos.X), float32(pos.Y), 1, entityColor, false)

		// Draw bounding box if size component exists
		if entry.HasComponent(core.Size) {
			size := core.Size.Get(entry)
			if size != nil {
				// Calculate bounds
				boundsX := pos.X - float64(size.Width)/2
				boundsY := pos.Y - float64(size.Height)/2

				// Draw bounding box - very thin colored outline
				vector.StrokeRect(screen, float32(boundsX), float32(boundsY), float32(size.Width), float32(size.Height),
					1, entityColor, false)

				// Only show entity info text if mouse is nearby
				if dr.shouldShowEntityInfo(pos) {
					entityInfo := fmt.Sprintf("Pos: (%.1f,%.1f)\nSize: %dx%d", pos.X, pos.Y, size.Width, size.Height)
					dr.drawTextWithBackground(screen, entityInfo, pos.X+10, pos.Y-20)
				}
			}
		}

		// Draw sprite bounds (only if mouse is nearby)
		if dr.shouldShowEntityInfo(pos) {
			dr.drawSpriteDebug(screen, *sprite, pos.X, pos.Y)
		}
	})
}

// getEntityColor returns the appropriate color for different entity types
func (dr *DebugRenderer) getEntityColor(entry *donburi.Entry) color.RGBA {
	// Player entities - Green
	if entry.HasComponent(core.PlayerTag) {
		return color.RGBA{0, 255, 0, 80}
	}

	// Enemy entities - Red
	if entry.HasComponent(core.EnemyTag) {
		return color.RGBA{255, 0, 0, 80}
	}

	// Projectile entities - Yellow
	if entry.HasComponent(core.ProjectileTag) {
		return color.RGBA{255, 255, 0, 60}
	}

	// Star entities - Blue
	if entry.HasComponent(core.StarTag) {
		return color.RGBA{0, 150, 255, 40}
	}

	// Default - White
	return color.RGBA{255, 255, 255, 50}
}

// drawSpriteDebug draws sprite boundaries and center points
func (dr *DebugRenderer) drawSpriteDebug(screen, sprite *ebiten.Image, x, y float64) {
	bounds := sprite.Bounds()

	// Calculate sprite position (assuming sprite is centered on entity)
	spriteX := x - float64(bounds.Dx())/2
	spriteY := y - float64(bounds.Dy())/2

	// Draw sprite boundary rectangle - very thin white outline
	vector.StrokeRect(screen, float32(spriteX), float32(spriteY), float32(bounds.Dx()), float32(bounds.Dy()),
		1, color.RGBA{255, 255, 255, 30}, false)

	// Draw sprite center point - tiny white dot
	centerX, centerY := x, y
	vector.DrawFilledCircle(screen, float32(centerX), float32(centerY), 1, color.RGBA{255, 255, 255, 80}, false)
}

// drawCollisionDebug draws collision boundaries and detection ranges
func (dr *DebugRenderer) drawCollisionDebug(screen *ebiten.Image, world donburi.World) {
	// Draw player collision area
	query.NewQuery(
		filter.And(
			filter.Contains(core.PlayerTag),
			filter.Contains(core.Position),
			filter.Contains(core.Size),
		),
	).Each(world, func(entry *donburi.Entry) {
		pos := core.Position.Get(entry)
		size := core.Size.Get(entry)

		if pos == nil || size == nil {
			return
		}

		// Only show collision debug if mouse is nearby
		if !dr.shouldShowEntityInfo(pos) {
			return
		}

		// Draw player collision box - very thin green outline
		boundsX := pos.X - float64(size.Width)/2
		boundsY := pos.Y - float64(size.Height)/2
		vector.StrokeRect(screen, float32(boundsX), float32(boundsY), float32(size.Width), float32(size.Height),
			1, color.RGBA{0, 255, 0, 60}, false)
	})

	// Draw enemy collision areas
	query.NewQuery(
		filter.And(
			filter.Contains(core.EnemyTag),
			filter.Contains(core.Position),
			filter.Contains(core.Size),
		),
	).Each(world, func(entry *donburi.Entry) {
		pos := core.Position.Get(entry)
		size := core.Size.Get(entry)

		if pos == nil || size == nil {
			return
		}

		// Only show collision debug if mouse is nearby
		if !dr.shouldShowEntityInfo(pos) {
			return
		}

		// Draw enemy collision box - very thin red outline
		boundsX := pos.X - float64(size.Width)/2
		boundsY := pos.Y - float64(size.Height)/2
		vector.StrokeRect(screen, float32(boundsX), float32(boundsY), float32(size.Width), float32(size.Height),
			1, color.RGBA{255, 0, 0, 60}, false)
	})
}
