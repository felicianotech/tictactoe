package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Board struct {
	tiles [3][3]string
}

func (this *Board) Draw(screen *ebiten.Image) {

	var output string

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			output = fmt.Sprintf("%s %s | ", output, this.tiles[i][j])
		}

		output = fmt.Sprintf("%s\n", output)
	}

	ebitenutil.DebugPrint(screen, output)
}

type Game struct {
	board      *Board
	boardImage *ebiten.Image
}

func (this *Game) Draw(screen *ebiten.Image) {

	if this.boardImage == nil {
		this.boardImage = ebiten.NewImage(500, 500)
	}

	this.board.Draw(this.boardImage)
	screen.DrawImage(this.boardImage, nil)
}

func (this *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func (this *Game) Update() error {
	return nil
}

func main() {

	game := &Game{
		board: &Board{},
	}

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Michael Threat")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
