package game

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/entity/player"
	"github.com/jonesrussell/gimbal/internal/entity/stars"
	"github.com/jonesrussell/gimbal/internal/input"
	"github.com/jonesrussell/gimbal/internal/logger"
)

const (
	// DebugTextMargin is the margin for debug text from screen edges
	DebugTextMargin = 10
	// DebugTextLineHeight is the vertical spacing between debug text lines
	DebugTextLineHeight = 20
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
			X: float64(config.ScreenSize.Width) / common.CenterDivisor,
			Y: float64(config.ScreenSize.Height) / common.CenterDivisor,
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

			logger.GlobalLogger.Debug("Updating player angle",
				"input_angle", inputAngle.ToRadians()/common.DegreesToRadians,
				"current_angle", currentAngle.ToRadians()/common.DegreesToRadians,
				"angle_step", g.config.AngleStep,
				"step", step.ToRadians()/common.DegreesToRadians,
				"new_angle", newAngle.ToRadians()/common.DegreesToRadians,
			)
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
	g.stars.Draw(screen)

	// Draw player
	g.player.Draw(screen)

	// Draw debug info if enabled
	if g.config.Debug {
		g.drawDebugInfo(screen)
	}
}

// GetPlayer returns the player entity
func (g *GimlarGame) GetPlayer() *player.Player {
	return g.player
}

// GetRadius returns the game's radius
func (g *GimlarGame) GetRadius() float64 {
	return g.config.Radius
}

// GetStars returns the stars from the star manager
func (g *GimlarGame) GetStars() []*stars.Star {
	return g.stars.GetStars()
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
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Position: (%.2f, %.2f)", pos.X, pos.Y),
		DebugTextMargin, DebugTextMargin)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Angle: %.2f°", angle),
		DebugTextMargin, DebugTextMargin+DebugTextLineHeight)
}

// SimulateKeyPress simulates a key press for testing
func (g *GimlarGame) SimulateKeyPress(key ebiten.Key) {
	g.input.SimulateKeyPress(key)
}

// SimulateKeyRelease simulates a key release for testing
func (g *GimlarGame) SimulateKeyRelease(key ebiten.Key) {
	g.input.SimulateKeyRelease(key)
}

// EnableTestMode enables test mode for input simulation
func (g *GimlarGame) EnableTestMode(enabled bool) {
	g.input.SetTestMode(enabled)
}
