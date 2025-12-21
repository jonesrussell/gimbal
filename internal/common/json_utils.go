package common

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jonesrussell/gimbal/assets"
	"github.com/jonesrussell/gimbal/internal/errors"
)

// ValidatorFunc is a function that validates loaded configuration
type ValidatorFunc[T any] func(T) error

// LoadAndValidateJSON loads JSON from embedded assets and validates it
func LoadAndValidateJSON[T any](ctx context.Context, path string, validator ValidatorFunc[T]) (T, error) {
	var zero T

	// Check context cancellation
	if err := CheckContextCancellation(ctx); err != nil {
		return zero, err
	}

	// Read file from embedded assets
	data, err := assets.Assets.ReadFile(path)
	if err != nil {
		return zero, errors.NewGameErrorWithCause(
			errors.AssetLoadFailed,
			fmt.Sprintf("failed to read %s", path),
			err,
		)
	}

	// Unmarshal JSON
	var config T
	if unmarshalErr := json.Unmarshal(data, &config); unmarshalErr != nil {
		return zero, errors.NewGameErrorWithCause(
			errors.ConfigInvalid,
			fmt.Sprintf("failed to parse %s", path),
			unmarshalErr,
		)
	}

	// Validate if validator provided
	if validator != nil {
		if err := validator(config); err != nil {
			return zero, errors.NewGameErrorWithCause(
				errors.ConfigInvalid,
				fmt.Sprintf("invalid configuration in %s", path),
				err,
			)
		}
	}

	return config, nil
}
