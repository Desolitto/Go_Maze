package main

import (
	"bufio"
	"flag"
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	maxMazeSize     = 50
	wallThickness   = 2
	mazeWidth       = 500
	mazeHeight      = 500 // Высота лабиринта
	buttonHeight    = 30
	borderThickness = float32(2)
)

type Cell struct {
	Visited bool
	Right   bool
	Bottom  bool
}

type Maze struct {
	Rows, Cols int
	Cells      [][]Cell
}

type Game struct {
	maze     *Maze
	cellSize float32
}

func NewMaze(rows, cols int) *Maze {
	ebiten.SetWindowSize(mazeWidth+int(borderThickness*2), mazeHeight+buttonHeight+int(borderThickness))
	cells := make([][]Cell, rows)
	for i := range cells {
		cells[i] = make([]Cell, cols)
		for j := range cells[i] {
			cells[i][j] = Cell{Visited: false, Right: true, Bottom: true}
		}
	}
	return &Maze{Rows: rows, Cols: cols, Cells: cells}
}

func (m *Maze) Generate(x, y int) {
	m.Cells[y][x].Visited = true

	directions := []struct {
		dx, dy int
	}{
		{1, 0},  // вправо
		{0, 1},  // вниз
		{-1, 0}, // влево
		{0, -1}, // вверх
	}

	rand.Shuffle(len(directions), func(i, j int) {
		directions[i], directions[j] = directions[j], directions[i]
	})

	for _, dir := range directions {
		newX, newY := x+dir.dx, y+dir.dy
		if newX >= 0 && newX < m.Cols && newY >= 0 && newY < m.Rows && !m.Cells[newY][newX].Visited {
			if dir.dx == 1 { // вправо
				m.Cells[y][x].Right = false
			} else if dir.dy == 1 { // вниз
				m.Cells[y][x].Bottom = false
			} else if dir.dx == -1 { // влево
				m.Cells[newY][newX].Right = false
			} else if dir.dy == -1 { // вверх
				m.Cells[newY][newX].Bottom = false
			}
			m.Generate(newX, newY)
		}
	}
}

// Game структура для ebiten

// Новый экземпляр игры
// func NewGame(rows, cols int) *Game {
// 	return &Game{
// 		maze:     NewMaze(rows, cols),
// 		cellSize: float32(mazeWidth) / float32(cols),
// 	}
// }

func NewGame(rows, cols int) *Game {
	maze := NewMaze(rows, cols)
	maze.Generate(0, 0)
	cellSize := float32(mazeWidth) / float32(cols)
	return &Game{maze: maze, cellSize: cellSize}
}

// Update обновляет состояние игры
func (g *Game) Update() error {
	// if g.maze.Cells == nil {
	// 	err := g.LoadMazeFromFile("/Users/calamarp/Desktop/go/Go_Maze/src/assets/maze_test.txt")
	// 	if err != nil {
	// 		return err
	// 	}

	// 	fmt.Println("Загруженная матрица:")
	// 	for _, row := range g.maze.Cells {
	// 		fmt.Println(row)
	// 	}
	// }
	return nil
}

// Draw отрисовывает лабиринт и кнопку
func (g *Game) Draw(screen *ebiten.Image) {
	strokeColor := color.RGBA{0, 0, 0, 255}
	fillColor := color.RGBA{255, 255, 255, 255}

	// Рисуем лабиринт
	for y := 0; y < g.maze.Rows; y++ {
		for x := 0; x < g.maze.Cols; x++ {
			// Рисуем ячейку
			if g.maze.Cells[y][x].Visited {
				vector.DrawFilledRect(screen, float32(x)*g.cellSize, float32(y)*g.cellSize, g.cellSize, g.cellSize, fillColor, false)
			}

			// Рисуем правую границу
			if x < g.maze.Cols-1 && g.maze.Cells[y][x].Right {
				vector.StrokeLine(screen, float32(x+1)*g.cellSize, float32(y)*g.cellSize, float32(x+1)*g.cellSize, float32(y+1)*g.cellSize, wallThickness, strokeColor, false)
			}

			// Рисуем нижнюю границу
			if y < g.maze.Rows-1 && g.maze.Cells[y][x].Bottom {
				vector.StrokeLine(screen, float32(x)*g.cellSize, float32(y+1)*g.cellSize, float32(x+1)*g.cellSize, float32(y+1)*g.cellSize, wallThickness, strokeColor, false)
			}
		}
	}

	// Рисуем кнопку под лабиринтом
	buttonY := mazeHeight
	buttonWidth := mazeWidth
	vector.DrawFilledRect(screen, 0, float32(buttonY), float32(buttonWidth), float32(buttonHeight), color.RGBA{200, 200, 200, 255}, false)
	vector.StrokeLine(screen, 0, float32(buttonY), float32(buttonWidth), float32(buttonY), wallThickness, strokeColor, false)
}

