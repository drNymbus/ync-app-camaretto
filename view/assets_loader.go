package view

import (
	"log"

	"os"
	"io"
	"bytes"
	"image"
	"image/color"

	"golang.org/x/text/language"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var (
	Sheet *ebiten.Image
	TileWidth int = 64//32
	TileHeight int = 64//32

	CardImage [13]*ebiten.Image
	JokerImage *ebiten.Image
	EmptyCardImage *ebiten.Image
	HiddenCardImage *ebiten.Image
	GraveImage *ebiten.Image

	BlueButtonFlatImage *ebiten.Image
	BlueButtonDepthImage *ebiten.Image
	RedButtonFlatImage *ebiten.Image
	RedButtonDepthImage *ebiten.Image
	GreenButtonFlatImage *ebiten.Image
	GreenButtonDepthImage *ebiten.Image
	YellowButtonFlatImage *ebiten.Image
	YellowButtonDepthImage *ebiten.Image
	ButtonWidth int = 192
	ButtonHeight int = 64

	CursorImage *ebiten.Image
	CursorWidth int = 32
	CursorHeight int = 32

	BarImage *ebiten.Image
	BarWidth int = 96
	BarHeight int = 16

	CoffeeImage *ebiten.Image
	CoffeeWidth int = 325
	CoffeeHeight int = 297

	AmarettoImage *ebiten.Image
	AmarettoWidth int = 325
	AmarettoHeight int = 297

	FaceSource *text.GoTextFaceSource
	TextFace *text.GoTextFace
	FontSize float64 = 24
)

// @desc: Given a file path returns the content of the file as a byte slice
func getFileByte(filename string) []byte {
	var err error

	// Opens the file
	var file *os.File
	file, err = os.Open(filename)
	if err != nil { log.Fatal("[parameters.getFileByte] Open file:", err) }

	// Get file size, to make sure we read it correctly/entirely
	var stat os.FileInfo
	stat, err = file.Stat()
	if err != nil { log.Fatal("[parameters.getFileByte] Get file info:", err) }

	// Turn the file into a byte slice
	var n int
	var fileByte []byte = make([]byte, stat.Size())
	n, err = file.Read(fileByte)
	if err != nil && err != io.EOF { log.Fatal("[parameters.getFileByte] Read file:", err) }
	if n != int(stat.Size()) { log.Println("Warning | [parameters.getFile] File read partially:", err) }

	return fileByte
}

func GetImage(filename string) *ebiten.Image {
	var err error
	var img *ebiten.Image

	img, _, err = ebitenutil.NewImageFromFile(filename)
	if err != nil {
		var msg string = "[view.GetImage] Load image (" + filename + "):"
		log.Fatal(msg, err)
	}

	return img
}

// @desc: Sets the global variables of the view package (a.k.a all images used throughout the game)
func InitAssets() {
	var err error
	var ogSheet *ebiten.Image = GetImage("assets/cards/cardsLarge_tilemap_packed.png")

	// Scale down the original sheet
	var width, height int = ogSheet.Size()
	var xScale, yScale float64 = 1, 1

	width, height = int(float64(width) * xScale), int(float64(height) * yScale)
	TileWidth, TileHeight = int(float64(TileWidth) * xScale), int(float64(TileHeight) * yScale)

	Sheet = ebiten.NewImage(width, height)
	op := &ebiten.DrawImageOptions{}; op.GeoM.Scale(xScale, yScale)
	Sheet.DrawImage(ogSheet, op)

	for i := 0; i < 13; i++ { // Init all cards image from Ace to King
		var sx int = i * TileWidth
		var img *ebiten.Image = Sheet.SubImage(image.Rect(sx, 0, sx+TileWidth, TileHeight)).(*ebiten.Image)
		CardImage[i] = img
	}

	// All other cards are not logically placed in the tilemap sheet
	JokerImage = Sheet.SubImage(image.Rect((13*TileWidth), (2*TileHeight), (13*TileWidth) + TileWidth, (2*TileHeight) + TileHeight)).(*ebiten.Image)
	EmptyCardImage = Sheet.SubImage(image.Rect((13*TileWidth), 0, (13*TileWidth) + TileWidth, TileHeight)).(*ebiten.Image)
	HiddenCardImage = Sheet.SubImage(image.Rect((13*TileWidth), TileHeight, (13*TileWidth) + TileWidth, TileHeight + TileHeight)).(*ebiten.Image)

	// Death Sprite
	var tmp *ebiten.Image = GetImage("assets/characters/jesus.jpg")

	width, height = tmp.Size()
	xScale, yScale = 0.1, 0.1
	width, height = int(float64(width) * xScale), int(float64(height) * yScale)
	GraveImage = ebiten.NewImage(width, height)
	op.GeoM.Reset(); op.GeoM.Scale(xScale, yScale)
	GraveImage.DrawImage(tmp, op)

	//Buttons
	BlueButtonFlatImage = GetImage("assets/buttons/blue_button_pressed_gloss.png")
	BlueButtonDepthImage = GetImage("assets/buttons/blue_button.png")
	RedButtonFlatImage = GetImage("assets/buttons/red_button_pressed_gloss.png")
	RedButtonDepthImage = GetImage("assets/buttons/red_button.png")
	GreenButtonFlatImage = GetImage("assets/buttons/green_button_pressed_gloss.png")
	GreenButtonDepthImage = GetImage("assets/buttons/green_button.png")
	YellowButtonFlatImage = GetImage("assets/buttons/yellow_button_pressed_gloss.png")
	YellowButtonDepthImage = GetImage("assets/buttons/yellow_button.png")

	//Cursor
	CursorImage = GetImage("assets/cursor.png")
	//Bar
	BarImage = GetImage("assets/black_bar.png")
	//Coffee
	// CoffeeImage = ebiten.NewImage(CoffeeWidth, CoffeeHeight)
	// CoffeeImage.Fill(color.RGBA{255,255,255,255})
	// CoffeeImage.DrawImage(GetImage("assets/cafe.png"), nil)
	CoffeeImage = GetImage("assets/cafe.png")
	//Amaretto
	AmarettoImage = ebiten.NewImage(AmarettoWidth, AmarettoHeight)
	AmarettoImage.Fill(color.RGBA{255,255,255,255})
	AmarettoImage.DrawImage(GetImage("assets/amaretto.png"), nil)

	// Load font file
	var fontByte []byte = getFileByte("assets/fonts/NaturalMono_Regular.ttf")
	// var fontByte []byte = getFileByte("assets/fonts/Kenney_Future_Narrow.ttf")
	FaceSource, err = text.NewGoTextFaceSource(bytes.NewReader(fontByte))
	if err != nil { log.Fatal("[parametersInitAssets] Set FaceSource:", err) }

	TextFace = &text.GoTextFace{
		Source: FaceSource,
		Direction: text.DirectionLeftToRight,
		Size: 24, Language: language.English,
	}
}

func InitIcon(filepath string) (image.Image, error) {
	var err error
	var file *os.File

	file, err = os.Open(filepath)
	if err != nil {
		return nil, err
	}

	image, _, err := image.Decode(file)
	return image, err
}