package weapon

import (
	"image/color"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/config"
	"github.com/jonesrussell/gimbal/internal/ecs/core"
	gameMath "github.com/jonesrussell/gimbal/internal/math"
)

// Weapon system constants
const (
	DefaultWeaponFireIntervalFrames = 10   // 10 frames between shots (6 shots/sec at 60fps)
	DefaultProjectileSpeed          = 8.0  // Pixels per frame
	DefaultProjectileSize           = 4    // Projectile size in pixels
	ProjectileOffset                = 20.0 // Distance from player center
	ProjectileMargin                = 50.0 // Screen margin for cleanup
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
	config          *config.GameConfig
	fireTimer       float64
	fireInterval    float64
	lastFireTime    time.Time
	projectileSpeed float64
	projectileSize  struct {
		Width, Height int
	}
	projectileSprites map[int]*ebiten.Image // Sprite cache
}

// NewWeaponSystem creates a new weapon management system with the provided configuration
func NewWeaponSystem(world donburi.World, gameConfig *config.GameConfig) *WeaponSystem {
	ws := &WeaponSystem{
		world:             world,
		config:            gameConfig,
		fireTimer:         0,
		fireInterval:      DefaultWeaponFireIntervalFrames, // Fire every 10 frames (6 shots per second at 60fps)
		lastFireTime:      time.Now(),
		projectileSpeed:   DefaultProjectileSpeed,
		projectileSize:    struct{ Width, Height int }{Width: DefaultProjectileSize, Height: DefaultProjectileSize},
		projectileSprites: make(map[int]*ebiten.Image),
	}
	ws.initializeProjectileSprites()
	return ws
}

// Update updates the weapon system
func (ws *WeaponSystem) Update(deltaTime float64) {
	ws.fireTimer += deltaTime

	// Update existing projectiles
	ws.updateProjectiles()
}

// FireWeapon fires a weapon if enough time has passed
func (ws *WeaponSystem) FireWeapon(weaponType int, playerPos common.Point, playerAngle gameMath.Angle) bool {
	if ws.fireTimer < ws.fireInterval {
		return false
	}

	ws.createProjectile(weaponType, playerPos, playerAngle)
	ws.fireTimer = 0
	ws.lastFireTime = time.Now()

	return true
}

// createProjectile creates a new projectile
func (ws *WeaponSystem) createProjectile(weaponType int, startPos common.Point, direction gameMath.Angle) {
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
	angleRad := float64(direction) * gameMath.DegreesToRadians
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

	// Calculate direction toward screen center (Gyruss-style shooting)
	centerX := float64(ws.config.ScreenSize.Width) / 2
	centerY := float64(ws.config.ScreenSize.Height) / 2

	// Direction vector from player to center
	dirX := centerX - pos.X
	dirY := centerY - pos.Y

	// Normalize the direction vector
	distance := math.Sqrt(dirX*dirX + dirY*dirY)
	if distance > 0 {
		dirX /= distance
		dirY /= distance
	}

	// Apply speed to create velocity
	velocity := common.Point{
		X: dirX * ws.projectileSpeed,
		Y: dirY * ws.projectileSpeed,
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
	if sprite, exists := ws.projectileSprites[weaponType]; exists {
		core.Sprite.SetValue(entry, sprite)
	} else {
		core.Sprite.SetValue(entry, ws.projectileSprites[WeaponTypePrimary])
	}
}

// initializeProjectileSprites pre-creates one image per weapon type
func (ws *WeaponSystem) initializeProjectileSprites() {
	ws.projectileSprites = make(map[int]*ebiten.Image)
	weaponConfigs := map[int]color.RGBA{
		WeaponTypePrimary:   {R: 255, G: 255, B: 0, A: 255}, // Yellow
		WeaponTypeSecondary: {R: 0, G: 255, B: 255, A: 255}, // Cyan
		WeaponTypeSpecial:   {R: 255, G: 0, B: 255, A: 255}, // Magenta
	}
	for weaponType, projectileColor := range weaponConfigs {
		img := ebiten.NewImage(ws.projectileSize.Width, ws.projectileSize.Height)
		img.Fill(projectileColor)
		ws.projectileSprites[weaponType] = img
	}
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
	return pos.X < -margin || pos.X > float64(ws.config.ScreenSize.Width)+margin ||
		pos.Y < -margin || pos.Y > float64(ws.config.ScreenSize.Height)+margin
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

// SetProjectileSize sets the projectile size
func (ws *WeaponSystem) SetProjectileSize(size struct {
	Width, Height int
},
) {
	ws.projectileSize = size
}
