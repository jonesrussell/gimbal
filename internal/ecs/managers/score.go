package managers

// ScoreManager manages the game score, high scores, and multipliers
type ScoreManager struct {
	score            int
	highScore        int
	multiplier       int
	bonusLifeScore   int
	bonusLifeAwarded bool
}

// NewScoreManager creates a new score manager with default settings
func NewScoreManager(bonusLifeScore int) *ScoreManager {
	return &ScoreManager{
		score:            0,
		highScore:        0,
		multiplier:       1,
		bonusLifeScore:   bonusLifeScore,
		bonusLifeAwarded: false,
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

// GetBonusLifeScore returns the score threshold for bonus lives
func (scoreMgr *ScoreManager) GetBonusLifeScore() int {
	return scoreMgr.bonusLifeScore
}

// AddScore adds points to the score with multiplier
func (scoreMgr *ScoreManager) AddScore(points int) {
	if points <= 0 {
		return
	}

	// Calculate score with multiplier
	totalPoints := points * scoreMgr.multiplier
	scoreMgr.score += totalPoints

	// Update high score if needed
	if scoreMgr.score > scoreMgr.highScore {
		scoreMgr.highScore = scoreMgr.score
	}
}

// SetScore sets the score to a specific value
func (scoreMgr *ScoreManager) SetScore(score int) {
	if score < 0 {
		score = 0
	}
	scoreMgr.score = score

	// Update high score if needed
	if score > scoreMgr.highScore {
		scoreMgr.highScore = score
	}
}

// Reset resets the score to 0 but keeps high score
func (scoreMgr *ScoreManager) Reset() {
	scoreMgr.score = 0
	scoreMgr.multiplier = 1
	scoreMgr.bonusLifeAwarded = false
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
}

// ShouldAwardBonusLife returns true if a bonus life should be awarded
func (scoreMgr *ScoreManager) ShouldAwardBonusLife() bool {
	return scoreMgr.score >= scoreMgr.bonusLifeScore && !scoreMgr.bonusLifeAwarded
}

// MarkBonusLifeAwarded marks that a bonus life has been awarded
func (scoreMgr *ScoreManager) MarkBonusLifeAwarded() {
	scoreMgr.bonusLifeAwarded = true
}

// GetBonusLifeCount returns the number of bonus lives earned
func (scoreMgr *ScoreManager) GetBonusLifeCount() int {
	return scoreMgr.score / scoreMgr.bonusLifeScore
}
