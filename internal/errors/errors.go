package errors

import (
	"context"
	"fmt"
	"time"
)

// GameError represents a game-specific error with code, message, and cause
type GameError struct {
	Code      ErrorCode
	Message   string
	Cause     error
	Timestamp time.Time
	Context   map[string]interface{}
}

// Error implements the error interface
func (e *GameError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap returns the underlying cause error
func (e *GameError) Unwrap() error {
	return e.Cause
}

// newContextMap creates a new context map
func newContextMap() map[string]interface{} {
	return make(map[string]interface{})
}

// WithContext adds context information to the error
func (e *GameError) WithContext(key string, value interface{}) *GameError {
	if e.Context == nil {
		e.Context = newContextMap()
	}
	e.Context[key] = value
	return e
}

// WithContextMap adds multiple context values to the error
func (e *GameError) WithContextMap(ctx map[string]interface{}) *GameError {
	if e.Context == nil {
		e.Context = newContextMap()
	}
	for k, v := range ctx {
		e.Context[k] = v
	}
	return e
}

// NewGameError creates a new game error with the specified code and message
func NewGameError(code ErrorCode, message string) *GameError {
	return &GameError{
		Code:      code,
		Message:   message,
		Timestamp: time.Now(),
		Context:   newContextMap(),
	}
}

// NewGameErrorWithCause creates a new game error with the specified code, message, and underlying cause
func NewGameErrorWithCause(code ErrorCode, message string, cause error) *GameError {
	return &GameError{
		Code:      code,
		Message:   message,
		Cause:     cause,
		Timestamp: time.Now(),
		Context:   newContextMap(),
	}
}

// Wrap wraps an existing error with game error context
func Wrap(err error, code ErrorCode, message string) *GameError {
	return NewGameErrorWithCause(code, message, err)
}

// Wrapf wraps an existing error with formatted message
func Wrapf(err error, code ErrorCode, format string, args ...interface{}) *GameError {
	return NewGameErrorWithCause(code, fmt.Sprintf(format, args...), err)
}

// FromContext creates an error from context cancellation
func FromContext(ctx context.Context, code ErrorCode, message string) *GameError {
	err := NewGameError(code, message)
	if ctx.Err() != nil {
		err.Cause = ctx.Err()
	}
	return err
}

// ErrorCode represents a typed error code for better type safety
type ErrorCode string

const (
	// Asset error codes
	AssetNotFound    ErrorCode = "ASSET_NOT_FOUND"
	AssetLoadFailed  ErrorCode = "ASSET_LOAD_FAILED"
	AssetInvalid     ErrorCode = "ASSET_INVALID"
	AssetCorrupted   ErrorCode = "ASSET_CORRUPTED"
	AssetUnsupported ErrorCode = "ASSET_UNSUPPORTED"

	// Entity error codes
	EntityNotFound  ErrorCode = "ENTITY_NOT_FOUND"
	EntityInvalid   ErrorCode = "ENTITY_INVALID"
	EntityDestroyed ErrorCode = "ENTITY_DESTROYED"
	EntityExists    ErrorCode = "ENTITY_EXISTS"

	// Component error codes
	ComponentMissing ErrorCode = "COMPONENT_MISSING"
	ComponentInvalid ErrorCode = "COMPONENT_INVALID"

	// System error codes
	SystemInitFailed    ErrorCode = "SYSTEM_INIT_FAILED"
	SystemUpdateFailed  ErrorCode = "SYSTEM_UPDATE_FAILED"
	SystemCleanupFailed ErrorCode = "SYSTEM_CLEANUP_FAILED"

	// Configuration error codes
	ConfigInvalid    ErrorCode = "CONFIG_INVALID"
	ConfigMissing    ErrorCode = "CONFIG_MISSING"
	ConfigValidation ErrorCode = "CONFIG_VALIDATION"

	// Input error codes
	InputInvalid     ErrorCode = "INPUT_INVALID"
	InputUnsupported ErrorCode = "INPUT_UNSUPPORTED"
	InputTimeout     ErrorCode = "INPUT_TIMEOUT"

	// Resource error codes
	ResourceNotFound   ErrorCode = "RESOURCE_NOT_FOUND"
	ResourceLoadFailed ErrorCode = "RESOURCE_LOAD_FAILED"
	ResourceExhausted  ErrorCode = "RESOURCE_EXHAUSTED"
	ResourceLocked     ErrorCode = "RESOURCE_LOCKED"

	// Game state error codes
	StateInvalid    ErrorCode = "STATE_INVALID"
	StateTransition ErrorCode = "STATE_TRANSITION"
	StateCorrupted  ErrorCode = "STATE_CORRUPTED"

	// Rendering error codes
	RenderFailed      ErrorCode = "RENDER_FAILED"
	RenderUnsupported ErrorCode = "RENDER_UNSUPPORTED"
	RenderTimeout     ErrorCode = "RENDER_TIMEOUT"

	// Validation error codes
	ValidationFailed  ErrorCode = "VALIDATION_FAILED"
	ValidationTimeout ErrorCode = "VALIDATION_TIMEOUT"

	// Scene error codes
	SceneNotFound   ErrorCode = "SCENE_NOT_FOUND"
	SceneTransition ErrorCode = "SCENE_TRANSITION"
	SceneLoadFailed ErrorCode = "SCENE_LOAD_FAILED"
)

// ErrorBuilder provides a fluent interface for building errors
type ErrorBuilder struct {
	code    ErrorCode
	message string
	cause   error
	context map[string]interface{}
}

// NewErrorBuilder creates a new error builder
func NewErrorBuilder(code ErrorCode, message string) *ErrorBuilder {
	return &ErrorBuilder{
		code:    code,
		message: message,
		context: newContextMap(),
	}
}

// WithCause sets the underlying cause error
func (b *ErrorBuilder) WithCause(cause error) *ErrorBuilder {
	b.cause = cause
	return b
}

// WithContext adds context information
func (b *ErrorBuilder) WithContext(key string, value interface{}) *ErrorBuilder {
	if b.context == nil {
		b.context = newContextMap()
	}
	b.context[key] = value
	return b
}

// WithContextMap adds multiple context values
func (b *ErrorBuilder) WithContextMap(ctx map[string]interface{}) *ErrorBuilder {
	if b.context == nil {
		b.context = newContextMap()
	}
	for k, v := range ctx {
		b.context[k] = v
	}
	return b
}

// Build creates the final error
func (b *ErrorBuilder) Build() *GameError {
	err := &GameError{
		Code:      b.code,
		Message:   b.message,
		Cause:     b.cause,
		Timestamp: time.Now(),
		Context:   b.context,
	}
	return err
}

// IsGameError checks if an error is a GameError
func IsGameError(err error) bool {
	_, ok := err.(*GameError)
	return ok
}

// GetGameError extracts GameError from an error chain
func GetGameError(err error) (*GameError, bool) {
	if err != nil {
		if gameErr, ok := err.(*GameError); ok {
			return gameErr, true
		}
	}
	return nil, false
}

// HasErrorCode checks if an error has a specific error code
func HasErrorCode(err error, code ErrorCode) bool {
	if gameErr, ok := GetGameError(err); ok {
		return gameErr.Code == code
	}
	return false
}
