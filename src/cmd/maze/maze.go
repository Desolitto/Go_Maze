// package main

// import (
// 	"flag"
// 	"fmt"
// 	"image/color"
// 	"log"
// 	"math/rand"

// 	"github.com/hajimehoshi/ebiten/v2"
// 	"github.com/hajimehoshi/ebiten/v2/vector"
// )

// type Maze struct {
// 	Rows              int
// 	Cols              int
// 	ActiveRow         int
// 	SetCounter        int
// 	SetMatrix         [][]int
// 	RightBorderMatrix [][]int
// 	LowBorderMatrix   [][]int
// }

// func NewMaze(rows int, cols int) *Maze {
// 	maze := &Maze{
// 		Rows:              rows,
// 		Cols:              cols,
// 		ActiveRow:         0,
// 		SetCounter:        1,
// 		SetMatrix:         make([][]int, rows),
// 		RightBorderMatrix: make([][]int, rows),
// 		LowBorderMatrix:   make([][]int, rows),
// 	}
// 	for i := range maze.SetMatrix {
// 		maze.SetMatrix[i] = make([]int, cols)
// 		maze.RightBorderMatrix[i] = make([]int, cols)
// 		maze.LowBorderMatrix[i] = make([]int, cols)
// 	}
// 	return maze
// }

// func (m *Maze) RowSize() int {
// 	return len(m.RightBorderMatrix)
// }

// func (m *Maze) ColSize() int {
// 	return len(m.RightBorderMatrix[0])
// }

// // Построчная генерация лабиринта с помощью алгоритма Эйлера
// func GenerateMaze(rows int, cols int) (*Maze, error) {
// 	if rows < 1 || cols < 1 {
// 		return nil, fmt.Errorf("rows and columns must be positive numbers")
// 	}
// 	if rows > 50 || cols > 50 {
// 		return nil, fmt.Errorf("rows and columns must be <= 50")
// 	}
// 	maze := NewMaze(rows, cols)
// 	for i := 0; i < rows; i++ {
// 		maze.assignUniqueSet()
// 		maze.addingVerticalWalls()
// 		maze.addingHorizontalWalls()
// 		maze.preparatingNewLine()
// 	}
// 	maze.addingEndLine()
// 	// maze.writeToFile()
// 	return maze, nil
// }

// // Присвоение ячейки множества
// func (m *Maze) assignUniqueSet() {
// 	for j := 0; j < m.Cols; j++ {
// 		if m.SetMatrix[m.ActiveRow][j] == 0 {
// 			m.SetMatrix[m.ActiveRow][j] = m.SetCounter
// 			m.SetCounter++
// 		}
// 	}
// }

// func (m *Maze) addingVerticalWalls() {
// 	for i := 0; i < m.Cols-1; i++ {
// 		choise := rand.Int() % 2
// 		if choise == 1 || m.SetMatrix[m.ActiveRow][i] == m.SetMatrix[m.ActiveRow][i+1] {
// 			m.RightBorderMatrix[m.ActiveRow][i] = 1
// 		} else {
// 			m.mergeSet(i)
// 		}
// 	}
// 	m.RightBorderMatrix[m.ActiveRow][m.Cols-1] = 1
// }

// // Объединение ячеек в одно множество
// func (m *Maze) mergeSet(i int) {
// 	x := m.SetMatrix[m.ActiveRow][i+1]
// 	for j := 0; j < m.Cols; j++ {
// 		if m.SetMatrix[m.ActiveRow][j] == x {
// 			m.SetMatrix[m.ActiveRow][j] = m.SetMatrix[m.ActiveRow][i]
// 		}
// 	}
// }

// // Добавление горизонтальных (нижних) стен
// func (m *Maze) addingHorizontalWalls() {
// 	for i := 0; i < m.Cols; i++ {
// 		choise := rand.Int() % 2
// 		check := m.checkedHorizontalWalls(i)
// 		if choise == 1 && check {
// 			m.LowBorderMatrix[m.ActiveRow][i] = 1
// 		}
// 	}
// }

// func (m *Maze) checkedHorizontalWalls(index int) bool {
// 	set := m.SetMatrix[m.ActiveRow][index]
// 	for i := 0; i < m.Cols; i++ {
// 		if m.SetMatrix[m.ActiveRow][i] == set && i != index && m.LowBorderMatrix[m.ActiveRow][i] == 0 {
// 			return true
// 		}
// 	}
// 	return false
// }

// func (m *Maze) preparatingNewLine() {
// 	if m.ActiveRow == m.Rows-1 {
// 		return
// 	}
// 	m.ActiveRow++
// 	for i := 0; i < m.Cols; i++ {
// 		if m.LowBorderMatrix[m.ActiveRow-1][i] == 0 {
// 			m.SetMatrix[m.ActiveRow][i] = m.SetMatrix[m.ActiveRow-1][i]
// 		} else {
// 			m.SetMatrix[m.ActiveRow][i] = 0
// 		}
// 	}
// }

