//go:build !dev

package app

// StartPprofServer is a no-op in production builds
func StartPprofServer() {}
