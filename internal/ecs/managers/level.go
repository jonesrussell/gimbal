package managers

import (
	"fmt"

	"github.com/jonesrussell/gimbal/internal/common"
)

// LevelEventEmitter defines the interface for emitting level events.
type LevelEventEmitter interface {
	EmitLevelChanged(oldLevel, newLevel int)
}

// LevelManager manages the game level progression
type LevelManager struct {
	level        int
	levels       []LevelConfig
	logger       common.Logger
	eventEmitter LevelEventEmitter
}

// NewLevelManager creates a new level management system with the provided logger
func NewLevelManager(logger common.Logger) *LevelManager {
	lm := &LevelManager{
		level:  1,
		levels: []LevelConfig{},
		logger: logger,
	}
	return lm
}

// SetEventEmitter sets the event emitter for level change notifications.
func (lm *LevelManager) SetEventEmitter(emitter LevelEventEmitter) {
	lm.eventEmitter = emitter
}

// LoadLevels loads level configurations from the provided slice
func (lm *LevelManager) LoadLevels(levels []LevelConfig) error {
	if len(levels) == 0 {
		return fmt.Errorf("no levels provided")
	}
	lm.levels = levels
	lm.logger.Info("Levels loaded", "count", len(levels))
	return nil
}

// GetLevel returns the current level number
func (lm *LevelManager) GetLevel() int {
	return lm.level
}

// GetCurrentLevelConfig returns the configuration for the current level
func (lm *LevelManager) GetCurrentLevelConfig() *LevelConfig {
	if lm.level < 1 || lm.level > len(lm.levels) {
		// Return first level as fallback
		if len(lm.levels) > 0 {
			lm.logger.Warn("Level out of bounds, using first level", "requested", lm.level, "total", len(lm.levels))
			return &lm.levels[0]
		}
		return nil
	}
	return &lm.levels[lm.level-1] // Convert to 0-based index
}

// SetLevel sets the level to a specific value
func (lm *LevelManager) SetLevel(level int) {
	if level < 1 {
		level = 1
	}
	if level > len(lm.levels) {
		level = len(lm.levels)
		lm.logger.Warn("Level exceeds available levels", "requested", level, "max", len(lm.levels))
	}
	oldLevel := lm.level
	if oldLevel == level {
		return // No change
	}
	lm.level = level
	lm.logger.Debug("Level changed", "old_level", oldLevel, "new_level", lm.level)

	// Emit level changed event
	if lm.eventEmitter != nil {
		lm.eventEmitter.EmitLevelChanged(oldLevel, lm.level)
	}
}

// IncrementLevel increases the level by 1
func (lm *LevelManager) IncrementLevel() {
	if lm.level < len(lm.levels) {
		lm.SetLevel(lm.level + 1)
	} else {
		lm.logger.Debug("Already at max level", "level", lm.level)
	}
}

// Reset resets the level to 1
func (lm *LevelManager) Reset() {
	oldLevel := lm.level
	if oldLevel == 1 {
		return // Already at level 1
	}
	lm.level = 1
	lm.logger.Debug("Level reset", "old_level", oldLevel, "new_level", lm.level)

	// Emit level changed event
	if lm.eventEmitter != nil {
		lm.eventEmitter.EmitLevelChanged(oldLevel, lm.level)
	}
}

// HasMoreLevels returns true if there are more levels available
func (lm *LevelManager) HasMoreLevels() bool {
	return lm.level < len(lm.levels)
}

// GetLevelCount returns the total number of levels
func (lm *LevelManager) GetLevelCount() int {
	return len(lm.levels)
}
