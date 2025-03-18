package player

import (
	"errors"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/logger"
	"github.com/jonesrussell/gimbal/internal/physics"
	"github.com/solarlune/resolv"
	"go.uber.org/zap"
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

	logger.GlobalLogger.Debug("Creating new player",
		zap.Any("position", map[string]float64{
			"x": config.Position.X,
			"y": config.Position.Y,
		}),
		zap.Any("size", map[string]int{
			"width":  config.Size.Width,
			"height": config.Size.Height,
		}),
		zap.Float64("speed", config.Speed),
		zap.Float64("radius", config.Radius),
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
	// Calculate center of the screen using the radius as reference for screen size
	screenHeight := p.config.Radius * 3       // Since radius is height/3
	screenWidth := screenHeight * (4.0 / 3.0) // Standard 4:3 aspect ratio
	centerX := screenWidth / HalfDivisor
	centerY := screenHeight / HalfDivisor

	// Get current position
	pos := p.coords.CalculateCircularPosition(p.posAngle)

	// Calculate angle to face center
	dx := centerX - pos.X
	dy := centerY - pos.Y
	facingAngle := math.Atan2(dy, dx) * 180 / math.Pi
	p.SetFacingAngle(common.Angle(facingAngle))

	// Update collision shape
	p.shape.SetPosition(pos.X-16, pos.Y-16) // Center the 32x32 collision box

	// Log position periodically
	if time.Since(p.lastLog) >= p.logInterval {
		logger.GlobalLogger.Debug("Player position updated",
			zap.Float64("x", pos.X),
			zap.Float64("y", pos.Y),
			zap.Float64("angle", float64(p.posAngle)),
			zap.Float64("facing_angle", float64(p.facingAngle)),
			zap.Float64("center_x", centerX),
			zap.Float64("center_y", centerY),
			zap.Float64("screen_width", screenWidth),
			zap.Float64("screen_height", screenHeight),
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

		// Set rotation based on facing angle
		// Add 90 degrees to make sprite face upward by default
		rotationAngle := float64(p.GetFacingAngle())*math.Pi/180 + math.Pi/2
		drawOp.GeoM.Rotate(rotationAngle)

		// Set position after rotation
		pos := p.GetPosition()
		drawOp.GeoM.Translate(pos.X, pos.Y)

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
