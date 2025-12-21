package errors //nolint:testpackage // Testing from same package to access unexported functions

import (
	"context"
	"errors"
	"testing"
)

func TestNewGameError(t *testing.T) {
	err := NewGameError(AssetNotFound, "test message")
	if err.Code != AssetNotFound {
		t.Errorf("Expected code %v, got %v", AssetNotFound, err.Code)
	}
	if err.Message != "test message" {
		t.Errorf("Expected message %q, got %q", "test message", err.Message)
	}
	if err.Cause != nil {
		t.Errorf("Expected nil cause, got %v", err.Cause)
	}
	if err.Context == nil {
		t.Error("Expected context map to be initialized")
	}
	if err.Timestamp.IsZero() {
		t.Error("Expected timestamp to be set")
	}
}

func TestNewGameErrorWithCause(t *testing.T) {
	originalErr := errors.New("original error")
	err := NewGameErrorWithCause(SystemInitFailed, "wrapper message", originalErr)
	if err.Code != SystemInitFailed {
		t.Errorf("Expected code %v, got %v", SystemInitFailed, err.Code)
	}
	if err.Message != "wrapper message" {
		t.Errorf("Expected message %q, got %q", "wrapper message", err.Message)
	}
	if err.Cause != originalErr {
		t.Errorf("Expected cause %v, got %v", originalErr, err.Cause)
	}
	if err.Context == nil {
		t.Error("Expected context map to be initialized")
	}
}

