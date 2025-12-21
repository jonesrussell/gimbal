package config

import (
	"fmt"
	"os"
	"runtime"

	"github.com/kelseyhightower/envconfig"
)

// AppConfig holds all application configuration
type AppConfig struct {
	LogLevel string `envconfig:"LOG_LEVEL" default:"DEBUG"`
	Logging  *LoggingConfig
	Game     *AppGameConfig
	Debug    *DebugConfig
}

// LoggingConfig holds logging-specific configuration
type LoggingConfig struct {
	LogFile    string `envconfig:"LOG_FILE" default:"logs/gimbal.log"`
	ConsoleOut bool   `envconfig:"LOG_CONSOLE_OUT" default:"true"`
	FileOut    bool   `envconfig:"LOG_FILE_OUT" default:"true"`
}

// AppGameConfig holds application-level game configuration
type AppGameConfig struct {
	WindowWidth  int    `envconfig:"GAME_WINDOW_WIDTH" default:"1280"`
	WindowHeight int    `envconfig:"GAME_WINDOW_HEIGHT" default:"720"`
	WindowTitle  string `envconfig:"GAME_TITLE" default:"Gimbal - ECS Version"`
	TPS          int    `envconfig:"GAME_TPS" default:"60"`
	Resizable    bool   `envconfig:"GAME_RESIZABLE" default:"true"`
	DefaultScene string `envconfig:"GAME_DEFAULT_SCENE" default:"menu"`
}

// DebugConfig holds debug-specific configuration
type DebugConfig struct {
	PprofPort  int  `envconfig:"DEBUG_PPROF_PORT" default:"6060"`
	Enabled    bool `envconfig:"DEBUG" default:"true"`
	ShowFPS    bool `envconfig:"DEBUG_SHOW_FPS" default:"false"`
	ShowMemory bool `envconfig:"DEBUG_SHOW_MEMORY" default:"false"`
}

// SystemInfo holds system information
type SystemInfo struct {
	Version   string
	GOOS      string
	GOARCH    string
	NumCPU    int
	GoVersion string
	LogLevel  string
}

// LoadAppConfig loads application configuration from environment and defaults
func LoadAppConfig() (*AppConfig, error) {
	config := &AppConfig{}

	// Load main app config
	if err := envconfig.Process("", config); err != nil {
		return nil, fmt.Errorf("failed to process app config: %w", err)
	}

	// Load nested configs
	config.Logging = &LoggingConfig{}
	if err := envconfig.Process("", config.Logging); err != nil {
		return nil, fmt.Errorf("failed to process logging config: %w", err)
	}

	config.Game = &AppGameConfig{}
	if err := envconfig.Process("", config.Game); err != nil {
		return nil, fmt.Errorf("failed to process game config: %w", err)
	}

	config.Debug = &DebugConfig{}
	if err := envconfig.Process("", config.Debug); err != nil {
		return nil, fmt.Errorf("failed to process debug config: %w", err)
	}

	return config, nil
}

// IsDevelopment returns true if running in development mode
func (c *AppConfig) IsDevelopment() bool {
	return c.Debug.Enabled
}

// GetSystemInfo returns current system information
func (c *AppConfig) GetSystemInfo() *SystemInfo {
	return &SystemInfo{
		Version:   getEnvWithDefault("APP_VERSION", "dev"),
		GOOS:      runtime.GOOS,
		GOARCH:    runtime.GOARCH,
		NumCPU:    runtime.NumCPU(),
		GoVersion: runtime.Version(),
		LogLevel:  c.LogLevel,
	}
}

// Validate validates the application configuration
func (c *AppConfig) Validate() error {
	if err := c.Logging.Validate(); err != nil {
		return err
	}

	if err := c.Game.Validate(); err != nil {
		return err
	}

	if err := c.Debug.Validate(); err != nil {
		return err
	}

	return nil
}

// Validate validates logging configuration
func (lc *LoggingConfig) Validate() error {
	if lc.LogFile == "" {
		return fmt.Errorf("log file path cannot be empty")
	}

	if !lc.ConsoleOut && !lc.FileOut {
		return fmt.Errorf("at least one output destination must be enabled")
	}

	return nil
}

// Validate validates application game configuration
func (gc *AppGameConfig) Validate() error {
	if gc.WindowWidth <= 0 || gc.WindowHeight <= 0 {
		return fmt.Errorf("invalid window dimensions: %dx%d", gc.WindowWidth, gc.WindowHeight)
	}

	if gc.TPS <= 0 || gc.TPS > 120 {
		return fmt.Errorf("invalid TPS: %d (must be 1-120)", gc.TPS)
	}

	if gc.WindowTitle == "" {
		return fmt.Errorf("window title cannot be empty")
	}

	return nil
}

// Validate validates debug configuration
func (dc *DebugConfig) Validate() error {
	if dc.PprofPort <= 0 || dc.PprofPort > 65535 {
		return fmt.Errorf("invalid pprof port: %d", dc.PprofPort)
	}

	return nil
}

// Helper function for system info (still needed for APP_VERSION)
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
