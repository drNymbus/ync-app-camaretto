package component

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
	var tc *TextCapture = &TextCapture{}
	tc.width, tc.height = w, h
	tc.textInput = ""
	tc.charLimit = limit

	tc.image = ebiten.NewImage(w, h)
	tc.image.Fill(color.RGBA{0, 0, 0, 255})

	tc.margin = margin
	var filling *ebiten.Image = ebiten.NewImage(w - tc.margin*2, h - tc.margin*2)
	filling.Fill(color.RGBA{255, 255, 255, 255})

	var op *ebiten.DrawImageOptions = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(tc.margin), float64(tc.margin))

	tc.image.DrawImage(filling, op)
	tc.SSprite = view.NewSprite(tc.image, false, color.RGBA{}, nil)

	return tc
}

func (tc *TextCapture) SetText(s string) { tc.textInput = s }
func (tc *TextCapture) GetText() string { return tc.textInput }

func (tc *TextCapture) render() {
	if len(tc.textInput) > 0 {
		var img, textImg *ebiten.Image
		img = ebiten.NewImage(tc.width, tc.height)
		img.DrawImage(tc.image, nil)
	
		var th float64
		textImg, _, th = view.TextToImage(tc.textInput, color.RGBA{0, 0, 0, 255})
		var op *ebiten.DrawImageOptions = &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(tc.margin) + 25, float64(tc.height/2) - th/2)
		img.DrawImage(textImg, op)
	
		tc.SSprite.SetImage(img)
	} else { tc.SSprite.SetImage(tc.image) }
}

func (tc *TextCapture) HandleEvent(e *event.KeyEvent, k *event.Keyboard) {
	if len(tc.textInput) < tc.charLimit && e.Key < 27 {
		tc.textInput = tc.textInput + string(e.Key + 65)
		tc.render()
	} else if len(tc.textInput) > 0 && e.Key == ebiten.KeyBackspace {
		tc.textInput = tc.textInput[:len(tc.textInput)-1]
		tc.render()
	}
}
