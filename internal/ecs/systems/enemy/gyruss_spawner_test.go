//nolint:testpackage // White box tests need access to internal methods
package enemy

import (
	"context"
	"testing"

	"github.com/yohamta/donburi"

	"github.com/jonesrussell/gimbal/internal/config"
	resources "github.com/jonesrussell/gimbal/internal/ecs/managers/resource"
)

type testLogger struct{}

func (m *testLogger) Debug(_ string, _ ...any)                           {}
func (m *testLogger) DebugContext(_ context.Context, _ string, _ ...any) {}
func (m *testLogger) Info(_ string, _ ...any)                            {}
func (m *testLogger) InfoContext(_ context.Context, _ string, _ ...any)  {}
func (m *testLogger) Warn(_ string, _ ...any)                            {}
func (m *testLogger) WarnContext(_ context.Context, _ string, _ ...any)  {}
func (m *testLogger) Error(_ string, _ ...any)                           {}
func (m *testLogger) ErrorContext(_ context.Context, _ string, _ ...any) {}
func (m *testLogger) Sync() error                                        { return nil }

func TestNewGyrussSpawner(t *testing.T) {
	world := donburi.NewWorld()
	gameConfig := &config.GameConfig{
		ScreenSize: config.Size{Width: 800, Height: 600},
	}
	logger := &testLogger{}
	ctx := context.Background()
	resourceMgr := resources.NewResourceManager(ctx, logger)

	spawner := NewGyrussSpawner(world, gameConfig, resourceMgr, logger)

	if spawner == nil {
		t.Fatal("Expected spawner to be created")
	}
	if spawner.world != world {
		t.Error("Expected world to be set")
	}
	if spawner.screenCenter.X != 400 || spawner.screenCenter.Y != 300 {
		t.Errorf("Expected screen center (400, 300), got (%f, %f)",
			spawner.screenCenter.X, spawner.screenCenter.Y)
	}
}

func TestGyrussSpawner_GetEnemyType(t *testing.T) {
	world := donburi.NewWorld()
	gameConfig := &config.GameConfig{
		ScreenSize: config.Size{Width: 800, Height: 600},
	}
	logger := &testLogger{}
	ctx := context.Background()
	resourceMgr := resources.NewResourceManager(ctx, logger)
	spawner := NewGyrussSpawner(world, gameConfig, resourceMgr, logger)

	tests := []struct {
		typeStr  string
		expected EnemyType
	}{
		{EnemyTypeStrBasic, EnemyTypeBasic},
		{EnemyTypeStrHeavy, EnemyTypeHeavy},
		{EnemyTypeStrBoss, EnemyTypeBoss},
		{EnemyTypeStrSatellite, EnemyTypeBasic},
		{"unknown", EnemyTypeBasic},
	}

	for _, tt := range tests {
		result := spawner.getEnemyType(tt.typeStr)
		if result != tt.expected {
			t.Errorf("getEnemyType(%q) = %d, expected %d", tt.typeStr, result, tt.expected)
		}
	}
}

func TestGyrussSpawner_GetHealthForType(t *testing.T) {
	world := donburi.NewWorld()
	gameConfig := &config.GameConfig{
		ScreenSize: config.Size{Width: 800, Height: 600},
	}
	logger := &testLogger{}
	ctx := context.Background()
	resourceMgr := resources.NewResourceManager(ctx, logger)
	spawner := NewGyrussSpawner(world, gameConfig, resourceMgr, logger)

	tests := []struct {
		typeStr  string
		expected int
	}{
		{EnemyTypeStrBasic, 1},
		{EnemyTypeStrHeavy, 3},
		{EnemyTypeStrSatellite, 1},
		{"unknown", 1},
	}

	for _, tt := range tests {
		result := spawner.getHealthForType(tt.typeStr)
		if result != tt.expected {
			t.Errorf("getHealthForType(%q) = %d, expected %d", tt.typeStr, result, tt.expected)
		}
	}
}

func TestGyrussSpawner_GetOrbitRadius(t *testing.T) {
	world := donburi.NewWorld()
	logger := &testLogger{}
	ctx := context.Background()
	resourceMgr := resources.NewResourceManager(ctx, logger)

	tests := []struct {
		width    int
		height   int
		expected float64
	}{
		{800, 600, 600 * 0.35}, // Height is smaller
		{600, 800, 600 * 0.35}, // Width is smaller
		{500, 500, 500 * 0.35}, // Equal
	}

	for _, tt := range tests {
		gameConfig := &config.GameConfig{
			ScreenSize: config.Size{Width: tt.width, Height: tt.height},
		}
		spawner := NewGyrussSpawner(world, gameConfig, resourceMgr, logger)
		result := spawner.getOrbitRadius()
		if result != tt.expected {
			t.Errorf("getOrbitRadius() with %dx%d = %f, expected %f",
				tt.width, tt.height, result, tt.expected)
		}
	}
}
