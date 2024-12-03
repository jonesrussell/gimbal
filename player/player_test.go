package player

import (
	"image"
	"testing"
)

func TestPlayer_Input(t *testing.T) {
	mock := NewMockHandler()
	speed := 5.0
	img := image.NewRGBA(image.Rect(0, 0, 32, 32))
	center := image.Pt(320, 240)

	player, err := NewPlayer(mock, speed, img, center)
	if err != nil {
		t.Fatalf("Failed to create player: %v", err)
	}

	if player == nil {
		t.Fatalf("Expected player to not be nil")
		return
	}

	if player.Object == nil {
		t.Fatal("Player Object should not be nil")
	}

	expectedX := float64(center.X) + radius*(-1.0)
	expectedY := float64(center.Y) + radius

	pos := player.Object.Position()
	if pos.X != expectedX || pos.Y != expectedY {
		t.Errorf("Unexpected initial position. Got (%f,%f), want (%f,%f)",
			pos.X, pos.Y,
			expectedX, expectedY)
	}
}

func TestPlayerPosition(t *testing.T) {
	mock := NewMockHandler()
	speed := 5.0
	img := image.NewRGBA(image.Rect(0, 0, 32, 32))
	center := image.Pt(320, 240)

	player, err := NewPlayer(mock, speed, img, center)
	if err != nil {
		t.Fatalf("Failed to create player: %v", err)
	}

	expectedX := float64(center.X) + radius*(-1.0)
	expectedY := float64(center.Y) + radius

	// Get the position vector
	pos := player.Object.Position()

	// Access X and Y from the position vector
	if pos.X != expectedX || pos.Y != expectedY {
		t.Errorf("Expected position (%f, %f), got (%f, %f)",
			expectedX, expectedY, pos.X, pos.Y)
	}
}
