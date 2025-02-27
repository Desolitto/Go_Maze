package game

import "github.com/hajimehoshi/ebiten/v2"

type GameInterface interface {
	Update() error
	Draw(screen *ebiten.Image)
	Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int)
	ResetGame()
}
