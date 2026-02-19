package managers

import (
	"context"
	"testing"

	"github.com/jonesrussell/gimbal/assets"
	"github.com/jonesrussell/gimbal/internal/common"
)

// testLogger satisfies common.Logger for tests
type testLogger struct{}

func (*testLogger) Debug(msg string, fields ...any)   {}
func (*testLogger) Info(msg string, fields ...any)   {}
func (*testLogger) Warn(msg string, fields ...any)   {}
func (*testLogger) Error(msg string, fields ...any)   {}
func (*testLogger) DebugContext(_ context.Context, _ string, _ ...any) {}
func (*testLogger) InfoContext(_ context.Context, _ string, _ ...any)  {}
func (*testLogger) WarnContext(_ context.Context, _ string, _ ...any)  {}
func (*testLogger) ErrorContext(_ context.Context, _ string, _ ...any) {}
func (*testLogger) Sync() error                       { return nil }

func TestStageLoader_LoadStage_AssetValidation(t *testing.T) {
	var logger common.Logger = &testLogger{}
	loader := NewStageLoader(logger, assets.Assets)
	total := loader.GetTotalStages()

	for stageNum := 1; stageNum <= total; stageNum++ {
		config, err := loader.LoadStage(stageNum)
		if err != nil {
			t.Fatalf("stage %d: LoadStage failed: %v", stageNum, err)
		}
		if config == nil {
			t.Fatalf("stage %d: config is nil", stageNum)
		}
		if len(config.Waves) == 0 {
			t.Errorf("stage %d: len(config.Waves) = 0, want > 0", stageNum)
		}
		for i, wave := range config.Waves {
			if len(wave.SpawnSequence) == 0 {
				t.Errorf("stage %d wave %d (wave_id=%q): len(SpawnSequence) = 0, want > 0", stageNum, i, wave.WaveID)
			}
			for j, group := range wave.SpawnSequence {
				if group.EnemyType == "" {
					t.Errorf("stage %d wave %d group %d: EnemyType is empty", stageNum, i, j)
				}
				if group.Count <= 0 {
					t.Errorf("stage %d wave %d group %d: Count = %d, want > 0", stageNum, i, j, group.Count)
				}
			}
		}
		if config.Boss.Enabled {
			if config.Boss.BossType == "" {
				t.Errorf("stage %d: Boss.Enabled=true but Boss.BossType is empty", stageNum)
			}
			if config.Boss.Health <= 0 {
				t.Errorf("stage %d: Boss.Enabled=true but Boss.Health = %d, want > 0", stageNum, config.Boss.Health)
			}
		}
	}
}
