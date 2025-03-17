package test

import (
	"os"
	"os/exec"
	"testing"
)

// EnsureXvfb ensures that xvfb is available for testing
func EnsureXvfb(t *testing.T) {
	t.Helper()

	// Check if xvfb is already running
	if os.Getenv("DISPLAY") != "" {
		return
	}

	// Check if xvfb-run is available
	_, err := exec.LookPath("xvfb-run")
	if err != nil {
		t.Skip("xvfb-run not available, skipping test")
		return
	}

	// Check if xvfb is available
	_, err = exec.LookPath("Xvfb")
	if err != nil {
		t.Skip("Xvfb not available, skipping test")
		return
	}
}
