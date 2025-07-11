package scenes

import (
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// updateFadeIn handles the fade-in animation
func (s *PausedScene) updateFadeIn(dt float64) {
	if s.fadeIn < 1.0 {
		s.fadeIn = math.Min(1.0, s.fadeIn+dt/fadeInDuration)
	}
}

// updateSelectionAnimation manages selection change animations
func (s *PausedScene) updateSelectionAnimation() {
	if s.selectionChanged && time.Since(s.lastSelectionTime).Seconds() > selectionDelay {
		s.selectionChanged = false
	}
}

// handleInput processes pause-specific input (ESC key)
func (s *PausedScene) handleInput() {
	currentEscPressed := ebiten.IsKeyPressed(ebiten.KeyEscape)
	escJustPressed := inpututil.IsKeyJustPressed(ebiten.KeyEscape)

	// If ESC was pressed when we entered, wait for it to be released
	if s.escWasPressed && currentEscPressed {
		return // ESC is still held down from the pause action
	}

	// ESC has been released (or wasn't pressed when we entered)
	if s.escWasPressed && !currentEscPressed {
		s.escWasPressed = false
		s.canUnpause = true
		return // Don't process input this frame, just mark as ready
	}

	// Now we can check for new ESC presses
	if s.canUnpause && escJustPressed {
		// Call resume callback to unpause game state
		if s.manager.onResume != nil {
			s.manager.onResume()
		}

		s.manager.SwitchScene(ScenePlaying)
	}

	// If we entered without ESC pressed, we can unpause immediately
	if !s.escWasPressed {
		s.canUnpause = true
		if escJustPressed {
			// Call resume callback to unpause game state
			if s.manager.onResume != nil {
				s.manager.onResume()
			}

			s.manager.SwitchScene(ScenePlaying)
		}
	}
}
