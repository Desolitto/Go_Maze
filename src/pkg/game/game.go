package game

import (
	"go-maze/pkg/maze"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	w, h int
	maze *maze.Maze
}

func NewGame(w, h int) *Game {
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Maze")
	maze := maze.NewMaze(w, h)
	return &Game{w, h, maze}
}

func (g *Game) Update() error {
	// Обновление логики игры
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Рисуем фон
	screen.Fill(color.RGBA{0, 0, 0, 255})
	// Рисуем лабиринт
	for y, row := range g.maze.Grid {
		for x, cell := range row {
			if cell == maze.Wall {
				// Рисуем стену (например, прямоугольник)
				ebitenutil.DrawRect(screen, float64(x*40), float64(y*40), 40, 40, color.RGBA{255, 0, 0, 255}) // Красная стена
			}
		}
	}

	ebitenutil.DebugPrint(screen, "Maze Game")
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}
