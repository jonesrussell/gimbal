package enemy

import (
	"context"
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/config"
	"github.com/jonesrussell/gimbal/internal/ecs/contracts"
	"github.com/jonesrussell/gimbal/internal/ecs/core"
	resources "github.com/jonesrussell/gimbal/internal/ecs/managers/resource"
)

// EnemySystem manages enemy spawning, movement, and behavior
type EnemySystem struct {
	world         donburi.World
	gameConfig    *config.GameConfig
	spawnTimer    float64
	spawnInterval float64
	resourceMgr   *resources.ResourceManager
	logger        common.Logger

	// Enemy sprites cache
	enemySprites map[core.EnemyType]*ebiten.Image

	// Wave management
	currentWave    int
	enemiesInWave  int
	enemiesSpawned int
	waveTimer      float64
	waveInterval   float64

	// Difficulty scaling
	difficultyLevel     int
	spawnRateMultiplier float64
}

// NewEnemySystem creates a new enemy management system with the provided dependencies
func NewEnemySystem(
	world donburi.World,
	gameConfig *config.GameConfig,
	resourceMgr *resources.ResourceManager,
	logger common.Logger,
) *EnemySystem {
	es := &EnemySystem{
		world:               world,
		gameConfig:          gameConfig,
		spawnTimer:          0,
		spawnInterval:       60, // Spawn every 60 frames (1 second at 60fps)
		resourceMgr:         resourceMgr,
		logger:              logger,
		enemySprites:        make(map[core.EnemyType]*ebiten.Image),
		currentWave:         1,
		enemiesInWave:       5,
		waveTimer:           0,
		waveInterval:        300, // 5 seconds between waves
		difficultyLevel:     1,
		spawnRateMultiplier: 1.0,
	}

	return es
}

func (es *EnemySystem) Update(ctx context.Context, deltaTime float64) error {
	// Check for cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	es.spawnTimer += deltaTime
	es.waveTimer += deltaTime

	// Check if it's time for a new wave
	if es.waveTimer >= es.waveInterval {
		es.startNewWave()
	}

	// Spawn enemies based on current wave and difficulty
	if es.spawnTimer >= es.spawnInterval && es.enemiesSpawned < es.enemiesInWave {
		es.spawnEnemy(ctx)
		es.spawnTimer = 0
		es.enemiesSpawned++
	}

	es.updateEnemies(deltaTime)
	return nil
}

func (es *EnemySystem) startNewWave() {
	es.currentWave++
	es.enemiesSpawned = 0
	es.waveTimer = 0

	// Increase difficulty
	es.difficultyLevel = (es.currentWave-1)/5 + 1
	es.spawnRateMultiplier = 1.0 + float64(es.difficultyLevel-1)*0.2
	es.enemiesInWave = 5 + es.currentWave*2

	// Adjust spawn interval based on difficulty
	es.spawnInterval = 60 / es.spawnRateMultiplier

	es.logger.Info("Starting new wave",
		"wave", es.currentWave,
		"enemies", es.enemiesInWave,
		"difficulty", es.difficultyLevel)
}

func (es *EnemySystem) spawnEnemy(ctx context.Context) {
	// Determine enemy type based on wave and difficulty
	enemyType := es.selectEnemyType()

	// Load enemy sprite if not already loaded
	if es.enemySprites[enemyType] == nil {
		es.loadEnemySprite(ctx, enemyType)
	}

	es.logger.Debug("[ENEMY_SPAWN] Spawning enemy", "type", enemyType)

	// Get enemy configuration
	enemyConfig := es.getEnemyConfig(enemyType)

	// Calculate spawn position (Gyruss-style: from center outward)
	spawnPos := es.calculateSpawnPosition(enemyType)

	// Create enemy entity
	entity := es.world.Create(
		core.EnemyTag, core.Position, core.Sprite, core.Movement,
		core.Size, core.Health, core.EnemyTypeComponent, core.EnemyBehavior,
	)
	entry := es.world.Entry(entity)

	// Set basic components
	core.Position.SetValue(entry, spawnPos)
	core.Sprite.SetValue(entry, es.enemySprites[enemyType])
	core.Size.SetValue(entry, enemyConfig.Size)
	core.Health.SetValue(entry, core.NewHealthData(enemyConfig.Health, enemyConfig.Health))

	// Set enemy type data
	core.EnemyTypeComponent.SetValue(entry, enemyConfig)

	// Set behavior data
	behavior := es.createEnemyBehavior(enemyType, spawnPos)
	core.EnemyBehavior.SetValue(entry, behavior)

	// Set movement
	movement := es.calculateInitialMovement(enemyType, spawnPos, behavior)
	core.Movement.SetValue(entry, movement)
}

