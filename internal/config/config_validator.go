// Package config provides game configuration management and validation
package config

import "fmt"

// Validator validates game configuration parameters and ensures they meet game requirements
type Validator struct {
	// Add any validation dependencies here
}

// NewValidator creates a new configuration validator
func NewValidator() *Validator {
	return &Validator{}
}

// ValidationError represents a configuration validation error
type ValidationError struct {
	Field   string
	Message string
}

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
func (v *Validator) Validate(config *GameConfig) *ValidationResult {
	result := NewValidationResult()

	if config == nil {
		result.AddError("config", "configuration cannot be nil")
		return result
	}

	// Validate each configuration section
	validators := []func(*GameConfig, *ValidationResult){
		v.validateScreenSize,
		v.validatePlayerSize,
		v.validateRadius,
		v.validateStarConfig,
		v.validateSpeedValues,
		v.validateAngleStep,
	}

	for _, validate := range validators {
		validate(config, result)
	}

	return result
}

// validateScreenSize validates screen dimensions
func (v *Validator) validateScreenSize(config *GameConfig, result *ValidationResult) {
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
func (v *Validator) validatePlayerSize(config *GameConfig, result *ValidationResult) {
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
func (v *Validator) validateRadius(config *GameConfig, result *ValidationResult) {
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

// validateStarConfig validates all star field configuration parameters
func (v *Validator) validateStarConfig(config *GameConfig, result *ValidationResult) {
	validators := []func(*GameConfig, *ValidationResult){
		v.validateStarCount,
		v.validateStarSize,
		v.validateStarSpeed,
		v.validateStarSpawnRadius,
		v.validateStarScale,
	}

	for _, validate := range validators {
		validate(config, result)
	}
}

// validateStarCount ensures star count is within performance limits
func (v *Validator) validateStarCount(config *GameConfig, result *ValidationResult) {
	if config.NumStars < 0 {
		result.AddError("num_stars", "cannot be negative")
	}
	if config.NumStars > 1000 {
		result.AddError("num_stars", "maximum 1000 stars allowed")
	}
}

// validateStarSize ensures star size is within reasonable bounds
func (v *Validator) validateStarSize(config *GameConfig, result *ValidationResult) {
	if config.StarSize <= 0 {
		result.AddError("star_size", "must be positive")
	}
	if config.StarSize > 20 {
		result.AddError("star_size", "maximum star size is 20")
	}
}

// validateStarSpeed ensures star movement speed is within game balance requirements
func (v *Validator) validateStarSpeed(config *GameConfig, result *ValidationResult) {
	if config.StarSpeed <= 0 {
		result.AddError("star_speed", "must be positive")
	}
	if config.StarSpeed > 500 {
		result.AddError("star_speed", "maximum star speed is 500")
	}
}

// validateStarSpawnRadius ensures spawn radius values are within acceptable bounds
func (v *Validator) validateStarSpawnRadius(config *GameConfig, result *ValidationResult) {
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
}

// validateStarScale ensures star scale values are reasonable
func (v *Validator) validateStarScale(config *GameConfig, result *ValidationResult) {
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
func (v *Validator) validateSpeedValues(config *GameConfig, result *ValidationResult) {
	if config.Speed <= 0 {
		result.AddError("speed", "must be positive")
	}
	if config.Speed > 1.0 {
		result.AddError("speed", "maximum speed is 1.0")
	}
}

// validateAngleStep validates angle step configuration
func (v *Validator) validateAngleStep(config *GameConfig, result *ValidationResult) {
	if config.AngleStep <= 0 {
		result.AddError("angle_step", "must be positive")
	}
	if config.AngleStep > 1.0 {
		result.AddError("angle_step", "maximum angle step is 1.0")
	}
}

// ValidateConfig is a convenience function for quick validation
func ValidateConfig(config *GameConfig) error {
	validator := NewValidator()
	result := validator.Validate(config)
	if !result.IsValid {
		return fmt.Errorf("%s", result.Error())
	}
	return nil
}
