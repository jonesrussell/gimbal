package player

import (
	"errors"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/logger"
	"github.com/jonesrussell/gimbal/internal/physics"
	"github.com/solarlune/resolv"
)

// Player represents the player entity in the game
type Player struct {
	coords *physics.CoordinateSystem
	config *common.EntityConfig
	sprite *ebiten.Image
	shape  resolv.IShape
	angle  common.Angle
	path   []resolv.Vector
}

// New creates a new player instance
func New(config *common.EntityConfig, sprite *ebiten.Image) (*Player, error) {
	if config == nil {
		return nil, errors.New("config cannot be nil")
	}

	coords := physics.NewCoordinateSystem(config.Position, config.Radius)

	// Create player collision shape as a rectangle
	shape := resolv.NewRectangle(
		config.Position.X,
		config.Position.Y,
		float64(config.Size.Width),
		float64(config.Size.Height),
	)

	// Start at the bottom of the screen (270 degrees)
	initialAngle := common.Angle(common.BottomAngle * common.DegreesToRadians)

	// Debug logging
	logger.GlobalLogger.Debug("Creating new player",
		"config", map[string]any{
			"position": config.Position,
			"size":     config.Size,
			"radius":   config.Radius,
			"speed":    config.Speed,
		},
		"sprite_size", map[string]any{
			"width":  sprite.Bounds().Dx(),
			"height": sprite.Bounds().Dy(),
		},
		"initial_angle", initialAngle,
	)

	return &Player{
		coords: coords,
		config: config,
		sprite: sprite,
		shape:  shape,
		angle:  initialAngle,
		path:   make([]resolv.Vector, 0),
	}, nil
}

// Update implements Entity interface
func (p *Player) Update() {
	pos := p.GetPosition()
	p.shape.SetPosition(pos.X, pos.Y)
}

// Draw implements Entity interface
func (p *Player) Draw(screen *ebiten.Image) {
	if p.sprite == nil {
		return
	}

	pos := p.GetPosition()
	op := &ebiten.DrawImageOptions{}

	// Calculate scaling to match configured size
	scaleX := float64(p.config.Size.Width) / float64(p.sprite.Bounds().Dx())
	scaleY := float64(p.config.Size.Height) / float64(p.sprite.Bounds().Dy())

	// Calculate offsets for rotation
	offsetX := -float64(p.config.Size.Width) / common.CenterDivisor
	offsetY := -float64(p.config.Size.Height) / common.CenterDivisor

	// Apply transformations
	op.GeoM.Translate(offsetX, offsetY)
	op.GeoM.Rotate(float64(p.angle.ToRadians()))
	op.GeoM.Scale(scaleX, scaleY)

	// Move to final position
	finalX := pos.X + float64(p.config.Size.Width)/common.CenterDivisor
	finalY := pos.Y + float64(p.config.Size.Height)/common.CenterDivisor
	op.GeoM.Translate(finalX, finalY)

	// Debug logging
	logger.GlobalLogger.Debug("Drawing player",
		"position", map[string]any{
			"x": pos.X,
			"y": pos.Y,
		},
		"angle", p.angle,
		"scale", map[string]any{
			"x": scaleX,
			"y": scaleY,
		},
		"final_position", map[string]any{
			"x": finalX,
			"y": finalY,
		},
	)

	screen.DrawImage(p.sprite, op)
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
	return p.config.Speed
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
