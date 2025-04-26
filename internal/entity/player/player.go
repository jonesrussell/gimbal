package player

import (
	"errors"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jonesrussell/gimbal/internal/common"
)

const (
	// Movement constants
	HalfDivisor        = 2
	DegreesToRadians   = math.Pi / 180
	RadiansToDegrees   = 180 / math.Pi
	DefaultFacingAngle = 270

	// Logging constants
	LogIntervalSeconds   = 5
	PositionLogThreshold = 0.1 // Log position changes greater than this threshold
)

// Drawable interface defines the methods required for drawing
type Drawable interface {
	Draw(screen any, op any)
}

// PlayerInterface defines the complete set of player behaviors
type PlayerInterface interface {
	Drawable
	Update()
	GetPosition() common.Point
	SetPosition(pos common.Point) error
	GetSpeed() float64
	GetFacingAngle() common.Angle
	SetFacingAngle(angle common.Angle)
	GetAngle() common.Angle
	SetAngle(angle common.Angle) error
	GetBounds() common.Size
	Config() *common.EntityConfig
}

// Player represents the player entity in the game
type Player struct {
	position    common.Point
	config      *common.EntityConfig
	sprite      Drawable
	facingAngle common.Angle
	lastLog     time.Time
	logInterval time.Duration
	logger      common.Logger
	bounds      common.Size // Cached bounds for collision detection
}

// Ensure Player implements PlayerInterface at compile time
var _ PlayerInterface = (*Player)(nil)

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

	// Validate config
	if config.Size.Width <= 0 || config.Size.Height <= 0 {
		return nil, errors.New("invalid player size")
	}
	if config.Speed < 0 {
		return nil, errors.New("invalid player speed")
	}

	// Set initial angle based on mode
	initialAngle := common.Angle(DefaultFacingAngle)
	if config.Radius > 0 {
		initialAngle = 0 // Start at top of circle when in orbital mode
	}

	// Create player with initial position
	player := &Player{
		position:    config.Position,
		config:      config,
		sprite:      sprite,
		facingAngle: initialAngle,
		lastLog:     time.Now(),
		logInterval: time.Second * LogIntervalSeconds,
		logger:      logger,
		bounds:      config.Size,
	}

	// Set initial orbital position if radius is set
	if config.Radius > 0 {
		if err := player.updateOrbitalPosition(); err != nil {
			return nil, err
		}
	}

	player.logger.Debug("Player initialization complete",
		"initial_position", player.position,
		"facing_angle", float64(player.facingAngle),
		"size", player.config.Size,
		"log_interval", player.logInterval.Seconds(),
	)

	return player, nil
}

// Update implements PlayerInterface
func (p *Player) Update() {
	// Update orbital position if in orbital mode
	if p.config.Radius > 0 {
		if err := p.updateOrbitalPosition(); err != nil {
			p.logger.Error("Failed to update orbital position", "error", err)
		}
	}

	// Log state periodically
	if time.Since(p.lastLog) >= p.logInterval {
		p.logger.Debug("Player state",
			"position", p.position,
			"facing_angle", float64(p.facingAngle),
		)
		p.lastLog = time.Now()
	}
}

// updateOrbitalPosition updates the player's position based on orbital movement
func (p *Player) updateOrbitalPosition() error {
	angleRad := float64(p.facingAngle) * DegreesToRadians
	center := p.config.Position
	radius := p.config.Radius

	// Calculate new position on circle
	newPos := common.Point{
		X: center.X + radius*math.Sin(angleRad),
		Y: center.Y - radius*math.Cos(angleRad), // Subtract because Y increases downward
	}

	// Only update and log if position changed significantly
	if math.Abs(newPos.X-p.position.X) > PositionLogThreshold ||
		math.Abs(newPos.Y-p.position.Y) > PositionLogThreshold {
		p.logger.Debug("Player orbital position changed",
			"old_position", p.position,
			"new_position", newPos,
			"angle", float64(p.facingAngle),
			"angle_rad", angleRad,
		)
		p.position = newPos
	}

	return nil
}

// Draw implements PlayerInterface
func (p *Player) Draw(screen, op any) {
	if p.sprite == nil {
		return
	}

	// Create draw options if none provided
	drawOp := &ebiten.DrawImageOptions{}
	if ebitenOp, ok := op.(*ebiten.DrawImageOptions); ok {
		drawOp = ebitenOp
	}

	// Apply transformations in order:
	// 1. Center the sprite
	centerOffsetX := -float64(p.config.Size.Width) / HalfDivisor
	centerOffsetY := -float64(p.config.Size.Height) / HalfDivisor
	drawOp.GeoM.Translate(centerOffsetX, centerOffsetY)

	// 2. Rotate
	rotationAngle := float64(p.facingAngle) * DegreesToRadians
	drawOp.GeoM.Rotate(rotationAngle)

	// 3. Move to final position
	drawOp.GeoM.Translate(p.position.X, p.position.Y)

	p.sprite.Draw(screen, drawOp)
}

// GetPosition implements PlayerInterface
func (p *Player) GetPosition() common.Point {
	return p.position
}

// SetPosition implements PlayerInterface
func (p *Player) SetPosition(pos common.Point) error {
	// Only update position if we're not in orbital mode
	if p.config.Radius > 0 {
		return errors.New("cannot set position directly in orbital mode")
	}

	// Validate position is within bounds
	if pos.X < 0 || pos.Y < 0 {
		return errors.New("position cannot be negative")
	}

	// Only update and log if position changed significantly
	if math.Abs(pos.X-p.position.X) > PositionLogThreshold ||
		math.Abs(pos.Y-p.position.Y) > PositionLogThreshold {
		p.logger.Debug("Player position changed",
			"old_position", p.position,
			"new_position", pos,
			"facing_angle", float64(p.facingAngle),
		)
		p.position = pos
	}

	return nil
}

// GetSpeed implements PlayerInterface
func (p *Player) GetSpeed() float64 {
	return p.config.Speed
}

// GetAngle implements PlayerInterface
func (p *Player) GetAngle() common.Angle {
	return p.facingAngle
}

// SetAngle implements PlayerInterface
func (p *Player) SetAngle(angle common.Angle) error {
	oldAngle := p.facingAngle
	p.facingAngle = angle.Normalize()

	// If we have a radius, update orbital position
	if p.config.Radius > 0 {
		if err := p.updateOrbitalPosition(); err != nil {
			return err
		}
	}

	// Log angle change
	if math.Abs(float64(oldAngle-p.facingAngle)) > PositionLogThreshold {
		p.logger.Debug("Player angle set",
			"old_angle", float64(oldAngle),
			"new_angle", float64(p.facingAngle),
			"position", p.position,
		)
	}

	return nil
}

// GetFacingAngle implements PlayerInterface
func (p *Player) GetFacingAngle() common.Angle {
	return p.facingAngle
}

// SetFacingAngle implements PlayerInterface
func (p *Player) SetFacingAngle(angle common.Angle) {
	oldAngle := p.facingAngle
	p.facingAngle = angle.Normalize()

	// Log angle change if significant
	if math.Abs(float64(oldAngle-p.facingAngle)) > PositionLogThreshold {
		p.logger.Debug("Player facing angle changed",
			"old_angle", float64(oldAngle),
			"new_angle", float64(p.facingAngle),
		)
	}
}

// GetBounds implements PlayerInterface
func (p *Player) GetBounds() common.Size {
	return p.bounds
}

// Config implements PlayerInterface
func (p *Player) Config() *common.EntityConfig {
	return p.config
}

// Sprite returns the player's sprite
func (p *Player) Sprite() Drawable {
	return p.sprite
}
