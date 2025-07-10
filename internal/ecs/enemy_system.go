package ecs

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/jonesrussell/gimbal/internal/common"
)

// Enemy types
const (
	EnemyTypeSwarmDrone = iota
	EnemyTypeHeavyCruiser
	EnemyTypeBoss
	EnemyTypeAsteroid
)

// Enemy spawn patterns
const (
	SpawnPatternCircle = iota
	SpawnPatternSpiral
	SpawnPatternWave
	SpawnPatternRandom
)

// EnemySystem manages enemy spawning, movement, and behavior
type EnemySystem struct {
	world          donburi.World
	config         *common.GameConfig
	spawnTimer     float64
	spawnInterval  float64
	currentWave    int
	enemiesInWave  int
	enemiesSpawned int
	waveComplete   bool
	level          int
	difficulty     float64
}

// NewEnemySystem creates a new enemy system
func NewEnemySystem(world donburi.World, config *common.GameConfig) *EnemySystem {
	return &EnemySystem{
		world:          world,
		config:         config,
		spawnTimer:     0,
		spawnInterval:  60, // Spawn every 60 frames (1 second at 60fps)
		currentWave:    1,
		enemiesInWave:  5,
		enemiesSpawned: 0,
		waveComplete:   false,
		level:          1,
		difficulty:     1.0,
	}
}

// Update updates the enemy system
func (es *EnemySystem) Update(deltaTime float64) {
	es.spawnTimer += deltaTime

	// Check if it's time to spawn enemies
	if es.spawnTimer >= es.spawnInterval && !es.waveComplete {
		es.spawnEnemy()
		es.spawnTimer = 0
		es.enemiesSpawned++

		// Check if wave is complete
		if es.enemiesSpawned >= es.enemiesInWave {
			es.waveComplete = true
			es.startNextWave()
		}
	}

	// Update existing enemies
	es.updateEnemies()
}

// spawnEnemy spawns a new enemy based on current wave and difficulty
func (es *EnemySystem) spawnEnemy() {
	enemyType := es.selectEnemyType()
	spawnPattern := es.selectSpawnPattern()

	// Calculate spawn position based on pattern
	spawnPos := es.calculateSpawnPosition(spawnPattern)

	// Create enemy entity
	es.createEnemy(enemyType, spawnPos)
}

// selectEnemyType selects enemy type based on current wave and difficulty
func (es *EnemySystem) selectEnemyType() int {
	// Simple probability-based selection
	//nolint:gosec // Game logic doesn't need cryptographic randomness
	r := rand.Float64()

	if es.currentWave >= 5 && r < 0.1 {
		return EnemyTypeBoss
	} else if es.currentWave >= 3 && r < 0.3 {
		return EnemyTypeHeavyCruiser
	} else if r < 0.1 {
		return EnemyTypeAsteroid
	} else {
		return EnemyTypeSwarmDrone
	}
}

// selectSpawnPattern selects spawn pattern based on enemy type and wave
func (es *EnemySystem) selectSpawnPattern() int {
	//nolint:gosec // Game logic doesn't need cryptographic randomness
	r := rand.Float64()

	if r < 0.4 {
		return SpawnPatternCircle
	} else if r < 0.7 {
		return SpawnPatternSpiral
	} else if r < 0.9 {
		return SpawnPatternWave
	} else {
		return SpawnPatternRandom
	}
}

