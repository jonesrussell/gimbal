package movement

import (
	"context"
	"math"
	"math/rand"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/config"
	"github.com/jonesrussell/gimbal/internal/ecs/core"
	internalmath "github.com/jonesrussell/gimbal/internal/math"
)

// MovementSystem updates entity positions based on velocity or input
// It is responsible for moving the player, starfield, and any other moving entities.
type MovementSystem struct {
	world        donburi.World
	config       *config.GameConfig
	logger       common.Logger
	inputHandler common.GameInputHandler
}

func NewMovementSystem(
	world donburi.World,
	cfg *config.GameConfig,
	logger common.Logger,
	inputHandler common.GameInputHandler,
) *MovementSystem {
	return &MovementSystem{
		world:        world,
		config:       cfg,
		logger:       logger,
		inputHandler: inputHandler,
	}
}

// Update updates positions of all relevant entities (player, stars, etc.)
func (ms *MovementSystem) Update(ctx context.Context, deltaTime float64) error {
	// Update player orbital movement based on input
	ms.updatePlayerMovement(deltaTime)

	// Update starfield radial movement
	ms.updateStarMovement(deltaTime)

	return nil
}

// updatePlayerMovement updates the player's position based on input
func (ms *MovementSystem) updatePlayerMovement(deltaTime float64) {
	// Query for player entities with orbital movement
	query.NewQuery(
		filter.And(
			filter.Contains(core.PlayerTag),
			filter.Contains(core.Orbital),
			filter.Contains(core.Position),
		),
	).Each(ms.world, func(entry *donburi.Entry) {
		// Get the orbital data
		orbital := core.Orbital.Get(entry)
		pos := core.Position.Get(entry)

		// Get movement input from input handler
		movementInput := ms.inputHandler.GetMovementInput()

		// Update orbital angle based on input
		orbital.OrbitalAngle += internalmath.Angle(movementInput)

		// Keep angle within 0-360 range
		if orbital.OrbitalAngle < 0 {
			orbital.OrbitalAngle += 360
		} else if orbital.OrbitalAngle >= 360 {
			orbital.OrbitalAngle -= 360
		}

		// Calculate new position on the orbital path
		radians := float64(orbital.OrbitalAngle) * math.Pi / 180.0
		newX := orbital.Center.X + math.Cos(radians)*orbital.Radius
		newY := orbital.Center.Y + math.Sin(radians)*orbital.Radius

		// Update position
		pos.X = newX
		pos.Y = newY

		// Update facing angle to point towards center (Gyruss-style)
		// Calculate angle from player to center
		dx := orbital.Center.X - pos.X
		dy := orbital.Center.Y - pos.Y
		facingAngle := math.Atan2(dy, dx) * 180.0 / math.Pi
		orbital.FacingAngle = internalmath.Angle(facingAngle)

		ms.logger.Debug("Player orbital movement updated",
			"input", movementInput,
			"orbital_angle", orbital.OrbitalAngle,
			"position", pos,
			"facing_angle", orbital.FacingAngle)
	})
}

// updateStarMovement updates the starfield positions to create a radial effect
func (ms *MovementSystem) updateStarMovement(deltaTime float64) {
	// Query for star entities with speed component
	query.NewQuery(
		filter.And(
			filter.Contains(core.StarTag),
			filter.Contains(core.Position),
			filter.Contains(core.Speed),
		),
	).Each(ms.world, func(entry *donburi.Entry) {
		pos := core.Position.Get(entry)
		speed := core.Speed.Get(entry)

		// Calculate screen center
		centerX := float64(ms.config.ScreenSize.Width) / 2
		centerY := float64(ms.config.ScreenSize.Height) / 2

		// Calculate direction from center (normalize the vector)
		dx := pos.X - centerX
		dy := pos.Y - centerY
		distance := math.Sqrt(dx*dx + dy*dy)

		// Handle edge case where star is exactly at center
		if distance == 0 {
			// Give it a random direction
			//nolint:gosec // math/rand is fine for non-crypto game randomness
			angle := rand.Float64() * 2 * math.Pi
			dx = math.Cos(angle)
			dy = math.Sin(angle)
		} else {
			// Normalize the direction vector
			dx /= distance
			dy /= distance
		}

		// Move star outward from center
		movementDistance := *speed * deltaTime * 5.0 // 5x faster for more engaging starfield
		pos.X += dx * movementDistance
		pos.Y += dy * movementDistance

		// FIXED: Check if star has moved off-screen (not distance from center!)
		screenWidth := float64(ms.config.ScreenSize.Width)
		screenHeight := float64(ms.config.ScreenSize.Height)

		// Reset star when it goes off any edge of the screen
		if pos.X < 0 || pos.X > screenWidth || pos.Y < 0 || pos.Y > screenHeight {
			// Reset star to a random position near the center
			ms.resetStarPosition(entry, centerX, centerY)
			ms.logger.Debug("Star reset to new position",
				"entity_id", entry.Entity().String(),
				"new_position", pos)
		}

		// Update scale based on distance from center (stars get bigger as they move outward)
		if entry.HasComponent(core.Scale) {
			scale := core.Scale.Get(entry)
			// Scale from min to max based on distance
			newDistance := math.Sqrt((pos.X-centerX)*(pos.X-centerX) + (pos.Y-centerY)*(pos.Y-centerY))
			scaleFactor := math.Min(1.0, newDistance/ms.config.StarScaleDistance)
			*scale = ms.config.StarMinScale + (ms.config.StarMaxScale-ms.config.StarMinScale)*scaleFactor
		}
	})
}

// resetStarPosition resets a star to a new random position near the center
func (ms *MovementSystem) resetStarPosition(entry *donburi.Entry, centerX, centerY float64) {
	pos := core.Position.Get(entry)

	// Generate random angle
	//nolint:gosec // math/rand is fine for non-crypto game randomness
	angle := rand.Float64() * 2 * math.Pi

	// Generate random radius within spawn range
	//nolint:gosec // math/rand is fine for non-crypto game randomness
	spawnRadius := ms.config.StarSpawnRadiusMin +
		rand.Float64()*(ms.config.StarSpawnRadiusMax-ms.config.StarSpawnRadiusMin)

	// Set new position
	pos.X = centerX + math.Cos(angle)*spawnRadius
	pos.Y = centerY + math.Sin(angle)*spawnRadius

	// Reset scale to minimum
	if entry.HasComponent(core.Scale) {
		scale := core.Scale.Get(entry)
		*scale = ms.config.StarMinScale
	}

	ms.logger.Debug("Star reset to new position",
		"entity_id", entry.Entity().String(),
		"new_position", pos)
}
