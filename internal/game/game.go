package game

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"image"
	"image/color"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/entity/player"
	ebitensprite "github.com/jonesrussell/gimbal/internal/entity/player/ebiten"
	"github.com/jonesrussell/gimbal/internal/entity/stars"
	"github.com/jonesrussell/gimbal/internal/input"
)

const (
	// DebugTextMargin is the margin for debug text from screen edges
	DebugTextMargin = 10
	// DebugTextLineHeight is the vertical spacing between debug text lines
	DebugTextLineHeight = 20
	// FacingAngleOffset is the angle offset to make the player face the center
	FacingAngleOffset = 180
	// RadiusDivisor is used to calculate the player's orbit radius as a fraction of screen height
	RadiusDivisor = 3
	// DefaultTPS is the default ticks per second for the game loop
	DefaultTPS = 60
)

//go:embed assets/*
var assets embed.FS

// GimlarGame represents the main game state
type GimlarGame struct {
	config       *common.GameConfig
	player       *player.Player
	stars        *stars.Manager
	inputHandler input.Interface
	logger       common.Logger
	isPaused     bool
}

// New creates a new game instance
func New(config *common.GameConfig, logger common.Logger) (*GimlarGame, error) {
	if config == nil {
		return nil, errors.New("config cannot be nil")
	}
	if logger == nil {
		return nil, errors.New("logger cannot be nil")
	}

	logger.Debug("Creating new game instance",
		"screen_size", config.ScreenSize,
		"player_size", config.PlayerSize,
		"num_stars", config.NumStars,
	)

	// Create input handler
	inputHandler := input.New(logger)
	logger.Debug("Input handler created")

	// Create star manager
	starManager := stars.NewManager(
		config.ScreenSize,
		config.NumStars,
		config.StarSize,
		config.StarSpeed,
	)
	logger.Debug("Star manager created",
		"num_stars", len(starManager.GetStars()),
	)

	// Load the player sprite
	imageData, err := assets.ReadFile("assets/player.png")
	if err != nil {
		return nil, fmt.Errorf("failed to load player image: %w", err)
	}
	logger.Debug("Player image loaded",
		"size", len(imageData),
	)

	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return nil, fmt.Errorf("failed to decode player image: %w", err)
	}
	logger.Debug("Player image decoded",
		"bounds", img.Bounds(),
		"color_model", img.ColorModel(),
	)

	// Create player entity
	playerConfig := &common.EntityConfig{
		Position: common.Point{
			X: float64(config.ScreenSize.Width) / common.CenterDivisor,
			Y: float64(config.ScreenSize.Height) / common.CenterDivisor,
		},
		Size:   config.ScreenSize,
		Speed:  config.Speed,
		Radius: float64(config.ScreenSize.Height) / RadiusDivisor,
	}

	// Create player sprite
	playerSprite := ebitensprite.NewSprite(ebiten.NewImageFromImage(img))
	player, err := player.New(playerConfig, playerSprite, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create player: %w", err)
	}
	logger.Debug("Player created",
		"position", player.GetPosition(),
		"angle", player.GetAngle(),
	)

	return &GimlarGame{
		config:       config,
		player:       player,
		stars:        starManager,
		inputHandler: inputHandler,
		logger:       logger,
		isPaused:     false,
	}, nil
}

// Layout implements ebiten.Game interface
func (g *GimlarGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.config.ScreenSize.Width, g.config.ScreenSize.Height
}

// Update implements ebiten.Game interface
func (g *GimlarGame) Update() error {
	// Handle input
	g.inputHandler.HandleInput()

	// Check for pause
	if g.inputHandler.IsPausePressed() {
		g.isPaused = !g.isPaused
		g.logger.Debug("Game paused",
			"is_paused", g.isPaused,
		)
	}

	// Check for quit
	if g.inputHandler.IsQuitPressed() {
		g.logger.Debug("Quit requested")
		return errors.New("game quit requested")
	}

	if !g.isPaused {
		// Update player angle based on input
		inputAngle := g.inputHandler.GetMovementInput()
		g.logger.Debug("Game update",
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

			g.logger.Debug("Player movement",
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
		g.logger.Debug("Skipping draw - screen is nil")
		return
	}

	// Clear the screen with a dark background
	screen.Fill(color.RGBA{0, 0, 0, 255})
	g.logger.Debug("Screen cleared")

	// Draw stars
	if g.stars != nil {
		g.stars.Draw(screen)
		g.logger.Debug("Stars drawn")
	}

	// Draw player
	if g.player != nil {
		g.player.Draw(screen, nil)
		g.logger.Debug("Player drawn",
			"position", g.player.GetPosition(),
			"angle", g.player.GetAngle(),
		)
	}

	// Draw debug info if enabled
	if g.config.Debug {
		g.drawDebugInfo(screen)
		g.logger.Debug("Debug info drawn")
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
	// Force stdout to be unbuffered
	os.Stdout.Sync()

	// Ensure logger is synced on exit
	if f, ok := g.logger.(interface{ Sync() error }); ok {
		defer func() {
			if err := f.Sync(); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to sync logger: %v\n", err)
			}
		}()
	}

	g.logger.Debug("Setting up game window",
		"width", g.config.ScreenSize.Width,
		"height", g.config.ScreenSize.Height,
	)

	ebiten.SetWindowSize(g.config.ScreenSize.Width, g.config.ScreenSize.Height)
	ebiten.SetWindowTitle("Gimbal Game")

	// Set window options for better visibility
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetFullscreen(false)
	ebiten.SetVsyncEnabled(true)
	ebiten.SetTPS(DefaultTPS)
	ebiten.SetMaxTPS(DefaultTPS)

	g.logger.Debug("Starting game loop")

	// Run the game loop
	if err := ebiten.RunGame(g); err != nil {
		g.logger.Debug("Game loop ended with error", "error", err)
		return fmt.Errorf("game loop error: %w", err)
	}

	return nil
}

// drawDebugInfo draws debug information on screen
func (g *GimlarGame) drawDebugInfo(screen *ebiten.Image) {
	pos := g.player.GetPosition()
	angle := g.player.GetAngle()
	g.logger.Debug("Debug info",
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
	g.inputHandler = handler
}

// IsPaused returns whether the game is paused
func (g *GimlarGame) IsPaused() bool {
	return g.isPaused
}
