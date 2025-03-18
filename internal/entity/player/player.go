package player

import (
	"errors"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/physics"
	"github.com/solarlune/resolv"
)

const (
	// HalfDivisor is used to calculate half of a dimension
	HalfDivisor = 2
	// LogIntervalSeconds is the interval in seconds between position logs
	LogIntervalSeconds = 5
	// DefaultPlayerSize is the default size of the player
	DefaultPlayerSize = 100
	// DegreesToRadians is used to convert degrees to radians
	DegreesToRadians = math.Pi / 180
	// RadiansToDegrees is used to convert radians to degrees
	RadiansToDegrees = 180 / math.Pi
	// FacingCenterOffset is the angle offset to make the player face the center
	FacingCenterOffset = 180
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
	logger      common.Logger
}

// New creates a new player instance
func New(config *common.EntityConfig, sprite Drawable, logger common.Logger) (*Player, error) {
	if config == nil {
		return nil, errors.New("config cannot be nil")
	}
	if sprite == nil {
		return nil, errors.New("sprite cannot be nil")
	}
	if logger == nil {
		return nil, errors.New("logger cannot be nil")
	}

	logger.Debug("Creating new player",
		"position", map[string]float64{
			"x": config.Position.X,
			"y": config.Position.Y,
		},
		"size", map[string]int{
			"width":  config.Size.Width,
			"height": config.Size.Height,
		},
		"speed", config.Speed,
		"radius", config.Radius,
	)

	// Create coordinate system
	coords := physics.NewCoordinateSystem(config.Position, config.Radius)

	// Create player with initial position
	player := &Player{
		coords:      coords,
		config:      config,
		sprite:      sprite,
		posAngle:    common.AngleRight, // Start at 0 degrees (right side)
		facingAngle: common.AngleLeft,  // Face the center
		lastLog:     time.Now(),
		logInterval: time.Second * LogIntervalSeconds,
		logger:      logger,
	}

	// Create collision shape
	player.shape = resolv.NewRectangle(
		config.Position.X-float64(config.Size.Width)/HalfDivisor,
		config.Position.Y-float64(config.Size.Height)/HalfDivisor,
		float64(config.Size.Width),
		float64(config.Size.Height),
	)

	return player, nil
}

// Update implements Entity interface
func (p *Player) Update() {
	// Use the actual screen dimensions from config
	centerX := float64(p.config.Size.Width) / HalfDivisor
	centerY := float64(p.config.Size.Height) / HalfDivisor

	// Get current position
	pos := p.coords.CalculateCircularPosition(p.posAngle)

	// Calculate angle to face center
	dx := centerX - pos.X
	dy := centerY - pos.Y
	facingAngle := math.Atan2(dy, dx) * RadiansToDegrees
	p.SetFacingAngle(common.Angle(facingAngle))

	// Update collision shape
	p.shape.SetPosition(pos.X-float64(DefaultPlayerSize)/HalfDivisor, pos.Y-float64(DefaultPlayerSize)/HalfDivisor)

	// Log position periodically
	if time.Since(p.lastLog) >= p.logInterval {
		p.logger.Info("Player position updated",
			"x", pos.X,
			"y", pos.Y,
			"angle", float64(p.posAngle),
			"facing_angle", float64(p.facingAngle),
			"center_x", centerX,
			"center_y", centerY,
		)
		p.lastLog = time.Now()
	}
}

// Draw implements the Drawable interface
func (p *Player) Draw(screen any, op any) {
	if p.sprite != nil {
		// Create draw options if none provided
		drawOp := &ebiten.DrawImageOptions{}
		if ebitenOp, ok := op.(*ebiten.DrawImageOptions); ok {
			drawOp = ebitenOp
		}

		// Get current position
		pos := p.GetPosition()

		// Center the sprite on its position
		drawOp.GeoM.Translate(-float64(DefaultPlayerSize)/HalfDivisor, -float64(DefaultPlayerSize)/HalfDivisor)

		// Set rotation based on facing angle
		rotationAngle := float64(p.GetFacingAngle()) * DegreesToRadians
		drawOp.GeoM.Rotate(rotationAngle)

		// Move to final position
		drawOp.GeoM.Translate(pos.X, pos.Y)

		// Draw with debug info
		p.logger.Debug("Drawing player",
			"x", pos.X,
			"y", pos.Y,
			"rotation", rotationAngle,
			"facing_angle", float64(p.facingAngle),
		)

		p.sprite.Draw(screen, drawOp)
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
