package ecs

import (
	"fmt"
	"sync"
)

// ComponentRegistry manages all ECS components
type ComponentRegistry struct {
	mu         sync.RWMutex
	components map[string]interface{}
	tags       map[string]interface{}
}

// NewComponentRegistry creates a new component registry
func NewComponentRegistry() *ComponentRegistry {
	return &ComponentRegistry{
		components: make(map[string]interface{}),
		tags:       make(map[string]interface{}),
	}
}

// RegisterComponent registers a component with the registry
func (cr *ComponentRegistry) RegisterComponent(name string, component interface{}) error {
	cr.mu.Lock()
	defer cr.mu.Unlock()

	if name == "" {
		return fmt.Errorf("component name cannot be empty")
	}

	if component == nil {
		return fmt.Errorf("component cannot be nil")
	}

	if _, exists := cr.components[name]; exists {
		return fmt.Errorf("component '%s' already registered", name)
	}

	cr.components[name] = component
	return nil
}

// RegisterTag registers a tag component with the registry
func (cr *ComponentRegistry) RegisterTag(name string, tag interface{}) error {
	cr.mu.Lock()
	defer cr.mu.Unlock()

	if name == "" {
		return fmt.Errorf("tag name cannot be empty")
	}

	if tag == nil {
		return fmt.Errorf("tag cannot be nil")
	}

	if _, exists := cr.tags[name]; exists {
		return fmt.Errorf("tag '%s' already registered", name)
	}

	cr.tags[name] = tag
	return nil
}

// GetComponent retrieves a component by name
func (cr *ComponentRegistry) GetComponent(name string) (interface{}, error) {
	cr.mu.RLock()
	defer cr.mu.RUnlock()

	component, exists := cr.components[name]
	if !exists {
		return nil, fmt.Errorf("component '%s' not found", name)
	}

	return component, nil
}

// GetTag retrieves a tag by name
func (cr *ComponentRegistry) GetTag(name string) (interface{}, error) {
	cr.mu.RLock()
	defer cr.mu.RUnlock()

	tag, exists := cr.tags[name]
	if !exists {
		return nil, fmt.Errorf("tag '%s' not found", name)
	}

	return tag, nil
}

// HasComponent checks if a component is registered
func (cr *ComponentRegistry) HasComponent(name string) bool {
	cr.mu.RLock()
	defer cr.mu.RUnlock()

	_, exists := cr.components[name]
	return exists
}

// HasTag checks if a tag is registered
func (cr *ComponentRegistry) HasTag(name string) bool {
	cr.mu.RLock()
	defer cr.mu.RUnlock()

	_, exists := cr.tags[name]
	return exists
}

// GetComponentNames returns all registered component names
func (cr *ComponentRegistry) GetComponentNames() []string {
	cr.mu.RLock()
	defer cr.mu.RUnlock()

	names := make([]string, 0, len(cr.components))
	for name := range cr.components {
		names = append(names, name)
	}
	return names
}

// GetTagNames returns all registered tag names
func (cr *ComponentRegistry) GetTagNames() []string {
	cr.mu.RLock()
	defer cr.mu.RUnlock()

	names := make([]string, 0, len(cr.tags))
	for name := range cr.tags {
		names = append(names, name)
	}
	return names
}

// GetComponentCount returns the number of registered components
func (cr *ComponentRegistry) GetComponentCount() int {
	cr.mu.RLock()
	defer cr.mu.RUnlock()

	return len(cr.components)
}

// GetTagCount returns the number of registered tags
func (cr *ComponentRegistry) GetTagCount() int {
	cr.mu.RLock()
	defer cr.mu.RUnlock()

	return len(cr.tags)
}

// Clear removes all registered components and tags
func (cr *ComponentRegistry) Clear() {
	cr.mu.Lock()
	defer cr.mu.Unlock()

	cr.components = make(map[string]interface{})
	cr.tags = make(map[string]interface{})
}

// InitializeDefaultComponents registers all default components
func (cr *ComponentRegistry) InitializeDefaultComponents() error {
	// Register components
	components := map[string]interface{}{
		"position": Position,
		"sprite":   Sprite,
		"speed":    Speed,
		"size":     Size,
		"angle":    Angle,
		"scale":    Scale,
		"orbital":  Orbital,
	}

	for name, component := range components {
		if err := cr.RegisterComponent(name, component); err != nil {
			return fmt.Errorf("failed to register component '%s': %w", name, err)
		}
	}

	// Register tags
	tags := map[string]interface{}{
		"player": PlayerTag,
		"star":   StarTag,
	}

	for name, tag := range tags {
		if err := cr.RegisterTag(name, tag); err != nil {
			return fmt.Errorf("failed to register tag '%s': %w", name, err)
		}
	}

	return nil
}
