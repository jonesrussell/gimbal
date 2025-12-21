package enemy

import (
	"context"
	"image/color"
	stdmath "math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/config"
	"github.com/jonesrussell/gimbal/internal/ecs/core"
)

// spawnWaveEnemy spawns a single enemy from the current wave
func (es *EnemySystem) spawnWaveEnemy(ctx context.Context, wave *WaveState) {
	enemyType := es.waveManager.GetNextEnemyType()
	enemyData, err := es.GetEnemyTypeData(enemyType)
	if err != nil {
		es.logger.Error("Failed to get enemy type data", "type", enemyType, "error", err)
		return // Skip this enemy type
	}

	// Override movement pattern with wave's pattern
	enemyData.MovementPattern = wave.Config.MovementPattern

	// Get formation data
	centerX := float64(es.gameConfig.ScreenSize.Width) / 2
	centerY := float64(es.gameConfig.ScreenSize.Height) / 2
	spawnRadius := GetSpawnRadius(es.gameConfig)

	// Calculate base angle for formation (random rotation)
	//nolint:gosec // Game logic randomness is acceptable
	baseAngle := rand.Float64() * 2 * stdmath.Pi

	// Get formation positions
	formationParams := FormationParams{
		FormationType: wave.Config.FormationType,
		EnemyCount:    wave.Config.EnemyCount,
		CenterX:       centerX,
		CenterY:       centerY,
		BaseAngle:     baseAngle,
		SpawnRadius:   spawnRadius,
	}
	formationData := CalculateFormation(formationParams)

	// Spawn enemy at the appropriate position in formation
	enemyIndex := wave.EnemiesSpawned
	if enemyIndex < len(formationData) {
		formData := formationData[enemyIndex]
		es.spawnEnemyAt(ctx, formData.Position, formData.Angle, enemyType, &enemyData)
	} else {
		// Fallback: spawn at center with random angle
		es.spawnEnemyAt(ctx, common.Point{X: centerX, Y: centerY}, baseAngle, enemyType, &enemyData)
	}
}

// spawnEnemyAt spawns an enemy at a specific position with specific movement
func (es *EnemySystem) spawnEnemyAt(
	ctx context.Context,
	position common.Point,
	angle float64,
	enemyType EnemyType,
	enemyData *EnemyTypeData,
) {
	// Get or create sprite
	sprite := es.getEnemySprite(ctx, enemyType, enemyData)

	entity := es.world.Create(
		core.EnemyTag, core.Position, core.Sprite, core.Movement,
		core.Size, core.Health, core.EnemyTypeID,
	)
	entry := es.world.Entry(entity)

	// Set position
	core.Position.SetValue(entry, position)

	// Set sprite
	core.Sprite.SetValue(entry, sprite)

	// Set size
	core.Size.SetValue(entry, config.Size{Width: enemyData.Size, Height: enemyData.Size})

	// Set health
	core.Health.SetValue(entry, core.NewHealthData(enemyData.Health, enemyData.Health))

	// Set enemy type for proper identification (avoids health-based heuristics)
	core.EnemyTypeID.SetValue(entry, int(enemyType))

	// Set movement based on type and pattern
	switch enemyData.MovementType {
	case "spiral":
		// Spiral movement: start with angle, then spiral outward
		es.setSpiralMovement(entry, angle, enemyData.Speed, enemyData.MovementPattern)
	case "orbital":
		// Orbital movement (for boss, handled separately)
		// This shouldn't be called for regular enemies
		es.setOutwardMovement(entry, angle, enemyData.Speed, enemyData.MovementPattern)
	default:
		// Default: outward movement
		es.setOutwardMovement(entry, angle, enemyData.Speed, enemyData.MovementPattern)
	}

	es.logger.Debug("Enemy spawned",
		"type", enemyType.String(),
		"sprite", enemyData.SpriteName,
		"health", enemyData.Health,
		"position", position,
		"angle", angle)
}

// setOutwardMovement sets simple outward movement with optional pattern
func (es *EnemySystem) setOutwardMovement(entry *donburi.Entry, angle, speed float64, pattern MovementPattern) {
	velocity := common.Point{
		X: stdmath.Cos(angle) * speed,
		Y: stdmath.Sin(angle) * speed,
	}

	core.Movement.SetValue(entry, core.MovementData{
		Velocity:    velocity,
		MaxSpeed:    speed,
		Pattern:     int(pattern),
		PatternTime: 0,
		BaseAngle:   angle,
		BaseSpeed:   speed,
	})
}

// setSpiralMovement sets spiral movement pattern
func (es *EnemySystem) setSpiralMovement(entry *donburi.Entry, baseAngle, speed float64, pattern MovementPattern) {
	// For now, use outward movement with slight variation
	// TODO: Implement actual spiral pattern with time-based angle change
	velocity := common.Point{
		X: stdmath.Cos(baseAngle) * speed,
		Y: stdmath.Sin(baseAngle) * speed,
	}

	core.Movement.SetValue(entry, core.MovementData{
		Velocity:    velocity,
		MaxSpeed:    speed,
		Pattern:     int(pattern),
		PatternTime: 0,
		BaseAngle:   baseAngle,
		BaseSpeed:   speed,
	})
}

// getEnemySprite gets or creates the sprite for an enemy type
func (es *EnemySystem) getEnemySprite(
	ctx context.Context,
	enemyType EnemyType,
	enemyData *EnemyTypeData,
) *ebiten.Image {
	// Check cache
	if sprite, ok := es.enemySprites[enemyType]; ok {
		es.logger.Debug("[ENEMY_SPRITE] Using cached sprite",
			"type", enemyType.String(),
			"sprite_name", enemyData.SpriteName)
		return sprite
	}

	// Try to load sprite (full size, will be scaled during rendering)
	sprite, exists := es.resourceMgr.GetSprite(ctx, enemyData.SpriteName)
	if !exists {
		es.logger.Warn("Enemy sprite not found, using placeholder",
			"type", enemyType.String(),
			"sprite", enemyData.SpriteName)
		// Create placeholder with different colors
		sprite = ebiten.NewImage(enemyData.Size, enemyData.Size)
		switch enemyType {
		case EnemyTypeHeavy:
			sprite.Fill(color.RGBA{255, 165, 0, 255}) // Orange
		case EnemyTypeBoss:
			sprite.Fill(color.RGBA{128, 0, 128, 255}) // Purple
		default:
			sprite.Fill(color.RGBA{255, 0, 0, 255}) // Red
		}
	} else {
		es.logger.Debug("[ENEMY_SPRITE] Loaded sprite from resource manager",
			"type", enemyType.String(),
			"sprite_name", enemyData.SpriteName)
	}

	// Cache sprite
	es.enemySprites[enemyType] = sprite
	return sprite
}
