package pause

import (
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/jonesrussell/gimbal/internal/scenes"
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

// PauseCommand defines the interface for pause menu input commands
// (Single method, single responsibility)
type PauseCommand interface {
	Execute(s *PausedScene)
}

// escHeldCommand does nothing while ESC is held down
type escHeldCommand struct{}

func (c *escHeldCommand) Execute(s *PausedScene) {}

// escReleasedCommand marks the scene as ready to unpause
type escReleasedCommand struct{}

func (c *escReleasedCommand) Execute(s *PausedScene) {
	s.escWasPressed = false
	s.canUnpause = true
}

// resumeCommand resumes the game from pause
type resumeCommand struct{}

func (c *resumeCommand) Execute(s *PausedScene) {
	// Note: Resume callback is handled in the menu action, not here
	s.manager.SwitchScene(scenes.ScenePlaying)
}

// handleInput processes pause-specific input (ESC key) using the Command pattern
func (s *PausedScene) handleInput() {
	currentEscPressed := ebiten.IsKeyPressed(ebiten.KeyEscape)
	escJustPressed := inpututil.IsKeyJustPressed(ebiten.KeyEscape)

	var cmd PauseCommand

	if s.escWasPressed && currentEscPressed {
		cmd = &escHeldCommand{} // ESC is still held down from the pause action
	} else if s.escWasPressed && !currentEscPressed {
		cmd = &escReleasedCommand{} // ESC has been released
	} else if s.canUnpause && escJustPressed {
		cmd = &resumeCommand{} // Resume from pause
	} else if !s.escWasPressed {
		s.canUnpause = true
		if escJustPressed {
			cmd = &resumeCommand{}
		}
	}

	if cmd != nil {
		cmd.Execute(s)
	}
}
