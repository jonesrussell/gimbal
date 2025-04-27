package game_test

import (
	"math"
	"testing"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/entity/stars"
	"github.com/jonesrussell/gimbal/internal/game"
	"github.com/jonesrussell/gimbal/internal/input"
	"github.com/jonesrussell/gimbal/internal/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	// Constants for initial angles
	InitialOrbitalAngle = 180.0 // Start at bottom of circle
	InitialFacingAngle  = 0.0   // Face upward
	RadiusDivisor       = 3.0   // Divisor for calculating orbit radius
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

// MockStarManager mocks the star manager for testing
type MockStarManager struct {
	mock.Mock
}

func (m *MockStarManager) Update()                   { m.Called() }
func (m *MockStarManager) Draw(screen any)           { m.Called(screen) }
func (m *MockStarManager) GetPosition() common.Point { return m.Called().Get(0).(common.Point) }
func (m *MockStarManager) GetStars() []*stars.Star   { return nil }
func (m *MockStarManager) Cleanup()                  { m.Called() }

// MockSprite mocks the player sprite for testing
type MockSprite struct {
	mock.Mock
}

func (m *MockSprite) Draw(screen any, op any) { m.Called(screen, op) }

// GameConfig holds the configuration for the game
type GameConfig struct {
	ScreenWidth  int
	ScreenHeight int
	PlayerSize   int
	NumStars     int
	StarSize     int
	StarSpeed    float64
	GameSpeed    time.Duration
	OrbitRadius  float64
}

// TestNew verifies game initialization
func TestNew(t *testing.T) {
	g := newTestGame(t)
	require.NotNil(t, g)

	// Verify initial state
	angle := float64(g.GetPlayer().GetAngle())
	facingAngle := float64(g.GetPlayer().GetFacingAngle())

	t.Logf("Player angle: %.2f (expected: %.2f)", angle, InitialOrbitalAngle)
	t.Logf("Player facing angle: %.2f (expected: %.2f)", facingAngle, InitialFacingAngle)

	// Compare angles with tolerance
	almostEqual := func(a, b float64) bool {
		diff := math.Abs(a - b)
		return diff < 0.0001 || diff > 359.9999
	}

	assert.True(t, almostEqual(angle, InitialOrbitalAngle),
		"Expected angle %v but got %v", InitialOrbitalAngle, angle)
	assert.True(t, almostEqual(facingAngle, InitialFacingAngle),
		"Expected facing angle %v but got %v", InitialFacingAngle, facingAngle)

	// Verify player position
	pos := g.GetPlayer().GetPosition()
	screenSize := g.GetScreenSize()
	expectedPos := common.Point{
		X: float64(screenSize.Width) / 2,
		Y: float64(screenSize.Height)/2 + float64(screenSize.Height)/RadiusDivisor,
	}

	t.Logf("Position: (%.2f, %.2f), Expected: (%.2f, %.2f)", pos.X, pos.Y, expectedPos.X, expectedPos.Y)
	t.Logf("Screen size: %dx%d, Radius: %.2f",
		screenSize.Width, screenSize.Height, g.GetRadius())

	assert.InDelta(t, expectedPos.X, pos.X, 0.1, "X position incorrect")
	assert.InDelta(t, expectedPos.Y, pos.Y, 0.1, "Y position incorrect")
}

func TestGame_Layout(t *testing.T) {
	g := newTestGame(t)
	screenSize := g.GetScreenSize()
	width, height := g.Layout(800, 600)
	assert.Equal(t, screenSize.Width, width)
	assert.Equal(t, screenSize.Height, height)
}

func TestGame_Update(t *testing.T) {
	t.Run("normal update", func(t *testing.T) {
		g := newTestGame(t)
		err := g.Update()
		assert.NoError(t, err)
	})

	t.Run("pause game", func(t *testing.T) {
		g := newTestGame(t)
		mockInput := new(MockInputHandler)
		mockInput.On("HandleInput").Return()
		mockInput.On("IsPausePressed").Return(true)
		mockInput.On("IsQuitPressed").Return(false)
		g.SetInputHandler(mockInput)

		err := g.Update()
		assert.NoError(t, err)
		assert.True(t, g.IsPaused())

		mockInput.AssertExpectations(t)
	})

	t.Run("quit game", func(t *testing.T) {
		g := newTestGame(t)
		mockInput := new(MockInputHandler)
		mockInput.On("HandleInput").Return()
		mockInput.On("IsPausePressed").Return(false)
		mockInput.On("IsQuitPressed").Return(true)
		g.SetInputHandler(mockInput)

		err := g.Update()
		assert.ErrorIs(t, err, game.ErrUserQuit)

		mockInput.AssertExpectations(t)
	})

	t.Run("movement update", func(t *testing.T) {
		g := newTestGame(t)
		mockInput := new(MockInputHandler)
		mockInput.On("HandleInput").Return()
		mockInput.On("IsPausePressed").Return(false)
		mockInput.On("IsQuitPressed").Return(false)
		mockInput.On("GetMovementInput").Return(common.Angle(5))
		g.SetInputHandler(mockInput)

		err := g.Update()
		assert.NoError(t, err)

		mockInput.AssertExpectations(t)
	})
}

func TestGame_Cleanup(t *testing.T) {
	g := newTestGame(t)
	g.Cleanup()
}

func TestGame_Input(t *testing.T) {
	// Create mock objects
	mockPlayer := new(MockPlayer)
	mockInput := new(MockInputHandler)
	mockStars := new(MockStarManager)
	mockLogger := new(mockLogger)

	// Configure screen dimensions and orbit radius
	screenWidth := 800
	screenHeight := 600
	orbitRadius := float64(screenHeight) / RadiusDivisor

	// Create game config
	config := &common.GameConfig{
		ScreenSize: common.Size{
			Width:  screenWidth,
			Height: screenHeight,
		},
		PlayerSize: common.Size{
			Width:  32,
			Height: 32,
		},
		Speed: 5.0,
	}

	// Calculate expected positions based on angle
	centerX := float64(screenWidth) / 2
	centerY := float64(screenHeight) / 2

	// Initial position at bottom of circle (180 degrees)
	initialPos := common.Point{
		X: centerX,
		Y: centerY + orbitRadius,
	}

	// Mock player expectations
	mockPlayer.On("GetPosition").Return(initialPos)
	mockPlayer.On("GetAngle").Return(common.Angle(180.0))
	mockPlayer.On("GetFacingAngle").Return(common.Angle(0.0))
	mockPlayer.On("SetAngle", mock.AnythingOfType("common.Angle")).Return(nil)
	mockPlayer.On("SetFacingAngle", mock.AnythingOfType("common.Angle")).Return()
	mockPlayer.On("Update").Return()

	// Mock input handler expectations
	mockInput.On("HandleInput").Return()
	mockInput.On("IsPausePressed").Return(false)
	mockInput.On("IsQuitPressed").Return(false)
	mockInput.On("GetMovementInput").Return(common.Angle(5.0))

	// Mock stars expectations
	mockStars.On("Update").Return()

	// Mock logger expectations
	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()

	// Create game instance with mocks
	game, err := game.NewWithDependencies(config, mockLogger, mockPlayer, mockStars, mockInput)
	require.NoError(t, err)
	require.NotNil(t, game)

	// Run update to simulate movement
	err = game.Update()
	assert.NoError(t, err)

	// Verify that the player's angle was updated
	mockPlayer.AssertCalled(t, "SetAngle", mock.AnythingOfType("common.Angle"))

	// Verify all expectations were met
	mockPlayer.AssertExpectations(t)
	mockInput.AssertExpectations(t)
	mockStars.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

// newTestGame creates a game instance with mocked components for testing
func newTestGame(t *testing.T) *game.GimlarGame {
	mockLogger := logger.NewMock()
	config := common.NewConfig()

	// Calculate initial position
	screenCenterX := float64(config.ScreenSize.Width) / 2
	screenCenterY := float64(config.ScreenSize.Height) / 2
	orbitRadius := float64(config.ScreenSize.Height) / game.RadiusDivisor

	// Initial position is at the bottom of the orbit (180 degrees)
	initialPos := common.Point{
		X: screenCenterX,
		Y: screenCenterY + orbitRadius,
	}

	// Create mock player
	mockPlayer := new(MockPlayer)
	mockPlayer.On("GetPosition").Return(initialPos)
	mockPlayer.On("GetAngle").Return(common.Angle(game.InitialOrbitalAngle))
	mockPlayer.On("GetFacingAngle").Return(common.Angle(game.InitialFacingAngle))
	mockPlayer.On("Update").Return()
	mockPlayer.On("Draw", mock.Anything, mock.Anything).Return()
	mockPlayer.On("SetAngle", mock.Anything).Return(nil)
	mockPlayer.On("SetFacingAngle", mock.Anything).Return()
	mockPlayer.On("Cleanup").Return()

	// Create mock star manager
	mockStars := new(MockStarManager)
	mockStars.On("Update").Return()
	mockStars.On("Draw", mock.Anything).Return()
	mockStars.On("GetStars").Return([]*stars.Star{})
	mockStars.On("Cleanup").Return()

	// Create game with mocks
	g, err := game.NewWithDependencies(config, mockLogger, mockPlayer, mockStars, input.New(mockLogger))
	require.NoError(t, err)

	return g
}
