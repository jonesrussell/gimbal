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

	// Provide config
	if err := container.Provide(func(logger *zap.Logger) (*config.Config, error) {
		config.Init(logger)
		return config.Load("development")
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
