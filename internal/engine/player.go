package engine

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Player struct {
	image *ebiten.Image
	x, y  float64
	speed float64
}

func NewPlayer(image *ebiten.Image, x, y float64) *Player {
	return &Player{
		image: image,
		x:     x,
		y:     y,
		speed: 4.0,
	}
}

func (p *Player) Update() {
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		p.x -= p.speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		p.x += p.speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		p.y -= p.speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		p.y += p.speed
	}
}

func (p *Player) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(p.x, p.y)
	screen.DrawImage(p.image, op)
}
