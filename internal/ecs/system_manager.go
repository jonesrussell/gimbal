package ecs

import (
	"fmt"

	"github.com/yohamta/donburi"
)

// System interface defines the contract for all systems
type System interface {
	Update(world donburi.World, args ...interface{}) error
	Name() string
}

// SystemManager manages the execution of systems
type SystemManager struct {
	systems []System
}

// NewSystemManager creates a new system manager
func NewSystemManager() *SystemManager {
	return &SystemManager{
		systems: make([]System, 0),
	}
}

// AddSystem adds a system to the manager
func (sm *SystemManager) AddSystem(system System) {
	sm.systems = append(sm.systems, system)
}

// UpdateAll updates all systems in order
func (sm *SystemManager) UpdateAll(world donburi.World, args ...interface{}) error {
	for _, system := range sm.systems {
		if err := system.Update(world, args...); err != nil {
			return fmt.Errorf("system %s failed: %w", system.Name(), err)
		}
	}
	return nil
}

// GetSystemCount returns the number of registered systems
func (sm *SystemManager) GetSystemCount() int {
	return len(sm.systems)
}

// GetSystemNames returns the names of all registered systems
func (sm *SystemManager) GetSystemNames() []string {
	names := make([]string, len(sm.systems))
	for i, system := range sm.systems {
		names[i] = system.Name()
	}
	return names
}

// Clear removes all systems from the manager
func (sm *SystemManager) Clear() {
	sm.systems = make([]System, 0)
}
