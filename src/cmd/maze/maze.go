package main

import (
	"bufio"
	"flag"
	"fmt"
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sqweek/dialog"
)

const (
	maxSize         = 50
	wallThickness   = 2
	sceneWidth      = 500
	sceneHeight     = 500 // Высота лабиринта
	buttonHeight    = 30
	borderThickness = float32(2)
)

type Cell struct {
	Right  bool
	Bottom bool
	Set    int
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
	ebiten.SetWindowSize(sceneWidth+int(borderThickness*2), sceneHeight+buttonHeight+int(borderThickness))
	cells := make([][]Cell, rows)
	for i := range cells {
		cells[i] = make([]Cell, cols)
		for j := range cells[i] {
			cells[i][j] = Cell{Right: false, Bottom: false, Set: -1}
		}
	}
	return &Maze{Rows: rows, Cols: cols, Cells: cells}
}

func (m *Maze) copyPreviousRow(row int, currentSetCount *int) {
	for col := 0; col < m.Cols; col++ {
		m.Cells[row][col].Right = m.Cells[row-1][col].Right
		m.Cells[row][col].Bottom = m.Cells[row-1][col].Bottom
		m.Cells[row][col].Set = m.Cells[row-1][col].Set
	}
	for col := 0; col < m.Cols; col++ {
		m.Cells[row][col].Right = false
		if m.Cells[row-1][col].Bottom {
			m.Cells[row][col].Set = 0        // Присваиваем пустое множество
			m.Cells[row][col].Bottom = false // Удаляем нижнюю стенку
		}
	}
	for col := 0; col < m.Cols; col++ {
		if m.Cells[row][col].Set == 0 {
			m.Cells[row][col].Set = (*currentSetCount)
			(*currentSetCount)++
		}
	}

}

func (m *Maze) initializeSets() {
	for row := 0; row < m.Rows; row++ {
		for col := 0; col < m.Cols; col++ {
			m.Cells[row][col].Set = row*m.Cols + col + 1
		}
	}
}

func (m *Maze) setFirstRowSets(currentSetCount *int) {
	for col := 0; col < m.Cols; col++ {
		m.Cells[0][col].Set = col + 1
		(*currentSetCount)++
	}
}

func (m *Maze) processRightWalls(row int, randomNumbers []int, index *int) {
	for col := 0; col < m.Cols-1; col++ {
		if randomNumbers[*index] == 1 {
			m.Cells[row][col].Right = true
		} else {
			set1 := m.Cells[row][col].Set
			set2 := m.Cells[row][col+1].Set
			if set1 != set2 {
				m.mergeSets(set1, set2)
			} else {
				m.Cells[row][col].Right = true
			}
		}
		(*index)++
	}
}

func (m *Maze) processBottomWalls(row int, randomNumbers []int, index *int) {
	for col := 0; col < m.Cols; col++ {
		set := m.Cells[row][col].Set
		count := m.countCellsWithoutBottom(set)

		if count > 1 && randomNumbers[*index] == 1 {
			m.Cells[row][col].Bottom = true
		}
		(*index)++
	}
}

func (m *Maze) countCellsWithoutBottom(set int) int {
	count := 0
	for c := 0; c < m.Cols; c++ {
		if m.Cells[m.Rows-1][c].Set == set && !m.Cells[m.Rows-1][c].Bottom {
			count++
		}
	}
	return count
}

func (m *Maze) addBottomWalls(row int) {
	for col := 0; col < m.Cols; col++ {
		m.Cells[row][col].Bottom = true
	}
}

func (m *Maze) GenerateEller(randomNumbers []int) {
	m.initializeSets()

	currentSetCount := 1
	m.setFirstRowSets(&currentSetCount)
	index := 0
	for row := 0; row < m.Rows; row++ {
		if row > 0 {
			m.copyPreviousRow(row, &currentSetCount)
		}
		m.processRightWalls(row, randomNumbers, &index)
		// // Обработка правых стенок
		// for col := 0; col < m.Cols-1; col++ {
		// 	if randomNumbers[index] == 1 {
		// 		// Ставим стенку
		// 		m.Cells[row][col].Right = true
		// 	} else {
		// 		set1 := m.Cells[row][col].Set
		// 		set2 := m.Cells[row][col+1].Set

		// 		if set1 != set2 {
		// 			for r := 0; r < m.Rows; r++ {
		// 				for c := 0; c < m.Cols; c++ {
		// 					if m.Cells[r][c].Set == set2 {
		// 						m.Cells[r][c].Set = set1
		// 					}
		// 				}
		// 			}
		// 		} else {
		// 			m.Cells[row][col].Right = true
		// 		}
		// 	}
		// 	index++
		// }

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
			m.addBottomWalls(row)
			m.mergeLastRowSets(row)
		}
	}

}

func (m *Maze) mergeLastRowSets(row int) {
	for col := 0; col < m.Cols-1; col++ {
		set1 := m.Cells[row][col].Set
		set2 := m.Cells[row][col+1].Set
		if set1 != set2 {
			m.Cells[row][col].Right = false
			m.mergeSets(set1, set2)
		}
	}
}

func (m *Maze) mergeSets(set1, set2 int) {
	for r := 0; r < m.Rows; r++ {
		for c := 0; c < m.Cols; c++ {
			if m.Cells[r][c].Set == set2 {
				m.Cells[r][c].Set = set1
			}
		}
	}
}

