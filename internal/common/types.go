package common

import "github.com/jonesrussell/gimbal/internal/config"

// Point represents a 2D point in the game world
type Point struct {
	X, Y float64
}

// GameState represents the current state of the game
type GameState struct {
	Center     Point
	ScreenSize config.Size
	Debug      bool
}

// EntityConfig holds common configuration for game entities
type EntityConfig struct {
	Position Point
	Size     config.Size
	Radius   float64
	Speed    float64
}
