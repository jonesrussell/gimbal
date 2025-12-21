package managers

// LevelConfig defines the complete configuration for a game level
// Note: WaveConfig types are defined as compatible structs here to avoid import cycles
// They will be converted to enemy.WaveConfig when loading into the enemy system
type LevelConfig struct {
	LevelNumber          int                  `json:"level_number"`
	Metadata             LevelMetadata        `json:"metadata"`
	Waves                []WaveConfig         `json:"waves"`
	Boss                 BossConfig           `json:"boss"`
	Difficulty           DifficultySettings   `json:"difficulty"`
	CompletionConditions CompletionConditions `json:"completion_conditions"`
}

// WaveConfig defines the configuration for a wave (compatible with enemy.WaveConfig)
// This is a duplicate definition to avoid import cycles
type WaveConfig struct {
	FormationType   int     `json:"formation_type"` // FormationType as int
	EnemyCount      int     `json:"enemy_count"`
	EnemyTypes      []int   `json:"enemy_types"` // EnemyType as int
	SpawnDelay      float64 `json:"spawn_delay"`
	Timeout         float64 `json:"timeout"`
	InterWaveDelay  float64 `json:"inter_wave_delay"`
	MovementPattern int     `json:"movement_pattern"` // MovementPattern as int
}

// LevelMetadata contains descriptive and presentation information about a level
type LevelMetadata struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	MusicTrack  string `json:"music_track"` // Path to music file (optional)
	Background  string `json:"background"`  // Background theme identifier (optional)
}

// BossConfig defines boss enemy configuration
type BossConfig struct {
	Enabled         bool    `json:"enabled"`
	EnemyType       int     `json:"enemy_type"` // EnemyType as int (2 = Boss)
	Health          int     `json:"health"`
	SpawnDelay      float64 `json:"spawn_delay"`   // Delay after all waves complete
	MovementType    string  `json:"movement_type"` // "orbital", "spiral", etc.
	CanShoot        bool    `json:"can_shoot"`
	FireRate        float64 `json:"fire_rate"` // Shots per second
	ProjectileSpeed float64 `json:"projectile_speed"`
	Size            int     `json:"size"`
	Points          int     `json:"points"`
}

// DifficultySettings contains level-specific difficulty multipliers
type DifficultySettings struct {
	EnemySpeedMultiplier     float64 `json:"enemy_speed_multiplier"`      // Default: 1.0
	EnemyHealthMultiplier    float64 `json:"enemy_health_multiplier"`     // Default: 1.0
	EnemySpawnRateMultiplier float64 `json:"enemy_spawn_rate_multiplier"` // Default: 1.0
	PlayerDamageMultiplier   float64 `json:"player_damage_multiplier"`    // Default: 1.0
	ScoreMultiplier          float64 `json:"score_multiplier"`            // Default: 1.0
}

// CompletionConditions defines what must happen for the level to be considered complete
type CompletionConditions struct {
	RequireBossKill         bool `json:"require_boss_kill"`          // If true, boss must be killed
	RequireAllWaves         bool `json:"require_all_waves"`          // If true, all waves must be completed
	RequireAllEnemiesKilled bool `json:"require_all_enemies_killed"` // If true, all enemies must be killed
}

// DefaultDifficultySettings returns default difficulty settings
func DefaultDifficultySettings() DifficultySettings {
	return DifficultySettings{
		EnemySpeedMultiplier:     1.0,
		EnemyHealthMultiplier:    1.0,
		EnemySpawnRateMultiplier: 1.0,
		PlayerDamageMultiplier:   1.0,
		ScoreMultiplier:          1.0,
	}
}

// DefaultBossConfig returns a default boss configuration
func DefaultBossConfig() BossConfig {
	return BossConfig{
		Enabled:         true,
		EnemyType:       2,
		Health:          10,
		SpawnDelay:      2.0,
		MovementType:    "orbital",
		CanShoot:        true,
		FireRate:        2.0,
		ProjectileSpeed: 6.0,
		Size:            64,
		Points:          1000,
	}
}

// DefaultCompletionConditions returns default completion conditions
func DefaultCompletionConditions() CompletionConditions {
	return CompletionConditions{
		RequireBossKill:         true,
		RequireAllWaves:         true,
		RequireAllEnemiesKilled: false,
	}
}
