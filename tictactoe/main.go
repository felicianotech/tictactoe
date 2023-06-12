package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct{}

func (this *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Hey Michael")
}

func (this *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func (this *Game) Update() error {
	return nil
}

func main() {

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Michael Threat")

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
