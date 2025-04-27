package input_test

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/input"
	"github.com/jonesrussell/gimbal/internal/logger"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Parallel()

	mockLogger := logger.NewMock()
	handler := input.New(mockLogger)
	assert.NotNil(t, handler)
}

func TestHandler_IsKeyPressed(t *testing.T) {
	t.Parallel()

	mockLogger := logger.NewMock()
	handler := input.New(mockLogger)

	// Test key press detection
	pressed := handler.IsKeyPressed(ebiten.KeySpace)
	assert.False(t, pressed)
}

func TestHandler_GetMovementInput(t *testing.T) {
	t.Parallel()

	mockLogger := logger.NewMock()
	handler := input.New(mockLogger)

	// Test movement input
	angle := handler.GetMovementInput()
	assert.InDelta(t, float64(common.Angle(0)), float64(angle), 0.0001, "Expected zero movement when no input")
}

func TestHandler_IsPausePressed(t *testing.T) {
	t.Parallel()

	mockLogger := logger.NewMock()
	handler := input.New(mockLogger)

	// Test pause detection
	paused := handler.IsPausePressed()
	assert.False(t, paused)
}

func TestHandler_IsQuitPressed(t *testing.T) {
	t.Parallel()

	mockLogger := logger.NewMock()
	handler := input.New(mockLogger)

	// Test quit detection
	quit := handler.IsQuitPressed()
	assert.False(t, quit)
}

func TestHandler_HandleInput(t *testing.T) {
	t.Parallel()

	mockLogger := logger.NewMock()
	handler := input.New(mockLogger)

	// Test input handling
	handler.HandleInput()
	// No assertions needed as this just logs input state
}
