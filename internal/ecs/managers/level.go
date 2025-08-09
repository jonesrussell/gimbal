package managers

import (
	"github.com/jonesrussell/gimbal/internal/common"
)

// LevelManager manages the game level progression
type LevelManager struct {
	level  int
	logger common.Logger
}

// NewLevelManager creates a new level management system with the provided logger
func NewLevelManager(logger common.Logger) *LevelManager {
	return &LevelManager{
		level:  1,
		logger: logger,
	}
}

// GetLevel returns the current level
func (lm *LevelManager) GetLevel() int {
	return lm.level
}

// SetLevel sets the level to a specific value
func (lm *LevelManager) SetLevel(level int) {
	if level < 1 {
		level = 1
	}
	oldLevel := lm.level
	lm.level = level
	lm.logger.Debug("Level changed", "old_level", oldLevel, "new_level", lm.level)
}

// IncrementLevel increases the level by 1
func (lm *LevelManager) IncrementLevel() {
	lm.SetLevel(lm.level + 1)
}

// Reset resets the level to 1
func (lm *LevelManager) Reset() {
	oldLevel := lm.level
	lm.level = 1
	lm.logger.Debug("Level reset", "old_level", oldLevel, "new_level", lm.level)
}
