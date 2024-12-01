package main

import (
	"fmt"
	"os"

	"go.uber.org/dig"
	"go.uber.org/zap"

	"github.com/jonesrussell/gimbal/game"
)

func main() {
	// Initialize logger first
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	// Create DI container
	container := dig.New()

	// Provide logger
	if err := container.Provide(func() *zap.Logger { return logger }); err != nil {
		logger.Fatal("Failed to provide logger", zap.Error(err))
	}

	// Provide game configuration
	if err := container.Provide(func() *game.Config {
		return &game.Config{
			Speed: 0.04,
		}
	}); err != nil {
		logger.Fatal("Failed to provide game config", zap.Error(err))
	}

	// Provide game instance
	if err := container.Provide(game.NewGimlarGame); err != nil {
		logger.Fatal("Failed to provide game instance", zap.Error(err))
	}

	// Invoke game runner
	if err := container.Invoke(func(g *game.GimlarGame) error {
		return g.Run()
	}); err != nil {
		logger.Fatal("Failed to run game", zap.Error(err))
	}
}
