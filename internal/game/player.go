package game

import (
	"errors"
	"log"
	"math"

	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/solarlune/resolv"
)

// Player represents a player in the game.
type Player struct {
	// Input is the input handler for the player.
	input InputHandlerInterface
	// Angle is the current angle of the player.
	angle float64
	// Speed is the speed of the player.
	speed float64
	// Direction is the current direction of the player.
	direction float64
	// Object is the game object representing the player.
	Object *resolv.Object
	// Sprite is the player's sprite.
	Sprite *ebiten.Image
}

const (
	initialAngle = -math.Pi / 2
	playerLabel  = "PlayeR"
)

func NewPlayer(input InputHandlerInterface, speed float64) (*Player, error) {
	if input == nil {
		return nil, errors.New("input handler cannot be nil")
	}
	if radius <= 0 {
		return nil, errors.New("radius must be greater than zero")
	}

	x := center.X + int(radius*math.Cos(initialAngle))
	y := center.Y - int(radius*math.Sin(initialAngle))
	width := playerWidth // replace with your player's width
	height := playerHeight

	// Load the sprite.
	spriteImage, _, err := ebitenutil.NewImageFromFile("assets/player.png")
	if err != nil {
		log.Fatal(err)
	}

	return &Player{
		input:  input,
		speed:  speed,
		angle:  initialAngle,
		Object: resolv.NewObject(float64(x), float64(y), float64(width), float64(height)),
		Sprite: spriteImage,
	}, nil
}

func (player *Player) Update() {
	if player.input.IsKeyPressed(ebiten.KeyLeft) {
		player.direction = -1
	} else if player.input.IsKeyPressed(ebiten.KeyRight) {
		player.direction = 1
	} else {
		player.direction = 0
	}

	player.angle += player.direction * player.speed

	x := center.X + int(radius*math.Cos(player.angle))
	y := center.Y - int(radius*math.Sin(player.angle))
	player.Object.Position.X = float64(x)
	player.Object.Position.Y = float64(y)
	player.Object.Update()
}

func (p *Player) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	// Scale the sprite to half its original size.
	op.GeoM.Scale(0.1, 0.1)
	op.GeoM.Translate(p.Object.Position.X, p.Object.Position.Y)
	p.Object.Update()
	screen.DrawImage(p.Sprite, op)
}
