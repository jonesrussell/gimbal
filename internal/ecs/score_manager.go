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
func (scoreMgr *ScoreManager) GetScore() int {
	return scoreMgr.score
}

// AddScore adds points to the score
func (scoreMgr *ScoreManager) AddScore(points int) {
	oldScore := scoreMgr.score
	scoreMgr.score += points
	scoreMgr.eventSystem.EmitScoreChanged(oldScore, scoreMgr.score)
	scoreMgr.logger.Debug("Score updated", "old_score", oldScore, "new_score", scoreMgr.score, "points", points)
}

// SetScore sets the score to a specific value
func (scoreMgr *ScoreManager) SetScore(score int) {
	if score < 0 {
		score = 0
	}
	oldScore := scoreMgr.score
	scoreMgr.score = score
	scoreMgr.eventSystem.EmitScoreChanged(oldScore, scoreMgr.score)
	scoreMgr.logger.Debug("Score set", "old_score", oldScore, "new_score", scoreMgr.score)
}

// Reset resets the score to 0
func (scoreMgr *ScoreManager) Reset() {
	oldScore := scoreMgr.score
	scoreMgr.score = 0
	scoreMgr.eventSystem.EmitScoreChanged(oldScore, scoreMgr.score)
	scoreMgr.logger.Debug("Score reset", "old_score", oldScore, "new_score", scoreMgr.score)
}
