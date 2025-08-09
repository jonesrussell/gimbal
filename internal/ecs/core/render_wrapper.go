package core

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"

	"github.com/jonesrussell/gimbal/internal/errors"
)

type RenderSystemWrapper struct {
	screen *ebiten.Image
}

func NewRenderSystemWrapper(screen *ebiten.Image) *RenderSystemWrapper {
	return &RenderSystemWrapper{
		screen: screen,
	}
}

func (rsw *RenderSystemWrapper) Update(world donburi.World, args ...interface{}) error {
	if rsw.screen == nil {
		return errors.NewGameError(errors.RenderFailed, "screen is nil")
	}
	RenderSystem(world, rsw.screen)
	return nil
}

// Name() method removed - no longer needed since SystemManager is unused
