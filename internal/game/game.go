package game

import (
	"fmt"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/solarlune/resolv"
)

var (
	width, height             = 640, 480
	radius                    = float64(height / 2)
	center                    = image.Point{X: width / 2, Y: height / 2}
	gridSpacing               = 32
	playerWidth, playerHeight = 16, 16
)

type GimlarGame struct {
	player *Player
	speed  float64
	debug  bool
	space  *resolv.Space
}

func NewGimlarGame(speed float64) (*GimlarGame, error) {
	g := &GimlarGame{
		speed: speed,
		debug: true,
	}

	handler := &InputHandler{}
	var err error
	g.player, err = NewPlayer(handler, g.speed, radius)
	if err != nil {
		return nil, err
	}

	g.space = resolv.NewSpace(width, height, playerWidth, playerHeight)
	g.space.Add(g.player.Object)

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
	if g.debug {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", ebiten.ActualFPS())) // Print the current FPS
		// Draw grid overlay
		for i := 0; i < width; i += gridSpacing {
			vector.StrokeLine(screen, float32(i), 0, float32(i), float32(height), 1, color.White, false)
		}
		for i := 0; i < height; i += gridSpacing {
			vector.StrokeLine(screen, 0, float32(i), float32(width), float32(i), 1, color.White, false)
		}
	}
}
