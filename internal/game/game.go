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
	ebitensprite "github.com/jonesrussell/gimbal/internal/entity/player/ebiten"
	"github.com/jonesrussell/gimbal/internal/entity/stars"
	"github.com/jonesrussell/gimbal/internal/input"
	"github.com/jonesrussell/gimbal/internal/logger"
)

const (
	// DebugTextMargin is the margin for debug text from screen edges
	DebugTextMargin = 10
	// DebugTextLineHeight is the vertical spacing between debug text lines
	DebugTextLineHeight = 20
	// FacingAngleOffset is the angle offset to make the player face the center
	FacingAngleOffset = 180
)

//go:embed assets/*
var assets embed.FS

// GimlarGame represents the main game state
type GimlarGame struct {
	config   *common.GameConfig
	player   *player.Player
	stars    *stars.Manager
	input    input.Interface
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

	// Create player sprite
	playerSprite := ebitensprite.NewSprite(ebiten.NewImageFromImage(img))
	player, err := player.New(playerConfig, playerSprite)
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
		logger.GlobalLogger.Debug("Game paused", "is_paused", g.isPaused)
	}

	// Check for quit
	if g.input.IsQuitPressed() {
		return errors.New("game quit requested")
	}

	if !g.isPaused {
		// Update player angle based on input
		inputAngle := g.input.GetMovementInput()
		logger.GlobalLogger.Debug("Game update",
			"input_angle", inputAngle,
			"is_paused", g.isPaused,
		)

		if inputAngle != 0 {
			currentAngle := g.player.GetAngle()
			// Add the input angle to the current position angle
			newAngle := currentAngle.Add(inputAngle)
			g.player.SetAngle(newAngle)

			// Make the player face the center by setting facing angle to 180 degrees from position angle
			// This ensures the player always points towards the center
			centerFacingAngle := newAngle.Add(common.Angle(FacingAngleOffset))
			g.player.SetFacingAngle(centerFacingAngle)

			logger.GlobalLogger.Debug("Player movement",
				"input_angle", inputAngle,
				"current_angle", currentAngle,
				"new_angle", newAngle,
				"position", g.player.GetPosition(),
				"facing_angle", g.player.GetFacingAngle(),
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
	// Skip drawing if screen is nil (testing)
	if screen == nil {
		return
	}

	// Draw stars
	g.stars.Draw(screen)

	// Draw player
	if g.player != nil {
		g.player.Draw(screen, nil)
	}

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

// SetInputHandler sets the input handler for the game
func (g *GimlarGame) SetInputHandler(handler input.Interface) {
	g.input = handler
}

// IsPaused returns whether the game is paused
func (g *GimlarGame) IsPaused() bool {
	return g.isPaused
}
