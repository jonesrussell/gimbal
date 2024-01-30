package game

import (
	"math"

	"github.com/solarlune/resolv"
)

const (
	MinAngle       = -math.Pi
	MaxAngle       = 3 * math.Pi / 2
	AngleStep      = 0.05
	RotationOffset = math.Pi / 2
)

func (player *Player) calculateCoordinates(angle float64) (int, int) {
	x := center.X + int(radius*math.Cos(angle))
	y := center.Y - int(radius*math.Sin(angle)) - playerHeight/2
	return x, y
}

func (player *Player) calculatePosition() resolv.Vector {
	x, y := player.calculateCoordinates(player.viewAngle)
	return resolv.Vector{X: float64(x), Y: float64(y)}
}

func (player *Player) calculateAngle() float64 {
	dx := float64(center.X) - player.Object.Position.X
	dy := float64(center.Y) - player.Object.Position.Y
	return math.Atan2(dy, dx) + RotationOffset
}

func (player *Player) calculatePath() []resolv.Vector {
	var path []resolv.Vector
	for angle := MinAngle; angle < MaxAngle; angle += AngleStep {
		x, y := player.calculateCoordinates(angle)
		path = append(path, resolv.Vector{X: float64(x), Y: float64(y)})
	}
	return path
}
