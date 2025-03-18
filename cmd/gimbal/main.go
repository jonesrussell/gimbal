package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/game"
	"github.com/jonesrussell/gimbal/internal/logger"
	"go.uber.org/zap"
)

// ExitCode represents the program's exit status
type ExitCode int

const (
	ExitSuccess ExitCode = 0
	ExitFailure ExitCode = 1
)

// run executes the main game logic and returns an error if something goes wrong
func run() error {
	// Force stdout to be unbuffered
	os.Stdout.Sync()

	// Set debug log level if not set
	if os.Getenv("LOG_LEVEL") == "" {
		os.Setenv("LOG_LEVEL", "DEBUG")
	}

	// Log system information
	logger.GlobalLogger.Info("Starting game",
		zap.String("goos", runtime.GOOS),
		zap.String("goarch", runtime.GOARCH),
		zap.Int("num_cpu", runtime.NumCPU()),
		zap.String("go_version", runtime.Version()),
		zap.String("log_level", os.Getenv("LOG_LEVEL")),
	)

	// Create game configuration with options
	config := common.NewConfig(
		common.WithDebug(true), // Force debug mode
		common.WithSpeed(common.DefaultSpeed),
		common.WithStarSettings(common.DefaultStarSize, common.DefaultStarSpeed),
		common.WithAngleStep(common.DefaultAngleStep),
	)

	logger.GlobalLogger.Info("Game configuration created",
		zap.Any("screen_size", config.ScreenSize),
		zap.Any("player_size", config.PlayerSize),
		zap.Int("num_stars", config.NumStars),
		zap.Bool("debug", config.Debug),
	)

	// Initialize game
	g, err := game.New(config)
	if err != nil {
		return fmt.Errorf("failed to initialize game: %w", err)
	}

	logger.GlobalLogger.Info("Game initialized successfully")

	// Run game
	return g.Run()
}

func main() {
	var exitCode ExitCode

	// Ensure logger is flushed on exit
	defer func() {
		if err := logger.Sync(); err != nil {
			fmt.Printf("Failed to sync logger: %v\n", err)
		}
		os.Exit(int(exitCode))
	}()

	if err := run(); err != nil {
		logger.GlobalLogger.Error("Game error", zap.Error(err))
		exitCode = ExitFailure
	}
}
