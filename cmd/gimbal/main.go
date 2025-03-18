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

func main() {
	// Force stdout to be unbuffered
	os.Stdout.Sync()

	// Set debug log level if not set
	if os.Getenv("LOG_LEVEL") == "" {
		os.Setenv("LOG_LEVEL", "DEBUG")
	}

	// Log system information
	logger.GlobalLogger.Debug("System information",
		zap.String("goos", runtime.GOOS),
		zap.String("goarch", runtime.GOARCH),
		zap.Int("num_cpu", runtime.NumCPU()),
		zap.String("go_version", runtime.Version()),
	)

	logger.GlobalLogger.Debug("Environment variables",
		zap.String("debug", os.Getenv("DEBUG")),
		zap.String("log_level", os.Getenv("LOG_LEVEL")),
		zap.String("goos", os.Getenv("GOOS")),
	)

	// Create game configuration with options
	config := common.NewConfig(
		common.WithDebug(true), // Force debug mode
		common.WithSpeed(common.DefaultSpeed),
		common.WithStarSettings(common.DefaultStarSize, common.DefaultStarSpeed),
		common.WithAngleStep(common.DefaultAngleStep),
	)

	logger.GlobalLogger.Debug("Game configuration created",
		zap.Any("screen_size", config.ScreenSize),
		zap.Any("player_size", config.PlayerSize),
		zap.Int("num_stars", config.NumStars),
		zap.Bool("debug", config.Debug),
	)

	// Initialize game
	g, initErr := game.New(config)
	if initErr != nil {
		logger.GlobalLogger.Error("Failed to initialize game",
			zap.Error(initErr),
			zap.String("error_type", fmt.Sprintf("%T", initErr)),
		)
		os.Exit(1)
	}

	logger.GlobalLogger.Debug("Game initialized successfully")

	// Run game
	if runErr := g.Run(); runErr != nil {
		logger.GlobalLogger.Error("Failed to run game",
			zap.Error(runErr),
			zap.String("error_type", fmt.Sprintf("%T", runErr)),
		)
		os.Exit(1)
	}
}
