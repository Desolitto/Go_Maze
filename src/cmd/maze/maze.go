package main

import (
	"flag"
	"fmt"
	"log"

	"go-maze/config"
	"go-maze/pkg/maze"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	w := flag.Int("w", config.MaxSize, "количество строк в лабиринте")
	h := flag.Int("h", config.MaxSize, "количество столбцов в лабиринте")
	flag.Parse()

	game := maze.NewGame(*w, *h)
	fmt.Println("Сгенерированный лабиринт:")

	err := game.Maze().SaveMaze("maze.txt")

	if err != nil {
		fmt.Println("Ошибка при сохранении лабиринта:", err)
	} else {
		fmt.Println("Лабиринт успешно сохранен в maze.txt")
	}
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
