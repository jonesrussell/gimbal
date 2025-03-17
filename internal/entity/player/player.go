package player

import (
	"errors"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/logger"
	"github.com/jonesrussell/gimbal/internal/physics"
	"github.com/solarlune/resolv"
)

// Player represents the player entity in the game
type Player struct {
	coords      *physics.CoordinateSystem
	config      *common.EntityConfig
	sprite      *ebiten.Image
	shape       resolv.IShape
	angle       common.Angle
	path        []resolv.Vector
	speed       float64
	size        common.Size
	lastLog     time.Time
	logInterval time.Duration
}

// New creates a new player instance
func New(config *common.EntityConfig, sprite *ebiten.Image) (*Player, error) {
	if config == nil {
		return nil, errors.New("config cannot be nil")
	}
	if sprite == nil {
		return nil, errors.New("sprite cannot be nil")
	}

	logger.GlobalLogger.Debug("Creating new player with config",
		"config", map[string]interface{}{
			"position": map[string]float64{
				"x": config.Position.X,
				"y": config.Position.Y,
			},
			"size": map[string]int{
				"width":  config.Size.Width,
				"height": config.Size.Height,
			},
			"radius": config.Radius,
			"speed":  config.Speed,
		},
	)

	// Create coordinate system for circular movement
	// Center point should be at the center of the screen
	center := common.Point{
		X: config.Position.X,
		Y: config.Position.Y,
	}
	coords := physics.NewCoordinateSystem(center, config.Radius)

	// Calculate initial position at the bottom of the screen (270 degrees)
	initialAngle := common.Angle(270 * common.DegreesToRadians)
	logger.GlobalLogger.Debug("Setting initial angle",
		"angle_rad", initialAngle.ToRadians(),
		"angle_deg", initialAngle.ToRadians()/common.DegreesToRadians,
	)

	initialPos := coords.CalculateCircularPosition(initialAngle)

	// Create player collision shape as a rectangle
	shape := resolv.NewRectangle(
		initialPos.X,
		initialPos.Y,
		float64(config.Size.Width),
		float64(config.Size.Height),
	)

	// Create player with initial position
	player := &Player{
		coords:      coords,
		config:      config,
		sprite:      sprite,
		shape:       shape,
		angle:       initialAngle,
		path:        make([]resolv.Vector, 0),
		speed:       config.Speed,
		size:        config.Size,
		lastLog:     time.Now(),
		logInterval: time.Second, // Only log once per second
	}

	logger.GlobalLogger.Debug("Player created",
		"sprite_size", map[string]int{
			"width":  sprite.Bounds().Dx(),
			"height": sprite.Bounds().Dy(),
		},
		"initial_angle", initialAngle.ToRadians(),
		"initial_position", map[string]float64{
			"x": initialPos.X,
			"y": initialPos.Y,
		},
		"center", map[string]float64{
			"x": center.X,
			"y": center.Y,
		},
		"radius", config.Radius,
	)

	return player, nil
}

// Update implements Entity interface
func (p *Player) Update() {
	pos := p.GetPosition()
	p.shape.SetPosition(pos.X, pos.Y)
}

// Draw implements Entity interface
func (p *Player) Draw(screen *ebiten.Image) {
	// Only log once per second
	now := time.Now()
	if now.Sub(p.lastLog) >= p.logInterval {
		pos := p.GetPosition()
		logger.GlobalLogger.Debug("Drawing player",
			"position", map[string]float64{
				"x": pos.X,
				"y": pos.Y,
			},
			"angle", p.GetAngle().ToRadians(),
			"scale", map[string]float64{
				"x": float64(p.size.Width) / float64(p.sprite.Bounds().Dx()),
				"y": float64(p.size.Height) / float64(p.sprite.Bounds().Dy()),
			},
			"final_position", map[string]float64{
				"x": pos.X + float64(p.size.Width)/2,
				"y": pos.Y + float64(p.size.Height)/2,
			},
		)
		p.lastLog = now
	}

	// Create GeoM for transformations
	geoM := ebiten.GeoM{}

	// Calculate scale based on sprite size
	scaleX := float64(p.size.Width) / float64(p.sprite.Bounds().Dx())
	scaleY := float64(p.size.Height) / float64(p.sprite.Bounds().Dy())

	// Move to origin for rotation
	offsetX := float64(p.sprite.Bounds().Dx()) / 2
	offsetY := float64(p.sprite.Bounds().Dy()) / 2
	geoM.Translate(-offsetX, -offsetY)

	// Apply rotation
	geoM.Rotate(p.GetAngle().ToRadians())

	// Apply scale
	geoM.Scale(scaleX, scaleY)

	// Move to final position
	finalX := p.GetPosition().X + float64(p.size.Width)/2
	finalY := p.GetPosition().Y + float64(p.size.Height)/2
	geoM.Translate(finalX, finalY)

	// Draw the sprite
	screen.DrawImage(p.sprite, &ebiten.DrawImageOptions{
		GeoM: geoM,
	})
}

// GetPosition implements Entity interface
func (p *Player) GetPosition() common.Point {
	return p.coords.CalculateCircularPosition(p.angle)
}

// SetPosition implements Movable interface
func (p *Player) SetPosition(pos common.Point) {
	p.coords.SetPosition(pos)
}

// GetSpeed implements Movable interface
func (p *Player) GetSpeed() float64 {
	return p.speed
}

// GetAngle returns the player's current angle
func (p *Player) GetAngle() common.Angle {
	return p.angle
}

// SetAngle sets the player's angle
func (p *Player) SetAngle(angle common.Angle) {
	p.angle = angle.Normalize()
}

// GetBounds implements Collidable interface
func (p *Player) GetBounds() common.Size {
	return p.size
}

// CheckCollision implements Collidable interface
func (p *Player) CheckCollision(other common.Collidable) bool {
	if p.shape == nil {
		return false
	}

	// Get the bounds of the other object
	otherBounds := other.GetBounds()
	otherPos := other.GetPosition()

	// Create a temporary shape for collision check
	otherShape := resolv.NewRectangle(
		otherPos.X,
		otherPos.Y,
		float64(otherBounds.Width),
		float64(otherBounds.Height),
	)

	// Check for intersection
	intersection := p.shape.Intersection(otherShape)
	return !intersection.IsEmpty()
}
