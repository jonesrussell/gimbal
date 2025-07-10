package ecs

import (
	"github.com/jonesrussell/gimbal/internal/common"
)

// ScoreManager manages the game score
type ScoreManager struct {
	score       int
	eventSystem *EventSystem
	logger      common.Logger
}

// NewScoreManager creates a new score manager
func NewScoreManager(eventSystem *EventSystem, logger common.Logger) *ScoreManager {
	return &ScoreManager{
		score:       0,
		eventSystem: eventSystem,
		logger:      logger,
	}
}

// GetScore returns the current score
func (sm *ScoreManager) GetScore() int {
	return sm.score
}

// AddScore adds points to the score
func (sm *ScoreManager) AddScore(points int) {
	oldScore := sm.score
	sm.score += points
	sm.eventSystem.EmitScoreChanged(oldScore, sm.score)
	sm.logger.Debug("Score updated", "old_score", oldScore, "new_score", sm.score, "points", points)
}

// SetScore sets the score to a specific value
func (sm *ScoreManager) SetScore(score int) {
	if score < 0 {
		score = 0
	}
	oldScore := sm.score
	sm.score = score
	sm.eventSystem.EmitScoreChanged(oldScore, sm.score)
	sm.logger.Debug("Score set", "old_score", oldScore, "new_score", sm.score)
}

// Reset resets the score to 0
func (sm *ScoreManager) Reset() {
	oldScore := sm.score
	sm.score = 0
	sm.eventSystem.EmitScoreChanged(oldScore, sm.score)
	sm.logger.Debug("Score reset", "old_score", oldScore, "new_score", sm.score)
}
