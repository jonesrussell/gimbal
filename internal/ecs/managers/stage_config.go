package managers

// StageConfig defines the complete configuration for a Gyruss-style stage
type StageConfig struct {
	StageNumber int            `json:"stage_number"`
	Planet      string         `json:"planet"`
	Metadata    StageMetadata  `json:"metadata"`
	Waves       []GyrussWave   `json:"waves"`
	Boss        StageBossConfig `json:"boss"`
	PowerUps    PowerUpConfig  `json:"power_ups"`
	Difficulty  DifficultySettings `json:"difficulty"`
}

// StageMetadata contains descriptive information about a stage
type StageMetadata struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	MusicTrack  string `json:"music_track"`
	Background  string `json:"background"`
}

// GyrussWave defines a wave in Gyruss format
type GyrussWave struct {
	WaveID        string             `json:"wave_id"`
	Description   string             `json:"description"`
	SpawnSequence []EnemyGroupConfig `json:"spawn_sequence"`
	OnClear       string             `json:"on_clear"` // "next_wave", "boss", etc.
	Timing        WaveTiming         `json:"timing"`
}

// EnemyGroupConfig defines a group of enemies in a wave
type EnemyGroupConfig struct {
	EnemyType         string            `json:"enemy_type"`
	Count             int               `json:"count"`
	SpawnDelay        float64           `json:"spawn_delay"`
	SpawnInterval     float64           `json:"spawn_interval"`
	EntryPath         EntryPathConfig   `json:"entry_path"`
	ScaleAnimation    ScaleAnimConfig   `json:"scale_animation"`
	Behavior          BehaviorConfig    `json:"behavior"`
	AttackPattern     AttackConfig      `json:"attack_pattern"`
	FirePattern       FireConfig        `json:"fire_pattern"`
	Retreat           RetreatConfig     `json:"retreat"`
	PowerUpTrigger    bool              `json:"powerup_trigger"`
	PowerUpType       string            `json:"powerup_type"`
}

// EntryPathConfig defines entry path parameters
type EntryPathConfig struct {
	Type       string             `json:"type"` // "spiral_in", "arc_sweep", "straight_in", "loop_entry"
	Duration   float64            `json:"duration"`
	Parameters EntryPathParams    `json:"parameters"`
}

// EntryPathParams contains path-specific parameters
type EntryPathParams struct {
	SpiralTurns       float64 `json:"spiral_turns"`
	ArcAngle          float64 `json:"arc_angle"`
	RotationDirection string  `json:"rotation_direction"` // "clockwise", "counter_clockwise"
	StartRadius       float64 `json:"start_radius"`
}

// ScaleAnimConfig defines scale animation parameters
type ScaleAnimConfig struct {
	StartScale float64 `json:"start_scale"`
	EndScale   float64 `json:"end_scale"`
	Easing     string  `json:"easing"` // "linear", "ease_in", "ease_out", "ease_in_out"
}

// BehaviorConfig defines post-entry behavior
type BehaviorConfig struct {
	PostEntry      string  `json:"post_entry"` // "orbit_only", "orbit_then_attack", "immediate_attack", "hover_center_then_orbit"
	OrbitDuration  float64 `json:"orbit_duration"`
	OrbitDirection string  `json:"orbit_direction"` // "clockwise", "counter_clockwise"
	OrbitSpeed     float64 `json:"orbit_speed"`
	MaxAttacks     int     `json:"max_attacks"`
}

// AttackConfig defines attack pattern parameters
type AttackConfig struct {
	Type       string  `json:"type"` // "none", "single_rush", "paired_rush", "loopback_rush", "suicide_dive"
	Cooldown   float64 `json:"cooldown"`
	RushSpeed  float64 `json:"rush_speed"`
	ReturnSpeed float64 `json:"return_speed"`
}

// FireConfig defines fire pattern parameters
type FireConfig struct {
	Type              string  `json:"type"` // "none", "single_shot", "burst", "spray"
	FireRate          float64 `json:"fire_rate"`
	BurstCount        int     `json:"burst_count"`
	SprayAngle        float64 `json:"spray_angle"`
	ProjectileCount   int     `json:"projectile_count"`
	FireWhileOrbit    bool    `json:"fire_while_orbit"`
	FireWhileAttack   bool    `json:"fire_while_attack"`
}