func (es *EnemySystem) selectEnemyType() core.EnemyType {
	// Weighted random selection based on wave and difficulty
	weights := es.getEnemyTypeWeights()

	totalWeight := 0.0
	for _, weight := range weights {
		totalWeight += weight
	}

	roll := rand.Float64() * totalWeight
	currentWeight := 0.0

	enemyTypes := []core.EnemyType{
		core.EnemyTypeSwarmDrone,
		core.EnemyTypeHeavyCruiser,
		core.EnemyTypeAsteroid,
		core.EnemyTypeZigzag,
		core.EnemyTypeShooter,
		core.EnemyTypeBoss,
	}

	for i, enemyType := range enemyTypes {
		currentWeight += weights[i]
		if roll <= currentWeight {
			return enemyType
		}
	}

	return core.EnemyTypeSwarmDrone // Default fallback
}

func (es *EnemySystem) getEnemyTypeWeights() []float64 {
	// Base weights for different enemy types
	baseWeights := []float64{
		0.4,  // SwarmDrone - 40% chance
		0.25, // HeavyCruiser - 25% chance
		0.2,  // Asteroid - 20% chance
		0.1,  // Zigzag - 10% chance
		0.04, // Shooter - 4% chance
		0.01, // Boss - 1% chance
	}

	// Adjust weights based on wave and difficulty
	if es.currentWave >= 5 {
		baseWeights[1] += 0.1  // More HeavyCruisers
		baseWeights[4] += 0.05 // More Shooters
	}

	if es.currentWave >= 10 {
		baseWeights[5] += 0.02 // More Bosses
	}

	return baseWeights
}

func (es *EnemySystem) getEnemyConfig(enemyType core.EnemyType) core.EnemyTypeData {
	configs := map[core.EnemyType]core.EnemyTypeData{
		core.EnemyTypeSwarmDrone: {
			Type:       core.EnemyTypeSwarmDrone,
			Size:       config.Size{Width: 16, Height: 16},
			Health:     1,
			Speed:      3.0,
			Points:     100,
			SpriteName: "enemy_swarm",
			Color:      contracts.Color{R: 255, G: 0, B: 0, A: 255}, // Red
		},
		core.EnemyTypeHeavyCruiser: {
			Type:       core.EnemyTypeHeavyCruiser,
			Size:       config.Size{Width: 32, Height: 32},
			Health:     3,
			Speed:      1.5,
			Points:     300,
			SpriteName: "enemy_cruiser",
			Color:      contracts.Color{R: 255, G: 165, B: 0, A: 255}, // Orange
		},
		core.EnemyTypeBoss: {
			Type:       core.EnemyTypeBoss,
			Size:       config.Size{Width: 64, Height: 64},
			Health:     10,
			Speed:      0.8,
			Points:     1000,
			SpriteName: "enemy_boss",
			Color:      contracts.Color{R: 128, G: 0, B: 128, A: 255}, // Purple
		},
		core.EnemyTypeAsteroid: {
			Type:       core.EnemyTypeAsteroid,
			Size:       config.Size{Width: 24, Height: 24},
			Health:     2,
			Speed:      2.0,
			Points:     200,
			SpriteName: "enemy_asteroid",
			Color:      contracts.Color{R: 139, G: 69, B: 19, A: 255}, // Brown
		},
		core.EnemyTypeShooter: {
			Type:       core.EnemyTypeShooter,
			Size:       config.Size{Width: 28, Height: 28},
			Health:     2,
			Speed:      1.8,
			Points:     250,
			SpriteName: "enemy_shooter",
			Color:      contracts.Color{R: 0, G: 255, B: 0, A: 255}, // Green
		},
		core.EnemyTypeZigzag: {
			Type:       core.EnemyTypeZigzag,
			Size:       config.Size{Width: 20, Height: 20},
			Health:     1,
			Speed:      2.5,
			Points:     150,
			SpriteName: "enemy_zigzag",
			Color:      contracts.Color{R: 255, G: 255, B: 0, A: 255}, // Yellow
		},
	}

	return configs[enemyType]
}

