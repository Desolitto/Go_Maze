package game

import (
	"bufio"
	"fmt"
	"go-maze/pkg/maze"
	"image/color"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	maxCaveSize = 50 // Максимальный размер лабиринта
	// cellSize      = 15  // Размер ячейки (можно изменять)
	wallThickness   = 2          // Толщина стен
	caveWidth       = 500        // Ширина области для лабиринта
	caveHeight      = 500        // Высота области для лабиринта
	buttonHeight    = 30         // Высота кнопки
	borderThickness = float32(2) // Толщина рамки
)

var colorMaze = color.RGBA{255, 255, 255, 255}

type Game struct {
	width, height    int
	cave             *maze.Maze
	cellSize         float32 // Размер ячейки
	stepMode         bool    // Режим пошаговой отрисовки
	birthLimit       int
	deathLimit       int
	initialChance    int
	autoStepInterval time.Duration
	autoStepActive   bool
}

func NewGame(w, h, birthLimit, deathLimit, initialChance int) *Game {
	if w > maxCaveSize || h > maxCaveSize {
		log.Fatalf("Размер лабиринта не должен превышать %d", maxCaveSize)
	}
	// Увеличиваем высоту окна на размер кнопки и рамки
	ebiten.SetWindowSize(caveWidth+int(borderThickness*2), caveHeight+buttonHeight*4+int(borderThickness*2))
	ebiten.SetWindowTitle("Cave Generator")
	// Вычисляем размер ячейки
	cellSize := float32(caveWidth) / float32(w)
	cave := maze.NewMaze(w, h)
	return &Game{
		width:            w,
		height:           h,
		cave:             cave,
		cellSize:         cellSize,
		stepMode:         false,
		birthLimit:       birthLimit,
		deathLimit:       deathLimit,
		initialChance:    initialChance,
		autoStepInterval: 100 * time.Millisecond,
		autoStepActive:   false,
	}
}

