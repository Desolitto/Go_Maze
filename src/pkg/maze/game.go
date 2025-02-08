package maze

import (
	"fmt"
	"go-maze/config"
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/exp/rand"
)

type Game struct {
	maze        *Maze
	cellSize    float32
	mazeSolving *MazeSolving
}

/**********************РЕШЕНИЕ************/

type MazeSolving struct {
	mazeInfo      *Maze
	startX        int
	startY        int
	endX          int
	endY          int
	solvingMatrix [][]int
	currentX      int
	currentY      int
	isStartSet    bool
	isEndSet      bool
	// path          []Point // Для хранения найденного пути
}

// type Point struct {
// 	x, y int
// }

func NewMazeSolving(maze *Maze, startX, startY, endX, endY int) *MazeSolving {
	solvingMatrix := make([][]int, maze.Rows)
	for i := range solvingMatrix {
		solvingMatrix[i] = make([]int, maze.Cols)
		for j := range solvingMatrix[i] {
			solvingMatrix[i][j] = config.MaxSize * config.MaxSize // Инициализируем максимальным значением
		}
	}
	return &MazeSolving{
		mazeInfo:      maze,
		startX:        startX,
		startY:        startY,
		endX:          endX,
		endY:          endY,
		solvingMatrix: solvingMatrix,
		currentX:      startX,
		currentY:      startY,
		isStartSet:    startX >= 0 && startY >= 0,
		isEndSet:      endX >= 0 && endY >= 0,
	}
}

// Метод Solve запускает решение лабиринта
func (m *MazeSolving) Solve() error {
	if !m.isStartSet || !m.isEndSet {
		return fmt.Errorf("начальная и конечная точки должны быть установлены")
	}
	return m.findPath(m.startX, m.startY)
}

// Метод findPath ищет путь в лабиринте
func (m *MazeSolving) findPath(x, y int) error {
	log.Printf("Проверяем клетку (%d, %d)", x, y)

	if x == m.endX && y == m.endY {
		m.solvingMatrix[y][x] = 1
		log.Printf("Достигнута конечная точка (%d, %d)", x, y)
		return nil
	}

	// Проверяем, если клетка уже посещена
	if m.solvingMatrix[y][x] == 1 {
		log.Printf("Клетка (%d, %d) уже посещена, пропускаем", x, y)
		return fmt.Errorf("путь не найден")
	}

	m.solvingMatrix[y][x] = 1 // Помечаем текущую клетку как посещенную

	// Перебираем возможные направления
	directions := [][3]int{
		{0, -1, 0},  // Вверх
		{1, 0, 1},   // Вправо
		{0, 1, 0},   // Вниз
		{-1, 0, -1}, // Влево
	}

	for _, dir := range directions {
		newX := x + dir[0]
		newY := y + dir[1]

		if m.isValidMove(newX, newY, dir[0], dir[1]) {
			log.Printf("Пытаемся двигаться в (%d, %d)", newX, newY)
			err := m.findPath(newX, newY)
			if err == nil {
				return nil
			}
		}
	}

	m.solvingMatrix[y][x] = 0 // Возвращаемся, если путь не найден
	log.Printf("Возврат из клетки (%d, %d) - путь не найден", x, y)
	return fmt.Errorf("путь не найден")
}

// Проверка, можно ли сделать ход
func (m *MazeSolving) isValidMove(x, y, dirX, dirY int) bool {
	// Проверяем границы
	if x < 0 || x >= m.mazeInfo.Cols || y < 0 || y >= m.mazeInfo.Rows {
		log.Printf("Недоступный ход в (%d, %d) - выход за границы", x, y)
		return false
	}

	// Проверяем наличие стен в зависимости от направления
	if dirX == 1 { // Движение вправо
		if m.mazeInfo.Cells[y][x].Right {
			log.Printf("Недоступный ход в (%d, %d) - есть правая стена", x, y)
			return false
		}
	} else if dirX == -1 { // Движение влево
		if x > 0 && m.mazeInfo.Cells[y][x-1].Right {
			log.Printf("Недоступный ход в (%d, %d) - есть левая стена", x, y)
			return false
		}
	}

	if dirY == 1 { // Движение вниз
		if m.mazeInfo.Cells[y][x].Bottom {
			log.Printf("Недоступный ход в (%d, %d) - есть нижняя стена", x, y)
			return false
		}
	} else if dirY == -1 { // Движение вверх
		if y > 0 && m.mazeInfo.Cells[y-1][x].Bottom {
			log.Printf("Недоступный ход в (%d, %d) - есть верхняя стена", x, y)
			return false
		}
	}

	log.Printf("Доступный ход в (%d, %d)", x, y)
	return true
}

/*************************************************************************************************/

