package game

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/entity/player"
	"github.com/jonesrussell/gimbal/internal/entity/stars"
	"github.com/jonesrussell/gimbal/internal/input"
	"github.com/jonesrussell/gimbal/internal/logger"
)

//go:embed assets/*
var assets embed.FS

// GimlarGame represents the main game state
type GimlarGame struct {
	config   *common.GameConfig
	player   *player.Player
	stars    *stars.Manager
	input    *input.Handler
	isPaused bool
}

// New creates a new game instance
func New(config *common.GameConfig) (*GimlarGame, error) {
	if config == nil {
		return nil, errors.New("config cannot be nil")
	}

	// Create input handler
	inputHandler := input.New()

	// Create star manager
	starManager := stars.NewManager(
		config.ScreenSize,
		config.NumStars,
		config.StarSize,
		config.StarSpeed,
	)

	// Load the player sprite
	imageData, err := assets.ReadFile("assets/player.png")
	if err != nil {
		return nil, fmt.Errorf("failed to load player image: %w", err)
	}

	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return nil, fmt.Errorf("failed to decode player image: %w", err)
	}

	// Create player entity
	playerConfig := &common.EntityConfig{
		Position: common.Point{
			X: float64(config.ScreenSize.Width / 2),
			Y: float64(config.ScreenSize.Height / 2),
		},
		Size: common.Size{
			Width:  config.PlayerSize.Width,
			Height: config.PlayerSize.Height,
		},
		Speed:  config.Speed,
		Radius: config.Radius,
	}

	player, err := player.New(playerConfig, ebiten.NewImageFromImage(img))
	if err != nil {
		return nil, fmt.Errorf("failed to create player: %w", err)
	}

	return &GimlarGame{
		config:   config,
		player:   player,
		stars:    starManager,
		input:    inputHandler,
		isPaused: false,
	}, nil
}

// Layout implements ebiten.Game interface
func (g *GimlarGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.config.ScreenSize.Width, g.config.ScreenSize.Height
}

// Update implements ebiten.Game interface
func (g *GimlarGame) Update() error {
	// Handle input
	g.input.HandleInput()

	// Check for pause
	if g.input.IsPausePressed() {
		g.isPaused = !g.isPaused
	}

	// Check for quit
	if g.input.IsQuitPressed() {
		return errors.New("game quit requested")
	}

	if !g.isPaused {
		// Update player angle based on input
		inputAngle := g.input.GetMovementInput()
		if inputAngle != 0 {
			currentAngle := g.player.GetAngle()
			// Use Angle.Mul for scalar multiplication
			step := inputAngle.Mul(g.config.AngleStep)
			// Use Angle.Add for angle addition
			newAngle := currentAngle.Add(step)
			g.player.SetAngle(newAngle)
		}

		// Update entities
		g.player.Update()
		g.stars.Update()
	}

	return nil
}

// Draw implements ebiten.Game interface
func (g *GimlarGame) Draw(screen *ebiten.Image) {
	// Draw stars
	g.drawStars(screen)

	// Draw player
	g.player.Draw(screen)

	// Draw debug info if enabled
	if g.config.Debug {
		g.drawDebugInfo(screen)
	}
}

// Run starts the game loop
func (g *GimlarGame) Run() error {
	ebiten.SetWindowSize(g.config.ScreenSize.Width, g.config.ScreenSize.Height)
	ebiten.SetWindowTitle("Gimbal Game")
	return ebiten.RunGame(g)
}

// drawDebugInfo draws debug information on screen
func (g *GimlarGame) drawDebugInfo(screen *ebiten.Image) {
	pos := g.player.GetPosition()
	angle := g.player.GetAngle()
	logger.GlobalLogger.Debug("Debug info",
		"position", fmt.Sprintf("(%.2f, %.2f)", pos.X, pos.Y),
		"angle", fmt.Sprintf("%.2f°", angle),
	)
}
