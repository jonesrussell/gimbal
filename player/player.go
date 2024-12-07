package player

import (
	"errors"
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
)

type PlayerInput struct {
	input InputHandlerInterface
}

type PlayerPosition struct {
	Object *resolv.ConvexPolygon
}

type PlayerSprite struct {
	Sprite *ebiten.Image
}

type Position struct {
	X, Y float64
}

type Player struct {
	PlayerInput
	PlayerPosition
	PlayerSprite
	viewAngle float64
	angle     float64
	Speed     float64
	center    image.Point
}

// NewPlayer creates a new instance of a player with the given input handler, speed, and sprite image.
// If any of the arguments are nil, an error is returned.
func NewPlayer(input InputHandlerInterface, speed float64, spriteImage image.Image, screenCenter image.Point) (*Player, error) {
	if input == nil {
		return nil, errors.New("input handler cannot be nil")
	}

	if speed <= 0 {
		return nil, errors.New("speed must be greater than zero")
	}

	playerImage := ebiten.NewImageFromImage(spriteImage)
	player := &Player{
		PlayerInput: PlayerInput{input: input},
		PlayerPosition: PlayerPosition{
			Object: resolv.NewConvexPolygon(
				float64(screenCenter.X),
				float64(screenCenter.Y),
				[]float64{
					-16, -16,
					16, -16,
					16, 16,
					-16, 16,
				},
			),
		},
		PlayerSprite: PlayerSprite{Sprite: playerImage},
		Speed:        speed,
		center:       screenCenter,
	}

	return player, nil
}

func (p *Player) Update() {
	// Update orbital position
	if p.input.IsKeyPressed(ebiten.KeyLeft) {
		p.angle -= 0.05
	}
	if p.input.IsKeyPressed(ebiten.KeyRight) {
		p.angle += 0.05
	}

	// Calculate new position
	newX := float64(p.center.X) + math.Cos(p.angle)*radius
	newY := float64(p.center.Y) + math.Sin(p.angle)*radius

	// Calculate rotation angle to face center
	// Add π to rotate 180 degrees (sprite faces inward)
	p.viewAngle = math.Atan2(
		float64(p.center.Y)-newY,
		float64(p.center.X)-newX,
	) - math.Pi/2 + math.Pi

	// Update position
	p.PlayerPosition.Object.SetPosition(newX, newY)
}

func (p *Player) Draw(screen *ebiten.Image) {
	p.drawSprite(screen)
}

func (p *Player) drawSprite(screen *ebiten.Image) {
	if p.Sprite != nil {
		spriteOp := p.createSpriteOptions()
		rotatedSprite := p.getRotatedSprite()
		screen.DrawImage(rotatedSprite, spriteOp)
	}
}

func (p *Player) createSpriteOptions() *ebiten.DrawImageOptions {
	width := p.Sprite.Bounds().Dx()
	height := float64(p.Sprite.Bounds().Dy())
	spriteOp := &ebiten.DrawImageOptions{}

	// Center the sprite
	spriteOp.GeoM.Translate(-float64(width)/2, -height/2)

	// Scale the sprite to 1/10th size
	spriteOp.GeoM.Scale(0.1, 0.1)

	// Use viewAngle for rotation
	spriteOp.GeoM.Rotate(p.viewAngle)

	// Move to player position
	pos := p.Object.Position()
	spriteOp.GeoM.Translate(pos.X, pos.Y)

	return spriteOp
}

func (p *Player) getRotatedSprite() *ebiten.Image {
	return p.Sprite.SubImage(
		image.Rect(0, 0, p.Sprite.Bounds().Dx(), p.Sprite.Bounds().Dy()),
	).(*ebiten.Image)
}
