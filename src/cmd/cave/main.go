package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"go-maze/pkg/cave/game"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/exp/rand"
)

func main() {
	rand.Seed(uint64(time.Now().UnixNano()))
	w := flag.Int("w", 20, "width of the cave")
	h := flag.Int("h", 20, "height of the cave")
	birthLimit := flag.Int("b", 4, "birth limit (0-7)")
	deathLimit := flag.Int("d", 3, "death limit (0-7)")
	initialChance := flag.Int("с", 55, "initial chance (0-100)")
	flag.Parse()

	game := game.NewGame(*w, *h, *birthLimit, *deathLimit, *initialChance)
	fmt.Println("Исходная матрица:")
	game.PrintCave()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