func (es *EnemySystem) calculateSpawnPosition(enemyType core.EnemyType) common.Point {
	centerX := float64(es.gameConfig.ScreenSize.Width) / 2
	centerY := float64(es.gameConfig.ScreenSize.Height) / 2

	// Different spawn patterns based on enemy type
	switch enemyType {
	case core.EnemyTypeBoss:
		// Boss spawns at center
		return common.Point{X: centerX, Y: centerY}
	case core.EnemyTypeSwarmDrone, core.EnemyTypeZigzag:
		// Small enemies spawn in a tight circle around center
		angle := rand.Float64() * 2 * math.Pi
		radius := 50 + rand.Float64()*50
		return common.Point{
			X: centerX + math.Cos(angle)*radius,
			Y: centerY + math.Sin(angle)*radius,
		}
	default:
		// Larger enemies spawn further out
		angle := rand.Float64() * 2 * math.Pi
		radius := 100 + rand.Float64()*100
		return common.Point{
			X: centerX + math.Cos(angle)*radius,
			Y: centerY + math.Sin(angle)*radius,
		}
	}
}

func (es *EnemySystem) createEnemyBehavior(enemyType core.EnemyType, spawnPos common.Point) core.EnemyBehaviorData {
	behavior := core.EnemyBehaviorData{
		Phase:       0,
		PhaseSpeed:  1.0,
		Wobble:      0,
		WobbleSpeed: 2.0,
		FireTimer:   0,
		FireRate:    2.0, // Fire every 2 seconds
	}

	switch enemyType {
	case core.EnemyTypeSwarmDrone:
		behavior.Pattern = core.MovementPatternLinear
		behavior.PhaseSpeed = 1.5
	case core.EnemyTypeHeavyCruiser:
		behavior.Pattern = core.MovementPatternSpiral
		behavior.PhaseSpeed = 0.8
	case core.EnemyTypeBoss:
		behavior.Pattern = core.MovementPatternOrbital
		behavior.PhaseSpeed = 0.5
	case core.EnemyTypeAsteroid:
		behavior.Pattern = core.MovementPatternLinear
		behavior.PhaseSpeed = 1.2
	case core.EnemyTypeShooter:
		behavior.Pattern = core.MovementPatternChase
		behavior.PhaseSpeed = 1.0
	case core.EnemyTypeZigzag:
		behavior.Pattern = core.MovementPatternZigzag
		behavior.PhaseSpeed = 2.0
		behavior.Wobble = 30.0
		behavior.WobbleSpeed = 3.0
	}

	return behavior
}

