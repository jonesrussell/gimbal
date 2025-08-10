package collision

import (
	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/config"
)

// checkCollision checks if two entities are colliding using AABB collision detection
// Entities are positioned at their center, so we need to calculate the bounding box from the center
func (cs *CollisionSystem) checkCollision(
	pos1 common.Point, size1 config.Size,
	pos2 common.Point, size2 config.Size,
) bool {
	// Calculate bounding boxes from center position
	left1 := pos1.X - float64(size1.Width)/2
	right1 := pos1.X + float64(size1.Width)/2
	top1 := pos1.Y - float64(size1.Height)/2
	bottom1 := pos1.Y + float64(size1.Height)/2

	left2 := pos2.X - float64(size2.Width)/2
	right2 := pos2.X + float64(size2.Width)/2
	top2 := pos2.Y - float64(size2.Height)/2
	bottom2 := pos2.Y + float64(size2.Height)/2

	// Check for overlap
	return left1 < right2 && right1 > left2 && top1 < bottom2 && bottom1 > top2
}
