package main

import (
	"bufio"
	"flag"
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sqweek/dialog"
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
	Set     int
}

type Maze struct {
	Rows, Cols int
	Cells      [][]Cell
}

type Game struct {
	maze       *Maze
	cellSize   float32
	loadButton bool
}

func NewMaze(rows, cols int) *Maze {
	ebiten.SetWindowSize(mazeWidth+int(borderThickness*2), mazeHeight+buttonHeight+int(borderThickness))
	cells := make([][]Cell, rows)
	for i := range cells {
		cells[i] = make([]Cell, cols)
		for j := range cells[i] {
			cells[i][j] = Cell{Visited: false, Right: false, Bottom: false, Set: -1}
		}
	}
	return &Maze{Rows: rows, Cols: cols, Cells: cells}
}

// Инициализация лабиринта
func (m *Maze) Initialize(rows, cols int) {
	m.Rows = rows
	m.Cols = cols
	m.Cells = make([][]Cell, rows)

	for y := 0; y < rows; y++ {
		m.Cells[y] = make([]Cell, cols)
		for x := 0; x < cols; x++ {
			// Устанавливаем все стенки по умолчанию
			m.Cells[y][x].Right = true
			m.Cells[y][x].Bottom = true
		}
	}
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

func debugPrintSets(m *Maze, y int) {
	fmt.Printf("Row %d sets: ", y)
	for x := 0; x < m.Cols; x++ {
		fmt.Printf("%d ", m.Cells[y][x].Set)
	}
	fmt.Println()
}

func (m *Maze) GenerateEller(randomNumbers []int) {
	// Инициализация ячеек
	for row := 0; row < m.Rows; row++ {
		for col := 0; col < m.Cols; col++ {
			m.Cells[row][col].Set = row*m.Cols + col + 1 // Множества начинаются с 1
		}
	}

	index := 0
	currentSetCount := 1 // Начинаем с 1, чтобы множества начинались с 1
	for col := 0; col < m.Cols; col++ {
		m.Cells[0][col].Set = col + 1 // Присваиваем множества начиная с 1
		currentSetCount++
		fmt.Printf(" set %d curr - %d\n", m.Cells[0][col].Set, currentSetCount)
	}

	for row := 0; row < m.Rows; row++ {
		fmt.Println(row)
		if row > 0 {
			for col := 0; col < m.Cols; col++ {
				m.Cells[row][col].Right = m.Cells[row-1][col].Right
				m.Cells[row][col].Bottom = m.Cells[row-1][col].Bottom
				m.Cells[row][col].Set = m.Cells[row-1][col].Set
			}
			// Удаляем правые стенки и нижние границы
			for col := 0; col < m.Cols; col++ {
				m.Cells[row][col].Right = false
				if m.Cells[row-1][col].Bottom {
					m.Cells[row][col].Set = 0        // Присваиваем пустое множество
					m.Cells[row][col].Bottom = false // Удаляем нижнюю стенку
				}
			}
			// Присваиваем новые множества для следующей строки
			for col := 0; col < m.Cols; col++ {
				if m.Cells[row][col].Set == 0 {
					m.Cells[row][col].Set = currentSetCount
					currentSetCount++
					fmt.Printf("1Присвоено новое множество ячейке (%d, %d): Set=%d\n", row, col, m.Cells[row][col].Set)
				}
			}
		}

		fmt.Printf("ПЕРЕД УСТАНОВКОЙ СТЕНКИ:\nСтрока %d", row)
		for col := 0; col < m.Cols; col++ {
			fmt.Printf("{R: %v, B: %v, Set: %d} ", m.Cells[row][col].Right, m.Cells[row][col].Bottom, m.Cells[row][col].Set)
		}
		fmt.Println()

		// Обработка правых стенок
		for col := 0; col < m.Cols-1; col++ {
			fmt.Printf("Перед установкой стенки: Cell(%d, %d) Set=%d\n\n", row, col, m.Cells[row][col].Set)
			fmt.Printf("randomNumbers[index] right = %d\n", randomNumbers[index])
			if randomNumbers[index] == 1 {
				// Ставим стенку
				m.Cells[row][col].Right = true
				fmt.Printf("После установкой стенки: Cell(%d, %d) Set=%d\n", row, col, m.Cells[row][col].Set)
			} else {
				// Не ставим стенку, объединяем множества
				set1 := m.Cells[row][col].Set
				set2 := m.Cells[row][col+1].Set

				if set1 != set2 {
					// Объединяем множества
					for r := 0; r < m.Rows; r++ {
						for c := 0; c < m.Cols; c++ {
							if m.Cells[r][c].Set == set2 {
								m.Cells[r][c].Set = set1
							}
						}
					}
				} else {
					// Ставим стенку, если множества совпадают
					m.Cells[row][col].Right = true
				}
			}
			index++
		}

		// Обработка нижних стенок
		for col := 0; col < m.Cols; col++ {
			set := m.Cells[row][col].Set
			count := 0

			// Подсчет ячеек без нижней границы
			for c := 0; c < m.Cols; c++ {
				if m.Cells[row][c].Set == set && !m.Cells[row][c].Bottom {
					count++
				}
			}

			if count > 1 {
				fmt.Printf("randomNumbers[index] bottom = %d\n", randomNumbers[index])
				if randomNumbers[index] == 1 {
					m.Cells[row][col].Bottom = true
				}
			}
			index++
		}

		// Если это последняя строка, добавляем нижние стенки
		if row == m.Rows-1 {
			for col := 0; col < m.Cols; col++ {
				m.Cells[row][col].Bottom = true
			}
			// Двигаясь слева направо, убираем стенки между ячейками, если множества не совпадают
			for col := 0; col < m.Cols-1; col++ {
				set1 := m.Cells[row][col].Set
				set2 := m.Cells[row][col+1].Set

				if set1 != set2 {
					// Убираем стенку между ячейками
					m.Cells[row][col].Right = false
					// Объединяем множества
					for r := 0; r < m.Rows; r++ {
						for c := 0; c < m.Cols; c++ {
							if m.Cells[r][c].Set == set2 {
								m.Cells[r][c].Set = set1
							}
						}
					}
				}
			}
		}

		// Печатаем измененную строку
		fmt.Print("Измененая строка: [")
		for col := 0; col < m.Cols; col++ {
			fmt.Printf("{Right: %v, Bottom: %v, Set: %d}", m.Cells[row][col].Right, m.Cells[row][col].Bottom, m.Cells[row][col].Set)
			if col < m.Cols-1 {
				fmt.Print(" ")
			}
		}
		fmt.Println("]")
	}

}

// func (m *Maze) GenerateEller() {
// 	r := rand.New(rand.NewSource(time.Now().UnixNano()))

// 	// Инициализация множеств для первой строки
// 	for x := 0; x < m.Cols; x++ {
// 		m.Cells[0][x].Set = x + 1
// 	}

// 	// Обработка всех строк, кроме последней
// 	for y := 0; y < m.Rows-1; y++ {
// 		processRow(m, y, r)
// 	}

// 	// Обработка последней строки
// 	processLastRow(m, r)

// 	// Проверка связности лабиринта
// 	if !isFullyConnected(m) {
// 		fmt.Println("Лабиринт не является идеальным!")
// 	} else {
// 		fmt.Println("Лабиринт идеальный!")
// 	}
// }

func processRow(m *Maze, y int, r *rand.Rand) {
	// Шаг 1: Удаление случайных правых стен
	for x := 0; x < m.Cols-1; x++ {
		if m.Cells[y][x].Set != m.Cells[y][x+1].Set && r.Float32() < 0.5 {
			m.Cells[y][x].Right = false
			mergeSets(m, y, x, x+1)
		}
	}

	// Шаг 2: Удаление случайных нижних стен
	nextSetID := maxSetID(m.Cells[y]) + 1
	for x := 0; x < m.Cols; x++ {
		if r.Float32() < 0.5 || isOnlyInSet(m, y, x) {
			m.Cells[y][x].Bottom = false
		} else {
			// Передача множества в следующую строку
			m.Cells[y+1][x].Set = m.Cells[y][x].Set
		}

		// Если множество в следующей строке еще не установлено, создаем новое
		if m.Cells[y+1][x].Set == 0 {
			m.Cells[y+1][x].Set = nextSetID
			nextSetID++
		}
	}
}

func processLastRow(m *Maze, r *rand.Rand) {
	for x := 0; x < m.Cols-1; x++ {
		if m.Cells[m.Rows-1][x].Set != m.Cells[m.Rows-1][x+1].Set {
			m.Cells[m.Rows-1][x].Right = false
			mergeSets(m, m.Rows-1, x, x+1)
		}
	}
}

func mergeSets(m *Maze, y, x1, x2 int) {
	setToMerge := m.Cells[y][x2].Set
	targetSet := m.Cells[y][x1].Set

	for x := 0; x < m.Cols; x++ {
		if m.Cells[y][x].Set == setToMerge {
			m.Cells[y][x].Set = targetSet
		}
	}
}

func isOnlyInSet(m *Maze, y, x int) bool {
	set := m.Cells[y][x].Set
	count := 0

	for i := 0; i < m.Cols; i++ {
		if m.Cells[y][i].Set == set {
			count++
			if count > 1 {
				return false
			}
		}
	}

	return true
}
func maxSetID(row []Cell) int {
	maxID := 0
	for _, cell := range row {
		if cell.Set > maxID {
			maxID = cell.Set
		}
	}
	return maxID
}

func isFullyConnected(m *Maze) bool {
	visited := make([][]bool, m.Rows)
	for i := range visited {
		visited[i] = make([]bool, m.Cols)
	}

	var dfs func(y, x int)
	dfs = func(y, x int) {
		if y < 0 || y >= m.Rows || x < 0 || x >= m.Cols || visited[y][x] {
			return
		}
		visited[y][x] = true
		if !m.Cells[y][x].Right {
			dfs(y, x+1)
		}
		if !m.Cells[y][x].Bottom {
			dfs(y+1, x)
		}
		if x > 0 && !m.Cells[y][x-1].Right {
			dfs(y, x-1)
		}
		if y > 0 && !m.Cells[y-1][x].Bottom {
			dfs(y-1, x)
		}
	}

	// Запускаем DFS из первой ячейки
	dfs(0, 0)

	// Проверяем, все ли ячейки посещены
	for y := 0; y < m.Rows; y++ {
		for x := 0; x < m.Cols; x++ {
			if !visited[y][x] {
				fmt.Printf("Ячейка не посещена: (%d, %d)\n", y, x)
				return false
			}
		}
	}
	fmt.Println("Все ячейки посещены!")
	return true
}

// func (m *Maze) GenerateEller() {
// 	for x := 0; x < m.Cols; x++ {
// 		m.Cells[0][x].Set = x + 1 // Каждая клетка начинает в своем множестве
// 	}
// 	debugPrintSets(m, 0)

// 	for y := 0; y < m.Rows-1; y++ {
// 		processRow(m, y)
// 		debugPrintSets(m, y+1)
// 	}

// 	processLastRow(m)
// }

// func processRow(m *Maze, y int) {
// 	// Инициализация множеств для следующей строки
// 	for x := 0; x < m.Cols; x++ {
// 		m.Cells[y+1][x].Set = 0 // Сбрасываем множества для следующей строки
// 	}

// 	// Шаг 1: Удаление случайных правых стен
// 	for x := 0; x < m.Cols-1; x++ {
// 		if m.Cells[y][x].Set != m.Cells[y][x+1].Set && rand.Float32() < 0.5 {
// 			m.Cells[y][x].Right = false
// 			mergeSets(m, y, x, x+1)
// 		}
// 	}

// 	// Шаг 2: Удаление случайных нижних стен
// 	for x := 0; x < m.Cols; x++ {
// 		if rand.Float32() < 0.5 || isOnlyInSet(m, y, x) {
// 			m.Cells[y][x].Bottom = false
// 		} else {
// 			// Если стена не удалена, клетка в следующей строке получает новый номер множества
// 			m.Cells[y+1][x].Set = m.Cells[y][x].Set
// 		}
// 	}
// }

// func processLastRow(m *Maze) {
// 	for x := 0; x < m.Cols-1; x++ {
// 		if m.Cells[m.Rows-1][x].Set != m.Cells[m.Rows-1][x+1].Set {
// 			m.Cells[m.Rows-1][x].Right = false
// 			mergeSets(m, m.Rows-1, x, x+1)
// 		}
// 	}
// }

// func mergeSets(m *Maze, y, x1, x2 int) {
// 	setToMerge := m.Cells[y][x2].Set
// 	targetSet := m.Cells[y][x1].Set

// 	for x := 0; x < m.Cols; x++ {
// 		if m.Cells[y][x].Set == setToMerge {
// 			m.Cells[y][x].Set = targetSet
// 		}
// 	}
// }

// func isOnlyInSet(m *Maze, y, x int) bool {
// 	set := m.Cells[y][x].Set
// 	count := 0

// 	for i := 0; i < m.Cols; i++ {
// 		if m.Cells[y][i].Set == set {
// 			count++
// 			if count > 1 {
// 				return false
// 			}
// 		}
// 	}

// 	return true
// }
// func isLastInSet(sets []int, x int) bool {
// 	setID := sets[x]
// 	for i := x + 1; i < len(sets); i++ {
// 		if sets[i] == setID {
// 			return false
// 		}
// 	}
// 	return true
// }

// func maxSetID(sets []int) int {
// 	maxID := sets[0]
// 	for _, id := range sets {
// 		if id > maxID {
// 			maxID = id
// 		}
// 	}
// 	return maxID
// }

// func isFullyConnected(m *Maze) bool {
// 	visited := make([][]bool, m.Rows)
// 	for i := range visited {
// 		visited[i] = make([]bool, m.Cols)
// 	}

// 	var dfs func(y, x int)
// 	dfs = func(y, x int) {
// 		if y < 0 || y >= m.Rows || x < 0 || x >= m.Cols || visited[y][x] {
// 			return
// 		}
// 		visited[y][x] = true
// 		if !m.Cells[y][x].Right {
// 			dfs(y, x+1)
// 		}
// 		if !m.Cells[y][x].Bottom {
// 			dfs(y+1, x)
// 		}
// 		if x > 0 && !m.Cells[y][x-1].Right {
// 			dfs(y, x-1)
// 		}
// 		if y > 0 && !m.Cells[y-1][x].Bottom {
// 			dfs(y-1, x)
// 		}
// 	}

// 	// Запускаем DFS из первой ячейки
// 	dfs(0, 0)

// 	// Проверяем, все ли ячейки посещены
// 	for y := 0; y < m.Rows; y++ {
// 		for x := 0; x < m.Cols; x++ {
// 			if !visited[y][x] {
// 				fmt.Printf("Ячейка не посещена: (%d, %d)\n", y, x)
// 				return false
// 			}
// 		}
// 	}
// 	fmt.Println("Все ячейки посещены!")
// 	return true
// }

// func createPassages(m *Maze) {
// 	visited := make([][]bool, m.Rows)
// 	for i := range visited {
// 		visited[i] = make([]bool, m.Cols)
// 	}

// 	var dfs func(y, x int)
// 	dfs = func(y, x int) {
// 		if y < 0 || y >= m.Rows || x < 0 || x >= m.Cols || visited[y][x] {
// 			return
// 		}
// 		visited[y][x] = true
// 		if !m.Cells[y][x].Right {
// 			dfs(y, x+1)
// 		}
// 		if !m.Cells[y][x].Bottom {
// 			dfs(y+1, x)
// 		}
// 		if x > 0 && !m.Cells[y][x-1].Right {
// 			dfs(y, x-1)
// 		}
// 		if y > 0 && !m.Cells[y-1][x].Bottom {
// 			dfs(y-1, x)
// 		}
// 	}

// 	// Запускаем DFS из первой ячейки
// 	dfs(0, 0)

// 	// Находим все изолированные ячейки и создаем проходы
// 	for y := 0; y < m.Rows; y++ {
// 		for x := 0; x < m.Cols; x++ {
// 			if !visited[y][x] {
// 				// Создаем проход вниз или вправо, но только один
// 				if y < m.Rows-1 {
// 					m.Cells[y][x].Bottom = false
// 					fmt.Printf("Создан проход вниз между (%d, %d) и (%d, %d)\n", y, x, y+1, x)
// 				} else if x < m.Cols-1 {
// 					m.Cells[y][x].Right = false
// 					fmt.Printf("Создан проход вправо между (%d, %d) и (%d, %d)\n", y, x, y, x+1)
// 				}
// 				// Помечаем ячейку как посещенную
// 				visited[y][x] = true
// 			}
// 		}
// 	}
// }

func NewGame(rows, cols int) *Game {
	maze := NewMaze(rows, cols)
	// maze.Initialize(rows, cols)
	// maze.Generate(0, 0)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	numRandomNumbers := rows * cols * 2
	randomNumbers := make([]int, numRandomNumbers)
	for i := range randomNumbers {
		randomNumbers[i] = r.Intn(2) // Генерация 0 или 1
	}
	maze.GenerateEller(randomNumbers)
	cellSize := float32(mazeWidth) / float32(cols)
	return &Game{maze: maze, cellSize: cellSize}
}

// Update обновляет состояние игры
func (g *Game) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()

		if g.isInsideButton(float32(x), float32(y), float32(mazeHeight+borderThickness), buttonHeight) {
			go g.ShowFileSelector()
		}
	}
	return nil
}

