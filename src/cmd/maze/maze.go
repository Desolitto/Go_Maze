package main

import (
	"flag"
	"fmt"
	"image/color"
	"log"
	"math/rand"

	// "github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Maze struct {
	Rows              int
	Cols              int
	ActiveRow         int
	SetCounter        int
	SetMatrix         [][]int
	RightBorderMatrix [][]int
	LowBorderMatrix   [][]int
}

func NewMaze(rows int, cols int) *Maze {
	maze := &Maze{
		Rows:              rows,
		Cols:              cols,
		ActiveRow:         0,
		SetCounter:        1,
		SetMatrix:         make([][]int, rows),
		RightBorderMatrix: make([][]int, rows),
		LowBorderMatrix:   make([][]int, rows),
	}
	for i := range maze.SetMatrix {
		maze.SetMatrix[i] = make([]int, cols)
		maze.RightBorderMatrix[i] = make([]int, cols)
		maze.LowBorderMatrix[i] = make([]int, cols)
	}
	return maze
}

func (m *Maze) RowSize() int {
	return len(m.RightBorderMatrix)
}

func (m *Maze) ColSize() int {
	return len(m.RightBorderMatrix[0])
}

// Построчная генерация лабиринта с помощью алгоритма Эйлера
func GenerateMaze(rows int, cols int) (*Maze, error) {
	if rows < 1 || cols < 1 {
		return nil, fmt.Errorf("rows and columns must be positive numbers")
	}
	if rows > 50 || cols > 50 {
		return nil, fmt.Errorf("rows and columns must be <= 50")
	}
	maze := NewMaze(rows, cols)
	for i := 0; i < rows; i++ {
		maze.assignUniqueSet()
		maze.addingVerticalWalls()
		maze.addingHorizontalWalls()
		maze.preparatingNewLine()
	}
	maze.addingEndLine()
	// maze.writeToFile()
	return maze, nil
}

// Присвоение ячейки множества
func (m *Maze) assignUniqueSet() {
	for j := 0; j < m.Cols; j++ {
		if m.SetMatrix[m.ActiveRow][j] == 0 {
			m.SetMatrix[m.ActiveRow][j] = m.SetCounter
			m.SetCounter++
		}
	}
}

func (m *Maze) addingVerticalWalls() {
	for i := 0; i < m.Cols-1; i++ {
		choise := rand.Int() % 2
		if choise == 1 || m.SetMatrix[m.ActiveRow][i] == m.SetMatrix[m.ActiveRow][i+1] {
			m.RightBorderMatrix[m.ActiveRow][i] = 1
		} else {
			m.mergeSet(i)
		}
	}
	m.RightBorderMatrix[m.ActiveRow][m.Cols-1] = 1
}

// Объединение ячеек в одно множество
func (m *Maze) mergeSet(i int) {
	x := m.SetMatrix[m.ActiveRow][i+1]
	for j := 0; j < m.Cols; j++ {
		if m.SetMatrix[m.ActiveRow][j] == x {
			m.SetMatrix[m.ActiveRow][j] = m.SetMatrix[m.ActiveRow][i]
		}
	}
}

// Добавление горизонтальных (нижних) стен
func (m *Maze) addingHorizontalWalls() {
	for i := 0; i < m.Cols; i++ {
		choise := rand.Int() % 2
		check := m.checkedHorizontalWalls(i)
		if choise == 1 && check {
			m.LowBorderMatrix[m.ActiveRow][i] = 1
		}
	}
}

func (m *Maze) checkedHorizontalWalls(index int) bool {
	set := m.SetMatrix[m.ActiveRow][index]
	for i := 0; i < m.Cols; i++ {
		if m.SetMatrix[m.ActiveRow][i] == set && i != index && m.LowBorderMatrix[m.ActiveRow][i] == 0 {
			return true
		}
	}
	return false
}

func (m *Maze) preparatingNewLine() {
	if m.ActiveRow == m.Rows-1 {
		return
	}
	m.ActiveRow++
	for i := 0; i < m.Cols; i++ {
		if m.LowBorderMatrix[m.ActiveRow-1][i] == 0 {
			m.SetMatrix[m.ActiveRow][i] = m.SetMatrix[m.ActiveRow-1][i]
		} else {
			m.SetMatrix[m.ActiveRow][i] = 0
		}
	}
}

func (m *Maze) addingEndLine() {
	for i := 0; i < m.Cols-1; i++ {
		m.LowBorderMatrix[m.ActiveRow][i] = 1
		if m.SetMatrix[m.ActiveRow][i] != m.SetMatrix[m.ActiveRow][i+1] {
			m.RightBorderMatrix[m.ActiveRow][i] = 0
			m.mergeSet(i)
		}
	}
	m.LowBorderMatrix[m.ActiveRow][m.Cols-1] = 1
	m.RightBorderMatrix[m.ActiveRow][m.Cols-1] = 1
}

// Game структура для ebiten
type Game struct {
	maze *Maze
}

// Новый экземпляр игры
func NewGame(rows, cols int) *Game {
	maze, err := GenerateMaze(rows, cols)
	if err != nil {
		log.Fatal(err)
	}
	return &Game{maze: maze}
}

// Update обновляет состояние игры
func (g *Game) Update() error {
	return nil
}

