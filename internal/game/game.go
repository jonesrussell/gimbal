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
	"github.com/jonesrussell/gimbal/internal/entity/player"
	ebitensprite "github.com/jonesrussell/gimbal/internal/entity/player/ebiten"
	"github.com/jonesrussell/gimbal/internal/entity/stars"
	"github.com/jonesrussell/gimbal/internal/input"
)

const (
	// Debug and UI constants
	DebugTextMargin     = 10
	DebugTextLineHeight = 20

	// Game configuration constants
	RadiusDivisor            = 3
	DefaultTPS               = 60
	LogInterval              = DefaultTPS * 5 // Log every 5 seconds
	SpeedNormalizationFactor = 60
	HalfDivisor              = 2

	// Angle constants
	// InitialFacingAngle is the starting angle for the player sprite
	InitialFacingAngle = 0 // Start facing upward
	// FullCircleDegrees represents a full circle in degrees
	FullCircleDegrees = 360
	// RadiansToDegrees is the multiplier to convert radians to degrees
	RadiansToDegrees = 180 / math.Pi
	// RightToUpwardOffset is the angle offset to convert from right-facing 0° to upward-facing 0°
	RightToUpwardOffset = 90
	InitialOrbitalAngle = 180 // Start at bottom of circle
)

// Error definitions
var (
	ErrNilConfig     = errors.New("config cannot be nil")
	ErrNilLogger     = errors.New("logger cannot be nil")
	ErrLoadingSprite = errors.New("failed to load player sprite")
	ErrUserQuit      = errors.New("user requested quit")
	ErrGameLoop      = errors.New("game loop error")
)

//go:embed assets/*
var assets embed.FS

// StarManager defines the interface for star management
type StarManager interface {
	Update()
	Draw(screen any)
	GetPosition() common.Point
	GetStars() []*stars.Star
	Cleanup()
}

// GimlarGame represents the main game state
type GimlarGame struct {
	config       *common.GameConfig
	player       player.PlayerInterface
	stars        StarManager
	inputHandler input.Interface
	logger       common.Logger
	isPaused     bool
	// State tracking for logging
	lastLoggedPos common.Point
	frameCount    int
	logInterval   int
	deltaTime     float64 // Time since last frame in seconds
}

// New creates a new game instance
func New(config *common.GameConfig, logger common.Logger) (*GimlarGame, error) {
	if config == nil {
		return nil, ErrNilConfig
	}
	if logger == nil {
		return nil, ErrNilLogger
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

	// Load and initialize player
	player, err := initializePlayer(config, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize player: %w", err)
	}

	return &GimlarGame{
		config:       config,
		player:       player,
		stars:        starManager,
		inputHandler: inputHandler,
		logger:       logger,
		isPaused:     false,
		frameCount:   0,
		logInterval:  LogInterval,
		deltaTime:    1.0 / float64(DefaultTPS),
	}, nil
}

// initializePlayer loads the sprite and creates the player entity
func initializePlayer(config *common.GameConfig, logger common.Logger) (player.PlayerInterface, error) {
	// Load the player sprite
	imageData, err := assets.ReadFile("assets/player.png")
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrLoadingSprite, err)
	}
	logger.Debug("Player image loaded", "size", len(imageData))

	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to decode image: %v", ErrLoadingSprite, err)
	}
	logger.Debug("Player image decoded",
		"bounds", img.Bounds(),
		"color_model", img.ColorModel(),
	)

	// Calculate player position and configuration
	screenCenterX := float64(config.ScreenSize.Width) / common.CenterDivisor
	screenCenterY := float64(config.ScreenSize.Height) / common.CenterDivisor
	orbitRadius := float64(config.ScreenSize.Height) / RadiusDivisor

	playerConfig := &common.EntityConfig{
		Position: common.Point{
			X: screenCenterX,
			Y: screenCenterY,
		},
		Size:   common.Size{Width: config.PlayerSize.Width, Height: config.PlayerSize.Height},
		Speed:  config.Speed,
		Radius: orbitRadius,
	}

	// Create player sprite and entity
	playerSprite := ebitensprite.NewSprite(ebiten.NewImageFromImage(img))
	player, err := player.New(playerConfig, playerSprite, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create player: %w", err)
	}

	logger.Debug("Player created",
		"position", player.GetPosition(),
		"center", playerConfig.Position,
		"radius", playerConfig.Radius,
		"facing_angle", player.GetFacingAngle(),
		"orbital_angle", player.GetAngle(),
		"screen_size", config.ScreenSize,
	)

	return player, nil
}

// Layout implements ebiten.Game interface
func (g *GimlarGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.config.ScreenSize.Width, g.config.ScreenSize.Height
}

// Update implements ebiten.Game interface
func (g *GimlarGame) Update() error {
	g.frameCount++
	g.deltaTime = 1.0 / ebiten.ActualTPS()

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
		return ErrUserQuit
	}

	if g.isPaused {
		return nil
	}

	if err := g.updateGameState(); err != nil {
		g.logger.Error("Failed to update game state", "error", err)
		return fmt.Errorf("game state update error: %w", err)
	}

	return nil
}

