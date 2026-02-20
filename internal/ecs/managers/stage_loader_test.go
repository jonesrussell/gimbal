package managers_test

import (
	"context"
	"testing"

	"github.com/jonesrussell/gimbal/assets"
	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/ecs/managers"
)

type noopLogger struct{}

func (noopLogger) Debug(string, ...any)                         {}
func (noopLogger) Info(string, ...any)                          {}
func (noopLogger) Warn(string, ...any)                          {}
func (noopLogger) Error(string, ...any)                         {}
func (noopLogger) DebugContext(context.Context, string, ...any) {}
func (noopLogger) InfoContext(context.Context, string, ...any)  {}
func (noopLogger) WarnContext(context.Context, string, ...any)  {}
func (noopLogger) ErrorContext(context.Context, string, ...any) {}
func (noopLogger) Sync() error                                  { return nil }

func TestStageLoader_LoadStage_AssetValidation(t *testing.T) {
	var logger common.Logger = noopLogger{}
	loader := managers.NewStageLoader(logger, assets.Assets)
	total := loader.GetTotalStages()

	for stageNum := 1; stageNum <= total; stageNum++ {
		config, err := loader.LoadStage(stageNum)
		if err != nil {
			t.Fatalf("stage %d: LoadStage failed: %v", stageNum, err)
		}
		if config == nil {
			t.Fatalf("stage %d: config is nil", stageNum)
		}
		validateWaves(t, stageNum, config.Waves)
		validateBoss(t, stageNum, &config.Boss)
	}
}

func validateWaves(t *testing.T, stageNum int, waves []managers.GyrussWave) {
	t.Helper()
	if len(waves) == 0 {
		t.Errorf("stage %d: len(config.Waves) = 0, want > 0", stageNum)
		return
	}
	for i, wave := range waves {
		validateSpawnSequence(t, stageNum, i, wave.WaveID, wave.SpawnSequence)
	}
}

func validateSpawnSequence(t *testing.T, stageNum, waveIdx int, waveID string, seq []managers.EnemyGroupConfig) {
	t.Helper()
	if len(seq) == 0 {
		t.Errorf("stage %d wave %d (wave_id=%q): len(SpawnSequence) = 0, want > 0", stageNum, waveIdx, waveID)
		return
	}
	for j, group := range seq {
		if group.EnemyType == "" {
			t.Errorf("stage %d wave %d group %d: EnemyType is empty", stageNum, waveIdx, j)
		}
		if group.Count <= 0 {
			t.Errorf("stage %d wave %d group %d: Count = %d, want > 0", stageNum, waveIdx, j, group.Count)
		}
	}
}

func validateBoss(t *testing.T, stageNum int, boss *managers.StageBossConfig) {
	t.Helper()
	if !boss.Enabled {
		return
	}
	if boss.BossType == "" {
		t.Errorf("stage %d: Boss.Enabled=true but Boss.BossType is empty", stageNum)
	}
	if boss.Health <= 0 {
		t.Errorf("stage %d: Boss.Enabled=true but Boss.Health = %d, want > 0", stageNum, boss.Health)
	}
}
