package main

import (
	"bufio"
	"flag"
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"os"

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
			cells[i][j] = Cell{Visited: false, Right: true, Bottom: true}
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

func (m *Maze) GenerateEller() {
	// Инициализация наборов для первой строки
	sets := make([]int, m.Cols)
	for j := range sets {
		sets[j] = j // Каждая ячейка начинает в своем собственном наборе
	}

	for y := 0; y < m.Rows; y++ {
		// Генерация горизонтальных соединений
		for x := 0; x < m.Cols-1; x++ {
			if rand.Intn(2) == 0 { // С вероятностью 50% удаляем стену справа
				if sets[x] != sets[x+1] {
					m.Cells[y][x].Right = false
					mergeSets(sets, sets[x], sets[x+1])
				}
			}
		}

		// Генерация вертикальных соединений
		nextRowSets := make([]int, m.Cols)
		for x := 0; x < m.Cols; x++ {
			m.Cells[y][x].Visited = true
			if y < m.Rows-1 { // Если это не последняя строка
				if rand.Intn(2) == 0 || isLastInSet(sets, x) {
					m.Cells[y][x].Bottom = false
					nextRowSets[x] = sets[x]
				} else {
					nextRowSets[x] = -1 // Новый набор для следующей строки
				}
			}
		}

		// Если это не последняя строка, обновляем наборы
		if y < m.Rows-1 {
			for x := 0; x < m.Cols; x++ {
				if nextRowSets[x] == -1 {
					nextRowSets[x] = maxSetID(sets) + 1
				}
			}
			sets = nextRowSets
		}
	}

	// Объединение всех наборов в последней строке
	for x := 0; x < m.Cols-1; x++ {
		if sets[x] != sets[x+1] {
			m.Cells[m.Rows-1][x].Right = false
			mergeSets(sets, sets[x], sets[x+1])
		}
	}
	if !isFullyConnected(m) {
		createPassages(m)
	}

}

// Вспомогательные функции

func createPassages(m *Maze) {
	// Логика для создания проходов между ячейками
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

	// Находим все ячейки и создаем проходы
	for y := 0; y < m.Rows; y++ {
		for x := 0; x < m.Cols; x++ {
			if visited[y][x] {
				// Проверяем соседние ячейки для создания проходов
				if y < m.Rows-1 && !visited[y+1][x] {
					m.Cells[y][x].Bottom = false // Создаем проход вниз
					fmt.Printf("Создан проход вниз между (%d, %d) и (%d, %d)\n", y, x, y+1, x)
				}
				if x < m.Cols-1 && !visited[y][x+1] {
					m.Cells[y][x].Right = false // Создаем проход вправо
					fmt.Printf("Создан проход вправо между (%d, %d) и (%d, %d)\n", y, x, y, x+1)
				}
			}
		}
	}
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

// mergeSets объединяет два набора
func mergeSets(sets []int, fromSet, toSet int) {
	for i := range sets {
		if sets[i] == fromSet {
			sets[i] = toSet
		}
	}
}

// isLastInSet проверяет, является ли ячейка последней в своем наборе
func isLastInSet(sets []int, x int) bool {
	for i := x + 1; i < len(sets); i++ {
		if sets[i] == sets[x] {
			return false
		}
	}
	return true
}

// maxSetID возвращает максимальный идентификатор набора
func maxSetID(sets []int) int {
	maxID := sets[0]
	for _, id := range sets {
		if id > maxID {
			maxID = id
		}
	}
	return maxID
}

// func (m *Maze) GenerateEller() {
// 	sets := make([]int, m.Cols)
// 	for j := range sets {
// 		sets[j] = j // Каждая ячейка начинает в своем собственном наборе
// 	}
// 	fmt.Println(m.Rows, m.Cols)
// 	for y := 0; y < m.Rows; y++ {
// 		for x := 0; x < m.Cols; x++ {
// 			m.Cells[y][x].Visited = true

// 			// Удаляем стенку справа с вероятностью 50%
// 			if x < m.Cols-1 && rand.Intn(2) == 0 {
// 				m.Cells[y][x].Right = false
// 				// 	// Объединяем наборы
// 				sets[x+1] = sets[x]
// 			}

// 			// // Удаляем стенку снизу с вероятностью 50%
// 			if y < m.Rows-1 && (x == 0 || rand.Intn(2) == 0) {
// 				m.Cells[y][x].Bottom = false
// 				// Объединяем наборы
// 				sets[x] = sets[x]
// 			}
// 		}

// 		// Объединяем наборы после обработки строки
// 		for x := 1; x < m.Cols; x++ {
// 			if sets[x] != sets[x-1] {
// 				oldSet := sets[x]
// 				for j := x; j < m.Cols; j++ {
// 					if sets[j] == oldSet {
// 						sets[j] = sets[x-1]
// 					}
// 				}
// 			}
// 		}
// 	}
// }

// 	for y := 0; y < m.Rows; y++ {
// 		for x := 0; x < m.Cols; x++ {
// 			if m.Cells[y][x].Right {
// 				fmt.Print("1 ")
// 			} else {
// 				fmt.Print("0 ")
// 			}
// 		}
// 		fmt.Println()
// 	}
// 	fmt.Println("")
// 	for y := 0; y < m.Rows; y++ {
// 		for x := 0; x < m.Cols; x++ {
// 			if m.Cells[y][x].Bottom {
// 				fmt.Print("1 ")
// 			} else {
// 				fmt.Print("0 ")
// 			}
// 		}
// 		fmt.Println()
// 	}
// }

func NewGame(rows, cols int) *Game {
	maze := NewMaze(rows, cols)
	maze.Initialize(rows, cols)
	// maze.Generate(0, 0)
	maze.GenerateEller()
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

// func (m *Maze) PrintMaze() {
// 	for y := 0; y < m.Rows; y++ {
// 		// Печатаем верхнюю границу ячейки
// 		for x := 0; x < m.Cols; x++ {
// 			if x == 0 {
// 				fmt.Print("1 ") // Левая граница
// 			}
// 			if m.Cells[y][x].Right {
// 				fmt.Print("1 ") // Стенка справа
// 			} else {
// 				fmt.Print("0 ") // Нет стенки справа
// 			}
// 		}
// 		fmt.Println("1") // Правая граница для последней ячейки

// 		// Печатаем нижнюю границу ячейки
// 		for x := 0; x < m.Cols; x++ {
// 			if m.Cells[y][x].Bottom {
// 				fmt.Print("1 ") // Стенка снизу
// 			} else {
// 				fmt.Print("0 ") // Нет стенки снизу
// 			}
// 		}
// 		fmt.Println("1") // Нижняя граница для последней ячейки
// 	}
// }
