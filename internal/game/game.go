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

const (
	screenWidth      = 640
	screenHeight     = 480
	playerWidth      = 16
	playerHeight     = 16
	debugGridSpacing = 32
)

var (
	assetsBasePath = "/home/russell/Development/gimbal/assets/"
	// assetsBasePath            = "assets/"
	radius = float64(screenHeight/2) * 0.75
	center = image.Point{X: screenWidth / 2, Y: screenHeight / 2}

	starImage *ebiten.Image

	gameStarted bool
	Debug       bool
)

type GimlarGame struct {
	player       *Player
	stars        []Star
	speed        float64
	space        *resolv.Space
	logger       slog.Logger
	prevX, prevY float64
}

func init() {
	// Create a single star image that will be used for all stars
	starImage = ebiten.NewImage(1, 1)
	starImage.Fill(color.White)
}

func NewGimlarGame(speed float64) (*GimlarGame, error) {
	Debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))

	var level slog.Level
	if Debug {
		level = slog.LevelDebug
	} else {
		level = slog.LevelInfo // Or whatever non-debug level you prefer
	}

	g := &GimlarGame{
		player: &Player{},
		stars:  []Star{},
		speed:  speed,
		space:  &resolv.Space{},
		logger: logger.NewSlogHandler(level),
		prevX:  0,
		prevY:  0,
	}

	// Initialize stars
	if starImage == nil {
		return nil, fmt.Errorf("starImage is not loaded")
	}
	g.stars = initializeStars(100, starImage)

	handler := &InputHandler{}

	// Load the player sprite.
	spriteImage, _, loadErr := ebitenutil.NewImageFromFile(assetsBasePath + "player.png")
	if loadErr != nil {
		g.logger.Error("Failed to load player sprite", "error", loadErr)
		return nil, loadErr // Return the error instead of exiting
	}

	var err error
	g.player, err = NewPlayer(handler, g.speed, spriteImage, g.logger)
	if err != nil {
		g.logger.Error("Failed to create player: %v", err)
		return nil, err // Return the error instead of exiting
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
	// Update the stars
	g.updateStars()

	// Update the player's state
	g.player.Update()
	g.player.updatePosition()

	// Log the player's position after updating if it has changed
	if g.player.Object.Position.X != g.prevX || g.player.Object.Position.Y != g.prevY {
		g.logger.Debug("Player position after update", "X", g.player.Object.Position.X, "Y", g.player.Object.Position.Y)
		g.prevX = g.player.Object.Position.X
		g.prevY = g.player.Object.Position.Y
	}

	return nil
}

func (g *GimlarGame) Draw(screen *ebiten.Image) {
	// Draw the stars
	g.drawStars(screen)

	// Draw the player
	g.drawPlayer(screen)

	// Draw debug info if debug is true
	if Debug {
		g.DrawDebugInfo(screen)
	}
}

func (g *GimlarGame) drawPlayer(screen *ebiten.Image) {
	// Assuming the player has a Draw method
	g.player.Draw(screen)
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
