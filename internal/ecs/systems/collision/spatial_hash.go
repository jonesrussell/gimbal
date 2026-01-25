package collision

import (
	"github.com/yohamta/donburi"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/config"
)

// SpatialHashCellSize is the size of each cell in the spatial hash grid.
// Chosen to balance granularity with overhead for typical game entity sizes.
const SpatialHashCellSize = 64

// SpatialHash provides O(1) spatial lookups for collision detection.
// Entities are stored in grid cells based on their position.
type SpatialHash struct {
	cellSize int
	cells    map[cellKey][]donburi.Entity
	width    int
	height   int
}

// cellKey uniquely identifies a cell in the grid
type cellKey struct {
	x, y int
}

// NewSpatialHash creates a new spatial hash for the given world dimensions.
func NewSpatialHash(worldWidth, worldHeight int) *SpatialHash {
	return &SpatialHash{
		cellSize: SpatialHashCellSize,
		cells:    make(map[cellKey][]donburi.Entity),
		width:    worldWidth,
		height:   worldHeight,
	}
}

// Clear removes all entities from the spatial hash.
func (sh *SpatialHash) Clear() {
	// Reuse the map but clear all slices
	for k := range sh.cells {
		delete(sh.cells, k)
	}
}

// Insert adds an entity to the spatial hash based on its position and size.
// Entities spanning multiple cells are added to all overlapping cells.
func (sh *SpatialHash) Insert(entity donburi.Entity, pos common.Point, size config.Size) {
	minX := int(pos.X) / sh.cellSize
	maxX := int(pos.X+float64(size.Width)) / sh.cellSize
	minY := int(pos.Y) / sh.cellSize
	maxY := int(pos.Y+float64(size.Height)) / sh.cellSize

	for x := minX; x <= maxX; x++ {
		for y := minY; y <= maxY; y++ {
			key := cellKey{x, y}
			sh.cells[key] = append(sh.cells[key], entity)
		}
	}
}

// Query returns all entities that could potentially collide with an entity
// at the given position and size. Returns entities from overlapping cells.
func (sh *SpatialHash) Query(pos common.Point, size config.Size) []donburi.Entity {
	minX := int(pos.X) / sh.cellSize
	maxX := int(pos.X+float64(size.Width)) / sh.cellSize
	minY := int(pos.Y) / sh.cellSize
	maxY := int(pos.Y+float64(size.Height)) / sh.cellSize

	// Use a map to deduplicate entities that span multiple cells
	seen := make(map[donburi.Entity]struct{})
	result := make([]donburi.Entity, 0)

	for x := minX; x <= maxX; x++ {
		for y := minY; y <= maxY; y++ {
			key := cellKey{x, y}
			for _, entity := range sh.cells[key] {
				if _, ok := seen[entity]; !ok {
					seen[entity] = struct{}{}
					result = append(result, entity)
				}
			}
		}
	}

	return result
}

// QueryNearby returns entities in the same cell and adjacent cells.
// Useful for broad-phase collision detection.
func (sh *SpatialHash) QueryNearby(pos common.Point) []donburi.Entity {
	centerX := int(pos.X) / sh.cellSize
	centerY := int(pos.Y) / sh.cellSize

	seen := make(map[donburi.Entity]struct{})
	result := make([]donburi.Entity, 0)

	// Check 3x3 grid of cells centered on entity position
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			key := cellKey{centerX + dx, centerY + dy}
			for _, entity := range sh.cells[key] {
				if _, ok := seen[entity]; !ok {
					seen[entity] = struct{}{}
					result = append(result, entity)
				}
			}
		}
	}

	return result
}
