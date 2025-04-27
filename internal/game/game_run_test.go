package game_test

import (
	"testing"
	"time"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/game"
	"github.com/jonesrussell/gimbal/internal/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGame_Run tests the game's main loop and shutdown
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
