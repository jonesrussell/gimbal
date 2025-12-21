package managers

import (
	"context"
	"fmt"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/errors"
)

// LoadPlayerConfig loads player configuration from embedded assets
// Returns an error if the file is missing or invalid (no fallback)
func LoadPlayerConfig(ctx context.Context, logger common.Logger) (*PlayerConfig, error) {
	// Validation function for player config
	validator := func(config PlayerConfig) error {
		if config.Health <= 0 {
			return errors.NewGameError(
				errors.AssetInvalid,
				fmt.Sprintf("invalid player health: %d (must be > 0)", config.Health),
			)
		}
		if config.Size <= 0 {
			return errors.NewGameError(
				errors.AssetInvalid,
				fmt.Sprintf("invalid player size: %d (must be > 0)", config.Size),
			)
		}
		if config.SpriteName == "" {
			return errors.NewGameError(
				errors.AssetInvalid,
				"player sprite_name is required",
			)
		}
		if config.InvincibilityDuration <= 0 {
			return errors.NewGameError(
				errors.AssetInvalid,
				fmt.Sprintf("invalid invincibility_duration: %f (must be > 0)", config.InvincibilityDuration),
			)
		}
		return nil
	}

	// Load and validate using generic loader
	config, err := common.LoadAndValidateJSON(ctx, "entities/player.json", validator)
	if err != nil {
		return nil, err
	}

	if logger != nil {
		logger.Info("Player config loaded from JSON", "health", config.Health, "size", config.Size)
	}

	return &config, nil
}

// LoadEnemyConfigs loads all enemy type configurations from embedded assets
// Returns an error if the file is missing or invalid (no fallback)
func LoadEnemyConfigs(ctx context.Context, logger common.Logger) (*EnemyConfigs, error) {
	// Validation function for enemy configs
	validator := func(configs EnemyConfigs) error {
		if len(configs.EnemyTypes) == 0 {
			return errors.NewGameError(
				errors.AssetInvalid,
				"enemies.json must contain at least one enemy type",
			)
		}

		// Validate each enemy type
		for i := range configs.EnemyTypes {
			if validateErr := validateEnemyType(&configs.EnemyTypes[i], i); validateErr != nil {
				return validateErr
			}
		}
		return nil
	}

	// Load and validate using generic loader
	configs, err := common.LoadAndValidateJSON(ctx, "entities/enemies.json", validator)
	if err != nil {
		return nil, err
	}

	if logger != nil {
		logger.Info("Enemy configs loaded from JSON", "count", len(configs.EnemyTypes))
	}

	return &configs, nil
}

// validateEnemyType validates a single enemy type configuration
func validateEnemyType(enemyType *EnemyTypeConfig, index int) error {
	if enemyType.Type == "" {
		return errors.NewGameError(
			errors.AssetInvalid,
			fmt.Sprintf("enemy type at index %d has empty type field", index),
		)
	}
	if enemyType.Health <= 0 {
		return errors.NewGameError(
			errors.AssetInvalid,
			fmt.Sprintf("enemy type '%s' has invalid health: %d (must be > 0)", enemyType.Type, enemyType.Health),
		)
	}
	if enemyType.Speed <= 0 {
		return errors.NewGameError(
			errors.AssetInvalid,
			fmt.Sprintf("enemy type '%s' has invalid speed: %f (must be > 0)", enemyType.Type, enemyType.Speed),
		)
	}
	if enemyType.Size <= 0 {
		return errors.NewGameError(
			errors.AssetInvalid,
			fmt.Sprintf("enemy type '%s' has invalid size: %d (must be > 0)", enemyType.Type, enemyType.Size),
		)
	}
	if enemyType.SpriteName == "" {
		return errors.NewGameError(
			errors.AssetInvalid,
			fmt.Sprintf("enemy type '%s' has empty sprite_name", enemyType.Type),
		)
	}
	if enemyType.MovementType == "" {
		return errors.NewGameError(
			errors.AssetInvalid,
			fmt.Sprintf("enemy type '%s' has empty movement_type", enemyType.Type),
		)
	}
	return nil
}
