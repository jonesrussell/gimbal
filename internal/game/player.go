package game

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	input "github.com/quasilyte/ebitengine-input"
)

type Player struct {
	input     *input.Handler
	angle     float64
	speed     float64 // Add a speed variable to the player struct
	direction float64 // Add a direction variable to the player struct
}

func NewPlayer(input *input.Handler, speed float64) *Player {
	return &Player{
		input: input,
		speed: speed,
		angle: -math.Pi / 2, // Initialize the angle to -math.Pi / 2 to start at the bottom
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
}

func (p *Player) Draw(screen *ebiten.Image) {
	x := center.X + int(radius*math.Cos(p.angle))
	y := center.Y - int(radius*math.Sin(p.angle))
	ebitenutil.DebugPrintAt(screen, "PlayeR", x, y)
}
