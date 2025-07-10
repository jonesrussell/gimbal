package common_test

import (
	"testing"

	"github.com/jonesrussell/gimbal/internal/common"
)

func TestConfigValidator_ValidConfig(t *testing.T) {
	validator := common.NewConfigValidator()
	config := common.NewConfig(
		common.WithScreenSize(640, 480),
		common.WithPlayerSize(32, 32),
		common.WithNumStars(100),
		common.WithDebug(true),
	)

	result := validator.Validate(config)
	if !result.IsValid {
		t.Errorf("Valid config should pass validation: %s", result.Error())
	}
}

func TestConfigValidator_NilConfig(t *testing.T) {
	validator := common.NewConfigValidator()
	result := validator.Validate(nil)

	if result.IsValid {
		t.Error("Nil config should fail validation")
	}

	if len(result.Errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(result.Errors))
	}

	if result.Errors[0].Field != "config" {
		t.Errorf("Expected error field 'config', got '%s'", result.Errors[0].Field)
	}
}

func TestConfigValidator_InvalidScreenSize(t *testing.T) {
	validator := common.NewConfigValidator()

	// Test negative width
	config := common.NewConfig(common.WithScreenSize(-100, 480))
	result := validator.Validate(config)
	if result.IsValid {
		t.Error("Negative screen width should fail validation")
	}

	// Test negative height
	config = common.NewConfig(common.WithScreenSize(640, -100))
	result = validator.Validate(config)
	if result.IsValid {
		t.Error("Negative screen height should fail validation")
	}

	// Test too small width
	config = common.NewConfig(common.WithScreenSize(100, 480))
	result = validator.Validate(config)
	if result.IsValid {
		t.Error("Too small screen width should fail validation")
	}

	// Test too small height
	config = common.NewConfig(common.WithScreenSize(640, 100))
	result = validator.Validate(config)
	if result.IsValid {
		t.Error("Too small screen height should fail validation")
	}

	// Test too large width
	config = common.NewConfig(common.WithScreenSize(3000, 480))
	result = validator.Validate(config)
	if result.IsValid {
		t.Error("Too large screen width should fail validation")
	}

	// Test too large height
	config = common.NewConfig(common.WithScreenSize(640, 3000))
	result = validator.Validate(config)
	if result.IsValid {
		t.Error("Too large screen height should fail validation")
	}
}

func TestConfigValidator_InvalidPlayerSize(t *testing.T) {
	validator := common.NewConfigValidator()

	// Test negative width
	config := common.NewConfig(
		common.WithScreenSize(640, 480),
		common.WithPlayerSize(-10, 32),
	)
	result := validator.Validate(config)
	if result.IsValid {
		t.Error("Negative player width should fail validation")
	}

	// Test negative height
	config = common.NewConfig(
		common.WithScreenSize(640, 480),
		common.WithPlayerSize(32, -10),
	)
	result = validator.Validate(config)
	if result.IsValid {
		t.Error("Negative player height should fail validation")
	}

	// Test player too large for screen
	config = common.NewConfig(
		common.WithScreenSize(640, 480),
		common.WithPlayerSize(200, 200), // More than 1/4 of screen
	)
	result = validator.Validate(config)
	if result.IsValid {
		t.Error("Player too large for screen should fail validation")
	}
}

func TestConfigValidator_InvalidRadius(t *testing.T) {
	validator := common.NewConfigValidator()

	// Test negative radius
	config := common.NewConfig(
		common.WithScreenSize(640, 480),
		common.WithPlayerSize(32, 32),
	)
	config.Radius = -10
	result := validator.Validate(config)
	if result.IsValid {
		t.Error("Negative radius should fail validation")
	}

	// Test radius too large for screen
	config = common.NewConfig(
		common.WithScreenSize(640, 480),
		common.WithPlayerSize(32, 32),
	)
	config.Radius = 500 // Too large for 640x480 screen
	result = validator.Validate(config)
	if result.IsValid {
		t.Error("Radius too large for screen should fail validation")
	}

	// Test radius too small for player
	config = common.NewConfig(
		common.WithScreenSize(640, 480),
		common.WithPlayerSize(32, 32),
	)
	config.Radius = 5 // Too small for 32x32 player
	result = validator.Validate(config)
	if result.IsValid {
		t.Error("Radius too small for player should fail validation")
	}
}

