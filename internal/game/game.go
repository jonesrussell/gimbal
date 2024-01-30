package game

import (
	"fmt"
	"image"
	"image/color"
	"log/slog"
	"os"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/jonesrussell/gimbal/internal/logger"
	"github.com/solarlune/resolv"
)

var (
	screenWidth, screenHeight = 640, 480
	radius                    = float64(screenHeight/2) * 0.75
	center                    = image.Point{X: screenWidth / 2, Y: screenHeight / 2}
	debugGridSpacing          = 32
	playerWidth, playerHeight = 16, 16
	gameStarted               bool // Debugging check if game started
	Debug                     bool
)

type GimlarGame struct {
	player *Player
	speed  float64
	space  *resolv.Space
	logger slog.Logger
}

func NewGimlarGame(speed float64) (*GimlarGame, error) {
	g := &GimlarGame{
		speed:  speed,
		logger: logger.NewSlogHandler(),
	}

	Debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))

	handler := &InputHandler{}

	// Load the player sprite.
	spriteImage, _, loadErr := ebitenutil.NewImageFromFile("assets/player.png")
	if loadErr != nil {
		g.logger.Error("Failed to load player sprite", "loadErr", loadErr)
		os.Exit(1)
	}

	var err error
	g.player, err = NewPlayer(handler, g.speed, spriteImage, g.logger)
	if err != nil {
		g.logger.Error("Failed to create player", "err", err)
		os.Exit(1)
	}

	g.space = resolv.NewSpace(screenWidth, screenHeight, playerWidth, playerHeight)
	g.space.Add(g.player.Object)

	return g, nil
}

func (g *GimlarGame) Run() error {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	return ebiten.RunGame(g)
}

func (g *GimlarGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *GimlarGame) Update() error {
	// Update the player's state
	g.player.Update()

	return nil
}

func (g *GimlarGame) Draw(screen *ebiten.Image) {
	g.player.Draw(screen)

	// Draw debug info if debug is true
	if Debug {
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
	for i := 0; i < screenWidth; i += debugGridSpacing {
		vector.StrokeLine(screen, float32(i), 0, float32(i), float32(screenHeight), 1, color.White, false)
	}
	for i := 0; i < screenHeight; i += debugGridSpacing {
		vector.StrokeLine(screen, 0, float32(i), float32(screenWidth), float32(i), 1, color.White, false)
	}
}
