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
	isMaze             bool
	isCave             bool
	gameInstance       *maze.Game // Ссылка на экземпляр Game
	inputWidth         string
	inputHeight        string
	inputBirthLimit    string
	inputDeathLimit    string
	inputInitialChance string
	mazeStarted        bool // Флаг для проверки, создан ли лабиринт
	cursorX            float32
	activeField        string // Новое поле для отслеживания активного поля ввода
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
				g.isCave = true
				g.activeField = "birthLimit" // Устанавливаем активное поле на "birthLimit"
			} else {
				g.handleCaveInputFieldClicks(float32(x), float32(y))
			}
		}

		// Обработка ввода с клавиатуры
		if ebiten.IsKeyPressed(ebiten.KeyBackspace) {
			g.handleKeyboardInput()
		} else {
			g.handleTextInput()
		}
	}
	return nil
}
func (g *LabGame) handleCaveInputFieldClicks(x, y float32) {
	if g.isInsideButton(float32(x), float32(y), 10, 50, 200, 40) {
		g.activeField = "height"
	} else if g.isInsideButton(float32(x), float32(y), 10, 100, 200, 40) {
		g.activeField = "width"
	} else if g.isInsideButton(float32(x), float32(y), 10, 290, 200, 40) {
		g.activeField = "birthLimit"
	} else if g.isInsideButton(float32(x), float32(y), 10, 360, 200, 40) {
		g.activeField = "deathLimit"
	} else if g.isInsideButton(float32(x), float32(y), 10, 430, 200, 40) {
		g.activeField = "initialChance"
	}
}
func (g *LabGame) handleTextInput() {
	inputChars := ebiten.AppendInputChars(nil)
	if len(inputChars) == 1 {
		if g.activeField == "height" {
			g.inputHeight += string(inputChars[0])
		} else if g.activeField == "width" {
			g.inputWidth += string(inputChars[0])
		} else if g.activeField == "birthLimit" {
			g.inputBirthLimit += string(inputChars[0])
		} else if g.activeField == "deathLimit" {
			g.inputDeathLimit += string(inputChars[0])
		} else if g.activeField == "initialChance" {
			g.inputInitialChance += string(inputChars[0])
		}
	}
}
func (g *LabGame) handleKeyboardInput() {
	if ebiten.IsKeyPressed(ebiten.KeyBackspace) {
		if g.activeField == "height" && len(g.inputHeight) > 0 {
			g.inputHeight = g.inputHeight[:len(g.inputHeight)-1] // Удаляем последний символ
		} else if g.activeField == "width" && len(g.inputWidth) > 0 {
			g.inputWidth = g.inputWidth[:len(g.inputWidth)-1] // Удаляем последний символ
		} else if g.activeField == "birthLimit" && len(g.inputBirthLimit) > 0 {
			g.inputBirthLimit = g.inputBirthLimit[:len(g.inputBirthLimit)-1] // Удаляем последний символ
		} else if g.activeField == "deathLimit" && len(g.inputDeathLimit) > 0 {
			g.inputDeathLimit = g.inputDeathLimit[:len(g.inputDeathLimit)-1] // Удаляем последний символ
		} else if g.activeField == "initialChance" && len(g.inputInitialChance) > 0 {
			g.inputInitialChance = g.inputInitialChance[:len(g.inputInitialChance)-1] // Удаляем последний символ
		}
	}
}
func (g *LabGame) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 0}) // Заливка фона белым цветом

	if g.isMaze {
		if g.gameInstance != nil {
			g.gameInstance.Draw(screen) // Отрисовка лабиринта
		}
	} else {
		// Отрисовка полей для ввода
		g.drawInputField(screen, "Rows: "+g.inputHeight, 50, g.activeField == "height")
		g.drawInputField(screen, "Cols: "+g.inputWidth, 100, g.activeField == "width")

		// Отрисовка кнопок
		g.drawButton(screen, "Start MAZE", 170, color.RGBA{0, 0, 155, 255})
		g.drawButton(screen, "Start CAVE", 230, color.RGBA{0, 155, 0, 255})
		if g.isCave {
			g.drawInputField(screen, "Birth Limit: "+g.inputBirthLimit, 290, g.activeField == "birthLimit")
			g.drawInputField(screen, "Death Limit: "+g.inputDeathLimit, 360, g.activeField == "deathLimit")
			g.drawInputField(screen, "Initial Chance: "+g.inputInitialChance, 430, g.activeField == "initialChance")
			g.drawButton(screen, "Start CAVE", 500, color.RGBA{0, 155, 0, 255})
		}
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

func (g *LabGame) drawInputField(screen *ebiten.Image, label string, fieldY float32, isActive bool) {
	fieldWidth := float32(config.SceneWidth - config.BorderThickness*10)
	fieldHeight := float32(40) // Увеличиваем высоту поля ввода

	// Отрисовка поля ввода
	if isActive {
		vector.DrawFilledRect(screen, 10, fieldY, fieldWidth, fieldHeight, color.RGBA{255, 255, 0, 255}, false) // Желтый фон с отступами
	} else {
		vector.DrawFilledRect(screen, 10, fieldY, fieldWidth, fieldHeight, color.RGBA{0, 255, 255, 0}, false) // Белый фон с отступами
	}
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
