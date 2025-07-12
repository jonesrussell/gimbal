package collision

import (
	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/config"
)

// checkCollision checks if two entities are colliding using AABB collision detection
func (cs *CollisionSystem) checkCollision(
	pos1 common.Point, size1 config.Size,
	pos2 common.Point, size2 config.Size,
) bool {
	// Calculate bounding boxes
	left1 := pos1.X
	right1 := pos1.X + float64(size1.Width)
	top1 := pos1.Y
	bottom1 := pos1.Y + float64(size1.Height)

	left2 := pos2.X
	right2 := pos2.X + float64(size2.Width)
	top2 := pos2.Y
	bottom2 := pos2.Y + float64(size2.Height)

	// Check for overlap
	return left1 < right2 && right1 > left2 && top1 < bottom2 && bottom1 > top2
}
