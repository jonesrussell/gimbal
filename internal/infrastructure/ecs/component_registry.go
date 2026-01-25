// Package ecs provides ECS infrastructure components following Clean Architecture.
// This package bridges the domain layer with the Donburi ECS framework.
package ecs

import (
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/config"
	"github.com/jonesrussell/gimbal/internal/domain/value"
	"github.com/jonesrussell/gimbal/internal/ecs/core"
	"github.com/jonesrussell/gimbal/internal/math"
)

// ComponentRegistry holds all ECS component type definitions.
// This allows components to be injected into systems for better testability.
// Each ComponentRegistry instance creates its own set of component types,
// enabling isolated testing without global state pollution.
type ComponentRegistry struct {
	// Tags for different entity types
	PlayerTag          *donburi.ComponentType[struct{}]
	StarTag            *donburi.ComponentType[struct{}]
	EnemyTag           *donburi.ComponentType[struct{}]
	ProjectileTag      *donburi.ComponentType[struct{}]
	EnemyProjectileTag *donburi.ComponentType[struct{}]

	// Components
	Position    *donburi.ComponentType[common.Point]
	Sprite      *donburi.ComponentType[*ebiten.Image]
	Movement    *donburi.ComponentType[core.MovementData]
	Orbital     *donburi.ComponentType[core.OrbitalData]
	Size        *donburi.ComponentType[config.Size]
	Speed       *donburi.ComponentType[float64]
	Angle       *donburi.ComponentType[math.Angle]
	Scale       *donburi.ComponentType[float64]
	Health      *donburi.ComponentType[core.HealthData]
	EnemyTypeID *donburi.ComponentType[int]
}

// NewComponentRegistry creates a new component registry with fresh component types.
// Use this for testing to create isolated component sets.
func NewComponentRegistry() *ComponentRegistry {
	return &ComponentRegistry{
		// Tags
		PlayerTag:          donburi.NewComponentType[struct{}](),
		StarTag:            donburi.NewComponentType[struct{}](),
		EnemyTag:           donburi.NewComponentType[struct{}](),
		ProjectileTag:      donburi.NewComponentType[struct{}](),
		EnemyProjectileTag: donburi.NewComponentType[struct{}](),

		// Components
		Position:    donburi.NewComponentType[common.Point](),
		Sprite:      donburi.NewComponentType[*ebiten.Image](),
		Movement:    donburi.NewComponentType[core.MovementData](),
		Orbital:     donburi.NewComponentType[core.OrbitalData](),
		Size:        donburi.NewComponentType[config.Size](),
		Speed:       donburi.NewComponentType[float64](),
		Angle:       donburi.NewComponentType[math.Angle](),
		Scale:       donburi.NewComponentType[float64](),
		Health:      donburi.NewComponentType[core.HealthData](),
		EnemyTypeID: donburi.NewComponentType[int](),
	}
}

// Global registry for backward compatibility
var (
	globalRegistry     *ComponentRegistry
	globalRegistryOnce sync.Once
)

// GlobalRegistry returns the singleton component registry.
// For production use - returns the same registry every time.
// New code should prefer receiving ComponentRegistry via dependency injection.
func GlobalRegistry() *ComponentRegistry {
	globalRegistryOnce.Do(func() {
		globalRegistry = NewComponentRegistry()
	})
	return globalRegistry
}

// MovementPattern re-exports the domain value type for convenience.
type MovementPattern = value.MovementPattern

// Movement pattern constants re-exported for convenience.
const (
	MovementPatternNormal       = value.MovementPatternNormal
	MovementPatternZigzag       = value.MovementPatternZigzag
	MovementPatternAccelerating = value.MovementPatternAccelerating
	MovementPatternPulsing      = value.MovementPatternPulsing
)
