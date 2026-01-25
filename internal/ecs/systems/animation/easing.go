package animation

import (
	"math"

	"github.com/jonesrussell/gimbal/internal/ecs/core"
)

// ApplyEasing applies the specified easing function to a progress value
func ApplyEasing(progress float64, easing core.EasingType) float64 {
	// Clamp progress to [0, 1]
	if progress < 0 {
		progress = 0
	}
	if progress > 1 {
		progress = 1
	}

	switch easing {
	case core.EasingEaseIn:
		return easeInQuad(progress)
	case core.EasingEaseOut:
		return easeOutQuad(progress)
	case core.EasingEaseInOut:
		return easeInOutQuad(progress)
	default: // EasingLinear
		return progress
	}
}

// easeInQuad - quadratic ease in (slow start)
func easeInQuad(t float64) float64 {
	return t * t
}

// easeOutQuad - quadratic ease out (slow end)
func easeOutQuad(t float64) float64 {
	return 1 - (1-t)*(1-t)
}

// easeInOutQuad - quadratic ease in and out
func easeInOutQuad(t float64) float64 {
	if t < 0.5 {
		return 2 * t * t
	}
	return 1 - math.Pow(-2*t+2, 2)/2
}

// easeInCubic - cubic ease in
func easeInCubic(t float64) float64 {
	return t * t * t
}

// easeOutCubic - cubic ease out
func easeOutCubic(t float64) float64 {
	return 1 - math.Pow(1-t, 3)
}

// easeInOutCubic - cubic ease in and out
func easeInOutCubic(t float64) float64 {
	if t < 0.5 {
		return 4 * t * t * t
	}
	return 1 - math.Pow(-2*t+2, 3)/2
}

// Lerp performs linear interpolation between two values
func Lerp(start, end, t float64) float64 {
	return start + (end-start)*t
}

// LerpWithEasing performs interpolation with easing
func LerpWithEasing(start, end, t float64, easing core.EasingType) float64 {
	easedT := ApplyEasing(t, easing)
	return Lerp(start, end, easedT)
}
