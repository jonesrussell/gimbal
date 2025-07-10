// Package assets provides embedded game assets such as sprites and sounds.
package assets

import "embed"

//go:embed sprites/* fonts/*
var Assets embed.FS
