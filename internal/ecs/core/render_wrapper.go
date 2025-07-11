package core

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"

	"github.com/jonesrussell/gimbal/internal/common"
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
		return common.NewGameError(common.ErrorCodeRenderingFailed, "screen is nil")
	}
	RenderSystem(world, rsw.screen)
	return nil
}

func (rsw *RenderSystemWrapper) Name() string {
	return "RenderSystem"
}
