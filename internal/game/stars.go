package game

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/jonesrussell/gimbal/internal/common"
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
	stars := g.stars.GetStars()
	for _, star := range stars {
		pos := star.GetPosition()
		speed := star.GetSpeed()
		angle := star.GetAngle()

		// Update position based on speed and angle
		pos.X += speed * math.Cos(angle)
		pos.Y += speed * math.Sin(angle)

		// Reset star if it goes off screen
		if pos.X < 0 || pos.X > float64(g.config.ScreenSize.Width) ||
			pos.Y < 0 || pos.Y > float64(g.config.ScreenSize.Height) {
			// Reset to center with random properties
			pos = common.Point{
				X: float64(g.config.ScreenSize.Width) / 2,
				Y: float64(g.config.ScreenSize.Height) / 2,
			}
			star.SetPosition(pos)
			star.SetSpeed(randomFloat64() * 2)
			star.SetAngle(randomFloat64() * 2 * math.Pi)
			star.SetSize(randomFloat64() * 2)
			continue
		}

		star.SetPosition(pos)
	}
}

func (g *GimlarGame) drawStars(screen *ebiten.Image) {
	stars := g.stars.GetStars()
	for _, star := range stars {
		pos := star.GetPosition()
		sprite := star.GetSprite()
		size := star.GetSize()

		if sprite == nil {
			continue
		}

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(size, size)
		op.GeoM.Translate(pos.X, pos.Y)
		screen.DrawImage(sprite, op)

		if g.config.Debug {
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Star: pos=(%v,%v) size=%v", pos.X, pos.Y, size), int(pos.X), int(pos.Y))
		}
	}
}
