package managers

import (
	"embed"
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/jonesrussell/gimbal/internal/common"
)

// StageLoader loads stage configurations from JSON files
type StageLoader struct {
	logger    common.Logger
	assetsFS  embed.FS
	stagePath string
}

// NewStageLoader creates a new stage loader
func NewStageLoader(logger common.Logger, assetsFS embed.FS) *StageLoader {
	return &StageLoader{
		logger:    logger,
		assetsFS:  assetsFS,
		stagePath: "stages", // Embedded paths are relative to the embed directive, not including "assets/"
	}
}

// LoadStage loads a stage configuration by number
func (sl *StageLoader) LoadStage(stageNumber int) (*StageConfig, error) {
	filename := fmt.Sprintf("stage_%02d.json", stageNumber)
	fullPath := filepath.Join(sl.stagePath, filename)

	data, err := sl.assetsFS.ReadFile(fullPath)
	if err != nil {
		sl.logger.Warn("Stage file not found, using default",
			"stage", stageNumber,
			"path", fullPath,
			"error", err)
		return sl.createDefaultStage(stageNumber), nil
	}

	var config StageConfig
	if unmarshalErr := json.Unmarshal(data, &config); unmarshalErr != nil {
		sl.logger.Error("Failed to parse stage config",
			"stage", stageNumber,
			"error", unmarshalErr)
		return nil, fmt.Errorf("failed to parse stage %d: %w", stageNumber, unmarshalErr)
	}

	sl.logger.Debug("Stage loaded",
		"stage", stageNumber,
		"name", config.Metadata.Name,
		"waves", len(config.Waves))

	return &config, nil
}

// LoadStageByName loads a stage configuration by planet name
func (sl *StageLoader) LoadStageByName(planetName string) (*StageConfig, error) {
	// Map planet names to stage numbers
	planetToStage := map[string]int{
		"earth":   1,
		"mars":    2,
		"jupiter": 3,
		"saturn":  4,
		"uranus":  5,
		"neptune": 6,
		"pluto":   7,
	}

	stageNum, exists := planetToStage[planetName]
	if !exists {
		return nil, fmt.Errorf("unknown planet: %s", planetName)
	}

	return sl.LoadStage(stageNum)
}

// createDefaultStage creates a default stage configuration
func (sl *StageLoader) createDefaultStage(stageNumber int) *StageConfig {
	planet := sl.getPlanetName(stageNumber)

	return &StageConfig{
		StageNumber: stageNumber,
		Planet:      planet,
		Metadata:    sl.createDefaultMetadata(stageNumber, planet),
		Waves:       []GyrussWave{sl.createDefaultWave()},
		Boss:        sl.createDefaultBoss(),
		PowerUps:    sl.createDefaultPowerUps(),
		Difficulty:  DefaultDifficultySettings(),
	}
}

func (sl *StageLoader) getPlanetName(stageNumber int) string {
	planets := []string{"Earth", "Mars", "Jupiter", "Saturn", "Uranus", "Neptune", "Pluto"}
	if stageNumber > 0 && stageNumber <= len(planets) {
		return planets[stageNumber-1]
	}
	return "Unknown"
}

func (sl *StageLoader) createDefaultMetadata(stageNumber int, planet string) StageMetadata {
	return StageMetadata{
		Name:        fmt.Sprintf("Stage %d - %s", stageNumber, planet),
		Description: fmt.Sprintf("Journey to %s", planet),
		MusicTrack:  "",
		Background:  "default",
	}
}

func (sl *StageLoader) createDefaultWave() GyrussWave {
	return GyrussWave{
		WaveID:        "default_wave_1",
		Description:   "Opening wave",
		SpawnSequence: []EnemyGroupConfig{sl.createDefaultEnemyGroup()},
		Timing:        WaveTiming{InterWaveDelay: 2.0, Timeout: 30.0},
	}
}

func (sl *StageLoader) createDefaultEnemyGroup() EnemyGroupConfig {
	return EnemyGroupConfig{
		EnemyType:     "basic",
		Count:         8,
		SpawnDelay:    0.0,
		SpawnInterval: 0.3,
		EntryPath: EntryPathConfig{
			Type: "spiral_in", Duration: 2.0,
			Parameters: EntryPathParams{SpiralTurns: 1.5, RotationDirection: "clockwise", StartRadius: 20},
		},
		ScaleAnimation: ScaleAnimConfig{StartScale: 0.1, EndScale: 1.0, Easing: "ease_out"},
		Behavior: BehaviorConfig{
			PostEntry: "orbit_then_attack", OrbitDuration: 3.0,
			OrbitDirection: "clockwise", OrbitSpeed: 45.0, MaxAttacks: 2,
		},
		AttackPattern: AttackConfig{Type: "single_rush", Cooldown: 5.0, RushSpeed: 300.0},
		FirePattern:   FireConfig{Type: "single_shot", FireRate: 0.5, FireWhileOrbit: true},
		Retreat:       RetreatConfig{HealthThreshold: 0.2, Timeout: 15.0, Speed: 200.0},
	}
}

func (sl *StageLoader) createDefaultBoss() StageBossConfig {
	return StageBossConfig{
		Enabled:   true,
		BossType:  "standard_boss",
		Health:    10,
		Size:      64,
		EntryPath: EntryPathConfig{Type: "straight_in", Duration: 3.0},
		Behavior: BehaviorConfig{
			PostEntry: "hover_center_then_orbit", OrbitDuration: 2.0,
			OrbitDirection: "clockwise", OrbitSpeed: 30.0, MaxAttacks: 0,
		},
		AttackPattern: AttackConfig{Type: "paired_rush", Cooldown: 4.0, RushSpeed: 250.0},
		FirePattern: FireConfig{
			Type: "spray", FireRate: 2.0, SprayAngle: 60.0, ProjectileCount: 5,
			FireWhileOrbit: true, FireWhileAttack: true,
		},
		SpawnDelay: 3.0,
		Points:     1000,
	}
}

func (sl *StageLoader) createDefaultPowerUps() PowerUpConfig {
	return PowerUpConfig{
		DropChance: 0.15,
		Types: []PowerUpTypeConfig{
			{Type: "double_shot", Weight: 0.7, Duration: 10.0},
			{Type: "extra_life", Weight: 0.3},
		},
	}
}

// GetTotalStages returns the total number of stages
func (sl *StageLoader) GetTotalStages() int {
	return 6 // Earth through Neptune/Pluto
}
