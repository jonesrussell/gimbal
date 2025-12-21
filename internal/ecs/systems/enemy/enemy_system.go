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
	"github.com/jonesrussell/gimbal/internal/ecs/managers"
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
	bossConfig     *managers.BossConfig // Current level's boss configuration

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

	// Check if we need to start a new wave (but not if we're waiting for level start or inter-wave delay)
	if es.waveManager.GetCurrentWave() == nil &&
		es.waveManager.HasMoreWaves() &&
		!es.waveManager.IsWaiting() &&
		!es.waveManager.IsWaitingForLevelStart() {
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

	// Handle boss spawning after all waves are complete
	es.handleBossSpawning(ctx, deltaTime)

	// Update enemy movement (including boss)
	es.updateEnemies(deltaTime)
	es.UpdateBossMovement(deltaTime)

	return nil
}

// handleWaveCompletion handles wave completion and advances to next wave
func (es *EnemySystem) handleWaveCompletion(ctx context.Context, deltaTime float64) {
	if !es.waveManager.IsWaveComplete() {
		return
	}

	es.waveManager.CompleteWave()
	if es.waveManager.HasMoreWaves() {
		es.waveManager.StartNextWave()
		return
	}

	// All waves complete - boss will be handled in handleBossSpawning
	if es.bossConfig != nil && es.bossConfig.Enabled && !es.bossSpawned {
		es.logger.Debug("All waves complete, boss will spawn soon", "spawn_delay", es.bossConfig.SpawnDelay)
	}
}

// handleBossSpawning handles boss spawn timer and spawning
func (es *EnemySystem) handleBossSpawning(ctx context.Context, deltaTime float64) {
	// Only spawn boss if all waves are complete and boss hasn't been spawned yet
	if es.waveManager.HasMoreWaves() {
		return // Still have waves to complete
	}

	if es.bossSpawned {
		return // Boss already spawned
	}

	if es.bossConfig == nil || !es.bossConfig.Enabled {
		return // No boss configured for this level
	}

	// Increment boss spawn timer
	es.bossSpawnTimer += deltaTime
	spawnDelay := es.bossConfig.SpawnDelay
	if spawnDelay <= 0 {
		spawnDelay = BossSpawnDelay // Fallback to default
	}

	if es.bossSpawnTimer >= spawnDelay {
		es.SpawnBoss(ctx)
		es.bossSpawned = true
		es.logger.Debug("Boss spawned", "delay", es.bossSpawnTimer)
	}
}