// updateGameState handles the main game state updates
func (g *GimlarGame) updateGameState() error {
	// Handle movement
	inputAngle := g.inputHandler.GetMovementInput()
	if inputAngle != 0 {
		// Apply frame rate independence to movement
		scaledInput := float64(inputAngle) * g.deltaTime * SpeedNormalizationFactor

		// Update orbital angle
		currentAngle := g.player.GetAngle()
		newAngle := currentAngle + common.Angle(scaledInput)
		if err := g.player.SetAngle(newAngle); err != nil {
			return fmt.Errorf("failed to set angle: %w", err)
		}

		// Log movement for debugging
		g.logger.Debug("Player moved",
			"position", g.player.GetPosition(),
			"orbital_angle", float64(newAngle),
			"facing_angle", g.player.GetFacingAngle(),
			"input_angle", scaledInput,
			"delta_time", g.deltaTime,
		)
	}

	// Calculate facing angle to always point towards the center
	pos := g.player.GetPosition()
	centerX := float64(g.config.ScreenSize.Width) / 2
	centerY := float64(g.config.ScreenSize.Height) / 2

	// Calculate angle from player to center using atan2
	// Note: We use centerY - pos.Y because our Y axis increases downward
	dx := centerX - pos.X
	dy := centerY - pos.Y

	// Handle case where player is at center point
	if dx == 0 && dy == 0 {
		// Keep current facing angle if at center
		return nil
	}

	baseAngle := math.Atan2(dy, dx) * RadiansToDegrees

	// Convert to game's angle system:
	// 1. Add 90° because sprite's 0° faces up while atan2's 0° faces right
	// 2. Normalize to [0, 360) range
	facingAngle := baseAngle + RightToUpwardOffset
	if facingAngle < 0 {
		facingAngle += FullCircleDegrees
	} else if facingAngle >= FullCircleDegrees {
		facingAngle -= FullCircleDegrees
	}

	g.player.SetFacingAngle(common.Angle(facingAngle))

	// Update entities
	g.player.Update()
	g.stars.Update()

	// Periodic state logging
	if g.shouldLogState() {
		g.logGameState()
	}

	return nil
}

// shouldLogState determines if the game state should be logged
func (g *GimlarGame) shouldLogState() bool {
	return g.frameCount%g.logInterval == 0 && g.config.Debug
}

// logGameState logs the current game state
func (g *GimlarGame) logGameState() {
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

// Draw implements ebiten.Game interface
func (g *GimlarGame) Draw(screen *ebiten.Image) {
	if screen == nil {
		g.logger.Debug("Draw skipped: screen is nil")
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

// Cleanup performs cleanup of game resources
func (g *GimlarGame) Cleanup() {
	g.logger.Debug("Cleaning up game resources")

	// Cleanup player resources
	if g.player != nil {
		if cleaner, ok := g.player.(interface{ Cleanup() }); ok {
			cleaner.Cleanup()
		}
	}

	// Cleanup star resources
	if g.stars != nil {
		g.stars.Cleanup()
	}

	// Sync logger before exit
	if f, ok := g.logger.(interface{ Sync() error }); ok {
		if err := f.Sync(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to sync logger: %v\n", err)
		}
	}
}

// Run starts the game loop
func (g *GimlarGame) Run() error {
	// Ensure cleanup is performed
	defer g.Cleanup()

	// Force stdout to be unbuffered
	os.Stdout.Sync()

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
		// Handle different error types
		switch {
		case errors.Is(err, ErrUserQuit):
			g.logger.Info("Game closed by user")
			return nil
		case errors.Is(err, ErrGameLoop):
			g.logger.Error("Game loop error", "error", err)
			return fmt.Errorf("game loop error: %w", err)
		default:
			g.logger.Error("Unexpected error", "error", err)
			return fmt.Errorf("unexpected error: %w", err)
		}
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

// NewWithDependencies creates a new game instance with injected dependencies (for testing)
func NewWithDependencies(
	config *common.GameConfig,
	logger common.Logger,
	player player.PlayerInterface,
	stars StarManager,
	inputHandler input.Interface,
) (*GimlarGame, error) {
	if config == nil {
		return nil, ErrNilConfig
	}
	if logger == nil {
		return nil, ErrNilLogger
	}
	if player == nil {
		return nil, errors.New("player cannot be nil")
	}
	if stars == nil {
		return nil, errors.New("stars cannot be nil")
	}
	if inputHandler == nil {
		return nil, errors.New("input handler cannot be nil")
	}

	return &GimlarGame{
		config:       config,
		player:       player,
		stars:        stars,
		inputHandler: inputHandler,
		logger:       logger,
		isPaused:     false,
		frameCount:   0,
		logInterval:  LogInterval,
		deltaTime:    1.0 / float64(DefaultTPS),
	}, nil
}

// GetScreenSize returns the game's screen size
func (g *GimlarGame) GetScreenSize() common.Size {
	return g.config.ScreenSize
}

// GetPlayerSize returns the player's size
func (g *GimlarGame) GetPlayerSize() common.Size {
	return g.config.PlayerSize
}

// GetSpeed returns the game's speed
func (g *GimlarGame) GetSpeed() float64 {
	return g.config.Speed
}

// GetDebug returns whether debug mode is enabled
func (g *GimlarGame) GetDebug() bool {
	return g.config.Debug
}
