package game

import (
	"crypto/rand"
	"encoding/binary"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	starSizeMultiplier = 2
	minStarSize        = 1
	maxStarSize        = 5
	maxStarSpeed       = 2
	twoPi              = 2 * math.Pi
	randomBitShift     = 64 // Number of bits to shift for random float64 generation
)

type Star struct {
	X, Y, Size, Angle, Speed float64
	Image                    *ebiten.Image
}

// randomFloat64 generates a random float64 between 0 and 1 using crypto/rand
func randomFloat64() float64 {
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		return 0
	}
	return float64(binary.LittleEndian.Uint64(b[:])) / (1 << randomBitShift)
}

func initializeStars(numStars int, starImage *ebiten.Image) []Star {
	stars := make([]Star, numStars)
	for i := range stars {
		stars[i] = Star{
			X:     float64(DefaultConfig().ScreenWidth) / screenCenterDivisor,
			Y:     float64(DefaultConfig().ScreenHeight) / screenCenterDivisor,
			Size:  randomFloat64()*maxStarSize + minStarSize,
			Angle: randomFloat64() * twoPi,
			Speed: randomFloat64() * maxStarSpeed,
			Image: starImage,
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
		if g.stars[i].X < 0 || g.stars[i].X > float64(g.config.ScreenWidth) ||
			g.stars[i].Y < 0 || g.stars[i].Y > float64(g.config.ScreenHeight) {
			g.stars[i].X = float64(g.config.ScreenWidth) / screenCenterDivisor
			g.stars[i].Y = float64(g.config.ScreenHeight) / screenCenterDivisor
			g.stars[i].Size = randomFloat64() * maxStarSize
			g.stars[i].Angle = randomFloat64() * twoPi
			g.stars[i].Speed = randomFloat64() * maxStarSpeed
		}
	}
}

func (g *GimlarGame) drawStars(screen *ebiten.Image) {
	for _, star := range g.stars {
		// Calculate the size of the star image
		size := int(star.Size * starSizeMultiplier)
		if size < minStarSize {
			size = minStarSize
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
