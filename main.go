package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"

	"github.com/jonesrussell/gimbal/internal/app"
	"github.com/jonesrussell/gimbal/internal/config"
	"github.com/jonesrussell/gimbal/internal/errors"
	"github.com/jonesrussell/gimbal/internal/game"

	// Scene packages - explicit imports for registration

	"github.com/jonesrussell/gimbal/internal/scenes/credits"
	"github.com/jonesrussell/gimbal/internal/scenes/gameover"
	"github.com/jonesrussell/gimbal/internal/scenes/gameplay"
	"github.com/jonesrussell/gimbal/internal/scenes/intro"
	"github.com/jonesrussell/gimbal/internal/scenes/mainmenu"
	"github.com/jonesrussell/gimbal/internal/scenes/pause"
	"github.com/jonesrussell/gimbal/internal/scenes/stageintro"
	"github.com/jonesrussell/gimbal/internal/scenes/stagetransition"
	"github.com/jonesrussell/gimbal/internal/scenes/victory"
)

// ExitCode represents the program's exit status
type ExitCode int

const (
	ExitSuccess ExitCode = 0
	ExitFailure ExitCode = 1
)

// Application represents the main application
type Application struct {
	container  *app.Container
	config     *config.AppConfig
	invincible bool
}

// NewApplication creates a new application instance
func NewApplication(invincible bool) (*Application, error) {
	cfg, err := config.LoadAppConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	// Validate configuration
	if err = cfg.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	container := app.NewContainer(cfg, invincible)

	return &Application{
		container:  container,
		config:     cfg,
		invincible: invincible,
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
			errors.SystemInitFailed,
			"failed to initialize application container",
			err,
		)
	}

	return nil
}

// Run starts the application
func (a *Application) Run() error {
	gameInstance := a.container.GetGame()

	// Log system information
	a.logSystemInfo()

	// Configure and run the game
	gameRunner := game.NewRunner(a.config.Game, gameInstance)

	if err := gameRunner.Run(); err != nil {
		// Check if the error is audio-related (ALSA, oto, etc.)
		// Audio is optional - if audio initialization fails, we should
		// log a warning but not fail the entire game
		errMsg := err.Error()
		if strings.Contains(strings.ToLower(errMsg), "alsa") ||
			strings.Contains(strings.ToLower(errMsg), "oto") ||
			strings.Contains(strings.ToLower(errMsg), "audio") ||
			strings.Contains(strings.ToLower(errMsg), "pulse") ||
			strings.Contains(strings.ToLower(errMsg), "jack") {
			// Audio is optional - log warning but don't fail
			log.Printf("[WARN] Audio initialization failed (audio is optional), game will run without audio: %v", err)
			// Return nil to indicate success (game can run without audio)
			return nil
		}
		// For non-audio errors, return them as system init failures
		return errors.NewGameErrorWithCause(
			errors.SystemInitFailed,
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

	// Auto-disable audio in containers or WSL (unless explicitly enabled)
	// WSL often has no ALSA device; avoid ALSA errors and game exit by disabling audio upfront
	if os.Getenv("DISABLE_AUDIO") == "" && (isContainer() || isWSL()) {
		if err := os.Setenv("DISABLE_AUDIO", "1"); err != nil {
			return fmt.Errorf("failed to set DISABLE_AUDIO: %w", err)
		}
	}

	return nil
}

// isContainer detects if the application is running in a container
func isContainer() bool {
	// Check for .dockerenv file (Docker)
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}

	// Check cgroup for container indicators
	if cgroup, err := os.ReadFile("/proc/self/cgroup"); err == nil {
		cgroupStr := string(cgroup)
		if strings.Contains(cgroupStr, "docker") ||
			strings.Contains(cgroupStr, "containerd") ||
			strings.Contains(cgroupStr, "kubepods") ||
			strings.Contains(cgroupStr, "container") {
			return true
		}
	}

	return false
}

// isWSL reports whether the process is running under Windows Subsystem for Linux.
// WSL typically has no ALSA sound device, so we disable audio to avoid init failures.
func isWSL() bool {
	data, err := os.ReadFile("/proc/version")
	if err != nil {
		return false
	}
	v := strings.ToLower(string(data))
	return strings.Contains(v, "microsoft") || strings.Contains(v, "wsl")
}

// logSystemInfo logs system and runtime information
func (a *Application) logSystemInfo() {
	info := a.config.GetSystemInfo()
	log.Printf("[INFO] Starting Gimbal version=%s goos=%s goarch=%s num_cpu=%d %dx%d",
		info.Version, info.GOOS, info.GOARCH, info.NumCPU,
		a.config.Game.WindowWidth, a.config.Game.WindowHeight)
}

// registerScenes explicitly registers all scene factories.
// This replaces the implicit init()-based registration for clean architecture.
func registerScenes() {
	intro.Register()
	mainmenu.Register()
	gameplay.Register()
	pause.Register()
	gameover.Register()
	// New scene registrations
	credits.Register()
	stageintro.Register()
	stagetransition.Register()
	victory.Register()
}

// run executes the main application logic
func run() error {
	// Explicitly register all scenes (Clean Architecture approach)
	registerScenes()

	// Parse command-line flags
	invincible := flag.Bool("invincible", false, "Enable player invincibility (only works when DEBUG=true)")
	flag.Parse()

	application, err := NewApplication(*invincible)
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
