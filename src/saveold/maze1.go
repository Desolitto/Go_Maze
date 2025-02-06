// package main

// import (
// 	"bufio"
// 	"flag"
// 	"fmt"
// 	"image/color"
// 	"log"
// 	"math/rand"
// 	"os"
// 	"time"

// 	"github.com/hajimehoshi/ebiten/v2"
// 	"github.com/hajimehoshi/ebiten/v2/inpututil"
// 	"github.com/hajimehoshi/ebiten/v2/vector"
// 	"github.com/sqweek/dialog"
// )

// const (
// 	maxMazeSize     = 50
// 	wallThickness   = 2
// 	mazeWidth       = 500
// 	mazeHeight      = 500 // Высота лабиринта
// 	buttonHeight    = 30
// 	borderThickness = float32(2)
// )

// type Cell struct {
// 	Visited bool
// 	Right   bool
// 	Bottom  bool
// 	Set     int
// }

// type Maze struct {
// 	Rows, Cols int
// 	Cells      [][]Cell
// }

// type Game struct {
// 	maze       *Maze
// 	cellSize   float32
// 	loadButton bool
// }

// // func NewMaze(rows, cols int) *Maze {
// // 	ebiten.SetWindowSize(mazeWidth+int(borderThickness*2), mazeHeight+buttonHeight+int(borderThickness))
// // 	cells := make([][]Cell, rows)
// // 	for i := range cells {
// // 		cells[i] = make([]Cell, cols)
// // 		for j := range cells[i] {
// // 			cells[i][j] = Cell{Visited: false, Right: true, Bottom: true}
// // 		}
// // 	}
// // 	return &Maze{Rows: rows, Cols: cols, Cells: cells}
// // }

// // func NewMaze(rows, cols int) *Maze {
// // 	ebiten.SetWindowSize(mazeWidth+int(borderThickness*2), mazeHeight+buttonHeight+int(borderThickness))

// //		cells := make([][]Cell, rows)
// //		for i := range cells {
// //			cells[i] = make([]Cell, cols)
// //		}
// //		return &Maze{Rows: rows, Cols: cols, Cells: cells}
// //	}
// func NewMaze(rows, cols int) *Maze {
// 	cells := make([][]Cell, rows)
// 	for i := range cells {
// 		cells[i] = make([]Cell, cols)
// 		for j := range cells[i] {
// 			cells[i][j] = Cell{Visited: false, Right: true, Bottom: true, Set: i*cols + j}
// 		}
// 	}
// 	return &Maze{Rows: rows, Cols: cols, Cells: cells}
// }

// // Инициализация лабиринта
// func (m *Maze) Initialize(rows, cols int) {
// 	m.Rows = rows
// 	m.Cols = cols
// 	m.Cells = make([][]Cell, rows)

// 	for y := 0; y < rows; y++ {
// 		m.Cells[y] = make([]Cell, cols)
// 		for x := 0; x < cols; x++ {
// 			// Устанавливаем все стенки по умолчанию
// 			m.Cells[y][x].Right = true
// 			m.Cells[y][x].Bottom = true
// 		}
// 	}
// }

// func (m *Maze) Generate(x, y int) {
// 	// m.Cells[y][x].Visited = true

// 	directions := []struct {
// 		dx, dy int
// 	}{
// 		{1, 0},  // вправо
// 		{0, 1},  // вниз
// 		{-1, 0}, // влево
// 		{0, -1}, // вверх
// 	}

// 	rand.Shuffle(len(directions), func(i, j int) {
// 		directions[i], directions[j] = directions[j], directions[i]
// 	})

// 	for _, dir := range directions {
// 		newX, newY := x+dir.dx, y+dir.dy
// 		if newX >= 0 && newX < m.Cols && newY >= 0 && newY < m.Rows && !m.Cells[newY][newX].Visited {
// 			if dir.dx == 1 { // вправо
// 				m.Cells[y][x].Right = false
// 			} else if dir.dy == 1 { // вниз
// 				m.Cells[y][x].Bottom = false
// 			} else if dir.dx == -1 { // влево
// 				m.Cells[newY][newX].Right = false
// 			} else if dir.dy == -1 { // вверх
// 				m.Cells[newY][newX].Bottom = false
// 			}
// 			m.Generate(newX, newY)
// 		}
// 	}
// }