func NewGame(rows, cols int) *Game {
	if rows > config.MaxSize || cols > config.MaxSize {
		log.Fatalf("Размер лабиринта не должен превышать %d", config.MaxSize)
	}
	ebiten.SetWindowSize(config.SceneWidth+int(config.BorderThickness*2), config.SceneHeight+config.ButtonHeight*3+int(config.BorderThickness))
	ebiten.SetWindowTitle("Cave Generator")
	cellSize := float32(config.SceneWidth) / float32(cols)
	maze := NewMaze(rows, cols)
	r := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))
	numRandomNumbers := rows * cols * 2
	randomNumbers := make([]int, numRandomNumbers)
	for i := range randomNumbers {
		randomNumbers[i] = r.Intn(2) // Генерация 0 или 1
	}
	/*
		// maze.Initialize(rows, cols)
		// maze.Generate(0, 0)
		// randomNumbers := make([]int, 0) // Для 4 строк по 4 столбца
		// randomNumbers = append(randomNumbers, 0, 1, 0, 0, 1, 1, 0, 1, 0, 0, 0, 0, 1, 1, 0, 1, 0, 1, 1, 0, 1, 0, 1, 0, 1, 1, 0, 0, 0, 0, 0, 1, 0, 1, 1, 0, 0, 1, 0, 1, 1, 0, 0, 0, 1, 0, 1, 1, 1, 0)
	*/
	maze.GenerateEller(randomNumbers)
	mazeSolving := NewMazeSolving(maze, -1, -1, -1, -1) // Начальные и конечные точки пока не установлены
	return &Game{maze: maze, cellSize: cellSize, mazeSolving: mazeSolving}
}

func (g *Game) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		cellX := int(float32(x) / g.cellSize)
		cellY := int(float32(y) / g.cellSize)

		if g.isInsideButton(float32(x), float32(y), float32(config.SceneHeight+config.BorderThickness), config.ButtonHeight) {
			go g.ShowFileSelector()
		} else if !g.mazeSolving.isStartSet {
			g.mazeSolving.startX = cellX
			g.mazeSolving.startY = cellY
			g.mazeSolving.isStartSet = true
			fmt.Printf("Начальная точка: (%d, %d)\n", cellX, cellY)
		} else if !g.mazeSolving.isEndSet {
			g.mazeSolving.endX = cellX
			g.mazeSolving.endY = cellY
			g.mazeSolving.isEndSet = true
			fmt.Printf("Конечная точка: (%d, %d)\n", cellX, cellY)
			// Запуск решения лабиринта после установки конечной точки
			if err := g.mazeSolving.Solve(); err != nil {
				log.Println("Ошибка при решении лабиринта:", err)
			}
		}
	}
	return nil
}

// Draw отрисовывает лабиринт и кнопку
func (g *Game) Draw(screen *ebiten.Image) {
	strokeColor := color.RGBA{0, 0, 0, 255}
	fillColor := color.RGBA{255, 255, 255, 255}
	mazeX := float32(0)
	mazeY := float32(0)
	// Рисуем лабиринт
	for y := 0; y < g.maze.Rows; y++ {
		for x := 0; x < g.maze.Cols; x++ {
			vector.DrawFilledRect(screen, mazeX+float32(x)*g.cellSize, mazeY+float32(y)*g.cellSize, g.cellSize, g.cellSize, fillColor, false)

			// Рисуем правую границу
			if x < g.maze.Cols-1 && g.maze.Cells[y][x].Right {
				vector.StrokeLine(screen, float32(x+1)*g.cellSize, float32(y)*g.cellSize, float32(x+1)*g.cellSize, float32(y+1)*g.cellSize, config.WallThickness, strokeColor, false)
			}

			// Рисуем нижнюю границу
			if y < g.maze.Rows-1 && g.maze.Cells[y][x].Bottom {
				vector.StrokeLine(screen, float32(x)*g.cellSize, float32(y+1)*g.cellSize, float32(x+1)*g.cellSize, float32(y+1)*g.cellSize, config.WallThickness, strokeColor, false)
			}
		}
	}
	g.drawButton(screen, "Open maze", float32(config.SceneHeight+config.BorderThickness), strokeColor)
	// Отрисовка начальной точки
	if g.mazeSolving.isStartSet {
		vector.DrawFilledRect(screen, float32(g.mazeSolving.startX)*g.cellSize, float32(g.mazeSolving.startY)*g.cellSize, g.cellSize, g.cellSize, color.RGBA{0, 255, 0, 255}, false) // Зеленый цвет
	}

	// Отрисовка конечной точки
	if g.mazeSolving.isEndSet {
		vector.DrawFilledRect(screen, float32(g.mazeSolving.endX)*g.cellSize, float32(g.mazeSolving.endY)*g.cellSize, g.cellSize, g.cellSize, color.RGBA{255, 0, 0, 255}, false) // Красный цвет
	}
	for y := 0; y < g.maze.Rows; y++ {
		for x := 0; x < g.maze.Cols; x++ {
			if g.mazeSolving.solvingMatrix[y][x] == 1 {
				vector.DrawFilledRect(screen, float32(x)*g.cellSize, float32(y)*g.cellSize, g.cellSize, g.cellSize, color.RGBA{0, 0, 255, 128}, false) // Синий цвет для пути
			}
		}
	}
}

func (g *Game) drawButton(screen *ebiten.Image, buttonText string, buttonY float32, color color.RGBA) {
	buttonWidth := float32(config.SceneWidth + config.BorderThickness*2)

	vector.DrawFilledRect(screen, 0, buttonY, buttonWidth, config.ButtonHeight, color, false)

	textWidth := float32(len(buttonText) * 8)
	textHeight := float32(16)

	textX := (buttonWidth - textWidth) / 2
	textY := buttonY + (config.ButtonHeight-textHeight)/2

	ebitenutil.DebugPrintAt(screen, buttonText, int(textX), int(textY))
}

// Layout определяет размер окна
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func (g *Game) Maze() *Maze {
	return g.maze
}
