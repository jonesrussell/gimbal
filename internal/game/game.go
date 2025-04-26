package game

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"image"
	"image/color"
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/entity/orbital"
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
	// LogInterval is the number of frames between periodic log messages
	LogInterval = DefaultTPS * 5 // Log every 5 seconds
	// FacingUpwardAngle is the angle representing facing upward
	FacingUpwardAngle = 270
	// SpeedNormalizationFactor normalizes speed with frame rate
	SpeedNormalizationFactor = 60
	// HalfDivisor is used for division by 2
	HalfDivisor = 2
)

//go:embed assets/*
var assets embed.FS

// GimlarGame represents the main game state
type GimlarGame struct {
	config       *common.GameConfig
	player       player.PlayerInterface
	stars        *stars.Manager
	inputHandler input.Interface
	logger       common.Logger
	isPaused     bool
	// State tracking for logging
	lastLoggedPos common.Point
	frameCount    int
	logInterval   int
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
	screenCenterX := float64(config.ScreenSize.Width) / common.CenterDivisor
	screenCenterY := float64(config.ScreenSize.Height) / common.CenterDivisor
	orbitRadius := float64(config.ScreenSize.Height) / RadiusDivisor

	playerConfig := &common.EntityConfig{
		Position: common.Point{
			X: screenCenterX, // Center X
			Y: screenCenterY, // Center Y
		},
		Size:   common.Size{Width: config.PlayerSize.Width, Height: config.PlayerSize.Height},
		Speed:  config.Speed,
		Radius: orbitRadius,
	}

	// Create player sprite
	playerSprite := ebitensprite.NewSprite(ebiten.NewImageFromImage(img))
	player, err := player.New(playerConfig, playerSprite, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create player: %w", err)
	}

	// Set initial angle to start at bottom of circle (180 degrees)
	initialAngle := common.Angle(180)
	// Set facing angle to 0 degrees (facing upward)
	player.SetAngle(initialAngle)          // Set orbital position first
	player.SetFacingAngle(common.Angle(0)) // Then set facing angle to 0 (up)

	logger.Debug("Player created",
		"position", player.GetPosition(),
		"center", playerConfig.Position,
		"radius", playerConfig.Radius,
		"facing_angle", player.GetFacingAngle(),
		"orbital_angle", player.GetAngle(),
		"screen_size", config.ScreenSize,
	)

	return &GimlarGame{
		config:       config,
		player:       player,
		stars:        starManager,
		inputHandler: inputHandler,
		logger:       logger,
		isPaused:     false,
		frameCount:   0,
		logInterval:  LogInterval,
	}, nil
}

// Layout implements ebiten.Game interface
func (g *GimlarGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.config.ScreenSize.Width, g.config.ScreenSize.Height
}

// Update implements ebiten.Game interface
func (g *GimlarGame) Update() error {
	g.frameCount++

	// Handle input
	g.inputHandler.HandleInput()

	// Check for pause
	if g.inputHandler.IsPausePressed() {
		g.isPaused = !g.isPaused
		g.logger.Debug("Game paused", "is_paused", g.isPaused)
	}

	// Check for quit
	if g.inputHandler.IsQuitPressed() {
		g.logger.Debug("Quit requested")
		return errors.New("game quit requested")
	}

	if g.isPaused {
		return nil
	}

	// Simplified movement logic
	inputAngle := g.inputHandler.GetMovementInput()
	if inputAngle != 0 {
		// Update orbital angle
		currentAngle := g.player.GetAngle()
		newAngle := currentAngle + inputAngle
		g.player.SetAngle(newAngle)

		// Calculate facing angle based on position relative to center
		playerPos := g.player.GetPosition()
		centerX := float64(g.config.ScreenSize.Width) / common.CenterDivisor
		centerY := float64(g.config.ScreenSize.Height) / common.CenterDivisor

		// Calculate angle from player to center
		dx := centerX - playerPos.X
		dy := centerY - playerPos.Y
		// atan2 gives angle in radians, convert to degrees
		baseAngle := math.Atan2(dy, dx) * orbital.RadiansToDegrees
		// Normalize to 0-360 range
		if baseAngle < 0 {
			baseAngle += 360
		}
		// Add 90 degrees to align sprite (sprite's 0° faces up, atan2's 0° faces right)
		facingAngle := common.Angle(baseAngle) + 90

		g.player.SetFacingAngle(facingAngle)

		// Log movement for debugging
		g.logger.Debug("Player moved",
			"position", playerPos,
			"orbital_angle", float64(newAngle),
			"facing_angle", float64(facingAngle),
			"base_angle", baseAngle,
			"input_angle", inputAngle,
			"center", common.Point{X: centerX, Y: centerY},
		)
	}

	// Update entities
	g.player.Update()
	g.stars.Update()

	// Log state periodically or when it changes significantly
	if g.frameCount%LogInterval == 0 {
		pos := g.player.GetPosition()
		if pos != g.lastLoggedPos {
			g.logger.Debug("Game state",
				"frame", g.frameCount,
				"position", pos,
				"tps", ebiten.ActualTPS(),
				"fps", ebiten.ActualFPS(),
			)
			g.lastLoggedPos = pos
		}
	}

	return nil
}

// Draw implements ebiten.Game interface
func (g *GimlarGame) Draw(screen *ebiten.Image) {
	// Skip drawing if screen is nil (testing)
	if screen == nil {
		return
	}

	// Clear the screen with a dark background
	screen.Fill(color.RGBA{0, 0, 0, 255})

	// Draw stars
	if g.stars != nil {
		g.stars.Draw(screen)
	}

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
func (g *GimlarGame) GetPlayer() player.PlayerInterface {
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
	facingAngle := g.player.GetFacingAngle()
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Position: (%.2f, %.2f)", pos.X, pos.Y),
		DebugTextMargin, DebugTextMargin)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Facing: %.2f°", float64(facingAngle)),
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
