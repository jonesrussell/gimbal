package game

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jonesrussell/gimbal/internal/config"
	"go.uber.org/zap"
)

func setupTest(t *testing.T) (*GimlarGame, *zap.Logger) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	cfg := &config.Config{
		Screen: struct {
			Title  string `json:"title"`
			Width  int    `json:"width"`
			Height int    `json:"height"`
		}{
			Title:  "Test Game",
			Width:  640,
			Height: 480,
		},
		Game: struct {
			Speed    float64 `json:"speed"`
			NumStars int     `json:"numStars"`
			Debug    bool    `json:"debug"`
		}{
			Speed:    1.0,
			NumStars: 100,
			Debug:    true,
		},
		Player: struct {
			Width  int `json:"width"`
			Height int `json:"height"`
		}{
			Width:  16,
			Height: 16,
		},
	}

	game, err := NewGimlarGame(logger, cfg)
	if err != nil {
		t.Fatalf("Failed to create game: %v", err)
	}

	return game, logger
}

func TestNewGimlarGame(t *testing.T) {
	game, _ := setupTest(t)
	if game == nil {
		t.Error("Expected game to not be nil")
	}
}

func TestUpdateStars(t *testing.T) {
	game, _ := setupTest(t)

	// Store initial positions
	initialPositions := make([]struct{ x, y float64 }, len(game.stars))
	for i, star := range game.stars {
		initialPositions[i].x = star.X
		initialPositions[i].y = star.Y
	}

	// Update stars
	game.updateStars()

	// Verify stars have moved
	for i, star := range game.stars {
		if star.X == initialPositions[i].x && star.Y == initialPositions[i].y {
			t.Errorf("Star %d did not move", i)
		}
	}
}

func TestDrawStars(t *testing.T) {
	game, _ := setupTest(t)
	screen := ebiten.NewImage(640, 480)

	// This should not panic
	game.drawStars(screen)
}

func TestLayout(t *testing.T) {
	game, _ := setupTest(t)
	width, height := game.Layout(800, 600)

	cfg, err := config.New()
	if err != nil {
		t.Fatalf("Failed to get config: %v", err)
	}

	if width != cfg.Screen.Width || height != cfg.Screen.Height {
		t.Errorf("Expected %dx%d, got %dx%d",
			cfg.Screen.Width, cfg.Screen.Height,
			width, height)
	}
}

func TestUpdate(t *testing.T) {
	game, _ := setupTest(t)

	err := game.Update()
	if err != nil {
		t.Errorf("Update returned error: %v", err)
	}
}
