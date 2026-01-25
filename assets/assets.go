// Package assets provides embedded game assets such as sprites and sounds.
package assets

import "embed"

//go:embed sprites/*.png fonts/*.ttf entities/*.json sounds/*.ogg stages/*.json levels/*.json ui/*.png cutscenes/*.png cutscenes/*.json ending/*.png planets/*.png bosses/*.png
var Assets embed.FS
