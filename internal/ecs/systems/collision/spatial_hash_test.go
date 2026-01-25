package collision

import (
	"testing"

	"github.com/yohamta/donburi"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/config"
)

func TestSpatialHash_InsertAndQuery(t *testing.T) {
	sh := NewSpatialHash(640, 480)

	// Insert an entity at position (100, 100) with size 32x32
	entity1 := donburi.Entity(1)
	sh.Insert(entity1, common.Point{X: 100, Y: 100}, config.Size{Width: 32, Height: 32})

	// Query at a nearby position - should find entity1
	results := sh.Query(common.Point{X: 110, Y: 110}, config.Size{Width: 16, Height: 16})
	if len(results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results))
	}
	if len(results) > 0 && results[0] != entity1 {
		t.Errorf("Expected entity1, got %v", results[0])
	}

	// Query at a distant position - should find nothing
	results = sh.Query(common.Point{X: 500, Y: 400}, config.Size{Width: 16, Height: 16})
	if len(results) != 0 {
		t.Errorf("Expected 0 results, got %d", len(results))
	}
}

func TestSpatialHash_Clear(t *testing.T) {
	sh := NewSpatialHash(640, 480)

	// Insert an entity
	entity1 := donburi.Entity(1)
	sh.Insert(entity1, common.Point{X: 100, Y: 100}, config.Size{Width: 32, Height: 32})

	// Verify it's there
	results := sh.Query(common.Point{X: 100, Y: 100}, config.Size{Width: 32, Height: 32})
	if len(results) != 1 {
		t.Errorf("Expected 1 result before clear, got %d", len(results))
	}

	// Clear the hash
	sh.Clear()

	// Verify it's gone
	results = sh.Query(common.Point{X: 100, Y: 100}, config.Size{Width: 32, Height: 32})
	if len(results) != 0 {
		t.Errorf("Expected 0 results after clear, got %d", len(results))
	}
}

func TestSpatialHash_MultipleEntities(t *testing.T) {
	sh := NewSpatialHash(640, 480)

	// Insert multiple entities in the same cell
	entity1 := donburi.Entity(1)
	entity2 := donburi.Entity(2)
	entity3 := donburi.Entity(3)

	sh.Insert(entity1, common.Point{X: 10, Y: 10}, config.Size{Width: 16, Height: 16})
	sh.Insert(entity2, common.Point{X: 20, Y: 20}, config.Size{Width: 16, Height: 16})
	sh.Insert(entity3, common.Point{X: 500, Y: 400}, config.Size{Width: 16, Height: 16}) // Different cell

	// Query near first two entities
	results := sh.Query(common.Point{X: 15, Y: 15}, config.Size{Width: 16, Height: 16})
	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}
}

func TestSpatialHash_EntitySpanningCells(t *testing.T) {
	sh := NewSpatialHash(640, 480)

	// Insert a large entity that spans multiple cells
	entity1 := donburi.Entity(1)
	// At cell size 64, an entity at (60, 60) with size 100x100 would span cells (0,0), (0,1), (1,0), (1,1)
	sh.Insert(entity1, common.Point{X: 60, Y: 60}, config.Size{Width: 100, Height: 100})

	// Query in different cells - should find the same entity
	results1 := sh.Query(common.Point{X: 10, Y: 10}, config.Size{Width: 16, Height: 16})
	results2 := sh.Query(common.Point{X: 100, Y: 100}, config.Size{Width: 16, Height: 16})

	// Both queries should find entity1
	found1 := containsEntity(results1, entity1)
	found2 := containsEntity(results2, entity1)

	if !found1 {
		t.Errorf("Expected to find entity1 in cell (0,0)")
	}
	if !found2 {
		t.Errorf("Expected to find entity1 in cell (1,1)")
	}
}

func TestSpatialHash_QueryNearby(t *testing.T) {
	sh := NewSpatialHash(640, 480)

	// Insert entity in center
	entity1 := donburi.Entity(1)
	sh.Insert(entity1, common.Point{X: 100, Y: 100}, config.Size{Width: 16, Height: 16})

	// QueryNearby should find it
	results := sh.QueryNearby(common.Point{X: 100, Y: 100})
	if !containsEntity(results, entity1) {
		t.Errorf("QueryNearby should find entity at same position")
	}

	// Query from adjacent cell should also find it
	results = sh.QueryNearby(common.Point{X: 150, Y: 100}) // Next cell over
	if !containsEntity(results, entity1) {
		t.Errorf("QueryNearby should find entity in adjacent cell")
	}
}

func containsEntity(entities []donburi.Entity, target donburi.Entity) bool {
	for _, e := range entities {
		if e == target {
			return true
		}
	}
	return false
}
