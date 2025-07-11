package common

import "fmt"

// GameError represents a game-specific error with code, message, and cause
type GameError struct {
	Code    string
	Message string
	Cause   error
}

// Error implements the error interface
func (e *GameError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the underlying cause error
func (e *GameError) Unwrap() error {
	return e.Cause
}

// NewGameError creates a new GameError
func NewGameError(code, message string) *GameError {
	return &GameError{
		Code:    code,
		Message: message,
	}
}

// NewGameErrorWithCause creates a new GameError with a cause
func NewGameErrorWithCause(code, message string, cause error) *GameError {
	return &GameError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

// Predefined error constants
var (
	// Asset errors
	ErrAssetNotFound   = NewGameError("ASSET_NOT_FOUND", "asset not found")
	ErrAssetLoadFailed = NewGameError("ASSET_LOAD_FAILED", "failed to load asset")
	ErrAssetInvalid    = NewGameError("ASSET_INVALID", "asset is invalid or corrupted")

	// Entity errors
	ErrEntityInvalid        = NewGameError("ENTITY_INVALID", "entity is invalid")
	ErrEntityNotFound       = NewGameError("ENTITY_NOT_FOUND", "entity not found")
	ErrEntityCreationFailed = NewGameError("ENTITY_CREATION_FAILED", "failed to create entity")

	// Component errors
	ErrComponentNotFound = NewGameError("COMPONENT_NOT_FOUND", "component not found")
	ErrComponentInvalid  = NewGameError("COMPONENT_INVALID", "component is invalid")

	// System errors
	ErrSystemFailed   = NewGameError("SYSTEM_FAILED", "system execution failed")
	ErrSystemNotFound = NewGameError("SYSTEM_NOT_FOUND", "system not found")

	// Configuration errors
	ErrConfigInvalid = NewGameError("CONFIG_INVALID", "configuration is invalid")
	ErrConfigMissing = NewGameError("CONFIG_MISSING", "required configuration missing")

	// Input errors
	ErrInputInvalid      = NewGameError("INPUT_INVALID", "input is invalid")
	ErrInputNotSupported = NewGameError("INPUT_NOT_SUPPORTED", "input not supported")

	// Resource errors
	ErrResourceNotFound   = NewGameError("RESOURCE_NOT_FOUND", "resource not found")
	ErrResourceLoadFailed = NewGameError("RESOURCE_LOAD_FAILED", "failed to load resource")
	ErrResourceInvalid    = NewGameError("RESOURCE_INVALID", "resource is invalid")

	// Game state errors
	ErrGameStateInvalid          = NewGameError("GAME_STATE_INVALID", "game state is invalid")
	ErrGameStateTransitionFailed = NewGameError("GAME_STATE_TRANSITION_FAILED", "failed to transition game state")

	// Rendering errors
	ErrRenderingFailed = NewGameError("RENDERING_FAILED", "rendering failed")
	ErrSpriteNotFound  = NewGameError("SPRITE_NOT_FOUND", "sprite not found")

	// Validation errors
	ErrValidationFailed = NewGameError("VALIDATION_FAILED", "validation failed")
	ErrValueOutOfRange  = NewGameError("VALUE_OUT_OF_RANGE", "value is out of valid range")
)

// Error codes for easy reference
const (
	// Asset error codes
	ErrorCodeAssetNotFound   = "ASSET_NOT_FOUND"
	ErrorCodeAssetLoadFailed = "ASSET_LOAD_FAILED"
	ErrorCodeAssetInvalid    = "ASSET_INVALID"

	// Entity error codes
	ErrorCodeEntityInvalid        = "ENTITY_INVALID"
	ErrorCodeEntityNotFound       = "ENTITY_NOT_FOUND"
	ErrorCodeEntityCreationFailed = "ENTITY_CREATION_FAILED"

	// Component error codes
	ErrorCodeComponentNotFound = "COMPONENT_NOT_FOUND"
	ErrorCodeComponentInvalid  = "COMPONENT_INVALID"

	// System error codes
	ErrorCodeSystemFailed   = "SYSTEM_FAILED"
	ErrorCodeSystemNotFound = "SYSTEM_NOT_FOUND"

	// Configuration error codes
	ErrorCodeConfigInvalid = "CONFIG_INVALID"
	ErrorCodeConfigMissing = "CONFIG_MISSING"

	// Input error codes
	ErrorCodeInputInvalid      = "INPUT_INVALID"
	ErrorCodeInputNotSupported = "INPUT_NOT_SUPPORTED"

	// Resource error codes
	ErrorCodeResourceNotFound   = "RESOURCE_NOT_FOUND"
	ErrorCodeResourceLoadFailed = "RESOURCE_LOAD_FAILED"
	ErrorCodeResourceInvalid    = "RESOURCE_INVALID"

	// Game state error codes
	ErrorCodeGameStateInvalid          = "GAME_STATE_INVALID"
	ErrorCodeGameStateTransitionFailed = "GAME_STATE_TRANSITION_FAILED"

	// Rendering error codes
	ErrorCodeRenderingFailed = "RENDERING_FAILED"
	ErrorCodeSpriteNotFound  = "SPRITE_NOT_FOUND"

	// Validation error codes
	ErrorCodeValidationFailed = "VALIDATION_FAILED"
	ErrorCodeValueOutOfRange  = "VALUE_OUT_OF_RANGE"

	// Scene error codes
	ErrorCodeSceneNotFound = "SCENE_NOT_FOUND"
)
