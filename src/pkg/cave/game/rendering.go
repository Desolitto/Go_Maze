package game_cave

import (
	"go-maze/config"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func (g *Game) drawCaveBorder(screen *ebiten.Image) {
	borderColor := color.RGBA{255, 255, 255, 255}
	vector.StrokeLine(screen, 0, 0, config.SceneWidth, 0, config.BorderThickness, borderColor, true)
	vector.StrokeLine(screen, 0, config.SceneHeight, config.SceneWidth, config.SceneHeight, config.BorderThickness, borderColor, true)
	vector.StrokeLine(screen, 0, 0, 0, config.SceneHeight, config.BorderThickness, borderColor, true)
	vector.StrokeLine(screen, config.SceneWidth, 0, config.SceneWidth, config.SceneHeight, config.BorderThickness, borderColor, true)
}

func (g *Game) drawButton(screen *ebiten.Image, buttonText string, buttonY float32, color color.RGBA) {
	buttonWidth := float32(config.SceneWidth + config.BorderThickness*2)
	buttonHeight := float32(30)

	vector.DrawFilledRect(screen, 0, buttonY, buttonWidth, buttonHeight, color, false)

	textWidth := float32(len(buttonText) * 8)
	textHeight := float32(16)

	textX := (buttonWidth - textWidth) / 2
	textY := buttonY + (buttonHeight-textHeight)/2

	ebitenutil.DebugPrintAt(screen, buttonText, int(textX), int(textY))
}

func (g *Game) drawControlButtons(screen *ebiten.Image) {
	buttonY := float32(config.SceneHeight + config.BorderThickness + config.ButtonHeight)

	g.drawButton(screen, "Generate Cave", float32(config.SceneHeight+config.BorderThickness), color.RGBA{0, 0, 155, 255})
	nextStepButtonY := buttonY
	g.drawButton(screen, "Next Step", nextStepButtonY, color.RGBA{0, 155, 0, 255})
	autoStepButtonY := nextStepButtonY + config.ButtonHeight
	g.drawButton(screen, "Auto Step", autoStepButtonY, color.RGBA{155, 0, 0, 255})
}
