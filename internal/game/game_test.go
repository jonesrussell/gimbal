package game_test

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jonesrussell/gimbal/internal/game"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewGimlarGame(t *testing.T) {
	// Setup
	config := game.NewConfig(
		game.WithScreenSize(640, 480),
		game.WithPlayerSize(16, 16),
		game.WithSpeed(1.0),
	)
	input := &game.InputHandler{}

	// Execute
	g, err := game.NewGimlarGame(config, input)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, g)
	assert.InEpsilon(t, config.Speed, g.GetSpeed(), 0.0001)
	assert.NotNil(t, g.GetPlayer())
	assert.NotNil(t, g.GetStars())
	assert.NotNil(t, g.GetSpace())
}

func TestUpdate(t *testing.T) {
	// Setup
	config := game.NewConfig(
		game.WithScreenSize(640, 480),
		game.WithSpeed(1.0),
	)
	input := &game.InputHandler{}
	g, err := game.NewGimlarGame(config, input)
	require.NoError(t, err)

	// Execute
	err = g.Update()

	// Assert
	require.NoError(t, err)
}

func TestDraw(t *testing.T) {
	// Setup
	config := game.NewConfig(
		game.WithScreenSize(640, 480),
	)
	input := &game.InputHandler{}
	g, err := game.NewGimlarGame(config, input)
	require.NoError(t, err)

	// Execute
	image := ebiten.NewImage(100, 100)
	g.Draw(image)

	// Assert
	assert.NotNil(t, image)
}

func TestPlayerMovement(t *testing.T) {
	// Setup
	config := game.NewConfig(
		game.WithScreenSize(640, 480),
		game.WithSpeed(1.0),
	)
	input := &game.InputHandler{}
	g, err := game.NewGimlarGame(config, input)
	require.NoError(t, err)

	// Initial player position
	player := g.GetPlayer()
	initialPos := player.Object.Position()
	initialX := initialPos.X

	// Execute
	player.Object.Move(g.GetSpeed(), 0)

	// Assert
	finalPos := player.Object.Position()
	assert.NotEqual(t, initialX, finalPos.X)
}

func TestStarMovement(t *testing.T) {
	// Setup
	config := game.NewConfig(
		game.WithScreenSize(640, 480),
		game.WithStarSettings(5.0, 2.0),
	)
	input := &game.InputHandler{}
	g, err := game.NewGimlarGame(config, input)
	require.NoError(t, err)

	// Get stars from the game
	stars := g.GetStars()

	// Initial star position
	initialPosition := stars[0].X

	// Execute
	stars[0].X += stars[0].Speed

	// Assert
	finalPosition := stars[0].X
	assert.NotEqual(t, initialPosition, finalPosition)
}

func TestGimlarGame_Update(t *testing.T) {
	// Setup
	config := game.NewConfig(
		game.WithScreenSize(640, 480),
		game.WithSpeed(1.0),
	)
	input := &game.InputHandler{}
	g, err := game.NewGimlarGame(config, input)
	require.NoError(t, err)

	// Execute
	err = g.Update()

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, g.GetPlayer())
}

func TestGimlarGame_Draw(t *testing.T) {
	// Setup
	config := game.NewConfig(
		game.WithScreenSize(640, 480),
	)
	input := &game.InputHandler{}
	g, err := game.NewGimlarGame(config, input)
	require.NoError(t, err)

	// Create a test screen
	screen := ebiten.NewImage(config.ScreenWidth, config.ScreenHeight)

	// Execute
	g.Draw(screen)

	// Assert
	assert.NotNil(t, g.GetPlayer())
	assert.NotNil(t, screen)
}

func TestGimlarGame_Layout(t *testing.T) {
	// Setup
	config := game.NewConfig(
		game.WithScreenSize(640, 480),
	)
	input := &game.InputHandler{}
	g, err := game.NewGimlarGame(config, input)
	require.NoError(t, err)

	// Execute
	width, height := g.Layout(800, 600)

	// Assert
	assert.Equal(t, config.ScreenWidth, width)
	assert.Equal(t, config.ScreenHeight, height)
}

func TestGimlarGame_GetRadius(t *testing.T) {
	// Setup
	config := game.NewConfig(
		game.WithScreenSize(640, 480),
	)
	input := &game.InputHandler{}
	g, err := game.NewGimlarGame(config, input)
	require.NoError(t, err)

	// Execute
	radius := g.GetRadius()

	// Assert
	assert.Greater(t, radius, 0.0)
	assert.Less(t, radius, float64(config.ScreenHeight))
}

func TestGameConfig_Options(t *testing.T) {
	// Test different configuration options
	tests := []struct {
		name     string
		opts     []game.GameOption
		validate func(*testing.T, *game.GameConfig)
	}{
		{
			name: "custom screen size",
			opts: []game.GameOption{
				game.WithScreenSize(800, 600),
			},
			validate: func(t *testing.T, c *game.GameConfig) {
				assert.Equal(t, 800, c.ScreenWidth)
				assert.Equal(t, 600, c.ScreenHeight)
				assert.InDelta(t, 225.0, c.Radius, 0.001) // 0.75 * 600/2
			},
		},
		{
			name: "custom player size",
			opts: []game.GameOption{
				game.WithPlayerSize(32, 32),
			},
			validate: func(t *testing.T, c *game.GameConfig) {
				assert.Equal(t, 32, c.PlayerWidth)
				assert.Equal(t, 32, c.PlayerHeight)
			},
		},
		{
			name: "custom star settings",
			opts: []game.GameOption{
				game.WithStarSettings(10.0, 3.0),
			},
			validate: func(t *testing.T, c *game.GameConfig) {
				assert.InDelta(t, 10.0, c.StarSize, 0.001)
				assert.InDelta(t, 3.0, c.StarSpeed, 0.001)
			},
		},
		{
			name: "debug mode",
			opts: []game.GameOption{
				game.WithDebug(true),
			},
			validate: func(t *testing.T, c *game.GameConfig) {
				assert.True(t, c.Debug)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := game.NewConfig(tt.opts...)
			tt.validate(t, config)
		})
	}
}
