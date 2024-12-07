package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/dig"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/jonesrussell/gimbal/internal/config"
	"github.com/jonesrussell/gimbal/internal/core"
	"github.com/jonesrussell/gimbal/internal/engine"
	"github.com/jonesrussell/gimbal/internal/game"
)

func main() {
	// Create root context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Create DI container
	container := dig.New()
	g, ctx := errgroup.WithContext(ctx)

	// Create logger instance
	logger, err := initLogger()
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	// Provide logger
	if err := container.Provide(func() *zap.Logger {
		return logger
	}); err != nil {
		logger.Fatal("Failed to provide logger", zap.Error(err))
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
		logger.Fatal("Failed to provide config manager", zap.Error(err))
	}

	// Provide config for backward compatibility
	if err := container.Provide(func(manager *config.Manager) *config.Config {
		return manager.Get()
	}); err != nil {
		logger.Fatal("Failed to provide config", zap.Error(err))
	}

	// Provide game state
	if err := container.Provide(game.NewGimlarGame); err != nil {
		logger.Fatal("Failed to provide game state", zap.Error(err))
	}

	// Provide asset manager
	if err := container.Provide(func(logger *zap.Logger, manager *config.Manager) (core.AssetManager, error) {
		cfg := manager.Get()
		return core.NewAssetManagerImpl(logger,
			core.WithBaseDir("assets"),
			core.WithCacheSize(1000),
			core.WithSound(true),
			core.WithAssetDebug(cfg.Game.Debug),
		)
	}); err != nil {
		logger.Fatal("Failed to provide asset manager", zap.Error(err))
	}

	// Provide game engine
	if err := container.Provide(func(
		logger *zap.Logger,
		cfg *config.Config,
		gameState *game.GimlarGame,
		assets core.AssetManager,
	) (*engine.Game, error) {
		return engine.NewGame(logger, cfg, gameState, assets)
	}); err != nil {
		logger.Fatal("Failed to provide game engine", zap.Error(err))
	}

	// Add shutdown handler
	g.Go(func() error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case sig := <-sigChan:
			logger.Info("received shutdown signal", zap.String("signal", sig.String()))
			cancel()
			return nil
		}
	})

	// Run the game engine
	if err := container.Invoke(func(game *engine.Game) error {
		return game.Run()
	}); err != nil {
		logger.Error("failed to run game", zap.Error(err))
		os.Exit(1)
	}

	// Wait for all goroutines to complete
	if err := g.Wait(); err != nil {
		logger.Error("error during shutdown", zap.Error(err))
		os.Exit(1)
	}
}

func initLogger() (*zap.Logger, error) {
	if os.Getenv("DEBUG") == "true" {
		return zap.NewDevelopment()
	}
	return zap.NewProduction()
}
