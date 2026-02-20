//nolint:testpackage // White box tests need access to internal fields
package enemy

import (
	"testing"

	"github.com/yohamta/donburi"

	"github.com/jonesrussell/gimbal/internal/ecs/core"
	"github.com/jonesrussell/gimbal/internal/ecs/managers"
)

func TestNewGyrussWaveManager(t *testing.T) {
	world := donburi.NewWorld()
	wm := NewGyrussWaveManager(world)

	if wm == nil {
		t.Fatal("Expected wave manager to be created")
	}
	if wm.world != world {
		t.Error("Expected world to be set")
	}
}

func TestGyrussWaveManager_LoadStage(t *testing.T) {
	world := donburi.NewWorld()
	wm := NewGyrussWaveManager(world)

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
	wm := NewGyrussWaveManager(world)

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
	wm.isSpawning = true

	wm.Reset()

	if wm.GetCurrentWaveIndex() != 0 {
		t.Error("Expected wave index to reset to 0")
	}
	if wm.isSpawning {
		t.Error("Expected isSpawning to reset to false")
	}
}

func TestGyrussWaveManager_HasMoreWaves(t *testing.T) {
	world := donburi.NewWorld()
	wm := NewGyrussWaveManager(world)

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

func TestGyrussWaveManager_GetBossConfig(t *testing.T) {
	world := donburi.NewWorld()
	wm := NewGyrussWaveManager(world)

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

func TestGyrussWaveManager_StartWave(t *testing.T) {
	world := donburi.NewWorld()
	wm := NewGyrussWaveManager(world)

	config := &managers.StageConfig{
		StageNumber: 1,
		Waves: []managers.GyrussWave{
			{
				WaveID: "wave_1",
				SpawnSequence: []managers.EnemyGroupConfig{
					{EnemyType: "basic", Count: 2, SpawnDelay: 0, SpawnInterval: 0},
				},
			},
			{WaveID: "wave_2"},
		},
	}
	wm.LoadStage(config)

	// Before StartWave, ShouldSpawnEnemy returns false (not spawning)
	if _, ok := wm.ShouldSpawnEnemy(); ok {
		t.Error("Expected ShouldSpawnEnemy false before StartWave")
	}

	wm.StartWave(0)
	if wm.GetCurrentWaveIndex() != 0 {
		t.Errorf("Expected current wave index 0, got %d", wm.GetCurrentWaveIndex())
	}
	// After StartWave(0), ShouldSpawnEnemy should return true (spawn delay 0)
	_, ok := wm.ShouldSpawnEnemy()
	if !ok {
		t.Error("Expected ShouldSpawnEnemy true after StartWave(0) with SpawnDelay 0")
	}

	// StartWave(1) switches to wave 1
	wm.StartWave(1)
	if wm.GetCurrentWaveIndex() != 1 {
		t.Errorf("Expected current wave index 1, got %d", wm.GetCurrentWaveIndex())
	}
}

func TestGyrussWaveManager_AllSpawnedForCurrentWave(t *testing.T) {
	world := donburi.NewWorld()
	wm := NewGyrussWaveManager(world)

	config := &managers.StageConfig{
		Waves: []managers.GyrussWave{
			{
				WaveID: "wave_1",
				SpawnSequence: []managers.EnemyGroupConfig{
					{EnemyType: "basic", Count: 2, SpawnDelay: 0, SpawnInterval: 0},
				},
			},
		},
	}
	wm.LoadStage(config)
	wm.StartWave(0)

	if wm.AllSpawnedForCurrentWave() {
		t.Error("Expected AllSpawnedForCurrentWave false before any spawns")
	}
	// Exhaust the single group: ShouldSpawnEnemy + MarkEnemySpawned for each of 2 enemies, then one more ShouldSpawnEnemy to advance past group
	for i := 0; i < 2; i++ {
		_, ok := wm.ShouldSpawnEnemy()
		if !ok {
			t.Fatalf("Expected ShouldSpawnEnemy true for spawn %d", i)
		}
		wm.MarkEnemySpawned()
	}
	_, _ = wm.ShouldSpawnEnemy() // advances currentGroupIndex past the only group
	if !wm.AllSpawnedForCurrentWave() {
		t.Error("Expected AllSpawnedForCurrentWave true after all groups spawned")
	}
}

func TestGyrussWaveManager_ActiveEnemyCount(t *testing.T) {
	world := donburi.NewWorld()
	wm := NewGyrussWaveManager(world)

	config := &managers.StageConfig{Waves: []managers.GyrussWave{{WaveID: "wave_1"}}}
	wm.LoadStage(config)

	if n := wm.ActiveEnemyCount(); n != 0 {
		t.Errorf("Expected ActiveEnemyCount 0 in empty world, got %d", n)
	}
	e1 := world.Create(core.EnemyTag)
	e2 := world.Create(core.EnemyTag)
	if n := wm.ActiveEnemyCount(); n != 2 {
		t.Errorf("Expected ActiveEnemyCount 2, got %d", n)
	}
	world.Remove(e1)
	if n := wm.ActiveEnemyCount(); n != 1 {
		t.Errorf("Expected ActiveEnemyCount 1 after remove, got %d", n)
	}
	world.Remove(e2)
	if n := wm.ActiveEnemyCount(); n != 0 {
		t.Errorf("Expected ActiveEnemyCount 0 after all removed, got %d", n)
	}
}
