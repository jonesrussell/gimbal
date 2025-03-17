package main

import (
	"log/slog"
	"os"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/game"
)

func main() {
	// Initialize logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Create game configuration with options
	config := common.NewConfig(
		common.WithDebug(os.Getenv("DEBUG") != ""),
		common.WithSpeed(common.DefaultSpeed),
		common.WithStarSettings(common.DefaultStarSize, common.DefaultStarSpeed),
		common.WithAngleStep(common.DefaultAngleStep),
	)

	// Initialize game
	g, initErr := game.New(config)
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
