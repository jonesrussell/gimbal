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
	radius    float64
	starImage *ebiten.Image
	Debug     bool
)

type GimlarGame struct {
	player *player.Player
	stars  []Star
	speed  float64
	space  *resolv.Space
	prevX  float64
	prevY  float64
	logger *zap.Logger
	config *config.Config
}

func NewGimlarGame(logger *zap.Logger, cfg *config.Config) (*GimlarGame, error) {
	logger.Debug("Initializing GimlarGame")
	Debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))

	radius = float64(cfg.Screen.Height/2) * 0.75
	starImage = ebiten.NewImage(1, 1)
	starImage.Fill(color.White)

	g := &GimlarGame{
		player: &player.Player{},
		stars:  []Star{},
		speed:  cfg.Game.Speed,
		space:  &resolv.Space{},
		prevX:  0,
		prevY:  0,
		logger: logger,
		config: cfg,
	}

	// Initialize stars
	if starImage == nil {
		return nil, fmt.Errorf("starImage is not loaded")
	}
	stars, err := initializeStars(100, starImage)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize stars: %w", err)
	}
	g.stars = stars

	handler := player.NewInputHandler()

	// Load the player sprite.
	imageData, rfErr := assets.ReadFile("assets/player.png")
	if rfErr != nil {
		logger.Error("Failed to load player image", zap.Error(rfErr))
	}

	if imageData == nil {
		logger.Error("Image data is nil")
		return nil, fmt.Errorf("failed to load player image data")
	}
	logger.Debug("Player image data loaded", zap.Int("bytes", len(imageData)))

	playerImg, imgFormat, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		logger.Error("Failed to decode player image", zap.Error(err))
		return nil, fmt.Errorf("failed to decode player image: %w", err)
	}
	logger.Debug("Player image decoded successfully",
		zap.String("format", imgFormat),
		zap.Int("width", playerImg.Bounds().Dx()),
		zap.Int("height", playerImg.Bounds().Dy()))

	var npErr error
	screenCenter := image.Pt(cfg.Screen.Width/2, cfg.Screen.Height/2)
	g.player, npErr = player.NewPlayer(handler, g.speed, playerImg, screenCenter)
	if npErr != nil {
		logger.Error("Failed to create player", zap.Error(npErr))
		return nil, npErr
	}

	g.space = resolv.NewSpace(
		cfg.Screen.Width,
		cfg.Screen.Height,
		cfg.Player.Width,
		cfg.Player.Height,
	)
	g.space.Add(g.player.Object)

	return g, nil
}

func (g *GimlarGame) Run() error {
	ebiten.SetWindowTitle("Gimbal")
	ebiten.SetWindowSize(g.config.Screen.Width, g.config.Screen.Height)
	g.logger.Debug("Starting game",
		zap.Int("window_width", g.config.Screen.Width),
		zap.Int("window_height", g.config.Screen.Height))
	return ebiten.RunGame(g)
}

func (g *GimlarGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.config.Screen.Width, g.config.Screen.Height
}

func (g *GimlarGame) Update() error {
	g.logger.Debug("Update frame")

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
	g.logger.Debug("Drawing frame",
		zap.Int("screen_width", screen.Bounds().Dx()),
		zap.Int("screen_height", screen.Bounds().Dy()))

	// Fill with a visible color to verify drawing is happening
	screen.Fill(color.RGBA{20, 20, 40, 255}) // Dark blue instead of black

	g.drawStars(screen)
	g.drawPlayer(screen)

	if Debug {
		g.DrawDebugInfo(screen)
	}
}

func (g *GimlarGame) drawPlayer(screen *ebiten.Image) {
	g.logger.Debug("Drawing player",
		zap.Float64("x", g.player.Object.Position.X),
		zap.Float64("y", g.player.Object.Position.Y))
	g.player.Draw(screen)
}

func (g *GimlarGame) GetRadius() float64 {
	return radius
}