func (g *Game) isInsideButton(x, y float32, buttonY float32, buttonHeight float32) bool {
	buttonX := float32(0)
	buttonWidth := float32(mazeHeight + borderThickness*2)
	return x >= buttonX && x <= buttonX+buttonWidth && y >= buttonY && y <= buttonY+buttonHeight
}

func (g *Game) ShowFileSelector() {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Println("Ошибка при получении текущей директории:", err)
		return
	}

	filename, err := dialog.File().
		Filter("Text files", "txt").
		SetStartDir(currentDir).
		Load()

	if err != nil {
		log.Println("Ошибка при выборе файла:", err)
		return
	}

	// Загружаем лабиринт из выбранного файла
	maze, err := LoadMaze(filename)
	if err != nil {
		log.Println("Ошибка при загрузке лабиринта:", err)
		return
	}

	// Обновляем состояние игры с новым лабиринтом
	g.maze = maze
}

// Draw отрисовывает лабиринт и кнопку
func (g *Game) Draw(screen *ebiten.Image) {
	strokeColor := color.RGBA{0, 0, 0, 255}
	fillColor := color.RGBA{255, 255, 255, 255}

	// Рисуем лабиринт
	for y := 0; y < g.maze.Rows; y++ {
		for x := 0; x < g.maze.Cols; x++ {
			// Рисуем ячейку
			// if g.maze.Cells[y][x].Visited {
			vector.DrawFilledRect(screen, float32(x)*g.cellSize, float32(y)*g.cellSize, g.cellSize, g.cellSize, fillColor, false)
			// }

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

func LoadMaze(filename string) (*Maze, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("ошибка при открытии файла: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Читаем размеры лабиринта
	if !scanner.Scan() {
		return nil, fmt.Errorf("ошибка при чтении размеров лабиринта: %v", scanner.Err())
	}
	var rows, cols int
	_, err = fmt.Sscanf(scanner.Text(), "%d %d", &rows, &cols)
	if err != nil {
		return nil, fmt.Errorf("ошибка при парсинге размеров лабиринта: %v", err)
	}
	fmt.Printf("Размеры лабиринта: %d строк, %d столбцов\n", rows, cols)

	maze := &Maze{
		Rows:  rows,
		Cols:  cols,
		Cells: make([][]Cell, rows),
	}

	for i := range maze.Cells {
		maze.Cells[i] = make([]Cell, cols)
	}

	// Читаем первую матрицу (стенки справа)
	for y := 0; y < rows; y++ {
		if !scanner.Scan() {
			return nil, fmt.Errorf("ошибка при чтении стенок справа в строке %d: %v", y, scanner.Err())
		}
		for x := 0; x < cols; x++ {
			var wall int
			_, err = fmt.Sscanf(scanner.Text()[x*2:x*2+1], "%d", &wall) // Предполагаем, что данные разделены пробелами
			if err != nil {
				return nil, fmt.Errorf("ошибка при парсинге стенки справа в строке %d, столбце %d: %v", y, x, err)
			}
			if wall == 1 {
				maze.Cells[y][x].Right = true

			}
			fmt.Printf("Строка %d, столбец %d: стенка справа = %d\n", y, x, wall)
		}
	}

	// Пропускаем пустую строку между матрицами
	if !scanner.Scan() {
		return nil, fmt.Errorf("ошибка при чтении пустой строки между матрицами: %v", scanner.Err())
	}

	// Читаем вторую матрицу (стенки снизу)
	for y := 0; y < rows; y++ {
		if !scanner.Scan() {
			return nil, fmt.Errorf("ошибка при чтении стенок снизу в строке %d: %v", y, scanner.Err())
		}
		for x := 0; x < cols; x++ {
			var wall int
			_, err = fmt.Sscanf(scanner.Text()[x*2:x*2+1], "%d", &wall) // Предполагаем, что данные разделены пробелами
			if err != nil {
				return nil, fmt.Errorf("ошибка при парсинге стенки снизу в строке %d, столбце %d: %v", y, x, err)
			}
			if wall == 1 {
				maze.Cells[y][x].Bottom = true
			}
			maze.Cells[y][x].Visited = true
			fmt.Printf("Строка %d, столбец %d: стенка снизу = %d\n", y, x, wall)
		}
	}

	fmt.Println("Загрузка лабиринта завершена успешно.")
	return maze, nil
}

// SaveMaze сохраняет лабиринт в файл в указанном формате
func (m *Maze) SaveMaze(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Записываем размеры лабиринта
	_, err = fmt.Fprintf(file, "%d %d\n", m.Rows, m.Cols)
	if err != nil {
		return err
	}

	// Записываем стенки справа
	for y := 0; y < m.Rows; y++ {
		for x := 0; x < m.Cols; x++ {
			if x < m.Cols {
				if m.Cells[y][x].Right {
					_, err = fmt.Fprintf(file, "1 ")
				} else {
					_, err = fmt.Fprintf(file, "0 ")
				}
			} else {
				// Для последнего элемента в строке добавляем "1", чтобы закрыть строку
				// _, err = fmt.Fprintf(file, "1")
			}
		}
		_, err = fmt.Fprintln(file) // Переход на новую строку
		if err != nil {
			return err
		}
	}

	// Добавляем пробел между матрицами
	_, err = fmt.Fprintln(file)
	if err != nil {
		return err
	}

	// Записываем стенки снизу
	for y := 0; y < m.Rows; y++ {
		for x := 0; x < m.Cols; x++ {
			if x < m.Cols {
				if m.Cells[y][x].Bottom {
					_, err = fmt.Fprintf(file, "1 ")
				} else {
					_, err = fmt.Fprintf(file, "0 ")
				}
			} else {
				// Для последнего элемента в строке добавляем "0", чтобы закрыть строку
				// _, err = fmt.Fprintf(file, "0")
			}
		}
		_, err = fmt.Fprintln(file) // Переход на новую строку
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	w := flag.Int("w", maxMazeSize, "количество строк в лабиринте")
	h := flag.Int("h", maxMazeSize, "количество столбцов в лабиринте")
	flag.Parse()

	game := NewGame(*w, *h)
	// Печатаем сгенерированный лабиринт в терминал
	fmt.Println("Сгенерированный лабиринт:")
	// game.maze.PrintMaze()
	err := game.maze.SaveMaze("maze.txt")
	if err != nil {
		fmt.Println("Ошибка при сохранении лабиринта:", err)
	} else {
		fmt.Println("Лабиринт успешно сохранен в maze.txt")
	}
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

// // func (m *Maze) PrintMaze() {
// // 	for y := 0; y < m.Rows; y++ {
// // 		// Печатаем верхнюю границу ячейки
// // 		for x := 0; x < m.Cols; x++ {
// // 			if x == 0 {
// // 				fmt.Print("1 ") // Левая граница
// // 			}
// // 			if m.Cells[y][x].Right {
// // 				fmt.Print("1 ") // Стенка справа
// // 			} else {
// // 				fmt.Print("0 ") // Нет стенки справа
// // 			}
// // 		}
// // 		fmt.Println("1") // Правая граница для последней ячейки

// // 		// Печатаем нижнюю границу ячейки
// // 		for x := 0; x < m.Cols; x++ {
// // 			if m.Cells[y][x].Bottom {
// // 				fmt.Print("1 ") // Стенка снизу
// // 			} else {
// // 				fmt.Print("0 ") // Нет стенки снизу
// // 			}
// // 		}
// // 		fmt.Println("1") // Нижняя граница для последней ячейки
// // 	}
// // }

// func (m *Maze) Generate(randomNumbers []int) {
// 	// Инициализация ячеек
// 	for row := 0; row < m.Rows; row++ {
// 		for col := 0; col < m.Cols; col++ {
// 			m.Cells[row][col].Set = row*m.Cols + col + 1 // Множества начинаются с 1
// 		}
// 	}

// 	index := 0
// 	currentSetCount := 1 // Начинаем с 1, чтобы множества начинались с 1
// 	for col := 0; col < 1; col++ {
// 		m.Cells[0][col].Set = col + 1 // Присваиваем множества начиная с 1
// 		currentSetCount++
// 		fmt.Printf(" set %d curr - %d\n", m.Cells[0][col].Set, currentSetCount)
// 	}
// 	for row := 0; row < m.Rows; row++ {
// 		// Создаем новую строку
// 		newRow := make([]Cell, m.Cols)

// 		// Обработка правых стенок
// 		for col := 0; col < m.Cols-1; col++ {
// 			fmt.Printf("randomNumbers[index] right = %d set %d\n", randomNumbers[index], m.Cells[row][col].Set)
// 			if randomNumbers[index] == 1 {
// 				// Ставим стенку
// 				m.Cells[row][col].RightWall = true
// 			} else {
// 				// Не ставим стенку, объединяем множества
// 				set1 := m.Cells[row][col].Set
// 				set2 := m.Cells[row][col+1].Set

// 				if set1 != set2 {
// 					// Объединяем множества
// 					for r := 0; r < m.Rows; r++ {
// 						for c := 0; c < m.Cols; c++ {
// 							if m.Cells[r][c].Set == set2 {
// 								m.Cells[r][c].Set = set1
// 							}
// 						}
// 					}
// 				} else {
// 					// Ставим стенку, если множества совпадают
// 					m.Cells[row][col].RightWall = true
// 				}
// 			}
// 			index++
// 		}

// 		// Обработка нижних стенок
// 		for col := 0; col < m.Cols; col++ {
// 			set := m.Cells[row][col].Set
// 			count := 0

// 			// Подсчет ячеек без нижней границы
// 			for c := 0; c < m.Cols; c++ {
// 				if m.Cells[row][c].Set == set && !m.Cells[row][c].BottomWall {
// 					count++
// 				}
// 			}

// 			if count > 1 {
// 				fmt.Printf("randomNumbers[index] bottom = %d\n", randomNumbers[index])
// 				if randomNumbers[index] == 1 {
// 					m.Cells[row][col].BottomWall = true
// 				}
// 			}
// 			index++
// 		}
// 		// fmt.Printf("Измененая строка: %v\n", newRow)
// 		// Если это последняя строка, добавляем нижние стенки
// 		if row == m.Rows-1 {
// 			for col := 0; col < m.Cols; col++ {
// 				m.Cells[row][col].BottomWall = true
// 			}
// 		} else {
// 			// Копируем текущую строку для следующей итерации
// 			for col := 0; col < m.Cols; col++ {
// 				newRow[col] = m.Cells[row][col]
// 				// Удаляем правые стенки и нижние границы
// 				newRow[col].RightWall = false
// 				if m.Cells[row][col].BottomWall {
// 					newRow[col].Set = 0            // Присваиваем пустое множество
// 					newRow[col].BottomWall = false // Удаляем нижнюю стенку
// 				}
// 			}

// 			// Присваиваем новые множества
// 			for col := 0; col < m.Cols; col++ {
// 				if newRow[col].Set == 0 {
// 					// Присваиваем новое множество
// 					newRow[col].Set = currentSetCount
// 					currentSetCount++
// 					fmt.Printf("Присвоено новое множество ячейке (%d, %d): Set=%d\n", row+1, col, newRow[col].Set)
// 				}
// 			}

// 			// Устанавливаем нижние стенки для новой строки
// 			for col := 0; col < m.Cols; col++ {
// 				if newRow[col].Set != 0 && m.Cells[row][col].BottomWall {
// 					newRow[col].BottomWall = true // Устанавливаем нижнюю стенку, если это необходимо
// 				}
// 			}
// 			fmt.Print("Измененая строка: [")
// 			for col := 0; col < m.Cols; col++ {
// 				fmt.Printf("{%v %v %d}", m.Cells[row][col].RightWall, m.Cells[row][col].BottomWall, m.Cells[row][col].Set)
// 				if col < m.Cols-1 {
// 					fmt.Print(" ")
// 				}
// 			}
// 			fmt.Println("]")
// 			// Добавляем новую строку в лабиринт
// 			m.Cells = append(m.Cells, newRow)
// 			fmt.Printf("Добавлена новая строка: %v\n", newRow)
// 		}
// 	}

// 	// Вывод состояния всех ячеек
// 	for row := 0; row < m.Rows; row++ {
// 		for col := 0; col < m.Cols; col++ {
// 			fmt.Printf("Cell(%d, %d): RightWall=%v, BottomWall=%v, Set=%d\n",
// 				row, col, m.Cells[row][col].RightWall, m.Cells[row][col].BottomWall, m.Cells[row][col].Set)
// 		}
// 	}
// }

// func (m *Maze) Generate(randomNumbers []int) {
// 	// Инициализация ячеек
// 	for row := 0; row < m.Rows; row++ {
// 		for col := 0; col < m.Cols; col++ {
// 			m.Cells[row][col].Set = row*m.Cols + col + 1 // Множества начинаются с 1
// 		}
// 	}

// 	index := 0
// 	currentSetCount := 1 // Начинаем с 1, чтобы множества начинались с 1
// 	for col := 0; col < m.Cols; col++ {
// 		m.Cells[0][col].Set = col + 1 // Присваиваем множества начиная с 1
// 		currentSetCount++
// 		fmt.Printf(" set %d curr - %d\n", m.Cells[0][col].Set, currentSetCount)
// 	}
// 	// Создаем новую строку вне цикла
// 	newRow := make([]Cell, m.Cols)
// 	for row := 0; row < m.Rows; row++ {
// 		// fmt.Printf("некст новая строка: %v\n", newRow)
// 		// Обработка правых стенок
// 		if row > 0 {
// 			m.Cells = append(m.Cells, newRow)
// 		}
// 		for col := 0; col < m.Cols-1; col++ {
// 			fmt.Printf("Перед установкой стенки: Cell(%d, %d) Set=%d\n\n", row, col, m.Cells[row][col].Set)
// 			fmt.Printf("randomNumbers[index] right = %d\n", randomNumbers[index])
// 			if randomNumbers[index] == 1 {
// 				// Ставим стенку
// 				m.Cells[row][col].RightWall = true
// 				fmt.Printf("Псоле установкой стенки: Cell(%d, %d) Set=%d\n", row, col, m.Cells[row][col].Set)
// 			} else {
// 				// Не ставим стенку, объединяем множества
// 				set1 := m.Cells[row][col].Set
// 				set2 := m.Cells[row][col+1].Set

// 				if set1 != set2 {
// 					// Объединяем множества
// 					for r := 0; r < m.Rows; r++ {
// 						for c := 0; c < m.Cols; c++ {
// 							if m.Cells[r][c].Set == set2 {
// 								m.Cells[r][c].Set = set1
// 							}
// 						}
// 					}
// 				} else {
// 					// Ставим стенку, если множества совпадают
// 					m.Cells[row][col].RightWall = true

// 				}
// 			}
// 			index++
// 		}

// 		// Обработка нижних стенок
// 		for col := 0; col < m.Cols; col++ {
// 			set := m.Cells[row][col].Set
// 			count := 0

// 			// Подсчет ячеек без нижней границы
// 			for c := 0; c < m.Cols; c++ {
// 				if m.Cells[row][c].Set == set && !m.Cells[row][c].BottomWall {
// 					count++
// 				}
// 			}

// 			if count > 1 {
// 				fmt.Printf("randomNumbers[index] bottom = %d\n", randomNumbers[index])
// 				if randomNumbers[index] == 1 {
// 					m.Cells[row][col].BottomWall = true
// 				}
// 			}
// 			index++
// 		}
// 		// fmt.Printf("Измененая строка: %v\n", newRow)
// 		// Если это последняя строка, добавляем нижние стенки
// 		if row == m.Rows-1 {
// 			for col := 0; col < m.Cols; col++ {
// 				m.Cells[row][col].BottomWall = true
// 			}
// 		} else {
// 			// Копируем текущую строку для следующей итерации
// 			for col := 0; col < m.Cols; col++ {
// 				newRow[col] = Cell{
// 					RightWall:  m.Cells[row][col].RightWall,
// 					BottomWall: m.Cells[row][col].BottomWall,
// 					Set:        m.Cells[row][col].Set,
// 				}
// 				// Удаляем правые стенки и нижние границы
// 				newRow[col].RightWall = false
// 				if m.Cells[row][col].BottomWall {
// 					newRow[col].Set = 0            // Присваиваем пустое множество
// 					newRow[col].BottomWall = false // Удаляем нижнюю стенку
// 				}
// 			}
// 			// Обновляем множества для следующей строки
// 			for col := 0; col < m.Cols; col++ {
// 				// Удаляем правые стенки и нижние границы
// 				m.Cells[row+1][col].RightWall = false
// 				if m.Cells[row][col].BottomWall {
// 					m.Cells[row+1][col].Set = 0            // Присваиваем пустое множество
// 					m.Cells[row+1][col].BottomWall = false // Удаляем нижнюю стенку
// 				}
// 			}
// 			// Присваиваем новые множества для следующей строки
// 			for col := 0; col < m.Cols; col++ {
// 				if m.Cells[row+1][col].Set == 0 {
// 					// Присваиваем новое множество
// 					m.Cells[row+1][col].Set = currentSetCount
// 					currentSetCount++
// 					fmt.Printf("Присвоено новое множество ячейке (%d, %d): Set=%d\n", row+1, col, m.Cells[row+1][col].Set)
// 				}
// 			}
// 			// for col := 0; col < m.Cols; col++ {
// 			// 	if m.Cells[row+1][col].Set != 0 && m.Cells[row][col].BottomWall {
// 			// 		m.Cells[row+1][col].BottomWall = true // Устанавливаем нижнюю стенку, если это необходимо
// 			// 	}
// 			// }
// 			// // Присваиваем новые множества
// 			// for col := 0; col < m.Cols; col++ {
// 			// 	if newRow[col].Set == 0 {
// 			// 		// Присваиваем новое множество
// 			// 		newRow[col].Set = currentSetCount
// 			// 		currentSetCount++
// 			// 		fmt.Printf("Присвоено новое множество ячейке (%d, %d): Set=%d\n", row+1, col, newRow[col].Set)
// 			// 	}
// 			// }

// 			// Устанавливаем нижние стенки для новой строки
// 			for col := 0; col < m.Cols; col++ {
// 				if newRow[col].Set != 0 && m.Cells[row][col].BottomWall {
// 					newRow[col].BottomWall = true // Устанавливаем нижнюю стенку, если это необходимо
// 				}
// 			}
// 			fmt.Print("Измененая строка: [")
// 			for col := 0; col < m.Cols; col++ {
// 				fmt.Printf("{%v %v %d}", m.Cells[row][col].RightWall, m.Cells[row][col].BottomWall, m.Cells[row][col].Set)
// 				if col < m.Cols-1 {
// 					fmt.Print(" ")
// 				}
// 			}
// 			fmt.Println("]")
// 			Добавляем новую строку в лабиринт
// 			m.Cells = append(m.Cells, newRow)
// 			fmt.Printf("Добавлена новая строка: %v\n\n %v\n", newRow, m.Cells)
// 		}
// 	}

// 	// Вывод состояния всех ячеек
// 	for row := 0; row < m.Rows; row++ {
// 		for col := 0; col < m.Cols; col++ {
// 			fmt.Printf("Cell(%d, %d): RightWall=%v, BottomWall=%v, Set=%d\n",
// 				row, col, m.Cells[row][col].RightWall, m.Cells[row][col].BottomWall, m.Cells[row][col].Set)
// 		}
// 	}
// }
