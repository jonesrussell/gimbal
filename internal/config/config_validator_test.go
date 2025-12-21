package config //nolint:testpackage // Testing from same package to access unexported functions

import (
	"testing"
)

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func TestNewValidator(t *testing.T) {
	v := NewValidator()
	if v == nil {
		t.Error("NewValidator() returned nil")
	}
}

func TestNewValidationResult(t *testing.T) {
	result := NewValidationResult()
	// NewValidationResult always returns a non-nil pointer, so no nil check needed
	if !result.IsValid {
		t.Error("Expected new validation result to be valid")
	}
	if len(result.Errors) != 0 {
		t.Errorf("Expected empty errors slice, got %d errors", len(result.Errors))
	}
}

func TestValidationResult_AddError(t *testing.T) {
	result := NewValidationResult()
	result.AddError("field1", "error message 1")

	if result.IsValid {
		t.Error("Expected validation result to be invalid after adding error")
	}
	if len(result.Errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(result.Errors))
	}
	if result.Errors[0].Field != "field1" {
		t.Errorf("Expected error field 'field1', got %q", result.Errors[0].Field)
	}
	if result.Errors[0].Message != "error message 1" {
		t.Errorf("Expected error message 'error message 1', got %q", result.Errors[0].Message)
	}

	result.AddError("field2", "error message 2")
	if len(result.Errors) != 2 {
		t.Errorf("Expected 2 errors, got %d", len(result.Errors))
	}
}

func TestValidationResult_Error(t *testing.T) {
	result := NewValidationResult()
	if result.Error() != "" {
		t.Errorf("Expected empty error message for valid result, got %q", result.Error())
	}

	result.AddError("field1", "error 1")
	result.AddError("field2", "error 2")

	errorMsg := result.Error()
	if errorMsg == "" {
		t.Error("Expected non-empty error message for invalid result")
	}
	if len(errorMsg) < 10 {
		t.Errorf("Expected detailed error message, got %q", errorMsg)
	}
}

func TestValidator_Validate_NilConfig(t *testing.T) {
	validator := NewValidator()
	result := validator.Validate(nil)

	if result.IsValid {
		t.Error("Expected validation to fail for nil config")
	}
	if len(result.Errors) == 0 {
		t.Error("Expected at least one error for nil config")
	}
}

func TestValidator_Validate_ValidConfig(t *testing.T) {
	validator := NewValidator()
	config := DefaultConfig()

	result := validator.Validate(config)

	if !result.IsValid {
		t.Errorf("Expected valid config, but got errors: %s", result.Error())
	}
}

//nolint:revive // High cognitive complexity is acceptable for comprehensive table-driven tests
func TestValidator_ValidateScreenSize(t *testing.T) {
	validator := NewValidator()
	tests := []struct {
		name      string
		width     int
		height    int
		wantValid bool
		wantErrs  int
	}{
		{"valid size", 1280, 720, true, 0},
		{"zero width", 0, 720, false, 1},
		{"zero height", 1280, 0, false, 1},
		{"negative width", -100, 720, false, 1},
		{"negative height", 1280, -100, false, 1},
		{"below min width", 319, 720, false, 1},
		{"below min height", 1280, 239, false, 1},
		{"at min width", 320, 720, true, 0},   // Valid after adjusting related configs
		{"at min height", 1280, 240, true, 0}, // Valid after adjusting related configs
		{"above max width", 1921, 720, false, 1},
		{"above max height", 1280, 1081, false, 1},
		{"at max width", 1920, 720, true, 0},
		{"at max height", 1280, 1080, true, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := DefaultConfig()
			config.ScreenSize = Size{Width: tt.width, Height: tt.height}
			// Adjust radius and star spawn radius to avoid unrelated validation errors
			if tt.name == "at min width" || tt.name == "at min height" {
				// For minimum sizes, adjust related configs to be valid
				config.Radius = float64(minInt(tt.width, tt.height)) / 2 * 0.8
				config.StarSpawnRadiusMax = config.Radius * 0.5
			}

			result := validator.Validate(config)

			if result.IsValid != tt.wantValid {
				t.Errorf("Validate() IsValid = %v, want %v. Errors: %s", result.IsValid, tt.wantValid, result.Error())
			}
			// Count screen size related errors
			screenSizeErrs := 0
			for _, err := range result.Errors {
				if err.Field == "screen_size.width" || err.Field == "screen_size.height" {
					screenSizeErrs++
				}
			}
			if screenSizeErrs < tt.wantErrs {
				t.Errorf("Expected at least %d screen size errors, got %d", tt.wantErrs, screenSizeErrs)
			}
		})
	}
}

