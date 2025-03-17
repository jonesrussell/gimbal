package game

import (
	"fmt"
	"math"

	"github.com/solarlune/resolv"
)

const (
	MinAngle            = -math.Pi
	MaxAngle            = 3 * math.Pi / 2
	RotationOffset      = math.Pi / 2
	playerHeightDivisor = 2
	degreesToRadians    = math.Pi / 180 // Conversion factor from degrees to radians
	startAngleDegrees   = 270           // Starting angle in degrees (bottom of circle)
)

func (player *Player) CalculateCoordinates(angle float64) (float64, float64) {
	centerX := float64(player.Config.Center.X)
	centerY := float64(player.Config.Center.Y)
	radius := player.Config.Radius

	fmt.Printf("Debug - Input values:\n")
	fmt.Printf("  angle: %.2f degrees\n", angle)
	fmt.Printf("  centerX: %.2f, centerY: %.2f\n", centerX, centerY)
	fmt.Printf("  radius: %.2f\n", radius)

	// Convert angle to radians
	// In this system:
	// - 0° points right
	// - Positive angles go counterclockwise
	angleInRadians := angle * degreesToRadians
	fmt.Printf("Debug - Angle calculations:\n")
	fmt.Printf("  angleInRadians: %.4f radians (%.2f degrees)\n", angleInRadians, angle)

	// Calculate raw coordinates before rounding
	rawX := centerX + radius*math.Cos(angleInRadians)
	rawY := centerY - radius*math.Sin(angleInRadians) - float64(player.Config.PlayerHeight)/playerHeightDivisor

	fmt.Printf("Debug - Coordinate calculations:\n")
	fmt.Printf("  cos(angle): %.4f\n", math.Cos(angleInRadians))
	fmt.Printf("  sin(angle): %.4f\n", math.Sin(angleInRadians))
	fmt.Printf("  playerHeight adjustment: %.2f\n", float64(player.Config.PlayerHeight)/playerHeightDivisor)
	fmt.Printf("  Raw coordinates: (%.2f, %.2f)\n", rawX, rawY)

	// Round to nearest integer
	finalX := math.Round(rawX)
	finalY := math.Round(rawY)
	fmt.Printf("Debug - Final rounded coordinates: (%.0f, %.0f)\n", finalX, finalY)

	return finalX, finalY
}

func (player *Player) CalculatePosition() resolv.Vector {
	x, y := player.CalculateCoordinates(player.ViewAngle)
	return resolv.Vector{X: x, Y: y}
}

func (player *Player) CalculateAngle() float64 {
	pos := player.Object.Position()
	centerX := float64(player.Config.Center.X)
	centerY := float64(player.Config.Center.Y)
	dx := centerX - pos.X
	dy := centerY - pos.Y
	return math.Atan2(dy, dx) + RotationOffset
}
