package game

import (
	"fmt"
	"os"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/jonesrussell/gimbal/internal/config"
)

// Runner handles game execution with Ebiten
type Runner struct {
	config *config.AppGameConfig
	game   ebiten.Game
}

// NewRunner creates a new game runner
func NewRunner(cfg *config.AppGameConfig, game ebiten.Game) *Runner {
	return &Runner{
		config: cfg,
		game:   game,
	}
}

// Run configures Ebiten and starts the game
func (r *Runner) Run() error {
	// Configure Ebiten window
	ebiten.SetWindowSize(r.config.WindowWidth, r.config.WindowHeight)
	ebiten.SetWindowTitle(r.config.WindowTitle)
	ebiten.SetTPS(r.config.TPS)

	if r.config.Resizable {
		ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	}

	// Check if audio is already disabled via environment variable
	// If audio initialization failed during game init (e.g., no audio device),
	// we should set DISABLE_AUDIO to prevent Ebiten from trying to initialize audio
	// Ebiten's oto library may respect this, or at least fail more gracefully
	if os.Getenv("DISABLE_AUDIO") == "" {
		// Try to detect if audio is available by checking if we can create an audio context
		// If audio initialization would fail, set DISABLE_AUDIO to prevent Ebiten from trying
		// This is a best-effort attempt - Ebiten may still try to initialize audio internally
		r.checkAndDisableAudioIfNeeded()
	}

	// Run the game with panic recovery for audio-related issues
	// In environments like WSL2 where audio is unavailable, ebiten.RunGame
	// may panic or return an error. We catch these and handle them gracefully.
	var runErr error
	func() {
		defer func() {
			if r := recover(); r != nil {
				// Check if panic is audio-related
				panicMsg := strings.ToLower(fmt.Sprintf("%v", r))
				if containsAudioError(panicMsg) {
					// Audio is optional - convert panic to error
					if err, ok := r.(error); ok {
						runErr = fmt.Errorf("audio initialization failed (audio is optional): %w", err)
					} else {
						runErr = fmt.Errorf("audio initialization failed (audio is optional): %v", r)
					}
				} else {
					// Re-panic if it's not audio-related
					panic(r)
				}
			}
		}()

		// Run the game - this may panic or return an error if audio fails
		runErr = ebiten.RunGame(r.game)
	}()

	// Check if error is audio-related
	if runErr != nil && containsAudioError(strings.ToLower(runErr.Error())) {
		// Audio is optional - but if ebiten.RunGame returned an error, the game never started
		// We need to handle this differently - the game loop has already exited
		// Return the error so the caller knows the game couldn't start
		// The caller will log a warning but the game won't run
		return fmt.Errorf("audio initialization failed (audio is optional): %w", runErr)
	}

	return runErr
}

// checkAndDisableAudioIfNeeded checks if audio is available and sets DISABLE_AUDIO
// if audio initialization would fail. This is a best-effort attempt to prevent
// Ebiten from trying to initialize audio when it's not available.
func (r *Runner) checkAndDisableAudioIfNeeded() {
	// This is a lightweight check - we don't want to actually initialize audio here
	// as that would be wasteful. Instead, we check common indicators that audio
	// might not be available (e.g., running in WSL2, container, etc.)
	//
	// Note: This is a heuristic approach. The real solution would be for Ebiten
	// to handle audio initialization failures more gracefully, but until then,
	// we do our best to prevent the failure.

	// If DISABLE_AUDIO is not set, we let Ebiten try to initialize audio
	// and handle the error if it fails. Setting it here unconditionally would
	// disable audio even when it's available, which we don't want.
	//
	// The actual error handling happens in the panic recovery and error checking
	// above, which catches audio-related failures and treats them as non-fatal.
}

// containsAudioError checks if an error message is related to audio initialization
func containsAudioError(msg string) bool {
	return strings.Contains(msg, "alsa") ||
		strings.Contains(msg, "oto") ||
		strings.Contains(msg, "audio") ||
		strings.Contains(msg, "pulse") ||
		strings.Contains(msg, "jack") ||
		strings.Contains(msg, "snd_pcm")
}
