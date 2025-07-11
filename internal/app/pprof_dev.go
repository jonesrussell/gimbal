//go:build dev

package app

import (
	"log"
	"net/http"
	"time"

	_ "net/http/pprof" //nolint:gosec // Only included in dev builds
)

// StartPprofServer starts the pprof server for development builds only
func StartPprofServer() {
	// Create server with proper timeouts for security
	server := &http.Server{
		Addr:         "localhost:6060",
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Println("pprof server running at http://localhost:6060/debug/pprof/")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("pprof server error: %v", err)
		}
	}()

	// Store server reference for graceful shutdown if needed
	// This could be added to the container for proper cleanup
}
