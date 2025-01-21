package game

import (
	"bufio"
	"go-maze/pkg/cave"
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
	"github.com/sqweek/dialog"
)

const (
	maxCaveSize     = 50
	wallThickness   = 2
	caveWidth       = 500
	caveHeight      = 500
	buttonHeight    = 30
	borderThickness = float32(2)
)

var colorAlive = color.RGBA{0, 0, 0, 255}
var colorDeath = color.RGBA{255, 255, 255, 255}

type Game struct {
	width, height    int
	cave             *cave.Cave
	cellSize         float32
	stepMode         bool
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
	cave := cave.NewCave(w, h)
	cave.GenerateCave(initialChance, birthLimit, deathLimit)
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
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()

		if g.isInsideButton(float32(x), float32(y), float32(caveHeight+borderThickness), buttonHeight) {
			go g.ShowFileSelector()
		}

		if g.isInsideButton(float32(x), float32(y), float32(caveHeight+borderThickness+buttonHeight), buttonHeight) {
			g.Step()
			g.autoStepActive = false
		}

		if g.isInsideButton(float32(x), float32(y), float32(caveHeight+borderThickness+buttonHeight*2), buttonHeight) {
			g.autoStepActive = !g.autoStepActive
		}
	}

	if g.autoStepActive {
		time.Sleep(g.autoStepInterval)
		g.Step()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(colorDeath)
	caveX := float32(0)
	caveY := float32(0)
	g.drawCaveBorder(screen)

	for y, row := range g.cave.Grid {
		for x, cell := range row {
			if cell == cave.Alive {
				vector.DrawFilledRect(screen, caveX+float32(x)*g.cellSize+2, caveY+float32(y)*g.cellSize+2, g.cellSize-wallThickness, g.cellSize-wallThickness, colorAlive, false)
			}
		}
	}

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
		g.cave = cave.NewCave(width, height)

		for y := 0; y < height; y++ {
			if scanner.Scan() {
				row := strings.Fields(scanner.Text())
				if len(row) != width {
					log.Fatal("Неверный формат файла: количество столбцов не совпадает с заданной шириной.")
				}
				for x, cell := range row {
					if cell == "0" {
						g.cave.Grid[y][x] = cave.Death
					} else if cell == "1" {
						g.cave.Grid[y][x] = cave.Alive
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
	newGrid := make([][]cave.Cell, g.cave.Height)
	for i := range newGrid {
		newGrid[i] = make([]cave.Cell, g.cave.Width)
	}
	for y := 0; y < g.cave.Height; y++ {
		for x := 0; x < g.cave.Width; x++ {
			wallCount := g.cave.CountAliveAround(x, y)
			if g.cave.Grid[y][x] == cave.Alive {
				if wallCount < g.deathLimit {
					newGrid[y][x] = cave.Death
				} else {
					newGrid[y][x] = cave.Alive
				}
			} else {
				if wallCount > g.birthLimit {
					newGrid[y][x] = cave.Alive
				} else {
					newGrid[y][x] = cave.Death
				}
			}
		}
	}

	g.cave.Grid = newGrid
}

func (g *Game) drawCaveBorder(screen *ebiten.Image) {
	borderColor := color.RGBA{255, 255, 255, 255}
	vector.StrokeLine(screen, 0, 0, caveWidth, 0, borderThickness, borderColor, true)
	vector.StrokeLine(screen, 0, caveHeight, caveWidth, caveHeight, borderThickness, borderColor, true)
	vector.StrokeLine(screen, 0, 0, 0, caveHeight, borderThickness, borderColor, true)
	vector.StrokeLine(screen, caveWidth, 0, caveWidth, caveHeight, borderThickness, borderColor, true)
}

func (g *Game) drawButton(screen *ebiten.Image, buttonText string, buttonY float32, color color.RGBA) {
	buttonWidth := float32(caveWidth + borderThickness*2)
	buttonHeight := float32(30)

	vector.DrawFilledRect(screen, 0, buttonY, buttonWidth, buttonHeight, color, false)

	textWidth := float32(len(buttonText) * 8)
	textHeight := float32(16)

	textX := (buttonWidth - textWidth) / 2
	textY := buttonY + (buttonHeight-textHeight)/2

	ebitenutil.DebugPrintAt(screen, buttonText, int(textX), int(textY))
}

func (g *Game) drawControlButtons(screen *ebiten.Image) {
	buttonY := float32(caveHeight + borderThickness + buttonHeight)

	g.drawButton(screen, "Generate Cave", float32(caveHeight+borderThickness), color.RGBA{0, 0, 155, 255})
	nextStepButtonY := buttonY
	g.drawButton(screen, "Next Step", nextStepButtonY, color.RGBA{0, 155, 0, 255})
	autoStepButtonY := nextStepButtonY + buttonHeight
	g.drawButton(screen, "Auto Step", autoStepButtonY, color.RGBA{155, 0, 0, 255})
}

func (g *Game) isInsideButton(x, y float32, buttonY float32, buttonHeight float32) bool {
	buttonX := float32(0)
	buttonWidth := float32(caveWidth + borderThickness*2)
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

	g.LoadCaveFromFile(filename)
}

func (g *Game) PrintCave() {
	g.cave.Print()
}
