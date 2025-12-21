package enemy

// EnemyType represents the type of enemy
type EnemyType int

const (
	// EnemyTypeBasic is the standard fast enemy (1 HP, fast, random outward movement)
	EnemyTypeBasic EnemyType = iota
	// EnemyTypeHeavy is a slower, tougher enemy (2 HP, slower, spiral movement)
	EnemyTypeHeavy
	// EnemyTypeBoss is the boss enemy (10+ HP, large, orbital movement)
	EnemyTypeBoss
)

// String returns a human-readable string representation of the enemy type
func (et EnemyType) String() string {
	switch et {
	case EnemyTypeBasic:
		return "Basic"
	case EnemyTypeHeavy:
		return "Heavy"
	case EnemyTypeBoss:
		return "Boss"
	default:
		return "Unknown"
	}
}

// MovementPattern represents the movement behavior pattern for enemies
type MovementPattern int

const (
	// MovementPatternNormal is standard outward movement
	MovementPatternNormal MovementPattern = iota
	// MovementPatternZigzag oscillates side-to-side while moving outward
	MovementPatternZigzag
	// MovementPatternAccelerating starts slow and speeds up
	MovementPatternAccelerating
	// MovementPatternPulsing moves in bursts (fast-slow-fast)
	MovementPatternPulsing
)

// EnemyTypeData contains configuration for each enemy type
type EnemyTypeData struct {
	Type            EnemyType
	Health          int
	Speed           float64
	Size            int
	Points          int
	SpriteName      string
	MovementType    string          // "outward", "spiral", "orbital"
	MovementPattern MovementPattern // Movement behavior pattern
	CanShoot        bool            // Whether this enemy type can shoot
	FireRate        float64         // Shots per second (0 = no shooting)
	ProjectileSpeed float64         // Speed of enemy projectiles
}

// GetEnemyTypeData returns the configuration for an enemy type
func GetEnemyTypeData(enemyType EnemyType) EnemyTypeData {
	switch enemyType {
	case EnemyTypeBasic:
		return EnemyTypeData{
			Type:            EnemyTypeBasic,
			Health:          1,
			Speed:           2.0,
			Size:            32,
			Points:          100,
			SpriteName:      "enemy",
			MovementType:    "outward",
			MovementPattern: MovementPatternNormal,
			CanShoot:        true,
			FireRate:        0.5, // Shoots every 2 seconds
			ProjectileSpeed: 4.0,
		}
	case EnemyTypeHeavy:
		return EnemyTypeData{
			Type:            EnemyTypeHeavy,
			Health:          2,
			Speed:           1.5,
			Size:            32,
			Points:          200,
			SpriteName:      "enemy_heavy",
			MovementType:    "spiral",
			MovementPattern: MovementPatternNormal,
			CanShoot:        true,
			FireRate:        1.0, // Shoots every second
			ProjectileSpeed: 5.0,
		}
	case EnemyTypeBoss:
		return EnemyTypeData{
			Type:            EnemyTypeBoss,
			Health:          10,
			Speed:           1.0,
			Size:            64,
			Points:          1000,
			SpriteName:      "enemy_boss",
			MovementType:    "orbital",
			MovementPattern: MovementPatternNormal,
			CanShoot:        true,
			FireRate:        2.0, // Shoots twice per second
			ProjectileSpeed: 6.0,
		}
	default:
		return GetEnemyTypeData(EnemyTypeBasic)
	}
}
