package game

import (
	"errors"
	"fmt"
	"image"
	"log/slog"

	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
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
	path      []resolv.Vector
	logger    slog.Logger
}

// Draw the players path
func (player *Player) drawPath(screen *ebiten.Image) {
	for i := 0; i < len(player.path)-1; i++ {
		vector.StrokeLine(
			screen,
			float32(player.path[i].X),
			float32(player.path[i].Y),
			float32(player.path[i+1].X),
			float32(player.path[i+1].Y),
			1.0,
			color.RGBA{255, 0, 0, 255},
			false,
		)
	}
}

func NewPlayer(
	input InputHandlerInterface,
	speed float64,
	spriteImage *ebiten.Image,
	logger slog.Logger,
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
		viewAngle: MaxAngle,
		logger:    logger,
		Object:    &resolv.Object{},
	}

	// Calculate the initial position
	position := player.calculatePosition()
	player.Object = resolv.NewObject(
		position.X,
		position.Y,
		float64(playerWidth),
		float64(playerHeight),
	)

	// Calculate the entire path
	player.path = player.calculatePath()

	return player, nil
}

func (player *Player) Update() {
	if !gameStarted {
		player.logger.Debug("Player", "viewAngle", player.viewAngle, "direction", player.direction, "angle", player.angle, "X", float64(player.Object.Position.X), "Y", float64(player.Object.Position.Y))
		gameStarted = true
	}

	oldOrientation := player.viewAngle
	oldDirection := player.direction
	oldAngle := player.angle
	oldX := player.Object.Position.X
	oldY := player.Object.Position.Y

	if player.input.IsKeyPressed(ebiten.KeyLeft) {
		player.direction = -1
		player.viewAngle -= AngleStep
	} else if player.input.IsKeyPressed(ebiten.KeyRight) {
		player.direction = 1
		player.viewAngle += AngleStep
	} else {
		player.direction = 0
	}

	player.Object.Position = player.calculatePosition()

	player.angle = player.calculateAngle()

	if player.viewAngle != oldOrientation || player.direction != oldDirection || player.angle != oldAngle || player.Object.Position.X != oldX || player.Object.Position.Y != oldY {
		player.logger.Debug("Player", "viewAngle", player.viewAngle, "direction", player.direction, "angle", player.angle, "X", float64(player.Object.Position.X), "Y", float64(player.Object.Position.Y))
	}

	// Add the current position to the path
	player.path = append(player.path, player.Object.Position)

	player.Object.Update()
}

var prevRectX, prevRectY float64

func (player *Player) Draw(screen *ebiten.Image) {
	// Draw the player's path
	player.drawPath(screen)
	// Draw the rectangle image onto the screen
	player.drawRectangle(screen)
	// Draw the player's sprite
	player.drawSprite(screen)
}

func (player *Player) drawRectangle(screen *ebiten.Image) {
	// Create a new image for the rectangle
	rectColor := color.RGBA{255, 0, 0, 255}
	img := ebiten.NewImage(int(player.Object.Size.X), int(player.Object.Size.Y))
	img.Fill(rectColor)

	// Calculate the rectangle's top-left corner position
	rectX := player.Object.Position.X - player.Object.Size.X/2
	rectY := player.Object.Position.Y - player.Object.Size.Y/2

	// Check if rectX or rectY has changed since the last call
	if rectX != prevRectX || rectY != prevRectY {
		fmt.Printf("rectX: %f, rectY: %f\n", rectX, rectY)
		prevRectX, prevRectY = rectX, rectY // Update the stored values
	}

	// Draw the rectangle image onto the screen
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(rectX, rectY)
	screen.DrawImage(img, op)
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
	// Calculate the sprite's top-left corner position
	width := player.Sprite.Bounds().Dx()
	someValue := float64(width) / 2 // Convert to float64 and divide by 2
	height := float64(player.Sprite.Bounds().Dy())

	spriteOp := &ebiten.DrawImageOptions{}

	// Translate the sprite so that its center is at the origin
	spriteOp.GeoM.Translate(-someValue, -height/2)

	// Scale the sprite to 1/10th size and rotate
	spriteOp.GeoM.Scale(0.1, 0.1)
	spriteOp.GeoM.Rotate(player.angle)

	// Translate the rotated and scaled sprite to the player's position
	spriteX := player.Object.Position.X
	spriteY := player.Object.Position.Y
	spriteOp.GeoM.Translate(spriteX, spriteY)

	return spriteOp
}

func (player *Player) getRotatedSprite() *ebiten.Image {
	return player.Sprite.SubImage(
		image.Rect(0, 0, player.Sprite.Bounds().Dx(), player.Sprite.Bounds().Dy()),
	).(*ebiten.Image)
}
