package effects

import (
	"image/color"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// Star represents a single star in the starfield
type Star struct {
	X      float64
	Y      float64
	Speed  float64
	Size   float64
	Brightness float64
}

// Starfield represents an animated starfield effect
type Starfield struct {
	stars       []Star
	width       int
	height      int
	speed       float64
	starCount   int
	elapsed     float64
	initialized bool
}

// NewStarfield creates a new starfield effect
func NewStarfield(width, height int, starCount int, speed float64) *Starfield {
	sf := &Starfield{
		width:     width,
		height:    height,
		speed:     speed,
		starCount: starCount,
		stars:     make([]Star, starCount),
	}
	sf.initialize()
	return sf
}

// initialize creates the initial star positions
func (sf *Starfield) initialize() {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	
	for i := 0; i < sf.starCount; i++ {
		sf.stars[i] = Star{
			X:          float64(rng.Intn(sf.width)),
			Y:          float64(rng.Intn(sf.height)),
			Speed:      0.5 + rng.Float64()*sf.speed,
			Size:       0.5 + rng.Float64()*1.5,
			Brightness: 0.5 + rng.Float64()*0.5,
		}
	}
	sf.initialized = true
}

// Update updates the starfield animation
func (sf *Starfield) Update(deltaTime float64) {
	sf.elapsed += deltaTime
	
	for i := range sf.stars {
		// Move stars downward (or in direction of travel)
		sf.stars[i].Y += sf.stars[i].Speed * deltaTime * 60.0 // Normalize to 60 FPS
		
		// Wrap stars around when they go off screen
		if sf.stars[i].Y >= float64(sf.height) {
			sf.stars[i].Y = 0
			sf.stars[i].X = float64(rand.Intn(sf.width))
		}
	}
}

// Draw draws the starfield
func (sf *Starfield) Draw(screen *ebiten.Image) {
	if !sf.initialized {
		return
	}

	for _, star := range sf.stars {
		// Draw star as a small filled rectangle (simple and fast)
		brightness := uint8(star.Brightness * 255)
		size := int(star.Size) + 1
		if size < 1 {
			size = 1
		}
		
		// Create star color
		starColor := color.RGBA{brightness, brightness, brightness, 255}
		
		// Draw star as a small square (simple and performant)
		starImg := ebiten.NewImage(size, size)
		starImg.Fill(starColor)
		
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(star.X-float64(size)/2, star.Y-float64(size)/2)
		screen.DrawImage(starImg, op)
	}
}

// SetSpeed sets the speed multiplier for the starfield
func (sf *Starfield) SetSpeed(speed float64) {
	sf.speed = speed
	for i := range sf.stars {
		sf.stars[i].Speed = 0.5 + (sf.stars[i].Speed-0.5)*(speed/sf.speed)
	}
}

// Reset resets the starfield to initial state
func (sf *Starfield) Reset() {
	sf.initialize()
	sf.elapsed = 0
}
