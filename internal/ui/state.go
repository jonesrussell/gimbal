package ui

import (
	"github.com/jonesrussell/gimbal/internal/ui/core"
)

// State represents the current UI state
type State struct {
	Lives  int
	Score  int
	Health float64
	Ammo   int
}

// NewState creates a new UI state with default values
func NewState() *State {
	return &State{
		Lives:  3,
		Score:  0,
		Health: 1.0,
		Ammo:   10,
	}
}

// Validate ensures the state values are valid
func (s *State) Validate() error {
	if s.Lives < 0 {
		return core.ErrInvalidLives
	}
	if s.Score < 0 {
		return core.ErrInvalidScore
	}
	if s.Health < 0 || s.Health > 1 {
		return core.ErrInvalidHealth
	}
	if s.Ammo < 0 {
		return core.ErrInvalidAmmo
	}
	return nil
}
