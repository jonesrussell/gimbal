package main

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"github.com/jonesrussell/gimbal/internal/app"
	"github.com/jonesrussell/gimbal/internal/config"
	"github.com/jonesrussell/gimbal/internal/errors"
	"github.com/jonesrussell/gimbal/internal/game"
)

// ExitCode represents the program's exit status
type ExitCode int

const (
	ExitSuccess ExitCode = 0
	ExitFailure ExitCode = 1
)

// Application represents the main application
type Application struct {
	container *app.Container
	config    *config.AppConfig
}

// NewApplication creates a new application instance
func NewApplication() (*Application, error) {
	cfg, err := config.LoadAppConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	// Validate configuration
	if err = cfg.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	container := app.NewContainer(cfg)

	return &Application{
		container: container,
		config:    cfg,
	}, nil
}

// Initialize initializes the application
func (a *Application) Initialize(ctx context.Context) error {
	// Setup environment if needed
	if err := a.setupEnvironment(); err != nil {
		return fmt.Errorf("failed to setup environment: %w", err)
	}

	// Start profiling server in development
	if a.config.IsDevelopment() {
		app.StartPprofServer()
	}

	// Initialize container
	if err := a.container.Initialize(ctx); err != nil {
		return errors.NewGameErrorWithCause(
			errors.ErrorCodeSystemFailed,
			"failed to initialize application container",
			err,
		)
	}

	return nil
}

// Run starts the application
func (a *Application) Run() error {
	logger := a.container.GetLogger()
	gameInstance := a.container.GetGame()

	// Log system information
	a.logSystemInfo(logger)

	// Configure and run the game
	gameRunner := game.NewRunner(a.config.Game, gameInstance)

	if err := gameRunner.Run(); err != nil {
		return errors.NewGameErrorWithCause(
			errors.ErrorCodeSystemFailed,
			"game execution failed",
			err,
		)
	}

	return nil
}

// Shutdown gracefully shuts down the application
func (a *Application) Shutdown(ctx context.Context) error {
	if a.container != nil {
		return a.container.Shutdown(ctx)
	}
	return nil
}

// setupEnvironment configures the runtime environment
func (a *Application) setupEnvironment() error {
	// Set default log level if not configured
	if os.Getenv("LOG_LEVEL") == "" && a.config.LogLevel != "" {
		if err := os.Setenv("LOG_LEVEL", a.config.LogLevel); err != nil {
			return fmt.Errorf("failed to set log level: %w", err)
		}
	}

	return nil
}

// logSystemInfo logs system and runtime information
func (a *Application) logSystemInfo(logger interface{ Info(string, ...interface{}) }) {
	info := a.config.GetSystemInfo()
	logger.Info("Starting Gimbal",
		"version", info.Version,
		"goos", info.GOOS,
		"goarch", info.GOARCH,
		"num_cpu", info.NumCPU,
		"go_version", info.GoVersion,
		"log_level", info.LogLevel,
		"window_width", a.config.Game.WindowWidth,
		"window_height", a.config.Game.WindowHeight,
	)
}

// run executes the main application logic
func run() error {
	application, err := NewApplication()
	if err != nil {
		return err
	}

	ctx := context.Background()

	// Initialize application
	initErr := application.Initialize(ctx)
	if initErr != nil {
		return initErr
	}

	// Ensure graceful shutdown
	defer func() {
		if shutdownErr := application.Shutdown(ctx); shutdownErr != nil {
			fmt.Fprintf(os.Stderr, "Failed to shutdown application: %v\n", shutdownErr)
		}
	}()

	// Run the application
	return application.Run()
}

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Fprintf(os.Stderr, "Error loading .env file: %v\n", err)
	}

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(int(ExitFailure))
	}
	os.Exit(int(ExitSuccess))
}
