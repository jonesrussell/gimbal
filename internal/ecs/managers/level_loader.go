package managers

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jonesrussell/gimbal/internal/common"
)

// LoadLevelsFromJSON loads level configurations from JSON files in the specified directory
// Returns the loaded levels or an error if loading fails
// Falls back to default levels if directory doesn't exist or files can't be loaded
func LoadLevelsFromJSON(dirPath string, logger common.Logger) ([]LevelConfig, error) {
	if logger == nil {
		// Create a no-op logger if none provided
		logger = &noOpLogger{}
	}

	// Check if directory exists
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		logger.Debug("Level directory does not exist, using default levels", "path", dirPath)
		return GetDefaultLevelDefinitions(), nil
	}

	// Read all JSON files in the directory
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		logger.Warn("Failed to read level directory, using default levels", "error", err, "path", dirPath)
		return GetDefaultLevelDefinitions(), nil
	}

	// Pre-allocate slice with estimated capacity
	levels := processLevelFiles(entries, dirPath, logger)

	// If no levels were loaded, fall back to defaults
	if len(levels) == 0 {
		logger.Debug("No valid level files found, using default levels")
		return GetDefaultLevelDefinitions(), nil
	}

	// Sort levels by level number (in case files are loaded out of order)
	sortLevelsByNumber(levels)

	logger.Info("Levels loaded from JSON", "count", len(levels), "path", dirPath)
	return levels, nil
}

// loadLevelFromFile loads a single level configuration from a JSON file
func loadLevelFromFile(filePath string, logger common.Logger) (LevelConfig, error) {
	var level LevelConfig

	data, err := os.ReadFile(filePath)
	if err != nil {
		return level, fmt.Errorf("failed to read file: %w", err)
	}

	if unmarshalErr := json.Unmarshal(data, &level); unmarshalErr != nil {
		return level, fmt.Errorf("failed to parse JSON: %w", unmarshalErr)
	}

	// Validate level configuration
	if level.LevelNumber < 1 {
		return level, fmt.Errorf("invalid level number: %d", level.LevelNumber)
	}

	if len(level.Waves) == 0 {
		return level, fmt.Errorf("level %d has no waves", level.LevelNumber)
	}

	// Set defaults for missing fields
	if level.Difficulty.EnemySpeedMultiplier == 0 {
		level.Difficulty = DefaultDifficultySettings()
	}

	// Validate enum values in waves
	for i := range level.Waves {
		if err = validateWaveConfig(&level.Waves[i], i); err != nil {
			return level, fmt.Errorf("wave %d: %w", i, err)
		}
	}

	// Validate boss config if enabled
	if level.Boss.Enabled {
		if err = validateBossConfig(&level.Boss); err != nil {
			return level, fmt.Errorf("boss config: %w", err)
		}
	}

	return level, nil
}

// validateWaveConfig validates enum values in a wave config
func validateWaveConfig(wave *WaveConfig, index int) error {
	// Validate FormationType (0-6)
	if wave.FormationType < 0 || wave.FormationType > 6 {
		return fmt.Errorf("invalid formation_type: %d (must be 0-6)", wave.FormationType)
	}

	// Validate EnemyTypes
	for j, et := range wave.EnemyTypes {
		if et < 0 || et > 2 {
			return fmt.Errorf("invalid enemy_type at index %d: %d (must be 0-2)", j, et)
		}
	}

	// Validate MovementPattern (0-3)
	if wave.MovementPattern < 0 || wave.MovementPattern > 3 {
		return fmt.Errorf("invalid movement_pattern: %d (must be 0-3)", wave.MovementPattern)
	}

	return nil
}

// validateBossConfig validates enum values in boss config
func validateBossConfig(boss *BossConfig) error {
	// Validate EnemyType (should be 2 for Boss)
	if boss.EnemyType < 0 || boss.EnemyType > 2 {
		return fmt.Errorf("invalid enemy_type: %d (must be 0-2)", boss.EnemyType)
	}
	return nil
}

// Note: Enum conversion functions are in game/game_level.go (convertWaveConfigs)
// to avoid import cycles between managers and enemy packages

// processLevelFiles processes directory entries and loads levels
func processLevelFiles(entries []os.DirEntry, dirPath string, logger common.Logger) []LevelConfig {
	levels := make([]LevelConfig, 0, len(entries))

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// Only process .json files
		if filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		filePath := filepath.Join(dirPath, entry.Name())
		level, loadErr := loadLevelFromFile(filePath, logger)
		if loadErr != nil {
			logger.Warn("Failed to load level file, skipping", "file", entry.Name(), "error", loadErr)
			continue
		}

		levels = append(levels, level)
		logger.Debug("Loaded level from JSON", "file", entry.Name(), "level", level.LevelNumber)
	}

	return levels
}

// sortLevelsByNumber sorts levels by level number using bubble sort
func sortLevelsByNumber(levels []LevelConfig) {
	for i := 0; i < len(levels)-1; i++ {
		for j := 0; j < len(levels)-i-1; j++ {
			if levels[j].LevelNumber > levels[j+1].LevelNumber {
				levels[j], levels[j+1] = levels[j+1], levels[j]
			}
		}
	}
}

// noOpLogger is a minimal logger implementation for when no logger is provided
type noOpLogger struct{}

func (n *noOpLogger) Debug(msg string, fields ...interface{})                             {}
func (n *noOpLogger) DebugContext(ctx context.Context, msg string, fields ...interface{}) {}
func (n *noOpLogger) Info(msg string, fields ...interface{})                              {}
func (n *noOpLogger) InfoContext(ctx context.Context, msg string, fields ...interface{})  {}
func (n *noOpLogger) Warn(msg string, fields ...interface{})                              {}
func (n *noOpLogger) WarnContext(ctx context.Context, msg string, fields ...interface{})  {}
func (n *noOpLogger) Error(msg string, fields ...interface{})                             {}
func (n *noOpLogger) ErrorContext(ctx context.Context, msg string, fields ...interface{}) {}
func (n *noOpLogger) Sync() error                                                         { return nil }
