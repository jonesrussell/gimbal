package main

import (
	"image"
	"image/color"
	"log"
	"math"
	"os"
	"strconv"

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
	radius        = 200.0 // Change radius to float64
	center        = image.Point{X: width / 2, Y: height / 2}
	debug         = false
)

func main() {
	// Check if DEBUG environment variable is set to true
	if value, exists := os.LookupEnv("DEBUG"); exists {
		debug, _ = strconv.ParseBool(value)
	}

	ebiten.SetWindowSize(width, height)
	if err := ebiten.RunGame(newExampleGame(0.02)); err != nil { // Pass the speed as an argument
		log.Fatal(err)
	}
}

type exampleGame struct {
	p           *player
	inputSystem input.System
	speed       float64 // Add a speed variable to the exampleGame struct
}

func newExampleGame(speed float64) *exampleGame { // Take speed as an argument
	g := &exampleGame{speed: speed} // Initialize the speed variable
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
		speed: g.speed, // Pass the speed to the player
	}
	return g
}

func (g *exampleGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return width, height
}

func (g *exampleGame) Draw(screen *ebiten.Image) {
	g.p.Draw(screen)

	// Draw debug info if debug is true
	if debug {
		ebitenutil.DebugPrint(screen, "Debugging Info")
		// Draw grid overlay
		for i := 0; i < width; i += 40 {
			ebitenutil.DrawLine(screen, float64(i), 0, float64(i), float64(height), color.White)
		}
		for i := 0; i < height; i += 40 {
			ebitenutil.DrawLine(screen, 0, float64(i), float64(width), float64(i), color.White)
		}
	}
}

func (g *exampleGame) Update() error {
	g.inputSystem.Update()
	g.p.Update()
	return nil
}

type player struct {
	input     *input.Handler
	angle     float64
	speed     float64 // Add a speed variable to the player struct
	direction float64 // Add a direction variable to the player struct
}

func (p *player) Update() {
	if p.input.ActionIsPressed(ActionMoveLeft) {
		p.direction = -1
	} else if p.input.ActionIsPressed(ActionMoveRight) {
		p.direction = 1
	} else {
		p.direction = 0
	}
	p.angle += p.direction * p.speed
}

func (p *player) Draw(screen *ebiten.Image) {
	x := center.X + int(radius*math.Cos(p.angle))
	y := center.Y - int(radius*math.Sin(p.angle))
	ebitenutil.DebugPrintAt(screen, "PlayeR", x, y)
}
