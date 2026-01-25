package bossintro

import (
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/jonesrussell/gimbal/internal/ecs/managers"
	resources "github.com/jonesrussell/gimbal/internal/ecs/managers/resource"
	"github.com/jonesrussell/gimbal/internal/scenes"
)

// Register registers the boss intro overlay (not a full scene, but a component)
// This is kept for consistency, though boss intro is an overlay, not a scene
func Register() {
	// Boss intro is an overlay, not a scene, so we don't register it as a scene
	// It's used directly from the playing scene
}

func init() {
	Register()
}

// Note: Boss intro is an overlay component, not a full scene
// It's used directly from the playing scene
