package managers

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jonesrussell/gimbal/assets"
)

// PlayerConfig defines player entity configuration
type PlayerConfig struct {
	Health                int     `json:"health"`
	Size                  int     `json:"size"`
	SpriteName            string  `json:"sprite_name"`
	InvincibilityDuration float64 `json:"invincibility_duration"` // Duration in seconds (JSON)
}

// GetInvincibilityDuration returns the invincibility duration as time.Duration
func (pc *PlayerConfig) GetInvincibilityDuration() time.Duration {
	return time.Duration(pc.InvincibilityDuration * float64(time.Second))
}

// EnemyTypeConfig defines configuration for a single enemy type (JSON representation)
// Uses int for MovementPattern to avoid import cycles
type EnemyTypeConfig struct {
	Type         string  `json:"type"` // "basic", "heavy", "boss"
	Health       int     `json:"health"`
	Speed        float64 `json:"speed"`
	Size         int     `json:"size"`
	Points       int     `json:"points"`
	SpriteName   string  `json:"sprite_name"`
	MovementType string  `json:"movement_type"` // "outward", "spiral", "orbital"
	// MovementPattern as int (0=normal, 1=zigzag, 2=accelerating, 3=pulsing)
	MovementPattern int     `json:"movement_pattern"`
	CanShoot        bool    `json:"can_shoot"`
	FireRate        float64 `json:"fire_rate"` // Shots per second
	ProjectileSpeed float64 `json:"projectile_speed"`
}

// EnemyConfigs contains all enemy type configurations
type EnemyConfigs struct {
	EnemyTypes []EnemyTypeConfig `json:"enemy_types"`
}

// LoadPlayerConfig loads player configuration from the embedded assets
func LoadPlayerConfig(ctx context.Context) (*PlayerConfig, error) {
	data, err := assets.Assets.ReadFile("entities/player.json")
	if err != nil {
		return defaultPlayerConfig(), nil
	}

	var config PlayerConfig
	if unmarshalErr := json.Unmarshal(data, &config); unmarshalErr != nil {
		return defaultPlayerConfig(), nil
	}

	return &config, nil
}

// defaultPlayerConfig returns default player configuration
func defaultPlayerConfig() *PlayerConfig {
	return &PlayerConfig{
		Health:                3,
		Size:                  32,
		SpriteName:            "player",
		InvincibilityDuration: 2.0,
	}
}
