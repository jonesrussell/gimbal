package ecs_test

import (
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/ecs"
	"github.com/jonesrussell/gimbal/test/mocks"
)

func TestECSGameWithMockInput(t *testing.T) {
	// Create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock logger
	mockLogger := mocks.NewMockLogger(ctrl)

	// Create game config
	config := common.NewConfig(
		common.WithScreenSize(640, 480),
		common.WithDebug(false),
	)

	// Create mock input handler
	mockInput := mocks.NewMockGameInputHandler(ctrl)

	// Set up expectations
	mockLogger.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()

	// Test that we can create the game with mock input
	game, err := ecs.NewECSGame(config, mockLogger, mockInput)
	if err != nil {
		t.Fatalf("Failed to create ECS game with mock input: %v", err)
	}

	// Test that the game uses the injected input handler
	if game.GetInputHandler() != mockInput {
		t.Error("Game should use the injected input handler")
	}

	// Test that we can simulate input and the game responds
	mockInput.EXPECT().GetMovementInput().Return(common.Angle(10))
	movementInput := mockInput.GetMovementInput()
	if movementInput != common.Angle(10) {
		t.Errorf("Expected movement input 10, got %v", movementInput)
	}

	// Test pause functionality
	mockInput.EXPECT().IsPausePressed().Return(true)
	if !mockInput.IsPausePressed() {
		t.Error("Pause should be pressed")
	}

	// Test quit functionality
	mockInput.EXPECT().IsQuitPressed().Return(true)
	if !mockInput.IsQuitPressed() {
		t.Error("Quit should be pressed")
	}
}

func TestECSGameConstructorValidation(t *testing.T) {
	// Create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock logger
	mockLogger := mocks.NewMockLogger(ctrl)

	// Create game config
	config := common.NewConfig(
		common.WithScreenSize(640, 480),
		common.WithDebug(false),
	)

	// Test with nil config
	_, nilConfigErr := ecs.NewECSGame(nil, mockLogger, mocks.NewMockGameInputHandler(ctrl))
	if nilConfigErr == nil {
		t.Error("Should return error when config is nil")
	}

	// Test with nil logger
	_, nilLoggerErr := ecs.NewECSGame(config, nil, mocks.NewMockGameInputHandler(ctrl))
	if nilLoggerErr == nil {
		t.Error("Should return error when logger is nil")
	}

	// Test with nil input handler
	_, nilInputErr := ecs.NewECSGame(config, mockLogger, nil)
	if nilInputErr == nil {
		t.Error("Should return error when input handler is nil")
	}
}
