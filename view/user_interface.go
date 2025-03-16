package view

import (
	"log"

	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ButtonWidth int = 192
	ButtonHeight int = 64

	CoffeeWidth int = 325
	CoffeeHeight int = 297
	AmarettoWidth int = 325
	AmarettoHeight int = 297
	BarWidth int = 96
	BarHeight int = 16

	CursorWidth int = 32
	CursorHeight int = 32
)

type ButtonImage struct {
	Pressed *ebiten.Image
	Released *ebiten.Image
}

func LoadButtonImage(c string) *ButtonImage {
	var bi *ButtonImage = &ButtonImage{nil, nil}

	if c == "RED" {
		bi.Pressed = GetImage("assets/buttons/red_button_pressed_gloss.png")
		bi.Released = GetImage("assets/buttons/red_button.png")
	} else if c == "BLUE" {
		bi.Pressed = GetImage("assets/buttons/blue_button_pressed_gloss.png")
		bi.Released = GetImage("assets/buttons/blue_button.png")
	} else if c == "GREEN" {
		bi.Pressed = GetImage("assets/buttons/green_button_pressed_gloss.png")
		bi.Released = GetImage("assets/buttons/green_button.png")	
	} else if c == "YELLOW" {
		bi.Pressed = GetImage("assets/buttons/yellow_button_pressed_gloss.png")
		bi.Released = GetImage("assets/buttons/yellow_button.png")
	} else { log.Fatal("[view.LoadButtonImage] Cannot load unknown button color:", c) }

	return bi
}

type IconImage struct {
	Coffee *ebiten.Image
	Amaretto *ebiten.Image
	Bar *ebiten.Image
}

func LoadIconImage() *IconImage {
	var ii *IconImage = &IconImage{}
	ii.Coffee = GetImage("assets/cafe.png")

	ii.Amaretto = ebiten.NewImage(AmarettoWidth, AmarettoHeight)
	ii.Amaretto.Fill(color.RGBA{255,255,255,255})
	ii.Amaretto.DrawImage(GetImage("assets/amaretto.png"), nil)

	ii.Bar = GetImage("assets/black_bar.png")
	return ii
}

func LoadCursorImage() *ebiten.Image {
	return GetImage("assets/cursor.png")
}