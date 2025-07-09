package ecs

import (
	"testing"
)

func TestComponentRegistry_New(t *testing.T) {
	registry := NewComponentRegistry()

	if registry == nil {
		t.Fatal("NewComponentRegistry should not return nil")
	}

	if registry.GetComponentCount() != 0 {
		t.Error("New registry should have 0 components")
	}

	if registry.GetTagCount() != 0 {
		t.Error("New registry should have 0 tags")
	}
}

func TestComponentRegistry_RegisterComponent(t *testing.T) {
	registry := NewComponentRegistry()

	// Test successful registration
	err := registry.RegisterComponent("test", Position)
	if err != nil {
		t.Errorf("Failed to register component: %v", err)
	}

	if !registry.HasComponent("test") {
		t.Error("Component should be registered")
	}

	if registry.GetComponentCount() != 1 {
		t.Error("Component count should be 1")
	}

	// Test duplicate registration
	err = registry.RegisterComponent("test", Sprite)
	if err == nil {
		t.Error("Should return error for duplicate registration")
	}

	// Test empty name
	err = registry.RegisterComponent("", Position)
	if err == nil {
		t.Error("Should return error for empty name")
	}

	// Test nil component
	err = registry.RegisterComponent("nil", nil)
	if err == nil {
		t.Error("Should return error for nil component")
	}
}

func TestComponentRegistry_RegisterTag(t *testing.T) {
	registry := NewComponentRegistry()

	// Test successful registration
	err := registry.RegisterTag("test", PlayerTag)
	if err != nil {
		t.Errorf("Failed to register tag: %v", err)
	}

	if !registry.HasTag("test") {
		t.Error("Tag should be registered")
	}

	if registry.GetTagCount() != 1 {
		t.Error("Tag count should be 1")
	}

	// Test duplicate registration
	err = registry.RegisterTag("test", StarTag)
	if err == nil {
		t.Error("Should return error for duplicate registration")
	}

	// Test empty name
	err = registry.RegisterTag("", PlayerTag)
	if err == nil {
		t.Error("Should return error for empty name")
	}

	// Test nil tag
	err = registry.RegisterTag("nil", nil)
	if err == nil {
		t.Error("Should return error for nil tag")
	}
}

func TestComponentRegistry_GetComponent(t *testing.T) {
	registry := NewComponentRegistry()

	// Test getting non-existent component
	_, err := registry.GetComponent("nonexistent")
	if err == nil {
		t.Error("Should return error for non-existent component")
	}

	// Test getting existing component
	registry.RegisterComponent("test", Position)
	component, err := registry.GetComponent("test")
	if err != nil {
		t.Errorf("Failed to get component: %v", err)
	}

	if component != Position {
		t.Error("Should return the correct component")
	}
}

func TestComponentRegistry_GetTag(t *testing.T) {
	registry := NewComponentRegistry()

	// Test getting non-existent tag
	_, err := registry.GetTag("nonexistent")
	if err == nil {
		t.Error("Should return error for non-existent tag")
	}

	// Test getting existing tag
	registry.RegisterTag("test", PlayerTag)
	tag, err := registry.GetTag("test")
	if err != nil {
		t.Errorf("Failed to get tag: %v", err)
	}

	if tag != PlayerTag {
		t.Error("Should return the correct tag")
	}
}

func TestComponentRegistry_GetComponentNames(t *testing.T) {
	registry := NewComponentRegistry()

	// Test empty registry
	names := registry.GetComponentNames()
	if len(names) != 0 {
		t.Error("Empty registry should return empty names list")
	}

	// Test with components
	registry.RegisterComponent("pos", Position)
	registry.RegisterComponent("sprite", Sprite)

	names = registry.GetComponentNames()
	if len(names) != 2 {
		t.Errorf("Expected 2 component names, got %d", len(names))
	}

	// Check that both names are present
	foundPos := false
	foundSprite := false
	for _, name := range names {
		if name == "pos" {
			foundPos = true
		}
		if name == "sprite" {
			foundSprite = true
		}
	}

	if !foundPos || !foundSprite {
		t.Error("Component names should include all registered components")
	}
}

func TestComponentRegistry_GetTagNames(t *testing.T) {
	registry := NewComponentRegistry()

	// Test empty registry
	names := registry.GetTagNames()
	if len(names) != 0 {
		t.Error("Empty registry should return empty names list")
	}

	// Test with tags
	registry.RegisterTag("player", PlayerTag)
	registry.RegisterTag("star", StarTag)

	names = registry.GetTagNames()
	if len(names) != 2 {
		t.Errorf("Expected 2 tag names, got %d", len(names))
	}

	// Check that both names are present
	foundPlayer := false
	foundStar := false
	for _, name := range names {
		if name == "player" {
			foundPlayer = true
		}
		if name == "star" {
			foundStar = true
		}
	}

	if !foundPlayer || !foundStar {
		t.Error("Tag names should include all registered tags")
	}
}

func TestComponentRegistry_Clear(t *testing.T) {
	registry := NewComponentRegistry()

	// Add some components and tags
	registry.RegisterComponent("pos", Position)
	registry.RegisterTag("player", PlayerTag)

	// Verify they exist
	if registry.GetComponentCount() != 1 {
		t.Error("Should have 1 component before clear")
	}
	if registry.GetTagCount() != 1 {
		t.Error("Should have 1 tag before clear")
	}

	// Clear
	registry.Clear()

	// Verify they're gone
	if registry.GetComponentCount() != 0 {
		t.Error("Should have 0 components after clear")
	}
	if registry.GetTagCount() != 0 {
		t.Error("Should have 0 tags after clear")
	}

	if registry.HasComponent("pos") {
		t.Error("Component should not exist after clear")
	}
	if registry.HasTag("player") {
		t.Error("Tag should not exist after clear")
	}
}

func TestComponentRegistry_InitializeDefaultComponents(t *testing.T) {
	registry := NewComponentRegistry()

	// Initialize default components
	err := registry.InitializeDefaultComponents()
	if err != nil {
		t.Fatalf("Failed to initialize default components: %v", err)
	}

	// Check that all expected components are registered
	expectedComponents := []string{"position", "sprite", "speed", "size", "angle", "scale", "orbital"}
	for _, name := range expectedComponents {
		if !registry.HasComponent(name) {
			t.Errorf("Expected component '%s' to be registered", name)
		}
	}

	// Check that all expected tags are registered
	expectedTags := []string{"player", "star"}
	for _, name := range expectedTags {
		if !registry.HasTag(name) {
			t.Errorf("Expected tag '%s' to be registered", name)
		}
	}

	// Check counts
	if registry.GetComponentCount() != len(expectedComponents) {
		t.Errorf("Expected %d components, got %d", len(expectedComponents), registry.GetComponentCount())
	}

	if registry.GetTagCount() != len(expectedTags) {
		t.Errorf("Expected %d tags, got %d", len(expectedTags), registry.GetTagCount())
	}
}

func TestComponentRegistry_ConcurrentAccess(t *testing.T) {
	registry := NewComponentRegistry()

	// Test concurrent registration
	done := make(chan bool, 2)

	go func() {
		for i := 0; i < 100; i++ {
			registry.RegisterComponent("comp"+string(rune(i)), Position)
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 100; i++ {
			registry.RegisterTag("tag"+string(rune(i)), PlayerTag)
		}
		done <- true
	}()

	<-done
	<-done

	// Verify no panics occurred and components were registered
	if registry.GetComponentCount() == 0 {
		t.Error("Components should have been registered")
	}

	if registry.GetTagCount() == 0 {
		t.Error("Tags should have been registered")
	}
}
