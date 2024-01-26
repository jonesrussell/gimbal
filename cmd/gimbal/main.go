package main

import (
	"log"

	"github.com/jonesrussell/gimbal/internal/game"
)

func main() {
	g, err := game.NewGimlarGame(0.04) // Pass the speed as an argument
	if err != nil {
		log.Fatal((err))
	}

	if err := g.Run(); err != nil {
		log.Fatal(err)
	}
}
