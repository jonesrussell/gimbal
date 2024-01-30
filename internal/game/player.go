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
	// Debugging check if game started
	gameStarted bool
	// Orientation of player's viewable sprite
	orientation float64
}

func NewPlayer(input InputHandlerInterface, speed float64, radius float64) (*Player, error) {
	if input == nil {
		return nil, errors.New("input handler cannot be nil")
	}

	if speed <= 0 {
		return nil, errors.New("speed must be greater than zero")
	}

	if radius <= 0 {
		return nil, errors.New("radius must be greater than zero")
	}

	center := GetCenter()
	x := center.X + int(radius*math.Cos(math.Pi/2))
	y := center.Y - int(radius*math.Sin(math.Pi/2))

	// Load the sprite.
	spriteImage, _, err := ebitenutil.NewImageFromFile("assets/player.png")
	if err != nil {
		log.Fatal(err)
	}

	return &Player{
		input:       input,
		speed:       speed,
		angle:       math.Pi / 2,
		orientation: 0.0,
		Object:      resolv.NewObject(float64(x), float64(y), float64(playerWidth), float64(playerHeight)),
		Sprite:      spriteImage,
		gameStarted: false,
	}, nil
}

func (player *Player) Update() {
	if !player.gameStarted {
		player.DebugPrint()
		player.gameStarted = true
	}

	oldOrientation := player.orientation
	oldDirection := player.direction
	oldAngle := player.angle
	oldX := player.Object.Position.X
	oldY := player.Object.Position.Y

	if player.input.IsKeyPressed(ebiten.KeyLeft) {
		player.direction = -1
		player.orientation -= 0.05
	} else if player.input.IsKeyPressed(ebiten.KeyRight) {
		player.direction = 1
		player.orientation += 0.05
	} else {
		player.direction = 0
	}

	// Calculate the x and y positions based on the current angle
	x := center.X + int(radius*math.Cos(player.orientation))
	y := center.Y - int(radius*math.Sin(player.orientation))
	player.Object.Position.X = float64(x)
	player.Object.Position.Y = float64(y)

	if player.orientation != oldOrientation {
		player.DebugPrintOrientation()
	}

	if player.direction != oldDirection {
		player.DebugPrintDirection()
	}

	if player.angle != oldAngle {
		player.DebugPrintAngle()
	}

	if player.Object.Position.X != oldX || player.Object.Position.Y != oldY {
		player.DebugPrintPosition()
	}

	player.Object.Update()
}

func (player *Player) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	// Scale the sprite to half its original size.
	op.GeoM.Scale(0.1, 0.1)
	op.GeoM.Rotate(player.angle)
	op.GeoM.Translate(player.Object.Position.X, player.Object.Position.Y)
	player.Object.Update()
	screen.DrawImage(player.Sprite, op)
}

func (player *Player) DebugPrint() {
	player.DebugPrintOrientation()
	player.DebugPrintDirection()
	player.DebugPrintAngle()
	player.DebugPrintPosition()
}

func (player *Player) DebugPrintOrientation() {
	log.Printf("Player orientation: %f", player.orientation)
}

func (player *Player) DebugPrintDirection() {
	log.Printf("Player direction: %f", player.direction)
}

func (player *Player) DebugPrintAngle() {
	log.Printf("Player angle: %f", player.angle)
}
func (player *Player) DebugPrintPosition() {
	log.Printf("Player position: (%f, %f)", player.Object.Position.X, player.Object.Position.Y)
}
