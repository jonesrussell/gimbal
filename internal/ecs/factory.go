package ecs

import (
	"math"
	"math/rand"
	"time"

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

	// Set up orbital movement - start at bottom (180 degrees)
	orbitalData := OrbitalData{
		Center:       center,
		Radius:       config.Radius,
		OrbitalAngle: common.HalfCircleDegrees, // 180 degrees
		FacingAngle:  0,                        // Will be calculated by PlayerInputSystem
	}
	Orbital.SetValue(entry, orbitalData)

	// Set size
	Size.SetValue(entry, config.PlayerSize)

	// Set initial angle
	Angle.SetValue(entry, common.Angle(0))

	return entity
}

// CreateStar creates a star entity with Gyruss-style movement
func CreateStar(w donburi.World, sprite *ebiten.Image, config *common.GameConfig, x, y float64) donburi.Entity {
	entity := w.Create(StarTag, Position, Sprite, Speed, Size, Scale)
	entry := w.Entry(entity)

	// Set position
	Position.SetValue(entry, common.Point{X: x, Y: y})
	Sprite.SetValue(entry, sprite)

	// Set speed
	Speed.SetValue(entry, config.StarSpeed)

	// Set size
	Size.SetValue(entry, common.Size{Width: int(config.StarSize), Height: int(config.StarSize)})

	// Set random initial scale (0.3 to 0.8)
	initialScale := 0.3 + float64(entry.Entity().Id()%6)*0.1
	Scale.SetValue(entry, initialScale)

	return entity
}

// CreateStarField creates multiple stars for the background in Gyruss-style pattern
func CreateStarField(w donburi.World, sprite *ebiten.Image, config *common.GameConfig) []donburi.Entity {
	entities := make([]donburi.Entity, config.NumStars)

	centerX := float64(config.ScreenSize.Width) / 2
	centerY := float64(config.ScreenSize.Height) / 2

	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < config.NumStars; i++ {
		// Create truly random positions along a small orbital path
		// Random angle around the circle (0 to 2π)
		angle := rand.Float64() * 2 * math.Pi

		// Random radius within the spawn range (30-80 pixels from center)
		spawnRadius := 30.0 + rand.Float64()*50.0

		x := centerX + math.Cos(angle)*spawnRadius
		y := centerY + math.Sin(angle)*spawnRadius

		entities[i] = CreateStar(w, sprite, config, x, y)
	}

	return entities
}
