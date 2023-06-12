package main

import (
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
)

var (
	uiFont           font.Face
	xOffset, yOffset int
	player2          bool
)

func init() {

	tt, err := opentype.Parse(goregular.TTF)
	if err != nil {
		log.Fatal(err)
	}

	uiFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    36,
		DPI:     96,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		log.Fatal(err)
	}
}

type tile struct {
	rect      image.Rectangle
	value     string
	mouseDown bool
	onPressed func(this *tile)
}

func (this *tile) Draw(screen *ebiten.Image) {

	m := uiFont.Metrics()
	w := font.MeasureString(uiFont, this.value).Floor()
	h := (m.Ascent + m.Descent).Floor()
	x := this.rect.Min.X + (this.rect.Dx()-w)/2
	y := this.rect.Min.Y + (this.rect.Dy()-h)/2 + m.Ascent.Floor()
	vector.DrawFilledRect(screen, float32(this.rect.Min.X), float32(this.rect.Min.Y), float32(this.rect.Dx()), float32(this.rect.Dy()), color.RGBA{0xff, 0x00, 0xff, 0xff}, true)
	text.Draw(screen, this.value, uiFont, x, y, color.Black)
}

func (this *tile) Update() {

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {

		x, y := ebiten.CursorPosition()
		x = x - xOffset
		y = y - yOffset

		if this.rect.Min.X <= x && x < this.rect.Max.X && this.rect.Min.Y <= y && y < this.rect.Max.Y {
			this.mouseDown = true
		} else {
			this.mouseDown = false
		}
	} else {

		if this.mouseDown && this.value == "" {
			if player2 {
				this.value = "o"
			} else {
				this.value = "x"
			}
			player2 = !player2
		}

		this.mouseDown = false
	}

}

type Board struct {
	tiles [3][3]*tile
}

func (this *Board) Draw(screen *ebiten.Image) {

	screen.Fill(color.RGBA{0x00, 0x00, 0x00, 0xff})

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			this.tiles[i][j].Draw(screen)
		}
	}
}

func (this *Board) Update() error {

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			this.tiles[i][j].Update()
		}
	}

	// check for winner
	if (this.tiles[0][0].value != "" && this.tiles[0][0].value == this.tiles[0][1].value && this.tiles[0][1].value == this.tiles[0][2].value) ||
		(this.tiles[1][0].value != "" && this.tiles[1][0].value == this.tiles[1][1].value && this.tiles[1][1].value == this.tiles[1][2].value) ||
		(this.tiles[2][0].value != "" && this.tiles[2][0].value == this.tiles[2][1].value && this.tiles[2][1].value == this.tiles[2][2].value) ||
		(this.tiles[0][0].value != "" && this.tiles[0][0].value == this.tiles[1][0].value && this.tiles[1][0].value == this.tiles[2][0].value) ||
		(this.tiles[0][1].value != "" && this.tiles[0][1].value == this.tiles[1][1].value && this.tiles[1][1].value == this.tiles[2][1].value) ||
		(this.tiles[0][2].value != "" && this.tiles[0][2].value == this.tiles[1][2].value && this.tiles[1][2].value == this.tiles[2][2].value) ||
		(this.tiles[0][0].value != "" && this.tiles[0][0].value == this.tiles[1][1].value && this.tiles[1][1].value == this.tiles[2][2].value) ||
		(this.tiles[2][2].value != "" && this.tiles[2][2].value == this.tiles[1][1].value && this.tiles[1][1].value == this.tiles[0][2].value) {
		println("You won!")
		return ebiten.Termination
	}

	var count int
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if this.tiles[i][j].value != "" {
				count++
			}
		}
	}

	if count == 9 {
		println("You tied!")
		return ebiten.Termination
	}

	return nil
}

func NewBoard() *Board {

	board := &Board{}

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {

			hGap := 20
			vGap := 20

			if i == 0 {
				hGap = 0
			}

			if j == 0 {
				vGap = 0
			}

			board.tiles[i][j] = &tile{
				rect: image.Rect(i*100+hGap, j*100+vGap, i*100+100, j*100+100),
			}
		}
	}

	return board
}

type Game struct {
	board      *Board
	boardImage *ebiten.Image
}

func (this *Game) Draw(screen *ebiten.Image) {

	if this.boardImage == nil {
		this.boardImage = ebiten.NewImage(300, 300)
	}

	screen.Fill(color.RGBA{0xff, 0x00, 0xff, 0xff})

	op := &ebiten.DrawImageOptions{}
	sw, sh := screen.Bounds().Dx(), screen.Bounds().Dy()
	bw, bh := this.boardImage.Bounds().Dx(), this.boardImage.Bounds().Dy()

	x := (sw - bw) / 2
	y := (sh - bh) / 2

	xOffset = x
	yOffset = y

	this.board.Draw(this.boardImage)

	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(this.boardImage, op)
}

func (this *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 700, 700
}

func (this *Game) Update() error {

	err := this.board.Update()
	if err != nil {
		return ebiten.Termination
	}

	return nil
}

func main() {

	board := NewBoard()

	game := &Game{
		board: board,
	}

	ebiten.SetWindowSize(700, 700)
	ebiten.SetWindowTitle("Michael Threat")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
