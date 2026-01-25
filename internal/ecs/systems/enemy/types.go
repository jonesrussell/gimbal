package enemy

import (
	"fmt"

	"github.com/jonesrussell/gimbal/internal/domain/value"
	"github.com/jonesrussell/gimbal/internal/ecs/managers"
)

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

// MovementPattern is a type alias for backward compatibility within this package.
// The canonical definition is in domain/value package.
type MovementPattern = value.MovementPattern

// Re-export movement pattern constants for backward compatibility.
const (
	MovementPatternNormal       = value.MovementPatternNormal
	MovementPatternZigzag       = value.MovementPatternZigzag
	MovementPatternAccelerating = value.MovementPatternAccelerating
	MovementPatternPulsing      = value.MovementPatternPulsing
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

// ConvertEnemyTypeConfig converts a managers.EnemyTypeConfig to enemy.EnemyTypeData
func ConvertEnemyTypeConfig(config *managers.EnemyTypeConfig, enemyType EnemyType) (EnemyTypeData, error) {
	// Convert movement pattern from int to MovementPattern
	var movementPattern MovementPattern
	switch config.MovementPattern {
	case 0:
		movementPattern = MovementPatternNormal
	case 1:
		movementPattern = MovementPatternZigzag
	case 2:
		movementPattern = MovementPatternAccelerating
	case 3:
		movementPattern = MovementPatternPulsing
	default:
		return EnemyTypeData{}, fmt.Errorf("invalid movement_pattern: %d (must be 0-3)", config.MovementPattern)
	}

	return EnemyTypeData{
		Type:            enemyType,
		Health:          config.Health,
		Speed:           config.Speed,
		Size:            config.Size,
		Points:          config.Points,
		SpriteName:      config.SpriteName,
		MovementType:    config.MovementType,
		MovementPattern: movementPattern,
		CanShoot:        config.CanShoot,
		FireRate:        config.FireRate,
		ProjectileSpeed: config.ProjectileSpeed,
	}, nil
}

// GetEnemyTypeFromString converts a string type name to EnemyType
func GetEnemyTypeFromString(typeStr string) (EnemyType, error) {
	switch typeStr {
	case "basic":
		return EnemyTypeBasic, nil
	case "heavy":
		return EnemyTypeHeavy, nil
	case "boss":
		return EnemyTypeBoss, nil
	default:
		return EnemyTypeBasic, fmt.Errorf("unknown enemy type: %s", typeStr)
	}
}
