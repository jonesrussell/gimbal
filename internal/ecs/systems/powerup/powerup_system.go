package powerup

import (
	"context"
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/config"
	"github.com/jonesrussell/gimbal/internal/ecs/core"
)

const (
	PowerUpSize       = 16
	PowerUpOrbitSpeed = 30.0 // degrees per second
	PowerUpLifetime   = 10.0 // seconds before despawn
	PowerUpDropChance = 0.15 // 15% chance to drop
	CollisionDistance = 30.0 // pixels for collection
	DoubleFireRate    = 2.0  // Fire rate multiplier for double shot
)

// PowerUpSystem manages power-up spawning, movement, and collection
type PowerUpSystem struct {
	world           donburi.World
	config          *config.GameConfig
	logger          common.Logger
	screenCenter    common.Point
	sprites         map[core.PowerUpType]*ebiten.Image
	playerHasDouble bool
	doubleRemaining time.Duration
}

// NewPowerUpSystem creates a new power-up system
func NewPowerUpSystem(
	world donburi.World,
	cfg *config.GameConfig,
	logger common.Logger,
) *PowerUpSystem {
	ps := &PowerUpSystem{
		world:  world,
		config: cfg,
		logger: logger,
		screenCenter: common.Point{
			X: float64(cfg.ScreenSize.Width) / 2,
			Y: float64(cfg.ScreenSize.Height) / 2,
		},
		sprites:         make(map[core.PowerUpType]*ebiten.Image),
		playerHasDouble: false,
		doubleRemaining: 0,
	}

	ps.createSprites()

	return ps
}

// createSprites creates power-up sprites
func (ps *PowerUpSystem) createSprites() {
	// Double shot sprite (yellow)
	doubleSprite := ebiten.NewImage(PowerUpSize, PowerUpSize)
	doubleSprite.Fill(color.RGBA{R: 255, G: 255, B: 0, A: 255})
	ps.sprites[core.PowerUpDoubleShot] = doubleSprite

	// Extra life sprite (green)
	lifeSprite := ebiten.NewImage(PowerUpSize, PowerUpSize)
	lifeSprite.Fill(color.RGBA{R: 0, G: 255, B: 0, A: 255})
	ps.sprites[core.PowerUpExtraLife] = lifeSprite
}

// Update processes power-up logic
func (ps *PowerUpSystem) Update(ctx context.Context, deltaTime float64) error {
	// Check for cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Update existing power-ups (movement, lifetime)
	ps.updatePowerUps(deltaTime)

	// Check for collection
	ps.checkCollection()

	// Update active power-up effects
	ps.updateEffects(deltaTime)

	return nil
}

// updatePowerUps updates all power-up entities
func (ps *PowerUpSystem) updatePowerUps(deltaTime float64) {
	query.NewQuery(
		filter.And(
			filter.Contains(core.PowerUpTag),
			filter.Contains(core.PowerUpData),
			filter.Contains(core.Position),
		),
	).Each(ps.world, func(entry *donburi.Entry) {
		powerUpData := core.PowerUpData.Get(entry)
		pos := core.Position.Get(entry)

		// Update lifetime
		powerUpData.LifeTime += time.Duration(deltaTime * float64(time.Second))

		// Check for despawn
		if powerUpData.MaxLifeTime > 0 && powerUpData.LifeTime >= powerUpData.MaxLifeTime {
			ps.world.Remove(entry.Entity())
			return
		}

		// Orbital movement around center
		powerUpData.OrbitalAngle += powerUpData.OrbitalSpeed * deltaTime

		// Calculate orbit radius
		dx := pos.X - ps.screenCenter.X
		dy := pos.Y - ps.screenCenter.Y
		radius := math.Sqrt(dx*dx + dy*dy)

		// Update position
		angleRad := powerUpData.OrbitalAngle * math.Pi / 180
		pos.X = ps.screenCenter.X + radius*math.Cos(angleRad)
		pos.Y = ps.screenCenter.Y + radius*math.Sin(angleRad)

		core.PowerUpData.SetValue(entry, *powerUpData)
	})
}

