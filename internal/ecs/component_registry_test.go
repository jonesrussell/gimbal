package ecs_test

import (
	"testing"

	"github.com/jonesrussell/gimbal/internal/ecs"
)

func TestComponentRegistry_New(t *testing.T) {
	registry := ecs.NewComponentRegistry()

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
	registry := ecs.NewComponentRegistry()

	// Test successful registration
	err := registry.RegisterComponent("test", ecs.Position)
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
	err = registry.RegisterComponent("test", ecs.Sprite)
	if err == nil {
		t.Error("Should return error for duplicate registration")
	}

	// Test empty name
	err = registry.RegisterComponent("", ecs.Position)
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
	registry := ecs.NewComponentRegistry()

	// Test successful registration
	err := registry.RegisterTag("test", ecs.PlayerTag)
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
	err = registry.RegisterTag("test", ecs.StarTag)
	if err == nil {
		t.Error("Should return error for duplicate registration")
	}

	// Test empty name
	err = registry.RegisterTag("", ecs.PlayerTag)
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
	registry := ecs.NewComponentRegistry()

	// Test getting non-existent component
	_, err := registry.GetComponent("nonexistent")
	if err == nil {
		t.Error("Should return error for non-existent component")
	}

	// Test getting existing component
	regErr := registry.RegisterComponent("test", ecs.Position)
	if regErr != nil {
		t.Fatalf("Failed to register component: %v", regErr)
	}
	component, err := registry.GetComponent("test")
	if err != nil {
		t.Errorf("Failed to get component: %v", err)
	}

	if component != ecs.Position {
		t.Error("Should return the correct component")
	}
}

func TestComponentRegistry_GetTag(t *testing.T) {
	registry := ecs.NewComponentRegistry()

	// Test getting non-existent tag
	_, err := registry.GetTag("nonexistent")
	if err == nil {
		t.Error("Should return error for non-existent tag")
	}

	// Test getting existing tag
	regErr := registry.RegisterTag("test", ecs.PlayerTag)
	if regErr != nil {
		t.Fatalf("Failed to register tag: %v", regErr)
	}
	tag, err := registry.GetTag("test")
	if err != nil {
		t.Errorf("Failed to get tag: %v", err)
	}

	if tag != ecs.PlayerTag {
		t.Error("Should return the correct tag")
	}
}

func TestComponentRegistry_GetComponentNames(t *testing.T) {
	registry := ecs.NewComponentRegistry()

	// Test empty registry
	names := registry.GetComponentNames()
	if len(names) != 0 {
		t.Error("Empty registry should return empty names list")
	}

	// Test with components
	regErr := registry.RegisterComponent("pos", ecs.Position)
	if regErr != nil {
		t.Fatalf("Failed to register position component: %v", regErr)
	}
	regErr = registry.RegisterComponent("sprite", ecs.Sprite)
	if regErr != nil {
		t.Fatalf("Failed to register sprite component: %v", regErr)
	}

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
	registry := ecs.NewComponentRegistry()

	// Test empty registry
	names := registry.GetTagNames()
	if len(names) != 0 {
		t.Error("Empty registry should return empty names list")
	}

	// Test with tags
	regErr := registry.RegisterTag("player", ecs.PlayerTag)
	if regErr != nil {
		t.Fatalf("Failed to register player tag: %v", regErr)
	}
	regErr = registry.RegisterTag("star", ecs.StarTag)
	if regErr != nil {
		t.Fatalf("Failed to register star tag: %v", regErr)
	}

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
	registry := ecs.NewComponentRegistry()

	// Add some components and tags
	regErr := registry.RegisterComponent("pos", ecs.Position)
	if regErr != nil {
		t.Fatalf("Failed to register position component: %v", regErr)
	}
	regErr = registry.RegisterTag("player", ecs.PlayerTag)
	if regErr != nil {
		t.Fatalf("Failed to register player tag: %v", regErr)
	}

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
	registry := ecs.NewComponentRegistry()

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
	registry := ecs.NewComponentRegistry()

	// Test concurrent registration
	done := make(chan bool, 2)

	go func() {
		for i := 0; i < 100; i++ {
			regErr := registry.RegisterComponent("comp"+string(rune(i)), ecs.Position)
			if regErr != nil {
				t.Errorf("Failed to register component in goroutine: %v", regErr)
			}
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 100; i++ {
			regErr := registry.RegisterTag("tag"+string(rune(i)), ecs.PlayerTag)
			if regErr != nil {
				t.Errorf("Failed to register tag in goroutine: %v", regErr)
			}
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