// RetreatConfig defines retreat behavior
type RetreatConfig struct {
	HealthThreshold float64 `json:"health_threshold"`
	Timeout         float64 `json:"timeout"`
	Speed           float64 `json:"speed"`
}

// WaveTiming defines wave timing parameters
type WaveTiming struct {
	InterWaveDelay float64 `json:"inter_wave_delay"`
	Timeout        float64 `json:"timeout"`
}

// StageBossConfig defines boss configuration for a stage
type StageBossConfig struct {
	Enabled       bool          `json:"enabled"`
	BossType      string        `json:"boss_type"`
	Health        int           `json:"health"`
	Size          int           `json:"size"`
	EntryPath     EntryPathConfig `json:"entry_path"`
	Behavior      BehaviorConfig  `json:"behavior"`
	AttackPattern AttackConfig    `json:"attack_pattern"`
	FirePattern   FireConfig      `json:"fire_pattern"`
	SpawnDelay    float64       `json:"spawn_delay"`
	Points        int           `json:"points"`
}

// PowerUpConfig defines power-up drop configuration
type PowerUpConfig struct {
	DropChance float64            `json:"drop_chance"`
	Types      []PowerUpTypeConfig `json:"types"`
}

// PowerUpTypeConfig defines a power-up type configuration
type PowerUpTypeConfig struct {
	Type     string  `json:"type"` // "double_shot", "extra_life"
	Weight   float64 `json:"weight"`
	Duration float64 `json:"duration"`
}

// DifficultySettings contains level-specific difficulty multipliers
type DifficultySettings struct {
	EnemySpeedMultiplier     float64 `json:"enemy_speed_multiplier"`
	EnemyHealthMultiplier    float64 `json:"enemy_health_multiplier"`
	EnemySpawnRateMultiplier float64 `json:"enemy_spawn_rate_multiplier"`
	PlayerDamageMultiplier   float64 `json:"player_damage_multiplier"`
	ScoreMultiplier          float64 `json:"score_multiplier"`
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

// ConvertDirection converts string direction to int (1=clockwise, -1=counter-clockwise)
func ConvertDirection(dir string) int {
	switch dir {
	case "counter_clockwise", "ccw":
		return -1
	default:
		return 1
	}
}

// PathTypeFromString converts string to path type integer
func PathTypeFromString(pathType string) int {
	switch pathType {
	case "spiral_in":
		return 0 // PathTypeSpiralIn
	case "arc_sweep":
		return 1 // PathTypeArcSweep
	case "straight_in":
		return 2 // PathTypeStraightIn
	case "loop_entry":
		return 3 // PathTypeLoopEntry
	case "random_outward":
		return 4 // PathTypeRandomOutward
	default:
		return 0
	}
}

// BehaviorFromString converts string to behavior integer
func BehaviorFromString(behavior string) int {
	switch behavior {
	case "orbit_only":
		return 0 // BehaviorOrbitOnly
	case "orbit_then_attack":
		return 1 // BehaviorOrbitThenAttack
	case "immediate_attack":
		return 2 // BehaviorImmediateAttack
	case "hover_center_then_orbit":
		return 3 // BehaviorHoverCenterThenOrbit
	default:
		return 1
	}
}

// AttackPatternFromString converts string to attack pattern integer
func AttackPatternFromString(attack string) int {
	switch attack {
	case "none":
		return 0 // AttackNone
	case "single_rush":
		return 1 // AttackSingleRush
	case "paired_rush":
		return 2 // AttackPairedRush
	case "loopback_rush":
		return 3 // AttackLoopbackRush
	case "suicide_dive":
		return 4 // AttackSuicideDive
	default:
		return 0
	}
}

// FirePatternFromString converts string to fire pattern integer
func FirePatternFromString(fire string) int {
	switch fire {
	case "none":
		return 0 // FireNone
	case "single_shot":
		return 1 // FireSingleShot
	case "burst":
		return 2 // FireBurst
	case "spray":
		return 3 // FireSpray
	default:
		return 0
	}
}

// EasingFromString converts string to easing integer
func EasingFromString(easing string) int {
	switch easing {
	case "ease_in":
		return 1 // EasingEaseIn
	case "ease_out":
		return 2 // EasingEaseOut
	case "ease_in_out":
		return 3 // EasingEaseInOut
	default:
		return 0 // EasingLinear
	}
}
