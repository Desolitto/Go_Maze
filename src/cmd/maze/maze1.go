package main

import (
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

func NewGame(rows, cols int) *Game {
	maze := NewMaze(rows, cols)
	maze.Initialize(rows, cols)
	maze.Generate(0, 0)
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

// func (g *Game) LoadMazeFromFile(filename string) error {
// 	file, err := os.Open(filename)
// 	if err != nil {
// 		return fmt.Errorf("ошибка открытия файла: %v", err)
// 	}
// 	defer file.Close()

// 	scanner := bufio.NewScanner(file)

// 	// Считываем первую строку с размерами лабиринта.
// 	if !scanner.Scan() {
// 		return fmt.Errorf("файл пуст или недоступен")
// 	}
// 	dimensions := strings.Fields(scanner.Text())
// 	fmt.Println("Считанные размеры лабиринта:", dimensions)
// 	if len(dimensions) != 2 {
// 		log.Fatal("Неверный формат файла: первая строка должна содержать размеры пещеры.")
// 	}

// 	width, err := strconv.Atoi(dimensions[0])
// 	if err != nil || width <= 0 || width > maxMazeSize {
// 		log.Fatal("Неверная ширина пещеры:", err)
// 	}

// 	height, err := strconv.Atoi(dimensions[1])
// 	if err != nil || height <= 0 || height > maxMazeSize {
// 		log.Fatal("Неверная высота пещеры:", err)
// 	}

// 	g.maze = NewMaze(height, width)
// 	fmt.Printf("Создан лабиринт размером %dx%d\n", width, height)

// 	// Считываем матрицу правых стенок.
// 	fmt.Println("Матрица со стенками справа:")
// 	for y := 0; y < height; y++ {
// 		if !scanner.Scan() {
// 			return fmt.Errorf("недостаточно строк для матрицы правых стенок")
// 		}
// 		row := strings.Fields(scanner.Text())
// 		fmt.Printf("Считанная строка %d: %v\n", y+1, row)
// 		if len(row) != width {
// 			log.Fatal("Неверный формат файла: количество столбцов не совпадает с заданной шириной.")
// 		}
// 		for x, cell := range row {
// 			if cell == "0" {
// 				g.maze.Cells[y][x].Right = false
// 			} else if cell == "1" {
// 				g.maze.Cells[y][x].Right = true
// 			} else {
// 				log.Fatal("Неверный символ в пещере: должен быть 0 или 1.")
// 			}
// 			fmt.Printf("%s ", cell)
// 		}
// 		fmt.Println()
// 	}

// 	// Пропускаем пустую строку между матрицами.
// 	if !scanner.Scan() || scanner.Text() != "" {
// 		log.Fatal("Ожидалась пустая строка между матрицами.")
// 	}
// 	fmt.Println("Пустая строка между матрицами успешно пропущена.")

// 	// Считываем матрицу нижних стенок.
// 	fmt.Println("\nМатрица со стенками снизу:")
// 	for y := 0; y < height; y++ {
// 		if !scanner.Scan() {
// 			return fmt.Errorf("недостаточно строк для матрицы нижних стенок")
// 		}
// 		row := strings.Fields(scanner.Text())
// 		fmt.Printf("Считанная строка %d: %v\n", y+1, row)
// 		if len(row) != width {
// 			log.Fatal("Неверный формат файла: количество столбцов не совпадает с заданной шириной.")
// 		}
// 		for x, cell := range row {
// 			if cell == "0" {
// 				g.maze.Cells[y][x].Bottom = false
// 			} else if cell == "1" {
// 				g.maze.Cells[y][x].Bottom = true
// 			} else {
// 				log.Fatal("Неверный символ в пещере: должен быть 0 или 1.")
// 			}
// 			fmt.Printf("%s ", cell)
// 		}
// 		fmt.Println()
// 	}

// 	return nil
// }

// LoadMaze загружает лабиринт из файла
// func LoadMaze(filename string) (*Maze, error) {
// 	file, err := os.Open(filename)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer file.Close()

// 	var rows, cols int
// 	_, err = fmt.Fscanf(file, "%d %d\n", &rows, &cols)
// 	if err != nil {
// 		return nil, err
// 	}

// 	maze := &Maze{
// 		Rows:  rows,
// 		Cols:  cols,
// 		Cells: make([][]Cell, rows),
// 	}

// 	for i := range maze.Cells {
// 		maze.Cells[i] = make([]Cell, cols)
// 	}

// 	// Читаем первую матрицу (стенки справа)
// 	for y := 0; y < rows; y++ {
// 		for x := 0; x < cols; x++ {
// 			var wall int
// 			_, err = fmt.Fscanf(file, "%d", &wall)
// 			if err != nil {
// 				return nil, err
// 			}
// 			if wall == 1 {
// 				maze.Cells[y][x].Right = true
// 			}
// 		}
// 	}

// 	// Читаем вторую матрицу (стенки снизу)
// 	for y := 0; y < rows; y++ {
// 		for x := 0; x < cols; x++ {
// 			var wall int
// 			_, err = fmt.Fscanf(file, "%d", &wall)
// 			if err != nil {
// 				return nil, err
// 			}
// 			if wall == 1 {
// 				maze.Cells[y][x].Bottom = true
// 			}
// 		}
// 	}

// 	return maze, nil
// }

func LoadMaze(filename string) (*Maze, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var rows, cols int
	_, err = fmt.Fscanf(file, "%d %d\n", &rows, &cols)
	if err != nil {
		return nil, err
	}

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
		for x := 0; x < cols; x++ {
			var wall int
			_, err = fmt.Fscanf(file, "%d", &wall)
			if err != nil {
				return nil, err
			}
			if wall == 1 {
				maze.Cells[y][x].Right = true
			}
		}
	}

	// Пропускаем пустую строку между матрицами
	var emptyLine string
	_, err = fmt.Fscanln(file, &emptyLine)
	if err != nil {
		return nil, err
	}

	// Читаем вторую матрицу (стенки снизу)
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			var wall int
			_, err = fmt.Fscanf(file, "%d", &wall)
			if err != nil {
				return nil, err
			}
			if wall == 1 {
				maze.Cells[y][x].Bottom = true
			}
		}
	}

	// Установим стенки для последнего ряда и последнего столбца
	for y := 0; y < rows; y++ {
		maze.Cells[y][cols-1].Right = true // Стенка справа для последнего столбца
	}
	for x := 0; x < cols; x++ {
		maze.Cells[rows-1][x].Bottom = true // Стенка снизу для последнего ряда
	}

	return maze, nil
}

