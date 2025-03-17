package game

import (
	"math"

	"github.com/solarlune/resolv"
)

const (
	MinAngle            = -math.Pi
	MaxAngle            = 3 * math.Pi / 2
	RotationOffset      = math.Pi / 2
	playerHeightDivisor = 2
)

func (player *Player) CalculateCoordinates(angle float64) (float64, float64) {
	centerX := float64(player.Config.ScreenWidth) / screenCenterDivisor
	centerY := float64(player.Config.ScreenHeight) / screenCenterDivisor
	radius := float64(player.Config.ScreenHeight) / radiusDivisor

	x := centerX + radius*math.Cos(angle)
	y := centerY - radius*math.Sin(angle) - float64(player.Config.PlayerHeight)/playerHeightDivisor
	return x, y
}

func (player *Player) CalculatePosition() resolv.Vector {
	x, y := player.CalculateCoordinates(player.ViewAngle)
	return resolv.Vector{X: x, Y: y}
}

func (player *Player) CalculateAngle() float64 {
	pos := player.Object.Position()
	centerX := float64(player.Config.ScreenWidth) / screenCenterDivisor
	centerY := float64(player.Config.ScreenHeight) / screenCenterDivisor
	dx := centerX - pos.X
	dy := centerY - pos.Y
	return math.Atan2(dy, dx) + RotationOffset
}