func TestGameError_Error(t *testing.T) {
	tests := []struct {
		name    string
		gameErr *GameError
		want    string
	}{
		{
			name:    "error without cause",
			gameErr: NewGameError(AssetNotFound, "test message"),
			want:    "[ASSET_NOT_FOUND] test message",
		},
		{
			name: "error with cause",
			gameErr: NewGameErrorWithCause(
				SystemInitFailed,
				"wrapper message",
				errors.New("underlying error"),
			),
			want: "[SYSTEM_INIT_FAILED] wrapper message: underlying error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.gameErr.Error()
			if got != tt.want {
				t.Errorf("Error() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestGameError_Unwrap(t *testing.T) {
	originalErr := errors.New("original error")
	err := NewGameErrorWithCause(AssetLoadFailed, "wrapper", originalErr)
	unwrapped := err.Unwrap()
	if unwrapped != originalErr {
		t.Errorf("Unwrap() = %v, want %v", unwrapped, originalErr)
	}

	errNoCause := NewGameError(AssetNotFound, "no cause")
	if errNoCause.Unwrap() != nil {
		t.Errorf("Expected nil unwrap for error without cause, got %v", errNoCause.Unwrap())
	}
}

func TestGameError_WithContext(t *testing.T) {
	err := NewGameError(AssetNotFound, "test")
	err = err.WithContext("file", "player.png")
	err = err.WithContext("path", "/assets/sprites")

	if err.Context["file"] != "player.png" {
		t.Errorf("Expected context file=%q, got %v", "player.png", err.Context["file"])
	}
	if err.Context["path"] != "/assets/sprites" {
		t.Errorf("Expected context path=%q, got %v", "/assets/sprites", err.Context["path"])
	}
}

func TestGameError_WithContextMap(t *testing.T) {
	err := NewGameError(AssetNotFound, "test")
	ctx := map[string]interface{}{
		"key1": "value1",
		"key2": 42,
	}
	err = err.WithContextMap(ctx)

	if err.Context["key1"] != "value1" {
		t.Errorf("Expected context key1=%q, got %v", "value1", err.Context["key1"])
	}
	if err.Context["key2"] != 42 {
		t.Errorf("Expected context key2=42, got %v", err.Context["key2"])
	}
}

func TestWrap(t *testing.T) {
	originalErr := errors.New("original")
	wrapped := Wrap(originalErr, SystemUpdateFailed, "wrapped message")

	if wrapped.Code != SystemUpdateFailed {
		t.Errorf("Expected code %v, got %v", SystemUpdateFailed, wrapped.Code)
	}
	if wrapped.Cause != originalErr {
		t.Errorf("Expected cause %v, got %v", originalErr, wrapped.Cause)
	}
}

func TestWrapf(t *testing.T) {
	originalErr := errors.New("original")
	wrapped := Wrapf(originalErr, ConfigInvalid, "failed to load %s", "config.json")

	if wrapped.Code != ConfigInvalid {
		t.Errorf("Expected code %v, got %v", ConfigInvalid, wrapped.Code)
	}
	expectedMsg := "failed to load config.json"
	if wrapped.Message != expectedMsg {
		t.Errorf("Expected message %q, got %q", expectedMsg, wrapped.Message)
	}
	if wrapped.Cause != originalErr {
		t.Errorf("Expected cause %v, got %v", originalErr, wrapped.Cause)
	}
}

func TestFromContext(t *testing.T) {
	tests := []struct {
		name    string
		ctx     context.Context
		wantErr bool
	}{
		{
			name:    "active context",
			ctx:     context.Background(),
			wantErr: false,
		},
		{
			name:    "canceled context",
			ctx:     canceledContext(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := FromContext(tt.ctx, InputTimeout, "operation timeout")
			if err.Code != InputTimeout {
				t.Errorf("Expected code %v, got %v", InputTimeout, err.Code)
			}
			if tt.wantErr && err.Cause == nil {
				t.Error("Expected error to have cause from canceled context")
			}
			if !tt.wantErr && err.Cause != nil {
				t.Errorf("Expected nil cause, got %v", err.Cause)
			}
		})
	}
}

func TestErrorBuilder(t *testing.T) {
	builder := NewErrorBuilder(EntityNotFound, "entity not found")
	builder = builder.WithCause(errors.New("underlying error"))
	builder = builder.WithContext("entity_id", "123")
	builder = builder.WithContext("entity_type", "player")

	err := builder.Build()

	if err.Code != EntityNotFound {
		t.Errorf("Expected code %v, got %v", EntityNotFound, err.Code)
	}
	if err.Message != "entity not found" {
		t.Errorf("Expected message %q, got %q", "entity not found", err.Message)
	}
	if err.Cause == nil {
		t.Error("Expected cause to be set")
	}
	if err.Context["entity_id"] != "123" {
		t.Errorf("Expected context entity_id=%q, got %v", "123", err.Context["entity_id"])
	}
	if err.Context["entity_type"] != "player" {
		t.Errorf("Expected context entity_type=%q, got %v", "player", err.Context["entity_type"])
	}
	if err.Timestamp.IsZero() {
		t.Error("Expected timestamp to be set")
	}
}

func TestErrorBuilder_WithContextMap(t *testing.T) {
	builder := NewErrorBuilder(ValidationFailed, "validation failed")
	ctx := map[string]interface{}{
		"field1": "value1",
		"field2": 99,
	}
	builder = builder.WithContextMap(ctx)

	err := builder.Build()

	if err.Context["field1"] != "value1" {
		t.Errorf("Expected context field1=%q, got %v", "value1", err.Context["field1"])
	}
	if err.Context["field2"] != 99 {
		t.Errorf("Expected context field2=99, got %v", err.Context["field2"])
	}
}

func TestIsGameError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "GameError",
			err:  NewGameError(AssetNotFound, "test"),
			want: true,
		},
		{
			name: "standard error",
			err:  errors.New("standard error"),
			want: false,
		},
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
		{
			name: "wrapped GameError",
			err:  NewGameErrorWithCause(AssetNotFound, "wrapped", errors.New("cause")),
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsGameError(tt.err)
			if got != tt.want {
				t.Errorf("IsGameError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetGameError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantErr  bool
		wantCode ErrorCode
	}{
		{
			name:     "GameError",
			err:      NewGameError(AssetNotFound, "test"),
			wantErr:  true,
			wantCode: AssetNotFound,
		},
		{
			name:    "standard error",
			err:     errors.New("standard error"),
			wantErr: false,
		},
		{
			name:    "nil error",
			err:     nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gameErr, ok := GetGameError(tt.err)
			if ok != tt.wantErr {
				t.Errorf("GetGameError() ok = %v, want %v", ok, tt.wantErr)
			}
			if tt.wantErr && gameErr.Code != tt.wantCode {
				t.Errorf("GetGameError() code = %v, want %v", gameErr.Code, tt.wantCode)
			}
		})
	}
}

func TestHasErrorCode(t *testing.T) {
	tests := []struct {
		name string
		err  error
		code ErrorCode
		want bool
	}{
		{
			name: "matching code",
			err:  NewGameError(AssetNotFound, "test"),
			code: AssetNotFound,
			want: true,
		},
		{
			name: "different code",
			err:  NewGameError(AssetNotFound, "test"),
			code: AssetLoadFailed,
			want: false,
		},
		{
			name: "standard error",
			err:  errors.New("standard error"),
			code: AssetNotFound,
			want: false,
		},
		{
			name: "nil error",
			err:  nil,
			code: AssetNotFound,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HasErrorCode(tt.err, tt.code)
			if got != tt.want {
				t.Errorf("HasErrorCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrorCodes(t *testing.T) {
	// Test that all error codes are defined
	codes := []ErrorCode{
		AssetNotFound, AssetLoadFailed, AssetInvalid, AssetCorrupted, AssetUnsupported,
		EntityNotFound, EntityInvalid, EntityDestroyed, EntityExists,
		ComponentMissing, ComponentInvalid,
		SystemInitFailed, SystemUpdateFailed, SystemCleanupFailed,
		ConfigInvalid, ConfigMissing, ConfigValidation,
		InputInvalid, InputUnsupported, InputTimeout,
		ResourceNotFound, ResourceLoadFailed, ResourceExhausted, ResourceLocked,
		StateInvalid, StateTransition, StateCorrupted,
		RenderFailed, RenderUnsupported, RenderTimeout,
		ValidationFailed, ValidationTimeout,
		SceneNotFound, SceneTransition, SceneLoadFailed,
	}

	for _, code := range codes {
		if code == "" {
			t.Errorf("Error code is empty")
		}
	}
}

func TestErrorBuilder_FluentInterface(t *testing.T) {
	// Test that the fluent interface works correctly
	err := NewErrorBuilder(EntityNotFound, "entity not found").
		WithCause(errors.New("not found")).
		WithContext("id", "123").
		WithContextMap(map[string]interface{}{"type": "player"}).
		Build()

	if err.Code != EntityNotFound {
		t.Errorf("Expected code %v, got %v", EntityNotFound, err.Code)
	}
	if err.Cause == nil {
		t.Error("Expected cause to be set")
	}
	if err.Context["id"] != "123" {
		t.Errorf("Expected context id=%q, got %v", "123", err.Context["id"])
	}
	if err.Context["type"] != "player" {
		t.Errorf("Expected context type=%q, got %v", "player", err.Context["type"])
	}
}

func TestGameError_ContextInitialization(t *testing.T) {
	// Test that context is initialized even if WithContext is never called
	err := NewGameError(AssetNotFound, "test")
	if err.Context == nil {
		t.Error("Expected context to be initialized")
	}

	// Test that WithContext initializes context if nil (shouldn't happen in practice)
	err2 := &GameError{
		Code:    AssetNotFound,
		Message: "test",
		Context: nil,
	}
	err2 = err2.WithContext("key", "value")
	if err2.Context == nil {
		t.Error("Expected context to be initialized by WithContext")
	}
}

func canceledContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	return ctx
}
