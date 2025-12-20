package enemy

import (
	"context"
	"image/color"
	stdmath "math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/config"
	"github.com/jonesrussell/gimbal/internal/ecs/core"
	resources "github.com/jonesrussell/gimbal/internal/ecs/managers/resource"
)

// Enemy system constants
const (
	// DefaultSpawnIntervalSeconds is the time between enemy spawns (legacy, now used for wave delays)
	DefaultSpawnIntervalSeconds = 1.0
	// DefaultEnemySpeed is the movement speed of enemies
	DefaultEnemySpeed = 2.0
	// DefaultEnemySize is the size of enemy sprites
	DefaultEnemySize = 32
	// BossSpawnDelay is the delay after last wave before boss spawns
	BossSpawnDelay = 2.0
)

// EnemySystem manages enemy spawning, movement, and behavior
type EnemySystem struct {
	world         donburi.World
	gameConfig    *config.GameConfig
	spawnTimer    float64
	spawnInterval float64
	resourceMgr   *resources.ResourceManager
	logger        common.Logger

	// Wave management
	waveManager *WaveManager

	// Boss spawning
	bossSpawnTimer float64
	bossSpawned    bool

	// Enemy sprites cache
	enemySprites map[EnemyType]*ebiten.Image
}

// NewEnemySystem creates a new enemy management system with the provided dependencies
func NewEnemySystem(
	world donburi.World,
	gameConfig *config.GameConfig,
	resourceMgr *resources.ResourceManager,
	logger common.Logger,
) *EnemySystem {
	es := &EnemySystem{
		world:         world,
		gameConfig:    gameConfig,
		spawnTimer:    0,
		spawnInterval: DefaultSpawnIntervalSeconds,
		resourceMgr:   resourceMgr,
		logger:        logger,
		enemySprites:  make(map[EnemyType]*ebiten.Image),
	}

	// Initialize wave manager
	es.waveManager = NewWaveManager(world, logger)

	return es
}

func (es *EnemySystem) Update(ctx context.Context, deltaTime float64) error {
	// Check for cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Update wave manager
	es.waveManager.Update(deltaTime)

	// Check if we need to start a new wave
	if es.waveManager.GetCurrentWave() == nil && es.waveManager.HasMoreWaves() {
		es.waveManager.StartNextWave()
	}

	// Spawn enemies from current wave
	if es.waveManager.ShouldSpawnEnemy(deltaTime) {
		wave := es.waveManager.GetCurrentWave()
		if wave != nil {
			es.spawnWaveEnemy(ctx, wave)
			es.waveManager.MarkEnemySpawned()
		}
	}

	// Check if wave is complete and start next
	es.handleWaveCompletion(ctx, deltaTime)

	// Update enemy movement (including boss)
	es.updateEnemies()
	es.UpdateBossMovement(deltaTime)

	return nil
}

// handleWaveCompletion handles wave completion and boss spawning
func (es *EnemySystem) handleWaveCompletion(ctx context.Context, deltaTime float64) {
	if !es.waveManager.IsWaveComplete() {
		return
	}

	es.waveManager.CompleteWave()
	if es.waveManager.HasMoreWaves() {
		es.waveManager.StartNextWave()
		return
	}

	// All waves complete, start boss spawn timer
	if !es.bossSpawned {
		es.bossSpawnTimer += deltaTime
		if es.bossSpawnTimer >= BossSpawnDelay {
			es.SpawnBoss(ctx)
			es.bossSpawned = true
		}
	}
}

// spawnWaveEnemy spawns a single enemy from the current wave
func (es *EnemySystem) spawnWaveEnemy(ctx context.Context, wave *WaveState) {
	enemyType := es.waveManager.GetNextEnemyType()
	enemyData := GetEnemyTypeData(enemyType)

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
		es.spawnEnemyAt(ctx, formData.Position, formData.Angle, enemyType, enemyData)
	} else {
		// Fallback: spawn at center with random angle
		es.spawnEnemyAt(ctx, common.Point{X: centerX, Y: centerY}, baseAngle, enemyType, enemyData)
	}
}