func NewGame(rows, cols int) *Game {
	maze := NewMaze(rows, cols)
	// maze.Initialize(rows, cols)
	// maze.Generate(0, 0)
	// r := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))
	// numRandomNumbers := rows * cols * 2
	// randomNumbers := make([]int, numRandomNumbers)
	// for i := range randomNumbers {
	// 	randomNumbers[i] = r.Intn(2) // Генерация 0 или 1
	// }
	randomNumbers := make([]int, 0) // Для 4 строк по 4 столбца
	randomNumbers = append(randomNumbers, 0, 1, 0, 0, 1, 1, 0, 1, 0, 0, 0, 0, 1, 1, 0, 1, 0, 1, 1, 0, 1, 0, 1, 0, 1, 1, 0, 0, 0, 0, 0, 1, 0, 1, 1, 0, 0, 1, 0, 1, 1, 0, 0, 0, 1, 0, 1, 1, 1, 0)
	maze.GenerateEller(randomNumbers)
	cellSize := float32(sceneWidth) / float32(cols)
	return &Game{maze: maze, cellSize: cellSize}
}

// Update обновляет состояние игры
func (g *Game) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()

		if g.isInsideButton(float32(x), float32(y), float32(sceneHeight+borderThickness), buttonHeight) {
			go g.ShowFileSelector()
		}
	}
	return nil
}

func (g *Game) isInsideButton(x, y float32, buttonY float32, buttonHeight float32) bool {
	buttonX := float32(0)
	buttonWidth := float32(sceneHeight + borderThickness*2)
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
			vector.DrawFilledRect(screen, float32(x)*g.cellSize, float32(y)*g.cellSize, g.cellSize, g.cellSize, fillColor, false)

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
	g.drawButton(screen, "Open maze", float32(sceneHeight+borderThickness), strokeColor)
}

func (g *Game) drawButton(screen *ebiten.Image, buttonText string, buttonY float32, color color.RGBA) {
	buttonWidth := float32(sceneWidth + borderThickness*2)
	buttonHeight := float32(30)

	vector.DrawFilledRect(screen, 0, buttonY, buttonWidth, buttonHeight, color, false)

	textWidth := float32(len(buttonText) * 8)
	textHeight := float32(16)

	textX := (buttonWidth - textWidth) / 2
	textY := buttonY + (buttonHeight-textHeight)/2

	ebitenutil.DebugPrintAt(screen, buttonText, int(textX), int(textY))
}

// Layout определяет размер окна
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
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

	if _, err = fmt.Fprintf(file, "%d %d\n", m.Rows, m.Cols); err != nil {
		return err
	}

	if err := m.writeWalls(file, true); err != nil {
		return err
	}

	if _, err = fmt.Fprintln(file); err != nil {
		return err
	}

	return m.writeWalls(file, false)
}

func (m *Maze) writeWalls(file *os.File, isRight bool) error {
	for y := 0; y < m.Rows; y++ {
		for x := 0; x < m.Cols; x++ {
			var wall bool
			if isRight {
				wall = m.Cells[y][x].Right
			} else {
				wall = m.Cells[y][x].Bottom
			}
			if _, err := fmt.Fprintf(file, "%d ", boolToInt(wall)); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintln(file); err != nil {
			return err
		}
	}
	return nil
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func main() {
	w := flag.Int("w", maxSize, "количество строк в лабиринте")
	h := flag.Int("h", maxSize, "количество столбцов в лабиринте")
	flag.Parse()

	game := NewGame(*w, *h)
	fmt.Println("Сгенерированный лабиринт:")
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

/* ================== old code ========================== */

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
// 	visited := make([][]bool, m.Rows)
// 	for i := range visited {
// 		visited[i] = make([]bool, m.Cols)
// 	}

// 	stack := []struct{ x, y int }{{x, y}}
// 	visited[y][x] = true

// 	directions := []struct {
// 		dx, dy int
// 	}{
// 		{1, 0},  // вправо
// 		{0, 1},  // вниз
// 		{-1, 0}, // влево
// 		{0, -1}, // вверх
// 	}

// 	for len(stack) > 0 {
// 		curr := stack[len(stack)-1]
// 		stack = stack[:len(stack)-1]

// 		rand.Shuffle(len(directions), func(i, j int) {
// 			directions[i], directions[j] = directions[j], directions[i]
// 		})

// 		for _, dir := range directions {
// 			newX, newY := curr.x+dir.dx, curr.y+dir.dy
// 			if newX >= 0 && newX < m.Cols && newY >= 0 && newY < m.Rows && !visited[newY][newX] {
// 				if dir.dx == 1 { // вправо
// 					m.Cells[curr.y][curr.x].Right = false
// 				} else if dir.dy == 1 { // вниз
// 					m.Cells[curr.y][curr.x].Bottom = false
// 				} else if dir.dx == -1 { // влево
// 					m.Cells[newY][newX].Right = false
// 				} else if dir.dy == -1 { // вверх
// 					m.Cells[newY][newX].Bottom = false
// 				}

// 				visited[newY][newX] = true
// 				stack = append(stack, struct{ x, y int }{newX, newY})
// 			}
// 		}
// 	}
// }
