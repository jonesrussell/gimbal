package core

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/config"
	"github.com/jonesrussell/gimbal/internal/math"
)

// CreatePlayer creates a player entity with orbital movement
func CreatePlayer(w donburi.World, sprite *ebiten.Image, gameConfig *config.GameConfig) donburi.Entity {
	entity := w.Create(PlayerTag, Position, Sprite, Orbital, Size, Angle, Health)
	entry := w.Entry(entity)

	// Set initial position at the center of the screen
	center := common.Point{
		X: float64(gameConfig.ScreenSize.Width) / 2,
		Y: float64(gameConfig.ScreenSize.Height) / 2,
	}

	Position.SetValue(entry, center)
	Sprite.SetValue(entry, sprite)

	// Set up orbital movement - start at bottom (180 degrees)
	orbitalData := OrbitalData{
		Center:       center,
		Radius:       gameConfig.Radius,
		OrbitalAngle: 180, // 180 degrees
		FacingAngle:  0,   // Will be calculated by PlayerInputSystem
	}
	Orbital.SetValue(entry, orbitalData)

	// Set size
	Size.SetValue(entry, gameConfig.PlayerSize)

	// Set initial angle
	Angle.SetValue(entry, math.Angle(0))

	// Set initial health (3 lives)
	healthData := NewHealthData(3, 3)
	Health.SetValue(entry, healthData)

	return entity
}

// CreateStar creates a star entity with Gyruss-style movement
func CreateStar(w donburi.World, sprite *ebiten.Image, gameConfig *config.GameConfig, x, y float64) donburi.Entity {
	entity := w.Create(StarTag, Position, Sprite, Speed, Size, Scale)
	entry := w.Entry(entity)

	// Set position
	Position.SetValue(entry, common.Point{X: x, Y: y})
	Sprite.SetValue(entry, sprite)

	// Set speed
	Speed.SetValue(entry, gameConfig.StarSpeed)

	// Set size
	starSize := struct {
		Width, Height int
	}{Width: int(gameConfig.StarSize), Height: int(gameConfig.StarSize)}
	Size.SetValue(entry, starSize)

	// Set random initial scale (0.3 to 0.8)
	initialScale := 0.3 + float64(entry.Entity().Id()%6)*0.1
	Scale.SetValue(entry, initialScale)

	return entity
}

// CreateStarField creates multiple stars for the background in Gyruss-style pattern
func CreateStarField(w donburi.World, sprite *ebiten.Image, gameConfig *config.GameConfig) []donburi.Entity {
	entities := make([]donburi.Entity, gameConfig.NumStars)

	// Create star field helper with configuration from game config
	starConfig := &StarFieldConfig{
		SpawnRadiusMin: gameConfig.StarSpawnRadiusMin,
		SpawnRadiusMax: gameConfig.StarSpawnRadiusMax,
		Speed:          gameConfig.StarSpeed,
		MinScale:       gameConfig.StarMinScale,
		MaxScale:       gameConfig.StarMaxScale,
		ScaleDistance:  gameConfig.StarScaleDistance,
		ResetMargin:    gameConfig.StarResetMargin,
		Seed:           time.Now().UnixNano(),
	}
	starHelper := NewStarFieldHelper(starConfig, gameConfig.ScreenSize)

	for i := 0; i < gameConfig.NumStars; i++ {
		// Generate random position using helper with offset
		pos := starHelper.GenerateRandomPositionWithOffset(int64(i))

		// Create star at the generated position
		entities[i] = CreateStar(w, sprite, gameConfig, pos.X, pos.Y)
	}

	return entities
}
