package model

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"camaretto/view"
	// "camaretto/event"
)

type Button struct {
	message string
	textColor color.RGBA
	backgroundColor color.RGBA
	SSprite *view.Sprite
}

func NewButton(msg string, textClr color.RGBA, bgClr color.RGBA) *Button {
	var b *Button = &Button{msg, textClr, bgClr, nil}
	b.Render()
	return b
}

func (b *Button) Render() {
	var tWidth, tHeight float64 = text.Measure(b.message, view.TextFace, 0.0)
	var img *ebiten.Image = ebiten.NewImage(int(tWidth), int(tHeight))

	op := &text.DrawOptions{}; op.ColorScale.ScaleWithColor(b.textColor)
	text.Draw(img, b.message, &text.GoTextFace{Source: view.FaceSource, Size: view.FontSize}, op)

	b.SSprite = view.NewSprite(img, true, b.backgroundColor)
}

func (b *Button) SetMessage(msg string) {
	b.message = msg
	b.Render()
}

func (b *Button) SetTextColor(c color.RGBA) {
	b.textColor = c
	b.Render()
}

func (b *Button) SetBackgroundColor(c color.RGBA) {
	b.backgroundColor = c
	b.Render()
}

func (b *Button) Hover(x, y float64) {
	var color color.RGBA = b.backgroundColor
	if b.SSprite.In(x, y) {
		color.A = 255
		b.SetBackgroundColor(color)
	} else {
		color.A = 127
		b.SetBackgroundColor(color)
	}
}