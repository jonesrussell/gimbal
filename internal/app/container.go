package app

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/config"
	gamepkg "github.com/jonesrussell/gimbal/internal/game"
	"github.com/jonesrussell/gimbal/internal/input"
)

// Container manages all application dependencies and their lifecycle
type Container struct {
	mu sync.RWMutex

	// Core dependencies
	appConfig    *config.AppConfig
	config       *config.GameConfig
	inputHandler common.GameInputHandler
	game         *gamepkg.ECSGame
	invincible   bool

	// State
	initialized bool
	shutdown    bool
}

// NewContainer creates a new application dependency container with the provided configuration
func NewContainer(appConfig *config.AppConfig, invincible bool) *Container {
	return &Container{
		appConfig:  appConfig,
		invincible: invincible,
	}
}

// Initialize sets up all dependencies in the correct order
func (c *Container) Initialize(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.initialized {
		return fmt.Errorf("container already initialized")
	}

	// Step 1: Initialize configuration
	if err := c.initializeConfig(); err != nil {
		return fmt.Errorf("failed to initialize config: %w", err)
	}

	// Step 2: Initialize input handler
	if err := c.initializeInputHandler(); err != nil {
		return fmt.Errorf("failed to initialize input handler: %w", err)
	}

	// Step 3: Initialize game
	if err := c.initializeGame(ctx); err != nil {
		return fmt.Errorf("failed to initialize game: %w", err)
	}

	c.initialized = true
	log.Printf("[INFO] Application container initialized successfully")
	return nil
}

// initializeConfig creates and validates the game configuration
func (c *Container) initializeConfig() error {
	// Use AppConfig's DEBUG value
	debugEnabled := c.appConfig.IsDevelopment()

	opts := []config.GameOption{
		config.WithDebug(debugEnabled),
		config.WithSpeed(config.DefaultSpeed),
		config.WithStarSettings(config.DefaultStarSize, config.DefaultStarSpeed),
		config.WithAngleStep(config.DefaultAngleStep),
	}

	// Only add invincible option if debug is enabled
	if debugEnabled && c.invincible {
		opts = append(opts, config.WithInvincible(true))
		log.Printf("[INFO] Invincible mode enabled (DEBUG mode required)")
	} else if c.invincible && !debugEnabled {
		log.Printf("[WARN] Invincible flag ignored: DEBUG must be true to use invincible mode")
	}

	gameConfig := config.NewConfig(opts...)

	// Safety check: ensure invincible is disabled if debug is disabled
	if !gameConfig.Debug && gameConfig.Invincible {
		gameConfig.Invincible = false
		log.Printf("[WARN] Invincible mode disabled: DEBUG mode is required")
	}

	// Validate configuration
	if err := config.ValidateConfig(gameConfig); err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	c.config = gameConfig
	log.Printf("[INFO] Game configuration created and validated (debug=%v invincible=%v)", gameConfig.Debug, gameConfig.Invincible)
	return nil
}

// initializeInputHandler creates the input handler
func (c *Container) initializeInputHandler() error {
	inputHandler := input.NewHandler()
	c.inputHandler = inputHandler
	log.Printf("[INFO] Input handler created")
	return nil
}

// initializeGame creates the ECS game instance
func (c *Container) initializeGame(ctx context.Context) error {
	game, err := gamepkg.NewECSGame(ctx, c.config, c.inputHandler)
	if err != nil {
		return err
	}
	c.game = game
	log.Printf("[INFO] ECS game initialized successfully")
	return nil
}

// GetGame returns the ECS game instance
func (c *Container) GetGame() *gamepkg.ECSGame {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.game
}

// Shutdown gracefully shuts down all dependencies
func (c *Container) Shutdown(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.shutdown {
		return nil
	}

	log.Printf("[INFO] Shutting down application container")

	// Shutdown in reverse order of initialization
	if c.game != nil {
		if ctx == nil {
			return fmt.Errorf("context must not be nil in Shutdown")
		}
		c.game.Cleanup(ctx)
	}

	c.shutdown = true
	log.Printf("[INFO] Application container shutdown complete")
	return nil
}