// calculateSpawnPosition calculates spawn position based on pattern
func (es *EnemySystem) calculateSpawnPosition(pattern int) common.Point {
	center := common.Point{
		X: float64(es.config.ScreenSize.Width) / 2,
		Y: float64(es.config.ScreenSize.Height) / 2,
	}

	switch pattern {
	case SpawnPatternCircle:
		// Spawn in a circle around the center
		//nolint:gosec // Game logic doesn't need cryptographic randomness
		angle := rand.Float64() * 2 * math.Pi
		//nolint:gosec // Game logic doesn't need cryptographic randomness
		radius := 50.0 + rand.Float64()*100.0
		return common.Point{
			X: center.X + radius*math.Cos(angle),
			Y: center.Y + radius*math.Sin(angle),
		}

	case SpawnPatternSpiral:
		// Spawn in a spiral pattern
		//nolint:gosec // Game logic doesn't need cryptographic randomness
		angle := rand.Float64() * 2 * math.Pi
		//nolint:gosec // Game logic doesn't need cryptographic randomness
		radius := 30.0 + rand.Float64()*150.0
		return common.Point{
			X: center.X + radius*math.Cos(angle),
			Y: center.Y + radius*math.Sin(angle),
		}

	case SpawnPatternWave:
		// Spawn in a wave pattern from one side
		//nolint:gosec // Game logic doesn't need cryptographic randomness
		side := rand.Intn(4) // 0: top, 1: right, 2: bottom, 3: left
		var x, y float64

		switch side {
		case 0: // top
			//nolint:gosec // Game logic doesn't need cryptographic randomness
			x = center.X + (rand.Float64()-0.5)*200
			y = -50
		case 1: // right
			x = float64(es.config.ScreenSize.Width) + 50
			//nolint:gosec // Game logic doesn't need cryptographic randomness
			y = center.Y + (rand.Float64()-0.5)*200
		case 2: // bottom
			//nolint:gosec // Game logic doesn't need cryptographic randomness
			x = center.X + (rand.Float64()-0.5)*200
			y = float64(es.config.ScreenSize.Height) + 50
		case 3: // left
			x = -50
			//nolint:gosec // Game logic doesn't need cryptographic randomness
			y = center.Y + (rand.Float64()-0.5)*200
		}

		return common.Point{X: x, Y: y}

	case SpawnPatternRandom:
		// Random position within screen bounds
		//nolint:gosec // Game logic doesn't need cryptographic randomness
		return common.Point{
			X: rand.Float64() * float64(es.config.ScreenSize.Width),
			//nolint:gosec // Game logic doesn't need cryptographic randomness
			Y: rand.Float64() * float64(es.config.ScreenSize.Height),
		}

	default:
		return center
	}
}

// createEnemy creates an enemy entity with appropriate components
func (es *EnemySystem) createEnemy(enemyType int, spawnPos common.Point) {
	entity := es.world.Create(EnemyTag, Position, Sprite, Movement, Size, Speed, Health)
	entry := es.world.Entry(entity)

	// Set position
	Position.SetValue(entry, spawnPos)

	// Create enemy sprite
	es.createEnemySprite(entry, enemyType)

	// Set enemy-specific properties
	switch enemyType {
	case EnemyTypeSwarmDrone:
		es.setupSwarmDrone(entry)
	case EnemyTypeHeavyCruiser:
		es.setupHeavyCruiser(entry)
	case EnemyTypeBoss:
		es.setupBoss(entry)
	case EnemyTypeAsteroid:
		es.setupAsteroid(entry)
	}
}

// createEnemySprite creates a sprite for the enemy based on type
func (es *EnemySystem) createEnemySprite(entry *donburi.Entry, enemyType int) {
	var size common.Size
	var enemyColor color.Color

	switch enemyType {
	case EnemyTypeSwarmDrone:
		size = common.Size{Width: 16, Height: 16}
		enemyColor = color.RGBA{255, 0, 0, 255} // Red
	case EnemyTypeHeavyCruiser:
		size = common.Size{Width: 32, Height: 32}
		enemyColor = color.RGBA{255, 100, 0, 255} // Orange
	case EnemyTypeBoss:
		size = common.Size{Width: 64, Height: 64}
		enemyColor = color.RGBA{128, 0, 128, 255} // Purple
	case EnemyTypeAsteroid:
		size = common.Size{Width: 24, Height: 24}
		enemyColor = color.RGBA{128, 128, 128, 255} // Gray
	default:
		size = common.Size{Width: 16, Height: 16}
		enemyColor = color.RGBA{255, 0, 0, 255} // Red
	}

	// Create sprite
	img := ebiten.NewImage(size.Width, size.Height)
	img.Fill(enemyColor)
	Sprite.SetValue(entry, img)
}

// movementTowardsCenter sets up movement for an enemy towards the center with given speed and maxSpeed
func (es *EnemySystem) movementTowardsCenter(entry *donburi.Entry, speed, maxSpeed float64) {
	center := common.Point{
		X: float64(es.config.ScreenSize.Width) / 2,
		Y: float64(es.config.ScreenSize.Height) / 2,
	}
	pos := Position.Get(entry)
	dx := center.X - pos.X
	dy := center.Y - pos.Y
	distance := math.Sqrt(dx*dx + dy*dy)
	if distance > 0 {
		velocity := common.Point{
			X: (dx / distance) * speed * es.difficulty,
			Y: (dy / distance) * speed * es.difficulty,
		}
		Movement.SetValue(entry, MovementData{
			Velocity: velocity,
			MaxSpeed: maxSpeed * es.difficulty,
		})
	}
}