func (es *EnemySystem) calculateInitialMovement(enemyType core.EnemyType, spawnPos common.Point, behavior core.EnemyBehaviorData) core.MovementData {
	centerX := float64(es.gameConfig.ScreenSize.Width) / 2
	centerY := float64(es.gameConfig.ScreenSize.Height) / 2

	config := es.getEnemyConfig(enemyType)
	speed := config.Speed

	// Calculate direction based on enemy type and behavior
	var velocity common.Point

	switch behavior.Pattern {
	case core.MovementPatternLinear:
		// Move outward from center
		angle := math.Atan2(spawnPos.Y-centerY, spawnPos.X-centerX)
		velocity = common.Point{
			X: math.Cos(angle) * speed,
			Y: math.Sin(angle) * speed,
		}
	case core.MovementPatternSpiral:
		// Spiral outward
		angle := math.Atan2(spawnPos.Y-centerY, spawnPos.X-centerX)
		velocity = common.Point{
			X: math.Cos(angle) * speed * 0.7,
			Y: math.Sin(angle) * speed * 0.7,
		}
	case core.MovementPatternZigzag:
		// Zigzag movement
		angle := math.Atan2(spawnPos.Y-centerY, spawnPos.X-centerX)
		velocity = common.Point{
			X: math.Cos(angle) * speed,
			Y: math.Sin(angle) * speed,
		}
	case core.MovementPatternOrbital:
		// Orbital movement around center
		angle := math.Atan2(spawnPos.Y-centerY, spawnPos.X-centerX) + math.Pi/2
		velocity = common.Point{
			X: math.Cos(angle) * speed,
			Y: math.Sin(angle) * speed,
		}
	case core.MovementPatternChase:
		// Chase player (simplified - move toward center)
		velocity = common.Point{
			X: (centerX - spawnPos.X) * speed * 0.01,
			Y: (centerY - spawnPos.Y) * speed * 0.01,
		}
	default:
		// Default linear movement
		angle := math.Atan2(spawnPos.Y-centerY, spawnPos.X-centerX)
		velocity = common.Point{
			X: math.Cos(angle) * speed,
			Y: math.Sin(angle) * speed,
		}
	}

	return core.MovementData{
		Velocity: velocity,
		MaxSpeed: speed,
	}
}

func (es *EnemySystem) loadEnemySprite(ctx context.Context, enemyType core.EnemyType) {
	config := es.getEnemyConfig(enemyType)

	// Try to load sprite from resource manager
	sprite, exists := es.resourceMgr.GetSprite(ctx, config.SpriteName)
	if !exists {
		es.logger.Warn("[ENEMY_SPAWN] Enemy sprite not found, creating placeholder", "type", enemyType)
		// Create a placeholder sprite with enemy color
		sprite = ebiten.NewImage(config.Size.Width, config.Size.Height)
		sprite.Fill(color.RGBA{config.Color.R, config.Color.G, config.Color.B, config.Color.A})
	} else {
		es.logger.Debug("[ENEMY_SPAWN] Enemy sprite loaded successfully", "type", enemyType, "bounds", sprite.Bounds())
	}

	es.enemySprites[enemyType] = sprite
}

func (es *EnemySystem) updateEnemies(deltaTime float64) {
	query.NewQuery(
		filter.And(
			filter.Contains(core.EnemyTag),
			filter.Contains(core.Position),
			filter.Contains(core.Movement),
			filter.Contains(core.EnemyBehavior),
		),
	).Each(es.world, func(entry *donburi.Entry) {
		pos := core.Position.Get(entry)
		mov := core.Movement.Get(entry)
		behavior := core.EnemyBehavior.Get(entry)
		enemyType := core.EnemyTypeComponent.Get(entry)

		// Update behavior phase
		behavior.Phase += behavior.PhaseSpeed * deltaTime

		// Apply movement pattern
		es.applyMovementPattern(pos, mov, behavior, enemyType, deltaTime)

		// Update fire timer for shooting enemies
		if enemyType.Type == core.EnemyTypeShooter {
			behavior.FireTimer += deltaTime
			if behavior.FireTimer >= behavior.FireRate {
				behavior.FireTimer = 0
				// TODO: Implement enemy shooting
			}
		}

		// Remove enemies when they move too far from center (Gyruss-style)
		centerX := float64(es.gameConfig.ScreenSize.Width) / 2
		centerY := float64(es.gameConfig.ScreenSize.Height) / 2
		distanceFromCenter := math.Sqrt((pos.X-centerX)*(pos.X-centerX) + (pos.Y-centerY)*(pos.Y-centerY))
		maxDistance := math.Max(float64(es.gameConfig.ScreenSize.Width), float64(es.gameConfig.ScreenSize.Height)) * 0.8

		if distanceFromCenter > maxDistance {
			es.world.Remove(entry.Entity())
		}
	})
}

