package contracts

import (
	"context"
	"fmt"
)

// SystemRegistryImpl implements SystemRegistry with proper type safety
type SystemRegistryImpl struct {
	healthSystem   HealthSystem
	eventSystem    EventSystem
	enemySystem    EnemySystem
	weaponSystem   WeaponSystem
	resourceSystem ResourceSystem
	scoreSystem    ScoreSystem
	stateSystem    StateSystem
}

// NewSystemRegistry creates a new system registry with all required systems
func NewSystemRegistry(
	health HealthSystem,
	events EventSystem,
	enemies EnemySystem,
	weapons WeaponSystem,
	resources ResourceSystem,
	score ScoreSystem,
	state StateSystem,
) SystemRegistry {
	return &SystemRegistryImpl{
		healthSystem:   health,
		eventSystem:    events,
		enemySystem:    enemies,
		weaponSystem:   weapons,
		resourceSystem: resources,
		scoreSystem:    score,
		stateSystem:    state,
	}
}

// Health returns the health system
func (r *SystemRegistryImpl) Health() HealthSystem {
	return r.healthSystem
}

// Events returns the event system
func (r *SystemRegistryImpl) Events() EventSystem {
	return r.eventSystem
}

// Enemies returns the enemy system
func (r *SystemRegistryImpl) Enemies() EnemySystem {
	return r.enemySystem
}

// Weapons returns the weapon system
func (r *SystemRegistryImpl) Weapons() WeaponSystem {
	return r.weaponSystem
}

// Resources returns the resource system
func (r *SystemRegistryImpl) Resources() ResourceSystem {
	return r.resourceSystem
}

// Score returns the score system
func (r *SystemRegistryImpl) Score() ScoreSystem {
	return r.scoreSystem
}

// State returns the state system
func (r *SystemRegistryImpl) State() StateSystem {
	return r.stateSystem
}

// Validate ensures all required systems are properly initialized
func (r *SystemRegistryImpl) Validate(ctx context.Context) error {
	if r.healthSystem == nil {
		return fmt.Errorf("health system is nil")
	}
	if r.eventSystem == nil {
		return fmt.Errorf("event system is nil")
	}
	if r.enemySystem == nil {
		return fmt.Errorf("enemy system is nil")
	}
	if r.weaponSystem == nil {
		return fmt.Errorf("weapon system is nil")
	}
	if r.resourceSystem == nil {
		return fmt.Errorf("resource system is nil")
	}
	if r.scoreSystem == nil {
		return fmt.Errorf("score system is nil")
	}
	if r.stateSystem == nil {
		return fmt.Errorf("state system is nil")
	}
	return nil
}
