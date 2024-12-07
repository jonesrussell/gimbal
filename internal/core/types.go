package core

import (
	"context"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"go.uber.org/zap"
)

// Game represents the main game interface
type Game interface {
	// Update handles game logic updates
	Update(ctx context.Context) error
	// Draw handles rendering
	Draw(screen *ebiten.Image) error
	// Layout returns the game's screen dimensions
	Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int)
	// Shutdown performs cleanup when the game exits
	Shutdown(ctx context.Context) error
}

// InputHandler manages all input processing
type InputHandler interface {
	// Process handles input events
	Process(ctx context.Context) error
	// IsKeyPressed returns true if the given key is currently pressed
	IsKeyPressed(key ebiten.Key) bool
	// RegisterBinding adds a new key binding
	RegisterBinding(name string, key ebiten.Key, handler func(ctx context.Context) error) error
}

// Renderer handles all drawing operations
type Renderer interface {
	// Draw renders the current frame
	Draw(screen *ebiten.Image) error
	// AddSprite adds a sprite to the render queue
	AddSprite(id string, sprite *ebiten.Image, pos image.Point, z int) error
	// Clear removes all sprites from the render queue
	Clear()
}

// AssetManager handles loading and managing game assets
type AssetManager interface {
	// LoadImage loads an image asset
	LoadImage(ctx context.Context, path string) (*ebiten.Image, error)
	// LoadSound loads a sound asset
	LoadSound(ctx context.Context, path string) ([]byte, error)
	// Preload loads all assets in the specified directory
	Preload(ctx context.Context, dir string) error
	// Cleanup frees resources
	Cleanup(ctx context.Context) error
}

// System represents a game subsystem (physics, AI, etc.)
type System interface {
	// Init initializes the system
	Init(ctx context.Context) error
	// Update updates the system state
	Update(ctx context.Context) error
	// Cleanup performs system cleanup
	Cleanup(ctx context.Context) error
}

// Entity represents a game object
type Entity interface {
	// ID returns the entity's unique identifier
	ID() string
	// Update updates the entity state
	Update(ctx context.Context) error
	// Draw renders the entity
	Draw(screen *ebiten.Image) error
	// Bounds returns the entity's bounding box
	Bounds() image.Rectangle
}

// Constructor types for dependency injection
type (
	// NewGame creates a new game instance
	NewGame func(logger *zap.Logger) (Game, error)

	// NewInputHandler creates a new input handler
	NewInputHandler func(logger *zap.Logger) (InputHandler, error)

	// NewRenderer creates a new renderer
	NewRenderer func(logger *zap.Logger, opts ...RenderOption) (Renderer, error)

	// NewAssetManager creates a new asset manager
	NewAssetManager func(logger *zap.Logger, opts ...AssetOption) (AssetManager, error)

	// AssetOption configures the asset manager
	AssetOption func(*AssetManagerConfig)
)

// Options for configurable components using the functional options pattern
type (
	// GameOption configures a Game instance
	GameOption func(*GameConfig)

	// InputOption configures an InputHandler instance
	InputOption func(*InputConfig)

	// RenderOption configures a Renderer instance
	RenderOption func(*RenderConfig)
)

// Configuration structs
type (
	GameConfig struct {
		ScreenWidth  int  // 8 bytes
		ScreenHeight int  // 8 bytes
		Debug        bool // 1 byte, padded to 8 bytes
	}

	InputConfig struct {
		EnableKeyboard bool
		EnableMouse    bool
		EnableGamepad  bool
	}

	RenderConfig struct {
		VSync        bool
		DoubleBuffer bool
		MaxSprites   int
		Debug        bool
	}

	// AssetManagerConfig holds the configuration for the asset manager
	AssetManagerConfig struct {
		BaseDir     string
		CacheSize   int
		EnableSound bool
		Debug       bool
	}
)
