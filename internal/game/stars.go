package game

import (
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	starSizeMultiplier = 2
)

type Star struct {
	X, Y, Size, Angle, Speed float64
	Image                    *ebiten.Image
}

func initializeStars(numStars int, starImage *ebiten.Image) []Star {
	stars := make([]Star, numStars)
	for i := range stars {
		stars[i] = Star{
			X:     float64(DefaultConfig().ScreenWidth) / 2,
			Y:     float64(DefaultConfig().ScreenHeight) / 2,
			Size:  rand.Float64()*5 + 1, // Add 1 to ensure the size is always greater than 0
			Angle: rand.Float64() * 2 * math.Pi,
			Speed: rand.Float64() * 2,
			Image: starImage, // Assign the global starImage to each Star
		}
	}
	return stars
}

func (g *GimlarGame) updateStars() {
	for i := range g.stars {
		// Update star position based on its angle and speed
		g.stars[i].X += g.stars[i].Speed * math.Cos(g.stars[i].Angle)
		g.stars[i].Y += g.stars[i].Speed * math.Sin(g.stars[i].Angle)

		// If star goes off screen, reset it to the center
		if g.stars[i].X < 0 || g.stars[i].X > float64(g.config.ScreenWidth) || g.stars[i].Y < 0 || g.stars[i].Y > float64(g.config.ScreenHeight) {
			g.stars[i].X = float64(g.config.ScreenWidth) / 2
			g.stars[i].Y = float64(g.config.ScreenHeight) / 2
			g.stars[i].Size = rand.Float64() * 5
			g.stars[i].Angle = rand.Float64() * 2 * math.Pi
			g.stars[i].Speed = rand.Float64() * 2
		}
	}
}

func (g *GimlarGame) drawStars(screen *ebiten.Image) {
	for _, star := range g.stars {
		// Calculate the size of the star image
		size := int(star.Size * starSizeMultiplier)
		if size < 1 {
			size = 1
		}

		// Create an option to position the star
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(float64(size), float64(size))
		op.GeoM.Translate(star.X-float64(size)/starSizeMultiplier, star.Y-float64(size)/starSizeMultiplier)

		// Draw the star using the Image field of the Star struct
		screen.DrawImage(star.Image, op)

		// Debugging output
		if g.config.Debug {
			g.DebugPrintStar(screen, star)
		}
	}
}
