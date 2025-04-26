package game_test

import (
	"math"
	"testing"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/game"
	"github.com/jonesrussell/gimbal/internal/input"
	"github.com/jonesrussell/gimbal/internal/logger"
	"github.com/jonesrussell/gimbal/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	// Constants for initial angles
	InitialOrbitalAngle = 180.0 // Start at bottom of circle
	InitialFacingAngle  = 0.0   // Face upward
	RadiusDivisor       = 4.0   // Divisor for calculating orbit radius
)

// MockPlayer is a mock implementation of player.PlayerInterface
type MockPlayer struct {
	mock.Mock
}

func (m *MockPlayer) Draw(screen, op any)                { m.Called(screen, op) }
func (m *MockPlayer) Update()                            { m.Called() }
func (m *MockPlayer) GetPosition() common.Point          { return m.Called().Get(0).(common.Point) }
func (m *MockPlayer) SetPosition(pos common.Point) error { return m.Called(pos).Error(0) }
func (m *MockPlayer) GetSpeed() float64                  { return m.Called().Get(0).(float64) }
func (m *MockPlayer) GetFacingAngle() common.Angle       { return m.Called().Get(0).(common.Angle) }
func (m *MockPlayer) SetFacingAngle(angle common.Angle)  { m.Called(angle) }
func (m *MockPlayer) GetAngle() common.Angle             { return m.Called().Get(0).(common.Angle) }
func (m *MockPlayer) SetAngle(angle common.Angle) error  { return m.Called(angle).Error(0) }
func (m *MockPlayer) GetBounds() common.Size             { return m.Called().Get(0).(common.Size) }
func (m *MockPlayer) Config() *common.EntityConfig       { return m.Called().Get(0).(*common.EntityConfig) }
func (m *MockPlayer) Cleanup()                           { m.Called() }

// MockInputHandler is a mock implementation of input.Interface
type MockInputHandler struct {
	mock.Mock
}

func (m *MockInputHandler) HandleInput()                     { m.Called() }
func (m *MockInputHandler) IsKeyPressed(key ebiten.Key) bool { return m.Called(key).Bool(0) }
func (m *MockInputHandler) GetMovementInput() common.Angle   { return m.Called().Get(0).(common.Angle) }
func (m *MockInputHandler) IsQuitPressed() bool              { return m.Called().Bool(0) }
func (m *MockInputHandler) IsPausePressed() bool             { return m.Called().Bool(0) }
func (m *MockInputHandler) GetTouchState() *input.TouchState {
	return m.Called().Get(0).(*input.TouchState)
}
func (m *MockInputHandler) GetMousePosition() common.Point { return m.Called().Get(0).(common.Point) }
func (m *MockInputHandler) IsMouseButtonPressed(button ebiten.MouseButton) bool {
	return m.Called(button).Bool(0)
}
func (m *MockInputHandler) GetLastEvent() input.InputEvent {
	return m.Called().Get(0).(input.InputEvent)
}
func (m *MockInputHandler) SimulateKeyPress(key ebiten.Key)   { m.Called(key) }
func (m *MockInputHandler) SimulateKeyRelease(key ebiten.Key) { m.Called(key) }

// mockLogger implements the logger interface for testing
type mockLogger struct {
	mock.Mock
}

func (m *mockLogger) Info(msg string, args ...any)  { m.Called(msg, args) }
func (m *mockLogger) Debug(msg string, args ...any) { m.Called(msg, args) }
func (m *mockLogger) Error(msg string, args ...any) { m.Called(msg, args) }
func (m *mockLogger) Warn(msg string, args ...any)  { m.Called(msg, args) }
func (m *mockLogger) Sync() error                   { return nil }

// GameConfig holds the configuration for the game
type GameConfig struct {
	ScreenWidth  int
	ScreenHeight int
	PlayerSize   int
	NumStars     int
	StarSize     int
	StarSpeed    float64
	GameSpeed    time.Duration
}

