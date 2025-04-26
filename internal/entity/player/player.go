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
	// DefaultFacingAngle is the default angle the player faces (upward)
	DefaultFacingAngle = 270
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
	SetPosition(pos common.Point)
	GetSpeed() float64
	GetFacingAngle() common.Angle
	SetFacingAngle(angle common.Angle)
	GetAngle() common.Angle
	SetAngle(angle common.Angle)
	GetBounds() common.Size
	Config() *common.EntityConfig
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

	logger.Debug("Creating new player with config",
		"position", config.Position,
		"size", config.Size,
		"speed", config.Speed,
		"radius", config.Radius,
	)

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
	}

	// Set initial orbital position if radius is set
	if config.Radius > 0 {
		angleRad := float64(player.facingAngle) * DegreesToRadians
		center := config.Position
		radius := config.Radius
		player.position = common.Point{
			X: center.X + radius*math.Sin(angleRad),
			Y: center.Y - radius*math.Cos(angleRad),
		}
	}

	logger.Debug("Player initialization complete",
		"initial_position", player.position,
		"facing_angle", float64(player.facingAngle),
		"size", player.config.Size,
		"log_interval", player.logInterval.Seconds(),
	)

	return player, nil
}

// Update implements PlayerInterface
func (p *Player) Update() {
	// Calculate orbital position if radius is set
	if p.config.Radius > 0 {
		angleRad := float64(p.facingAngle) * DegreesToRadians
		center := p.config.Position
		radius := p.config.Radius

		// Calculate position on circle
		p.position = common.Point{
			X: center.X + radius*math.Sin(angleRad),
			Y: center.Y - radius*math.Cos(angleRad), // Subtract because Y increases downward in screen coordinates
		}
	}

	// Log position periodically
	if time.Since(p.lastLog) >= p.logInterval {
		p.logger.Debug("Player state",
			"position", p.position,
			"facing_angle", float64(p.facingAngle),
			"radius", p.config.Radius,
		)
		p.lastLog = time.Now()
	}
}

// Draw implements PlayerInterface
func (p *Player) Draw(screen, op any) {
	if p.sprite != nil {
		// Create draw options if none provided
		drawOp := &ebiten.DrawImageOptions{}
		if ebitenOp, ok := op.(*ebiten.DrawImageOptions); ok {
			drawOp = ebitenOp
		}

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
			"transform", map[string]any{
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

// GetPosition implements PlayerInterface
func (p *Player) GetPosition() common.Point {
	return p.position
}

// SetPosition implements PlayerInterface
// This is used for direct movement (left/right controls)
func (p *Player) SetPosition(pos common.Point) {
	// Only update position if we're not in orbital mode
	if p.config.Radius == 0 {
		p.position = pos
	}
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
func (p *Player) SetAngle(angle common.Angle) {
	p.facingAngle = angle.Normalize()
	// If we have a radius, immediately update position for orbital movement
	if p.config.Radius > 0 {
		angleRad := float64(p.facingAngle) * DegreesToRadians
		center := p.config.Position
		radius := p.config.Radius
		p.position = common.Point{
			X: center.X + radius*math.Sin(angleRad),
			Y: center.Y - radius*math.Cos(angleRad),
		}
	}
	p.logger.Debug("Player angle set",
		"angle", float64(angle),
		"normalized_angle", float64(p.facingAngle),
		"position", p.position,
	)
}

// GetFacingAngle implements PlayerInterface
func (p *Player) GetFacingAngle() common.Angle {
	return p.facingAngle
}

// SetFacingAngle implements PlayerInterface
func (p *Player) SetFacingAngle(angle common.Angle) {
	p.facingAngle = angle.Normalize()
}

// GetBounds implements PlayerInterface
func (p *Player) GetBounds() common.Size {
	return p.config.Size
}

// Config implements PlayerInterface
func (p *Player) Config() *common.EntityConfig {
	return p.config
}

// Sprite returns the player's sprite
func (p *Player) Sprite() Drawable {
	return p.sprite
}
