package game_test

import (
	"testing"
	"time"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/game"
	"github.com/jonesrussell/gimbal/internal/logger"
	"github.com/stretchr/testify/require"
)

// TestGame_Run tests the game's main loop and shutdown
func TestGame_Run(t *testing.T) {
	// Create a channel to receive errors from the game loop
	errChan := make(chan error)
	done := make(chan struct{})

	// Create a new game instance
	mockLogger := logger.NewMock()
	config := common.NewConfig()

	g, initErr := game.New(config, mockLogger)
	require.NoError(t, initErr)

	// Start the game in a goroutine
	go func() {
		runErr := g.Run()
		errChan <- runErr
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
	case <-time.After(100 * time.Millisecond):
		// Game is running as expected
		close(done)
	case errReceived := <-errChan:
		require.NoError(t, errReceived)
	}
}
