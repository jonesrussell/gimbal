package common

// Point represents a 2D point in the game world
type Point struct {
	X, Y float64
}

// Size represents dimensions
type Size struct {
	Width, Height int
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
