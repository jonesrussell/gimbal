package main

import (
	"log/slog"

	"github.com/jonesrussell/gimbal/internal/game"
)

func main() {
	speed := 0.04
	g, err := game.NewGimlarGame(speed)
	if err != nil {
		slog.Error("Failed to initialize game", "error", err)
	}

	if err := g.Run(); err != nil {
		slog.Error("Failed to run game", "error", err)
	}
}
