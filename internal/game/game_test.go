package game_test

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/game"
	inputtest "github.com/jonesrussell/gimbal/internal/input/test"
	"github.com/jonesrussell/gimbal/internal/logger"
	"github.com/jonesrussell/gimbal/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()

	mockLogger := logger.NewMock()

	tests := []struct {
		name    string
		config  *common.GameConfig
		wantErr bool
	}{
		{
			name:    "nil config",
			config:  nil,
			wantErr: true,
		},
		{
			name:    "valid config",
			config:  common.NewConfig(),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := game.New(tt.config, mockLogger)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, got)
		})
	}
}

func TestGame_Layout(t *testing.T) {
	t.Parallel()

	mockLogger := logger.NewMock()

	config := common.NewConfig()
	g, err := game.New(config, mockLogger)
	require.NoError(t, err)

	width, height := g.Layout(800, 600)
	assert.Equal(t, config.ScreenSize.Width, width)
	assert.Equal(t, config.ScreenSize.Height, height)
}

func TestGame_Update(t *testing.T) {
	test.EnsureXvfb(t)
	t.Parallel()

	mockLogger := logger.NewMock()

	config := common.NewConfig(
		common.WithScreenSize(640, 480),
		common.WithDebug(true),
	)

	g, err := game.New(config, mockLogger)
	require.NoError(t, err)

	// Test player movement
	testInput := inputtest.NewHandler()
	g.SetInputHandler(testInput)

	// Get initial position
	initialPos := g.GetPlayer().GetPosition()

	// Test right movement
	testInput.SimulateKeyPress(ebiten.KeyRight)
	g.Update()
	rightPos := g.GetPlayer().GetPosition()
	// Check that the position has changed
	assert.NotEqual(t, initialPos, rightPos)

	// Test no movement after release
	testInput.SimulateKeyRelease(ebiten.KeyRight)
	g.Update()
	releasePos := g.GetPlayer().GetPosition()
	assert.Equal(t, rightPos, releasePos)
}

func TestGame_Draw(t *testing.T) {
	test.EnsureXvfb(t)
	t.Parallel()

	mockLogger := logger.NewMock()

	config := common.NewConfig(
		common.WithScreenSize(640, 480),
		common.WithDebug(true),
	)

	g, err := game.New(config, mockLogger)
	require.NoError(t, err)

	// Test drawing with nil screen
	g.Draw(nil)
	// No panic expected
}

func TestGame_Input(t *testing.T) {
	test.EnsureXvfb(t)
	t.Parallel()

	mockLogger := logger.NewMock()

	config := common.NewConfig(
		common.WithScreenSize(640, 480),
		common.WithDebug(true),
	)

	g, err := game.New(config, mockLogger)
	require.NoError(t, err)

	testHandler := inputtest.NewHandler()
	g.SetInputHandler(testHandler)

	// Get initial position
	initialPos := g.GetPlayer().GetPosition()

	// Test left movement
	testHandler.SimulateKeyPress(ebiten.KeyLeft)
	g.Update()
	leftPos := g.GetPlayer().GetPosition()
	// Check that the position has changed
	assert.NotEqual(t, initialPos, leftPos)

	// Test no movement after release
	testHandler.SimulateKeyRelease(ebiten.KeyLeft)
	g.Update()
	releasePos := g.GetPlayer().GetPosition()
	assert.Equal(t, leftPos, releasePos)

	// Test right movement
	testHandler.SimulateKeyPress(ebiten.KeyRight)
	g.Update()
	rightPos := g.GetPlayer().GetPosition()
	// Check that the position has changed
	assert.NotEqual(t, leftPos, rightPos)

	// Test no movement after release
	testHandler.SimulateKeyRelease(ebiten.KeyRight)
	g.Update()
	rightReleasePos := g.GetPlayer().GetPosition()
	assert.Equal(t, rightPos, rightReleasePos)

	// Test space key (pause)
	testHandler.SimulateKeyPress(ebiten.KeySpace)
	g.Update()
	assert.True(t, g.IsPaused())

	testHandler.SimulateKeyRelease(ebiten.KeySpace)
	g.Update()
	assert.True(t, g.IsPaused()) // Should stay paused until space is pressed again
}
