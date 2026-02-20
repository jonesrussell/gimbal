package scenes //nolint:testpackage // Testing from same package to access unexported fields

import (
	"testing"
)

func TestSceneManager_RequestQuit(t *testing.T) {
	// Create a minimal SceneManager for testing
	sm := &SceneManager{
		quitRequested: false,
	}

	// Initially should not be requested
	if sm.IsQuitRequested() {
		t.Error("Expected quit not to be requested initially")
	}

	// Request quit
	sm.RequestQuit()

	// Should now be requested
	if !sm.IsQuitRequested() {
		t.Error("Expected quit to be requested after RequestQuit()")
	}

	// Request quit again (idempotent)
	sm.RequestQuit()

	// Should still be requested
	if !sm.IsQuitRequested() {
		t.Error("Expected quit to remain requested after second RequestQuit()")
	}
}

func TestSceneManager_IsQuitRequested(t *testing.T) {
	tests := []struct {
		name          string
		quitRequested bool
		want          bool
	}{
		{
			name:          "quit not requested",
			quitRequested: false,
			want:          false,
		},
		{
			name:          "quit requested",
			quitRequested: true,
			want:          true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := &SceneManager{
				quitRequested: tt.quitRequested,
			}
			got := sm.IsQuitRequested()
			if got != tt.want {
				t.Errorf("IsQuitRequested() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSceneManager_RequestQuit_WithLogger(t *testing.T) {
	// Test RequestQuit sets the flag (logging is via log.Printf)
	sm := &SceneManager{
		quitRequested: false,
	}

	if sm.IsQuitRequested() {
		t.Error("Expected quit not to be requested initially")
	}

	sm.RequestQuit()

	if !sm.IsQuitRequested() {
		t.Error("Expected quit to be requested after RequestQuit()")
	}
}
