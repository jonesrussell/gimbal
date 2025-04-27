package game_test

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

func TestGame_Draw(t *testing.T) {
	t.Run("nil screen", func(t *testing.T) {
		g := newTestGame(t)
		g.Draw(nil)
	})

	t.Run("normal draw", func(t *testing.T) {
		g := newTestGame(t)
		screen := ebiten.NewImage(800, 600)
		g.Draw(screen)
	})
}
