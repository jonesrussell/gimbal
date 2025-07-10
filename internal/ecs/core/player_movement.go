package core

import "github.com/jonesrussell/gimbal/internal/common"

// PlayerMovement handles all player movement logic
type PlayerMovement struct{}

// UpdateOrbitalAngle updates player's orbital angle based on input
func (pm *PlayerMovement) UpdateOrbitalAngle(orb *OrbitalData, inputAngle common.Angle) {
	if inputAngle != 0 {
		orb.OrbitalAngle += inputAngle
		pm.normalizeAngle(&orb.OrbitalAngle)
	}
}

// UpdateFacingAngle calculates facing angle for Gyruss-style gameplay
func (pm *PlayerMovement) UpdateFacingAngle(orb *OrbitalData) {
	orb.FacingAngle = orb.OrbitalAngle + common.FacingAngleOffset
	pm.normalizeAngle(&orb.FacingAngle)
}

// normalizeAngle keeps angle in [0, 360) range
func (pm *PlayerMovement) normalizeAngle(angle *common.Angle) {
	if *angle < 0 {
		*angle += common.FullCircleDegrees
	} else if *angle >= common.FullCircleDegrees {
		*angle -= common.FullCircleDegrees
	}
}
