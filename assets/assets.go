// Package assets provides embedded game assets such as sprites and sounds.
package assets

import "embed"

//go:embed sprites/* fonts/* entities/*
var Assets embed.FS
