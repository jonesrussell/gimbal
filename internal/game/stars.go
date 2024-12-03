package game

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jonesrussell/gimbal/internal/config"
)

type Star struct {
	X, Y, Size, Angle, Speed float64
	Image                    *ebiten.Image
}

func initializeStars(numStars int, starImage *ebiten.Image) ([]Star, error) {
	cfg, err := config.New()
	if err != nil {
		return nil, fmt.Errorf("failed to get config for star initialization: %w", err)
	}

	stars := make([]Star, numStars)
	for i := range stars {
		stars[i] = Star{
			X:     float64(cfg.Screen.Width) / 2,
			Y:     float64(cfg.Screen.Height) / 2,
			Size:  rand.Float64()*5 + 1,
			Angle: rand.Float64() * 2 * math.Pi,
			Speed: rand.Float64() * 2,
			Image: starImage,
		}
	}
	return stars, nil
}

func (g *GimlarGame) updateStars() error {
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("failed to get config for star update: %w", err)
	}

	for i := range g.stars {
		// Update star position based on its angle and speed
		g.stars[i].X += g.stars[i].Speed * math.Cos(g.stars[i].Angle)
		g.stars[i].Y += g.stars[i].Speed * math.Sin(g.stars[i].Angle)

		// If star goes off screen, reset it to the center
		if g.stars[i].X < 0 || g.stars[i].X > float64(cfg.Screen.Width) ||
			g.stars[i].Y < 0 || g.stars[i].Y > float64(cfg.Screen.Height) {
			g.stars[i].X = float64(cfg.Screen.Width) / 2
			g.stars[i].Y = float64(cfg.Screen.Height) / 2
			g.stars[i].Size = rand.Float64() * 5
			g.stars[i].Angle = rand.Float64() * 2 * math.Pi
			g.stars[i].Speed = rand.Float64() * 2
		}
	}
	return nil
}

func (g *GimlarGame) drawStars(screen *ebiten.Image) {
	for _, star := range g.stars {
		// Calculate the size of the star image
		size := int(star.Size * 2)
		if size < 1 {
			size = 1
		}

		// Create an option to position the star
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(float64(size), float64(size))
		op.GeoM.Translate(star.X-float64(size)/2, star.Y-float64(size)/2)

		// Draw the star using the Image field of the Star struct
		screen.DrawImage(star.Image, op)

		// Debugging output
		if Debug {
			fmt.Printf("Star: X=%.2f, Y=%.2f, Size=%.2f\n", star.X, star.Y, star.Size)
		}
	}
}
