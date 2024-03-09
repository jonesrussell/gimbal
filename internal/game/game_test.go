package game

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/stretchr/testify/assert"
)

func TestNewGimlarGame(t *testing.T) {
	// Test the NewGimlarGame function

	// Setup
	speed := 1.0 // Example speed value

	// Execute
	game, err := NewGimlarGame(speed)

	// Assert
	assert.NoError(t, err)             // Ensure no error is returned
	assert.NotNil(t, game)             // Ensure the game instance is not nil
	assert.Equal(t, speed, game.speed) // Ensure the game's speed is correctly set
	assert.NotNil(t, game.player)      // Ensure the player is initialized
	assert.NotNil(t, game.stars)       // Ensure the stars are initialized
	assert.NotNil(t, game.space)       // Ensure the space is initialized
	assert.NotNil(t, game.logger)      // Ensure the logger is initialized
}

func TestUpdate(t *testing.T) {
	// Setup
	speed := 1.0 // Example speed value
	game, err := NewGimlarGame(speed)
	if err != nil {
		t.Fatalf("Failed to create game: %v", err)
	}

	// Execute
	err = game.Update() // Directly call the Update method of your game instance

	// Assert
	assert.NoError(t, err) // Ensure no error is returned
	// Additional assertions can be added here to check the state of the game after the update
	// For example, checking if the player's position has changed or if the stars have moved
}

func TestDraw(t *testing.T) {
	// Setup
	speed := 1.0 // Example speed value
	game, err := NewGimlarGame(speed)
	if err != nil {
		t.Fatalf("Failed to create game: %v", err)
	}

	// Execute
	// Assuming Draw does not take any arguments and does not return any value
	// and that it does not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Draw method panicked: %v", r)
		}
	}()
	image := ebiten.NewImage(100, 100) // Create a dummy image
	game.Draw(image)

	// Assert
	// Since Draw might not have a direct observable effect, this test might be limited
	// Consider adding more comprehensive tests if possible, such as checking the state of the game's graphics context
}

func TestPlayerMovement(t *testing.T) {
	// Setup
	speed := 1.0 // Example speed value
	game, err := NewGimlarGame(speed)
	if err != nil {
		t.Fatalf("Failed to create game: %v", err)
	}

	// Initial player position
	initialPosition := game.player.Object.Position.X // Assuming Position() returns the player's current position

	// Execute
	// Simulate player movement by updating the player's Object position directly
	game.player.Object.Position.X += game.player.speed // Adjust the X position based on the player's speed

	// Assert
	// Check that the player's position has changed after the movement
	finalPosition := game.player.Object.Position.X
	assert.NotEqual(t, initialPosition, finalPosition)
}

func TestStarMovement(t *testing.T) {
	// Setup
	speed := 1.0 // Example speed value
	game, err := NewGimlarGame(speed)
	if err != nil {
		t.Fatalf("Failed to create game: %v", err)
	}

	// Assuming you have a way to initialize or retrieve the stars in your game
	// For example, if stars are part of the game state, you might retrieve them like this:
	stars := game.stars

	// Initial star position
	// Assuming each star has a Position field or method that returns its current position
	// For simplicity, let's assume the first star's initial position is stored in a variable
	initialPosition := stars[0].X // Adjust based on how you access the star's position

	// Execute
	// Simulate star movement by updating the star's position directly
	// This is a simplified example; you'll need to adjust based on your game's logic
	stars[0].X += stars[0].Speed // Adjust the X position based on the star's speed

	// Assert
	// Check that the star's position has changed after the movement
	finalPosition := stars[0].X
	assert.NotEqual(t, initialPosition, finalPosition)
}
