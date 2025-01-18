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
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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

var colorAlive = color.RGBA{0, 0, 0, 255}
var colorDeath = color.RGBA{255, 255, 255, 255}

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
	ebiten.SetWindowSize(caveWidth+int(borderThickness*2), caveHeight+buttonHeight*3+int(borderThickness))
	ebiten.SetWindowTitle("Cave Generator")
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

// func (g *Game) Update() error {
// 	// Обработка ввода
// 	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
// 		x, y := ebiten.CursorPosition()

// 		// Проверка нажатия кнопки "Загрузить лабиринт"
// 		if g.isInsideButton(float32(x), float32(y)) {
// 			g.LoadCaveFromFile("/Users/calamarp/Desktop/go/Go_Maze/src/example.txt")
// 			fmt.Println("Обновленная матрица:1")
// 			g.PrintMaze()
// 		}

// 		// Проверка нажатия кнопки "Следующий шаг"
// 		if g.isInsideControlButton(float32(x), float32(y), float32(caveHeight+borderThickness+buttonHeight)) {
// 			g.Step()                 // Вызываем шаг
// 			g.autoStepActive = false // Отключаем автошаг после шага
// 		}

// 		// Проверка нажатия кнопки "Автошаг"
// 		if g.isInsideControlButton(float32(x), float32(y), float32(caveHeight+borderThickness+buttonHeight*2)) {
// 			g.autoStepActive = !g.autoStepActive // Переключаем автошаг
// 		}
// 	}

// 	// Выполняем автоматический шаг, если активен
// 	if g.autoStepActive {
// 		time.Sleep(g.autoStepInterval)
// 		g.Step()
// 	}

// 	return nil
// }

func (g *Game) Update() error {
	// Обработка ввода
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()

		// Проверка нажатия кнопки "Загрузить лабиринт"
		if g.isInsideButton(float32(x), float32(y), float32(caveHeight+borderThickness), buttonHeight) {
			g.LoadCaveFromFile("/Users/calamarp/Desktop/go/Go_Maze/src/example.txt")
			fmt.Println("Обновленная матрица:1")
			g.PrintMaze()
		}

		// Проверка нажатия кнопки "Следующий шаг"
		if g.isInsideButton(float32(x), float32(y), float32(caveHeight+borderThickness+buttonHeight), buttonHeight) {
			g.Step()                 // Вызываем шаг
			g.autoStepActive = false // Отключаем автошаг после шага
		}

		// Проверка нажатия кнопки "Автошаг"
		if g.isInsideButton(float32(x), float32(y), float32(caveHeight+borderThickness+buttonHeight*2), buttonHeight) {
			g.autoStepActive = !g.autoStepActive // Переключаем автошаг
		}
	}

	// Выполняем автоматический шаг, если активен
	if g.autoStepActive {
		time.Sleep(g.autoStepInterval)
		g.Step()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Рисуем фон приложения (весь экран)
	screen.Fill(colorDeath)

	// Определяем размеры и координаты области для пещеры
	caveX := float32(0) // Начальная позиция по X
	caveY := float32(0) // Начальная позиция по Y

	// Рисуем рамку для области лабиринта
	g.drawCaveBorder(screen)

	// Рисуем лабиринт в области
	for y, row := range g.cave.Grid {
		for x, cell := range row {
			if cell == maze.Alive {
				vector.DrawFilledRect(screen, caveX+float32(x)*g.cellSize+2, caveY+float32(y)*g.cellSize+2, g.cellSize-wallThickness, g.cellSize-wallThickness, colorAlive, false)
			}
		}
	}

	// Рисуем кнопки
	// g.drawButton(screen)
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
						g.cave.Grid[y][x] = maze.Death
					} else if cell == "1" {
						g.cave.Grid[y][x] = maze.Alive
					} else {
						log.Fatal("Неверный символ в пещере: должен быть 0 или 1.")
					}
				}
			}
		}
	}
	g.autoStepActive = false
}

func (g *Game) Step() {
	newGrid := make([][]maze.Cell, g.cave.Height)
	for i := range newGrid {
		newGrid[i] = make([]maze.Cell, g.cave.Width)
	}

	fmt.Println("Состояние матрицы перед шагом:")
	g.PrintMaze()

	for y := 0; y < g.cave.Height; y++ {
		for x := 0; x < g.cave.Width; x++ {
			wallCount := g.cave.CountAliveAround(x, y)
			fmt.Printf("Клетка (%d, %d), wallCount: %d\n", x, y, wallCount)

			if g.cave.Grid[y][x] == maze.Alive { // Если клетка живая
				if wallCount < g.deathLimit {
					newGrid[y][x] = maze.Death // Клетка умирает
				} else {
					newGrid[y][x] = maze.Alive // Остается живой
				}
			} else { // Если клетка мертвая
				if wallCount > g.birthLimit {
					newGrid[y][x] = maze.Alive // Клетка становится живой
				} else {
					newGrid[y][x] = maze.Death // Остается мертвой
				}
			}
		}
	}

	g.cave.Grid = newGrid
	fmt.Println("Состояние матрицы после шага:")
	g.PrintMaze()
}

