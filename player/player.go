package player

import (
	"errors"
	"fmt"
	"image"
	"math"

	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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
	// Update orbital position
	if p.input.IsKeyPressed(ebiten.KeyLeft) {
		p.angle -= 0.05
	}
	if p.input.IsKeyPressed(ebiten.KeyRight) {
		p.angle += 0.05
	}

	// Calculate new position
	newX := float64(p.center.X) + math.Cos(p.angle)*radius
	newY := float64(p.center.Y) + math.Sin(p.angle)*radius

	// Calculate rotation angle to face center
	// Add π to rotate 180 degrees (sprite faces inward)
	p.viewAngle = math.Atan2(
		float64(p.center.Y)-newY,
		float64(p.center.X)-newX,
	) - math.Pi/2 + math.Pi

	// Update position
	p.PlayerPosition.Object.SetPosition(newX, newY)
}

var prevRectX, prevRectY float64

func (player *Player) Draw(screen *ebiten.Image) {
	// Draw the regular sprite
	player.drawSprite(screen)

	// Draw debug visualization if debug mode is enabled
	if Debug {
		player.drawDebugRotation(screen)
	}
}

func (p *Player) drawDebugRotation(screen *ebiten.Image) {
	pos := p.Object.Position()

	// Draw line from player to center (direction vector)
	vector.StrokeLine(
		screen,
		float32(pos.X),
		float32(pos.Y),
		float32(p.center.X),
		float32(p.center.Y),
		1,
		color.RGBA{R: 255, G: 0, B: 0, A: 255}, // Red line
		false,
	)

	// Draw angle arc
	radius := 20.0 // Small radius for angle visualization
	startAngle := p.angle - math.Pi/4
	endAngle := p.angle + math.Pi/4

	for angle := startAngle; angle <= endAngle; angle += 0.1 {
		x := pos.X + math.Cos(angle)*radius
		y := pos.Y + math.Sin(angle)*radius
		vector.StrokeLine(
			screen,
			float32(pos.X),
			float32(pos.Y),
			float32(x),
			float32(y),
			1,
			color.RGBA{R: 0, G: 255, B: 0, A: 255}, // Green arc
			false,
		)
	}

	// Draw text showing angle values
	text := fmt.Sprintf("Angle: %.2f°\nView: %.2f°",
		p.angle*180/math.Pi,
		p.viewAngle*180/math.Pi,
	)
	ebitenutil.DebugPrintAt(
		screen,
		text,
		int(pos.X)+10,
		int(pos.Y)-20,
	)
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
	width := player.Sprite.Bounds().Dx()
	height := float64(player.Sprite.Bounds().Dy())
	spriteOp := &ebiten.DrawImageOptions{}

	// Center the sprite
	spriteOp.GeoM.Translate(-float64(width)/2, -height/2)

	// Scale the sprite to 1/10th size
	spriteOp.GeoM.Scale(0.1, 0.1)

	// Use viewAngle for rotation instead of angle
	spriteOp.GeoM.Rotate(player.viewAngle)

	// Move to player position
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
