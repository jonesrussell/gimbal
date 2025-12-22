//go:build windows || js || wasm

package resources

// suppressStderrFD is a no-op on non-Unix platforms.
// File descriptor manipulation is not available or needed on these platforms.
func suppressStderrFD(originalFd, devNullFd int) func() {
	return func() {}
}
