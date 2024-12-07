package engine

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"go.uber.org/zap"

	"github.com/jonesrussell/gimbal/internal/config"
)

// Star represents a star in the game
type Star struct {
	// 8-byte pointer field first
	image *ebiten.Image // 8 bytes

	// 8-byte aligned fields next
	x     float64 // 8 bytes
	y     float64 // 8 bytes
	speed float64 // 8 bytes
	angle float64 // 8 bytes

	// 4-byte field
	size  int32      // 4 bytes
	color color.RGBA // 4 bytes
}

// Memory layout visualization:
// |-----------------------------------------------|
// | x, y (16 bytes)                               |
// |-----------------------------------------------|
// | speed (8 bytes)                               |
// |-----------------------------------------------|
// | angle (8 bytes)                               |
// |-----------------------------------------------|
// | distance (8 bytes)                            |
// |-----------------------------------------------|
// | size (4) | color (4)                          |
// |-----------------------------------------------|

func initializeStars(numStars int, starImage *ebiten.Image) ([]Star, error) {
	cfg, err := config.New()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize stars: %w", err)
	}

	stars := make([]Star, numStars)
	for i := range stars {
		stars[i] = Star{
			x:     float64(cfg.Screen.Width) / 2,
			y:     float64(cfg.Screen.Height) / 2,
			size:  int32(rand.Float64()*5 + 1),
			angle: rand.Float64() * 2 * math.Pi,
			speed: rand.Float64() * 2,
			color: color.RGBA{
				R: uint8(rand.Float64() * 256),
				G: uint8(rand.Float64() * 256),
				B: uint8(rand.Float64() * 256),
				A: 255,
			},
			image: starImage,
		}
	}
	return stars, nil
}

func (g *Game) updateStars() error {
	for i := range g.stars {
		// Update star position based on its angle and speed
		g.stars[i].x += g.stars[i].speed * math.Cos(g.stars[i].angle)
		g.stars[i].y += g.stars[i].speed * math.Sin(g.stars[i].angle)

		// If star goes off screen, reset it to the center
		if g.stars[i].x < 0 || g.stars[i].x > float64(g.config.Screen.Width) ||
			g.stars[i].y < 0 || g.stars[i].y > float64(g.config.Screen.Height) {
			g.stars[i].x = float64(g.config.Screen.Width) / 2
			g.stars[i].y = float64(g.config.Screen.Height) / 2
			g.stars[i].size = int32(rand.Float64()*5 + 1)
			g.stars[i].angle = rand.Float64() * 2 * math.Pi
			g.stars[i].speed = rand.Float64() * 2
		}
	}
	return nil
}

func (g *Game) drawStars(screen *ebiten.Image) {
	for _, star := range g.stars {
		// Calculate the size of the star image
		size := int(star.size * 2)
		if size < 1 {
			size = 1
		}

		// Create an option to position the star
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(float64(size), float64(size))
		op.GeoM.Translate(star.x-float64(size)/2, star.y-float64(size)/2)

		// Draw the star using the image field of the Star struct
		screen.DrawImage(star.image, op)

		// Debugging output
		if g.config.Game.Debug {
			g.logger.Debug("Star position",
				zap.Float64("X", star.x),
				zap.Float64("Y", star.y),
				zap.Float64("Size", float64(star.size))) // Convert int32 to float64
		}
	}
}
