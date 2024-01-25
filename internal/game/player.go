package game

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	input "github.com/quasilyte/ebitengine-input"
	"github.com/solarlune/resolv"
)

type Player struct {
	input     *input.Handler
	angle     float64
	speed     float64 // Add a speed variable to the player struct
	direction float64 // Add a direction variable to the player struct
	Object    *resolv.Object
}

func NewPlayer(input *input.Handler, speed float64) *Player {
	x := center.X + int(radius*math.Cos(-math.Pi/2))
	y := center.Y - int(radius*math.Sin(-math.Pi/2))
	width := 20  // replace with your player's width
	height := 20 // replace with your player's height

	return &Player{
		input:  input,
		speed:  speed,
		angle:  -math.Pi / 2, // Initialize the angle to -math.Pi / 2 to start at the bottom
		Object: resolv.NewObject(float64(x), float64(y), float64(width), float64(height)),
	}
}

func (p *Player) Update() {
	if p.input.ActionIsPressed(ActionMoveLeft) {
		p.direction = -1
	} else if p.input.ActionIsPressed(ActionMoveRight) {
		p.direction = 1
	} else {
		p.direction = 0
	}
	p.angle += p.direction * p.speed

	x := center.X + int(radius*math.Cos(p.angle))
	y := center.Y - int(radius*math.Sin(p.angle))
	p.Object.X = float64(x)
	p.Object.Y = float64(y)
}

func (p *Player) Draw(screen *ebiten.Image) {
	x := center.X + int(radius*math.Cos(p.angle))
	y := center.Y - int(radius*math.Sin(p.angle))
	ebitenutil.DebugPrintAt(screen, "PlayeR", x, y)
}