// Draw отрисовывает лабиринт
func (g *Game) Draw(screen *ebiten.Image) {
	cellWidth := float32(10.0)  // Изменено на float32
	cellHeight := float32(10.0) // Изменено на float32

	// Устанавливаем цвет стен
	strokeColor := color.RGBA{0, 0, 0, 255}     // Черный цвет
	fillColor := color.RGBA{255, 255, 255, 255} // Белый цвет

	// Рисуем ячейки и стены
	for y := 0; y < g.maze.Rows; y++ {
		for x := 0; x < g.maze.Cols; x++ {
			if g.maze.SetMatrix[y][x] == 1 {
				vector.DrawFilledRect(screen, float32(x)*cellWidth, float32(y)*cellHeight, cellWidth, cellHeight, fillColor, false)
			}
			if g.maze.RightBorderMatrix[y][x] == 1 {
				vector.StrokeLine(screen, float32(x+1)*cellWidth, float32(y)*cellHeight, float32(x+1)*cellWidth, float32(y+1)*cellHeight, 1, strokeColor, false)
			}
			if g.maze.LowBorderMatrix[y][x] == 1 {
				vector.StrokeLine(screen, float32(x)*cellWidth, float32(y+1)*cellHeight, float32(x+1)*cellWidth, float32(y+1)*cellHeight, 1, strokeColor, false)
			}
		}
	}
}

// Layout определяет размер окна
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.maze.Cols * 10, g.maze.Rows * 10 // Умножаем на размер ячейки
}

func main() {
	// rand.Seed(time.Now().UnixNano())
	w := flag.Int("w", 20, "width of the maze")
	h := flag.Int("h", 20, "height of the maze")
	flag.Parse()

	game := NewGame(*h, *w)
	fmt.Println("Исходная матрица:")
	for _, row := range game.maze.SetMatrix {
		fmt.Println(row)
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

// func (m *Maze) LoadCaveMazeFile(filename string) {
// 	file, err := os.Open(filename)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer file.Close()

// 	scanner := bufio.NewScanner(file)

// 	if scanner.Scan() {
// 		dimensions := strings.Fields(scanner.Text())
// 		if len(dimensions) != 2 {
// 			log.Fatal("Неверный формат файла: первая строка должна содержать размеры пещеры.")
// 		}

// 		width, err := strconv.Atoi(dimensions[0])
// 		if err != nil {
// 			log.Fatal("Неверная ширина лабиринта.")
// 		}

// 		height, err := strconv.Atoi(dimensions[1])
// 		if err != nil {
// 			log.Fatal("Неверная высота лабиринта.")
// 		}

// 		m.Rows, m.Cols = height, width
// 		m.Cave = cave.NewCave(width, height)
// 		// Загрузка первой матрицы (стена справа)
// 		m.RightBorderMatrix = make([][]int, height)

// 		for y := 0; y < height; y++ {
// 			if scanner.Scan() {
// 				row := strings.Fields(scanner.Text())
// 				if len(row) != width {
// 					log.Fatal("Неверный формат файла: количество столбцов не совпадает с заданной шириной.")
// 				}
// 				m.RightBorderMatrix[y] = make([]int, width)
// 				for x, cell := range row {
// 					if cell == "0" {
// 						m.Cave.Grid[y][x] = cave.Death
// 					} else if cell == "1" {
// 						m.Cave.Grid[y][x] = cave.Alive
// 					} else {
// 						log.Fatal("Неверный символ в пещере: должен быть 0 или 1.")
// 					}
// 				}
// 			}
// 		}
// 		m.LowBorderMatrix = make([][]int, height)
// 		for y := 0; y < height; y++ {
// 			if scanner.Scan() {
// 				row := strings.Fields(scanner.Text())
// 				if len(row) != width {
// 					log.Fatal("Неверный формат файла: количество столбцов не совпадает с заданной шириной.")
// 				}
// 				m.LowBorderMatrix[y] = make([]int, width)
// 				for x, cell := range row {
// 					value, err := strconv.Atoi(cell)
// 					if err != nil {
// 						log.Fatal("Неверный символ в матрице стен: должен быть целым числом.")
// 					}
// 					m.LowBorderMatrix[y][x] = value
// 				}
// 			}
// 		}
// 	}
// }

// func (m *Maze) PrintMaze() {
// 	for y := 0; y < m.Rows; y++ {
// 		for x := 0; x < m.Cols; x++ {
// 			if m.Cave.Grid[y][x] == cave.Alive {
// 				fmt.Print("1 ")
// 			} else {
// 				fmt.Print("0 ")
// 			}
// 		}
// 		fmt.Println()
// 	}
// }

// func main() {
// 	m := &Maze{}
// 	m.LoadCaveMazeFile("/Users/calamarp/Desktop/go/Go_Maze/src/assets/maze_test.txt")
// 	m.PrintMaze()
// }

// func main() {
// 	rand.Seed(uint64(time.Now().UnixNano()))
// 	w := flag.Int("w", 20, "width of the cave")
// 	h := flag.Int("h", 20, "height of the cave")
// 	initialChance := flag.Int("с", 55, "initial chance (0-100)")
// 	flag.Parse()

// 	game := game.NewGame(*w, *h, *initialChance)
// 	fmt.Println("Исходная матрица:")
// 	// game.PrintCave()
// 	if err := ebiten.RunGame(game); err != nil {
// 		log.Fatal(err)
// 	}
// }
