package main

import (
	"go-maze/pkg/maze"
	"image/color"
	"log"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type LabGame struct {
	isMaze       bool
	gameInstance *maze.Game // Ссылка на экземпляр Game
	inputWidth   string
	inputHeight  string
	mazeStarted  bool // Флаг для проверки, создан ли лабиринт
}

func (g *LabGame) Update() error {
	if g.isMaze {
		// Если мы в лабиринте, обновляем экземпляр игры
		if g.gameInstance != nil {
			if err := g.gameInstance.Update(); err != nil {
				return err
			}
		}
	} else {
		// Если мы в меню, обрабатываем нажатия кнопок
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			// Обработка нажатий кнопок
			if g.isInsideButton(float32(x), float32(y), 0, 100) && !g.mazeStarted { // "Запустить лабиринт"
				g.startMaze()
				g.mazeStarted = true // Устанавливаем флаг, что лабиринт создан
			} else if g.isInsideButton(float32(x), float32(y), 0, 150) { // "Запустить пещеру"
				g.isMaze = false
				// g.startCave()
			}
		}
	}
	return nil
}

func (g *LabGame) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{255, 255, 255, 255}) // Заливка фона белым цветом

	if g.isMaze {
		if g.gameInstance != nil {
			g.gameInstance.Draw(screen) // Отрисовка лабиринта
		}
	} else {
		// Отрисовка кнопок
		g.drawButton(screen, "Запустить лабиринт", 100, color.RGBA{0, 0, 0, 255})
		g.drawButton(screen, "Запустить пещеру", 150, color.RGBA{200, 200, 200, 255})
	}
}

func (g *LabGame) drawButton(screen *ebiten.Image, buttonText string, buttonY float32, bcolor color.RGBA) {
	buttonWidth := float32(200)
	buttonHeight := float32(40)

	// Отрисовка кнопки
	vector.DrawFilledRect(screen, 0, buttonY, buttonWidth, buttonHeight, bcolor, false)

	// Отрисовка текста на кнопке
	textWidth := float32(len(buttonText) * 8)
	textHeight := float32(16)

	textX := (buttonWidth - textWidth) / 2
	textY := buttonY + (buttonHeight-textHeight)/2

	ebitenutil.DebugPrintAt(screen, buttonText, int(textX), int(textY))
}

func (g *LabGame) startMaze() {
	rows, _ := strconv.Atoi(g.inputHeight)
	cols, _ := strconv.Atoi(g.inputWidth)
	g.gameInstance = maze.NewGame(rows, cols) // Создаем новый экземпляр Game
	g.isMaze = true                           // Устанавливаем состояние в true, чтобы отрисовать лабиринт
}

func (g *LabGame) isInsideButton(x, y float32, buttonX, buttonY float32) bool {
	buttonWidth := float32(200)
	buttonHeight := float32(40)
	return x >= buttonX && x <= buttonX+buttonWidth && y >= buttonY && y <= buttonY+buttonHeight
}

// Добавьте метод Layout
func (g *LabGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func (g *LabGame) ResetGame() {
	g.isMaze = false
	g.mazeStarted = false // Сбрасываем флаг
	if g.gameInstance != nil {
		g.gameInstance.ResetGame() // Сбрасываем игру в экземпляре maze.Game
	}
}

func main() {
	game := &LabGame{
		inputWidth:  "10", // Установите значения по умолчанию или получите их от пользователя
		inputHeight: "10",
	}
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
