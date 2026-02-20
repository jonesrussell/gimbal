package stage

import (
	"context"

	"github.com/yohamta/donburi"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/dbg"
	"github.com/jonesrussell/gimbal/internal/ecs/events"
	"github.com/jonesrussell/gimbal/internal/ecs/managers"
	"github.com/jonesrussell/gimbal/internal/ecs/systems/enemy"
)

// StageState represents the authoritative stage lifecycle state
type StageState int

const (
	StageStatePreWave StageState = iota
	StageStateWaveInProgress
	StageStateWaveCompleted
	StageStateBossSpawning
	StageStateBossActive
	StageStateBossDefeated
	StageStateStageCompleted
)

// GyrussSystemForStage is the minimal interface StageStateMachine needs from GyrussSystem
type GyrussSystemForStage interface {
	LoadStage(stageNumber int) error
	Reset()
	SpawnBoss(ctx context.Context)
}

// Config holds dependencies for creating a StageStateMachine
type Config struct {
	EventSystem  *events.EventSystem
	WaveManager  *enemy.GyrussWaveManager
	GyrussSystem GyrussSystemForStage
	Logger       common.Logger
}

// StageStateMachine owns stage progression, wave index, and boss lifecycle
type StageStateMachine struct {
	eventSystem    *events.EventSystem
	waveManager    *enemy.GyrussWaveManager
	gyrussSystem   GyrussSystemForStage
	logger         common.Logger
	state          StageState
	waveIndex      int
	stageNumber    int
	preWaveTimer   float64 // seconds
	preWaveDelay   float64 // seconds (level start or inter-wave)
	bossSpawnTimer float64
}

const levelStartDelaySec = 3.5

// NewStageStateMachine creates a new stage state machine
func NewStageStateMachine(cfg *Config) *StageStateMachine {
	ssm := &StageStateMachine{
		eventSystem:  cfg.EventSystem,
		waveManager:  cfg.WaveManager,
		gyrussSystem: cfg.GyrussSystem,
		logger:       cfg.Logger,
		state:        StageStatePreWave,
		preWaveDelay: levelStartDelaySec,
	}
	cfg.EventSystem.SubscribeToBossDefeated(ssm.onBossDefeated)
	return ssm
}

func (ssm *StageStateMachine) onBossDefeated(_ donburi.World, _ events.BossDefeatedEvent) {
	dbg.Log(dbg.Event, "StageStateMachine.onBossDefeated fired (state=%v)", ssm.state)
	// Accept both BossActive and BossSpawning so we never get stuck if the boss
	// is defeated while state is still BossSpawning (e.g. same-frame or ordering edge case).
	if ssm.state != StageStateBossActive && ssm.state != StageStateBossSpawning {
		return
	}
	old := ssm.state
	ssm.state = StageStateBossDefeated
	dbg.Log(dbg.State, "StageStateMachine: %v → %v", old, StageStateBossDefeated)
	ssm.eventSystem.EmitStageCompleted(ssm.stageNumber)
	ssm.state = StageStateStageCompleted
	dbg.Log(dbg.State, "StageStateMachine: %v → %v", StageStateBossDefeated, StageStateStageCompleted)
	ssm.logger.Debug("Stage completed", "stage", ssm.stageNumber)
}

// LoadStage delegates to GyrussSystem then resets state to PreWave
func (ssm *StageStateMachine) LoadStage(stageNumber int) error {
	if err := ssm.gyrussSystem.LoadStage(stageNumber); err != nil {
		return err
	}
	ssm.stageNumber = stageNumber
	ssm.waveIndex = 0
	ssm.state = StageStatePreWave
	ssm.preWaveTimer = 0
	ssm.preWaveDelay = levelStartDelaySec
	ssm.bossSpawnTimer = 0
	return nil
}

// LoadNextStage advances to the next stage and loads it
func (ssm *StageStateMachine) LoadNextStage() error {
	next := ssm.stageNumber + 1
	return ssm.LoadStage(next)
}

// Reset clears state for a new game
func (ssm *StageStateMachine) Reset() {
	ssm.gyrussSystem.Reset()
	ssm.state = StageStatePreWave
	ssm.waveIndex = 0
	ssm.preWaveTimer = 0
	ssm.preWaveDelay = levelStartDelaySec
	ssm.bossSpawnTimer = 0
}

// Update drives state transitions; call after GyrussSystem.Update so wave manager counts are current
func (ssm *StageStateMachine) Update(ctx context.Context, deltaTime float64) {
	select {
	case <-ctx.Done():
		return
	default:
	}
	ssm.logger.Debug("StageStateMachine.Update", "state", ssm.state)

	switch ssm.state {
	case StageStatePreWave:
		ssm.updatePreWave(deltaTime)
	case StageStateWaveInProgress:
		ssm.updateWaveInProgress()
	case StageStateWaveCompleted:
		ssm.updateWaveCompleted()
	case StageStateBossSpawning:
		ssm.updateBossSpawning(ctx, deltaTime)
	case StageStateBossActive, StageStateBossDefeated, StageStateStageCompleted:
		// No timers; BossDefeated handled by event subscriber
	}
}

