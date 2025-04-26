package player

import (
	"errors"
	"image"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/entity/orbital"
)

const (
	// Movement constants
	HalfDivisor        = 2
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
	orbital     *orbital.Calculator
	position    orbital.Position
	config      *common.EntityConfig
	sprite      Drawable
	lastLog     time.Time
	logInterval time.Duration
	logger      common.Logger
	bounds      common.Size
}

// Ensure Player implements PlayerInterface at compile time
var _ PlayerInterface = (*Player)(nil)

// New creates a new player instance
func New(config *common.EntityConfig, sprite Drawable, logger common.Logger) (*Player, error) {
	if err := validateConfig(config, sprite, logger); err != nil {
		return nil, err
	}

	// Create orbital calculator
	orbCalc := orbital.NewCalculator(orbital.Config{
		Center: config.Position,
		Radius: config.Radius,
	})

	// Create player with initial position
	player := &Player{
		orbital:     orbCalc,
		config:      config,
		sprite:      sprite,
		lastLog:     time.Now(),
		logInterval: time.Second * LogIntervalSeconds,
		logger:      logger,
		bounds:      config.Size,
	}

	// Initialize position
	player.position = orbCalc.NewPosition(180, 0) // Start at bottom, facing up

	player.logger.Debug("Player initialization complete",
		"position", player.position.Point,
		"orbital_angle", float64(player.position.Orbital),
		"facing_angle", float64(player.position.Facing),
		"size", player.config.Size,
		"log_interval", player.logInterval.Seconds(),
	)

	return player, nil
}

func validateConfig(config *common.EntityConfig, sprite Drawable, logger common.Logger) error {
	if config == nil {
		return errors.New("config cannot be nil")
	}
	if sprite == nil {
		return errors.New("sprite cannot be nil")
	}
	if logger == nil {
		return errors.New("logger cannot be nil")
	}
	if config.Size.Width <= 0 || config.Size.Height <= 0 {
		return errors.New("invalid player size")
	}
	if config.Speed < 0 {
		return errors.New("invalid player speed")
	}
	return nil
}

// Update implements PlayerInterface
func (p *Player) Update() {
	// Log state periodically
	if time.Since(p.lastLog) >= p.logInterval {
		p.logger.Debug("Player state",
			"position", p.position.Point,
			"orbital_angle", float64(p.position.Orbital),
			"facing_angle", float64(p.position.Facing),
		)
		p.lastLog = time.Now()
	}
}

// Draw implements PlayerInterface
func (p *Player) Draw(screen, op any) {
	if p.sprite == nil {
		return
	}

	drawOp := createDrawOptions(op)
	applyTransformations(p, drawOp)
	p.sprite.Draw(screen, drawOp)
}

func createDrawOptions(op any) *ebiten.DrawImageOptions {
	if ebitenOp, ok := op.(*ebiten.DrawImageOptions); ok {
		return ebitenOp
	}
	return &ebiten.DrawImageOptions{}
}

func applyTransformations(p *Player, drawOp *ebiten.DrawImageOptions) {
	// 1. Scale the sprite to match the configured size
	if img, ok := p.sprite.(interface{ Bounds() image.Rectangle }); ok {
		bounds := img.Bounds()
		scaleX := float64(p.config.Size.Width) / float64(bounds.Dx())
		scaleY := float64(p.config.Size.Height) / float64(bounds.Dy())
		drawOp.GeoM.Scale(scaleX, scaleY)
	}

	// 2. Move to the center of the sprite (for rotation origin)
	centerOffsetX := float64(p.config.Size.Width) / HalfDivisor
	centerOffsetY := float64(p.config.Size.Height) / HalfDivisor
	drawOp.GeoM.Translate(-centerOffsetX, -centerOffsetY)

	// 3. Rotate around the center
	rotationAngle := float64(p.position.Facing) * orbital.DegreesToRadians
	drawOp.GeoM.Rotate(rotationAngle)

	// 4. Move to final position
	drawOp.GeoM.Translate(p.position.Point.X, p.position.Point.Y)
}

// GetPosition implements PlayerInterface
func (p *Player) GetPosition() common.Point {
	return p.position.Point
}

// SetPosition implements PlayerInterface
func (p *Player) SetPosition(pos common.Point) error {
	return errors.New("cannot set position directly in orbital mode")
}

// GetSpeed implements PlayerInterface
func (p *Player) GetSpeed() float64 {
	return p.config.Speed
}

// GetAngle implements PlayerInterface
func (p *Player) GetAngle() common.Angle {
	return p.position.Orbital
}

// SetAngle implements PlayerInterface
func (p *Player) SetAngle(angle common.Angle) error {
	p.orbital.UpdatePosition(&p.position, angle, p.position.Facing)
	return nil
}

// GetFacingAngle implements PlayerInterface
func (p *Player) GetFacingAngle() common.Angle {
	return p.position.Facing
}

// SetFacingAngle implements PlayerInterface
func (p *Player) SetFacingAngle(angle common.Angle) {
	p.orbital.UpdatePosition(&p.position, p.position.Orbital, angle)
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