// spawnWaveEnemy spawns a single enemy from the current wave
func (es *EnemySystem) spawnWaveEnemy(ctx context.Context, wave *WaveState) {
	enemyType := es.waveManager.GetNextEnemyType()
	enemyData := GetEnemyTypeData(enemyType)

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

func (es *EnemySystem) updateEnemies(deltaTime float64) {
	query.NewQuery(
		filter.And(
			filter.Contains(core.EnemyTag),
			filter.Contains(core.Position),
			filter.Contains(core.Movement),
		),
	).Each(es.world, func(entry *donburi.Entry) {
		pos := core.Position.Get(entry)
		mov := core.Movement.Get(entry)

		// Update pattern time
		mov.PatternTime += deltaTime

		// Apply movement pattern
		velocity := es.applyMovementPattern(*mov)

		// Velocity is in pixels per frame (at 60fps), scale by deltaTime
		// deltaTime is typically 1/60 seconds, so multiply by 60 to get frame-equivalent
		frameScale := deltaTime * 60.0
		pos.X += velocity.X * frameScale
		pos.Y += velocity.Y * frameScale

		// Update movement component with new pattern time
		core.Movement.SetValue(entry, *mov)

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

// applyMovementPattern applies the movement pattern to calculate velocity
func (es *EnemySystem) applyMovementPattern(mov core.MovementData) common.Point {
	pattern := MovementPattern(mov.Pattern)

	switch pattern {
	case MovementPatternZigzag:
		return es.calculateZigzagVelocity(mov)
	case MovementPatternAccelerating:
		return es.calculateAcceleratingVelocity(mov)
	case MovementPatternPulsing:
		return es.calculatePulsingVelocity(mov)
	default:
		// Normal movement
		return mov.Velocity
	}
}

// calculateZigzagVelocity calculates zigzag movement (oscillates side-to-side)
func (es *EnemySystem) calculateZigzagVelocity(mov core.MovementData) common.Point {
	// Zigzag frequency (how fast it oscillates)
	zigzagFreq := 3.0      // oscillations per second
	zigzagAmplitude := 0.3 // how much it deviates

	// Calculate perpendicular angle (90 degrees to base direction)
	perpendicularAngle := mov.BaseAngle + stdmath.Pi/2

	// Oscillate perpendicular to movement direction
	oscillation := stdmath.Sin(mov.PatternTime*zigzagFreq*2*stdmath.Pi) * zigzagAmplitude

	// Base velocity
	baseVelX := stdmath.Cos(mov.BaseAngle) * mov.BaseSpeed
	baseVelY := stdmath.Sin(mov.BaseAngle) * mov.BaseSpeed

	// Add perpendicular oscillation
	perpendicularX := stdmath.Cos(perpendicularAngle) * oscillation * mov.BaseSpeed
	perpendicularY := stdmath.Sin(perpendicularAngle) * oscillation * mov.BaseSpeed

	return common.Point{
		X: baseVelX + perpendicularX,
		Y: baseVelY + perpendicularY,
	}
}

// calculateAcceleratingVelocity calculates accelerating movement (starts slow, speeds up)
func (es *EnemySystem) calculateAcceleratingVelocity(mov core.MovementData) common.Point {
	// Acceleration factor (0 to 1, where 1 is max speed)
	accelTime := 2.0 // seconds to reach max speed
	accelFactor := stdmath.Min(1.0, mov.PatternTime/accelTime)

	// Start at 30% speed, accelerate to 100%
	speedMultiplier := 0.3 + (accelFactor * 0.7)
	currentSpeed := mov.BaseSpeed * speedMultiplier

	return common.Point{
		X: stdmath.Cos(mov.BaseAngle) * currentSpeed,
		Y: stdmath.Sin(mov.BaseAngle) * currentSpeed,
	}
}

// calculatePulsingVelocity calculates pulsing movement (fast-slow-fast bursts)
func (es *EnemySystem) calculatePulsingVelocity(mov core.MovementData) common.Point {
	// Pulse frequency (how often it pulses)
	pulseFreq := 2.0 // pulses per second
	pulsePhase := mov.PatternTime * pulseFreq * 2 * stdmath.Pi

	// Use sine wave to create smooth pulsing (0.5 to 1.0 speed multiplier)
	speedMultiplier := 0.5 + 0.5*stdmath.Sin(pulsePhase)
	currentSpeed := mov.BaseSpeed * speedMultiplier

	return common.Point{
		X: stdmath.Cos(mov.BaseAngle) * currentSpeed,
		Y: stdmath.Sin(mov.BaseAngle) * currentSpeed,
	}
}

// DestroyEnemy destroys an enemy entity and returns points based on type
func (es *EnemySystem) DestroyEnemy(entity donburi.Entity) int {
	entry := es.world.Entry(entity)
	if !entry.Valid() {
		return 0
	}

	// Get enemy type from component (preferred) or fall back to health heuristic
	var enemyType EnemyType
	if entry.HasComponent(core.EnemyTypeID) {
		typeID := core.EnemyTypeID.Get(entry)
		enemyType = EnemyType(*typeID)
	} else if entry.HasComponent(core.Health) {
		// Fallback for legacy entities without EnemyTypeID
		health := core.Health.Get(entry)
		if health.Maximum >= 10 {
			enemyType = EnemyTypeBoss
		} else if health.Maximum >= 2 {
			enemyType = EnemyTypeHeavy
		} else {
			enemyType = EnemyTypeBasic
		}
	} else {
		enemyType = EnemyTypeBasic
	}

	points := GetEnemyTypeData(enemyType).Points

	// Mark enemy killed in wave manager
	es.waveManager.MarkEnemyKilled()

	// Remove the entity from the world
	es.world.Remove(entity)

	return points
}

// LoadLevelConfig loads the waves and boss configuration for a level
func (es *EnemySystem) LoadLevelConfig(waves []WaveConfig, bossConfig *managers.BossConfig) {
	es.waveManager.LoadWaves(waves)
	es.bossConfig = bossConfig
	es.bossSpawned = false
	es.bossSpawnTimer = 0

	bossHealth := 0
	if bossConfig != nil {
		bossHealth = bossConfig.Health
	}

	es.logger.Debug("Level config loaded",
		"waves", len(waves),
		"boss_enabled", bossConfig != nil && bossConfig.Enabled,
		"boss_health", bossHealth)
}

// Reset resets the enemy system for a new level
func (es *EnemySystem) Reset() {
	es.waveManager.Reset()
	es.bossSpawned = false
	es.bossSpawnTimer = 0
	es.spawnTimer = 0
	// Note: bossConfig is not reset here as it should be set via LoadLevelConfig
}

// IsBossActive checks if there's an active boss
func (es *EnemySystem) IsBossActive() bool {
	if es.bossSpawned {
		// Check if boss still exists using EnemyTypeID component
		count := 0
		query.NewQuery(
			filter.And(
				filter.Contains(core.EnemyTag),
				filter.Contains(core.EnemyTypeID),
			),
		).Each(es.world, func(entry *donburi.Entry) {
			typeID := core.EnemyTypeID.Get(entry)
			if EnemyType(*typeID) == EnemyTypeBoss {
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

// GetWaveManager returns the wave manager
func (es *EnemySystem) GetWaveManager() *WaveManager {
	return es.waveManager
}