func (g *Game) drawCaveBorder(screen *ebiten.Image) {
	borderColor := color.RGBA{255, 255, 255, 255} // Цвет рамки (белый)

	// Рисуем рамку для области пещеры
	vector.StrokeLine(screen, 0, 0, caveWidth, 0, borderThickness, borderColor, true)
	vector.StrokeLine(screen, 0, caveHeight, caveWidth, caveHeight, borderThickness, borderColor, true)
	vector.StrokeLine(screen, 0, 0, 0, caveHeight, borderThickness, borderColor, true)
	vector.StrokeLine(screen, caveWidth, 0, caveWidth, caveHeight, borderThickness, borderColor, true)
}

// func (g *Game) drawButton(screen *ebiten.Image) {
// 	buttonWidth := float32(caveWidth + borderThickness*2) // Кнопка на всю ширину
// 	buttonY := float32(caveHeight + borderThickness)      // Позиция Y

// 	// Рисуем кнопку
// 	vector.DrawFilledRect(screen, 0, buttonY, buttonWidth, buttonHeight, color.RGBA{0, 0, 155, 255}, false)

// 	// Определяем текст и его размеры
// 	buttonText := "Generate Cave"
// 	textWidth := float32(len(buttonText) * 8) // Оценка ширины текста
// 	textHeight := float32(16)                 // Высота текста

// 	// Вычисляем координаты для центрирования текста
// 	textX := (buttonWidth - textWidth) / 2
// 	textY := buttonY + (buttonHeight-textHeight)/2

// 	// Рисуем текст на кнопке
// 	ebitenutil.DebugPrintAt(screen, buttonText, int(textX), int(textY)) // Отрисовка текста
// }

// func (g *Game) isInsideButton(x, y float32) bool {
// 	buttonX := float32(0)                                 // Начинаем с нуля по X
// 	buttonY := float32(caveHeight + borderThickness)      // Позиция Y
// 	buttonWidth := float32(caveWidth + borderThickness*2) // Кнопка на всю ширину
// 	return x >= buttonX && x <= buttonX+buttonWidth && y >= buttonY && y <= buttonY+buttonHeight
// }

// func (g *Game) drawControlButtons(screen *ebiten.Image) {
// 	buttonWidth := float32(caveWidth + borderThickness*2)
// 	buttonY := float32(caveHeight + borderThickness + buttonHeight) // Начальная позиция Y для первой кнопки

// 	// Кнопка для следующего шага
// 	nextStepButtonY := buttonY
// 	vector.DrawFilledRect(screen, 0, nextStepButtonY, buttonWidth, buttonHeight, color.RGBA{0, 155, 0, 255}, false)
// 	ebitenutil.DebugPrintAt(screen, "Next Step", 10, int(nextStepButtonY)+5)

// 	// Кнопка для автоматического шага
// 	autoStepButtonY := nextStepButtonY + buttonHeight // Кнопка сразу под первой
// 	vector.DrawFilledRect(screen, 0, autoStepButtonY, buttonWidth, buttonHeight, color.RGBA{155, 0, 0, 255}, false)
// 	ebitenutil.DebugPrintAt(screen, "Auto Step", 10, int(autoStepButtonY)+5)
// }

// func (g *Game) isInsideControlButton(x, y float32, buttonY float32) bool {
// 	buttonX := float32(0)
// 	buttonWidth := float32(caveWidth + borderThickness*2)
// 	return x >= buttonX && x <= buttonX+buttonWidth && y >= buttonY && y <= buttonY+buttonHeight
// }

func (g *Game) PrintMaze() {
	g.cave.Print() // Вызываем метод Print на структуре Maze
}

func (g *Game) drawButton(screen *ebiten.Image, buttonText string, buttonY float32, color color.RGBA) {
	buttonWidth := float32(caveWidth + borderThickness*2) // Кнопка на всю ширину
	buttonHeight := float32(30)                           // Высота кнопки (установите нужное значение)

	// Рисуем кнопку
	vector.DrawFilledRect(screen, 0, buttonY, buttonWidth, buttonHeight, color, false)

	// Определяем размеры текста
	textWidth := float32(len(buttonText) * 8) // Оценка ширины текста
	textHeight := float32(16)                 // Высота текста

	// Вычисляем координаты для центрирования текста
	textX := (buttonWidth - textWidth) / 2
	textY := buttonY + (buttonHeight-textHeight)/2

	// Рисуем текст на кнопке
	ebitenutil.DebugPrintAt(screen, buttonText, int(textX), int(textY)) // Отрисовка текста
}

func (g *Game) drawControlButtons(screen *ebiten.Image) {
	buttonY := float32(caveHeight + borderThickness + buttonHeight) // Начальная позиция Y для первой кнопки

	// Кнопка для генерации лабиринта
	g.drawButton(screen, "Generate Cave", float32(caveHeight+borderThickness), color.RGBA{0, 0, 155, 255})

	// Кнопка для следующего шага
	nextStepButtonY := buttonY
	g.drawButton(screen, "Next Step", nextStepButtonY, color.RGBA{0, 155, 0, 255})

	// Кнопка для автоматического шага
	autoStepButtonY := nextStepButtonY + buttonHeight // Кнопка сразу под следующей
	g.drawButton(screen, "Auto Step", autoStepButtonY, color.RGBA{155, 0, 0, 255})
}

func (g *Game) isInsideButton(x, y float32, buttonY float32, buttonHeight float32) bool {
	buttonX := float32(0)
	buttonWidth := float32(caveWidth + borderThickness*2)
	return x >= buttonX && x <= buttonX+buttonWidth && y >= buttonY && y <= buttonY+buttonHeight
}
