package ecs

import (
	"time"

	"github.com/jonesrussell/gimbal/internal/common"
)

// GameState represents the current state of the game
type GameState struct {
	IsPaused    bool
	IsGameOver  bool
	IsVictory   bool
	StartTime   time.Time
	LastUpdate  time.Time
	FrameCount  int64
	PlayerSpeed float64
}

// NewGameState creates a new game state
func NewGameState() *GameState {
	now := time.Now()
	return &GameState{
		IsPaused:   false,
		StartTime:  now,
		LastUpdate: now,
		FrameCount: 0,
		IsGameOver: false,
		IsVictory:  false,
	}
}

// GameStateManager manages game state and state transitions
type GameStateManager struct {
	state       *GameState
	eventSystem *EventSystem
	logger      common.Logger
}

// NewGameStateManager creates a new game state manager
func NewGameStateManager(eventSystem *EventSystem, logger common.Logger) *GameStateManager {
	return &GameStateManager{
		state:       NewGameState(),
		eventSystem: eventSystem,
		logger:      logger,
	}
}

// GetState returns the current game state
func (gsm *GameStateManager) GetState() *GameState {
	return gsm.state
}

// IsPaused returns whether the game is paused
func (gsm *GameStateManager) IsPaused() bool {
	return gsm.state.IsPaused
}

// TogglePause toggles the pause state
func (gsm *GameStateManager) TogglePause() {
	gsm.state.IsPaused = !gsm.state.IsPaused
	if gsm.state.IsPaused {
		gsm.eventSystem.EmitGamePaused()
		gsm.logger.Debug("Game paused")
	} else {
		gsm.eventSystem.EmitGameResumed()
		gsm.logger.Debug("Game resumed")
	}
}

// SetPaused sets the pause state explicitly
func (gsm *GameStateManager) SetPaused(paused bool) {
	if gsm.state.IsPaused != paused {
		gsm.state.IsPaused = paused
		if paused {
			gsm.eventSystem.EmitGamePaused()
			gsm.logger.Debug("Game paused")
		} else {
			gsm.eventSystem.EmitGameResumed()
			gsm.logger.Debug("Game resumed")
		}
	}
}

// GetGameTime returns the total time the game has been running
func (gsm *GameStateManager) GetGameTime() time.Duration {
	return time.Since(gsm.state.StartTime)
}

// GetFrameCount returns the total number of frames processed
func (gsm *GameStateManager) GetFrameCount() int64 {
	return gsm.state.FrameCount
}

// IncrementFrameCount increases the frame count
func (gsm *GameStateManager) IncrementFrameCount() {
	gsm.state.FrameCount++
}

// UpdateLastUpdateTime updates the last update timestamp
func (gsm *GameStateManager) UpdateLastUpdateTime() {
	gsm.state.LastUpdate = time.Now()
}

// IsGameOver returns whether the game is over
func (gsm *GameStateManager) IsGameOver() bool {
	return gsm.state.IsGameOver
}

// SetGameOver sets the game over state
func (gsm *GameStateManager) SetGameOver(gameOver bool) {
	gsm.state.IsGameOver = gameOver
	if gameOver {
		gsm.logger.Debug("Game over")
	}
}

// IsVictory returns whether the player has won
func (gsm *GameStateManager) IsVictory() bool {
	return gsm.state.IsVictory
}

// SetVictory sets the victory state
func (gsm *GameStateManager) SetVictory(victory bool) {
	gsm.state.IsVictory = victory
	if victory {
		gsm.logger.Debug("Victory achieved")
	}
}

// Reset resets the game state to initial values
func (gsm *GameStateManager) Reset() {
	gsm.state = NewGameState()
	gsm.logger.Debug("Game state reset")
}

// GetStateInfo returns a summary of the current game state
func (gsm *GameStateManager) GetStateInfo() map[string]interface{} {
	return map[string]interface{}{
		"is_paused":    gsm.state.IsPaused,
		"game_time":    gsm.GetGameTime().String(),
		"frame_count":  gsm.state.FrameCount,
		"is_game_over": gsm.state.IsGameOver,
		"is_victory":   gsm.state.IsVictory,
	}
}
