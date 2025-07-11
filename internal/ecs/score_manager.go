package ecs

import (
	"github.com/jonesrussell/gimbal/internal/common"
)

// ScoreManager manages the game score, high scores, and multipliers
type ScoreManager struct {
	score          int
	highScore      int
	multiplier     int
	bonusLifeScore int
	eventSystem *EventSystem
	logger      common.Logger
}

// NewScoreManager creates a new score manager with default settings
func NewScoreManager(eventSystem *EventSystem, logger common.Logger) *ScoreManager {
	return &ScoreManager{
		score:          0,
		highScore:      0,
		multiplier:     1,
		bonusLifeScore: 10000, // Bonus life every 10,000 points
		eventSystem:    eventSystem,
		logger:         logger,
	}
}

// GetScore returns the current score
func (scoreMgr *ScoreManager) GetScore() int {
	return scoreMgr.score
}

// GetHighScore returns the highest score achieved
func (scoreMgr *ScoreManager) GetHighScore() int {
	return scoreMgr.highScore
}

// GetMultiplier returns the current score multiplier
func (scoreMgr *ScoreManager) GetMultiplier() int {
	return scoreMgr.multiplier
}

// AddScore adds points to the score with multiplier
func (scoreMgr *ScoreManager) AddScore(points int) {
	if points <= 0 {
		return
	}

	// Calculate score with multiplier
	totalPoints := points * scoreMgr.multiplier
	oldScore := scoreMgr.score
	newScore := oldScore + totalPoints

	// Update high score if needed
	if newScore > scoreMgr.highScore {
		scoreMgr.highScore = newScore
	}

	// Check for bonus life
	scoreMgr.checkBonusLife(oldScore, newScore)

	scoreMgr.score = newScore
	scoreMgr.eventSystem.EmitScoreChanged(oldScore, scoreMgr.score)
	scoreMgr.logger.Debug("Score updated", 
		"old_score", oldScore, 
		"new_score", scoreMgr.score, 
		"points", points,
		"multiplier", scoreMgr.multiplier)
}

// SetScore sets the score to a specific value
func (scoreMgr *ScoreManager) SetScore(score int) {
	if score < 0 {
		score = 0
	}
	oldScore := scoreMgr.score
	scoreMgr.score = score
	
	// Update high score if needed
	if score > scoreMgr.highScore {
		scoreMgr.highScore = score
	}

	scoreMgr.eventSystem.EmitScoreChanged(oldScore, scoreMgr.score)
	scoreMgr.logger.Debug("Score set", "old_score", oldScore, "new_score", scoreMgr.score)
}

// Reset resets the score to 0 but keeps high score
func (scoreMgr *ScoreManager) Reset() {
	oldScore := scoreMgr.score
	scoreMgr.score = 0
	scoreMgr.multiplier = 1
	scoreMgr.eventSystem.EmitScoreChanged(oldScore, scoreMgr.score)
	scoreMgr.logger.Debug("Score reset", "old_score", oldScore, "new_score", scoreMgr.score)
}

// SetMultiplier sets the score multiplier
func (scoreMgr *ScoreManager) SetMultiplier(multiplier int) {
	if multiplier < 1 {
		multiplier = 1
	}
	if multiplier > 10 {
		multiplier = 10 // Cap multiplier at 10x
	}
	
	scoreMgr.multiplier = multiplier
	scoreMgr.logger.Debug("Multiplier updated", "new_multiplier", multiplier)
}

// checkBonusLife checks if the score has reached a bonus life threshold
func (scoreMgr *ScoreManager) checkBonusLife(oldScore, newScore int) {
	if oldScore/scoreMgr.bonusLifeScore < newScore/scoreMgr.bonusLifeScore {
		// Player has earned a bonus life
		scoreMgr.eventSystem.EmitBonusLife()
		scoreMgr.logger.Info("Bonus life earned", "score", newScore)
	}
}
