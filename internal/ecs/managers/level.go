package managers

import (
	"github.com/jonesrussell/gimbal/internal/common"
)

// LevelManager manages the game level progression
type LevelManager struct {
	level     int
	maxLevels int
	logger    common.Logger
}

// NewLevelManager creates a new level management system with the provided logger
func NewLevelManager(logger common.Logger) *LevelManager {
	return &LevelManager{
		level:     1,
		maxLevels: 6, // Default to 6 stages
		logger:    logger,
	}
}

// SetMaxLevels sets the maximum number of levels
func (lm *LevelManager) SetMaxLevels(maxLevels int) {
	if maxLevels < 1 {
		maxLevels = 1
	}
	lm.maxLevels = maxLevels
}

// GetLevel returns the current level number
func (lm *LevelManager) GetLevel() int {
	return lm.level
}

// SetLevel sets the level to a specific value
func (lm *LevelManager) SetLevel(level int) {
	if level < 1 {
		level = 1
	}
	if level > lm.maxLevels {
		level = lm.maxLevels
		lm.logger.Warn("Level exceeds available levels", "requested", level, "max", lm.maxLevels)
	}
	oldLevel := lm.level
	lm.level = level
	lm.logger.Debug("Level changed", "old_level", oldLevel, "new_level", lm.level)
}

// IncrementLevel increases the level by 1
func (lm *LevelManager) IncrementLevel() {
	if lm.level < lm.maxLevels {
		lm.SetLevel(lm.level + 1)
	} else {
		lm.logger.Debug("Already at max level", "level", lm.level)
	}
}

// Reset resets the level to 1
func (lm *LevelManager) Reset() {
	oldLevel := lm.level
	lm.level = 1
	lm.logger.Debug("Level reset", "old_level", oldLevel, "new_level", lm.level)
}

// HasMoreLevels returns true if there are more levels available
func (lm *LevelManager) HasMoreLevels() bool {
	return lm.level < lm.maxLevels
}

// GetLevelCount returns the total number of levels
func (lm *LevelManager) GetLevelCount() int {
	return lm.maxLevels
}
