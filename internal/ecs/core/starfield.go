package core

import (
	"math"
	"math/rand"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/config"
)

// StarFieldConfig holds configuration for the star field behavior
type StarFieldConfig struct {
	// Spawn settings
	SpawnRadiusMin float64
	SpawnRadiusMax float64

	// Movement settings
	Speed float64

	// Scaling settings
	MinScale      float64
	MaxScale      float64
	ScaleDistance float64 // Distance at which stars reach max scale

	// Screen bounds for reset
	ResetMargin float64

	// Random seed for consistent behavior
	Seed int64
}

// StarFieldHelper provides helper functions for star field operations
type StarFieldHelper struct {
	config       *StarFieldConfig
	center       common.Point
	screenBounds config.Size
}

// NewStarFieldHelper creates a new star field helper
func NewStarFieldHelper(starConfig *StarFieldConfig, screenBounds config.Size) *StarFieldHelper {
	return &StarFieldHelper{
		config:       starConfig,
		center:       common.Point{X: float64(screenBounds.Width) / 2, Y: float64(screenBounds.Height) / 2},
		screenBounds: screenBounds,
	}
}

// GenerateRandomPositionWithOffset generates a random position with a seed offset
func (h *StarFieldHelper) GenerateRandomPositionWithOffset(offset int64) common.Point {
	// Create a local random generator with the configured seed plus offset
	//nolint:gosec // deterministic behavior and performance (not security-critical)
	r := rand.New(rand.NewSource(h.config.Seed + offset))

	// Random angle around the circle (0 to 2π)
	angle := r.Float64() * 2 * math.Pi

	// Random radius within the spawn range
	spawnRadius := h.config.SpawnRadiusMin + r.Float64()*(h.config.SpawnRadiusMax-h.config.SpawnRadiusMin)

	return common.Point{
		X: h.center.X + math.Cos(angle)*spawnRadius,
		Y: h.center.Y + math.Sin(angle)*spawnRadius,
	}
}
