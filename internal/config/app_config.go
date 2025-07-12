package config

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
)

// AppConfig holds all application configuration
type AppConfig struct {
	LogLevel string
	Game     *AppGameConfig
	Debug    *DebugConfig
}

// AppGameConfig holds application-level game configuration
type AppGameConfig struct {
	WindowWidth  int
	WindowHeight int
	WindowTitle  string
	TPS          int
	Resizable    bool
	DefaultScene string
}

// DebugConfig holds debug-specific configuration
type DebugConfig struct {
	Enabled    bool
	PprofPort  int
	ShowFPS    bool
	ShowMemory bool
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
	config := &AppConfig{
		LogLevel: getEnvWithDefault("LOG_LEVEL", "DEBUG"),
		Game:     loadAppGameConfig(),
		Debug:    loadDebugConfig(),
	}

	return config, nil
}

// loadAppGameConfig loads application-level game configuration
func loadAppGameConfig() *AppGameConfig {
	return &AppGameConfig{
		WindowWidth:  getEnvIntWithDefault("GAME_WINDOW_WIDTH", 1280),
		WindowHeight: getEnvIntWithDefault("GAME_WINDOW_HEIGHT", 720),
		WindowTitle:  getEnvWithDefault("GAME_TITLE", "Gimbal - ECS Version"),
		TPS:          getEnvIntWithDefault("GAME_TPS", 60),
		Resizable:    getEnvBoolWithDefault("GAME_RESIZABLE", true),
		DefaultScene: getEnvWithDefault("GAME_DEFAULT_SCENE", "menu"),
	}
}

// loadDebugConfig loads debug configuration
func loadDebugConfig() *DebugConfig {
	return &DebugConfig{
		Enabled:    getEnvBoolWithDefault("DEBUG_ENABLED", true),
		PprofPort:  getEnvIntWithDefault("DEBUG_PPROF_PORT", 6060),
		ShowFPS:    getEnvBoolWithDefault("DEBUG_SHOW_FPS", false),
		ShowMemory: getEnvBoolWithDefault("DEBUG_SHOW_MEMORY", false),
	}
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
	if err := c.Game.Validate(); err != nil {
		return err
	}

	if err := c.Debug.Validate(); err != nil {
		return err
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

// Helper functions for environment variable parsing
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvIntWithDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvBoolWithDefault(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
