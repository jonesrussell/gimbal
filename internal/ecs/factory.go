package ecs

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"

	"github.com/jonesrussell/gimbal/internal/common"
)

// CreatePlayer creates a player entity with orbital movement
func CreatePlayer(w donburi.World, sprite *ebiten.Image, config *common.GameConfig) donburi.Entity {
	entity := w.Create(PlayerTag, Position, Sprite, Orbital, Size, Angle)
	entry := w.Entry(entity)

	// Set initial position at the center of the screen
	center := common.Point{
		X: float64(config.ScreenSize.Width) / 2,
		Y: float64(config.ScreenSize.Height) / 2,
	}

	Position.SetValue(entry, center)
	Sprite.SetValue(entry, sprite)

	// Set up orbital movement
	orbitalData := OrbitalData{
		Center:       center,
		Radius:       config.Radius,
		OrbitalAngle: 0,
		FacingAngle:  0,
	}
	Orbital.SetValue(entry, orbitalData)

	// Set size
	Size.SetValue(entry, config.PlayerSize)

	// Set initial angle
	Angle.SetValue(entry, common.Angle(0))

	return entity
}

// CreateStar creates a star entity with falling movement
func CreateStar(w donburi.World, sprite *ebiten.Image, config *common.GameConfig, x, y float64) donburi.Entity {
	entity := w.Create(StarTag, Position, Sprite, Speed, Size)
	entry := w.Entry(entity)

	// Set position
	Position.SetValue(entry, common.Point{X: x, Y: y})
	Sprite.SetValue(entry, sprite)

	// Set speed
	Speed.SetValue(entry, config.StarSpeed)

	// Set size
	Size.SetValue(entry, common.Size{Width: int(config.StarSize), Height: int(config.StarSize)})

	return entity
}

// CreateStarField creates multiple stars for the background
func CreateStarField(w donburi.World, sprite *ebiten.Image, config *common.GameConfig) []donburi.Entity {
	entities := make([]donburi.Entity, config.NumStars)

	for i := 0; i < config.NumStars; i++ {
		// Random position across the screen width
		x := float64(i) * float64(config.ScreenSize.Width) / float64(config.NumStars)
		y := float64(i%10) * 50 // Staggered vertical positions

		entities[i] = CreateStar(w, sprite, config, x, y)
	}

	return entities
}
