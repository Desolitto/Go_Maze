package main

import (
	"bufio"
	"flag"
	"fmt"
	"go-maze/pkg/cave"
	"go-maze/pkg/game"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten"
	"golang.org/x/exp/rand"
)

// Структура лабиринта
type Maze struct {
	Rows              int
	Cols              int
	Cave              *cave.Cave
	RightBorderMatrix [][]int
	LowBorderMatrix   [][]int
}

func (m *Maze) LoadCaveMazeFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	if scanner.Scan() {
		dimensions := strings.Fields(scanner.Text())
		if len(dimensions) != 2 {
			log.Fatal("Неверный формат файла: первая строка должна содержать размеры пещеры.")
		}

		width, err := strconv.Atoi(dimensions[0])
		if err != nil {
			log.Fatal("Неверная ширина лабиринта.")
		}

		height, err := strconv.Atoi(dimensions[1])
		if err != nil {
			log.Fatal("Неверная высота лабиринта.")
		}

		m.Rows, m.Cols = height, width
		m.Cave = cave.NewCave(width, height)
		// Загрузка первой матрицы (стена справа)
		m.RightBorderMatrix = make([][]int, height)

		for y := 0; y < height; y++ {
			if scanner.Scan() {
				row := strings.Fields(scanner.Text())
				if len(row) != width {
					log.Fatal("Неверный формат файла: количество столбцов не совпадает с заданной шириной.")
				}
				m.RightBorderMatrix[y] = make([]int, width)
				for x, cell := range row {
					if cell == "0" {
						m.Cave.Grid[y][x] = cave.Death
					} else if cell == "1" {
						m.Cave.Grid[y][x] = cave.Alive
					} else {
						log.Fatal("Неверный символ в пещере: должен быть 0 или 1.")
					}
				}
			}
		}
		m.LowBorderMatrix = make([][]int, height)
		for y := 0; y < height; y++ {
			if scanner.Scan() {
				row := strings.Fields(scanner.Text())
				if len(row) != width {
					log.Fatal("Неверный формат файла: количество столбцов не совпадает с заданной шириной.")
				}
				m.LowBorderMatrix[y] = make([]int, width)
				for x, cell := range row {
					value, err := strconv.Atoi(cell)
					if err != nil {
						log.Fatal("Неверный символ в матрице стен: должен быть целым числом.")
					}
					m.LowBorderMatrix[y][x] = value
				}
			}
		}
	}
}

func (m *Maze) PrintMaze() {
	for y := 0; y < m.Rows; y++ {
		for x := 0; x < m.Cols; x++ {
			if m.Cave.Grid[y][x] == cave.Alive {
				fmt.Print("1 ")
			} else {
				fmt.Print("0 ")
			}
		}
		fmt.Println()
	}
}

// func main() {
// 	m := &Maze{}
// 	m.LoadCaveMazeFile("/Users/calamarp/Desktop/go/Go_Maze/src/assets/maze_test.txt")
// 	m.PrintMaze()
// }

func main() {
	rand.Seed(uint64(time.Now().UnixNano()))
	w := flag.Int("w", 20, "width of the cave")
	h := flag.Int("h", 20, "height of the cave")
	initialChance := flag.Int("с", 55, "initial chance (0-100)")
	flag.Parse()

	game := game.NewGame(*w, *h, *initialChance)
	fmt.Println("Исходная матрица:")
	// game.PrintCave()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
