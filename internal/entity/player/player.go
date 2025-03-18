package player

import (
	"errors"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jonesrussell/gimbal/internal/common"
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
	position    common.Point
	config      *common.EntityConfig
	sprite      Drawable
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

	logger.Debug("Creating new player with config",
		"position", config.Position,
		"size", config.Size,
		"speed", config.Speed,
		"radius", config.Radius,
	)

	// Create player with initial position
	player := &Player{
		position:    config.Position,
		config:      config,
		sprite:      sprite,
		facingAngle: common.Angle(270), // Face upward by default
		lastLog:     time.Now(),
		logInterval: time.Second * LogIntervalSeconds,
		logger:      logger,
	}

	logger.Debug("Player initialization complete",
		"initial_position", player.position,
		"facing_angle", float64(player.facingAngle),
		"size", player.config.Size,
		"log_interval", player.logInterval.Seconds(),
	)

	return player, nil
}

// Update implements Entity interface
func (p *Player) Update() {
	// Log position periodically only if it has changed
	if time.Since(p.lastLog) >= p.logInterval {
		p.logger.Debug("Player state",
			"position", p.position,
			"facing_angle", float64(p.facingAngle),
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

		// Log initial state before transformations
		p.logger.Debug("Player draw start",
			"initial_state", map[string]interface{}{
				"position":     p.position,
				"size":         p.config.Size,
				"facing_angle": float64(p.facingAngle),
			},
		)

		// Order of transformations:
		// 1. Center the sprite on its origin point
		centerOffsetX := -float64(p.config.Size.Width) / HalfDivisor
		centerOffsetY := -float64(p.config.Size.Height) / HalfDivisor
		drawOp.GeoM.Translate(centerOffsetX, centerOffsetY)

		// 2. Rotation based on facing angle
		rotationAngle := float64(p.facingAngle) * DegreesToRadians
		drawOp.GeoM.Rotate(rotationAngle)

		// 3. Move to final position
		drawOp.GeoM.Translate(p.position.X, p.position.Y)

		// Log transformation details
		p.logger.Debug("Player transformations",
			"transform", map[string]interface{}{
				"center_offset": map[string]float64{
					"x": centerOffsetX,
					"y": centerOffsetY,
				},
				"rotation_angle": rotationAngle,
				"final_position": map[string]float64{
					"x": p.position.X,
					"y": p.position.Y,
				},
			},
		)

		p.sprite.Draw(screen, drawOp)
	}
}

// GetPosition returns the player's current position
func (p *Player) GetPosition() common.Point {
	return p.position
}

// SetPosition sets the player's position
func (p *Player) SetPosition(pos common.Point) {
	p.position = pos
}

// GetSpeed returns the player's speed
func (p *Player) GetSpeed() float64 {
	return p.config.Speed
}

// GetFacingAngle returns the direction the player is facing
func (p *Player) GetFacingAngle() common.Angle {
	return p.facingAngle
}

// SetFacingAngle sets the direction the player is facing
func (p *Player) SetFacingAngle(angle common.Angle) {
	p.facingAngle = angle.Normalize()
}

// GetBounds returns the player's size
func (p *Player) GetBounds() common.Size {
	return p.config.Size
}

// Config returns the player's configuration
func (p *Player) Config() *common.EntityConfig {
	return p.config
}

// Sprite returns the player's sprite
func (p *Player) Sprite() Drawable {
	return p.sprite
}
