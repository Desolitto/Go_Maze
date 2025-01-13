package game

import (
	"go-maze/pkg/maze"
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	maxMazeSize = 50 // Максимальный размер лабиринта
	// cellSize      = 15  // Размер ячейки (можно изменять)
	wallThickness   = 2          // Толщина стен
	mazeWidth       = 500        // Ширина области для лабиринта
	mazeHeight      = 500        // Высота области для лабиринта
	buttonHeight    = 30         // Высота кнопки
	borderThickness = float32(2) // Толщина рамки
)

var colorMaze = color.RGBA{255, 255, 255, 255}

type Game struct {
	w, h     int
	maze     *maze.Maze
	cellSize float32 // Размер ячейки
}

func NewGame(w, h int) *Game {
	if w > maxMazeSize || h > maxMazeSize {
		log.Fatalf("Размер лабиринта не должен превышать %d", maxMazeSize)
	}
	// Увеличиваем высоту окна на размер кнопки и рамки
	ebiten.SetWindowSize(mazeWidth+int(borderThickness*2), mazeHeight+buttonHeight+int(borderThickness))
	ebiten.SetWindowTitle("Maze")
	// Вычисляем размер ячейки
	cellSize := float32(mazeWidth) / float32(w)
	maze := maze.NewMaze(w, h)
	return &Game{w, h, maze, cellSize}
}

func (g *Game) Update() error {
	// Обработка ввода
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if g.isInsideButton(float32(x), float32(y)) {
			g.LoadMazeFromFile("path/to/maze/file.txt") // Укажите путь к файлу
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Рисуем фон
	screen.Fill(color.RGBA{0, 0, 0, 255})
	// Рисуем лабиринт
	for y, row := range g.maze.Grid {
		for x, cell := range row {
			if cell == maze.Wall {

				vector.DrawFilledRect(screen, float32(x)*g.cellSize+2, float32(y)*g.cellSize+2, g.cellSize-wallThickness, g.cellSize-wallThickness, colorMaze, false)
			}
		}
	}
	// Рисуем рамку для области лабиринта
	g.drawMazeBorder(screen)
	// Рисуем кнопку
	g.drawButton(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func (g *Game) LoadMazeFromFile(filename string) {
	// Загрузка лабиринта из файла
	//...
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	// Загрузка и проверка лабиринта
	// g.maze = maze.LoadMaze(file)
}

func (g *Game) drawMazeBorder(screen *ebiten.Image) {
	borderColor := color.RGBA{255, 255, 255, 255} // Цвет рамки (белый)

	// Рисуем верхнюю линию
	vector.StrokeLine(screen, 0, 0, mazeWidth+borderThickness*2, 0, borderThickness, borderColor, true)

	// Рисуем нижнюю линию
	vector.StrokeLine(screen, 0, mazeHeight+borderThickness, mazeWidth+borderThickness*2, mazeHeight+borderThickness, borderThickness, borderColor, true)

	// Рисуем левую линию
	vector.StrokeLine(screen, 0, 0, 0, mazeHeight+borderThickness, borderThickness, borderColor, true)

	// Рисуем правую линию
	vector.StrokeLine(screen, mazeWidth+borderThickness*2, 0, mazeWidth+borderThickness*2, mazeHeight+borderThickness, borderThickness, borderColor, true)
}

func (g *Game) drawButton(screen *ebiten.Image) {
	buttonWidth := float32(mazeWidth + borderThickness*2) // Кнопка на всю ширину
	buttonY := float32(mazeHeight + borderThickness)      // Позиция Y

	// Рисуем кнопку
	vector.DrawFilledRect(screen, 0, buttonY, buttonWidth, buttonHeight, color.RGBA{0, 0, 155, 255}, false)

	// Определяем текст и его размеры
	buttonText := "Load Maze"
	textWidth := float32(len(buttonText) * 8) // Оценка ширины текста
	textHeight := float32(16)                 // Высота текста

	// Вычисляем координаты для центрирования текста
	textX := (buttonWidth - textWidth) / 2
	textY := buttonY + (buttonHeight-textHeight)/2

	// Рисуем текст на кнопке
	ebitenutil.DebugPrintAt(screen, buttonText, int(textX), int(textY)) // Отрисовка текста
}

func (g *Game) isInsideButton(x, y float32) bool {
	buttonX := float32(0)                                 // Начинаем с нуля по X
	buttonY := float32(mazeHeight + borderThickness)      // Позиция Y
	buttonWidth := float32(mazeWidth + borderThickness*2) // Кнопка на всю ширину
	return x >= buttonX && x <= buttonX+buttonWidth && y >= buttonY && y <= buttonY+buttonHeight
}
