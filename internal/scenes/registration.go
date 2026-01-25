// Package scenes provides scene management for the game.
// This file documents the explicit scene registration approach.
package scenes

// SceneRegistrationDoc documents the registration approach.
//
// For explicit scene registration, import and call Register() from each scene package:
//
//	import (
//		"github.com/jonesrussell/gimbal/internal/scenes/gameover"
//		"github.com/jonesrussell/gimbal/internal/scenes/gameplay"
//		"github.com/jonesrussell/gimbal/internal/scenes/intro"
//		"github.com/jonesrussell/gimbal/internal/scenes/mainmenu"
//		"github.com/jonesrussell/gimbal/internal/scenes/pause"
//	)
//
//	func registerScenes() {
//		intro.Register()
//		mainmenu.Register()
//		gameplay.Register()
//		pause.Register()
//		gameover.Register()
//	}
//
// Alternatively, use blank imports to trigger init() auto-registration:
//
//	import (
//		_ "github.com/jonesrussell/gimbal/internal/scenes/gameover"
//		_ "github.com/jonesrussell/gimbal/internal/scenes/gameplay"
//		// ...
//	)
const SceneRegistrationDoc = "See package documentation for registration patterns"
