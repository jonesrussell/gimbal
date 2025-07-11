package common

import "github.com/hajimehoshi/ebiten/v2"

// Logger represents a logging interface
type Logger interface {
	Debug(msg string, fields ...any)
	Info(msg string, fields ...any)
	Warn(msg string, fields ...any)
	Error(msg string, fields ...any)
	Sync() error
}

// Entity represents any game object that can be updated and drawn
type Entity interface {
	Update()
	Draw(screen any)
	GetPosition() Point
}

// Movable represents an entity that can move
type Movable interface {
	Entity
	SetPosition(pos Point)
	GetSpeed() float64
}

// Rotatable represents an entity that can rotate
type Rotatable interface {
	Movable
	SetAngle(angle Angle)
	GetAngle() Angle
}

// Collidable represents an entity that can collide with others
type Collidable interface {
	Entity
	GetBounds() Size
	CheckCollision(other Collidable) bool
}

// InputHandler represents a component that can handle input
type InputHandler interface {
	HandleInput()
	IsKeyPressed(key ebiten.Key) bool
}

// GameInputHandler represents the main input interface for the game
// This interface should be implemented by input adapters and used by the ECS system
type GameInputHandler interface {
	HandleInput()
	IsKeyPressed(key ebiten.Key) bool
	GetMovementInput() Angle
	IsQuitPressed() bool
	IsPausePressed() bool
	IsShootPressed() bool
	GetTouchState() *TouchState
	GetMousePosition() Point
	IsMouseButtonPressed(button ebiten.MouseButton) bool
	GetLastEvent() InputEvent
	// Simulation methods for testing
	SimulateKeyPress(key ebiten.Key)
	SimulateKeyRelease(key ebiten.Key)
}

// TouchState tracks touch input state
type TouchState struct {
	ID       ebiten.TouchID
	StartPos Point
	LastPos  Point
	Duration int
}

// InputEvent represents a game input event
type InputEvent int

const (
	InputEventNone InputEvent = iota
	InputEventMove
	InputEventPause
	InputEventQuit
	InputEventTouch
	InputEventMouseMove
)

// PhysicsSystem represents a system that handles physics calculations
type PhysicsSystem interface {
	CalculatePosition(angle Angle, radius float64) Point
	ValidatePosition(pos Point, bounds Size) Point
}

// RenderSystem represents a system that handles rendering
type RenderSystem interface {
	Draw(screen *ebiten.Image, sprite *ebiten.Image, pos Point, angle Angle)
	DrawDebug(screen *ebiten.Image)
}
