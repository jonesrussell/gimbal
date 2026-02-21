package common

import (
	"context"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/jonesrussell/gimbal/internal/math"
)

// GameUI represents the game user interface system.
// It manages HUD updates, score display, lives, and pause menu visibility.
type GameUI interface {
	Update() error
	Draw(screen *ebiten.Image)
	UpdateScore(score int)
	UpdateLives(lives int)
	ShowPauseMenu(visible bool)
	SetDeviceClass(deviceClass string)
}

// Logger represents a logging interface
type Logger interface {
	Debug(msg string, fields ...any)
	Info(msg string, fields ...any)
	Warn(msg string, fields ...any)
	Error(msg string, fields ...any)
	DebugContext(ctx context.Context, msg string, fields ...any)
	InfoContext(ctx context.Context, msg string, fields ...any)
	WarnContext(ctx context.Context, msg string, fields ...any)
	ErrorContext(ctx context.Context, msg string, fields ...any)
	Sync() error
}

// Entity represents any game object that can be updated and drawn
type Entity interface {
	Update()
	Draw(screen any)
	GetPosition() Point
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

// RenderSystem represents a system that handles rendering
type RenderSystem interface {
	Draw(screen *ebiten.Image, sprite *ebiten.Image, pos Point, angle math.Angle)
	DrawDebug(screen *ebiten.Image)
}

// HealthProvider provides access to player health information
// This interface enables type-safe health system access across packages
type HealthProvider interface {
	GetPlayerHealth() (current, maximum int)
}
