package ecs

import (
	"math"
	"math/rand"
	"time"

	"github.com/jonesrussell/gimbal/internal/common"
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

// DefaultStarFieldConfig returns a default star field configuration
func DefaultStarFieldConfig() *StarFieldConfig {
	return &StarFieldConfig{
		SpawnRadiusMin: 30.0,
		SpawnRadiusMax: 80.0,
		Speed:          2.0,
		MinScale:       0.3,
		MaxScale:       1.0,
		ScaleDistance:  200.0,
		ResetMargin:    50.0,
		Seed:           time.Now().UnixNano(),
	}
}

// DenseStarFieldConfig returns a configuration for a dense star field
func DenseStarFieldConfig() *StarFieldConfig {
	return &StarFieldConfig{
		SpawnRadiusMin: 20.0,
		SpawnRadiusMax: 60.0,
		Speed:          1.5,
		MinScale:       0.2,
		MaxScale:       0.8,
		ScaleDistance:  150.0,
		ResetMargin:    40.0,
		Seed:           time.Now().UnixNano(),
	}
}

// SparseStarFieldConfig returns a configuration for a sparse star field
func SparseStarFieldConfig() *StarFieldConfig {
	return &StarFieldConfig{
		SpawnRadiusMin: 40.0,
		SpawnRadiusMax: 100.0,
		Speed:          3.0,
		MinScale:       0.4,
		MaxScale:       1.2,
		ScaleDistance:  250.0,
		ResetMargin:    60.0,
		Seed:           time.Now().UnixNano(),
	}
}

// FastStarFieldConfig returns a configuration for a fast-moving star field
func FastStarFieldConfig() *StarFieldConfig {
	return &StarFieldConfig{
		SpawnRadiusMin: 25.0,
		SpawnRadiusMax: 70.0,
		Speed:          4.0,
		MinScale:       0.3,
		MaxScale:       1.1,
		ScaleDistance:  180.0,
		ResetMargin:    45.0,
		Seed:           time.Now().UnixNano(),
	}
}

// StarFieldHelper provides helper functions for star field operations
type StarFieldHelper struct {
	config       *StarFieldConfig
	center       common.Point
	screenBounds common.Size
}

// NewStarFieldHelper creates a new star field helper
func NewStarFieldHelper(config *StarFieldConfig, screenBounds common.Size) *StarFieldHelper {
	return &StarFieldHelper{
		config:       config,
		center:       common.Point{X: float64(screenBounds.Width) / 2, Y: float64(screenBounds.Height) / 2},
		screenBounds: screenBounds,
	}
}

// GenerateRandomPosition generates a random position along the spawn orbital path
func (h *StarFieldHelper) GenerateRandomPosition() common.Point {
	// Create a local random generator with the configured seed
	//nolint:gosec // deterministic behavior and performance (not security-critical)
	r := rand.New(rand.NewSource(h.config.Seed))

	// Random angle around the circle (0 to 2π)
	angle := r.Float64() * 2 * math.Pi

	// Random radius within the spawn range
	spawnRadius := h.config.SpawnRadiusMin + r.Float64()*(h.config.SpawnRadiusMax-h.config.SpawnRadiusMin)

	return common.Point{
		X: h.center.X + math.Cos(angle)*spawnRadius,
		Y: h.center.Y + math.Sin(angle)*spawnRadius,
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

// GenerateRandomScale generates a random scale within the configured range
func (h *StarFieldHelper) GenerateRandomScale() float64 {
	//nolint:gosec // deterministic behavior and performance (not security-critical)
	r := rand.New(rand.NewSource(h.config.Seed))
	return h.config.MinScale + r.Float64()*(h.config.MaxScale-h.config.MinScale)
}

// CalculateScale calculates the scale based on distance from center
func (h *StarFieldHelper) CalculateScale(distance float64) float64 {
	scaleRatio := math.Min(distance/h.config.ScaleDistance, 1.0)
	return h.config.MinScale + scaleRatio*(h.config.MaxScale-h.config.MinScale)
}

// IsOffScreen checks if a position is off screen (with margin)
func (h *StarFieldHelper) IsOffScreen(pos common.Point) bool {
	margin := h.config.ResetMargin
	return pos.X < -margin ||
		pos.X > float64(h.screenBounds.Width)+margin ||
		pos.Y < -margin ||
		pos.Y > float64(h.screenBounds.Height)+margin
}

// CalculateMovementDirection calculates the normalized direction vector from center to position
func (h *StarFieldHelper) CalculateMovementDirection(pos common.Point) (dx, dy float64) {
	dx = pos.X - h.center.X
	dy = pos.Y - h.center.Y
	distance := math.Sqrt(dx*dx + dy*dy)

	// Normalize direction vector
	if distance > 0 {
		dx /= distance
		dy /= distance
	}

	return dx, dy
}

// ResetStar resets a star to a new random position and scale
func (h *StarFieldHelper) ResetStar(pos *common.Point, scale *float64) {
	newPos := h.GenerateRandomPosition()
	pos.X = newPos.X
	pos.Y = newPos.Y
	*scale = h.GenerateRandomScale()
}

// UpdateStar updates a star's position and scale based on movement
func (h *StarFieldHelper) UpdateStar(pos *common.Point, scale *float64) {
	h.UpdateStarWithSpeed(pos, scale, h.config.Speed)
}

// UpdateStarWithSpeed updates a star's position and scale based on movement with custom speed
func (h *StarFieldHelper) UpdateStarWithSpeed(pos *common.Point, scale *float64, speed float64) {
	// Calculate movement direction
	dx, dy := h.CalculateMovementDirection(*pos)

	// Move star outward from center
	pos.X += dx * speed
	pos.Y += dy * speed

	// Calculate distance from center for scaling
	distance := math.Sqrt((pos.X-h.center.X)*(pos.X-h.center.X) + (pos.Y-h.center.Y)*(pos.Y-h.center.Y))

	// Update scale based on distance
	*scale = h.CalculateScale(distance)

	// Reset if off screen
	if h.IsOffScreen(*pos) {
		h.ResetStar(pos, scale)
	}
}
