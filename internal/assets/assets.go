package assets

import (
	"bytes"
	_ "embed"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed player.png
var data []byte

func LoadPlayerSprite() *ebiten.Image {
	m, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}

	return ebiten.NewImageFromImage(m)
}
