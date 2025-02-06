package game

import (
	"go-maze/pkg/cave"
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	maxSize         = 50
	wallThickness   = 2
	sceneWidth      = 500
	sceneHeight     = 500
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
	if w > maxSize || h > maxSize {
		log.Fatalf("Размер лабиринта не должен превышать %d", maxSize)
	}
	ebiten.SetWindowSize(sceneWidth+int(borderThickness*2), sceneHeight+buttonHeight*3+int(borderThickness))
	ebiten.SetWindowTitle("Cave Generator")
	cellSize := float32(sceneWidth) / float32(w)
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

		if g.isInsideButton(float32(x), float32(y), float32(sceneHeight+borderThickness), buttonHeight) {
			go g.ShowFileSelector()
		}

		if g.isInsideButton(float32(x), float32(y), float32(sceneHeight+borderThickness+buttonHeight), buttonHeight) {
			g.Step()
			g.autoStepActive = false
		}

		if g.isInsideButton(float32(x), float32(y), float32(sceneHeight+borderThickness+buttonHeight*2), buttonHeight) {
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
func (g *Game) PrintCave() {
	g.cave.Print()
}