func TestValidator_ValidatePlayerSize(t *testing.T) {
	validator := NewValidator()
	tests := []struct {
		name      string
		playerW   int
		playerH   int
		screenW   int
		screenH   int
		wantValid bool
	}{
		{"valid player size", 48, 48, 1280, 720, true},
		{"zero player width", 0, 48, 1280, 720, false},
		{"zero player height", 48, 0, 1280, 720, false},
		{"negative player width", -10, 48, 1280, 720, false},
		{"negative player height", 48, -10, 1280, 720, false},
		{"player width exceeds 1/4 screen", 321, 48, 1280, 720, false},
		{"player height exceeds 1/4 screen", 48, 181, 1280, 720, false},
		{"player at max width", 320, 48, 1280, 720, true},
		{"player at max height", 48, 180, 1280, 720, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := DefaultConfig()
			config.ScreenSize = Size{Width: tt.screenW, Height: tt.screenH}
			config.PlayerSize = Size{Width: tt.playerW, Height: tt.playerH}

			result := validator.Validate(config)

			if result.IsValid != tt.wantValid {
				t.Errorf("Validate() IsValid = %v, want %v. Errors: %s", result.IsValid, tt.wantValid, result.Error())
			}
		})
	}
}

func TestValidator_ValidateRadius(t *testing.T) {
	validator := NewValidator()
	tests := []struct {
		name      string
		radius    float64
		screenW   int
		screenH   int
		playerW   int
		wantValid bool
	}{
		{"valid radius", 300, 1280, 720, 48, true},
		{"zero radius", 0, 1280, 720, 48, false},
		{"negative radius", -100, 1280, 720, 48, false},
		{"radius too large", 500, 1280, 720, 48, false},
		{"radius too small", 20, 1280, 720, 48, false},
		{"radius at max", 360, 1280, 720, 48, true},
		{"radius at min", 24, 1280, 720, 48, false}, // May fail star spawn radius validation since default max is 80
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := DefaultConfig()
			config.ScreenSize = Size{Width: tt.screenW, Height: tt.screenH}
			config.PlayerSize = Size{Width: tt.playerW, Height: tt.playerW}
			config.Radius = tt.radius
			// Adjust star spawn radius to avoid unrelated validation errors
			if tt.name == "radius at min" {
				config.StarSpawnRadiusMax = tt.radius * 0.5 // Make sure star spawn radius fits within radius
			}

			result := validator.Validate(config)

			if result.IsValid != tt.wantValid {
				t.Errorf("Validate() IsValid = %v, want %v. Errors: %s", result.IsValid, tt.wantValid, result.Error())
			}
		})
	}
}

func TestValidator_ValidateStarCount(t *testing.T) {
	validator := NewValidator()
	tests := []struct {
		name      string
		numStars  int
		wantValid bool
	}{
		{"valid star count", 100, true},
		{"zero stars", 0, true},
		{"max stars", 1000, true},
		{"negative stars", -1, false},
		{"too many stars", 1001, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := DefaultConfig()
			config.NumStars = tt.numStars

			result := validator.Validate(config)

			if result.IsValid != tt.wantValid {
				t.Errorf("Validate() IsValid = %v, want %v. Errors: %s", result.IsValid, tt.wantValid, result.Error())
			}
		})
	}
}

func TestValidator_ValidateStarSize(t *testing.T) {
	validator := NewValidator()
	tests := []struct {
		name      string
		starSize  float64
		wantValid bool
	}{
		{"valid star size", 5.0, true},
		{"zero star size", 0, false},
		{"negative star size", -1, false},
		{"max star size", 20, true},
		{"too large star size", 21, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := DefaultConfig()
			config.StarSize = tt.starSize

			result := validator.Validate(config)

			if result.IsValid != tt.wantValid {
				t.Errorf("Validate() IsValid = %v, want %v. Errors: %s", result.IsValid, tt.wantValid, result.Error())
			}
		})
	}
}

func TestValidator_ValidateStarSpeed(t *testing.T) {
	validator := NewValidator()
	tests := []struct {
		name      string
		starSpeed float64
		wantValid bool
	}{
		{"valid star speed", 40.0, true},
		{"zero star speed", 0, false},
		{"negative star speed", -1, false},
		{"max star speed", 500, true},
		{"too fast star speed", 501, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := DefaultConfig()
			config.StarSpeed = tt.starSpeed

			result := validator.Validate(config)

			if result.IsValid != tt.wantValid {
				t.Errorf("Validate() IsValid = %v, want %v. Errors: %s", result.IsValid, tt.wantValid, result.Error())
			}
		})
	}
}

