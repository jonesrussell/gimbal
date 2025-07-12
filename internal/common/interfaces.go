package common

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/jonesrussell/gimbal/internal/config"
	"github.com/jonesrussell/gimbal/internal/math"
)

type GameUI interface {
	Update() error
	Draw(screen *ebiten.Image)
	UpdateScore(score int)
	UpdateLives(lives int)
	ShowPauseMenu(visible bool)
	SetDeviceClass(deviceClass string)
}

type UIData struct {
	Score       int
	Lives       int
	IsPaused    bool
	DeviceClass string
}

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
	SetAngle(angle math.Angle)
	GetAngle() math.Angle
}

// Collidable represents an entity that can collide with others
type Collidable interface {
	Entity
	GetBounds() config.Size
	CheckCollision(other Collidable) bool
}

// InputHandler represents a component that can handle input
type InputHandler interface {
	HandleInput()
	IsKeyPressed(key ebiten.Key) bool
}

// MovementInputHandler handles movement-specific input
type MovementInputHandler interface {
	GetMovementInput() math.Angle
}

// ActionInputHandler handles action-specific input (pause, shoot, quit)
type ActionInputHandler interface {
	IsQuitPressed() bool
	IsPausePressed() bool
	IsShootPressed() bool
}

// TouchInputHandler handles touch-specific input
type TouchInputHandler interface {
	GetTouchState() *TouchState
}

// MouseInputHandler handles mouse-specific input
type MouseInputHandler interface {
	GetMousePosition() Point
	IsMouseButtonPressed(button ebiten.MouseButton) bool
}

// EventInputHandler handles input event tracking
type EventInputHandler interface {
	GetLastEvent() InputEvent
}

// TestableInputHandler provides simulation methods for testing
type TestableInputHandler interface {
	SimulateKeyPress(key ebiten.Key)
	SimulateKeyRelease(key ebiten.Key)
}

// GameInputHandler represents the main input interface for the game
// This interface should be implemented by input adapters and used by the ECS system
// It composes smaller interfaces for better separation of concerns
type GameInputHandler interface {
	InputHandler
	MovementInputHandler
	ActionInputHandler
	TouchInputHandler
	MouseInputHandler
	EventInputHandler
	TestableInputHandler
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
	InputEventAny // Any key or mouse input - perfect for title screens
)

// PhysicsSystem represents a system that handles physics calculations
type PhysicsSystem interface {
	CalculatePosition(angle math.Angle, radius float64) Point
	ValidatePosition(pos Point, bounds config.Size) Point
}

// RenderSystem represents a system that handles rendering
type RenderSystem interface {
	Draw(screen *ebiten.Image, sprite *ebiten.Image, pos Point, angle math.Angle)
	DrawDebug(screen *ebiten.Image)
}

// Result represents a value that can either be a success or an error
type Result[T any] struct {
	value T
	err   error
}

// Ok creates a successful Result
func Ok[T any](value T) Result[T] {
	return Result[T]{value: value, err: nil}
}

// Err creates a failed Result
func Err[T any](err error) Result[T] {
	return Result[T]{err: err}
}

// Unwrap returns the value or panics if there's an error
func (r Result[T]) Unwrap() T {
	if r.err != nil {
		panic(r.err)
	}
	return r.value
}

// UnwrapOr returns the value or a default if there's an error
func (r Result[T]) UnwrapOr(defaultValue T) T {
	if r.err != nil {
		return defaultValue
	}
	return r.value
}

// IsOk returns true if the Result is successful
func (r Result[T]) IsOk() bool {
	return r.err == nil
}

// IsErr returns true if the Result contains an error
func (r Result[T]) IsErr() bool {
	return r.err != nil
}

// Error returns the error if present
func (r Result[T]) Error() error {
	return r.err
}

// Value returns the value if successful
func (r Result[T]) Value() T {
	return r.value
}

// Updatable represents an entity that can be updated with a delta time
type Updatable interface {
	Update(deltaTime float64) error
}

// Drawable represents an entity that can be drawn to a screen
type Drawable interface {
	Draw(screen *ebiten.Image) error
}

// Identifiable represents an entity with a unique identifier
type Identifiable interface {
	GetID() string
	SetID(id string)
}

// Configurable represents an entity that can be configured
type Configurable[T any] interface {
	GetConfig() T
	SetConfig(config T) error
	ValidateConfig(config T) error
}

// Lifecycle represents an entity with lifecycle management
type Lifecycle interface {
	Initialize() error
	Start() error
	Stop() error
	Cleanup() error
}

// Observable represents an entity that can notify observers of changes
type Observable[T any] interface {
	Subscribe(observer func(T)) func() // Returns unsubscribe function
	Notify(event T)
}

// Repository represents a data access interface
type Repository[T any, ID comparable] interface {
	Get(id ID) (T, error)
	GetAll() ([]T, error)
	Create(entity T) error
	Update(entity T) error
	Delete(id ID) error
	Exists(id ID) bool
}
