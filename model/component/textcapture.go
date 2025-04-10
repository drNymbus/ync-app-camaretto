package component

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"camaretto/view"
)

type TextCapture struct {
	width, height int
	margin int

	textInput string
	charLimit int

	background *ebiten.Image
	SSprite *view.Sprite

	count int
}

func NewTextCapture(limit, w, h, margin int) *TextCapture {
	var tc *TextCapture = &TextCapture{}
	tc.width, tc.height = w, h
	tc.textInput = ""
	tc.charLimit = limit

	tc.background = ebiten.NewImage(w, h)
	tc.background.Fill(color.RGBA{0, 0, 0, 255})

	tc.margin = margin
	var filling *ebiten.Image = ebiten.NewImage(w - tc.margin*2, h - tc.margin*2)
	filling.Fill(color.RGBA{255, 255, 255, 255})

	var op *ebiten.DrawImageOptions = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(tc.margin), float64(tc.margin))

	tc.background.DrawImage(filling, op)
	tc.SSprite = view.NewSprite(tc.background, nil)

	tc.count = 0

	return tc
}

func (tc *TextCapture) render() {
	if len(tc.textInput) > 0 {
		var img, textImg *ebiten.Image
		img = ebiten.NewImage(tc.width, tc.height)
		img.DrawImage(tc.background, nil)
	
		var th float64
		textImg, _, th = view.TextToImage(tc.textInput, color.RGBA{0, 0, 0, 255})
		var op *ebiten.DrawImageOptions = &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(tc.margin) + 25, float64(tc.height/2) - th/2)
		img.DrawImage(textImg, op)
	
		tc.SSprite.SetImage(img)
	} else { tc.SSprite.SetImage(tc.background) }
}

func (tc *TextCapture) SetText(s string) {
	tc.textInput = s
	tc.render()
}
func (tc *TextCapture) GetText() string { return tc.textInput }

func (tc *TextCapture) Update() error {
	var shiftModifier int = 0

	var keys []ebiten.Key = make([]ebiten.Key, 50)

	keys = inpututil.AppendPressedKeys(keys[:0])
	for _, k := range keys {
		if k == ebiten.KeyShiftLeft || k == ebiten.KeyShiftRight {
			shiftModifier = -32 // 97(='a') - 65(='A')
		}
	}

	keys = inpututil.AppendJustPressedKeys(keys[:0])
	for _, k := range keys {
		if k == ebiten.KeyBackspace {
			if len(tc.textInput) > 0 {
				tc.textInput = tc.textInput[:len(tc.textInput)-1]
			}
		} else if k < 27 { // A letter key is pressed
			tc.textInput = tc.textInput + string(int(k) + 97 + shiftModifier)
		}
	}

	tc.render()
	return nil
}

func (tc *TextCapture) Draw(screen *ebiten.Image) {
	tc.SSprite.Draw(screen)
}
