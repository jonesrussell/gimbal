//nolint:testpackage // White box tests need access to internal fields
package gyruss

import (
	"context"
	"testing"

	"github.com/yohamta/donburi"

	"github.com/jonesrussell/gimbal/assets"
	"github.com/jonesrussell/gimbal/internal/config"
	"github.com/jonesrussell/gimbal/internal/ecs/core"
	resources "github.com/jonesrussell/gimbal/internal/ecs/managers/resource"
	"github.com/jonesrussell/gimbal/internal/ecs/systems/enemy"
)

func createTestGyrussSystem(t *testing.T) *GyrussSystem {
	t.Helper()
	world := donburi.NewWorld()
	gameConfig := &config.GameConfig{
		ScreenSize: config.Size{Width: 800, Height: 600},
	}
	ctx := context.Background()
	resourceMgr := resources.NewResourceManager(ctx)

	return NewGyrussSystem(&GyrussSystemConfig{
		World:       world,
		GameConfig:  gameConfig,
		ResourceMgr: resourceMgr,
		AssetsFS:    assets.Assets,
		EventSystem: nil,
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

func TestGyrussSystem_Reset(t *testing.T) {
	gs := createTestGyrussSystem(t)
	if err := gs.LoadStage(1); err != nil {
		t.Fatalf("LoadStage(1): %v", err)
	}
	gs.Reset()
	// Wave manager should be reset (no direct boss state to clear)
	if gs.GetWaveManager().GetStageConfig() == nil {
		t.Error("Reset should not clear loaded stage config")
	}
}

func TestGyrussSystem_SpawnBoss_CreatesBossWhenEnabled(t *testing.T) {
	gs := createTestGyrussSystem(t)
	if err := gs.LoadStage(1); err != nil {
		t.Fatalf("LoadStage(1): %v", err)
	}
	ctx := context.Background()
	gs.SpawnBoss(ctx)
	// Stage 1 has boss enabled; should have created a boss entity
	entries := core.GetEnemyEntries(gs.world)
	var foundBoss bool
	for _, e := range entries {
		if e.HasComponent(core.EnemyTypeID) {
			typeID := core.EnemyTypeID.Get(e)
			if enemy.EnemyType(*typeID) == enemy.EnemyTypeBoss {
				foundBoss = true
				break
			}
		}
	}
	if !foundBoss {
		t.Error("Expected SpawnBoss to create a boss entity when stage has boss enabled")
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

func TestGyrussSystem_DestroyEnemy_RemovesEntity(t *testing.T) {
	gs := createTestGyrussSystem(t)

	entity := gs.world.Create(core.EnemyTag)
	if !gs.world.Valid(entity) {
		t.Fatal("Expected enemy entity to exist after Create")
	}

	points := gs.DestroyEnemy(entity)

	if gs.world.Valid(entity) {
		t.Error("Expected entity to be removed after DestroyEnemy")
	}
	if points <= 0 {
		t.Errorf("Expected positive points from DestroyEnemy, got %d", points)
	}
}
