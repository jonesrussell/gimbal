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
		game.WithSpeed(game.DefaultSpeed),
		game.WithStarSettings(game.DefaultStarSize, game.DefaultStarSpeed),
		game.WithAngleStep(game.DefaultAngleStep),
	)

	// Create input handler
	input := &game.InputHandler{}

	// Initialize game
	g, initErr := game.NewGimlarGame(config, input)
	if initErr != nil {
		logger.Error("Failed to initialize game", "error", initErr)
		os.Exit(1)
	}

	// Run game
	if runErr := g.Run(); runErr != nil {
		logger.Error("Failed to run game", "error", runErr)
		os.Exit(1)
	}
}
