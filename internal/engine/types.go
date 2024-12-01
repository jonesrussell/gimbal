package engine

import "github.com/hajimehoshi/ebiten/v2"

// System represents a game subsystem
type System interface {
	Update() error
}

// Drawable represents anything that can be drawn to the screen
type Drawable interface {
	Draw(screen *ebiten.Image)
}

// Entity represents a game object with position and state
type Entity interface {
	System
	Drawable
	GetID() uint64
}
