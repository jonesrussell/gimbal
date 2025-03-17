package game_test

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/game"
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

	// Get initial player position
	initialPos := g.GetPlayer().GetPosition()
	initialX := initialPos.X

	// Execute
	g.Update()

	// Assert
	finalPos := g.GetPlayer().GetPosition()
	assert.NotEqual(t, initialX, finalPos.X)
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
