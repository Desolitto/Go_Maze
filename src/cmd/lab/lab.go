package main

import (
	"go-maze/config"
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
	cursorX      float32
}

func (g *LabGame) Update() error {
	if g.isMaze {
		if g.gameInstance != nil {
			if err := g.gameInstance.Update(); err != nil {
				return err
			}
		}
	} else {
		// Обработка нажатий кнопок
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			if g.isInsideButton(float32(x), float32(y), 10, 170, 200, 40) && !g.mazeStarted {
				g.startMaze()
				g.mazeStarted = true
			} else if g.isInsideButton(float32(x), float32(y), 10, 230, 200, 40) {
				g.isMaze = false
				// g.startCave()
			}
		}

		// Обработка ввода с клавиатуры для строк
		if ebiten.IsKeyPressed(ebiten.KeyBackspace) {
			if len(g.inputHeight) > 0 {
				g.inputHeight = g.inputHeight[:len(g.inputHeight)-1] // Удаляем последний символ
			}
		} else {
			inputChars := ebiten.AppendInputChars(nil)
			if len(inputChars) == 1 {
				g.inputHeight += string(inputChars[0])
			}
		}

		// Обработка ввода с клавиатуры для столбцов
		if ebiten.IsKeyPressed(ebiten.KeyBackspace) {
			if len(g.inputWidth) > 0 {
				g.inputWidth = g.inputWidth[:len(g.inputWidth)-1] // Удаляем последний символ
			}
		} else {
			inputChars := ebiten.AppendInputChars(nil)
			if len(inputChars) == 1 {
				g.inputWidth += string(inputChars[0])
			}
		}
	}
	return nil
}
func (g *LabGame) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 0}) // Заливка фона белым цветом

	if g.isMaze {
		if g.gameInstance != nil {
			g.gameInstance.Draw(screen) // Отрисовка лабиринта
		}
	} else {
		// Отрисовка полей для ввода
		g.drawInputField(screen, "Rows: "+g.inputHeight, 50)
		g.drawInputField(screen, "Cols: "+g.inputWidth, 100)

		// Отрисовка кнопок
		g.drawButton(screen, "Start MAZE", 170, color.RGBA{0, 0, 155, 255})
		g.drawButton(screen, "Start CAVE", 230, color.RGBA{0, 155, 0, 255})
	}
}

func (g *LabGame) startMaze() {
	rows, err := strconv.Atoi(g.inputHeight)
	if err != nil {
		rows = 50 // Установите значение по умолчанию, если ввод некорректен
	}
	cols, err := strconv.Atoi(g.inputWidth)
	if err != nil {
		cols = 50 // Установите значение по умолчанию, если ввод некорректен
	}
	g.gameInstance = maze.NewGame(rows, cols) // Создаем новый экземпляр Game
	g.isMaze = true                           // Устанавливаем состояние в true, чтобы отрисовать лабиринт
}

func (g *LabGame) isInsideButton(x, y, buttonX, buttonY, buttonWidth, buttonHeight float32) bool {
	return x >= buttonX && x <= buttonX+buttonWidth && y >= buttonY && y <= buttonY+buttonHeight
}

func (g *LabGame) drawInputField(screen *ebiten.Image, label string, fieldY float32) {
	fieldWidth := float32(config.SceneWidth - config.BorderThickness*10)
	fieldHeight := float32(40) // Увеличиваем высоту поля ввода

	// Отрисовка поля ввода
	vector.DrawFilledRect(screen, 10, fieldY, fieldWidth, fieldHeight, color.RGBA{0, 255, 255, 0}, false) // Белый фон с отступами
	ebitenutil.DebugPrintAt(screen, label, int(15), int(fieldY+10))
}

func (g *LabGame) drawButton(screen *ebiten.Image, buttonText string, buttonY float32, buttonColor color.RGBA) {
	buttonWidth := float32(config.SceneWidth - config.BorderThickness*10)
	buttonHeight := float32(50) // Увеличиваем высоту кнопки

	// Отрисовка кнопки
	vector.DrawFilledRect(screen, 10, buttonY, buttonWidth, buttonHeight, buttonColor, false) // Оставляем отступ 10 пикселей по бокам

	textWidth := float32(len(buttonText) * 8)
	textHeight := float32(16)

	textX := (buttonWidth - textWidth) / 2
	textY := buttonY + (buttonHeight-textHeight)/2

	ebitenutil.DebugPrintAt(screen, buttonText, int(textX), int(textY))
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
	ebiten.SetWindowSize(config.SceneWidth, config.SceneHeight)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}

}
