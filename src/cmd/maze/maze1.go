package main

import (
	"flag"
	"fmt"
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	maxMazeSize     = 50
	wallThickness   = 2
	mazeWidth       = 500
	mazeHeight      = 500 // Высота лабиринта
	buttonHeight    = 30
	borderThickness = float32(2)
)

type Cell struct {
	Visited bool
	Top     bool
	Right   bool
	Bottom  bool
	Left    bool
}

type Maze struct {
	Rows, Cols int
	Cells      [][]Cell
}

func NewMaze(rows, cols int) *Maze {
	ebiten.SetWindowSize(mazeWidth+int(borderThickness*2), mazeWidth+buttonHeight+int(borderThickness))
	cells := make([][]Cell, rows)
	for i := range cells {
		cells[i] = make([]Cell, cols)
		for j := range cells[i] {
			cells[i][j] = Cell{Visited: false, Top: true, Right: true, Bottom: true, Left: true}
		}
	}
	return &Maze{Rows: rows, Cols: cols, Cells: cells}
}

func (m *Maze) Generate(x, y int) {
	m.Cells[y][x].Visited = true

	directions := []struct {
		dx, dy int
	}{
		{0, -1}, // вверх
		{1, 0},  // вправо
		{0, 1},  // вниз
		{-1, 0}, // влево
	}

	rand.Shuffle(len(directions), func(i, j int) {
		directions[i], directions[j] = directions[j], directions[i]
	})

	for _, dir := range directions {
		newX, newY := x+dir.dx, y+dir.dy
		if newX >= 0 && newX < m.Cols && newY >= 0 && newY < m.Rows && !m.Cells[newY][newX].Visited {
			if dir.dy == -1 { // вверх
				m.Cells[y][x].Top = false
				m.Cells[newY][newX].Bottom = false
			} else if dir.dy == 1 { // вниз
				m.Cells[y][x].Bottom = false
				m.Cells[newY][newX].Top = false
			} else if dir.dx == 1 { // вправо
				m.Cells[y][x].Right = false
				m.Cells[newY][newX].Left = false
			} else if dir.dx == -1 { // влево
				m.Cells[y][x].Left = false
				m.Cells[newY][newX].Right = false
			}
			m.Generate(newX, newY)
		}
	}
}

// Game структура для ebiten
type Game struct {
	maze     *Maze
	cellSize float32
}

// Новый экземпляр игры
func NewGame(rows, cols int) *Game {
	maze := NewMaze(rows, cols)
	maze.Generate(0, 0)                            // Начинаем генерацию с верхнего левого угла
	cellSize := float32(mazeWidth) / float32(cols) // Размер ячейки
	return &Game{maze: maze, cellSize: cellSize}
}

// Update обновляет состояние игры
func (g *Game) Update() error {
	return nil
}

// Draw отрисовывает лабиринт и кнопку
func (g *Game) Draw(screen *ebiten.Image) {
	strokeColor := color.RGBA{0, 0, 0, 255}     // Черный цвет
	fillColor := color.RGBA{255, 255, 255, 255} // Белый цвет

	// Рисуем лабиринт
	for y := 0; y < g.maze.Rows; y++ {
		for x := 0; x < g.maze.Cols; x++ {
			// Рисуем ячейку
			if g.maze.Cells[y][x].Visited {
				vector.DrawFilledRect(screen, float32(x)*g.cellSize, float32(y)*g.cellSize, g.cellSize, g.cellSize, fillColor, false)
			}
			// Рисуем стены
			if g.maze.Cells[y][x].Top {
				vector.StrokeLine(screen, float32(x)*g.cellSize, float32(y)*g.cellSize, float32(x+1)*g.cellSize, float32(y)*g.cellSize, wallThickness, strokeColor, false)
			}
			if g.maze.Cells[y][x].Right {
				vector.StrokeLine(screen, float32(x+1)*g.cellSize, float32(y)*g.cellSize, float32(x+1)*g.cellSize, float32(y+1)*g.cellSize, wallThickness, strokeColor, false)
			}
			if g.maze.Cells[y][x].Bottom {
				vector.StrokeLine(screen, float32(x)*g.cellSize, float32(y+1)*g.cellSize, float32(x+1)*g.cellSize, float32(y+1)*g.cellSize, wallThickness, strokeColor, false)
			}
			if g.maze.Cells[y][x].Left {
				vector.StrokeLine(screen, float32(x)*g.cellSize, float32(y)*g.cellSize, float32(x)*g.cellSize, float32(y+1)*g.cellSize, wallThickness, strokeColor, false)
			}
		}
	}

	// Рисуем кнопку под лабиринтом
	buttonY := mazeHeight
	buttonWidth := mazeWidth
	vector.DrawFilledRect(screen, 0, float32(buttonY), float32(buttonWidth), float32(buttonHeight), color.RGBA{200, 200, 200, 255}, false) // Цвет кнопки
	vector.StrokeLine(screen, 0, float32(buttonY), float32(buttonWidth), float32(buttonY), wallThickness, strokeColor, false)              // Контур кнопки
}

// Layout определяет размер окна
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return mazeWidth, mazeHeight + buttonHeight // Общая высота с кнопкой
}

func main() {
	w := flag.Int("w", maxMazeSize, "количество строк в лабиринте")
	h := flag.Int("h", maxMazeSize, "количество столбцов в лабиринте")
	flag.Parse()

	game := NewGame(*w, *h)
	fmt.Println("Исходная матрица:")
	for _, row := range game.maze.Cells {
		fmt.Println(row)
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
