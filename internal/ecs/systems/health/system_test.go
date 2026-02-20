package health

import (
	"context"
	"testing"

	"github.com/yohamta/donburi"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/config"
	"github.com/jonesrussell/gimbal/internal/ecs/core"
)

type mockEventSystem struct{}

func (mockEventSystem) EmitGameOver() {}

type mockGameStateManager struct{}

func (mockGameStateManager) SetGameOver(bool) {}

type testLogger struct{}

func (testLogger) Debug(_ string, _ ...any)                           {}
func (testLogger) DebugContext(_ context.Context, _ string, _ ...any) {}
func (testLogger) Info(_ string, _ ...any)                            {}
func (testLogger) InfoContext(_ context.Context, _ string, _ ...any)  {}
func (testLogger) Warn(_ string, _ ...any)                            {}
func (testLogger) WarnContext(_ context.Context, _ string, _ ...any)  {}
func (testLogger) Error(_ string, _ ...any)                           {}
func (testLogger) ErrorContext(_ context.Context, _ string, _ ...any) {}
func (testLogger) Sync() error                                        { return nil }

func TestDamageEntity_DevInvinciblePreventsDamage(t *testing.T) {
	world := donburi.NewWorld()
	entity := world.Create(core.PlayerTag, core.Health, core.Position, core.Size)
	entry := world.Entry(entity)
	core.Position.SetValue(entry, common.Point{X: 0, Y: 0})
	core.Size.SetValue(entry, config.Size{Width: 48, Height: 48})
	health := core.NewHealthData(3, 3)
	core.Health.SetValue(entry, health)

	cfg := config.NewConfig(config.WithDebug(true), config.WithInvincible(true))
	hs := NewHealthSystem(world, cfg, mockEventSystem{}, mockGameStateManager{}, testLogger{})
	ctx := context.Background()

	err := hs.DamageEntity(ctx, entity, 1)
	if err != nil {
		t.Fatalf("DamageEntity: %v", err)
	}

	current, max, ok := hs.GetHealth(ctx, entity)
	if !ok {
		t.Fatal("GetHealth failed")
	}
	if current != 3 || max != 3 {
		t.Errorf("Dev invincible should prevent damage: got health %d/%d, want 3/3", current, max)
	}
}

func TestDamageEntity_IsInvinciblePreventsDamage(t *testing.T) {
	world := donburi.NewWorld()
	entity := world.Create(core.PlayerTag, core.Health, core.Position, core.Size)
	entry := world.Entry(entity)
	core.Position.SetValue(entry, common.Point{X: 0, Y: 0})
	core.Size.SetValue(entry, config.Size{Width: 48, Height: 48})
	health := core.NewHealthData(3, 3)
	health.IsInvincible = true
	health.InvincibilityTime = 1
	core.Health.SetValue(entry, health)

	cfg := config.NewConfig(config.WithDebug(false))
	hs := NewHealthSystem(world, cfg, mockEventSystem{}, mockGameStateManager{}, testLogger{})
	ctx := context.Background()

	err := hs.DamageEntity(ctx, entity, 1)
	if err != nil {
		t.Fatalf("DamageEntity: %v", err)
	}

	current, max, ok := hs.GetHealth(ctx, entity)
	if !ok {
		t.Fatal("GetHealth failed")
	}
	if current != 3 || max != 3 {
		t.Errorf("IsInvincible should prevent damage: got health %d/%d, want 3/3", current, max)
	}
}

func TestDamageEntity_AppliesDamageWhenNotInvincible(t *testing.T) {
	world := donburi.NewWorld()
	entity := world.Create(core.PlayerTag, core.Health, core.Position, core.Size)
	entry := world.Entry(entity)
	core.Position.SetValue(entry, common.Point{X: 0, Y: 0})
	core.Size.SetValue(entry, config.Size{Width: 48, Height: 48})
	health := core.NewHealthData(3, 3)
	core.Health.SetValue(entry, health)

	cfg := config.NewConfig(config.WithDebug(false))
	hs := NewHealthSystem(world, cfg, mockEventSystem{}, mockGameStateManager{}, testLogger{})
	ctx := context.Background()

	err := hs.DamageEntity(ctx, entity, 1)
	if err != nil {
		t.Fatalf("DamageEntity: %v", err)
	}

	current, max, ok := hs.GetHealth(ctx, entity)
	if !ok {
		t.Fatal("GetHealth failed")
	}
	if current != 2 || max != 3 {
		t.Errorf("Damage should apply when not invincible: got health %d/%d, want 2/3", current, max)
	}

	// Entity should now be invincible (i-frames)
	entry = world.Entry(entity)
	healthPtr := core.Health.Get(entry)
	if !healthPtr.IsInvincible {
		t.Error("Expected entity to be invincible after taking damage (i-frames)")
	}
}
