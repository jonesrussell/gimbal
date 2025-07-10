package ecs

import (
	"image/color"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/ecs/core"
)

// Weapon types
const (
	WeaponTypePrimary = iota
	WeaponTypeSecondary
	WeaponTypeSpecial
)

// Projectile types
const (
	ProjectileTypeEnergy = iota
	ProjectileTypeMissile
	ProjectileTypeLaser
)

// WeaponSystem manages player weapons and projectiles
type WeaponSystem struct {
	world           donburi.World
	config          *common.GameConfig
	fireTimer       float64
	fireInterval    float64
	lastFireTime    time.Time
	projectileSpeed float64
	projectileSize  common.Size
}

// NewWeaponSystem creates a new weapon system
func NewWeaponSystem(world donburi.World, config *common.GameConfig) *WeaponSystem {
	return &WeaponSystem{
		world:           world,
		config:          config,
		fireTimer:       0,
		fireInterval:    DefaultWeaponFireIntervalFrames, // Fire every 10 frames (6 shots per second at 60fps)
		lastFireTime:    time.Now(),
		projectileSpeed: DefaultProjectileSpeed,
		projectileSize:  common.Size{Width: DefaultProjectileSize, Height: DefaultProjectileSize},
	}
}

// Update updates the weapon system
func (ws *WeaponSystem) Update(deltaTime float64) {
	ws.fireTimer += deltaTime

	// Update existing projectiles
	ws.updateProjectiles()
}

// FireWeapon fires a weapon if enough time has passed
func (ws *WeaponSystem) FireWeapon(weaponType int, playerPos common.Point, playerAngle common.Angle) bool {
	if ws.fireTimer < ws.fireInterval {
		return false
	}

	ws.createProjectile(weaponType, playerPos, playerAngle)
	ws.fireTimer = 0
	ws.lastFireTime = time.Now()

	return true
}

// createProjectile creates a new projectile
func (ws *WeaponSystem) createProjectile(weaponType int, startPos common.Point, direction common.Angle) {
	entity := ws.world.Create(
		core.ProjectileTag,
		core.Position,
		core.Sprite,
		core.Movement,
		core.Size,
		core.Speed,
		core.Angle,
	)
	entry := ws.world.Entry(entity)

	// Set position (slightly in front of player)
	angleRad := float64(direction) * common.DegreesToRadians
	offset := ProjectileOffset // Distance from player center
	pos := common.Point{
		X: startPos.X + offset*math.Cos(angleRad),
		Y: startPos.Y - offset*math.Sin(angleRad), // Subtract because Y increases downward
	}
	core.Position.SetValue(entry, pos)

	// Set size
	core.Size.SetValue(entry, ws.projectileSize)

	// Set speed
	core.Speed.SetValue(entry, ws.projectileSpeed)

	// Set angle
	core.Angle.SetValue(entry, direction)

	// Calculate velocity based on direction
	velocity := common.Point{
		X: ws.projectileSpeed * math.Cos(angleRad),
		Y: -ws.projectileSpeed * math.Sin(angleRad), // Negative because Y increases downward
	}

	core.Movement.SetValue(entry, core.MovementData{
		Velocity: velocity,
		MaxSpeed: ws.projectileSpeed,
	})

	// Create simple projectile sprite (white square)
	ws.createProjectileSprite(entry, weaponType)
}

// createProjectileSprite creates a simple sprite for the projectile
func (ws *WeaponSystem) createProjectileSprite(entry *donburi.Entry, weaponType int) {
	// Create a simple colored square based on weapon type
	img := ebiten.NewImage(ws.projectileSize.Width, ws.projectileSize.Height)

	var projectileColor color.Color
	switch weaponType {
	case WeaponTypePrimary:
		projectileColor = color.RGBA{R: 255, G: 255, B: 0, A: 255} // Yellow
	case WeaponTypeSecondary:
		projectileColor = color.RGBA{R: 0, G: 255, B: 255, A: 255} // Cyan
	case WeaponTypeSpecial:
		projectileColor = color.RGBA{R: 255, G: 0, B: 255, A: 255} // Magenta
	default:
		projectileColor = color.RGBA{R: 255, G: 255, B: 255, A: 255} // White
	}

	img.Fill(projectileColor)
	core.Sprite.SetValue(entry, img)
}

// updateProjectiles updates all projectile entities
func (ws *WeaponSystem) updateProjectiles() {
	query.NewQuery(
		filter.And(
			filter.Contains(core.ProjectileTag),
			filter.Contains(core.Position),
			filter.Contains(core.Movement),
		),
	).Each(ws.world, func(entry *donburi.Entry) {
		pos := core.Position.Get(entry)
		mov := core.Movement.Get(entry)

		// Update position based on velocity
		pos.X += mov.Velocity.X
		pos.Y += mov.Velocity.Y

		// Check if projectile is off screen
		if ws.isOffScreen(*pos) {
			// Remove projectile from world
			ws.world.Remove(entry.Entity())
		}
	})
}

// isOffScreen checks if a position is off screen
func (ws *WeaponSystem) isOffScreen(pos common.Point) bool {
	margin := ProjectileMargin
	return pos.X < -margin ||
		pos.X > float64(ws.config.ScreenSize.Width)+margin ||
		pos.Y < -margin ||
		pos.Y > float64(ws.config.ScreenSize.Height)+margin
}

// GetFireInterval returns the current fire interval
func (ws *WeaponSystem) GetFireInterval() float64 {
	return ws.fireInterval
}

// SetFireInterval sets the fire interval
func (ws *WeaponSystem) SetFireInterval(interval float64) {
	ws.fireInterval = interval
}

// GetProjectileSpeed returns the current projectile speed
func (ws *WeaponSystem) GetProjectileSpeed() float64 {
	return ws.projectileSpeed
}

// SetProjectileSpeed sets the projectile speed
func (ws *WeaponSystem) SetProjectileSpeed(speed float64) {
	ws.projectileSpeed = speed
}
