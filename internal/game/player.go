package game

import (
	"errors"
	"image"
	"math"

	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
)

// Player represents a player in the game.
type Player struct {
	// Input is the input handler for the player.
	input InputHandlerInterface
	// Angle is the current angle of the player.
	angle float64
	// Speed is the speed of the player.
	speed float64
	// Direction is the current direction of the player.
	direction float64
	// Object is the game object representing the player.
	Object *resolv.Object
	// Sprite is the player's sprite.
	Sprite *ebiten.Image
	// Orientation of player's viewable sprite
	viewAngle float64
	// Debugger
	debug *Debugger
}

func NewPlayer(
	input InputHandlerInterface,
	speed float64,
	debugger *Debugger,
	spriteImage *ebiten.Image,
) (*Player, error) {
	if input == nil {
		return nil, errors.New("input handler cannot be nil")
	}

	if speed <= 0 {
		return nil, errors.New("speed must be greater than zero")
	}

	if spriteImage == nil {
		return nil, errors.New("sprite image cannot be nil")
	}

	player := &Player{
		input:     input,
		angle:     0,
		speed:     speed,
		direction: 0,
		Sprite:    spriteImage,
		viewAngle: 3 * math.Pi / 2,
		debug:     debugger,
	}

	// Calculate the initial position
	position := player.calculatePosition()
	player.Object = resolv.NewObject(
		position.X,
		position.Y,
		float64(playerWidth),
		float64(playerHeight),
	)

	return player, nil
}

func (player *Player) calculatePosition() resolv.Vector {
	// Calculate the x and y positions based on the current angle
	x := center.X + int(radius*math.Cos(player.viewAngle))
	y := center.Y - int(radius*math.Sin(player.viewAngle))

	return resolv.Vector{X: float64(x), Y: float64(y)}
}

func (player *Player) calculateAngle() float64 {
	// Calculate the angle between the sprite and the center of the screen
	dx := float64(center.X) - player.Object.Position.X
	dy := float64(center.Y) - player.Object.Position.Y
	return math.Atan2(dy, dx) + math.Pi/2 // Add Pi/2 to rotate the sprite by 90 degrees
}

func (player *Player) Update() {
	if !gameStarted {
		player.debug.DebugPlayer(player)
		gameStarted = true
	}

	oldOrientation := player.viewAngle
	oldDirection := player.direction
	oldAngle := player.angle
	oldX := player.Object.Position.X
	oldY := player.Object.Position.Y

	if player.input.IsKeyPressed(ebiten.KeyLeft) {
		player.direction = -1
		player.viewAngle -= 0.05
	} else if player.input.IsKeyPressed(ebiten.KeyRight) {
		player.direction = 1
		player.viewAngle += 0.05
	} else {
		player.direction = 0
	}

	player.Object.Position = player.calculatePosition()

	player.angle = player.calculateAngle()

	if player.viewAngle != oldOrientation {
		player.debug.DebugPrintOrientation(player.viewAngle)
	}

	if player.direction != oldDirection {
		player.debug.DebugPrintDirection(player.direction)
	}

	if player.angle != oldAngle {
		player.debug.DebugPrintAngle(player.angle)
	}

	if player.Object.Position.X != oldX || player.Object.Position.Y != oldY {
		pos := image.Point{X: int(player.Object.Position.X), Y: int(player.Object.Position.Y)}
		player.debug.DebugPrintPosition(pos)
	}

	player.Object.Update()
}

func (player *Player) Draw(screen *ebiten.Image) {
	player.updatePosition()
	player.drawSprite(screen)
}

func (player *Player) updatePosition() {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(player.Object.Position.X, player.Object.Position.Y)
	player.Object.Update()
}

func (player *Player) drawSprite(screen *ebiten.Image) {
	if player.Sprite != nil {
		spriteOp := player.createSpriteOptions()
		rotatedSprite := player.getRotatedSprite()
		screen.DrawImage(rotatedSprite, spriteOp)
	}
}

func (player *Player) createSpriteOptions() *ebiten.DrawImageOptions {
	spriteOp := &ebiten.DrawImageOptions{}
	// Scale the sprite to 1/10th size.
	spriteOp.GeoM.Scale(0.1, 0.1)
	spriteOp.GeoM.Rotate(player.angle)
	// Translate the rotated sprite to the player's position
	spriteOp.GeoM.Translate(player.Object.Position.X, player.Object.Position.Y)
	return spriteOp
}

func (player *Player) getRotatedSprite() *ebiten.Image {
	return player.Sprite.SubImage(
		image.Rect(0, 0, player.Sprite.Bounds().Dx(), player.Sprite.Bounds().Dy()),
	).(*ebiten.Image)
}
