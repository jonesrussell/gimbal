package stage_test

import (
	"context"
	"testing"

	"github.com/yohamta/donburi"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/ecs/events"
	"github.com/jonesrussell/gimbal/internal/ecs/systems/enemy"
	"github.com/jonesrussell/gimbal/internal/ecs/systems/stage"
)

type testLogger struct{}

func (testLogger) Debug(msg string, keysAndValues ...any)                             {}
func (testLogger) Info(msg string, keysAndValues ...any)                              {}
func (testLogger) Warn(msg string, keysAndValues ...any)                              {}
func (testLogger) Error(msg string, keysAndValues ...any)                             {}
func (testLogger) DebugContext(ctx context.Context, msg string, keysAndValues ...any) {}
func (testLogger) InfoContext(ctx context.Context, msg string, keysAndValues ...any)  {}
func (testLogger) WarnContext(ctx context.Context, msg string, keysAndValues ...any)  {}
func (testLogger) ErrorContext(ctx context.Context, msg string, keysAndValues ...any) {}
func (testLogger) Sync() error                                                        { return nil }

func TestStageStateMachine_InitialState(t *testing.T) {
	world := donburi.NewWorld()
	evt := events.NewEventSystem(world)
	wm := enemy.NewGyrussWaveManager(world, testLogger{})
	gs := &mockGyrussSystem{}
	ssm := stage.NewStageStateMachine(&stage.Config{
		EventSystem:  evt,
		WaveManager:  wm,
		GyrussSystem: gs,
		Logger:       testLogger{},
	})
	if ssm.State() != stage.StageStatePreWave {
		t.Errorf("initial state want PreWave, got %v", ssm.State())
	}
	if ssm.IsStageCompleted() {
		t.Error("initial IsStageCompleted should be false")
	}
	if ssm.CurrentWaveIndex() != 0 {
		t.Errorf("initial wave index want 0, got %d", ssm.CurrentWaveIndex())
	}
}

func TestStageStateMachine_IsStageCompleted_OnlyWhenStageCompleted(t *testing.T) {
	world := donburi.NewWorld()
	evt := events.NewEventSystem(world)
	wm := enemy.NewGyrussWaveManager(world, testLogger{})
	gs := &mockGyrussSystem{}
	ssm := stage.NewStageStateMachine(&stage.Config{
		EventSystem:  evt,
		WaveManager:  wm,
		GyrussSystem: gs,
		Logger:       testLogger{},
	})
	if ssm.IsStageCompleted() {
		t.Error("PreWave should not be stage completed")
	}
	// Manually transition to BossActive and then trigger BossDefeated via event
	// (we can't easily do full flow without loading a stage)
	evt.EmitBossDefeated()
	evt.ProcessEvents()
	// Handler only runs when state is BossActive, so still not completed
	if ssm.IsStageCompleted() {
		t.Error("after BossDefeated with state PreWave, should not be stage completed")
	}
}

type mockGyrussSystem struct {
	loadStageErr   error
	resetCalls     int
	spawnBossCalls int
}

func (m *mockGyrussSystem) LoadStage(stageNumber int) error {
	return m.loadStageErr
}

func (m *mockGyrussSystem) Reset() {
	m.resetCalls++
}

func (m *mockGyrussSystem) SpawnBoss(ctx context.Context) {
	m.spawnBossCalls++
}

var _ stage.GyrussSystemForStage = (*mockGyrussSystem)(nil)

func TestStageStateMachine_LoadStage_ResetsState(t *testing.T) {
	world := donburi.NewWorld()
	evt := events.NewEventSystem(world)
	wm := enemy.NewGyrussWaveManager(world, testLogger{})
	gs := &mockGyrussSystem{}
	ssm := stage.NewStageStateMachine(&stage.Config{
		EventSystem:  evt,
		WaveManager:  wm,
		GyrussSystem: gs,
		Logger:       testLogger{},
	})
	if err := ssm.LoadStage(1); err != nil {
		t.Fatalf("LoadStage(1): %v", err)
	}
	if ssm.State() != stage.StageStatePreWave {
		t.Errorf("after LoadStage state want PreWave, got %v", ssm.State())
	}
	if ssm.StageNumber() != 1 {
		t.Errorf("StageNumber want 1, got %d", ssm.StageNumber())
	}
	if ssm.CurrentWaveIndex() != 0 {
		t.Errorf("wave index want 0, got %d", ssm.CurrentWaveIndex())
	}
}

func TestStageStateMachine_Reset_CallsGyrussReset(t *testing.T) {
	world := donburi.NewWorld()
	evt := events.NewEventSystem(world)
	wm := enemy.NewGyrussWaveManager(world, testLogger{})
	gs := &mockGyrussSystem{}
	ssm := stage.NewStageStateMachine(&stage.Config{
		EventSystem:  evt,
		WaveManager:  wm,
		GyrussSystem: gs,
		Logger:       testLogger{},
	})
	ssm.Reset()
	if gs.resetCalls != 1 {
		t.Errorf("Reset should call GyrussSystem.Reset once, got %d", gs.resetCalls)
	}
}

// common.Logger compatibility
var _ common.Logger = testLogger{}