func (g *Game) Update() error {
	// Обработка ввода
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if g.isInsideButton(float32(x), float32(y)) {
			g.LoadCaveFromFile("/Users/calamarp/Desktop/go/Go_Maze/src/example.txt") // Укажите путь к файлу
			// g.cave.GenerateCave(45, 4, 3)
		}
	}

	if g.autoStepActive {
		time.Sleep(g.autoStepInterval)
		g.Step()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Рисуем фон
	screen.Fill(color.RGBA{0, 0, 0, 255})
	// Рисуем лабиринт
	for y, row := range g.cave.Grid {
		for x, cell := range row {
			if cell == maze.Wall {

				vector.DrawFilledRect(screen, float32(x)*g.cellSize+2, float32(y)*g.cellSize+2, g.cellSize-wallThickness, g.cellSize-wallThickness, colorMaze, false)
			}
		}
	}
	// Рисуем рамку для области лабиринта
	g.drawCaveBorder(screen)
	// Рисуем кнопку
	g.drawButton(screen)
	g.drawControlButtons(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func (g *Game) LoadCaveFromFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Читаем размеры пещеры
	if scanner.Scan() {
		dimensions := strings.Fields(scanner.Text())
		if len(dimensions) != 2 {
			log.Fatal("Неверный формат файла: первая строка должна содержать размеры пещеры.")
		}

		width, err := strconv.Atoi(dimensions[0])
		if err != nil || width > maxCaveSize {
			log.Fatal("Неверная ширина пещеры.")
		}

		height, err := strconv.Atoi(dimensions[1])
		if err != nil || height > maxCaveSize {
			log.Fatal("Неверная высота пещеры.")
		}

		g.width, g.height = width, height
		g.cellSize = float32(caveWidth) / float32(width)
		g.cave = maze.NewMaze(width, height)

		// Читаем содержимое пещеры
		for y := 0; y < height; y++ {
			if scanner.Scan() {
				row := strings.Fields(scanner.Text())
				if len(row) != width {
					log.Fatal("Неверный формат файла: количество столбцов не совпадает с заданной шириной.")
				}
				for x, cell := range row {
					if cell == "0" {
						g.cave.Grid[y][x] = maze.Wall
					} else if cell == "1" {
						g.cave.Grid[y][x] = maze.Empty
					} else {
						log.Fatal("Неверный символ в пещере: должен быть 0 или 1.")
					}
				}
			}
		}
	}
}

func (g *Game) Step() {
	// Логика для одного шага генерации пещеры
	newGrid := make([][]maze.Cell, g.cave.Height)
	for i := range newGrid {
		newGrid[i] = make([]maze.Cell, g.cave.Width)
		copy(newGrid[i], g.cave.Grid[i]) // Копируем текущее состояние
	}

	for y := 0; y < g.cave.Height; y++ {
		for x := 0; x < g.cave.Width; x++ {
			wallCount := g.cave.CountAliveAround(x, y)
			fmt.Println("wall count: ", wallCount)
			if g.cave.Grid[y][x] == maze.Wall { // Если клетка мертвая
				if wallCount <= g.birthLimit {
					newGrid[y][x] = maze.Empty // Клетка становится живой
				}
			} else { // Если клетка живая
				if wallCount >= g.deathLimit {
					newGrid[y][x] = maze.Wall // Клетка умирает
				}
			}
		}
	}

	g.cave.Grid = newGrid // Обновляем состояние карты
}

func (g *Game) drawCaveBorder(screen *ebiten.Image) {
	borderColor := color.RGBA{255, 255, 255, 255} // Цвет рамки (белый)

	// Рисуем верхнюю линию
	vector.StrokeLine(screen, 0, 0, caveWidth+borderThickness*2, 0, borderThickness, borderColor, true)

	// Рисуем нижнюю линию
	vector.StrokeLine(screen, 0, caveHeight+borderThickness, caveWidth+borderThickness*2, caveHeight+borderThickness, borderThickness, borderColor, true)

	// Рисуем левую линию
	vector.StrokeLine(screen, 0, 0, 0, caveHeight+borderThickness, borderThickness, borderColor, true)

	// Рисуем правую линию
	vector.StrokeLine(screen, caveWidth+borderThickness*2, 0, caveWidth+borderThickness*2, caveHeight+borderThickness, borderThickness, borderColor, true)
}

func (g *Game) drawButton(screen *ebiten.Image) {
	buttonWidth := float32(caveWidth + borderThickness*2) // Кнопка на всю ширину
	buttonY := float32(caveHeight + borderThickness)      // Позиция Y

	// Рисуем кнопку
	vector.DrawFilledRect(screen, 0, buttonY, buttonWidth, buttonHeight, color.RGBA{0, 0, 155, 255}, false)

	// Определяем текст и его размеры
	buttonText := "Generate Cave"
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
	buttonY := float32(caveHeight + borderThickness)      // Позиция Y
	buttonWidth := float32(caveWidth + borderThickness*2) // Кнопка на всю ширину
	return x >= buttonX && x <= buttonX+buttonWidth && y >= buttonY && y <= buttonY+buttonHeight
}

func (g *Game) drawControlButtons(screen *ebiten.Image) {
	buttonWidth := float32(caveWidth + borderThickness*2)
	buttonY := float32(caveHeight + borderThickness + buttonHeight + 10) // Сдвигаем на высоту кнопки + отступ

	// Кнопка для следующего шага
	nextStepButtonY := buttonY
	vector.DrawFilledRect(screen, 0, nextStepButtonY, buttonWidth, buttonHeight, color.RGBA{0, 155, 0, 255}, false)
	ebitenutil.DebugPrintAt(screen, "Next Step", 10, int(nextStepButtonY)+5)

	// Кнопка для автоматического шага
	autoStepButtonY := buttonY + buttonHeight + 10 // Добавляем отступ между кнопками
	vector.DrawFilledRect(screen, 0, autoStepButtonY, buttonWidth, buttonHeight, color.RGBA{155, 0, 0, 255}, false)
	ebitenutil.DebugPrintAt(screen, "Auto Step", 10, int(autoStepButtonY)+5)

	// Обработка нажатий на кнопки
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if g.isInsideControlButton(float32(x), float32(y), nextStepButtonY) {
			g.Step() // Вызываем шаг
		}
		if g.isInsideControlButton(float32(x), float32(y), autoStepButtonY) {
			g.autoStepActive = !g.autoStepActive // Переключаем автошаг
		}
	}
}

func (g *Game) isInsideControlButton(x, y float32, buttonY float32) bool {
	buttonX := float32(0)
	buttonWidth := float32(caveWidth + borderThickness*2)
	return x >= buttonX && x <= buttonX+buttonWidth && y >= buttonY && y <= buttonY+buttonHeight
}