// spawnEnemyAt spawns an enemy at a specific position with specific movement
func (es *EnemySystem) spawnEnemyAt(
	ctx context.Context,
	position common.Point,
	angle float64,
	enemyType EnemyType,
	enemyData EnemyTypeData,
) {
	// Get or create sprite
	sprite := es.getEnemySprite(ctx, enemyType, enemyData)

	entity := es.world.Create(
		core.EnemyTag, core.Position, core.Sprite, core.Movement,
		core.Size, core.Health,
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

	// Set movement based on type
	switch enemyData.MovementType {
	case "spiral":
		// Spiral movement: start with angle, then spiral outward
		es.setSpiralMovement(entry, angle, enemyData.Speed)
	case "orbital":
		// Orbital movement (for boss, handled separately)
		// This shouldn't be called for regular enemies
		es.setOutwardMovement(entry, angle, enemyData.Speed)
	default:
		// Default: outward movement
		es.setOutwardMovement(entry, angle, enemyData.Speed)
	}

	es.logger.Debug("Enemy spawned", "type", enemyType, "position", position, "angle", angle)
}

// setOutwardMovement sets simple outward movement
func (es *EnemySystem) setOutwardMovement(entry *donburi.Entry, angle, speed float64) {
	velocity := common.Point{
		X: stdmath.Cos(angle) * speed,
		Y: stdmath.Sin(angle) * speed,
	}

	core.Movement.SetValue(entry, core.MovementData{
		Velocity: velocity,
		MaxSpeed: speed,
	})
}

// setSpiralMovement sets spiral movement pattern
func (es *EnemySystem) setSpiralMovement(entry *donburi.Entry, baseAngle, speed float64) {
	// For now, use outward movement with slight variation
	// TODO: Implement actual spiral pattern with time-based angle change
	velocity := common.Point{
		X: stdmath.Cos(baseAngle) * speed,
		Y: stdmath.Sin(baseAngle) * speed,
	}

	core.Movement.SetValue(entry, core.MovementData{
		Velocity: velocity,
		MaxSpeed: speed,
	})
}

// getEnemySprite gets or creates the sprite for an enemy type
func (es *EnemySystem) getEnemySprite(ctx context.Context, enemyType EnemyType, enemyData EnemyTypeData) *ebiten.Image {
	// Check cache
	if sprite, ok := es.enemySprites[enemyType]; ok {
		return sprite
	}

	// Try to load sprite
	sprite, exists := es.resourceMgr.GetSprite(ctx, enemyData.SpriteName)
	if !exists {
		es.logger.Warn("Enemy sprite not found, using placeholder", "type", enemyType, "sprite", enemyData.SpriteName)
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
	}

	// Cache sprite
	es.enemySprites[enemyType] = sprite
	return sprite
}

func (es *EnemySystem) updateEnemies() {
	query.NewQuery(
		filter.And(
			filter.Contains(core.EnemyTag),
			filter.Contains(core.Position),
			filter.Contains(core.Movement),
		),
	).Each(es.world, func(entry *donburi.Entry) {
		pos := core.Position.Get(entry)
		mov := core.Movement.Get(entry)
		pos.X += mov.Velocity.X
		pos.Y += mov.Velocity.Y

		// Remove enemies when they move too far from center (Gyruss-style)
		centerX := float64(es.gameConfig.ScreenSize.Width) / 2
		centerY := float64(es.gameConfig.ScreenSize.Height) / 2
		distanceFromCenter := stdmath.Sqrt((pos.X-centerX)*(pos.X-centerX) + (pos.Y-centerY)*(pos.Y-centerY))
		screenWidth := float64(es.gameConfig.ScreenSize.Width)
		screenHeight := float64(es.gameConfig.ScreenSize.Height)
		maxDistance := stdmath.Max(screenWidth, screenHeight) * 0.8

		if distanceFromCenter > maxDistance {
			es.world.Remove(entry.Entity())
		}
	})
}

// DestroyEnemy destroys an enemy entity and returns points based on type
func (es *EnemySystem) DestroyEnemy(entity donburi.Entity) int {
	entry := es.world.Entry(entity)
	if !entry.Valid() {
		return 0
	}

	// Determine enemy type from health (heuristic)
	health := core.Health.Get(entry)
	points := 100 // Default

	if health != nil {
		if health.Maximum >= 10 {
			points = GetEnemyTypeData(EnemyTypeBoss).Points
		} else if health.Maximum >= 2 {
			points = GetEnemyTypeData(EnemyTypeHeavy).Points
		} else {
			points = GetEnemyTypeData(EnemyTypeBasic).Points
		}
	}

	// Mark enemy killed in wave manager
	es.waveManager.MarkEnemyKilled()

	// Remove the entity from the world
	es.world.Remove(entity)

	return points
}

// Reset resets the enemy system for a new level
func (es *EnemySystem) Reset() {
	es.waveManager.Reset()
	es.bossSpawned = false
	es.bossSpawnTimer = 0
	es.spawnTimer = 0
}

// IsBossActive checks if there's an active boss
func (es *EnemySystem) IsBossActive() bool {
	if es.bossSpawned {
		// Check if boss still exists
		count := 0
		query.NewQuery(
			filter.And(
				filter.Contains(core.EnemyTag),
				filter.Contains(core.Orbital),
			),
		).Each(es.world, func(entry *donburi.Entry) {
			health := core.Health.Get(entry)
			if health != nil && health.Maximum >= 10 {
				count++
			}
		})
		return count > 0
	}
	return false
}

// WasBossSpawned returns true if boss was spawned (even if now killed)
func (es *EnemySystem) WasBossSpawned() bool {
	return es.bossSpawned
}
