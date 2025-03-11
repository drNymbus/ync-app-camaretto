package model

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"camaretto/view"
)

type TextBox struct {
	width, height int
	message string
	textColor color.RGBA
	count int

	backgroundColor color.RGBA
	background *ebiten.Image

	actualImg *ebiten.Image

	SSprite *view.Sprite
}

func NewTextBox(w, h int, msg string, tClr color.RGBA, bgClr color.RGBA) *TextBox {
	var tb *TextBox = &TextBox{}

	tb.width, tb.height = w, h
	tb.message = msg
	tb.textColor = tClr
	tb.count = 0

	tb.background = ebiten.NewImage(w, h)
	tb.background.Fill(bgClr)

	return tb
}

func (tb *TextBox) Render() {
	tb.actualImg = ebiten.NewImage(tb.width, tb.height)
	tb.actualImg.DrawImage(tb.background, nil)

	var op *text.DrawOptions = &text.DrawOptions{}
	op.GeoM.Reset()
	op.GeoM.Translate(float64(w)/2 - (tw/2), float64(h)/2 - (th/2))
	text.Draw(tb.actualImg, tb.message[:tb.count], &text.GoTextFace{Source: view.FaceSource, Size: view.FontSize}, op)

	tb.SSprite.SetImage(tb.actualImg, nil)
}

func (tb *TextBox) SetMessage(msg string) {}
func (tb *TextBox) SetTextColor(clr color.RGBA) {}
func (tb *TextBox) SetBackgroundColor(clr color.RGBA) {}