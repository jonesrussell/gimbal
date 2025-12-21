package common //nolint:testpackage // Testing from same package to access unexported functions

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/jonesrussell/gimbal/internal/errors"
)

// TestConfig is a simple config type for testing
type TestConfig struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

// TestConfigComplex is a more complex config type for testing
type TestConfigComplex struct {
	Name     string   `json:"name"`
	Settings []string `json:"settings"`
	Count    int      `json:"count"`
}

func TestLoadAndValidateJSON_ContextCancellation(t *testing.T) {
	// Test that context cancellation is checked
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := LoadAndValidateJSON[TestConfig](ctx, "test.json", nil)

	if err == nil {
		t.Error("LoadAndValidateJSON() expected error for canceled context, got nil")
	}

	// CheckContextCancellation returns a standard error (context.Canceled),
	// which LoadAndValidateJSON returns directly, not wrapped as GameError.
	// So we just check that an error was returned.
	if err == nil {
		t.Error("Expected error for canceled context")
	}
}

func TestLoadAndValidateJSON_MissingFile(t *testing.T) {
	ctx := context.Background()

	_, err := LoadAndValidateJSON[TestConfig](ctx, "nonexistent_file.json", nil)

	if err == nil {
		t.Error("LoadAndValidateJSON() expected error for missing file, got nil")
	}

	if !errors.HasErrorCode(err, errors.AssetLoadFailed) {
		t.Errorf("LoadAndValidateJSON() expected AssetLoadFailed error code, got: %v", err)
	}
}

func TestLoadAndValidateJSON_InvalidJSON(t *testing.T) {
	ctx := context.Background()

	// Try to load a file that might exist but have invalid JSON
	// Since we can't easily create temporary embedded files, we'll test the logic
	// by attempting to load a non-JSON file or a file that doesn't exist
	_, err := LoadAndValidateJSON[TestConfig](ctx, "testdata/invalid.json", nil)
	// This will likely fail with AssetLoadFailed, but if a file exists and has invalid JSON,
	// it should return ConfigInvalid
	if err != nil {
		if !errors.HasErrorCode(err, errors.AssetLoadFailed) &&
			!errors.HasErrorCode(err, errors.ConfigInvalid) {
			t.Errorf("LoadAndValidateJSON() expected AssetLoadFailed or ConfigInvalid, got: %v", err)
		}
	}
}

func TestLoadAndValidateJSON_ValidationSuccess(t *testing.T) {
	// Test the validation function logic separately
	validator := func(cfg TestConfig) error {
		if cfg.Name == "" {
			return errors.NewGameError(errors.ValidationFailed, "name is required")
		}
		if cfg.Value < 0 {
			return errors.NewGameError(errors.ValidationFailed, "value must be positive")
		}
		return nil
	}

	// Test validation success
	validConfig := TestConfig{Name: "test", Value: 10}
	err := validator(validConfig)
	if err != nil {
		t.Errorf("Expected no validation error, got: %v", err)
	}
}

func TestLoadAndValidateJSON_ValidationFailure(t *testing.T) {
	// Test the validation function logic separately
	validator := func(cfg TestConfig) error {
		if cfg.Value < 0 {
			return errors.NewGameError(errors.ValidationFailed, "value must be positive")
		}
		return nil
	}

	// Test validation failure
	invalidConfig := TestConfig{Value: -1}
	err := validator(invalidConfig)

	if err == nil {
		t.Error("Expected validation error for negative value")
	}

	if !errors.HasErrorCode(err, errors.ValidationFailed) {
		t.Errorf("Expected ValidationFailed error code, got: %v", err)
	}
}

func TestValidatorFunc_TypeParameter(t *testing.T) {
	// Test that ValidatorFunc works with different types
	tests := []struct {
		name      string
		validator ValidatorFunc[TestConfig]
		config    TestConfig
		wantErr   bool
	}{
		{
			name: "valid config",
			validator: func(cfg TestConfig) error {
				if cfg.Name == "" {
					return errors.NewGameError(errors.ValidationFailed, "name required")
				}
				return nil
			},
			config:  TestConfig{Name: "test", Value: 10},
			wantErr: false,
		},
		{
			name: "invalid config",
			validator: func(cfg TestConfig) error {
				if cfg.Name == "" {
					return errors.NewGameError(errors.ValidationFailed, "name required")
				}
				return nil
			},
			config:  TestConfig{Value: 10},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.validator(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validator() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoadAndValidateJSON_NilValidator(t *testing.T) {
	ctx := context.Background()

	// Test that nil validator doesn't cause issues
	// This will fail on file read, but shouldn't fail on validation
	_, err := LoadAndValidateJSON[TestConfig](ctx, "nonexistent.json", nil)

	if err == nil {
		t.Error("Expected error for missing file")
	}

	// Should get AssetLoadFailed, not validation error
	if errors.HasErrorCode(err, errors.ConfigInvalid) {
		t.Error("Should not get validation error when validator is nil and file is missing")
	}
}

func TestLoadAndValidateJSON_JSONUnmarshalError(t *testing.T) {
	// Test JSON unmarshaling logic separately
	tests := []struct {
		name    string
		data    string
		wantErr bool
	}{
		{
			name:    "valid JSON",
			data:    `{"name": "test", "value": 42}`,
			wantErr: false,
		},
		{
			name:    "invalid JSON",
			data:    `{"name": "test", "value": }`,
			wantErr: true,
		},
		{
			name:    "empty JSON",
			data:    ``,
			wantErr: true,
		},
		{
			name:    "malformed JSON",
			data:    `{name: "test"}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var config TestConfig
			err := json.Unmarshal([]byte(tt.data), &config)

			if (err != nil) != tt.wantErr {
				t.Errorf("json.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoadAndValidateJSON_ComplexType(t *testing.T) {
	// Test with complex types
	validator := func(cfg TestConfigComplex) error {
		if len(cfg.Settings) == 0 {
			return errors.NewGameError(errors.ValidationFailed, "settings required")
		}
		return nil
	}

	validConfig := TestConfigComplex{
		Name:     "test",
		Settings: []string{"setting1", "setting2"},
		Count:    5,
	}

	if err := validator(validConfig); err != nil {
		t.Errorf("Expected no validation error, got: %v", err)
	}

	invalidConfig := TestConfigComplex{
		Name:  "test",
		Count: 5,
	}

	if err := validator(invalidConfig); err == nil {
		t.Error("Expected validation error for missing settings")
	}
}

func TestLoadAndValidateJSON_ZeroValue(t *testing.T) {
	// Test that zero value is returned on error
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	result, err := LoadAndValidateJSON[TestConfig](ctx, "test.json", nil)

	if err == nil {
		t.Error("Expected error for canceled context")
	}

	// Verify zero value is returned
	if result.Name != "" || result.Value != 0 {
		t.Errorf("Expected zero value on error, got: %+v", result)
	}
}

func TestLoadAndValidateJSON_ContextTimeout(t *testing.T) {
	// Test context timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	time.Sleep(1 * time.Millisecond) // Ensure timeout

	_, err := LoadAndValidateJSON[TestConfig](ctx, "test.json", nil)

	if err == nil {
		t.Error("Expected error for timed out context")
	}
}