// Layout определяет размер окна
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return mazeWidth, mazeHeight + buttonHeight // Общая высота с кнопкой
}

func (g *Game) LoadMazeFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("ошибка открытия файла: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Считываем первую строку с размерами лабиринта.
	if !scanner.Scan() {
		return fmt.Errorf("файл пуст или недоступен")
	}
	dimensions := strings.Fields(scanner.Text())
	fmt.Println("Считанные размеры лабиринта:", dimensions)
	if len(dimensions) != 2 {
		log.Fatal("Неверный формат файла: первая строка должна содержать размеры пещеры.")
	}

	width, err := strconv.Atoi(dimensions[0])
	if err != nil || width <= 0 || width > maxMazeSize {
		log.Fatal("Неверная ширина пещеры:", err)
	}

	height, err := strconv.Atoi(dimensions[1])
	if err != nil || height <= 0 || height > maxMazeSize {
		log.Fatal("Неверная высота пещеры:", err)
	}

	g.maze = NewMaze(height, width)
	fmt.Printf("Создан лабиринт размером %dx%d\n", width, height)

	// Считываем матрицу правых стенок.
	fmt.Println("Матрица со стенками справа:")
	for y := 0; y < height; y++ {
		if !scanner.Scan() {
			return fmt.Errorf("недостаточно строк для матрицы правых стенок")
		}
		row := strings.Fields(scanner.Text())
		fmt.Printf("Считанная строка %d: %v\n", y+1, row)
		if len(row) != width {
			log.Fatal("Неверный формат файла: количество столбцов не совпадает с заданной шириной.")
		}
		for x, cell := range row {
			if cell == "0" {
				g.maze.Cells[y][x].Right = false
			} else if cell == "1" {
				g.maze.Cells[y][x].Right = true
			} else {
				log.Fatal("Неверный символ в пещере: должен быть 0 или 1.")
			}
			fmt.Printf("%s ", cell)
		}
		fmt.Println()
	}

	// Пропускаем пустую строку между матрицами.
	if !scanner.Scan() || scanner.Text() != "" {
		log.Fatal("Ожидалась пустая строка между матрицами.")
	}
	fmt.Println("Пустая строка между матрицами успешно пропущена.")

	// Считываем матрицу нижних стенок.
	fmt.Println("\nМатрица со стенками снизу:")
	for y := 0; y < height; y++ {
		if !scanner.Scan() {
			return fmt.Errorf("недостаточно строк для матрицы нижних стенок")
		}
		row := strings.Fields(scanner.Text())
		fmt.Printf("Считанная строка %d: %v\n", y+1, row)
		if len(row) != width {
			log.Fatal("Неверный формат файла: количество столбцов не совпадает с заданной шириной.")
		}
		for x, cell := range row {
			if cell == "0" {
				g.maze.Cells[y][x].Bottom = false
			} else if cell == "1" {
				g.maze.Cells[y][x].Bottom = true
			} else {
				log.Fatal("Неверный символ в пещере: должен быть 0 или 1.")
			}
			fmt.Printf("%s ", cell)
		}
		fmt.Println()
	}

	return nil
}
func main() {
	w := flag.Int("w", maxMazeSize, "количество строк в лабиринте")
	h := flag.Int("h", maxMazeSize, "количество столбцов в лабиринте")
	flag.Parse()

	game := NewGame(*w, *h)
	// err := game.LoadMazeFromFile("/Users/calamarp/Desktop/go/Go_Maze/src/assets/maze_test.txt")
	// if err != nil {
	// 	log.Fatalf("Failed to load maze from file: %v", err)
	// }
	fmt.Println("Исходная матрица:")
	for _, row := range game.maze.Cells {
		fmt.Println(row)
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
