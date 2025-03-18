package game

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/entity/player"
	ebitensprite "github.com/jonesrussell/gimbal/internal/entity/player/ebiten"
	"github.com/jonesrussell/gimbal/internal/entity/stars"
	"github.com/jonesrussell/gimbal/internal/input"
	"github.com/jonesrussell/gimbal/internal/logger"
	"go.uber.org/zap"
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

	logger.GlobalLogger.Debug("Creating new game instance",
		zap.Any("screen_size", config.ScreenSize),
		zap.Any("player_size", config.PlayerSize),
		zap.Int("num_stars", config.NumStars),
	)

	// Create input handler
	inputHandler := input.New()
	logger.GlobalLogger.Debug("Input handler created")

	// Create star manager
	starManager := stars.NewManager(
		config.ScreenSize,
		config.NumStars,
		config.StarSize,
		config.StarSpeed,
	)
	logger.GlobalLogger.Debug("Star manager created",
		zap.Int("num_stars", len(starManager.GetStars())),
	)

	// Load the player sprite
	imageData, err := assets.ReadFile("assets/player.png")
	if err != nil {
		return nil, fmt.Errorf("failed to load player image: %w", err)
	}
	logger.GlobalLogger.Debug("Player image loaded",
		zap.Int("size", len(imageData)),
	)

	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return nil, fmt.Errorf("failed to decode player image: %w", err)
	}
	logger.GlobalLogger.Debug("Player image decoded",
		zap.Any("bounds", img.Bounds()),
		zap.Any("color_model", img.ColorModel()),
	)

	// Create player entity
	playerConfig := &common.EntityConfig{
		Position: common.Point{
			X: float64(config.ScreenSize.Width) / common.CenterDivisor,
			Y: float64(config.ScreenSize.Height) / common.CenterDivisor,
		},
		Size:   config.ScreenSize, // Pass screen size for center calculations
		Speed:  config.Speed,
		Radius: float64(config.ScreenSize.Height) / 3, // Radius for circular movement
	}

	// Create player sprite
	playerSprite := ebitensprite.NewSprite(ebiten.NewImageFromImage(img))
	player, err := player.New(playerConfig, playerSprite)
	if err != nil {
		return nil, fmt.Errorf("failed to create player: %w", err)
	}
	logger.GlobalLogger.Debug("Player created",
		zap.Any("position", player.GetPosition()),
		zap.Any("angle", player.GetAngle()),
	)

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
		logger.GlobalLogger.Debug("Game paused",
			zap.Bool("is_paused", g.isPaused),
		)
	}

	// Check for quit
	if g.input.IsQuitPressed() {
		logger.GlobalLogger.Debug("Quit requested")
		return errors.New("game quit requested")
	}

	if !g.isPaused {
		// Update player angle based on input
		inputAngle := g.input.GetMovementInput()
		logger.GlobalLogger.Debug("Game update",
			zap.Any("input_angle", inputAngle),
			zap.Bool("is_paused", g.isPaused),
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
				zap.Any("input_angle", inputAngle),
				zap.Any("current_angle", currentAngle),
				zap.Any("new_angle", newAngle),
				zap.Any("position", g.player.GetPosition()),
				zap.Any("facing_angle", g.player.GetFacingAngle()),
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
		logger.GlobalLogger.Debug("Skipping draw - screen is nil")
		return
	}

	// Clear the screen with a dark background
	screen.Fill(color.RGBA{0, 0, 0, 255})
	logger.GlobalLogger.Debug("Screen cleared")

	// Draw stars
	if g.stars != nil {
		g.stars.Draw(screen)
		logger.GlobalLogger.Debug("Stars drawn")
	}

	// Draw player
	if g.player != nil {
		g.player.Draw(screen, nil)
		logger.GlobalLogger.Debug("Player drawn",
			zap.Any("position", g.player.GetPosition()),
			zap.Any("angle", g.player.GetAngle()),
		)
	}

	// Draw debug info if enabled
	if g.config.Debug {
		g.drawDebugInfo(screen)
		logger.GlobalLogger.Debug("Debug info drawn")
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
	logger.GlobalLogger.Debug("Setting up game window",
		zap.Int("width", g.config.ScreenSize.Width),
		zap.Int("height", g.config.ScreenSize.Height),
	)

	ebiten.SetWindowSize(g.config.ScreenSize.Width, g.config.ScreenSize.Height)
	ebiten.SetWindowTitle("Gimbal Game")

	// Set window options for better visibility
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetFullscreen(false)

	// Set FPS to 60
	ebiten.SetMaxTPS(60)

	logger.GlobalLogger.Debug("Starting game loop")
	return ebiten.RunGame(g)
}

// drawDebugInfo draws debug information on screen
func (g *GimlarGame) drawDebugInfo(screen *ebiten.Image) {
	pos := g.player.GetPosition()
	angle := g.player.GetAngle()
	logger.GlobalLogger.Debug("Debug info",
		zap.Any("position", fmt.Sprintf("(%.2f, %.2f)", pos.X, pos.Y)),
		zap.Any("angle", fmt.Sprintf("%.2f°", angle)),
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
