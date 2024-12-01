package player

import (
	"image"
	"testing"
)

func TestPlayer_Input(t *testing.T) {
	mock := NewMockHandler()
	speed := 5.0
	img := image.NewRGBA(image.Rect(0, 0, 32, 32))
	player, err := NewPlayer(mock, speed, img)
	if err != nil {
		t.Fatalf("Failed to create player: %v", err)
	}
	// Use player in tests...
}
