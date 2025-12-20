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

// InputState represents the current state of pause input
type InputState struct {
	CurrentEscPressed bool
	EscJustPressed    bool
	CanUnpause        bool
	EscWasPressed     bool
}

// detectInputState determines the current input state for pause handling
func (s *PausedScene) detectInputState() InputState {
	return InputState{
		CurrentEscPressed: ebiten.IsKeyPressed(ebiten.KeyEscape),
		EscJustPressed:    inpututil.IsKeyJustPressed(ebiten.KeyEscape),
		CanUnpause:        s.canUnpause,
		EscWasPressed:     s.escWasPressed,
	}
}

// selectCommand determines which command to execute based on input state
func (s *PausedScene) selectCommand(state InputState) PauseCommand {
	// Handle ESC key state transitions
	if state.EscWasPressed && state.CurrentEscPressed {
		return &escHeldCommand{} // ESC is still held down from the pause action
	}

	if state.EscWasPressed && !state.CurrentEscPressed {
		return &escReleasedCommand{} // ESC has been released
	}

	// Handle resume conditions
	if state.CanUnpause && state.EscJustPressed {
		return &resumeCommand{} // Resume from pause
	}

	// Handle initial unpause setup
	if !state.EscWasPressed {
		s.canUnpause = true
		if state.EscJustPressed {
			return &resumeCommand{}
		}
	}

	return nil // No command to execute
}

// handleInput processes pause-specific input (ESC key) using the Command pattern
func (s *PausedScene) handleInput() {
	state := s.detectInputState()
	cmd := s.selectCommand(state)

	if cmd != nil {
		cmd.Execute(s)
	}
}
