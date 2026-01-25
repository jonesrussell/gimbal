// Package value provides domain value objects used across the game.
// These types are pure domain concepts with no framework dependencies.
package value

// MovementPattern represents the movement behavior pattern for entities.
// This is a domain value object that defines how entities move through space.
type MovementPattern int

const (
	// MovementPatternNormal is standard outward movement
	MovementPatternNormal MovementPattern = iota
	// MovementPatternZigzag oscillates side-to-side while moving outward
	MovementPatternZigzag
	// MovementPatternAccelerating starts slow and speeds up
	MovementPatternAccelerating
	// MovementPatternPulsing moves in bursts (fast-slow-fast)
	MovementPatternPulsing
)

// String returns a human-readable string representation of the movement pattern.
func (mp MovementPattern) String() string {
	switch mp {
	case MovementPatternNormal:
		return "Normal"
	case MovementPatternZigzag:
		return "Zigzag"
	case MovementPatternAccelerating:
		return "Accelerating"
	case MovementPatternPulsing:
		return "Pulsing"
	default:
		return "Unknown"
	}
}

// IsValid returns true if the movement pattern is a known valid value.
func (mp MovementPattern) IsValid() bool {
	return mp >= MovementPatternNormal && mp <= MovementPatternPulsing
}

// ParseMovementPattern converts an integer to a MovementPattern.
// Returns MovementPatternNormal and false if the value is invalid.
func ParseMovementPattern(v int) (MovementPattern, bool) {
	mp := MovementPattern(v)
	if !mp.IsValid() {
		return MovementPatternNormal, false
	}
	return mp, true
}
