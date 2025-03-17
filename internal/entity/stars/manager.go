package stars

import (
	"crypto/rand"
	"encoding/binary"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jonesrussell/gimbal/internal/common"
)

// Manager handles multiple stars and implements Entity interface
type Manager struct {
	stars        []*Star
	screenBounds common.Size
	baseSprite   *ebiten.Image
	config       struct {
		starSize  float64
		starSpeed float64
	}
}

// randomFloat64 generates a random float64 between 0 and 1 using crypto/rand
func randomFloat64() float64 {
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		return 0
	}
	return float64(binary.LittleEndian.Uint64(b[:])) / (1 << 64)
}

// NewManager creates a new star manager
func NewManager(bounds common.Size, numStars int, starSize, starSpeed float64) *Manager {
	// Create base sprite for stars
	baseSprite := ebiten.NewImage(1, 1)
	baseSprite.Fill(color.White)

	m := &Manager{
		stars:        make([]*Star, numStars),
		screenBounds: bounds,
		baseSprite:   baseSprite,
		config: struct {
			starSize  float64
			starSpeed float64
		}{
			starSize:  starSize,
			starSpeed: starSpeed,
		},
	}

	m.initializeStars()
	return m
}

// initializeStars creates the initial set of stars
func (m *Manager) initializeStars() {
	for i := range m.stars {
		pos := common.Point{
			X: randomFloat64() * float64(m.screenBounds.Width),
			Y: randomFloat64() * float64(m.screenBounds.Height),
		}
		star := New(pos, m.config.starSpeed, m.config.starSize, m.baseSprite)
		star.SetBounds(m.screenBounds)
		m.stars[i] = star
	}
}

// Update implements Entity interface
func (m *Manager) Update() {
	for _, star := range m.stars {
		star.Update()
	}
}

// Draw implements Entity interface
func (m *Manager) Draw(screen *ebiten.Image) {
	for _, star := range m.stars {
		star.Draw(screen)
	}
}

// GetPosition implements Entity interface
func (m *Manager) GetPosition() common.Point {
	// Manager doesn't have a position, return center of screen
	return common.Point{
		X: float64(m.screenBounds.Width) / common.CenterDivisor,
		Y: float64(m.screenBounds.Height) / common.CenterDivisor,
	}
}

// GetStars returns all stars
func (m *Manager) GetStars() []*Star {
	return m.stars
}
