package managers

// PlayerConfig defines player entity configuration
type PlayerConfig struct {
	Health                int     `json:"health"`
	Size                  int     `json:"size"`
	SpriteName            string  `json:"sprite_name"`
	InvincibilityDuration float64 `json:"invincibility_duration"`
}

// EnemyTypeConfig defines configuration for a single enemy type (JSON representation)
// Uses int for MovementPattern to avoid import cycles
type EnemyTypeConfig struct {
	Type            string  `json:"type"` // "basic", "heavy", "boss"
	Health          int     `json:"health"`
	Speed           float64 `json:"speed"`
	Size            int     `json:"size"`
	Points          int     `json:"points"`
	SpriteName      string  `json:"sprite_name"`
	MovementType    string  `json:"movement_type"`    // "outward", "spiral", "orbital"
	MovementPattern int     `json:"movement_pattern"` // MovementPattern as int (0=normal, 1=zigzag, 2=accelerating, 3=pulsing)
	CanShoot        bool    `json:"can_shoot"`
	FireRate        float64 `json:"fire_rate"` // Shots per second
	ProjectileSpeed float64 `json:"projectile_speed"`
}

// EnemyConfigs contains all enemy type configurations
type EnemyConfigs struct {
	EnemyTypes []EnemyTypeConfig `json:"enemy_types"`
}
