package game

import (
	"bytes"
	"embed"
	"fmt"
	"image"
	"image/color"
	"os"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jonesrussell/gimbal/internal/config"
	"github.com/jonesrussell/gimbal/player"
	"github.com/solarlune/resolv"
	"go.uber.org/zap"
)

//go:embed assets/*
var assets embed.FS

var (
	radius float64
	center image.Point

	starImage *ebiten.Image

	gameStarted bool
	Debug       bool
)

type GimlarGame struct {
	player *player.Player
	stars  []Star
	speed  float64
	space  *resolv.Space
	prevX  float64
	prevY  float64
	logger *zap.Logger
}

func init() {
	cfg := config.New()
	radius = float64(cfg.Screen.Height/2) * 0.75
	center = image.Point{X: cfg.Screen.Width / 2, Y: cfg.Screen.Height / 2}

	// Create a single star image that will be used for all stars
	starImage = ebiten.NewImage(1, 1)
	starImage.Fill(color.White)
}

func NewGimlarGame(logger *zap.Logger, cfg *config.Config) (*GimlarGame, error) {
	Debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))

	g := &GimlarGame{
		player: &player.Player{},
		stars:  []Star{},
		speed:  cfg.Speed,
		space:  &resolv.Space{},
		prevX:  0,
		prevY:  0,
		logger: logger,
	}

	// Initialize stars
	if starImage == nil {
		return nil, fmt.Errorf("starImage is not loaded")
	}
	g.stars = initializeStars(100, starImage)

	handler := player.NewInputHandler()

	// Load the player sprite.
	imageData, rfErr := assets.ReadFile("assets/player.png")
	if rfErr != nil {
		logger.Error("Failed to load player image", zap.Error(rfErr))
	}

	image, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		logger.Error("Failed to decode player image", zap.Error(err))
	}

	var npErr error
	g.player, npErr = player.NewPlayer(handler, g.speed, image)
	if npErr != nil {
		logger.Error("Failed to create player", zap.Error(npErr))
		return nil, npErr // Return the error instead of exiting
	}

	g.space = resolv.NewSpace(cfg.Screen.Width, cfg.Screen.Height, cfg.Player.Width, cfg.Player.Height)
	g.space.Add(g.player.Object)

	return g, nil
}

func (g *GimlarGame) Run() error {
	ebiten.SetWindowSize(cfg.Screen.Width, cfg.Screen.Height)
	return ebiten.RunGame(g)
}

func (g *GimlarGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	cfg := config.New()
	return cfg.Screen.Width, cfg.Screen.Height
}

func (g *GimlarGame) Update() error {
	// Update the stars
	g.updateStars()

	// Update the player's state
	g.player.Update()
	g.player.UpdatePosition()

	// Log the player's position after updating if it has changed
	if g.player.Object.Position.X != g.prevX || g.player.Object.Position.Y != g.prevY {
		g.logger.Debug("Player position after update",
			zap.Float64("X", g.player.Object.Position.X),
			zap.Float64("Y", g.player.Object.Position.Y))
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
