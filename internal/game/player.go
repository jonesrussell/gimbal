package game

import (
	"errors"
	"image"
	"math"

	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
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
	// Debugger
	debug *Debugger
}

func NewPlayer(input InputHandlerInterface, speed float64, debugger *Debugger, spriteImage *ebiten.Image) (*Player, error) {
	if input == nil {
		return nil, errors.New("input handler cannot be nil")
	}

	if speed <= 0 {
		return nil, errors.New("speed must be greater than zero")
	}

	if spriteImage == nil {
		return nil, errors.New("sprite image cannot be nil")
	}

	center := Center()
	x := center.X + int(radius*math.Cos(math.Pi/2))
	y := center.Y - int(radius*math.Sin(math.Pi/2))

	return &Player{
		input:       input,
		speed:       speed,
		angle:       math.Pi / 2,
		orientation: 0.0,
		Object:      resolv.NewObject(float64(x), float64(y), float64(playerWidth), float64(playerHeight)),
		Sprite:      spriteImage,
		gameStarted: false,
		debug:       debugger,
	}, nil
}

func (player *Player) Update() {
	if !player.gameStarted {
		player.debug.DebugPlayer(player)
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
		player.debug.DebugPrintOrientation(player.orientation)
	}

	if player.direction != oldDirection {
		player.debug.DebugPrintDirection(player.direction)
	}

	if player.angle != oldAngle {
		player.debug.DebugPrintAngle(player.angle)
	}

	if player.Object.Position.X != oldX || player.Object.Position.Y != oldY {
		pos := image.Point{X: int(player.Object.Position.X), Y: int(player.Object.Position.Y)}
		player.debug.DebugPrintPosition(pos)
	}

	player.Object.Update()
}

func (player *Player) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(player.Object.Position.X, player.Object.Position.Y)
	player.Object.Update()

	if player.Sprite != nil {
		// Create a separate DrawImageOptions for the sprite rotation
		spriteOp := &ebiten.DrawImageOptions{}
		// Scale the sprite to half its original size.
		spriteOp.GeoM.Scale(0.1, 0.1)
		spriteOp.GeoM.Rotate(player.angle)
		// Translate the rotated sprite to the player's position
		spriteOp.GeoM.Translate(player.Object.Position.X, player.Object.Position.Y)
		rotatedSprite := player.Sprite.SubImage(image.Rect(0, 0, player.Sprite.Bounds().Dx(), player.Sprite.Bounds().Dy())).(*ebiten.Image)
		screen.DrawImage(rotatedSprite, spriteOp)
	}
}
