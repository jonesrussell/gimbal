package common

import (
	"math"
)

// Point represents a 2D point in the game world
type Point struct {
	X, Y float64
}

// Size represents dimensions
type Size struct {
	Width, Height int
}

// Angle represents an angle in degrees
type Angle float64

// ToRadians converts the angle from degrees to radians
func (a Angle) ToRadians() float64 {
	return float64(a) * (math.Pi / 180)
}

// FromRadians creates an Angle from radians
func FromRadians(rad float64) Angle {
	return Angle(rad * (180 / math.Pi))
}

// Add returns the sum of two angles
func (a Angle) Add(b Angle) Angle {
	return a + b
}

// Sub returns the difference between two angles
func (a Angle) Sub(b Angle) Angle {
	return a - b
}

// Mul returns the product of an angle and a scalar
func (a Angle) Mul(scalar float64) Angle {
	return Angle(float64(a) * scalar)
}

// Div returns the quotient of an angle and a scalar
func (a Angle) Div(scalar float64) Angle {
	return Angle(float64(a) / scalar)
}

// Normalize returns the angle normalized to the range [0, 360)
func (a Angle) Normalize() Angle {
	angle := float64(a)
	for angle < 0 {
		angle += 360
	}
	for angle >= 360 {
		angle -= 360
	}
	return Angle(angle)
}

// GameState represents the current state of the game
type GameState struct {
	Center     Point
	ScreenSize Size
	Debug      bool
}

// EntityConfig holds common configuration for game entities
type EntityConfig struct {
	Position Point
	Size     Size
	Radius   float64
	Speed    float64
}
