package managers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jonesrussell/gimbal/assets"
	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/errors"
)

// LoadPlayerConfig loads player configuration from embedded assets
// Returns an error if the file is missing or invalid (no fallback)
func LoadPlayerConfig(ctx context.Context, logger common.Logger) (*PlayerConfig, error) {
	// Check for cancellation
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Load from embedded assets
	data, err := assets.Assets.ReadFile("entities/player.json")
	if err != nil {
		return nil, errors.NewGameErrorWithCause(
			errors.AssetNotFound,
			"failed to read player.json from embedded assets",
			err,
		)
	}

	var config PlayerConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, errors.NewGameErrorWithCause(
			errors.AssetInvalid,
			"failed to parse player.json",
			err,
		)
	}

	// Validate required fields
	if config.Health <= 0 {
		return nil, errors.NewGameError(
			errors.AssetInvalid,
			fmt.Sprintf("invalid player health: %d (must be > 0)", config.Health),
		)
	}
	if config.Size <= 0 {
		return nil, errors.NewGameError(
			errors.AssetInvalid,
			fmt.Sprintf("invalid player size: %d (must be > 0)", config.Size),
		)
	}
	if config.SpriteName == "" {
		return nil, errors.NewGameError(
			errors.AssetInvalid,
			"player sprite_name is required",
		)
	}
	if config.InvincibilityDuration <= 0 {
		return nil, errors.NewGameError(
			errors.AssetInvalid,
			fmt.Sprintf("invalid invincibility_duration: %f (must be > 0)", config.InvincibilityDuration),
		)
	}

	if logger != nil {
		logger.Info("Player config loaded from JSON", "health", config.Health, "size", config.Size)
	}

	return &config, nil
}

// LoadEnemyConfigs loads all enemy type configurations from embedded assets
// Returns an error if the file is missing or invalid (no fallback)
func LoadEnemyConfigs(ctx context.Context, logger common.Logger) (*EnemyConfigs, error) {
	// Check for cancellation
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Load from embedded assets
	data, err := assets.Assets.ReadFile("entities/enemies.json")
	if err != nil {
		return nil, errors.NewGameErrorWithCause(
			errors.AssetNotFound,
			"failed to read enemies.json from embedded assets",
			err,
		)
	}

	var configs EnemyConfigs
	if err := json.Unmarshal(data, &configs); err != nil {
		return nil, errors.NewGameErrorWithCause(
			errors.AssetInvalid,
			"failed to parse enemies.json",
			err,
		)
	}

	// Validate required fields
	if len(configs.EnemyTypes) == 0 {
		return nil, errors.NewGameError(
			errors.AssetInvalid,
			"enemies.json must contain at least one enemy type",
		)
	}

	// Validate each enemy type
	for i, enemyType := range configs.EnemyTypes {
		if enemyType.Type == "" {
			return nil, errors.NewGameError(
				errors.AssetInvalid,
				fmt.Sprintf("enemy type at index %d has empty type field", i),
			)
		}
		if enemyType.Health <= 0 {
			return nil, errors.NewGameError(
				errors.AssetInvalid,
				fmt.Sprintf("enemy type '%s' has invalid health: %d (must be > 0)", enemyType.Type, enemyType.Health),
			)
		}
		if enemyType.Speed <= 0 {
			return nil, errors.NewGameError(
				errors.AssetInvalid,
				fmt.Sprintf("enemy type '%s' has invalid speed: %f (must be > 0)", enemyType.Type, enemyType.Speed),
			)
		}
		if enemyType.Size <= 0 {
			return nil, errors.NewGameError(
				errors.AssetInvalid,
				fmt.Sprintf("enemy type '%s' has invalid size: %d (must be > 0)", enemyType.Type, enemyType.Size),
			)
		}
		if enemyType.SpriteName == "" {
			return nil, errors.NewGameError(
				errors.AssetInvalid,
				fmt.Sprintf("enemy type '%s' has empty sprite_name", enemyType.Type),
			)
		}
		if enemyType.MovementType == "" {
			return nil, errors.NewGameError(
				errors.AssetInvalid,
				fmt.Sprintf("enemy type '%s' has empty movement_type", enemyType.Type),
			)
		}
	}

	if logger != nil {
		logger.Info("Enemy configs loaded from JSON", "count", len(configs.EnemyTypes))
	}

	return &configs, nil
}
