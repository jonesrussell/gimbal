// Package core provides core game engine types and interfaces
package core

import (
	"github.com/jonesrussell/gimbal/internal/core/types"
)

// Re-export all types
type (
	// Interfaces
	Game         = types.Game
	InputHandler = types.InputHandler
	Renderer     = types.Renderer
	AssetManager = types.AssetManager
	System       = types.System
	Entity       = types.Entity

	// Constructors
	NewGame         = types.NewGame
	NewInputHandler = types.NewInputHandler
	NewRenderer     = types.NewRenderer
	NewAssetManager = types.NewAssetManager

	// Options
	GameOption   = types.GameOption
	InputOption  = types.InputOption
	RenderOption = types.RenderOption
	AssetOption  = types.AssetOption

	// Configs
	GameConfig         = types.GameConfig
	InputConfig        = types.InputConfig
	RenderConfig       = types.RenderConfig
	AssetManagerConfig = types.AssetManagerConfig
)