// func (m *Maze) printRow(row int) {
// 	fmt.Printf("Row %d:\n", row)
// 	for col := 0; col < m.Cols; col++ {
// 		fmt.Printf("Cell (%d, %d): Set=%d, Right=%t, Bottom=%t\n",
// 			row, col, m.Cells[row][col].Set, m.Cells[row][col].Right, m.Cells[row][col].Bottom)
// 	}
// 	fmt.Println()
// }

// func (m *Maze) GenerateEller(randomNumbers []int) {
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

// 	for row := 0; row < m.Rows; row++ {
// 		fmt.Println(row)
// 		if row > 0 {
// 			for col := 0; col < m.Cols; col++ {
// 				m.Cells[row][col].Right = m.Cells[row-1][col].Right
// 				m.Cells[row][col].Bottom = m.Cells[row-1][col].Bottom
// 				m.Cells[row][col].Set = m.Cells[row-1][col].Set
// 			}
// 			// Удаляем правые стенки и нижние границы
// 			for col := 0; col < m.Cols; col++ {
// 				m.Cells[row][col].Right = false
// 				if m.Cells[row-1][col].Bottom {
// 					m.Cells[row][col].Set = 0        // Присваиваем пустое множество
// 					m.Cells[row][col].Bottom = false // Удаляем нижнюю стенку
// 				}
// 			}
// 			// Присваиваем новые множества для следующей строки
// 			for col := 0; col < m.Cols; col++ {
// 				if m.Cells[row][col].Set == 0 {
// 					m.Cells[row][col].Set = currentSetCount
// 					currentSetCount++
// 					fmt.Printf("1Присвоено новое множество ячейке (%d, %d): Set=%d\n", row, col, m.Cells[row][col].Set)
// 				}
// 			}
// 		}

// 		fmt.Printf("ПЕРЕД УСТАНОВКОЙ СТЕНКИ:\nСтрока %d", row)
// 		for col := 0; col < m.Cols; col++ {
// 			fmt.Printf("{R: %v, B: %v, Set: %d} ", m.Cells[row][col].Right, m.Cells[row][col].Bottom, m.Cells[row][col].Set)
// 		}
// 		fmt.Println()

// 		// Обработка правых стенок
// 		for col := 0; col < m.Cols-1; col++ {
// 			fmt.Printf("Перед установкой стенки: Cell(%d, %d) Set=%d\n\n", row, col, m.Cells[row][col].Set)
// 			fmt.Printf("randomNumbers[index] right = %d\n", randomNumbers[index])
// 			if randomNumbers[index] == 1 {
// 				// Ставим стенку
// 				m.Cells[row][col].Right = true
// 				fmt.Printf("После установкой стенки: Cell(%d, %d) Set=%d\n", row, col, m.Cells[row][col].Set)
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
// 					m.Cells[row][col].Right = true
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
// 				if m.Cells[row][c].Set == set && !m.Cells[row][c].Bottom {
// 					count++
// 				}
// 			}

// 			if count > 1 {
// 				fmt.Printf("randomNumbers[index] bottom = %d\n", randomNumbers[index])
// 				if randomNumbers[index] == 1 {
// 					m.Cells[row][col].Bottom = true
// 				}
// 			}
// 			index++
// 		}

// 		// Если это последняя строка, добавляем нижние стенки
// 		if row == m.Rows-1 {
// 			for col := 0; col < m.Cols; col++ {
// 				m.Cells[row][col].Bottom = true
// 			}
// 			// Двигаясь слева направо, убираем стенки между ячейками, если множества не совпадают
// 			for col := 0; col < m.Cols-1; col++ {
// 				set1 := m.Cells[row][col].Set
// 				set2 := m.Cells[row][col+1].Set