// func (m *Maze) addingEndLine() {
// 	for i := 0; i < m.Cols-1; i++ {
// 		m.LowBorderMatrix[m.ActiveRow][i] = 1
// 		if m.SetMatrix[m.ActiveRow][i] != m.SetMatrix[m.ActiveRow][i+1] {
// 			m.RightBorderMatrix[m.ActiveRow][i] = 0
// 			m.mergeSet(i)
// 		}
// 	}
// 	m.LowBorderMatrix[m.ActiveRow][m.Cols-1] = 1
// 	m.RightBorderMatrix[m.ActiveRow][m.Cols-1] = 1
// }

// // Game структура для ebiten
// type Game struct {
// 	maze *Maze
// }

// // Новый экземпляр игры
// func NewGame(rows, cols int) *Game {
// 	maze, err := GenerateMaze(rows, cols)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return &Game{maze: maze}
// }

// // Update обновляет состояние игры
// func (g *Game) Update() error {
// 	return nil
// }

// // Draw отрисовывает лабиринт
// func (g *Game) Draw(screen *ebiten.Image) {
// 	cellWidth := float32(10.0)  // Изменено на float32
// 	cellHeight := float32(10.0) // Изменено на float32

// 	// Устанавливаем цвет стен
// 	strokeColor := color.RGBA{0, 0, 0, 255}     // Черный цвет
// 	fillColor := color.RGBA{255, 255, 255, 255} // Белый цвет

// 	// Рисуем ячейки и стены
// 	for y := 0; y < g.maze.Rows; y++ {
// 		for x := 0; x < g.maze.Cols; x++ {
// 			if g.maze.SetMatrix[y][x] == 1 {
// 				vector.DrawFilledRect(screen, float32(x)*cellWidth, float32(y)*cellHeight, cellWidth, cellHeight, fillColor, false)
// 			}
// 			if g.maze.RightBorderMatrix[y][x] == 1 {
// 				vector.StrokeLine(screen, float32(x+1)*cellWidth, float32(y)*cellHeight, float32(x+1)*cellWidth, float32(y+1)*cellHeight, 1, strokeColor, false)
// 			}
// 			if g.maze.LowBorderMatrix[y][x] == 1 {
// 				vector.StrokeLine(screen, float32(x)*cellWidth, float32(y+1)*cellHeight, float32(x+1)*cellWidth, float32(y+1)*cellHeight, 1, strokeColor, false)
// 			}
// 		}
// 	}
// }

// // Layout определяет размер окна
// func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
// 	return 500, 500 // Размер окна 500x500
// }

// func main() {

// 	w := flag.Int("w", 20, "width of the maze")
// 	h := flag.Int("h", 20, "height of the maze")
// 	flag.Parse()

// 	game := NewGame(*h, *w)
// 	fmt.Println("Исходная матрица:")
// 	for _, row := range game.maze.SetMatrix {
// 		fmt.Println(row)
// 	}

// 	if err := ebiten.RunGame(game); err != nil {
// 		log.Fatal(err)
// 	}
// }

// // func (m *Maze) LoadCaveMazeFile(filename string) {
// // 	file, err := os.Open(filename)
// // 	if err != nil {
// // 		log.Fatal(err)
// // 	}
// // 	defer file.Close()

// // 	scanner := bufio.NewScanner(file)

// // 	if scanner.Scan() {
// // 		dimensions := strings.Fields(scanner.Text())
// // 		if len(dimensions) != 2 {
// // 			log.Fatal("Неверный формат файла: первая строка должна содержать размеры пещеры.")
// // 		}

// // 		width, err := strconv.Atoi(dimensions[0])
// // 		if err != nil {
// // 			log.Fatal("Неверная ширина лабиринта.")
// // 		}

// // 		height, err := strconv.Atoi(dimensions[1])
// // 		if err != nil {
// // 			log.Fatal("Неверная высота лабиринта.")
// // 		}

// // 		m.Rows, m.Cols = height, width
// // 		m.Cave = cave.NewCave(width, height)
// // 		// Загрузка первой матрицы (стена справа)
// // 		m.RightBorderMatrix = make([][]int, height)

// // 		for y := 0; y < height; y++ {
// // 			if scanner.Scan() {
// // 				row := strings.Fields(scanner.Text())
// // 				if len(row) != width {
// // 					log.Fatal("Неверный формат файла: количество столбцов не совпадает с заданной шириной.")
// // 				}
// // 				m.RightBorderMatrix[y] = make([]int, width)
// // 				for x, cell := range row {
// // 					if cell == "0" {
// // 						m.Cave.Grid[y][x] = cave.Death
// // 					} else if cell == "1" {
// // 						m.Cave.Grid[y][x] = cave.Alive
// // 					} else {
// // 						log.Fatal("Неверный символ в пещере: должен быть 0 или 1.")
// // 					}
// // 				}
// // 			}
// // 		}
// // 		m.LowBorderMatrix = make([][]int, height)
// // 		for y := 0; y < height; y++ {
// // 			if scanner.Scan() {
// // 				row := strings.Fields(scanner.Text())
// // 				if len(row) != width {
// // 					log.Fatal("Неверный формат файла: количество столбцов не совпадает с заданной шириной.")
// // 				}
// // 				m.LowBorderMatrix[y] = make([]int, width)
// // 				for x, cell := range row {
// // 					value, err := strconv.Atoi(cell)
// // 					if err != nil {
// // 						log.Fatal("Неверный символ в матрице стен: должен быть целым числом.")
// // 					}
// // 					m.LowBorderMatrix[y][x] = value
// // 				}
// // 			}
// // 		}
// // 	}
// // }

