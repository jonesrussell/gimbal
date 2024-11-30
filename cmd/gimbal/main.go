package main

import (
	"log/slog"

	"go.uber.org/zap"

	"github.com/jonesrussell/gimbal/game"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	speed := 0.04
	g, err := game.NewGimlarGame(speed)
	if err != nil {
		slog.Error("Failed to initialize game", "error", err)
	}

	if err := g.Run(); err != nil {
		slog.Error("Failed to run game", "error", err)
	}
}
