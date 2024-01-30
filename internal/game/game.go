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
	debugGridSpacing          = 32
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
	ebiten.SetWindowSize(width, height)
	return ebiten.RunGame(g)
}

func (g *GimlarGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return width, height
}

func (g *GimlarGame) Update() error {
	// Update the player's state
	g.player.Update()

	return nil
}

func (g *GimlarGame) Draw(screen *ebiten.Image) {
	g.player.Draw(screen)

	// Draw debug info if debug is true
	if g.IsDebugMode() {
		g.DrawDebugInfo(screen)
	}
}

func (g *GimlarGame) GetRadius() float64 {
	return radius
}

func (g *GimlarGame) DrawDebugInfo(screen *ebiten.Image) {
	// Print the current FPS
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", ebiten.ActualFPS()))

	// Draw grid overlay
	g.DrawGridOverlay(screen)
}

func (g *GimlarGame) DrawGridOverlay(screen *ebiten.Image) {
	// Draw grid overlay
	for i := 0; i < width; i += debugGridSpacing {
		vector.StrokeLine(screen, float32(i), 0, float32(i), float32(height), 1, color.White, false)
	}
	for i := 0; i < height; i += debugGridSpacing {
		vector.StrokeLine(screen, 0, float32(i), float32(width), float32(i), 1, color.White, false)
	}
}

func (g *GimlarGame) SetDebugMode(debug bool) {
	g.debug = debug
}

func (g *GimlarGame) IsDebugMode() bool {
	return g.debug
}

func GetWidth() int {
	return width
}

func GetHeight() int {
	return height
}

func GetCenter() image.Point {
	return center
}
