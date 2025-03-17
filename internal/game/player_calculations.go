package game

import (
	"math"

	"github.com/solarlune/resolv"
)

const (
	MinAngle       = -math.Pi
	MaxAngle       = 3 * math.Pi / 2
	RotationOffset = math.Pi / 2
)

func (player *Player) calculateCoordinates(angle float64) (float64, float64) {
	centerX := float64(player.config.ScreenWidth) / 2
	centerY := float64(player.config.ScreenHeight) / 2
	radius := float64(player.config.ScreenHeight) / 4

	x := centerX + radius*math.Cos(angle)
	y := centerY - radius*math.Sin(angle) - float64(player.config.PlayerHeight)/2
	return x, y
}

func (player *Player) calculatePosition() resolv.Vector {
	x, y := player.calculateCoordinates(player.viewAngle)
	return resolv.Vector{X: x, Y: y}
}

func (player *Player) calculateAngle() float64 {
	pos := player.Object.Position()
	centerX := float64(player.config.ScreenWidth) / 2
	centerY := float64(player.config.ScreenHeight) / 2
	dx := centerX - pos.X
	dy := centerY - pos.Y
	return math.Atan2(dy, dx) + RotationOffset
}
