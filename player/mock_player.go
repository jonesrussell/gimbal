package player

import "github.com/hajimehoshi/ebiten/v2"

type MockImage struct {
	drawn bool
}

func (mi *MockImage) DrawImage(_ *ebiten.Image, _ *ebiten.DrawImageOptions) {
	mi.drawn = true
}

// Implement other necessary methods...
