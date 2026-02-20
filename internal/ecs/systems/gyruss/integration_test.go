//go:build integration

package gyruss

import (
	"context"
	"testing"
	"time"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/jonesrussell/gimbal/internal/ecs/core"
)

func TestStage1Completes(t *testing.T) {
	gs := createTestGyrussSystem(t)
	if err := gs.LoadStage(1); err != nil {
		t.Fatalf("LoadStage(1): %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	dt := 1.0 / 60.0
	enemyQuery := query.NewQuery(filter.Contains(core.EnemyTag))

	for {
		select {
		case <-ctx.Done():
			if gs.IsStageComplete() {
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

		// Collect all enemy entities (avoid modifying world during Each)
		var toDestroy []donburi.Entity
		enemyQuery.Each(gs.world, func(entry *donburi.Entry) {
			toDestroy = append(toDestroy, entry.Entity())
		})
		for _, e := range toDestroy {
			gs.DestroyEnemy(e)
		}

		if gs.IsStageComplete() {
			return
		}
	}
}
