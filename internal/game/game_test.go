package game_test

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/game"
	"github.com/jonesrussell/gimbal/internal/game/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewGimlarGame(t *testing.T) {
	// Setup
	config := common.NewConfig(
		common.WithScreenSize(640, 480),
		common.WithPlayerSize(16, 16),
		common.WithSpeed(1.0),
	)

	// Execute
	g, err := game.New(config)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, g)
	assert.NotNil(t, g.GetStars())
}

func TestUpdate(t *testing.T) {
	// Setup
	config := common.NewConfig(
		common.WithScreenSize(640, 480),
		common.WithSpeed(1.0),
	)
	g, err := game.New(config)
	require.NoError(t, err)

	// Execute
	err = g.Update()

	// Assert
	require.NoError(t, err)
}

func TestDraw(t *testing.T) {
	// Setup
	config := common.NewConfig(
		common.WithScreenSize(640, 480),
	)
	g, err := game.New(config)
	require.NoError(t, err)

	// Execute
	image := ebiten.NewImage(100, 100)
	g.Draw(image)

	// Assert
	assert.NotNil(t, image)
}

func TestPlayerMovement(t *testing.T) {
	// Setup
	config := common.NewConfig(
		common.WithScreenSize(640, 480),
		common.WithSpeed(1.0),
	)
	g, err := game.New(config)
	require.NoError(t, err)

	// Create test input handler
	testInput := test.NewTestHandler()
	g.SetInputHandler(testInput)

	// Get initial player position and angle
	initialPos := g.GetPlayer().GetPosition()
	initialAngle := g.GetPlayer().GetAngle()

	// Simulate right movement input
	testInput.SimulateKeyPress(ebiten.KeyRight)

	// Execute update
	err = g.Update()
	require.NoError(t, err)

	// Get final position and angle
	finalPos := g.GetPlayer().GetPosition()
	finalAngle := g.GetPlayer().GetAngle()

	// The player should have moved in a circular path
	assert.NotEqual(t, initialPos, finalPos, "Player position should change")
	assert.NotEqual(t, initialAngle, finalAngle, "Player angle should change")

	// Cleanup
	testInput.SimulateKeyRelease(ebiten.KeyRight)
}

func TestStarMovement(t *testing.T) {
	// Setup
	config := common.NewConfig(
		common.WithScreenSize(640, 480),
		common.WithStarSettings(5.0, 2.0),
	)
	g, err := game.New(config)
	require.NoError(t, err)

	// Get initial star position
	stars := g.GetStars()
	require.NotEmpty(t, stars)
	initialPos := stars[0].GetPosition()

	// Execute
	g.Update()

	// Assert
	finalPos := stars[0].GetPosition()
	assert.NotEqual(t, initialPos, finalPos)
}

func TestGimlarGame_Update(t *testing.T) {
	// Setup
	config := common.NewConfig(
		common.WithScreenSize(640, 480),
		common.WithSpeed(1.0),
	)
	g, err := game.New(config)
	require.NoError(t, err)

	// Execute
	err = g.Update()

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, g.GetPlayer())
}

func TestGimlarGame_Draw(t *testing.T) {
	// Setup
	config := common.NewConfig(
		common.WithScreenSize(640, 480),
	)
	g, err := game.New(config)
	require.NoError(t, err)

	// Create a test screen
	screen := ebiten.NewImage(640, 480)

	// Execute
	g.Draw(screen)

	// Assert
	assert.NotNil(t, g.GetPlayer())
	assert.NotNil(t, screen)
}

func TestGimlarGame_Layout(t *testing.T) {
	// Setup
	config := common.NewConfig(
		common.WithScreenSize(640, 480),
	)
	g, err := game.New(config)
	require.NoError(t, err)

	// Execute
	width, height := g.Layout(800, 600)

	// Assert
	assert.Equal(t, config.ScreenSize.Width, width)
	assert.Equal(t, config.ScreenSize.Height, height)
}

func TestGameConfig_Options(t *testing.T) {
	// Test different configuration options
	tests := []struct {
		name     string
		opts     []common.GameOption
		validate func(*testing.T, *common.GameConfig)
	}{
		{
			name: "custom screen size",
			opts: []common.GameOption{
				common.WithScreenSize(800, 600),
			},
			validate: func(t *testing.T, c *common.GameConfig) {
				assert.Equal(t, 800, c.ScreenSize.Width)
				assert.Equal(t, 600, c.ScreenSize.Height)
				assert.InDelta(t, 225.0, c.Radius, 0.001) // 0.75 * 600/2
			},
		},
		{
			name: "custom player size",
			opts: []common.GameOption{
				common.WithPlayerSize(32, 32),
			},
			validate: func(t *testing.T, c *common.GameConfig) {
				assert.Equal(t, 32, c.PlayerSize.Width)
				assert.Equal(t, 32, c.PlayerSize.Height)
			},
		},
		{
			name: "custom star settings",
			opts: []common.GameOption{
				common.WithStarSettings(10.0, 3.0),
			},
			validate: func(t *testing.T, c *common.GameConfig) {
				assert.InDelta(t, 10.0, c.StarSize, 0.001)
				assert.InDelta(t, 3.0, c.StarSpeed, 0.001)
			},
		},
		{
			name: "debug mode",
			opts: []common.GameOption{
				common.WithDebug(true),
			},
			validate: func(t *testing.T, c *common.GameConfig) {
				assert.True(t, c.Debug)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := common.NewConfig(tt.opts...)
			tt.validate(t, config)
		})
	}
}

func TestGame_Update(t *testing.T) {
	t.Parallel()

	// Create a new game instance with default config
	config := common.DefaultConfig()
	g, err := game.New(config)
	require.NoError(t, err)

	// Create a test input handler
	testHandler := test.NewTestHandler()
	g.SetInputHandler(testHandler)

	// Test initial state
	assert.InDelta(t, float64(common.Angle(0)), float64(g.GetPlayer().GetAngle()), 0.001)
	assert.InDelta(t, float64(common.Angle(0)), float64(g.GetPlayer().GetFacingAngle()), 0.001)

	// Test left movement
	testHandler.SimulateKeyPress(ebiten.KeyLeft)
	g.Update()
	assert.InDelta(t, float64(common.Angle(-2)), float64(g.GetPlayer().GetAngle()), 0.001)
	assert.InDelta(t, float64(common.Angle(0)), float64(g.GetPlayer().GetFacingAngle()), 0.001)

	// Test right movement
	testHandler.SimulateKeyRelease(ebiten.KeyLeft)
	testHandler.SimulateKeyPress(ebiten.KeyRight)
	g.Update()
	assert.InDelta(t, float64(common.Angle(0)), float64(g.GetPlayer().GetAngle()), 0.001)
	assert.InDelta(t, float64(common.Angle(0)), float64(g.GetPlayer().GetFacingAngle()), 0.001)

	// Test pause
	testHandler.SimulateKeyRelease(ebiten.KeyRight)
	testHandler.SimulateKeyPress(ebiten.KeySpace)
	g.Update()
	require.True(t, g.IsPaused())

	// Test unpause
	testHandler.SimulateKeyRelease(ebiten.KeySpace)
	g.Update()
	require.False(t, g.IsPaused())
}
