//nolint:testpackage // White box tests need access to internal fields
package enemy

import (
	"context"
	"testing"

	"github.com/yohamta/donburi"

	"github.com/jonesrussell/gimbal/internal/ecs/managers"
)

type mockLogger struct{}

func (m *mockLogger) Debug(_ string, _ ...any)                           {}
func (m *mockLogger) DebugContext(_ context.Context, _ string, _ ...any) {}
func (m *mockLogger) Info(_ string, _ ...any)                            {}
func (m *mockLogger) InfoContext(_ context.Context, _ string, _ ...any)  {}
func (m *mockLogger) Warn(_ string, _ ...any)                            {}
func (m *mockLogger) WarnContext(_ context.Context, _ string, _ ...any)  {}
func (m *mockLogger) Error(_ string, _ ...any)                           {}
func (m *mockLogger) ErrorContext(_ context.Context, _ string, _ ...any) {}
func (m *mockLogger) Sync() error                                        { return nil }

func TestNewGyrussWaveManager(t *testing.T) {
	world := donburi.NewWorld()
	logger := &mockLogger{}

	wm := NewGyrussWaveManager(world, logger)

	if wm == nil {
		t.Fatal("Expected wave manager to be created")
	}
	if wm.world != world {
		t.Error("Expected world to be set")
	}
}

func TestGyrussWaveManager_LoadStage(t *testing.T) {
	world := donburi.NewWorld()
	logger := &mockLogger{}
	wm := NewGyrussWaveManager(world, logger)

	config := &managers.StageConfig{
		StageNumber: 1,
		Planet:      "Earth",
		Metadata: managers.StageMetadata{
			Name: "Test Stage",
		},
		Waves: []managers.GyrussWave{
			{
				WaveID:      "wave_1",
				Description: "Test wave",
				SpawnSequence: []managers.EnemyGroupConfig{
					{EnemyType: "basic", Count: 5},
				},
			},
		},
	}

	wm.LoadStage(config)

	if wm.GetStageConfig() != config {
		t.Error("Expected stage config to be set")
	}
	if wm.GetWaveCount() != 1 {
		t.Errorf("Expected 1 wave, got %d", wm.GetWaveCount())
	}
	if wm.GetCurrentWaveIndex() != 0 {
		t.Errorf("Expected wave index 0, got %d", wm.GetCurrentWaveIndex())
	}
}

func TestGyrussWaveManager_Reset(t *testing.T) {
	world := donburi.NewWorld()
	logger := &mockLogger{}
	wm := NewGyrussWaveManager(world, logger)

	config := &managers.StageConfig{
		StageNumber: 1,
		Waves: []managers.GyrussWave{
			{WaveID: "wave_1"},
			{WaveID: "wave_2"},
		},
	}
	wm.LoadStage(config)

	// Advance state
	wm.currentWaveIndex = 1
	wm.bossTriggered = true

	wm.Reset()

	if wm.GetCurrentWaveIndex() != 0 {
		t.Error("Expected wave index to reset to 0")
	}
	if wm.IsBossTriggered() {
		t.Error("Expected boss triggered to reset to false")
	}
}

func TestGyrussWaveManager_HasMoreWaves(t *testing.T) {
	world := donburi.NewWorld()
	logger := &mockLogger{}
	wm := NewGyrussWaveManager(world, logger)

	// No stage loaded
	if wm.HasMoreWaves() {
		t.Error("Expected no more waves when no stage loaded")
	}

	config := &managers.StageConfig{
		Waves: []managers.GyrussWave{
			{WaveID: "wave_1"},
			{WaveID: "wave_2"},
		},
	}
	wm.LoadStage(config)

	if !wm.HasMoreWaves() {
		t.Error("Expected more waves at start")
	}

	wm.currentWaveIndex = 2
	if wm.HasMoreWaves() {
		t.Error("Expected no more waves at end")
	}
}

func TestGyrussWaveManager_IsBossTriggered(t *testing.T) {
	world := donburi.NewWorld()
	logger := &mockLogger{}
	wm := NewGyrussWaveManager(world, logger)

	if wm.IsBossTriggered() {
		t.Error("Expected boss not triggered initially")
	}

	wm.bossTriggered = true
	if !wm.IsBossTriggered() {
		t.Error("Expected boss to be triggered")
	}
}

func TestGyrussWaveManager_GetBossConfig(t *testing.T) {
	world := donburi.NewWorld()
	logger := &mockLogger{}
	wm := NewGyrussWaveManager(world, logger)

	// No stage loaded
	if wm.GetBossConfig() != nil {
		t.Error("Expected nil boss config when no stage loaded")
	}

	config := &managers.StageConfig{
		Boss: managers.StageBossConfig{
			Enabled:  true,
			BossType: "test_boss",
			Health:   100,
		},
	}
	wm.LoadStage(config)

	bossConfig := wm.GetBossConfig()
	if bossConfig == nil {
		t.Fatal("Expected boss config to be returned")
	}
	if bossConfig.Health != 100 {
		t.Errorf("Expected boss health 100, got %d", bossConfig.Health)
	}
}

func TestGyrussWaveManager_IsWaitingForLevelStart(t *testing.T) {
	world := donburi.NewWorld()
	logger := &mockLogger{}
	wm := NewGyrussWaveManager(world, logger)

	config := &managers.StageConfig{
		Waves: []managers.GyrussWave{{WaveID: "wave_1"}},
	}
	wm.LoadStage(config)

	if !wm.IsWaitingForLevelStart() {
		t.Error("Expected to be waiting for level start after loading stage")
	}
}
