package player

import (
	"errors"
	"fmt"
	"image"
	"math"

	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
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

	playerImage := ebiten.NewImageFromImage(spriteImage)
	player := &Player{
		PlayerInput: PlayerInput{input: input},
		PlayerPosition: PlayerPosition{
			Object: resolv.NewConvexPolygon(
				float64(screenCenter.X),
				float64(screenCenter.Y),
				[]float64{
					-16, -16,
					16, -16,
					16, 16,
					-16, 16,
				},
			),
		},
		PlayerSprite: PlayerSprite{Sprite: playerImage},
		Speed:        speed,
		center:       screenCenter,
	}

	return player, nil
}

func (p *Player) Update() {
	// Update the player's angle based on input
	if p.input.IsKeyPressed(ebiten.KeyLeft) {
		p.angle -= 0.05
	}
	if p.input.IsKeyPressed(ebiten.KeyRight) {
		p.angle += 0.05
	}

	// Calculate new position based on angle
	p.PlayerPosition.Object.SetPosition(
		float64(p.center.X)+math.Cos(p.angle)*radius,
		float64(p.center.Y)+math.Sin(p.angle)*radius,
	)
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
