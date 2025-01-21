package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

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