// func (m *Maze) GeneratePerfectMaze() {
// 	type Edge struct {
// 		x1, y1, x2, y2 int // Координаты соединяемых клеток
// 	}

// 	edges := []Edge{}
// 	for y := 0; y < m.Rows; y++ {
// 		for x := 0; x < m.Cols; x++ {
// 			if x < m.Cols-1 {
// 				edges = append(edges, Edge{x, y, x + 1, y}) // Правое соединение
// 			}
// 			if y < m.Rows-1 {
// 				edges = append(edges, Edge{x, y, x, y + 1}) // Нижнее соединение
// 			}
// 		}
// 	}

// 	// Перемешиваем ребра
// 	rand.Shuffle(len(edges), func(i, j int) {
// 		edges[i], edges[j] = edges[j], edges[i]
// 	})

// 	// Создаем структуру для отслеживания соединенных компонентов
// 	parent := make([]int, m.Rows*m.Cols)
// 	for i := range parent {
// 		parent[i] = i
// 	}

// 	find := func(x int) int {
// 		if parent[x] != x {
// 			parent[x] = find(parent[x])
// 		}
// 		return parent[x]
// 	}

// 	union := func(x, y int) {
// 		rootX := find(x)
// 		rootY := find(y)
// 		if rootX != rootY {
// 			parent[rootX] = rootY
// 		}
// 	}

// 	// Генерация лабиринта
// 	for _, edge := range edges {
// 		cell1 := edge.y1*m.Cols + edge.x1
// 		cell2 := edge.y2*m.Cols + edge.x2

// 		if find(cell1) != find(cell2) {
// 			// Если клетки не соединены, соединяем их и убираем стену
// 			union(cell1, cell2)
// 			if edge.x1 == edge.x2 { // Если соединение вертикальное
// 				m.Cells[edge.y1][edge.x1].Bottom = false
// 			} else { // Если соединение горизонтальное
// 				m.Cells[edge.y1][edge.x1].Right = false
// 			}
// 		}
// 	}
// }

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
			if x < m.Cols-1 {
				if m.Cells[y][x].Right {
					_, err = fmt.Fprintf(file, "1 ")
				} else {
					_, err = fmt.Fprintf(file, "0 ")
				}
			}
		}
		_, err = fmt.Fprintln(file, "1") // Стенка справа для последней ячейки
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
			if x < m.Cols-1 {
				if m.Cells[y][x].Bottom {
					_, err = fmt.Fprintf(file, "1 ")
				} else {
					_, err = fmt.Fprintf(file, "0 ")
				}
			}
		}
		_, err = fmt.Fprintln(file, "0") // Нижняя граница для последней ячейки
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
	// err := game.LoadMazeFromFile("/Users/calamarp/Desktop/go/Go_Maze/src/assets/maze_test.txt")
	// if err != nil {
	// 	log.Fatalf("Failed to load maze from file: %v", err)
	// }
	// fmt.Println("Исходная матрица:")
	// for _, row := range game.maze.Cells {
	// 	fmt.Println(row)
	// }
	// Печатаем сгенерированный лабиринт в терминал
	fmt.Println("Сгенерированный лабиринт:")
	game.maze.PrintMaze()
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

func (m *Maze) PrintMaze() {
	for y := 0; y < m.Rows; y++ {
		// Печатаем верхнюю границу ячейки
		for x := 0; x < m.Cols; x++ {
			if x == 0 {
				fmt.Print("1 ") // Левая граница
			}
			if m.Cells[y][x].Right {
				fmt.Print("1 ") // Стенка справа
			} else {
				fmt.Print("0 ") // Нет стенки справа
			}
		}
		fmt.Println("1") // Правая граница для последней ячейки

		// Печатаем нижнюю границу ячейки
		for x := 0; x < m.Cols; x++ {
			if m.Cells[y][x].Bottom {
				fmt.Print("1 ") // Стенка снизу
			} else {
				fmt.Print("0 ") // Нет стенки снизу
			}
		}
		fmt.Println("1") // Нижняя граница для последней ячейки
	}
}
