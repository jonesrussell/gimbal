package input_test

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/input"
	"github.com/jonesrussell/gimbal/internal/input/test"
	"github.com/stretchr/testify/assert"
)

func TestHandler_GetMovementInput(t *testing.T) {
	h := input.New()
	angle := h.GetMovementInput()
	assert.InDelta(t, float64(common.Angle(0)), float64(angle), 0.001)
}

func TestHandler_KeyPress(t *testing.T) {
	t.Parallel()

	h := test.NewHandler()

	// Test key press
	assert.False(t, h.IsKeyPressed(ebiten.KeyA))
	h.SimulateKeyPress(ebiten.KeyA)
	assert.True(t, h.IsKeyPressed(ebiten.KeyA))

	// Test key release
	h.SimulateKeyRelease(ebiten.KeyA)
	assert.False(t, h.IsKeyPressed(ebiten.KeyA))
}

func TestHandler_IsQuitPressed(t *testing.T) {
	h := input.New()
	// Test with escape key not pressed
	assert.False(t, h.IsQuitPressed())
}

func TestHandler_IsPausePressed(t *testing.T) {
	h := input.New()
	// Test with space key not pressed
	assert.False(t, h.IsPausePressed())
}
