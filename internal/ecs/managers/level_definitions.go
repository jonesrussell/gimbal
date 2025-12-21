package managers

// GetDefaultLevelDefinitions returns the default level configurations
// These can be overridden by JSON files if they exist
func GetDefaultLevelDefinitions() []LevelConfig {
	return []LevelConfig{
		createLevel1(),
		createLevel2(),
		createLevel3(),
		createLevel4(),
		createLevel5(),
	}
}

// Constants for enemy types (matching enemy.EnemyType)
const (
	EnemyTypeBasic = 0
	EnemyTypeHeavy = 1
	EnemyTypeBoss  = 2
)

// Constants for formation types (matching enemy.FormationType)
const (
	FormationLine     = 0
	FormationCircle   = 1
	FormationV        = 2
	FormationDiamond  = 3
	FormationDiagonal = 4
	FormationSpiral   = 5
	FormationRandom   = 6
)

// Constants for movement patterns (matching enemy.MovementPattern)
const (
	MovementPatternNormal       = 0
	MovementPatternZigzag       = 1
	MovementPatternAccelerating = 2
	MovementPatternPulsing      = 3
)

// createLevel1 creates the first level (tutorial/easy)
func createLevel1() LevelConfig {
	return LevelConfig{
		LevelNumber: 1,
		Metadata: LevelMetadata{
			Name:        "First Contact",
			Description: "Welcome to the fight. Learn the basics.",
			MusicTrack:  "",
			Background:  "default",
		},
		Waves: []WaveConfig{
			createWave(FormationLine, 6, []int{EnemyTypeBasic},
				waveParams{0.25, 0.0, MovementPatternNormal}),
			createWave(FormationCircle, 10, []int{EnemyTypeBasic},
				waveParams{0.1, 1.5, MovementPatternNormal}),
			createWave(FormationV, 9, []int{EnemyTypeBasic, EnemyTypeHeavy},
				waveParams{0.18, 1.5, MovementPatternZigzag}),
		},
		Boss: BossConfig{
			Enabled:         true,
			EnemyType:       EnemyTypeBoss,
			Health:          8,
			SpawnDelay:      2.0,
			MovementType:    "orbital",
			CanShoot:        true,
			FireRate:        1.5,
			ProjectileSpeed: 5.0,
			Size:            64,
			Points:          1000,
		},
		Difficulty:           DefaultDifficultySettings(),
		CompletionConditions: DefaultCompletionConditions(),
	}
}

// createLevel2 creates the second level (moderate difficulty)
func createLevel2() LevelConfig {
	return LevelConfig{
		LevelNumber: 2,
		Metadata: LevelMetadata{
			Name:        "Escalation",
			Description: "The enemy adapts. More formations, more danger.",
			MusicTrack:  "",
			Background:  "default",
		},
		Waves: []WaveConfig{
			createWave(FormationDiamond, 8, []int{EnemyTypeHeavy},
				waveParams{0.2, 1.5, MovementPatternAccelerating}),
			createWave(FormationSpiral, 12,
				[]int{EnemyTypeBasic, EnemyTypeHeavy},
				waveParams{0.12, 1.5, MovementPatternPulsing}),
			createWave(FormationDiagonal, 10,
				[]int{EnemyTypeHeavy, EnemyTypeBasic},
				waveParams{0.15, 1.5, MovementPatternNormal}),
			createWave(FormationRandom, 14,
				[]int{EnemyTypeBasic, EnemyTypeHeavy},
				waveParams{0.1, 1.5, MovementPatternZigzag}),
		},
		Boss: BossConfig{
			Enabled:         true,
			EnemyType:       EnemyTypeBoss,
			Health:          12,
			SpawnDelay:      2.0,
			MovementType:    "orbital",
			CanShoot:        true,
			FireRate:        2.0,
			ProjectileSpeed: 6.0,
			Size:            64,
			Points:          1500,
		},
		Difficulty: DifficultySettings{
			EnemySpeedMultiplier:     1.1,
			EnemyHealthMultiplier:    1.0,
			EnemySpawnRateMultiplier: 1.0,
			PlayerDamageMultiplier:   1.0,
			ScoreMultiplier:          1.2,
		},
		CompletionConditions: DefaultCompletionConditions(),
	}
}

// createLevel3 creates the third level (increased difficulty)
func createLevel3() LevelConfig {
	return LevelConfig{
		LevelNumber: 3,
		Metadata: LevelMetadata{
			Name:        "Pressure Point",
			Description: "Enemies come faster and hit harder.",
			MusicTrack:  "",
			Background:  "default",
		},
		Waves: []WaveConfig{
			createWave(FormationCircle, 12,
				[]int{EnemyTypeHeavy, EnemyTypeBasic},
				waveParams{0.12, 1.5, MovementPatternAccelerating}),
			createWave(FormationSpiral, 15,
				[]int{EnemyTypeBasic, EnemyTypeHeavy},
				waveParams{0.1, 1.5, MovementPatternPulsing}),
			createWave(FormationDiamond, 10, []int{EnemyTypeHeavy},
				waveParams{0.15, 1.5, MovementPatternZigzag}),
			createWave(FormationRandom, 16,
				[]int{EnemyTypeBasic, EnemyTypeHeavy},
				waveParams{0.08, 1.5, MovementPatternAccelerating}),
		},
		Boss: BossConfig{
			Enabled:         true,
			EnemyType:       EnemyTypeBoss,
			Health:          15,
			SpawnDelay:      2.0,
			MovementType:    "orbital",
			CanShoot:        true,
			FireRate:        2.5,
			ProjectileSpeed: 6.5,
			Size:            64,
			Points:          2000,
		},
		Difficulty: DifficultySettings{
			EnemySpeedMultiplier:     1.2,
			EnemyHealthMultiplier:    1.1,
			EnemySpawnRateMultiplier: 1.1,
			PlayerDamageMultiplier:   1.0,
			ScoreMultiplier:          1.5,
		},
		CompletionConditions: DefaultCompletionConditions(),
	}
}

