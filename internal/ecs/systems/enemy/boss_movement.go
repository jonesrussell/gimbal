package enemy

import (
	stdmath "math"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/jonesrussell/gimbal/internal/ecs/core"
	"github.com/jonesrussell/gimbal/internal/math"
)

// UpdateBossMovement updates the boss's orbital movement
func (es *EnemySystem) UpdateBossMovement(deltaTime float64) {
	query.NewQuery(
		filter.And(
			filter.Contains(core.EnemyTag),
			filter.Contains(core.Orbital),
		),
	).Each(es.world, func(entry *donburi.Entry) {
		orbital := core.Orbital.Get(entry)
		health := core.Health.Get(entry)

		// Check if this is a boss (has orbital movement and high health)
		if health != nil && health.Maximum >= 10 {
			// Update orbital angle (convert speed from radians/sec to degrees/sec)
			angleDelta := math.Angle(BossOrbitalSpeed * deltaTime * 180.0 / stdmath.Pi)
			orbital.OrbitalAngle += angleDelta

			// Keep angle in 0-360 range
			orbital.OrbitalAngle = orbital.OrbitalAngle.Normalize()

			// Update position based on orbital angle
			radians := orbital.OrbitalAngle.ToRadians()
			pos := core.Position.Get(entry)
			pos.X = orbital.Center.X + stdmath.Cos(radians)*orbital.Radius
			pos.Y = orbital.Center.Y + stdmath.Sin(radians)*orbital.Radius

			// Update orbital component
			core.Orbital.SetValue(entry, *orbital)
		}
	})
}
