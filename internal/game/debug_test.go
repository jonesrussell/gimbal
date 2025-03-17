package game_test

import (
	"testing"

	"github.com/jonesrussell/gimbal/internal/game"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDebug(t *testing.T) {
	t.Parallel()

	d := game.NewDebug()
	require.NotNil(t, d)
}

func TestDebug_Draw(t *testing.T) {
	t.Parallel()

	d := game.NewDebug()
	require.NotNil(t, d)

	// Test drawing with nil screen
	d.Draw(nil)
	// No panic expected
}

func TestDebug_Update(t *testing.T) {
	t.Parallel()

	d := game.NewDebug()
	require.NotNil(t, d)

	d.SetFPS(60)
	d.SetEntityCount(10)

	d.Update()
	assert.Equal(t, 60, d.GetFPS())
	assert.Equal(t, 10, d.GetEntityCount())
}