// 				if set1 != set2 {
// 					// Убираем стенку между ячейками
// 					m.Cells[row][col].Right = false
// 					// Объединяем множества
// 					for r := 0; r < m.Rows; r++ {
// 						for c := 0; c < m.Cols; c++ {
// 							if m.Cells[r][c].Set == set2 {
// 								m.Cells[r][c].Set = set1
// 							}
// 						}
// 					}
// 				}
// 			}
// 		}

// 		// Печатаем измененную строку
// 		fmt.Print("Измененая строка: [")
// 		for col := 0; col < m.Cols; col++ {
// 			fmt.Printf("{RightWall: %v, Bottom: %v, Set: %d}", m.Cells[row][col].Right, m.Cells[row][col].Bottom, m.Cells[row][col].Set)
// 			if col < m.Cols-1 {
// 				fmt.Print(" ")
// 			}
// 		}
// 		fmt.Println("]")
// 	}

// }

// // func (m *Maze) GenerateEller() {
// // 	r := rand.New(rand.NewSource(time.Now().UnixNano()))
// // 	for row := 0; row < m.Rows; row++ {
// // 		// Step 1: Initialize the sets for each cell in the current row
// // 		sets := make(map[int]int)
// // 		nextSetID := 1
// // 		for col := 0; col < m.Cols; col++ {
// // 			if row == 0 || m.Cells[row-1][col].Set == 0 {
// // 				sets[col] = nextSetID
// // 				m.Cells[row][col].Set = nextSetID
// // 				nextSetID++
// // 			} else {
// // 				m.Cells[row][col].Set = m.Cells[row-1][col].Set
// // 			}
// // 		}
// // 		// Print initial sets for the row
// // 		fmt.Printf("Initial sets for Row %d:\n", row)
// // 		for col := 0; col < m.Cols; col++ {
// // 			fmt.Printf("Cell (%d, %d): Set=%d\n", row, col, m.Cells[row][col].Set)
// // 		}
// // 		fmt.Println()

// // 		// Step 2: Decide whether to create vertical walls between cells in the same row
// // 		for col := 0; col < m.Cols-1; col++ {
// // 			if r.Float32() < 0.5 {
// // 				if m.Cells[row][col].Set != m.Cells[row][col+1].Set {
// // 					mergeSets(sets, m.Cells[row][col].Set, m.Cells[row][col+1].Set)
// // 					m.Cells[row][col].Right = false
// // 					fmt.Printf("No right wall between (%d, %d) and (%d, %d), merged sets %d and %d\n",
// // 						row, col, row, col+1, m.Cells[row][col].Set, m.Cells[row][col+1].Set)
// // 				} else {
// // 					m.Cells[row][col].Right = true
// // 					fmt.Printf("Right wall between (%d, %d) and (%d, %d)\n", row, col, row, col+1)
// // 				}
// // 			} else {
// // 				m.Cells[row][col].Right = true
// // 				fmt.Printf("Right wall between (%d, %d) and (%d, %d)\n", row, col, row, col+1)
// // 			}
// // 		}

// // 		// Step 3: Decide whether to create horizontal walls between cells in the same column
// // 		if row < m.Rows-1 {
// // 			horizWalls := make([]bool, m.Cols)
// // 			setCounts := make(map[int]int)
// // 			for col := 0; col < m.Cols; col++ {
// // 				setCounts[m.Cells[row][col].Set]++
// // 			}
// // 			for col := 0; col < m.Cols; col++ {
// // 				if setCounts[m.Cells[row][col].Set] > 1 && r.Float32() < 0.5 {
// // 					horizWalls[col] = true
// // 				} else {
// // 					horizWalls[col] = false
// // 					m.Cells[row][col].Bottom = false
// // 					fmt.Printf("No bottom wall for cell (%d, %d)\n", row, col)
// // 				}
// // 			}
// // 			for col := 0; col < m.Cols; col++ {
// // 				if horizWalls[col] {
// // 					m.Cells[row][col].Bottom = true
// // 					fmt.Printf("Bottom wall for cell (%d, %d)\n", row, col)
// // 				}
// // 			}
// // 		} else {
// // 			// Last row: Remove all closed areas
// // 			for col := 0; col < m.Cols-1; col++ {
// // 				if m.Cells[row][col].Set != m.Cells[row][col+1].Set {
// // 					m.Cells[row][col].Right = false
// // 					mergeSets(sets, m.Cells[row][col].Set, m.Cells[row][col+1].Set)
// // 					fmt.Printf("No right wall between (%d, %d) and (%d, %d), merged sets %d and %d\n",
// // 						row, col, row, col+1, m.Cells[row][col].Set, m.Cells[row][col+1].Set)
// // 				} else {
// // 					m.Cells[row][col].Right = true
// // 					fmt.Printf("Right wall between (%d, %d) and (%d, %d)\n", row, col, row, col+1)
// // 				}
// // 			}
// // 		}

