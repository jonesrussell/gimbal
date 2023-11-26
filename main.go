package main

import (
	"image"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	input "github.com/quasilyte/ebitengine-input"
)

const (
	ActionMoveLeft input.Action = iota
	ActionMoveRight
)

var (
	width, height = 640, 480
	radius        = 200.0
	center        = image.Point{X: width / 2, Y: height / 2}
)

func main() {
	ebiten.SetWindowSize(width, height)
	if err := ebiten.RunGame(newExampleGame()); err != nil {
		log.Fatal(err)
	}
}

type exampleGame struct {
	p           *player
	inputSystem input.System
}

func newExampleGame() *exampleGame {
	g := &exampleGame{}
	g.inputSystem.Init(input.SystemConfig{
		DevicesEnabled: input.AnyDevice,
	})
	keymap := input.Keymap{
		ActionMoveLeft:  {input.KeyGamepadLeft, input.KeyLeft, input.KeyA},
		ActionMoveRight: {input.KeyGamepadRight, input.KeyRight, input.KeyD},
	}
	g.p = &player{
		input: g.inputSystem.NewHandler(0, keymap),
		angle: math.Pi / 2,
	}
	return g
}

func (g *exampleGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return width, height
}

func (g *exampleGame) Draw(screen *ebiten.Image) {
	g.p.Draw(screen)
}

func (g *exampleGame) Update() error {
	g.inputSystem.Update()
	g.p.Update()
	return nil
}

type player struct {
	input *input.Handler
	angle float64
	speed float64
}

func (p *player) Update() {
	if p.input.ActionIsPressed(ActionMoveLeft) {
		p.speed = -0.01
	} else if p.input.ActionIsPressed(ActionMoveRight) {
		p.speed = 0.01
	} else {
		p.speed = 0
	}
	p.angle += p.speed
}

func (p *player) Draw(screen *ebiten.Image) {
	x := center.X + int(radius*math.Cos(p.angle))
	y := center.Y - int(radius*math.Sin(p.angle))
	ebitenutil.DebugPrintAt(screen, "PlayeR", x, y)
}
