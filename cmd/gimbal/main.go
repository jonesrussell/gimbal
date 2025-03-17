package main

import (
	"log/slog"
	"os"

	"github.com/jonesrussell/gimbal/internal/game"
)

func main() {
	// Initialize logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Create game configuration with options
	config := game.NewConfig(
		game.WithDebug(os.Getenv("DEBUG") != ""),
		game.WithSpeed(0.04),
		game.WithStarSettings(5.0, 2.0),
		game.WithAngleStep(0.05),
	)

	// Create input handler
	input := &game.InputHandler{}

	// Initialize game
	g, err := game.NewGimlarGame(config, input)
	if err != nil {
		logger.Error("Failed to initialize game", "error", err)
		os.Exit(1)
	}

	// Run game
	if err := g.Run(); err != nil {
		logger.Error("Failed to run game", "error", err)
		os.Exit(1)
	}
}