// checkCollection checks if player collected any power-ups
func (ps *PowerUpSystem) checkCollection() {
	// Find player position
	var playerEntry *donburi.Entry
	var playerPos common.Point

	query.NewQuery(
		filter.Contains(core.PlayerTag),
	).Each(ps.world, func(entry *donburi.Entry) {
		playerEntry = entry
		if entry.HasComponent(core.Position) {
			pos := core.Position.Get(entry)
			playerPos = *pos
		}
	})

	if playerEntry == nil {
		return
	}

	// Check each power-up for collection
	query.NewQuery(
		filter.And(
			filter.Contains(core.PowerUpTag),
			filter.Contains(core.PowerUpData),
			filter.Contains(core.Position),
		),
	).Each(ps.world, func(entry *donburi.Entry) {
		pos := core.Position.Get(entry)
		powerUpData := core.PowerUpData.Get(entry)

		// Check distance to player
		dx := pos.X - playerPos.X
		dy := pos.Y - playerPos.Y
		distance := math.Sqrt(dx*dx + dy*dy)

		if distance < CollisionDistance {
			ps.collectPowerUp(playerEntry, powerUpData)
			ps.world.Remove(entry.Entity())
		}
	})
}

// collectPowerUp applies the power-up effect
func (ps *PowerUpSystem) collectPowerUp(playerEntry *donburi.Entry, data *core.PowerUpTypeData) {
	ps.logger.Debug("Power-up collected", "type", data.Type)

	switch data.Type {
	case core.PowerUpDoubleShot:
		ps.playerHasDouble = true
		ps.doubleRemaining = data.Duration
		if ps.doubleRemaining == 0 {
			ps.doubleRemaining = 10 * time.Second // Default 10 seconds
		}
		ps.logger.Info("Double shot activated", "duration", ps.doubleRemaining)

	case core.PowerUpExtraLife:
		if playerEntry.HasComponent(core.Health) {
			health := core.Health.Get(playerEntry)
			health.Current++
			if health.Current > health.Maximum {
				health.Maximum = health.Current
			}
			ps.logger.Info("Extra life collected", "lives", health.Current)
		}
	}
}

// updateEffects updates active power-up effects
func (ps *PowerUpSystem) updateEffects(deltaTime float64) {
	if ps.playerHasDouble && ps.doubleRemaining > 0 {
		ps.doubleRemaining -= time.Duration(deltaTime * float64(time.Second))
		if ps.doubleRemaining <= 0 {
			ps.playerHasDouble = false
			ps.doubleRemaining = 0
			ps.logger.Info("Double shot expired")
		}
	}
}

// SpawnPowerUp spawns a power-up at the given position
func (ps *PowerUpSystem) SpawnPowerUp(position common.Point, powerUpType core.PowerUpType) {
	entity := ps.world.Create(
		core.PowerUpTag,
		core.PowerUpData,
		core.Position,
		core.Sprite,
		core.Size,
	)
	entry := ps.world.Entry(entity)

	// Set position
	core.Position.SetValue(entry, position)

	// Set sprite
	sprite := ps.sprites[powerUpType]
	if sprite == nil {
		sprite = ps.sprites[core.PowerUpDoubleShot] // Default
	}
	core.Sprite.SetValue(entry, sprite)

	// Set size
	core.Size.SetValue(entry, config.Size{Width: PowerUpSize, Height: PowerUpSize})

	// Set power-up data
	// Calculate initial orbital angle from position
	dx := position.X - ps.screenCenter.X
	dy := position.Y - ps.screenCenter.Y
	angle := math.Atan2(dy, dx) * 180 / math.Pi

	duration := time.Duration(0)
	if powerUpType == core.PowerUpDoubleShot {
		duration = 10 * time.Second
	}

	core.PowerUpData.SetValue(entry, core.PowerUpTypeData{
		Type:         powerUpType,
		Duration:     duration,
		OrbitalAngle: angle,
		OrbitalSpeed: PowerUpOrbitSpeed,
		LifeTime:     0,
		MaxLifeTime:  time.Duration(PowerUpLifetime * float64(time.Second)),
	})

	ps.logger.Debug("Power-up spawned",
		"type", powerUpType,
		"position", position)
}

// TrySpawnPowerUp attempts to spawn a power-up with configured drop chance
func (ps *PowerUpSystem) TrySpawnPowerUp(position common.Point) {
	//nolint:gosec // Game logic randomness is acceptable
	if rand.Float64() < PowerUpDropChance {
		//nolint:gosec // Weighted random selection (70% double shot, 30% extra life)
		if rand.Float64() < 0.7 {
			ps.SpawnPowerUp(position, core.PowerUpDoubleShot)
		} else {
			ps.SpawnPowerUp(position, core.PowerUpExtraLife)
		}
	}
}

// HasDoubleShot returns whether the player has double shot active
func (ps *PowerUpSystem) HasDoubleShot() bool {
	return ps.playerHasDouble
}

// GetDoubleRemainingTime returns the remaining time for double shot
func (ps *PowerUpSystem) GetDoubleRemainingTime() time.Duration {
	return ps.doubleRemaining
}
