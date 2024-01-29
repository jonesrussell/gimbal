package game

import (
	"fmt"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	width, height = 640, 480
	radius        = 200.0
	center        = image.Point{X: width / 2, Y: height / 2}
	debug         = true
)

type GimlarGame struct {
	player *Player
	speed  float64
}

func NewGimlarGame(speed float64) (*GimlarGame, error) {
	g := &GimlarGame{
		speed: speed,
	}

	handler := &InputHandler{}
	var err error
	g.player, err = NewPlayer(handler, g.speed)
	if err != nil {
		return nil, err
	}
	return g, nil
}

func (g *GimlarGame) Run() error {
	ebiten.SetWindowSize(width, height) // Sets the window size to the width and height defined in your game.
	return ebiten.RunGame(g)            // Runs your game. The game's logic is executed in the Update method of your GimlarGame struct.
}

func (g *GimlarGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return width, height
}

func (g *GimlarGame) Update() error {
	// Update the player's state
	g.player.Update()

	// Add any other game state updates here

	return nil
}

func (g *GimlarGame) Draw(screen *ebiten.Image) {
	g.player.Draw(screen)

	// Draw debug info if debug is true
	if debug {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", ebiten.ActualFPS())) // Print the current FPS
		// Draw grid overlay
		for i := 0; i < width; i += 40 {
			vector.StrokeLine(screen, float32(i), 0, float32(i), float32(height), 1, color.White, false)
		}
		for i := 0; i < height; i += 40 {
			vector.StrokeLine(screen, 0, float32(i), float32(width), float32(i), 1, color.White, false)
		}
	}
}