// // 		// Print the state of the row after processing
// // 		m.printRow(row)
// // 	}
// // }

// func mergeSets(sets map[int]int, a, b int) {
// 	target := sets[a]
// 	replacement := sets[b]
// 	if target == replacement {
// 		return
// 	}
// 	for k, v := range sets {
// 		if v == replacement {
// 			sets[k] = target
// 		}
// 	}
// }

// func NewGame(rows, cols int) *Game {
// 	r := rand.New(rand.NewSource(time.Now().UnixNano()))
// 	numRandomNumbers := rows * cols * 2
// 	randomNumbers := make([]int, numRandomNumbers)
// 	for i := range randomNumbers {
// 		randomNumbers[i] = r.Intn(2) // Генерация 0 или 1
// 	}
// 	maze := NewMaze(rows, cols)
// 	// maze.Generate(rows, cols)
// 	maze.GenerateEller(randomNumbers)
// 	cellSize := float32(mazeWidth) / float32(cols)
// 	return &Game{maze: maze, cellSize: cellSize}
// }

// // Update обновляет состояние игры
// func (g *Game) Update() error {
// 	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
// 		x, y := ebiten.CursorPosition()

// 		if g.isInsideButton(float32(x), float32(y), float32(mazeHeight+borderThickness), buttonHeight) {
// 			go g.ShowFileSelector()
// 		}
// 	}
// 	return nil
// }

// func (g *Game) isInsideButton(x, y float32, buttonY float32, buttonHeight float32) bool {
// 	buttonX := float32(0)
// 	buttonWidth := float32(mazeHeight + borderThickness*2)
// 	return x >= buttonX && x <= buttonX+buttonWidth && y >= buttonY && y <= buttonY+buttonHeight
// }

// func (g *Game) ShowFileSelector() {
// 	currentDir, err := os.Getwd()
// 	if err != nil {
// 		log.Println("Ошибка при получении текущей директории:", err)
// 		return
// 	}

// 	filename, err := dialog.File().
// 		Filter("Text files", "txt").
// 		SetStartDir(currentDir).
// 		Load()

// 	if err != nil {
// 		log.Println("Ошибка при выборе файла:", err)
// 		return
// 	}

// 	// Загружаем лабиринт из выбранного файла
// 	maze, err := LoadMaze(filename)
// 	if err != nil {
// 		log.Println("Ошибка при загрузке лабиринта:", err)
// 		return
// 	}

// 	// Обновляем состояние игры с новым лабиринтом
// 	g.maze = maze
// }

// // Draw отрисовывает лабиринт и кнопку
// func (g *Game) Draw(screen *ebiten.Image) {
// 	strokeColor := color.RGBA{0, 0, 0, 255}
// 	fillColor := color.RGBA{255, 255, 255, 255}

// 	// Рисуем лабиринт
// 	for y := 0; y < g.maze.Rows; y++ {
// 		for x := 0; x < g.maze.Cols; x++ {
// 			// Рисуем ячейку
// 			// if g.maze.Cells[y][x].Visited {
// 			vector.DrawFilledRect(screen, float32(x)*g.cellSize, float32(y)*g.cellSize, g.cellSize, g.cellSize, fillColor, false)
// 			// }