func TestValidator_ValidateStarSpawnRadius(t *testing.T) {
	validator := NewValidator()
	tests := []struct {
		name       string
		minRadius  float64
		maxRadius  float64
		gameRadius float64
		wantValid  bool
	}{
		{"valid spawn radius", 30, 80, 300, true},
		{"negative min radius", -1, 80, 300, false},
		{"negative max radius", 30, -1, 300, false},
		{"min equals max", 50, 50, 300, false},
		{"min greater than max", 80, 30, 300, false},
		{"max exceeds game radius", 30, 400, 300, false},
		{"max at game radius", 30, 300, 300, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := DefaultConfig()
			config.Radius = tt.gameRadius
			config.StarSpawnRadiusMin = tt.minRadius
			config.StarSpawnRadiusMax = tt.maxRadius

			result := validator.Validate(config)

			if result.IsValid != tt.wantValid {
				t.Errorf("Validate() IsValid = %v, want %v. Errors: %s", result.IsValid, tt.wantValid, result.Error())
			}
		})
	}
}

func TestValidator_ValidateStarScale(t *testing.T) {
	validator := NewValidator()
	tests := []struct {
		name      string
		minScale  float64
		maxScale  float64
		wantValid bool
	}{
		{"valid star scale", 0.3, 1.0, true},
		{"zero min scale", 0, 1.0, false},
		{"negative min scale", -0.1, 1.0, false},
		{"zero max scale", 0.3, 0, false},
		{"negative max scale", 0.3, -0.1, false},
		{"min equals max", 0.5, 0.5, false},
		{"min greater than max", 1.0, 0.3, false},
		{"max at limit", 0.3, 2.0, true},
		{"max exceeds limit", 0.3, 2.1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := DefaultConfig()
			config.StarMinScale = tt.minScale
			config.StarMaxScale = tt.maxScale

			result := validator.Validate(config)

			if result.IsValid != tt.wantValid {
				t.Errorf("Validate() IsValid = %v, want %v. Errors: %s", result.IsValid, tt.wantValid, result.Error())
			}
		})
	}
}

func TestValidator_ValidateSpeed(t *testing.T) {
	validator := NewValidator()
	tests := []struct {
		name      string
		speed     float64
		wantValid bool
	}{
		{"valid speed", 0.04, true},
		{"zero speed", 0, false},
		{"negative speed", -0.1, false},
		{"max speed", 1.0, true},
		{"too fast speed", 1.1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := DefaultConfig()
			config.Speed = tt.speed

			result := validator.Validate(config)

			if result.IsValid != tt.wantValid {
				t.Errorf("Validate() IsValid = %v, want %v. Errors: %s", result.IsValid, tt.wantValid, result.Error())
			}
		})
	}
}

func TestValidator_ValidateAngleStep(t *testing.T) {
	validator := NewValidator()
	tests := []struct {
		name      string
		angleStep float64
		wantValid bool
	}{
		{"valid angle step", 0.05, true},
		{"zero angle step", 0, false},
		{"negative angle step", -0.1, false},
		{"max angle step", 1.0, true},
		{"too large angle step", 1.1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := DefaultConfig()
			config.AngleStep = tt.angleStep

			result := validator.Validate(config)

			if result.IsValid != tt.wantValid {
				t.Errorf("Validate() IsValid = %v, want %v. Errors: %s", result.IsValid, tt.wantValid, result.Error())
			}
		})
	}
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  *GameConfig
		wantErr bool
	}{
		{"valid config", DefaultConfig(), false},
		{"nil config", nil, true},
		{"invalid screen size", func() *GameConfig {
			c := DefaultConfig()
			c.ScreenSize.Width = -1
			return c
		}(), true},
		{"invalid player size", func() *GameConfig {
			c := DefaultConfig()
			c.PlayerSize.Width = 0
			return c
		}(), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateConfig(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidator_Validate_MultipleErrors(t *testing.T) {
	validator := NewValidator()
	config := DefaultConfig()
	// Make multiple fields invalid
	config.ScreenSize.Width = -1
	config.PlayerSize.Height = 0
	config.Radius = -100
	config.NumStars = -5
	config.Speed = -1

	result := validator.Validate(config)

	if result.IsValid {
		t.Error("Expected validation to fail with multiple errors")
	}
	if len(result.Errors) < 5 {
		t.Errorf("Expected at least 5 errors, got %d. Errors: %s", len(result.Errors), result.Error())
	}
}
