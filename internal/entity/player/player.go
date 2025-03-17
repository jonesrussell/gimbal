package player

import (
	"errors"
	"time"

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

// Drawable interface defines the methods required for drawing
type Drawable interface {
	Draw(screen any, op any)
}

// Player represents the player entity in the game
type Player struct {
	coords      *physics.CoordinateSystem
	config      *common.EntityConfig
	sprite      Drawable
	shape       resolv.IShape
	posAngle    common.Angle // Angle around the circle (position)
	facingAngle common.Angle // Angle the player is facing
	lastLog     time.Time
	logInterval time.Duration
}

// New creates a new player instance
func New(config *common.EntityConfig, sprite Drawable) (*Player, error) {
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
			"speed":  config.Speed,
			"radius": config.Radius,
		},
	)

	// Create coordinate system
	coords := physics.NewCoordinateSystem(config.Position, config.Radius)

	// Create player with initial position
	player := &Player{
		coords:      coords,
		config:      config,
		sprite:      sprite,
		posAngle:    common.Angle(0),
		facingAngle: common.Angle(0),
		lastLog:     time.Now(),
		logInterval: time.Second * LogIntervalSeconds,
	}

	// Create collision shape
	player.shape = resolv.NewRectangle(
		config.Position.X-float64(config.Size.Width)/2,
		config.Position.Y-float64(config.Size.Height)/2,
		float64(config.Size.Width),
		float64(config.Size.Height),
	)

	return player, nil
}

// Update implements Entity interface
func (p *Player) Update() {
	pos := p.GetPosition()
	p.shape.SetPosition(pos.X, pos.Y)
}

// Draw implements the Drawable interface
func (p *Player) Draw(screen any, op any) {
	if p.sprite != nil {
		p.sprite.Draw(screen, op)
	}
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
	return p.config.Speed
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
	return p.config.Size
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

// Config returns the player's configuration
func (p *Player) Config() *common.EntityConfig {
	return p.config
}

// Sprite returns the player's sprite
func (p *Player) Sprite() Drawable {
	return p.sprite
}
