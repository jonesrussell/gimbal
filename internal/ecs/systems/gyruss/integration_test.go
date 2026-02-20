//go:build integration

package gyruss

import (
	"context"
	"testing"
	"time"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/jonesrussell/gimbal/assets"
	"github.com/jonesrussell/gimbal/internal/config"
	"github.com/jonesrussell/gimbal/internal/ecs/core"
	"github.com/jonesrussell/gimbal/internal/ecs/events"
	resources "github.com/jonesrussell/gimbal/internal/ecs/managers/resource"
	"github.com/jonesrussell/gimbal/internal/ecs/systems/stage"
)

func TestStage1Completes(t *testing.T) {
	world := donburi.NewWorld()
	eventSystem := events.NewEventSystem(world)
	gs := createTestGyrussSystemWithEventSystem(t, world, eventSystem)
	ssm := stage.NewStageStateMachine(&stage.Config{
		EventSystem:  eventSystem,
		WaveManager:  gs.GetWaveManager(),
		GyrussSystem: gs,
		Logger:       &testLogger{},
	})
	if err := ssm.LoadStage(1); err != nil {
		t.Fatalf("LoadStage(1): %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	dt := 1.0 / 60.0
	enemyQuery := query.NewQuery(filter.Contains(core.EnemyTag))

	for {
		select {
		case <-ctx.Done():
			if ssm.IsStageCompleted() {
				return
			}
			t.Fatalf("timeout before stage complete: %v", ctx.Err())
		default:
		}

		if err := gs.Update(ctx, dt); err != nil {
			if ctx.Err() != nil {
				break
			}
			t.Fatalf("Update: %v", err)
		}
		ssm.Update(ctx, dt)
		eventSystem.ProcessEvents()

		// Collect all enemy entities (avoid modifying world during Each)
		var toDestroy []donburi.Entity
		enemyQuery.Each(gs.world, func(entry *donburi.Entry) {
			toDestroy = append(toDestroy, entry.Entity())
		})
		for _, e := range toDestroy {
			gs.DestroyEnemy(e)
		}
		eventSystem.ProcessEvents()

		if ssm.IsStageCompleted() {
			return
		}
	}
}

// createTestGyrussSystemWithEventSystem creates a GyrussSystem with EventSystem for integration tests
func createTestGyrussSystemWithEventSystem(t *testing.T, world donburi.World, eventSystem *events.EventSystem) *GyrussSystem {
	t.Helper()
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
		EventSystem: eventSystem,
	})
}
