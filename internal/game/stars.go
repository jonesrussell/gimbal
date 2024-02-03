package game

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

type Star struct {
	X, Y, Size, Angle, Speed float64
}

var stars []Star

func updateStars() {
	for i := range stars {
		// Update star position based on its angle and speed
		stars[i].X += stars[i].Speed * math.Cos(stars[i].Angle)
		stars[i].Y += stars[i].Speed * math.Sin(stars[i].Angle)

		// If star goes off screen, reset it to the center
		if stars[i].X < 0 || stars[i].X > float64(screenWidth) || stars[i].Y < 0 || stars[i].Y > float64(screenHeight) {
			stars[i].X = float64(screenWidth) / 2
			stars[i].Y = float64(screenHeight) / 2
			stars[i].Size = rand.Float64() * 5
			stars[i].Angle = rand.Float64() * 2 * math.Pi
			stars[i].Speed = rand.Float64() * 2
		}
	}
}

func drawStars(screen *ebiten.Image) {
	for i, star := range stars {
		// Create a new image for the star
		starImage := ebiten.NewImage(int(star.Size*2), int(star.Size*2))
		starImage.Fill(color.White)

		// Create an option to position the star
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(star.X-star.Size, star.Y-star.Size)

		// Draw the star
		screen.DrawImage(starImage, op)

		// Print debugging information
		fmt.Printf("Star %d: X=%.2f, Y=%.2f, Size=%.2f\n", i, star.X, star.Y, star.Size)
	}
}
