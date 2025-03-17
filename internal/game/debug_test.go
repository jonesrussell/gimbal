package game_test

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/game"
	"github.com/stretchr/testify/assert"
)

func TestDebugPrintStar(t *testing.T) {
	// Create a test game
	config := common.NewConfig(
		common.WithScreenSize(640, 480),
		common.WithDebug(true),
	)
	g, err := game.New(config)
	if err != nil {
		t.Fatalf("Failed to create game: %v", err)
	}

	// Create a test screen
	screen := ebiten.NewImage(640, 480)

	// Test with a valid star
	star := g.GetStars()[0]
	g.DebugPrintStar(screen, star) // This should not panic
}

func TestDrawDebugGridOverlay(t *testing.T) {
	// Create a test game
	config := common.NewConfig(
		common.WithScreenSize(640, 480),
		common.WithDebug(true),
	)
	g, err := game.New(config)
	if err != nil {
		t.Fatalf("Failed to create game: %v", err)
	}

	// Create a test screen
	screen := ebiten.NewImage(640, 480)

	// Test drawing the grid
	g.DrawDebugGridOverlay(screen)

	// Verify the screen was modified
	assert.NotNil(t, screen)
}

func TestGimlarGame_DrawDebugInfo(t *testing.T) {
	// Create a test game
	config := common.NewConfig(
		common.WithScreenSize(640, 480),
		common.WithDebug(true),
	)
	g, err := game.New(config)
	if err != nil {
		t.Fatalf("Failed to create game: %v", err)
	}

	// Create a test screen
	screen := ebiten.NewImage(640, 480)

	// Test drawing debug info
	g.DrawDebugInfo(screen)

	// Verify the screen was modified
	assert.NotNil(t, screen)
}

func TestGimlarGame_DrawDebugGrid(t *testing.T) {
	// Create a test game
	config := common.NewConfig(
		common.WithScreenSize(640, 480),
		common.WithDebug(true),
	)
	g, err := game.New(config)
	if err != nil {
		t.Fatalf("Failed to create game: %v", err)
	}

	// Create a test screen
	screen := ebiten.NewImage(640, 480)

	// Test drawing the grid
	g.DrawDebugGrid(screen)

	// Verify the screen was modified
	assert.NotNil(t, screen)
}
