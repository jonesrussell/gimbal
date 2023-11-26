package main

import (
	"log"

	"github.com/jonesrussell/gimbal/internal/game"
)

func main() {
	g := game.NewGimlarGame(0.04) // Pass the speed as an argument
	if err := g.Run(); err != nil {
		log.Fatal(err)
	}
}
