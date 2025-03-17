package game

import (
	"errors"
	"fmt"
	"image"
	"math"

	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/jonesrussell/gimbal/internal/logger"
	"github.com/solarlune/resolv"
)

const (
	initialAngleMultiplier = 1.5 // 270 degrees or bottom of the screen
	radiusDivisor          = 4
	spriteScaleFactor      = 0.1
	spriteCenterDivisor    = 2
	pathStrokeWidth        = 1.0
	debugAlpha             = 255
	positionDivisor        = 2 // Used for dividing screen dimensions and player dimensions
)

type PlayerInput struct {
	Input InputHandlerInterface
}

type PlayerPosition struct {
	Object *resolv.ConvexPolygon
}

type PlayerSprite struct {
	Sprite *ebiten.Image
}

type PlayerPath struct {
	Path []resolv.Vector
}

type Player struct {
	PlayerInput
	PlayerPosition
	PlayerSprite
	PlayerPath
	ViewAngle float64
	Direction float64
	Angle     float64
	Config    *GameConfig
}

// NewPlayer creates a new instance of a player with the given input handler, speed, and sprite image.
// If any of the arguments are nil, an error is returned.
func NewPlayer(input InputHandlerInterface, config *GameConfig, spriteImage image.Image) (*Player, error) {
	if input == nil {
		return nil, errors.New("input handler cannot be nil")
	}

	if config == nil {
		return nil, errors.New("config cannot be nil")
	}

	if spriteImage == nil {
		return nil, errors.New("sprite image cannot be nil")
	}

	// calculate the initial angle of the player (270 degrees)
	initialAngle := math.Pi * initialAngleMultiplier

	// calculate the initial X and Y positions of the player based on the center point and the initial angle
	centerX := float64(config.ScreenWidth) / positionDivisor
	centerY := float64(config.ScreenHeight) / positionDivisor
	radius := float64(config.ScreenHeight) / radiusDivisor

	initialX := centerX + radius*math.Cos(initialAngle)
	initialY := centerY - radius*math.Sin(initialAngle) - float64(config.PlayerHeight)/positionDivisor

	// create a new instance of a player with the given input handler, initial position, and sprite image
	player := &Player{
		PlayerInput: PlayerInput{
			Input: input,
		},
		PlayerPosition: PlayerPosition{
			Object: resolv.NewRectangle(initialX, initialY, float64(config.PlayerWidth), float64(config.PlayerHeight)),
		},
		PlayerSprite: PlayerSprite{
			Sprite: ebiten.NewImageFromImage(spriteImage),
		},
		PlayerPath: PlayerPath{},
		ViewAngle:  initialAngle,
		Config:     config,
	}

	return player, nil
}

func (player *Player) Update() {
	pos := player.Object.Position()
	logger.GlobalLogger.Debug(
		"Player",
		"viewAngle", player.ViewAngle,
		"direction", player.Direction,
		"angle", player.Angle,
		"X", pos.X,
		"Y", pos.Y,
	)

	oldOrientation := player.ViewAngle
	oldDirection := player.Direction
	oldAngle := player.Angle
	oldPos := player.Object.Position()

	switch {
	case player.Input.IsKeyPressed(ebiten.KeyLeft):
		player.Direction = -1
		player.ViewAngle -= player.Config.AngleStep
	case player.Input.IsKeyPressed(ebiten.KeyRight):
		player.Direction = 1
		player.ViewAngle += player.Config.AngleStep
	default:
		player.Direction = 0
	}

	position := player.CalculatePosition()
	logger.GlobalLogger.Info("position", "full", position)

	player.Object = resolv.NewRectangle(
		position.X,
		position.Y,
		float64(player.Config.PlayerWidth),
		float64(player.Config.PlayerHeight),
	)

	player.Angle = player.CalculateAngle()

	newPos := player.Object.Position()
	if player.ViewAngle != oldOrientation || player.Direction != oldDirection ||
		player.Angle != oldAngle || newPos.X != oldPos.X || newPos.Y != oldPos.Y {
		logger.GlobalLogger.Debug(
			"Player",
			"viewAngle", player.ViewAngle,
			"direction", player.Direction,
			"angle", player.Angle,
			"X", newPos.X,
			"Y", newPos.Y,
		)
	}

	// Add the current position to the path
	player.Path = append(player.Path, newPos)
}

var prevRectX, prevRectY float64

func (player *Player) Draw(screen *ebiten.Image) {
	// Draw the player's sprite
	player.drawSprite(screen)

	if player.Config.Debug {
		// Draw the player's path
		player.drawPath(screen)
		// Draw the rectangle image onto the screen
		player.drawRectangle(screen)
	}
}

// Draw the players path
func (player *Player) drawPath(screen *ebiten.Image) {
	for i := range player.Path[:len(player.Path)-1] {
		vector.StrokeLine(
			screen,
			float32(player.Path[i].X),
			float32(player.Path[i].Y),
			float32(player.Path[i+1].X),
			float32(player.Path[i+1].Y),
			pathStrokeWidth,
			color.RGBA{255, 0, 0, debugAlpha},
			false,
		)
	}
}

func (player *Player) drawRectangle(screen *ebiten.Image) {
	// Create a new image for the rectangle
	rectColor := color.RGBA{255, 0, 0, debugAlpha}
	img := ebiten.NewImage(player.Config.PlayerWidth, player.Config.PlayerHeight)
	img.Fill(rectColor)

	// Get the position and calculate the rectangle's top-left corner position
	pos := player.Object.Position()
	rectX := pos.X - float64(player.Config.PlayerWidth)/positionDivisor
	rectY := pos.Y - float64(player.Config.PlayerHeight)/positionDivisor

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
	pos := player.Object.Position()
	op.GeoM.Translate(pos.X, pos.Y)
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
	someValue := float64(width) / spriteCenterDivisor
	height := float64(player.Sprite.Bounds().Dy())

	spriteOp := &ebiten.DrawImageOptions{}

	// Translate the sprite so that its center is at the origin
	spriteOp.GeoM.Translate(-someValue, -height/spriteCenterDivisor)

	// Scale the sprite to 1/10th size and rotate
	spriteOp.GeoM.Scale(spriteScaleFactor, spriteScaleFactor)
	spriteOp.GeoM.Rotate(player.Angle)

	// Translate the rotated and scaled sprite to the player's position
	pos := player.Object.Position()
	spriteOp.GeoM.Translate(pos.X, pos.Y)

	return spriteOp
}

func (player *Player) getRotatedSprite() *ebiten.Image {
	img, ok := player.Sprite.SubImage(
		image.Rect(0, 0, player.Sprite.Bounds().Dx(), player.Sprite.Bounds().Dy()),
	).(*ebiten.Image)
	if !ok {
		logger.GlobalLogger.Error("Failed to get rotated sprite: type assertion failed")
		return nil
	}
	return img
}

// GetAngle returns the player's current angle
func (player *Player) GetAngle() float64 {
	return player.Angle
}

// SetAngle sets the player's current angle
func (player *Player) SetAngle(angle float64) {
	player.Angle = angle
}

// GetDirection returns the player's current direction
func (player *Player) GetDirection() float64 {
	return player.Direction
}

// SetDirection sets the player's current direction
func (player *Player) SetDirection(direction float64) {
	player.Direction = direction
}

// GetViewAngle returns the player's current view angle
func (player *Player) GetViewAngle() float64 {
	return player.ViewAngle
}

// SetViewAngle sets the player's current view angle
func (player *Player) SetViewAngle(viewAngle float64) {
	player.ViewAngle = viewAngle
}
