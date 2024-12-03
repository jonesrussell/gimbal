package config

import (
	"encoding/json"
	"fmt"
	"os"

	"go.uber.org/zap"
)

// Manager handles configuration loading and management
type Manager struct {
	logger *zap.Logger
	config *Config
	env    string
}

// NewManager creates a new configuration manager
func NewManager(logger *zap.Logger, env string) (*Manager, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger is required")
	}

	if env == "" {
		env = "development" // Default environment
	}

	return &Manager{
		logger: logger,
		env:    env,
	}, nil
}

// Load reads and parses the configuration file
func (m *Manager) Load() error {
	configPath := fmt.Sprintf("internal/config/config.%s.json", m.env)
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("reading config file: %w", err)
	}

	config := &Config{}
	if err := json.Unmarshal(data, config); err != nil {
		return fmt.Errorf("parsing config file: %w", err)
	}

	m.config = config
	m.logger.Info("Configuration loaded",
		zap.String("environment", m.env),
		zap.Int("screen.width", config.Screen.Width),
		zap.Int("screen.height", config.Screen.Height))

	return nil
}

// Get returns the current configuration
func (m *Manager) Get() *Config {
	return m.config
}

// For backward compatibility during transition
func New() (*Config, error) {
	logConfig := zap.NewDevelopmentConfig()
	logger, err := logConfig.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	manager, err := NewManager(logger, "development")
	if err != nil {
		return nil, fmt.Errorf("failed to create config manager: %w", err)
	}

	if err := manager.Load(); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return manager.Get(), nil
}
