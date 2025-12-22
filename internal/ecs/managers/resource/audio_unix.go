//go:build !windows && !js && !wasm

package resources

import "syscall"

// suppressStderrFD redirects stderr file descriptor to /dev/null and returns a restore function.
// This is necessary because C libraries (like ALSA) write directly to the file descriptor,
// bypassing Go's os.Stderr.
func suppressStderrFD(originalFd, devNullFd int) func() {
	// Save original stderr file descriptor by duplicating it
	savedFd, err := syscall.Dup(originalFd)
	if err != nil {
		// If dup fails, return no-op
		return func() {}
	}

	// Redirect stderr to /dev/null
	err = syscall.Dup2(devNullFd, originalFd)
	if err != nil {
		// If dup2 fails, close saved fd and return no-op
		syscall.Close(savedFd)
		return func() {}
	}

	// Return restore function
	return func() {
		// Restore original stderr
		//nolint:errcheck // Best effort restore - if it fails, there's nothing we can do
		_ = syscall.Dup2(savedFd, originalFd)
		//nolint:errcheck // Best effort cleanup - if it fails, there's nothing we can do
		_ = syscall.Close(savedFd)
	}
}
