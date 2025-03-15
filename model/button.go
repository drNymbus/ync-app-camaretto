package model

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	// "github.com/hajimehoshi/ebiten/v2/text/v2"

	"camaretto/view"
)

type Button struct {
	// width, height int
	message string
	textColor color.RGBA

	pushed bool
	sourcePressedImg *ebiten.Image
	sourceReleasedImg *ebiten.Image
	pressedImg *ebiten.Image
	releasedImg *ebiten.Image

	SSprite *view.Sprite
}

func NewButton(msg string, textClr color.RGBA, buttonColor string) *Button {
	var b *Button = &Button{}
	// b.width, b.height = w, h
	b.message = msg
	b.textColor = textClr

	b.pushed = false

	// if buttonColor == "BLUE" {
	// 	b.sourcePressedImg = view.BlueButtonFlatImage
	// 	b.sourceReleasedImg = view.BlueButtonDepthImage
	// } else if buttonColor == "RED" {
	// 	b.sourcePressedImg = view.RedButtonFlatImage
	// 	b.sourceReleasedImg = view.RedButtonDepthImage
	// } else if buttonColor == "GREEN" {
	// 	b.sourcePressedImg = view.GreenButtonFlatImage
	// 	b.sourceReleasedImg = view.GreenButtonDepthImage
	// } else if buttonColor == "YELLOW" {
	// 	b.sourcePressedImg = view.YellowButtonFlatImage
	// 	b.sourceReleasedImg = view.YellowButtonDepthImage
	// }
	var bi *view.ButtonImage = view.LoadButtonImage(buttonColor)
	b.sourcePressedImg = bi.Pressed
	b.sourceReleasedImg = bi.Released

	b.Render()
	if b.pushed {
		b.SSprite = view.NewSprite(b.pressedImg, false, color.RGBA{0, 0, 0, 0}, nil)
	} else {
		b.SSprite = view.NewSprite(b.releasedImg, false, color.RGBA{0, 0, 0, 0}, nil)
	}

	return b
}

func (b *Button) Render() {
	// var tw, th float64 = text.Measure(b.message, view.TextFace, 0.0)
	// var op *text.DrawOptions = &text.DrawOptions{}
	var op *ebiten.DrawImageOptions = &ebiten.DrawImageOptions{}
	// op.ColorScale.ScaleWithColor(b.textColor)
	var textImage, tw, th = view.TextToImage(b.message, b.textColor)

	var w, h int = b.sourcePressedImg.Size()
	b.pressedImg = ebiten.NewImage(w, h)
	b.pressedImg.DrawImage(b.sourcePressedImg, nil)
	op.GeoM.Translate(float64(w)/2 - (tw/2), float64(h)/2 - (th/2))
	// text.Draw(b.pressedImg, b.message, &text.GoTextFace{Source: view.FaceSource, Size: view.FontSize}, op)
	b.pressedImg.DrawImage(textImage, op)

	w, h = b.sourceReleasedImg.Size()
	b.releasedImg = ebiten.NewImage(w, h)
	b.releasedImg.DrawImage(b.sourceReleasedImg, nil)
	op.GeoM.Reset()
	op.GeoM.Translate(float64(w)/2 - (tw/2), float64(h)/2 - (th/2))
	// text.Draw(b.releasedImg, b.message, &text.GoTextFace{Source: view.FaceSource, Size: view.FontSize}, op)
	b.releasedImg.DrawImage(textImage, op)
}

func (b *Button) SetMessage(msg string) {
	b.message = msg
	b.Render()
	if b.pushed {
		b.SSprite.SetImage(b.pressedImg)
	} else {
		b.SSprite.SetImage(b.releasedImg)
	}
}

func (b *Button) SetTextColor(c color.RGBA) {
	b.textColor = c
	b.Render()
	if b.pushed {
		b.SSprite.SetImage(b.pressedImg)
	} else {
		b.SSprite.SetImage(b.releasedImg)
	}
}

func (b *Button) Pressed() {
	b.pushed = true
	b.SSprite.SetImage(b.pressedImg)
}

func (b *Button) Released() {
	b.pushed = false
	b.SSprite.SetImage(b.releasedImg)
}