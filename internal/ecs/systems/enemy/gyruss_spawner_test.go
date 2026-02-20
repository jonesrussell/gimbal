//nolint:testpackage // White box tests need access to internal methods
package enemy

import (
	"context"
	"testing"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/jonesrussell/gimbal/internal/config"
	"github.com/jonesrussell/gimbal/internal/ecs/core"
	"github.com/jonesrussell/gimbal/internal/ecs/managers"
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

func TestGyrussSpawner_SpawnIndexOrbitAngle(t *testing.T) {
	world := donburi.NewWorld()
	gameConfig := &config.GameConfig{
		ScreenSize: config.Size{Width: 800, Height: 600},
	}
	logger := &testLogger{}
	ctx := context.Background()
	resourceMgr := resources.NewResourceManager(ctx, logger)
	spawner := NewGyrussSpawner(world, gameConfig, resourceMgr, logger)

	groupConfig := managers.EnemyGroupConfig{
		EnemyType: "basic",
		Count:     3,
		EntryPath: managers.EntryPathConfig{
			Type:     "spiral_in",
			Duration: 2.0,
			Parameters: managers.EntryPathParams{
				SpiralTurns:       1.5,
				RotationDirection: "clockwise",
				StartRadius:       20,
			},
		},
		Behavior: managers.BehaviorConfig{
			PostEntry:      "orbit_then_attack",
			OrbitDuration:  3.0,
			OrbitDirection: "clockwise",
			OrbitSpeed:     45.0,
			MaxAttacks:     2,
		},
		ScaleAnimation: managers.ScaleAnimConfig{
			StartScale: 0.1,
			EndScale:   1.0,
			Easing:     "ease_out",
		},
		AttackPattern: managers.AttackConfig{
			Type:        "single_rush",
			Cooldown:    5.0,
			RushSpeed:   300.0,
			ReturnSpeed: 200.0,
		},
		FirePattern: managers.FireConfig{
			Type:            "single_shot",
			FireRate:        0.5,
			BurstCount:      0,
			SprayAngle:      0,
			ProjectileCount: 1,
			FireWhileOrbit:  true,
			FireWhileAttack: false,
		},
		Retreat: managers.RetreatConfig{
			Timeout: 15.0,
			Speed:   200.0,
		},
	}

	_ = spawner.SpawnEnemy(ctx, &groupConfig, 0)
	_ = spawner.SpawnEnemy(ctx, &groupConfig, 1)
	_ = spawner.SpawnEnemy(ctx, &groupConfig, 2)

	var endPositions []struct{ X, Y float64 }
	query.NewQuery(
		filter.And(
			filter.Contains(core.EnemyTag),
			filter.Contains(core.EntryPath),
		),
	).Each(world, func(entry *donburi.Entry) {
		path := core.EntryPath.Get(entry)
		endPositions = append(endPositions, struct{ X, Y float64 }{path.EndPosition.X, path.EndPosition.Y})
	})

	if len(endPositions) != 3 {
		t.Fatalf("expected 3 spawned enemies, got %d", len(endPositions))
	}

	// Spawn index i gives baseAngle = i * (2π/3), so end positions must differ
	seen := make(map[struct{ X, Y float64 }]bool)
	for i, p := range endPositions {
		key := struct{ X, Y float64 }{p.X, p.Y}
		if seen[key] {
			t.Errorf("spawn index %d: duplicate end position (%.2f, %.2f)", i, p.X, p.Y)
		}
		seen[key] = true
	}
}