func (es *EnemySystem) applyMovementPattern(pos *common.Point, mov *core.MovementData, behavior *core.EnemyBehaviorData, enemyType *core.EnemyTypeData, deltaTime float64) {
	centerX := float64(es.gameConfig.ScreenSize.Width) / 2
	centerY := float64(es.gameConfig.ScreenSize.Height) / 2

	switch behavior.Pattern {
	case core.MovementPatternLinear:
		// Simple linear movement (already set in initial movement)
		pos.X += mov.Velocity.X
		pos.Y += mov.Velocity.Y

	case core.MovementPatternSpiral:
		// Spiral outward movement
		angle := math.Atan2(pos.Y-centerY, pos.X-centerX)
		radius := math.Sqrt((pos.X-centerX)*(pos.X-centerX) + (pos.Y-centerY)*(pos.Y-centerY))

		// Increase radius over time
		radius += enemyType.Speed * deltaTime

		// Add spiral rotation
		angle += behavior.Phase * 0.5

		pos.X = centerX + math.Cos(angle)*radius
		pos.Y = centerY + math.Sin(angle)*radius

	case core.MovementPatternZigzag:
		// Zigzag movement
		baseAngle := math.Atan2(pos.Y-centerY, pos.X-centerX)
		wobbleAngle := math.Sin(behavior.Phase*behavior.WobbleSpeed) * behavior.Wobble * math.Pi / 180

		angle := baseAngle + wobbleAngle
		pos.X += math.Cos(angle) * enemyType.Speed
		pos.Y += math.Sin(angle) * enemyType.Speed

	case core.MovementPatternOrbital:
		// Orbital movement around center
		radius := math.Sqrt((pos.X-centerX)*(pos.X-centerX) + (pos.Y-centerY)*(pos.Y-centerY))
		angle := math.Atan2(pos.Y-centerY, pos.X-centerX)

		// Orbit around center
		angle += behavior.Phase * 0.3

		pos.X = centerX + math.Cos(angle)*radius
		pos.Y = centerY + math.Sin(angle)*radius

	case core.MovementPatternChase:
		// Chase player (simplified - move toward center)
		dirX := centerX - pos.X
		dirY := centerY - pos.Y
		distance := math.Sqrt(dirX*dirX + dirY*dirY)

		if distance > 0 {
			dirX /= distance
			dirY /= distance
		}

		pos.X += dirX * enemyType.Speed
		pos.Y += dirY * enemyType.Speed

	default:
		// Default linear movement
		pos.X += mov.Velocity.X
		pos.Y += mov.Velocity.Y
	}
}

// DestroyEnemy destroys an enemy entity and returns points
func (es *EnemySystem) DestroyEnemy(entity donburi.Entity) int {
	entry := es.world.Entry(entity)
	if !entry.Valid() {
		return 0
	}

	// Get enemy type for points
	enemyType := core.EnemyTypeComponent.Get(entry)
	points := enemyType.Points

	// Remove the entity from the world
	es.world.Remove(entity)

	es.logger.Debug("Enemy destroyed", "type", enemyType.Type, "points", points)
	return points
}

// GetActiveCount returns the number of active enemies
func (es *EnemySystem) GetActiveCount() int {
	count := 0
	query.NewQuery(filter.Contains(core.EnemyTag)).Each(es.world, func(entry *donburi.Entry) {
		count++
	})
	return count
}
