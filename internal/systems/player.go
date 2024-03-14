package systems

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jonesrussell/gimbal/internal/components"
	"github.com/jonesrussell/gimbal/internal/config"
	"github.com/jonesrussell/gimbal/internal/logger"
	dresolv "github.com/jonesrussell/gimbal/internal/resolv"
	"github.com/jonesrussell/gimbal/internal/tags"
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

const AngleStep = 0.05

const (
	playerWidth    = 16
	playerHeight   = 16
	RotationOffset = math.Pi / 2
)

var radius = float64(config.C.Height/2) * 0.75
var center = image.Point{X: config.C.Width / 2, Y: config.C.Height / 2}

func UpdatePlayer(ecs *ecs.ECS) {
	// Now we update the Player's movement. This is the real bread-and-butter of this example, naturally.
	playerEntry, _ := components.Player.First(ecs.World)
	player := components.Player.Get(playerEntry)
	playerObject := dresolv.GetObject(playerEntry)

	oldOrientation := player.ViewAngle
	oldDirection := player.Direction
	oldAngle := player.Angle
	oldX := playerObject.Position.X
	oldY := playerObject.Position.Y

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		player.Direction = -1
		player.ViewAngle -= AngleStep
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		player.Direction = 1
		player.ViewAngle += AngleStep
	} else {
		player.Direction = 0
	}

	position := calculatePosition(player)
	logger.GlobalLogger.Info("position", "full", position)

	player.Angle = calculateAngle(playerObject)

	if player.ViewAngle != oldOrientation || player.Direction != oldDirection || player.Angle != oldAngle || playerObject.Position.X != oldX || playerObject.Position.Y != oldY {
		logger.GlobalLogger.Debug("Player", "viewAngle", player.ViewAngle, "direction", player.Direction, "angle", player.Angle, "X", float64(playerObject.Position.X), "Y", float64(playerObject.Position.Y))
	}

	// Add the current position to the path
	//player.Path = append(player.Path, player.Object.Position)
}

func DrawPlayer(ecs *ecs.ECS, screen *ebiten.Image) {
	tags.Player.Each(ecs.World, func(e *donburi.Entry) {
		o := dresolv.GetObject(e)
		player := components.Player.Get(e) // get the PlayerData

		if player.Sprite != nil {
			op := &ebiten.DrawImageOptions{}
			// Scale the sprite down to 10% of its original size
			op.GeoM.Scale(0.1, 0.1)
			op.GeoM.Translate(float64(o.Position.X), float64(o.Position.Y))
			screen.DrawImage(player.Sprite, op)
		}
	})
}

func calculateCoordinates(angle float64) (int, int) {
	x := center.X + int(radius*math.Cos(angle))
	y := center.Y - int(radius*math.Sin(angle)) - playerHeight/2
	return x, y
}

func calculatePosition(data *components.PlayerData) resolv.Vector {
	x, y := calculateCoordinates(data.ViewAngle)
	return resolv.Vector{X: float64(x), Y: float64(y)}
}

func calculateAngle(playerObject *resolv.Object) float64 {
	dx := float64(center.X) - playerObject.Position.X
	dy := float64(center.Y) - playerObject.Position.Y
	return math.Atan2(dy, dx) + RotationOffset
}