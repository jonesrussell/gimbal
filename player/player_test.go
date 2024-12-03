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

	if player.Object.Position.X != expectedX || player.Object.Position.Y != expectedY {
		t.Errorf("Unexpected initial position. Got (%f,%f), want (%f,%f)",
			player.Object.Position.X, player.Object.Position.Y,
			expectedX, expectedY)
	}
}
