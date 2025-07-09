package app

import (
	"context"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/jonesrussell/gimbal/test/mocks"
)

func TestContainerInitialization(t *testing.T) {
	// Create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create container
	container := NewContainer()

	// Test initial state
	if container.IsInitialized() {
		t.Error("Container should not be initialized initially")
	}

	// Test initialization
	ctx := context.Background()
	err := container.Initialize(ctx)
	if err != nil {
		t.Fatalf("Failed to initialize container: %v", err)
	}

	// Test initialized state
	if !container.IsInitialized() {
		t.Error("Container should be initialized after Initialize()")
	}

	// Test that dependencies are available
	if container.GetLogger() == nil {
		t.Error("Logger should be available after initialization")
	}

	if container.GetConfig() == nil {
		t.Error("Config should be available after initialization")
	}

	if container.GetInputHandler() == nil {
		t.Error("Input handler should be available after initialization")
	}

	if container.GetGame() == nil {
		t.Error("Game should be available after initialization")
	}
}

func TestContainerDoubleInitialization(t *testing.T) {
	// Create container
	container := NewContainer()

	// Initialize once
	ctx := context.Background()
	err := container.Initialize(ctx)
	if err != nil {
		t.Fatalf("Failed to initialize container: %v", err)
	}

	// Try to initialize again
	err = container.Initialize(ctx)
	if err == nil {
		t.Error("Should return error when initializing already initialized container")
	}
}

func TestContainerShutdown(t *testing.T) {
	// Create container
	container := NewContainer()

	// Initialize
	ctx := context.Background()
	err := container.Initialize(ctx)
	if err != nil {
		t.Fatalf("Failed to initialize container: %v", err)
	}

	// Test shutdown
	err = container.Shutdown(ctx)
	if err != nil {
		t.Fatalf("Failed to shutdown container: %v", err)
	}

	// Test double shutdown (should not error)
	err = container.Shutdown(ctx)
	if err != nil {
		t.Errorf("Double shutdown should not error: %v", err)
	}
}

func TestContainerSetInputHandler(t *testing.T) {
	// Create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create container
	container := NewContainer()

	// Initialize
	ctx := context.Background()
	err := container.Initialize(ctx)
	if err != nil {
		t.Fatalf("Failed to initialize container: %v", err)
	}

	// Create mock input handler
	mockInput := mocks.NewMockGameInputHandler(ctrl)

	// Test setting custom input handler
	container.SetInputHandler(mockInput)

	// Verify the input handler was set
	if container.GetInputHandler() != mockInput {
		t.Error("Input handler should be updated after SetInputHandler")
	}
}

func TestContainerDependencyOrder(t *testing.T) {
	// Create container
	container := NewContainer()

	// Initialize
	ctx := context.Background()
	err := container.Initialize(ctx)
	if err != nil {
		t.Fatalf("Failed to initialize container: %v", err)
	}

	// Test that dependencies are properly initialized in order
	logger := container.GetLogger()
	config := container.GetConfig()
	inputHandler := container.GetInputHandler()
	game := container.GetGame()

	// All dependencies should be non-nil
	if logger == nil || config == nil || inputHandler == nil || game == nil {
		t.Error("All dependencies should be properly initialized")
	}

	// Test that game has the correct input handler
	// This tests that the dependency injection is working correctly
	if game.GetInputHandler() != inputHandler {
		t.Error("Game should have the same input handler as the container")
	}
}