func TestConfigValidator_InvalidStarConfig(t *testing.T) {
	validator := common.NewConfigValidator()

	// Test negative number of stars
	config := common.NewConfig(
		common.WithScreenSize(640, 480),
		common.WithNumStars(-10),
	)
	result := validator.Validate(config)
	if result.IsValid {
		t.Error("Negative number of stars should fail validation")
	}

	// Test too many stars
	config = common.NewConfig(
		common.WithScreenSize(640, 480),
		common.WithNumStars(2000),
	)
	result = validator.Validate(config)
	if result.IsValid {
		t.Error("Too many stars should fail validation")
	}

	// Test negative star size
	config = common.NewConfig(
		common.WithScreenSize(640, 480),
		common.WithStarSettings(-5, 2.0),
	)
	result = validator.Validate(config)
	if result.IsValid {
		t.Error("Negative star size should fail validation")
	}

	// Test too large star size
	config = common.NewConfig(
		common.WithScreenSize(640, 480),
		common.WithStarSettings(30, 2.0),
	)
	result = validator.Validate(config)
	if result.IsValid {
		t.Error("Too large star size should fail validation")
	}

	// Test negative star speed
	config = common.NewConfig(
		common.WithScreenSize(640, 480),
		common.WithStarSettings(5, -2.0),
	)
	result = validator.Validate(config)
	if result.IsValid {
		t.Error("Negative star speed should fail validation")
	}

	// Test too large star speed
	config = common.NewConfig(
		common.WithScreenSize(640, 480),
		common.WithStarSettings(5, 15.0),
	)
	result = validator.Validate(config)
	if result.IsValid {
		t.Error("Too large star speed should fail validation")
	}
}

func TestConfigValidator_InvalidSpeedValues(t *testing.T) {
	validator := common.NewConfigValidator()

	// Test negative speed
	config := common.NewConfig(
		common.WithScreenSize(640, 480),
		common.WithSpeed(-0.1),
	)
	result := validator.Validate(config)
	if result.IsValid {
		t.Error("Negative speed should fail validation")
	}

	// Test zero speed
	config = common.NewConfig(
		common.WithScreenSize(640, 480),
		common.WithSpeed(0),
	)
	result = validator.Validate(config)
	if result.IsValid {
		t.Error("Zero speed should fail validation")
	}

	// Test too large speed
	config = common.NewConfig(
		common.WithScreenSize(640, 480),
		common.WithSpeed(2.0),
	)
	result = validator.Validate(config)
	if result.IsValid {
		t.Error("Too large speed should fail validation")
	}
}

func TestConfigValidator_InvalidAngleStep(t *testing.T) {
	validator := common.NewConfigValidator()

	// Test negative angle step
	config := common.NewConfig(
		common.WithScreenSize(640, 480),
		common.WithAngleStep(-0.1),
	)
	result := validator.Validate(config)
	if result.IsValid {
		t.Error("Negative angle step should fail validation")
	}

	// Test zero angle step
	config = common.NewConfig(
		common.WithScreenSize(640, 480),
		common.WithAngleStep(0),
	)
	result = validator.Validate(config)
	if result.IsValid {
		t.Error("Zero angle step should fail validation")
	}

	// Test too large angle step
	config = common.NewConfig(
		common.WithScreenSize(640, 480),
		common.WithAngleStep(2.0),
	)
	result = validator.Validate(config)
	if result.IsValid {
		t.Error("Too large angle step should fail validation")
	}
}

func TestValidateConfig_ConvenienceFunction(t *testing.T) {
	// Test valid config
	config := common.NewConfig(
		common.WithScreenSize(640, 480),
		common.WithPlayerSize(32, 32),
		common.WithNumStars(100),
	)

	err := common.ValidateConfig(config)
	if err != nil {
		t.Errorf("Valid config should not return error: %v", err)
	}

	// Test invalid config
	config = common.NewConfig(
		common.WithScreenSize(-100, 480), // Invalid
		common.WithPlayerSize(32, 32),
	)

	err = common.ValidateConfig(config)
	if err == nil {
		t.Error("Invalid config should return error")
	}

	// Check error message format
	if err.Error() == "" {
		t.Error("Error message should not be empty")
	}
}

func TestValidationResult_ErrorFormatting(t *testing.T) {
	result := common.NewValidationResult()

	// Add some errors
	result.AddError("field1", "error message 1")
	result.AddError("field2", "error message 2")

	// Test error message format
	errorMsg := result.Error()
	if errorMsg == "" {
		t.Error("Error message should not be empty when there are errors")
	}

	// Check that both errors are included
	if len(result.Errors) != 2 {
		t.Errorf("Expected 2 errors, got %d", len(result.Errors))
	}

	// Test valid result
	validResult := common.NewValidationResult()
	validErrorMsg := validResult.Error()
	if validErrorMsg != "" {
		t.Error("Valid result should return empty error message")
	}
}
