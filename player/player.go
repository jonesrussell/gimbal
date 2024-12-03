package player

import (
	"errors"
	"fmt"
	"image"
	"math"

	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/jonesrussell/gimbal/logger"
	"github.com/solarlune/resolv"
)

type PlayerInput struct {
	input InputHandlerInterface
}

type PlayerPosition struct {
	Object *resolv.ConvexPolygon
}

type PlayerSprite struct {
	Sprite *ebiten.Image
}

type PlayerPath struct {
	path []resolv.Vector
}

type Position struct {
	X, Y float64
}

type Player struct {
	PlayerInput
	PlayerPosition
	PlayerSprite
	PlayerPath
	viewAngle   float64
	direction   float64
	angle       float64
	Speed       float64
	gameStarted bool
	center      image.Point
}

// NewPlayer creates a new instance of a player with the given input handler, speed, and sprite image.
// If any of the arguments are nil, an error is returned.
func NewPlayer(input InputHandlerInterface, speed float64, spriteImage image.Image, screenCenter image.Point) (*Player, error) {
	if input == nil {
		return nil, errors.New("input handler cannot be nil")
	}

	if speed <= 0 {
		return nil, errors.New("speed must be greater than zero")
	}

	if spriteImage == nil {
		return nil, errors.New("sprite image cannot be nil")
	}

	// calculate the initial angle of the player (270 degrees)
	initialAngle := math.Pi * 1.5 // 270 degrees or bottom of the screen

	// calculate the initial X and Y positions using the provided center
	initialX := float64(screenCenter.X) + radius*(-1.0)
	initialY := float64(screenCenter.Y) + radius

	// Create a rectangular polygon for collision
	playerShape := resolv.NewRectangle(
		float64(initialX),
		float64(initialY),
		float64(playerWidth),
		float64(playerHeight),
	)

	// create a new instance of a player with the given input handler, initial position, and sprite image
	player := &Player{
		PlayerInput: PlayerInput{
			input: input,
		},
		PlayerPosition: PlayerPosition{
			Object: playerShape,
		},
		PlayerSprite: PlayerSprite{
			Sprite: ebiten.NewImageFromImage(spriteImage),
		},
		PlayerPath: PlayerPath{},
		viewAngle:  initialAngle,
		Speed:      speed,
		center:     screenCenter,
	}

	return player, nil
}

func (player *Player) Update() {
	if !player.gameStarted {
		logger.GlobalLogger.Debug("Player", "viewAngle", player.viewAngle, "direction", player.direction, "angle", player.angle, "X", player.Object.Position().X, "Y", player.Object.Position().Y)
		player.gameStarted = true
	}

	oldOrientation := player.viewAngle
	oldDirection := player.direction
	oldAngle := player.angle
	oldPosition := player.Object.Position()

	if player.input.IsKeyPressed(ebiten.KeyLeft) {
		player.direction = -1
		player.viewAngle -= AngleStep
	} else if player.input.IsKeyPressed(ebiten.KeyRight) {
		player.direction = 1
		player.viewAngle += AngleStep
	} else {
		player.direction = 0
	}

	position := player.calculatePosition()
	logger.GlobalLogger.Info("position", "full", position)

	// Update the polygon's position
	player.Object.SetPosition(position.X, position.Y)

	player.angle = player.calculateAngle()

	newPosition := player.Object.Position()
	if player.viewAngle != oldOrientation || player.direction != oldDirection || player.angle != oldAngle || newPosition != oldPosition {
		logger.GlobalLogger.Debug("Player", "viewAngle", player.viewAngle, "direction", player.direction, "angle", player.angle, "X", newPosition.X, "Y", newPosition.Y)
	}

	// Add the current position to the path
	player.path = append(player.path, newPosition)

	// Remove the call to Update, as ConvexPolygon does not have an Update method
	// player.Object.Update()
}

var prevRectX, prevRectY float64

func (player *Player) Draw(screen *ebiten.Image) {
	// Draw the player's sprite
	player.drawSprite(screen)

	if Debug {
		// Draw the player's path
		player.drawPath(screen)
		// Draw the rectangle image onto the screen
		player.drawRectangle(screen)
	}
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

func (player *Player) drawRectangle(screen *ebiten.Image) {
	// Create a new image for the rectangle
	rectColor := color.RGBA{255, 0, 0, 255}
	bounds := player.Object.Bounds()
	img := ebiten.NewImage(int(bounds.Width()), int(bounds.Height()))
	img.Fill(rectColor)

	// Calculate the rectangle's top-left corner position
	pos := player.Object.Position()
	rectX := pos.X - bounds.Width()/2
	rectY := pos.Y - bounds.Height()/2

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

func (player *Player) UpdatePosition() {
	pos := player.Object.Position()
	op := &ebiten.DrawImageOptions{}
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
	someValue := float64(width) / 2 // Convert to float64 and divide by 2
	height := float64(player.Sprite.Bounds().Dy())

	spriteOp := &ebiten.DrawImageOptions{}

	// Translate the sprite so that its center is at the origin
	spriteOp.GeoM.Translate(-someValue, -height/2)

	// Scale the sprite to 1/10th size and rotate
	spriteOp.GeoM.Scale(0.1, 0.1)
	spriteOp.GeoM.Rotate(player.angle)

	// Translate the rotated and scaled sprite to the player's position
	pos := player.Object.Position()
	spriteOp.GeoM.Translate(pos.X, pos.Y)

	return spriteOp
}

func (player *Player) getRotatedSprite() *ebiten.Image {
	return player.Sprite.SubImage(
		image.Rect(0, 0, player.Sprite.Bounds().Dx(), player.Sprite.Bounds().Dy()),
	).(*ebiten.Image)
}

// calculateAngle calculates the rotation angle for the player sprite
func (player *Player) calculateAngle() float64 {
	return player.viewAngle - math.Pi/2
}

// Add this method if it's needed
func (p *Player) calculateCoordinates() (float64, float64) {
	// Use the Object's position and player's angle
	pos := p.Object.Position()
	x := pos.X + (p.Speed * math.Cos(p.angle))
	y := pos.Y + (p.Speed * math.Sin(p.angle))

	return x, y
}

func (p *Player) calculatePosition() resolv.Vector {
	x, y := p.calculateCoordinates()

	// Add boundary checks using the screen dimensions
	screenWidth := float64(p.center.X * 2)  // Assuming center is half the screen width
	screenHeight := float64(p.center.Y * 2) // Assuming center is half the screen height

	x = math.Max(0, math.Min(x, screenWidth))
	y = math.Max(0, math.Min(y, screenHeight))

	return resolv.Vector{X: x, Y: y}
}