func (ssm *StageStateMachine) updatePreWave(deltaTime float64) {
	ssm.preWaveTimer += deltaTime
	if ssm.preWaveTimer < ssm.preWaveDelay {
		return
	}
	ssm.waveManager.StartWave(ssm.waveIndex)
	ssm.eventSystem.EmitWaveStarted(ssm.waveIndex)
	ssm.state = StageStateWaveInProgress
	dbg.Log(dbg.State, "StageStateMachine: %v → %v", StageStatePreWave, StageStateWaveInProgress)
	ssm.logger.Debug("Wave started", "wave_index", ssm.waveIndex)
}

func (ssm *StageStateMachine) updateWaveInProgress() {
	if !ssm.waveManager.AllSpawnedForCurrentWave() || ssm.waveManager.ActiveEnemyCount() != 0 {
		return
	}
	ssm.eventSystem.EmitWaveCompleted(ssm.waveIndex)
	ssm.state = StageStateWaveCompleted
	dbg.Log(dbg.State, "StageStateMachine: %v → %v", StageStateWaveInProgress, StageStateWaveCompleted)
	ssm.logger.Debug("Wave completed", "wave_index", ssm.waveIndex)
}

func (ssm *StageStateMachine) updateWaveCompleted() {
	waveCount := ssm.waveManager.GetWaveCount()
	if ssm.waveIndex+1 < waveCount {
		ssm.waveIndex++
		interWaveDelay := ssm.getWaveInterWaveDelay(ssm.waveIndex)
		if interWaveDelay > 0 {
			ssm.preWaveTimer = 0
			ssm.preWaveDelay = interWaveDelay
			ssm.state = StageStatePreWave
			dbg.Log(dbg.State, "StageStateMachine: %v → %v", StageStateWaveCompleted, StageStatePreWave)
		} else {
			ssm.waveManager.StartWave(ssm.waveIndex)
			ssm.eventSystem.EmitWaveStarted(ssm.waveIndex)
			ssm.state = StageStateWaveInProgress
			dbg.Log(dbg.State, "StageStateMachine: %v → %v", StageStateWaveCompleted, StageStateWaveInProgress)
		}
		return
	}
	bossConfig := ssm.waveManager.GetBossConfig()
	if bossConfig != nil && bossConfig.Enabled {
		ssm.eventSystem.EmitBossSpawnRequested()
		ssm.bossSpawnTimer = 0
		ssm.state = StageStateBossSpawning
		dbg.Log(dbg.State, "StageStateMachine: %v → %v", StageStateWaveCompleted, StageStateBossSpawning)
	} else {
		ssm.state = StageStateStageCompleted
		dbg.Log(dbg.State, "StageStateMachine: %v → %v", StageStateWaveCompleted, StageStateStageCompleted)
		ssm.eventSystem.EmitStageCompleted(ssm.stageNumber)
	}
}

func (ssm *StageStateMachine) updateBossSpawning(ctx context.Context, deltaTime float64) {
	bossConfig := ssm.waveManager.GetBossConfig()
	if bossConfig == nil || !bossConfig.Enabled {
		ssm.state = StageStateStageCompleted
		dbg.Log(dbg.State, "StageStateMachine: %v → %v", StageStateBossSpawning, StageStateStageCompleted)
		ssm.eventSystem.EmitStageCompleted(ssm.stageNumber)
		return
	}
	ssm.bossSpawnTimer += deltaTime
	if ssm.bossSpawnTimer < bossConfig.SpawnDelay {
		return
	}
	ssm.gyrussSystem.SpawnBoss(ctx)
	ssm.eventSystem.EmitBossSpawned()
	ssm.state = StageStateBossActive
	dbg.Log(dbg.State, "StageStateMachine: %v → %v", StageStateBossSpawning, StageStateBossActive)
	ssm.logger.Debug("Boss spawned")
}

func (ssm *StageStateMachine) getWaveInterWaveDelay(index int) float64 {
	cfg := ssm.waveManager.GetStageConfig()
	if cfg == nil || index < 0 || index >= len(cfg.Waves) {
		return 0
	}
	return cfg.Waves[index].Timing.InterWaveDelay
}

// IsStageCompleted returns true when the stage is complete (boss defeated or no boss)
func (ssm *StageStateMachine) IsStageCompleted() bool {
	return ssm.state == StageStateStageCompleted
}

// State returns the current stage state
func (ssm *StageStateMachine) State() StageState {
	return ssm.state
}

// CurrentWaveIndex returns the current wave index
func (ssm *StageStateMachine) CurrentWaveIndex() int {
	return ssm.waveIndex
}

// StageNumber returns the current stage number
func (ssm *StageStateMachine) StageNumber() int {
	return ssm.stageNumber
}

// StageConfig returns the current stage config from the wave manager (for debug overlay)
func (ssm *StageStateMachine) StageConfig() *managers.StageConfig {
	return ssm.waveManager.GetStageConfig()
}
