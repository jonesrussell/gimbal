package main

import (
	"fmt"
	"os"

	"go.uber.org/dig"
	"go.uber.org/zap"

	"github.com/jonesrussell/gimbal/internal/config"
	"github.com/jonesrussell/gimbal/internal/engine"
	"github.com/jonesrussell/gimbal/internal/game"
)

func main() {
	container := dig.New()

	// Provide logger
	if err := container.Provide(func() (*zap.Logger, error) {
		if os.Getenv("DEBUG") == "true" {
			return zap.NewDevelopment()
		}
		return zap.NewProduction()
	}); err != nil {
		fmt.Printf("Failed to provide logger: %v\n", err)
		os.Exit(1)
	}

	// Provide config manager
	if err := container.Provide(func(logger *zap.Logger) (*config.Manager, error) {
		env := os.Getenv("ENV")
		if env == "" {
			env = "development"
		}
		manager, err := config.NewManager(logger, env)
		if err != nil {
			return nil, err
		}
		if err := manager.Load(); err != nil {
			return nil, err
		}
		return manager, nil
	}); err != nil {
		fmt.Printf("Failed to provide config manager: %v\n", err)
		os.Exit(1)
	}

	// Provide config for backward compatibility
	if err := container.Provide(func(manager *config.Manager) *config.Config {
		return manager.Get()
	}); err != nil {
		fmt.Printf("Failed to provide config: %v\n", err)
		os.Exit(1)
	}

	// Provide game state
	if err := container.Provide(game.NewGimlarGame); err != nil {
		fmt.Printf("Failed to provide game state: %v\n", err)
		os.Exit(1)
	}

	// Provide game engine
	if err := container.Provide(func(logger *zap.Logger, cfg *config.Config, gameState *game.GimlarGame) (*engine.Game, error) {
		return engine.NewGame(logger, cfg, gameState)
	}); err != nil {
		fmt.Printf("Failed to provide game engine: %v\n", err)
		os.Exit(1)
	}

	// Run the game engine
	if err := container.Invoke(func(g *engine.Game) error {
		return g.Run()
	}); err != nil {
		fmt.Printf("Failed to run game: %v\n", err)
		os.Exit(1)
	}
}
