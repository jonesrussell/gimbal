//nolint:testpackage // White box tests need access to internal fields
package gyruss

import (
	"context"
	"testing"

	"github.com/yohamta/donburi"

	"github.com/jonesrussell/gimbal/assets"
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

func createTestGyrussSystem(t *testing.T) *GyrussSystem {
	t.Helper()
	world := donburi.NewWorld()
	gameConfig := &config.GameConfig{
		ScreenSize: config.Size{Width: 800, Height: 600},
	}
	logger := &testLogger{}
	ctx := context.Background()
	resourceMgr := resources.NewResourceManager(ctx, logger)

	return NewGyrussSystem(&GyrussSystemConfig{
		World:       world,
		GameConfig:  gameConfig,
		ResourceMgr: resourceMgr,
		Logger:      logger,
		AssetsFS:    assets.Assets,
	})
}

func TestNewGyrussSystem(t *testing.T) {
	gs := createTestGyrussSystem(t)

	if gs == nil {
		t.Fatal("Expected GyrussSystem to be created")
	}
	if gs.currentStage != 1 {
		t.Errorf("Expected initial stage 1, got %d", gs.currentStage)
	}
	if gs.bossSpawned {
		t.Error("Expected bossSpawned to be false initially")
	}
}

func TestGyrussSystem_LoadStage(t *testing.T) {
	gs := createTestGyrussSystem(t)

	err := gs.LoadStage(1)
	if err != nil {
		t.Fatalf("Failed to load stage 1: %v", err)
	}

	if gs.currentStage != 1 {
		t.Errorf("Expected current stage 1, got %d", gs.currentStage)
	}

	waveManager := gs.GetWaveManager()
	if waveManager == nil {
		t.Fatal("Expected wave manager to be set")
	}
	if waveManager.GetStageConfig() == nil {
		t.Error("Expected stage config to be loaded")
	}
}

func TestGyrussSystem_GetCurrentStage(t *testing.T) {
	gs := createTestGyrussSystem(t)

	if gs.GetCurrentStage() != 1 {
		t.Errorf("Expected initial stage 1, got %d", gs.GetCurrentStage())
	}

	if err := gs.LoadStage(2); err != nil {
		t.Fatalf("Failed to load stage 2: %v", err)
	}
	if gs.GetCurrentStage() != 2 {
		t.Errorf("Expected stage 2 after load, got %d", gs.GetCurrentStage())
	}
}

func TestGyrussSystem_IsBossActive(t *testing.T) {
	gs := createTestGyrussSystem(t)

	if gs.IsBossActive() {
		t.Error("Expected boss to not be active initially")
	}

	// Boss spawned but no boss entity = defeated
	gs.bossSpawned = true
	if gs.IsBossActive() {
		t.Error("Expected boss to not be active when spawned but no entity exists")
	}
}

func TestGyrussSystem_WasBossSpawned(t *testing.T) {
	gs := createTestGyrussSystem(t)

	if gs.WasBossSpawned() {
		t.Error("Expected boss not spawned initially")
	}

	gs.bossSpawned = true
	if !gs.WasBossSpawned() {
		t.Error("Expected boss spawned to be true")
	}
}

func TestGyrussSystem_IsStageComplete(t *testing.T) {
	gs := createTestGyrussSystem(t)

	// Stage not complete - boss not spawned
	if gs.IsStageComplete() {
		t.Error("Expected stage not complete when boss not spawned")
	}

	// Boss spawned and defeated (no boss entity) = stage complete
	gs.bossSpawned = true
	if !gs.IsStageComplete() {
		t.Error("Expected stage complete when boss spawned and defeated")
	}
}

func TestGyrussSystem_Reset(t *testing.T) {
	gs := createTestGyrussSystem(t)

	gs.bossSpawned = true
	gs.bossTimer = 5.0

	gs.Reset()

	if gs.bossSpawned {
		t.Error("Expected bossSpawned to be reset to false")
	}
	if gs.bossTimer != 0 {
		t.Error("Expected bossTimer to be reset to 0")
	}
}

func TestGyrussSystem_LoadNextStage(t *testing.T) {
	gs := createTestGyrussSystem(t)

	if err := gs.LoadStage(1); err != nil {
		t.Fatalf("Failed to load stage 1: %v", err)
	}
	err := gs.LoadNextStage()
	if err != nil {
		t.Fatalf("Failed to load next stage: %v", err)
	}

	if gs.GetCurrentStage() != 2 {
		t.Errorf("Expected stage 2 after LoadNextStage, got %d", gs.GetCurrentStage())
	}
}

func TestGyrussSystem_GetWaveManager(t *testing.T) {
	gs := createTestGyrussSystem(t)

	wm := gs.GetWaveManager()
	if wm == nil {
		t.Fatal("Expected wave manager to be returned")
	}
}

func TestGyrussSystem_GetPowerUpSystem(t *testing.T) {
	gs := createTestGyrussSystem(t)

	ps := gs.GetPowerUpSystem()
	if ps == nil {
		t.Fatal("Expected power-up system to be returned")
	}
}
