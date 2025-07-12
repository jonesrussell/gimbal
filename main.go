package main

import (
	"context"
	"fmt"
	"os"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/jonesrussell/gimbal/internal/app"
	"github.com/jonesrussell/gimbal/internal/errors"
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

	// Start pprof server (only in dev builds)
	app.StartPprofServer()

	// Create and initialize application container
	container := app.NewContainer()

	// Initialize all dependencies
	if err := container.Initialize(context.Background()); err != nil {
		return errors.NewGameErrorWithCause(
			errors.ErrorCodeSystemFailed,
			"failed to initialize application container",
			err,
		)
	}

	// Get dependencies from container
	logger := container.GetLogger()
	config := container.GetConfig()
	game := container.GetGame()

	// Ensure graceful shutdown
	defer func() {
		if err := container.Shutdown(context.Background()); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to shutdown container: %v\n", err)
		}
	}()

	// Log system information
	logger.Info("Starting game",
		"goos", runtime.GOOS,
		"goarch", runtime.GOARCH,
		"num_cpu", runtime.NumCPU(),
		"go_version", runtime.Version(),
		"log_level", os.Getenv("LOG_LEVEL"),
	)

	// Run game with Ebiten
	ebiten.SetWindowSize(config.ScreenSize.Width, config.ScreenSize.Height)
	ebiten.SetWindowTitle("Gimbal - ECS Version")
	ebiten.SetTPS(60)

	if runErr := ebiten.RunGame(game); runErr != nil {
		return errors.NewGameErrorWithCause(errors.ErrorCodeSystemFailed, "game error", runErr)
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
