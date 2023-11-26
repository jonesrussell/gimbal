package main

import (
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	input "github.com/quasilyte/ebitengine-input"
)

const (
	ActionMoveLeft input.Action = iota
	ActionMoveRight
)

func main() {
	ebiten.SetWindowSize(640, 480)
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
		pos:   image.Point{X: 640 / 2, Y: 400},
	}
	return g
}

func (g *exampleGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 480
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
	pos   image.Point
}

func (p *player) Update() {
	if p.input.ActionIsPressed(ActionMoveLeft) {
		p.pos.X -= 4
	}
	if p.input.ActionIsPressed(ActionMoveRight) {
		p.pos.X += 4
	}
}

func (p *player) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(screen, "PlayeR", p.pos.X, p.pos.Y)
}