// createLevel4 creates the fourth level (high difficulty)
func createLevel4() LevelConfig {
	return LevelConfig{
		LevelNumber: 4,
		Metadata: LevelMetadata{
			Name:        "The Gauntlet",
			Description: "Survive the onslaught. Every wave counts.",
			MusicTrack:  "",
			Background:  "default",
		},
		Waves: []WaveConfig{
			createWave(FormationV, 12,
				[]int{EnemyTypeBasic, EnemyTypeHeavy},
				waveParams{0.1, 1.0, MovementPatternZigzag}),
			createWave(FormationSpiral, 18,
				[]int{EnemyTypeBasic, EnemyTypeHeavy},
				waveParams{0.08, 1.0, MovementPatternPulsing}),
			createWave(FormationDiamond, 12, []int{EnemyTypeHeavy},
				waveParams{0.12, 1.0, MovementPatternAccelerating}),
			createWave(FormationRandom, 20,
				[]int{EnemyTypeBasic, EnemyTypeHeavy},
				waveParams{0.06, 1.0, MovementPatternZigzag}),
			createWave(FormationCircle, 15,
				[]int{EnemyTypeHeavy, EnemyTypeBasic},
				waveParams{0.1, 1.0, MovementPatternPulsing}),
		},
		Boss: BossConfig{
			Enabled:         true,
			EnemyType:       EnemyTypeBoss,
			Health:          20,
			SpawnDelay:      1.5,
			MovementType:    "orbital",
			CanShoot:        true,
			FireRate:        3.0,
			ProjectileSpeed: 7.0,
			Size:            64,
			Points:          2500,
		},
		Difficulty: DifficultySettings{
			EnemySpeedMultiplier:     1.3,
			EnemyHealthMultiplier:    1.2,
			EnemySpawnRateMultiplier: 1.2,
			PlayerDamageMultiplier:   1.0,
			ScoreMultiplier:          2.0,
		},
		CompletionConditions: DefaultCompletionConditions(),
	}
}

// createLevel5 creates the fifth level (very high difficulty)
func createLevel5() LevelConfig {
	return LevelConfig{
		LevelNumber: 5,
		Metadata: LevelMetadata{
			Name:        "Final Stand",
			Description: "The ultimate test. Prove your worth.",
			MusicTrack:  "",
			Background:  "default",
		},
		Waves: []WaveConfig{
			createWave(FormationSpiral, 20,
				[]int{EnemyTypeBasic, EnemyTypeHeavy},
				waveParams{0.08, 0.8, MovementPatternPulsing}),
			createWave(FormationDiamond, 15, []int{EnemyTypeHeavy},
				waveParams{0.1, 0.8, MovementPatternAccelerating}),
			createWave(FormationRandom, 25,
				[]int{EnemyTypeBasic, EnemyTypeHeavy},
				waveParams{0.05, 0.8, MovementPatternZigzag}),
			createWave(FormationCircle, 18,
				[]int{EnemyTypeHeavy, EnemyTypeBasic},
				waveParams{0.08, 0.8, MovementPatternPulsing}),
			createWave(FormationV, 16,
				[]int{EnemyTypeBasic, EnemyTypeHeavy},
				waveParams{0.1, 0.8, MovementPatternAccelerating}),
		},
		Boss: BossConfig{
			Enabled:         true,
			EnemyType:       EnemyTypeBoss,
			Health:          25,
			SpawnDelay:      1.0,
			MovementType:    "orbital",
			CanShoot:        true,
			FireRate:        3.5,
			ProjectileSpeed: 8.0,
			Size:            64,
			Points:          5000,
		},
		Difficulty: DifficultySettings{
			EnemySpeedMultiplier:     1.5,
			EnemyHealthMultiplier:    1.3,
			EnemySpawnRateMultiplier: 1.3,
			PlayerDamageMultiplier:   1.0,
			ScoreMultiplier:          3.0,
		},
		CompletionConditions: DefaultCompletionConditions(),
	}
}

// waveParams holds parameters for creating a wave (helper type)
type waveParams struct {
	spawnDelay     float64
	interWaveDelay float64
	pattern        int // MovementPattern as int
}

// createWave creates a wave configuration with the given parameters
func createWave(
	formation int,
	count int,
	types []int,
	params waveParams,
) WaveConfig {
	return WaveConfig{
		FormationType:   formation,
		EnemyCount:      count,
		EnemyTypes:      types,
		SpawnDelay:      params.spawnDelay,
		Timeout:         12.0,
		InterWaveDelay:  params.interWaveDelay,
		MovementPattern: params.pattern,
	}
}
