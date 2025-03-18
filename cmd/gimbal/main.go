package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/game"
	"github.com/jonesrussell/gimbal/internal/logger"
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

	// Create logger
	log, err := logger.New()
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}

	// Ensure logger is flushed on exit
	defer func() {
		if err := log.Sync(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to sync logger: %v\n", err)
		}
	}()

	// Log system information
	log.Info("Starting game",
		"goos", runtime.GOOS,
		"goarch", runtime.GOARCH,
		"num_cpu", runtime.NumCPU(),
		"go_version", runtime.Version(),
		"log_level", os.Getenv("LOG_LEVEL"),
	)

	// Create game configuration with options
	config := common.NewConfig(
		common.WithDebug(true), // Force debug mode
		common.WithSpeed(common.DefaultSpeed),
		common.WithStarSettings(common.DefaultStarSize, common.DefaultStarSpeed),
		common.WithAngleStep(common.DefaultAngleStep),
	)

	log.Info("Game configuration created",
		"screen_size", config.ScreenSize,
		"player_size", config.PlayerSize,
		"num_stars", config.NumStars,
		"debug", config.Debug,
	)

	// Initialize game
	g, err := game.New(config, log)
	if err != nil {
		return fmt.Errorf("failed to initialize game: %w", err)
	}

	log.Info("Game initialized successfully")

	// Run game
	return g.Run()
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
