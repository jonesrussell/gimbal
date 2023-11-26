package game

import (
	"fmt"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	input "github.com/quasilyte/ebitengine-input"
)

var (
	width, height = 640, 480
	radius        = 200.0
	center        = image.Point{X: width / 2, Y: height / 2}
	debug         = true
)

const (
	ActionMoveLeft input.Action = iota
	ActionMoveRight
)

type GimlarGame struct {
	p           *Player
	inputSystem input.System
	speed       float64
}

func NewGimlarGame(speed float64) *GimlarGame { // Take speed as an argument
	g := &GimlarGame{
		p:           &Player{},
		inputSystem: input.System{},
		speed:       speed,
	} // Initialize the speed variable
	g.inputSystem.Init(input.SystemConfig{
		DevicesEnabled: input.AnyDevice,
	})
	keymap := input.Keymap{
		ActionMoveLeft:  {input.KeyGamepadLeft, input.KeyLeft, input.KeyA},
		ActionMoveRight: {input.KeyGamepadRight, input.KeyRight, input.KeyD},
	}
	g.p = NewPlayer(g.inputSystem.NewHandler(0, keymap), g.speed) // Pass the speed to the player
	return g
}

func (g *GimlarGame) Run() error {
	ebiten.SetWindowSize(width, height)
	return ebiten.RunGame(g)
}

func (g *GimlarGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return width, height
}

func (g *GimlarGame) Draw(screen *ebiten.Image) {
	g.p.Draw(screen)

	// Draw debug info if debug is true
	if debug {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", ebiten.CurrentFPS())) // Print the current FPS
		// Draw grid overlay
		for i := 0; i < width; i += 40 {
			ebitenutil.DrawLine(screen, float64(i), 0, float64(i), float64(height), color.White)
		}
		for i := 0; i < height; i += 40 {
			ebitenutil.DrawLine(screen, 0, float64(i), float64(width), float64(i), color.White)
		}
	}
}

func (g *GimlarGame) Update() error {
	g.inputSystem.Update()
	g.p.Update()
	return nil
}
