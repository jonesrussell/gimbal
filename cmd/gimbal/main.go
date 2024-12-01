package main

import (
	"fmt"
	"os"

	"go.uber.org/dig"
	"go.uber.org/zap"

	"github.com/jonesrussell/gimbal/internal/config"
	"github.com/jonesrussell/gimbal/internal/engine"
)

func main() {
	// Create DI container
	container := dig.New()

	// Provide logger
	if err := container.Provide(func() (*zap.Logger, error) {
		return zap.NewProduction()
	}); err != nil {
		fmt.Printf("Failed to provide logger: %v\n", err)
		os.Exit(1)
	}

	// Provide config
	if err := container.Provide(func(logger *zap.Logger) (*config.Config, error) {
		config.Init(logger)
		return config.Load("development")
	}); err != nil {
		fmt.Printf("Failed to provide config: %v\n", err)
		os.Exit(1)
	}

	// Provide game engine
	if err := container.Provide(engine.NewGame); err != nil {
		fmt.Printf("Failed to provide game engine: %v\n", err)
		os.Exit(1)
	}

	// Run the game
	if err := container.Invoke(func(game *engine.Game) error {
		return game.Run()
	}); err != nil {
		fmt.Printf("Failed to run game: %v\n", err)
		os.Exit(1)
	}
}
