package model

import (
	"math"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"camaretto/view"
)

type TextBox struct {
	width, height float64
	message string
	textColor color.RGBA
	length, count, speed int

	topMargin, bottomMargin float64
	leftMargin, rightMargin float64

	barScale float64
	barWidth, barHeight float64

	backgroundColor color.RGBA
	background *ebiten.Image

	actualImg *ebiten.Image

	SSprite *view.Sprite
}

func wrapText(message string, width int) string {
	var step int = width/int(view.FontSize*2/3)
	for i := step; i < len(message); i = i + step {
		message = message[:i] + "\n" + message[i:]
	}
	return message
}

func NewTextBox(w, h float64, msg string, tClr color.RGBA, bgClr color.RGBA) *TextBox {
	var tb *TextBox = &TextBox{}

	tb.width, tb.height = w, h
	tb.message = wrapText(msg, int(w))
	tb.textColor = tClr

	tb.length = 0
	tb.count = 0
	tb.speed = 5

	tb.barWidth, tb.barHeight = float64(view.BarWidth), float64(view.BarHeight)
	tb.barScale = 0.5
	tb.barHeight = tb.barHeight * tb.barScale

	tb.topMargin, tb.bottomMargin = tb.barHeight + 25, tb.barHeight + 25
	tb.leftMargin, tb.rightMargin = tb.barHeight + 30, tb.barHeight + 30

	tb.backgroundColor = bgClr
	tb.SSprite = nil
	tb.RenderBackground()

	return tb
}

func (tb *TextBox) RenderBackground() {
	var scaleWidth float64 = (tb.width - (tb.leftMargin+tb.rightMargin) - 50) / tb.barWidth
	var scaleHeight float64 = (tb.height - (tb.topMargin+tb.bottomMargin) - 50) / tb.barWidth

	// Init background
	tb.background = ebiten.NewImage(int(tb.width), int(tb.height))
	tb.background.Fill(tb.backgroundColor)

	// Top border
	var op *ebiten.DrawImageOptions = &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scaleWidth, tb.barScale)
	op.GeoM.Translate(-(tb.barWidth*scaleWidth)/2, -tb.barHeight/2)
	op.GeoM.Translate(tb.width/2, tb.topMargin/2)
	tb.background.DrawImage(view.BarImage, op)
	// Bottom border
	op.GeoM.Reset()
	op.GeoM.Scale(scaleWidth, tb.barScale)
	op.GeoM.Translate(-(tb.barWidth*scaleWidth)/2, -tb.barHeight/2)
	op.GeoM.Translate(tb.width/2, tb.height - (tb.bottomMargin/2))
	tb.background.DrawImage(view.BarImage, op)
	// Left border
	op.GeoM.Reset()
	op.GeoM.Scale(scaleHeight, tb.barScale)
	op.GeoM.Translate(-(tb.barWidth*scaleHeight)/2, -tb.barHeight/2)
	op.GeoM.Rotate(math.Pi/2)
	op.GeoM.Translate(tb.leftMargin/2, tb.height/2)
	tb.background.DrawImage(view.BarImage, op)
	// Right border
	op.GeoM.Reset()
	op.GeoM.Scale(scaleHeight, tb.barScale)
	op.GeoM.Translate(-(tb.barWidth*scaleHeight)/2, -tb.barHeight/2)
	op.GeoM.Rotate(math.Pi/2)
	op.GeoM.Translate(tb.width - (tb.rightMargin/2), tb.height/2)
	tb.background.DrawImage(view.BarImage, op)

	var iconScale float64 = 0.1
	var w, h float64 = 0, 0

	w, h = float64(view.CoffeeWidth)*iconScale, float64(view.CoffeeHeight)*iconScale
	// Top left coffee icon
	op.GeoM.Reset()
	op.GeoM.Scale(iconScale, iconScale)
	op.GeoM.Translate(-w/2, -h/2)
	op.GeoM.Translate(tb.leftMargin/2, tb.topMargin/2)
	tb.background.DrawImage(view.CoffeeImage, op)
	// Bottom right coffee icon
	op.GeoM.Reset()
	op.GeoM.Scale(iconScale, iconScale)
	op.GeoM.Translate(-w/2, -h/2)
	op.GeoM.Translate(tb.width - tb.rightMargin/2, tb.height - tb.bottomMargin/2)
	tb.background.DrawImage(view.CoffeeImage, op)

	w, h = float64(view.AmarettoWidth)*iconScale, float64(view.AmarettoHeight)*iconScale
	// Top right amaretto icon
	op.GeoM.Reset()
	op.GeoM.Scale(iconScale, iconScale)
	op.GeoM.Translate(-w/2, -h/2)
	op.GeoM.Translate(tb.width - tb.rightMargin/2, tb.topMargin/2)
	tb.background.DrawImage(view.AmarettoImage, op)
	// Bottom left amaretto icon
	op.GeoM.Reset()
	op.GeoM.Scale(iconScale, iconScale)
	op.GeoM.Translate(-w/2, -h/2)
	op.GeoM.Translate(tb.leftMargin/2, tb.height - tb.bottomMargin/2)
	tb.background.DrawImage(view.AmarettoImage, op)

	if tb.SSprite == nil {
		tb.SSprite = view.NewSprite(tb.background, true, tb.backgroundColor, nil)
	} else {
		tb.SSprite.SetImage(tb.background)
	}
}

func (tb *TextBox) Render() {
	if tb.length < len(tb.message) {
		tb.count++
		if tb.count == tb.speed {
			tb.length++
			tb.count = 0
		}
	}

	tb.actualImg = ebiten.NewImage(int(tb.width), int(tb.height))
	tb.actualImg.DrawImage(tb.background, nil)

	var op *text.DrawOptions = &text.DrawOptions{}
	op.ColorScale.ScaleWithColor(tb.textColor)
	op.LayoutOptions.LineSpacing = view.FontSize + 3
	op.GeoM.Reset()
	op.GeoM.Translate(tb.leftMargin*3/2, tb.topMargin)
	text.Draw(tb.actualImg, tb.message[:tb.length], &text.GoTextFace{Source: view.FaceSource, Size: view.FontSize}, op)

	tb.SSprite.SetImage(tb.actualImg)
}

func (tb *TextBox) SetMessage(msg string) {
	tb.message = wrapText(msg, int(tb.width))
	tb.length = 0
	tb.count = 0
}

func (tb *TextBox) SetTextColor(clr color.RGBA) {
	tb.textColor = clr
	tb.Render()
}

func (tb *TextBox) SetBackgroundColor(clr color.RGBA) {
	tb.backgroundColor = clr
	tb.RenderBackground()
}