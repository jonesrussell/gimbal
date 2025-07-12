package app

import (
	"context"
	"fmt"
	"sync"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/config"
	gamepkg "github.com/jonesrussell/gimbal/internal/game"
	"github.com/jonesrussell/gimbal/internal/input"
	"github.com/jonesrussell/gimbal/internal/logger"
	"github.com/jonesrussell/gimbal/internal/ui"
)

// Container manages all application dependencies and their lifecycle
type Container struct {
	mu sync.RWMutex

	// Core dependencies
	logger       common.Logger
	config       *config.GameConfig
	inputHandler common.GameInputHandler
	game         *gamepkg.ECSGame

	// State
	initialized bool
	shutdown    bool
}

// NewContainer creates a new application container
func NewContainer() *Container {
	return &Container{}
}

// Initialize sets up all dependencies in the correct order
func (c *Container) Initialize(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.initialized {
		return fmt.Errorf("container already initialized")
	}

	// Step 1: Initialize logger (needed by everything else)
	if err := c.initializeLogger(); err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	// Step 2: Initialize configuration
	if err := c.initializeConfig(); err != nil {
		return fmt.Errorf("failed to initialize config: %w", err)
	}

	// Step 3: Initialize input handler
	if err := c.initializeInputHandler(); err != nil {
		return fmt.Errorf("failed to initialize input handler: %w", err)
	}

	// Step 4: Initialize game
	if err := c.initializeGame(ctx); err != nil {
		return fmt.Errorf("failed to initialize game: %w", err)
	}

	c.initialized = true
	c.logger.Info("Application container initialized successfully")
	return nil
}

// initializeLogger creates and configures the logger
func (c *Container) initializeLogger() error {
	log, err := logger.New()
	if err != nil {
		return err
	}
	c.logger = log
	return nil
}

// initializeConfig creates and validates the game configuration
func (c *Container) initializeConfig() error {
	gameConfig := config.NewConfig(
		config.WithDebug(true), // Force debug mode
		config.WithSpeed(config.DefaultSpeed),
		config.WithStarSettings(config.DefaultStarSize, config.DefaultStarSpeed),
		config.WithAngleStep(config.DefaultAngleStep),
	)

	// Validate configuration
	if err := config.ValidateConfig(gameConfig); err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	c.config = gameConfig
	c.logger.Info("Game configuration created and validated",
		"screen_size", gameConfig.ScreenSize,
		"player_size", gameConfig.PlayerSize,
		"num_stars", gameConfig.NumStars,
		"debug", gameConfig.Debug,
	)
	return nil
}

// initializeInputHandler creates the input handler
func (c *Container) initializeInputHandler() error {
	inputHandler := input.NewHandler()
	c.inputHandler = inputHandler
	c.logger.Info("Input handler created")
	return nil
}

// initializeGame creates the ECS game instance
func (c *Container) initializeGame(ctx context.Context) error {
	// Create UI factory
	uiFactory := &ui.EbitenUIFactory{}

	game, err := gamepkg.NewECSGame(c.config, c.logger, c.inputHandler, uiFactory)
	if err != nil {
		return err
	}
	c.game = game
	c.logger.Info("ECS game initialized successfully")
	return nil
}

// GetLogger returns the logger instance
func (c *Container) GetLogger() common.Logger {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.logger
}

// GetConfig returns the game configuration
func (c *Container) GetConfig() *config.GameConfig {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.config
}

// GetInputHandler removed - dead code

// GetGame returns the ECS game instance
func (c *Container) GetGame() *gamepkg.ECSGame {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.game
}

// IsInitialized removed - dead code

// Shutdown gracefully shuts down all dependencies
func (c *Container) Shutdown(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.shutdown {
		return nil
	}

	c.logger.Info("Shutting down application container")

	// Shutdown in reverse order of initialization
	if c.game != nil {
		if ctx == nil {
			ctx = context.Background()
		}
		c.game.Cleanup(ctx)
		c.logger.Debug("Game cleaned up")
	}

	if c.logger != nil {
		if err := c.logger.Sync(); err != nil {
			c.logger.Error("Failed to sync logger during shutdown", "error", err)
		}
	}

	c.shutdown = true
	c.logger.Info("Application container shutdown complete")
	return nil
}

// SetInputHandler removed - dead code
