package health

import (
	"github.com/yohamta/donburi"

	"github.com/jonesrussell/gimbal/internal/common"
	"github.com/jonesrussell/gimbal/internal/ecs/core"
	"github.com/jonesrussell/gimbal/internal/math"
)

// respawnPlayer respawns the player at the center bottom of the screen
func (hs *HealthSystem) respawnPlayer(playerEntity donburi.Entity) {
	playerEntry := hs.world.Entry(playerEntity)
	if !playerEntry.Valid() {
		return
	}

	// Reset position to center bottom
	center := common.Point{
		X: float64(hs.gameConfig.ScreenSize.Width) / 2,
		Y: float64(hs.gameConfig.ScreenSize.Height) / 2,
	}

	// Update position
	core.Position.SetValue(playerEntry, center)

	// Reset orbital data to bottom position (180 degrees)
	orbitalData := core.Orbital.Get(playerEntry)
	orbitalData.Center = center
	orbitalData.OrbitalAngle = 180 // 180 degrees
	core.Orbital.SetValue(playerEntry, *orbitalData)

	// Reset angle
	core.Angle.SetValue(playerEntry, math.Angle(0))

	hs.logger.Debug("Player respawned at center bottom")
}
