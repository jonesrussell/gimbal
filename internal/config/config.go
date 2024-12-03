package config

import (
	"encoding/json"
	"fmt"
	"os"

	"go.uber.org/zap"
)

var (
	config *Config
	logger *zap.Logger
)

func Init(l *zap.Logger) {
	logger = l
}

func Load(env string) (*Config, error) {
	if config != nil {
		return config, nil
	}

	configPath := fmt.Sprintf("internal/config/config.%s.json", env)
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	config = &Config{}
	if err := json.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	logger.Info("Configuration loaded",
		zap.String("environment", env),
		zap.Int("screen.width", config.Screen.Width),
		zap.Int("screen.height", config.Screen.Height))

	return config, nil
}

func New() (*Config, error) {
	if logger == nil {
		logConfig := zap.NewDevelopmentConfig()
		log, err := logConfig.Build()
		if err != nil {
			return nil, fmt.Errorf("failed to initialize logger: %w", err)
		}
		logger = log
	}

	if config == nil {
		var err error
		config, err = Load("development") // Default to development
		if err != nil {
			return nil, fmt.Errorf("failed to load config: %w", err)
		}
	}
	return config, nil
}
