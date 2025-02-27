package game_cave

import (
	"go-maze/config"
	"go-maze/pkg/cave"
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
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
	if w > config.MaxSize || h > config.MaxSize {
		log.Fatalf("Размер лабиринта не должен превышать %d", config.MaxSize)
	}
	ebiten.SetWindowSize(config.SceneWidth+int(config.BorderThickness*2), config.SceneHeight+config.ButtonHeight*3+int(config.BorderThickness))
	ebiten.SetWindowTitle("Cave Generator")
	cellSize := float32(config.SceneWidth) / float32(w)
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

		if g.isInsideButton(float32(x), float32(y), float32(config.SceneHeight+config.BorderThickness), config.ButtonHeight) {
			go g.ShowFileSelector()
		}

		if g.isInsideButton(float32(x), float32(y), float32(config.SceneHeight+config.BorderThickness+config.ButtonHeight), config.ButtonHeight) {
			g.Step()
			g.autoStepActive = false
		}

		if g.isInsideButton(float32(x), float32(y), float32(config.SceneHeight+config.BorderThickness+config.ButtonHeight*2), config.ButtonHeight) {
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
				vector.DrawFilledRect(screen, caveX+float32(x)*g.cellSize+2, caveY+float32(y)*g.cellSize+2, g.cellSize-config.WallThickness, g.cellSize-config.WallThickness, colorAlive, false)
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

// game_cave.go
func (g *Game) ResetGame() {
	// Реализация метода ResetGame() для cave.Game
	g.cave.GenerateCave(g.initialChance, g.birthLimit, g.deathLimit)
}

func (g *Game) PrintCave() {
	g.cave.Print()
}
