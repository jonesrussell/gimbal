package core

import "errors"

var (
	ErrInvalidFont        = errors.New("invalid font type: expected textv2.Face")
	ErrInvalidHeartSprite = errors.New("invalid heart sprite: expected *ebiten.Image")
	ErrNilContainer       = errors.New("container cannot be nil")
	ErrInvalidLives       = errors.New("lives must be non-negative")
	ErrInvalidScore       = errors.New("score must be non-negative")
	ErrInvalidHealth      = errors.New("health must be between 0 and 1")
	ErrInvalidAmmo        = errors.New("ammo must be non-negative")
)
