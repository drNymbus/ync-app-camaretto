package model

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"camaretto/view"
	"camaretto/event"
)

type TextCapture struct {
	width, height int
	margin int

	textInput string
	charLimit int

	image *ebiten.Image
	SSprite *view.Sprite
}

func NewTextCapture(limit, w, h, margin int) *TextCapture {
	var ta *TextCapture = &TextCapture{}
	ta.width, ta.height = w, h
	ta.textInput = ""
	ta.charLimit = limit

	ta.image = ebiten.NewImage(w, h)
	ta.image.Fill(color.RGBA{0, 0, 0, 255})

	ta.margin = margin
	var filling *ebiten.Image = ebiten.NewImage(w - ta.margin*2, h - ta.margin*2)
	filling.Fill(color.RGBA{255, 255, 255, 255})

	var op *ebiten.DrawImageOptions = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(ta.margin), float64(ta.margin))

	ta.image.DrawImage(filling, op)
	ta.SSprite = view.NewSprite(ta.image, false, color.RGBA{}, nil)

	return ta
}

func (ta *TextCapture) GetText() string { return ta.textInput }

func (ta *TextCapture) render() {
	if len(ta.textInput) > 0 {
		var img, textImg *ebiten.Image
		img = ebiten.NewImage(ta.width, ta.height)
		img.DrawImage(ta.image, nil)
	
		var th float64
		textImg, _, th = view.TextToImage(ta.textInput, color.RGBA{0, 0, 0, 255})
		var op *ebiten.DrawImageOptions = &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(ta.margin) + 25, float64(ta.height/2) - th/2)
		img.DrawImage(textImg, op)
	
		ta.SSprite.SetImage(img)
	} else { ta.SSprite.SetImage(ta.image) }
}

func (ta *TextCapture) HandleEvent(e *event.KeyEvent, k *event.Keyboard) {
	if len(ta.textInput) < ta.charLimit && e.Key < 27 {
		ta.textInput = ta.textInput + string(e.Key + 65)
		ta.render()
	} else if len(ta.textInput) > 0 && e.Key == ebiten.KeyBackspace {
		ta.textInput = ta.textInput[:len(ta.textInput)-1]
		ta.render()
	}
}