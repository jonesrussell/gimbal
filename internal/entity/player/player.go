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

const (
	// HalfDivisor is used to calculate half of a dimension
	HalfDivisor = 2
	// LogIntervalSeconds is the interval in seconds between position logs
	LogIntervalSeconds = 5
)

// Player represents the player entity in the game
type Player struct {
	coords      *physics.CoordinateSystem
	config      *common.EntityConfig
	sprite      *ebiten.Image
	shape       resolv.IShape
	posAngle    common.Angle // Angle around the circle (position)
	facingAngle common.Angle // Direction the player is facing
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
		"config", map[string]any{
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

	// Start at bottom (180 degrees) and face center (0 degrees)
	posAngle := common.Angle(common.AngleDown) // Position at bottom
	facingAngle := common.Angle(0)             // Face center

	initialPos := coords.CalculateCircularPosition(posAngle)

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
		posAngle:    posAngle,
		facingAngle: facingAngle,
		path:        make([]resolv.Vector, 0),
		speed:       config.Speed,
		size:        config.Size,
		lastLog:     time.Now(),
		logInterval: time.Second * LogIntervalSeconds,
	}

	logger.GlobalLogger.Debug("Player created",
		"position", map[string]float64{
			"x": initialPos.X,
			"y": initialPos.Y,
		},
		"angle", posAngle.ToRadians()/common.DegreesToRadians,
	)

	return player, nil
}

// Update implements Entity interface
func (p *Player) Update() {
	pos := p.GetPosition()
	p.shape.SetPosition(pos.X, pos.Y)
}

// Draw draws the player sprite
func (p *Player) Draw(screen *ebiten.Image) {
	if p.sprite == nil {
		return
	}

	// Calculate sprite offset to center it
	offsetX := float64(p.sprite.Bounds().Dx()) / HalfDivisor
	offsetY := float64(p.sprite.Bounds().Dy()) / HalfDivisor

	// Create transformation options
	op := &ebiten.DrawImageOptions{}

	// Translate to center of sprite
	op.GeoM.Translate(-offsetX, -offsetY)

	// Rotate sprite based on facing angle
	op.GeoM.Rotate(p.facingAngle.ToRadians())

	// Translate to final position
	finalX := p.GetPosition().X + float64(p.size.Width)/HalfDivisor
	finalY := p.GetPosition().Y + float64(p.size.Height)/HalfDivisor
	op.GeoM.Translate(finalX, finalY)

	// Draw the sprite
	screen.DrawImage(p.sprite, op)
}

// GetPosition implements Entity interface
func (p *Player) GetPosition() common.Point {
	return p.coords.CalculateCircularPosition(p.posAngle)
}

// SetPosition implements Movable interface
func (p *Player) SetPosition(pos common.Point) {
	p.coords.SetPosition(pos)
}

// GetSpeed implements Movable interface
func (p *Player) GetSpeed() float64 {
	return p.speed
}

// GetAngle returns the player's current position angle
func (p *Player) GetAngle() common.Angle {
	return p.posAngle
}

// SetAngle sets the player's position angle
func (p *Player) SetAngle(angle common.Angle) {
	p.posAngle = angle.Normalize()
}

// GetFacingAngle returns the direction the player is facing
func (p *Player) GetFacingAngle() common.Angle {
	return p.facingAngle
}

// SetFacingAngle sets the direction the player is facing
func (p *Player) SetFacingAngle(angle common.Angle) {
	p.facingAngle = angle.Normalize()
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
