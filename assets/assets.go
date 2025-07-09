// Package assets provides embedded game assets such as sprites and sounds.
package assets

import "embed"

//go:embed sprites/*
var Assets embed.FS
