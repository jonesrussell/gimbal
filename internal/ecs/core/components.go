package core

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/config"
	"github.com/jonesrussell/gimbal/internal/math"
)

// Component tags for different entity types
var (
	// PlayerTag marks an entity as a player
	PlayerTag = donburi.NewTag()
	// StarTag marks an entity as a star
	StarTag = donburi.NewTag()
	// EnemyTag marks an entity as an enemy
	EnemyTag = donburi.NewTag()
	// ProjectileTag marks an entity as a projectile
	ProjectileTag = donburi.NewTag()
)

// Components
var (
	// Position component stores entity position
	Position = donburi.NewComponentType[common.Point]()
	// Sprite component stores the entity's sprite
	Sprite = donburi.NewComponentType[*ebiten.Image]()
	// Movement component stores movement data
	Movement = donburi.NewComponentType[MovementData]()
	// Orbital component stores orbital movement data
	Orbital = donburi.NewComponentType[OrbitalData]()
	// Size component stores entity dimensions
	Size = donburi.NewComponentType[config.Size]()
	// Speed component stores movement speed
	Speed = donburi.NewComponentType[float64]()
	// Angle component stores rotation angle
	Angle = donburi.NewComponentType[math.Angle]()
	// Scale component stores scaling factor
	Scale = donburi.NewComponentType[float64]()
	// Health component stores entity health data
	Health = donburi.NewComponentType[HealthData]()
)

// MovementData represents movement information
type MovementData struct {
	Velocity common.Point
	MaxSpeed float64
}

// OrbitalData represents orbital movement information
type OrbitalData struct {
	Center       common.Point
	Radius       float64
	OrbitalAngle math.Angle
	FacingAngle  math.Angle
}

// HealthData represents health and invincibility information
type HealthData struct {
	Current               int     // Current health/lives
	Maximum               int     // Maximum health/lives
	InvincibilityTime     float64 // Time remaining for invincibility (seconds)
	IsInvincible          bool    // Whether entity is currently invincible
	InvincibilityDuration float64 // Duration of invincibility when hit (seconds)
}

// NewHealthData creates a new health data with default invincibility duration
func NewHealthData(current, maximum int) HealthData {
	return HealthData{
		Current:               current,
		Maximum:               maximum,
		InvincibilityTime:     0,
		IsInvincible:          false,
		InvincibilityDuration: 2.0, // 2 seconds of invincibility
	}
}
