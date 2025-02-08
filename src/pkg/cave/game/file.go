package game

import (
	"bufio"
	"go-maze/pkg/cave"
	"log"
	"os"
	"strconv"
	"strings"
)

func (g *Game) LoadCaveFromFile(filename string) {
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
		if err != nil || width > maxSize {
			log.Fatal("Неверная ширина пещеры.")
		}

		height, err := strconv.Atoi(dimensions[1])
		if err != nil || height > maxSize {
			log.Fatal("Неверная высота пещеры.")
		}

		g.width, g.height = width, height
		g.cellSize = float32(sceneWidth) / float32(width)
		g.cave = cave.NewCave(width, height)

		for y := 0; y < height; y++ {
			if scanner.Scan() {
				row := strings.Fields(scanner.Text())
				if len(row) != width {
					log.Fatal("Неверный формат файла: количество столбцов не совпадает с заданной шириной.")
				}
				for x, cell := range row {
					if cell == "0" {
						g.cave.Grid[y][x] = cave.Death
					} else if cell == "1" {
						g.cave.Grid[y][x] = cave.Alive
					} else {
						log.Fatal("Неверный символ в пещере: должен быть 0 или 1.")
					}
				}
			}
		}
	}
	g.autoStepActive = false
}