// setupSwarmDrone configures a swarm drone enemy
func (es *EnemySystem) setupSwarmDrone(entry *donburi.Entry) {
	// Small, fast, weak enemy
	Size.SetValue(entry, common.Size{Width: 16, Height: 16})
	Speed.SetValue(entry, 2.0*es.difficulty)
	Health.SetValue(entry, 1)
	// Movement towards center
	es.movementTowardsCenter(entry, 2.0, 3.0)
}

// setupHeavyCruiser configures a heavy cruiser enemy
func (es *EnemySystem) setupHeavyCruiser(entry *donburi.Entry) {
	// Large, slow, strong enemy
	Size.SetValue(entry, common.Size{Width: 32, Height: 32})
	Speed.SetValue(entry, 1.0*es.difficulty)
	Health.SetValue(entry, 3)
	// Movement towards center
	es.movementTowardsCenter(entry, 1.0, 1.5)
}

// setupBoss configures a boss enemy
func (es *EnemySystem) setupBoss(entry *donburi.Entry) {
	// Very large, powerful boss enemy
	Size.SetValue(entry, common.Size{Width: 64, Height: 64})
	Speed.SetValue(entry, 0.5*es.difficulty)
	Health.SetValue(entry, 10)

	// Complex movement pattern (orbit around center)
	center := common.Point{
		X: float64(es.config.ScreenSize.Width) / 2,
		Y: float64(es.config.ScreenSize.Height) / 2,
	}

	// Start at a random angle
	//nolint:gosec // Game logic doesn't need cryptographic randomness
	angle := rand.Float64() * 2 * math.Pi
	radius := 100.0

	pos := common.Point{
		X: center.X + radius*math.Cos(angle),
		Y: center.Y + radius*math.Sin(angle),
	}
	Position.SetValue(entry, pos)

	// Orbital movement
	Movement.SetValue(entry, MovementData{
		Velocity: common.Point{X: 0, Y: 0}, // Will be calculated in update
		MaxSpeed: 1.0 * es.difficulty,
	})
}

// setupAsteroid configures an asteroid enemy
func (es *EnemySystem) setupAsteroid(entry *donburi.Entry) {
	// Medium-sized environmental hazard
	Size.SetValue(entry, common.Size{Width: 24, Height: 24})
	Speed.SetValue(entry, 1.5*es.difficulty)
	Health.SetValue(entry, 2)

	// Random movement direction
	//nolint:gosec // Game logic doesn't need cryptographic randomness
	angle := rand.Float64() * 2 * math.Pi
	velocity := common.Point{
		X: math.Cos(angle) * 1.5 * es.difficulty,
		Y: math.Sin(angle) * 1.5 * es.difficulty,
	}

	Movement.SetValue(entry, MovementData{
		Velocity: velocity,
		MaxSpeed: 2.0 * es.difficulty,
	})
}

// updateEnemies updates all enemy entities
func (es *EnemySystem) updateEnemies() {
	query.NewQuery(
		filter.And(
			filter.Contains(EnemyTag),
			filter.Contains(Position),
			filter.Contains(Movement),
		),
	).Each(es.world, func(entry *donburi.Entry) {
		pos := Position.Get(entry)
		mov := Movement.Get(entry)

		// Update position based on velocity
		pos.X += mov.Velocity.X
		pos.Y += mov.Velocity.Y

		// Check if enemy is off screen
		if es.isOffScreen(*pos) {
			// Remove enemy from world
			es.world.Remove(entry.Entity())
		}
	})
}

// isOffScreen checks if a position is off screen
func (es *EnemySystem) isOffScreen(pos common.Point) bool {
	margin := 50.0
	return pos.X < -margin ||
		pos.X > float64(es.config.ScreenSize.Width)+margin ||
		pos.Y < -margin ||
		pos.Y > float64(es.config.ScreenSize.Height)+margin
}

// startNextWave starts the next wave of enemies
func (es *EnemySystem) startNextWave() {
	es.currentWave++
	es.enemiesSpawned = 0
	es.waveComplete = false

	// Increase difficulty
	es.difficulty += 0.1

	// Increase enemies per wave
	es.enemiesInWave = 5 + es.currentWave*2

	// Decrease spawn interval (faster spawning)
	es.spawnInterval = math.Max(20, es.spawnInterval-5)
}

// GetCurrentWave returns the current wave number
func (es *EnemySystem) GetCurrentWave() int {
	return es.currentWave
}

// GetDifficulty returns the current difficulty multiplier
func (es *EnemySystem) GetDifficulty() float64 {
	return es.difficulty
}

// GetEnemiesRemaining returns the number of enemies remaining in current wave
func (es *EnemySystem) GetEnemiesRemaining() int {
	return es.enemiesInWave - es.enemiesSpawned
}
