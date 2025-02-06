package main

import (
	"flag"
	"fmt"
	"image/color"
	"log"
	"os"
	"time"

	"go-maze/config"
	"go-maze/pkg/maze"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sqweek/dialog"
	"golang.org/x/exp/rand"
)

// const (
// 	maxSize         = 50
// 	wallThickness   = 2
// 	SceneWidth      = 500
// 	SceneHeight     = 500 // Высота лабиринта
// 	ButtonHeight    = 30
// 	BorderThickness = float32(2)
// )

type Game struct {
	maze     *maze.Maze
	cellSize float32
}

func NewGame(rows, cols int) *Game {
	maze := maze.NewMaze(rows, cols)
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
	cellSize := float32(config.SceneWidth) / float32(cols)
	return &Game{maze: maze, cellSize: cellSize}
}

func (g *Game) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()

		if g.isInsideButton(float32(x), float32(y), float32(config.SceneHeight+config.BorderThickness), config.ButtonHeight) {
			go g.ShowFileSelector()
		}
	}
	return nil
}

func (g *Game) isInsideButton(x, y float32, buttonY float32, buttonHeight float32) bool {
	buttonX := float32(0)
	buttonWidth := float32(config.SceneHeight + config.BorderThickness*2)
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
	mazeNew, err := maze.LoadMaze(filename)
	if err != nil {
		log.Println("Ошибка при загрузке лабиринта:", err)
		return
	}

	// Обновляем состояние игры с новым лабиринтом
	g.maze = mazeNew
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
				vector.StrokeLine(screen, float32(x+1)*g.cellSize, float32(y)*g.cellSize, float32(x+1)*g.cellSize, float32(y+1)*g.cellSize, config.WallThickness, strokeColor, false)
			}

			// Рисуем нижнюю границу
			if y < g.maze.Rows-1 && g.maze.Cells[y][x].Bottom {
				vector.StrokeLine(screen, float32(x)*g.cellSize, float32(y+1)*g.cellSize, float32(x+1)*g.cellSize, float32(y+1)*g.cellSize, config.WallThickness, strokeColor, false)
			}
		}
	}
	g.drawButton(screen, "Open maze", float32(config.SceneHeight+config.BorderThickness), strokeColor)
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

func main() {
	w := flag.Int("w", config.MaxSize, "количество строк в лабиринте")
	h := flag.Int("h", config.MaxSize, "количество столбцов в лабиринте")
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
