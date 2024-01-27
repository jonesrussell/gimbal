package game

import (
	"errors"
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	input "github.com/quasilyte/ebitengine-input"
	"github.com/solarlune/resolv"
)

// Define the interface
type HandlerInterface interface {
	ActionIsPressed(action input.Action) bool
}

type HandlerWrapper struct {
	HandlerInterface
}

func (hw *HandlerWrapper) ActionIsPressed(action input.Action) bool {
	return hw.HandlerInterface.ActionIsPressed(action)
}

// Player represents a player in the game.
type Player struct {
	// Input is the input handler for the player.
	input HandlerInterface
	// Angle is the current angle of the player.
	angle float64
	// Speed is the speed of the player.
	speed float64
	// Direction is the current direction of the player.
	direction float64
	// Object is the game object representing the player.
	Object *resolv.Object
}

const (
	initialAngle = -math.Pi / 2
	playerLabel  = "PlayeR"
)

func NewPlayer(input HandlerInterface, speed float64) (*Player, error) {
	if input == nil {
		return nil, errors.New("input handler cannot be nil")
	}
	if radius <= 0 {
		return nil, errors.New("radius must be greater than zero")
	}

	x := center.X + int(radius*math.Cos(initialAngle))
	y := center.Y - int(radius*math.Sin(initialAngle))
	width := 20  // replace with your player's width
	height := 20 // replace with your player's height

	handlerWrapper := &HandlerWrapper{HandlerInterface: input}

	return &Player{
		input:  handlerWrapper,
		speed:  speed,
		angle:  -math.Pi / 2, // Initialize the angle to -math.Pi / 2 to start at the bottom
		Object: resolv.NewObject(float64(x), float64(y), float64(width), float64(height)),
	}, nil
}

func (p *Player) Update() {
	if p.input.ActionIsPressed(ActionMoveLeft) {
		p.direction = -1
	} else if p.input.ActionIsPressed(ActionMoveRight) {
		p.direction = 1
	} else {
		p.direction = 0
	}
	fmt.Println("Direction after Update: ", p.direction) // Debugging print statement

	p.angle += p.direction * p.speed

	x := center.X + int(radius*math.Cos(p.angle))
	y := center.Y - int(radius*math.Sin(p.angle))
	p.Object.X = float64(x)
	p.Object.Y = float64(y)
}

func (p *Player) Draw(screen *ebiten.Image) {
	x := center.X + int(radius*math.Cos(p.angle))
	y := center.Y - int(radius*math.Sin(p.angle))
	ebitenutil.DebugPrintAt(screen, playerLabel, x, y)
}
