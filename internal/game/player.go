package game

import (
	"github.com/jonesrussell/gimbal/player"
)

type PlayerInput struct {
	input player.InputHandlerInterface
}

const (
	AngleStep = player.AngleStep
)

func (g *Game) Update() error {
	g.player.UpdatePosition()
	return nil
}
