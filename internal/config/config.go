package config

import (
	"encoding/json"
	"fmt"
	"os"

	"go.uber.org/zap"
)

type Config struct {
	Screen struct {
		Width  int    `json:"width"`
		Height int    `json:"height"`
		Title  string `json:"title"`
	} `json:"screen"`
	Game struct {
		NumStars int  `json:"numStars"`
		Debug    bool `json:"debug"`
	} `json:"game"`
}

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

func New() *Config {
	if config == nil {
		var err error
		config, err = Load("development") // Default to development
		if err != nil {
			logger.Fatal("Failed to load config", zap.Error(err))
		}
	}
	return config
}
