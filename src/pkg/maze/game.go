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
	path          []Point // Для хранения найденного пути
}

type Point struct {
	X int
	Y int
}

// type Point struct {
// 	x, y int
// }

// Конструктор MazeSolving
func NewMazeSolving(maze *Maze, startX, startY, endX, endY int) *MazeSolving {
	solvingMatrix := make([][]int, maze.Rows)
	for i := range solvingMatrix {
		solvingMatrix[i] = make([]int, maze.Cols)
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
		path:          []Point{},
	}
}

// Метод Solve запускает решение лабиринта
func (m *MazeSolving) Solve() error {
	if !m.isStartSet || !m.isEndSet {
		return fmt.Errorf("начальная и конечная точки должны быть установлены")
	}
	m.path = []Point{} // Очищаем предыдущий путь
	log.Println("Запуск решения лабиринта...")
	if m.dfs(m.startX, m.startY) {
		log.Println("Путь найден:", m.path)
		return nil
	}
	return fmt.Errorf("не удалось найти путь")
}

func (m *MazeSolving) dfs(x, y int) bool {
	if x < 0 || x >= m.mazeInfo.Cols || y < 0 || y >= m.mazeInfo.Rows {
		log.Printf("Выход за границы: (%d, %d)\n", x, y)
		return false // Выход за границы
	}
	if m.solvingMatrix[y][x] == 1 {
		log.Printf("Уже посещенная клетка: (%d, %d)\n", x, y)
		return false // Уже посещенная клетка
	}
	if x == m.endX && y == m.endY {
		m.path = append(m.path, Point{X: x, Y: y}) // Добавляем конечную точку в путь
		log.Printf("Найдена конечная точка: (%d, %d)\n", x, y)
		return true // Найден путь
	}

	// Помечаем клетку как посещенную
	m.solvingMatrix[y][x] = 1
	m.path = append(m.path, Point{X: x, Y: y}) // Добавляем текущую точку в путь
	log.Printf("Проверяем клетку: (%d, %d)\n", x, y)

	// Пробуем двигаться в четырех направлениях, проверяя наличие стен
	if x < m.mazeInfo.Cols-1 && !m.mazeInfo.Cells[y][x].Right { // вправо
		if m.dfs(x+1, y) {
			return true
		}
	}
	if x > 0 && !m.mazeInfo.Cells[y][x-1].Right { // влево
		if m.dfs(x-1, y) {
			return true
		}
	}
	if y < m.mazeInfo.Rows-1 && !m.mazeInfo.Cells[y][x].Bottom { // вниз
		if m.dfs(x, y+1) {
			return true
		}
	}
	if y > 0 && !m.mazeInfo.Cells[y-1][x].Bottom { // вверх
		if m.dfs(x, y-1) {
			return true
		}
	}

	m.path = m.path[:len(m.path)-1] // Удаляем текущую точку из пути, если путь не найден
	log.Printf("Возвращаемся из клетки: (%d, %d)\n", x, y)
	return false
}

// Метод для получения найденного пути
func (m *MazeSolving) GetPath() []Point {
	return m.path
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
	for _, point := range g.mazeSolving.GetPath() {
		vector.DrawFilledRect(screen, float32(point.X)*g.cellSize, float32(point.Y)*g.cellSize, g.cellSize, g.cellSize, color.RGBA{0, 0, 255, 128}, false) // Синий цвет для пути
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
