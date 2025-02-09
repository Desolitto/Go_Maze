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

func NewGame(rows, cols int) *Game {
	if rows > config.MaxSize || cols > config.MaxSize {
		log.Fatalf("Размер лабиринта не должен превышать %d", config.MaxSize)
	}
	// Пересчитываем размеры окна
	windowWidth := config.SceneWidth + int(config.BorderThickness*2)
	windowHeight := config.SceneHeight + config.ButtonHeight*2 + int(config.BorderThickness)*3

	ebiten.SetWindowSize(windowWidth, windowHeight)

	// ebiten.SetWindowSize(config.SceneWidth+int(config.BorderThickness*2), config.SceneHeight+config.ButtonHeight*6+int(config.BorderThickness))
	ebiten.SetWindowTitle("Maze")
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
	log.Printf("Window size: %d x %d", windowWidth, windowHeight)

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
		} else if g.isInsideButton(float32(x), float32(y), float32(config.SceneHeight+config.BorderThickness+config.ButtonHeight), config.ButtonHeight) {
			g.ResetGame() // Сброс игры
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
	g.drawButton(screen, "Reset", float32(config.SceneHeight+config.BorderThickness*2+config.ButtonHeight), strokeColor)
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

func (g *Game) drawButton(screen *ebiten.Image, buttonText string, buttonY float32, bcolor color.RGBA) {
	buttonWidth := float32(config.SceneWidth + config.BorderThickness*2)

	// Отрисовка кнопки
	vector.DrawFilledRect(screen, 0, buttonY, buttonWidth, config.ButtonHeight, bcolor, false)

	// Отрисовка текста на кнопке
	textWidth := float32(len(buttonText) * 8)
	textHeight := float32(16)

	textX := (buttonWidth - textWidth) / 2
	textY := buttonY + (config.ButtonHeight-textHeight)/2

	ebitenutil.DebugPrintAt(screen, buttonText, int(textX), int(textY))

	// Отрисовка границы под кнопкой
	borderY := buttonY + config.ButtonHeight
	borderHeight := float32(4)                  // Увеличенная высота границы
	borderColor := color.RGBA{192, 192, 192, 0} // Красный цвет границы

	vector.DrawFilledRect(screen, 0, borderY, buttonWidth, borderHeight, borderColor, false)
}

// Layout определяет размер окна
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func (g *Game) Maze() *Maze {
	return g.maze
}

func (g *Game) ResetGame() {
	g.mazeSolving.startX = -1
	g.mazeSolving.startY = -1
	g.mazeSolving.endX = -1
	g.mazeSolving.endY = -1
	g.mazeSolving.isStartSet = false
	g.mazeSolving.isEndSet = false
	g.mazeSolving.path = []Point{} // Очищаем путь
	// Сбрасываем матрицу решения
	for y := range g.mazeSolving.solvingMatrix {
		for x := range g.mazeSolving.solvingMatrix[y] {
			g.mazeSolving.solvingMatrix[y][x] = 0 // Сбрасываем посещенные клетки
		}
	}
	g.mazeSolving = NewMazeSolving(g.maze, -1, -1, -1, -1) // Обновляем решение
	log.Println("Игра сброшена")
}
