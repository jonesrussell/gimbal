package config

import "fmt"

// ConfigValidator validates game configuration
type ConfigValidator struct{}

// NewConfigValidator creates a new configuration validator
func NewConfigValidator() *ConfigValidator {
	return &ConfigValidator{}
}

// ValidationError represents a configuration validation error
type ValidationError struct {
	Field   string
	Message string
}

// Error() method removed - dead code

// ValidationResult contains validation results
type ValidationResult struct {
	IsValid bool
	Errors  []ValidationError
}

// NewValidationResult creates a new validation result
func NewValidationResult() *ValidationResult {
	return &ValidationResult{
		IsValid: true,
		Errors:  make([]ValidationError, 0),
	}
}

// AddError adds a validation error
func (r *ValidationResult) AddError(field, message string) {
	r.IsValid = false
	r.Errors = append(r.Errors, ValidationError{
		Field:   field,
		Message: message,
	})
}

// Error returns the combined error message
func (r *ValidationResult) Error() string {
	if r.IsValid {
		return ""
	}

	msg := "configuration validation failed:\n"
	for _, err := range r.Errors {
		msg += fmt.Sprintf("  - %s: %s\n", err.Field, err.Message)
	}
	return msg
}

// Validate validates the game configuration
func (v *ConfigValidator) Validate(config *GameConfig) *ValidationResult {
	result := NewValidationResult()

	if config == nil {
		result.AddError("config", "configuration cannot be nil")
		return result
	}

	// Validate screen size
	v.validateScreenSize(config, result)

	// Validate player size
	v.validatePlayerSize(config, result)

	// Validate radius
	v.validateRadius(config, result)

	// Validate star configuration
	v.validateStarConfig(config, result)

	// Validate speed values
	v.validateSpeedValues(config, result)

	// Validate angle step
	v.validateAngleStep(config, result)

	return result
}

// validateScreenSize validates screen dimensions
func (v *ConfigValidator) validateScreenSize(config *GameConfig, result *ValidationResult) {
	if config.ScreenSize.Width <= 0 {
		result.AddError("screen_size.width", "must be positive")
	}
	if config.ScreenSize.Height <= 0 {
		result.AddError("screen_size.height", "must be positive")
	}
	if config.ScreenSize.Width < 320 {
		result.AddError("screen_size.width", "minimum width is 320")
	}
	if config.ScreenSize.Height < 240 {
		result.AddError("screen_size.height", "minimum height is 240")
	}
	if config.ScreenSize.Width > 1920 {
		result.AddError("screen_size.width", "maximum width is 1920")
	}
	if config.ScreenSize.Height > 1080 {
		result.AddError("screen_size.height", "maximum height is 1080")
	}
}

// validatePlayerSize validates player dimensions
func (v *ConfigValidator) validatePlayerSize(config *GameConfig, result *ValidationResult) {
	if config.PlayerSize.Width <= 0 {
		result.AddError("player_size.width", "must be positive")
	}
	if config.PlayerSize.Height <= 0 {
		result.AddError("player_size.height", "must be positive")
	}
	if config.PlayerSize.Width > config.ScreenSize.Width/4 {
		result.AddError("player_size.width", "player width cannot exceed 1/4 of screen width")
	}
	if config.PlayerSize.Height > config.ScreenSize.Height/4 {
		result.AddError("player_size.height", "player height cannot exceed 1/4 of screen height")
	}
}

// validateRadius validates the orbital radius
func (v *ConfigValidator) validateRadius(config *GameConfig, result *ValidationResult) {
	if config.Radius <= 0 {
		result.AddError("radius", "must be positive")
		return
	}

	// Radius should fit within screen bounds
	smallerDim := config.ScreenSize.Width
	if config.ScreenSize.Height < config.ScreenSize.Width {
		smallerDim = config.ScreenSize.Height
	}

	maxRadius := float64(smallerDim) / 2
	if config.Radius > maxRadius {
		result.AddError("radius", fmt.Sprintf("radius %.1f exceeds maximum allowed %.1f", config.Radius, maxRadius))
	}

	// Minimum radius should accommodate player size
	minRadius := float64(config.PlayerSize.Width) / 2
	if config.Radius < minRadius {
		result.AddError("radius", fmt.Sprintf("radius %.1f is too small for player size", config.Radius))
	}
}

// validateStarConfig validates star-related configuration
func (v *ConfigValidator) validateStarConfig(config *GameConfig, result *ValidationResult) {
	if config.NumStars < 0 {
		result.AddError("num_stars", "cannot be negative")
	}
	if config.NumStars > 1000 {
		result.AddError("num_stars", "maximum 1000 stars allowed")
	}

	if config.StarSize <= 0 {
		result.AddError("star_size", "must be positive")
	}
	if config.StarSize > 20 {
		result.AddError("star_size", "maximum star size is 20")
	}

	if config.StarSpeed <= 0 {
		result.AddError("star_speed", "must be positive")
	}
	if config.StarSpeed > 500 {
		result.AddError("star_speed", "maximum star speed is 500")
	}

	// Validate star spawn radius
	if config.StarSpawnRadiusMin < 0 {
		result.AddError("star_spawn_radius_min", "cannot be negative")
	}
	if config.StarSpawnRadiusMax < 0 {
		result.AddError("star_spawn_radius_max", "cannot be negative")
	}
	if config.StarSpawnRadiusMin >= config.StarSpawnRadiusMax {
		result.AddError("star_spawn_radius", "min radius must be less than max radius")
	}
	if config.StarSpawnRadiusMax > config.Radius {
		result.AddError("star_spawn_radius_max", "cannot exceed game radius")
	}

	// Validate star scale
	if config.StarMinScale <= 0 {
		result.AddError("star_min_scale", "must be positive")
	}
	if config.StarMaxScale <= 0 {
		result.AddError("star_max_scale", "must be positive")
	}
	if config.StarMinScale >= config.StarMaxScale {
		result.AddError("star_scale", "min scale must be less than max scale")
	}
	if config.StarMaxScale > 2.0 {
		result.AddError("star_max_scale", "maximum scale is 2.0")
	}
}

// validateSpeedValues validates speed-related configuration
func (v *ConfigValidator) validateSpeedValues(config *GameConfig, result *ValidationResult) {
	if config.Speed <= 0 {
		result.AddError("speed", "must be positive")
	}
	if config.Speed > 1.0 {
		result.AddError("speed", "maximum speed is 1.0")
	}
}

// validateAngleStep validates angle step configuration
func (v *ConfigValidator) validateAngleStep(config *GameConfig, result *ValidationResult) {
	if config.AngleStep <= 0 {
		result.AddError("angle_step", "must be positive")
	}
	if config.AngleStep > 1.0 {
		result.AddError("angle_step", "maximum angle step is 1.0")
	}
}

// ValidateConfig is a convenience function for quick validation
func ValidateConfig(config *GameConfig) error {
	validator := NewConfigValidator()
	result := validator.Validate(config)
	if !result.IsValid {
		return fmt.Errorf("%s", result.Error())
	}
	return nil
}