// TestNew verifies game initialization
func TestNew(t *testing.T) {
	// Helper function to compare float values with tolerance
	almostEqual := func(a, b float64) bool {
		return math.Abs(a-b) < 0.0001
	}

	// Create a mock logger
	logger := new(mockLogger)
	logger.On("Debug", mock.Anything, mock.Anything).Return()

	// Create game configuration
	cfg := &common.GameConfig{
		ScreenSize: common.Size{
			Width:  800,
			Height: 600,
		},
		PlayerSize: common.Size{
			Width:  32,
			Height: 32,
		},
		NumStars:  100,
		StarSize:  2.0,
		StarSpeed: 1.0,
		Speed:     1.0,
		Debug:     true,
		Radius:    float64(600) / common.CenterDivisor * common.DefaultRadiusRatio,
	}

	// Create game instance
	g, err := game.New(cfg, logger)
	require.NoError(t, err)
	require.NotNil(t, g)

	// Verify player initialization
	player := g.GetPlayer()
	require.NotNil(t, player)

	// Verify initial angles
	assert.True(t, almostEqual(float64(player.GetAngle()), InitialOrbitalAngle))
	assert.True(t, almostEqual(float64(player.GetFacingAngle()), InitialFacingAngle))

	// Verify initial position
	radius := float64(cfg.ScreenSize.Height) / RadiusDivisor
	expectedX := float64(cfg.ScreenSize.Width)/2 + radius*math.Cos(math.Pi)
	expectedY := float64(cfg.ScreenSize.Height)/2 + radius*math.Sin(math.Pi)

	pos := player.GetPosition()
	assert.True(t, almostEqual(pos.X, expectedX))
	assert.True(t, almostEqual(pos.Y, expectedY))
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
	mockLogger := logger.NewMock()
	config := common.NewConfig()

	tests := []struct {
		name           string
		setupInput     func(*MockInputHandler)
		expectedError  error
		expectedPaused bool
	}{
		{
			name: "normal update",
			setupInput: func(mi *MockInputHandler) {
				mi.On("HandleInput").Return()
				mi.On("IsPausePressed").Return(false)
				mi.On("IsQuitPressed").Return(false)
				mi.On("GetMovementInput").Return(common.Angle(0))
			},
			expectedError:  nil,
			expectedPaused: false,
		},
		{
			name: "pause game",
			setupInput: func(mi *MockInputHandler) {
				mi.On("HandleInput").Return()
				mi.On("IsPausePressed").Return(true)
				mi.On("IsQuitPressed").Return(false)
			},
			expectedError:  nil,
			expectedPaused: true,
		},
		{
			name: "quit game",
			setupInput: func(mi *MockInputHandler) {
				mi.On("HandleInput").Return()
				mi.On("IsPausePressed").Return(false)
				mi.On("IsQuitPressed").Return(true)
			},
			expectedError:  game.ErrUserQuit,
			expectedPaused: false,
		},
		{
			name: "movement update",
			setupInput: func(mi *MockInputHandler) {
				mi.On("HandleInput").Return()
				mi.On("IsPausePressed").Return(false)
				mi.On("IsQuitPressed").Return(false)
				mi.On("GetMovementInput").Return(common.Angle(5))
			},
			expectedError:  nil,
			expectedPaused: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockInput := new(MockInputHandler)

			if tt.setupInput != nil {
				tt.setupInput(mockInput)
			}

			g, err := game.New(config, mockLogger)
			require.NoError(t, err)

			// Set the input handler
			g.SetInputHandler(mockInput)

			err = g.Update()
			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedPaused, g.IsPaused())
			mockInput.AssertExpectations(t)
		})
	}
}

func TestGame_Draw(t *testing.T) {
	mockLogger := logger.NewMock()
	config := common.NewConfig()

	tests := []struct {
		name   string
		screen *ebiten.Image
	}{
		{
			name:   "nil screen",
			screen: nil,
		},
		{
			name:   "normal draw",
			screen: ebiten.NewImage(640, 480),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, err := game.New(config, mockLogger)
			require.NoError(t, err)

			g.Draw(tt.screen)
		})
	}
}

func TestGame_Cleanup(t *testing.T) {
	mockLogger := logger.NewMock()
	config := common.NewConfig()

	g, err := game.New(config, mockLogger)
	require.NoError(t, err)

	g.Cleanup()
}

func TestGame_Run(t *testing.T) {
	mockLogger := logger.NewMock()
	config := common.NewConfig()

	g, err := game.New(config, mockLogger)
	require.NoError(t, err)

	// Start the game in a goroutine
	errChan := make(chan error)
	go func() {
		errChan <- g.Run()
	}()

	// Wait a short time to let the game initialize
	time.Sleep(100 * time.Millisecond)

	// Simulate a clean shutdown
	mockInput := new(MockInputHandler)
	mockInput.On("HandleInput").Return()
	mockInput.On("IsPausePressed").Return(false)
	mockInput.On("IsQuitPressed").Return(true)
	g.SetInputHandler(mockInput)

	// Check the result
	select {
	case err := <-errChan:
		assert.NoError(t, err)
	case <-time.After(time.Second):
		t.Fatal("game.Run did not complete in time")
	}
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

	testHandler := input.New(mockLogger)
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
	assert.False(t, g.IsPaused())
}