// 			// Рисуем правую границу
// 			if x < g.maze.Cols-1 && g.maze.Cells[y][x].Right {
// 				vector.StrokeLine(screen, float32(x+1)*g.cellSize, float32(y)*g.cellSize, float32(x+1)*g.cellSize, float32(y+1)*g.cellSize, wallThickness, strokeColor, false)
// 			}

// 			// Рисуем нижнюю границу
// 			if y < g.maze.Rows-1 && g.maze.Cells[y][x].Bottom {
// 				vector.StrokeLine(screen, float32(x)*g.cellSize, float32(y+1)*g.cellSize, float32(x+1)*g.cellSize, float32(y+1)*g.cellSize, wallThickness, strokeColor, false)
// 			}
// 		}
// 	}

// 	// Рисуем кнопку под лабиринтом
// 	buttonY := mazeHeight
// 	buttonWidth := mazeWidth
// 	vector.DrawFilledRect(screen, 0, float32(buttonY), float32(buttonWidth), float32(buttonHeight), color.RGBA{200, 200, 200, 255}, false)
// 	vector.StrokeLine(screen, 0, float32(buttonY), float32(buttonWidth), float32(buttonY), wallThickness, strokeColor, false)
// }

// // Layout определяет размер окна
// func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
// 	return mazeWidth, mazeHeight + buttonHeight // Общая высота с кнопкой
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

// 			}
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

// // SaveMaze сохраняет лабиринт в файл в указанном формате
// func (m *Maze) SaveMaze(filename string) error {
// 	file, err := os.Create(filename)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	// Записываем размеры лабиринта
// 	_, err = fmt.Fprintf(file, "%d %d\n", m.Rows, m.Cols)
// 	if err != nil {
// 		return err
// 	}

// 	// Записываем стенки справа
// 	for y := 0; y < m.Rows; y++ {
// 		for x := 0; x < m.Cols; x++ {
// 			if x < m.Cols {
// 				if m.Cells[y][x].Right {
// 					_, err = fmt.Fprintf(file, "1 ")
// 				} else {
// 					_, err = fmt.Fprintf(file, "0 ")
// 				}
// 			} else {
// 				// Для последнего элемента в строке добавляем "1", чтобы закрыть строку
// 				// _, err = fmt.Fprintf(file, "1")
// 			}
// 		}
// 		_, err = fmt.Fprintln(file) // Переход на новую строку
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	// Добавляем пробел между матрицами
// 	_, err = fmt.Fprintln(file)
// 	if err != nil {
// 		return err
// 	}

// 	// Записываем стенки снизу
// 	for y := 0; y < m.Rows; y++ {
// 		for x := 0; x < m.Cols; x++ {
// 			if x < m.Cols {
// 				if m.Cells[y][x].Bottom {
// 					_, err = fmt.Fprintf(file, "1 ")
// 				} else {
// 					_, err = fmt.Fprintf(file, "0 ")
// 				}
// 			} else {
// 				// Для последнего элемента в строке добавляем "0", чтобы закрыть строку
// 				// _, err = fmt.Fprintf(file, "0")
// 			}
// 		}
// 		_, err = fmt.Fprintln(file) // Переход на новую строку
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// func main() {
// 	w := flag.Int("w", maxMazeSize, "количество строк в лабиринте")
// 	h := flag.Int("h", maxMazeSize, "количество столбцов в лабиринте")
// 	flag.Parse()

// 	game := NewGame(*w, *h)
// 	// Печатаем сгенерированный лабиринт в терминал
// 	fmt.Println("Сгенерированный лабиринт:")
// 	// game.maze.PrintMaze()
// 	err := game.maze.SaveMaze("maze.txt")
// 	if err != nil {
// 		fmt.Println("Ошибка при сохранении лабиринта:", err)
// 	} else {
// 		fmt.Println("Лабиринт успешно сохранен в maze.txt")
// 	}
// 	if err := ebiten.RunGame(game); err != nil {
// 		log.Fatal(err)
// 	}
// }

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