// // func (m *Maze) PrintMaze() {
// // 	for y := 0; y < m.Rows; y++ {
// // 		for x := 0; x < m.Cols; x++ {
// // 			if m.Cave.Grid[y][x] == cave.Alive {
// // 				fmt.Print("1 ")
// // 			} else {
// // 				fmt.Print("0 ")
// // 			}
// // 		}
// // 		fmt.Println()
// // 	}
// // }

// // func main() {
// // 	m := &Maze{}
// // 	m.LoadCaveMazeFile("/Users/calamarp/Desktop/go/Go_Maze/src/assets/maze_test.txt")
// // 	m.PrintMaze()
// // }

// // func main() {
// // 	rand.Seed(uint64(time.Now().UnixNano()))
// // 	w := flag.Int("w", 20, "width of the cave")
// // 	h := flag.Int("h", 20, "height of the cave")
// // 	initialChance := flag.Int("с", 55, "initial chance (0-100)")
// // 	flag.Parse()

// // 	game := game.NewGame(*w, *h, *initialChance)
// // 	fmt.Println("Исходная матрица:")
// // 	// game.PrintCave()
// // 	if err := ebiten.RunGame(game); err != nil {
// // 		log.Fatal(err)
// // 	}
// // }

/*********************************************/

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

// func LoadMaze(filename string) (*Maze, error) {
// 	file, err := os.Open(filename)
// 	if err != nil {
// 		return nil, fmt.Errorf("ошибка при открытии файла: %v", err)
// 	}
// 	defer file.Close()

// 	scanner := bufio.NewScanner(file)

// 	// Читаем размеры лабиринта
// 	if !scanner.Scan() {
// 		return nil, fmt.Errorf("ошибка при чтении размеров лабиринта: %v", scanner.Err())
// 	}
// 	var rows, cols int
// 	_, err = fmt.Sscanf(scanner.Text(), "%d %d", &rows, &cols)
// 	if err != nil {
// 		return nil, fmt.Errorf("ошибка при парсинге размеров лабиринта: %v", err)
// 	}
// 	fmt.Printf("Размеры лабиринта: %d строк, %d столбцов\n", rows, cols)

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
// 		if !scanner.Scan() {
// 			return nil, fmt.Errorf("ошибка при чтении стенок справа в строке %d: %v", y, scanner.Err())
// 		}
// 		for x := 0; x < cols; x++ {
// 			var wall int
// 			_, err = fmt.Sscanf(scanner.Text()[x*2:x*2+1], "%d", &wall) // Предполагаем, что данные разделены пробелами
// 			if err != nil {
// 				return nil, fmt.Errorf("ошибка при парсинге стенки справа в строке %d, столбце %d: %v", y, x, err)
// 			}
// 			if wall == 1 {
// 				maze.Cells[y][x].Right = true

// 			} // } else {
// 			// 	// Если стенки нет, помечаем ячейку как посещенную
// 			// 	maze.Cells[y][x].Visited = true
// 			// }
// 			fmt.Printf("Строка %d, столбец %d: стенка справа = %d\n", y, x, wall)
// 		}
// 	}

// 	// Пропускаем пустую строку между матрицами
// 	if !scanner.Scan() {
// 		return nil, fmt.Errorf("ошибка при чтении пустой строки между матрицами: %v", scanner.Err())
// 	}

// 	// Читаем вторую матрицу (стенки снизу)
// 	for y := 0; y < rows; y++ {
// 		if !scanner.Scan() {
// 			return nil, fmt.Errorf("ошибка при чтении стенок снизу в строке %d: %v", y, scanner.Err())
// 		}
// 		for x := 0; x < cols; x++ {
// 			var wall int
// 			_, err = fmt.Sscanf(scanner.Text()[x*2:x*2+1], "%d", &wall) // Предполагаем, что данные разделены пробелами
// 			if err != nil {
// 				return nil, fmt.Errorf("ошибка при парсинге стенки снизу в строке %d, столбце %d: %v", y, x, err)
// 			}
// 			if wall == 1 {
// 				maze.Cells[y][x].Bottom = true
// 			}
// 			maze.Cells[y][x].Visited = true
// 			fmt.Printf("Строка %d, столбец %d: стенка снизу = %d\n", y, x, wall)
// 		}
// 	}

// 	fmt.Println("Загрузка лабиринта завершена успешно.")
// 	return maze, nil
// }
