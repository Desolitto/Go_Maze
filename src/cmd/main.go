package main

import (
	"flag"
	"log"

	"go-maze/pkg/game"

	// "github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	w := flag.Int("w", 20, "width of the window")
	h := flag.Int("h", 20, "height of the window")
	flag.Parse()
	game := game.NewGame(*w, *h)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
