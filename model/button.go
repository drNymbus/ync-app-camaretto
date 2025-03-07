package model

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"camaretto/view"
)

type Button struct {
	width, height int
	message string
	textColor color.RGBA
	BackgroundColor color.RGBA
	SSprite *view.Sprite
}

func NewButton(width int, height int, msg string, textClr color.RGBA, bgClr color.RGBA) *Button {
	var b *Button = &Button{width, height, msg, textClr, bgClr, nil}
	b.Render()
	return b
}

func (b *Button) Render() {
	var img *ebiten.Image = ebiten.NewImage(b.width, b.height)
	
	var tWidth, tHeight float64 = text.Measure(b.message, view.TextFace, 0.0)
	op := &text.DrawOptions{}; op.ColorScale.ScaleWithColor(b.textColor)
	op.GeoM.Translate(float64(b.width)/2 - (tWidth/2), float64(b.height)/2 - (tHeight/2))
	text.Draw(img, b.message, &text.GoTextFace{Source: view.FaceSource, Size: view.FontSize}, op)

	b.SSprite = view.NewSprite(img, true, b.BackgroundColor, nil)
}

func (b *Button) SetMessage(msg string) {
	b.message = msg
	b.Render()
}

func (b *Button) SetTextColor(c color.RGBA) {
	b.textColor = c
	b.Render()
}