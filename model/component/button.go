package component

import (
	"log"

	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"camaretto/view"
)

type Button struct {
	message string
	textColor color.RGBA

	Trigger func()

	pushed bool
	sourcePressedImg *ebiten.Image
	sourceReleasedImg *ebiten.Image
	pressedImg *ebiten.Image
	releasedImg *ebiten.Image

	SSprite *Sprite
}

func NewButton(msg string, textClr color.RGBA, buttonColor string, onClick func()) *Button {
	var b *Button = &Button{}
	// b.width, b.height = w, h
	b.message = msg
	b.textColor = textClr

	b.Trigger = onClick

	b.pushed = false

	var bi *view.ButtonImage = view.LoadButtonImage(buttonColor)
	b.sourcePressedImg = bi.Pressed
	b.sourceReleasedImg = bi.Released

	b.render()
	b.SSprite = NewSprite(b.releasedImg, nil)

	return b
}

// @desc: Compute button image then set it to the sprite attribute
func (b *Button) render() {
	var op *ebiten.DrawImageOptions = &ebiten.DrawImageOptions{}
	var textImage, tw, th = view.TextToImage(b.message, b.textColor)

	var w, h int = b.sourcePressedImg.Size()
	b.pressedImg = ebiten.NewImage(w, h)
	b.pressedImg.DrawImage(b.sourcePressedImg, nil)
	op.GeoM.Translate(float64(w)/2 - (tw/2), float64(h)/2 - (th/2))
	b.pressedImg.DrawImage(textImage, op)

	w, h = b.sourceReleasedImg.Size()
	b.releasedImg = ebiten.NewImage(w, h)
	b.releasedImg.DrawImage(b.sourceReleasedImg, nil)
	op.GeoM.Reset(); op.GeoM.Translate(float64(w)/2 - (tw/2), float64(h)/2 - (th/2))
	b.releasedImg.DrawImage(textImage, op)
}

func (b *Button) pressed() {
	b.pushed = true
	b.SSprite.SetImage(b.pressedImg)
}

func (b *Button) released() {
	b.pushed = false
	b.SSprite.SetImage(b.releasedImg)
}

// @desc: Modify text on button
func (b *Button) SetMessage(msg string) {
	b.message = msg
	b.render()
	if b.pushed {
		b.SSprite.SetImage(b.pressedImg)
	} else {
		b.SSprite.SetImage(b.releasedImg)
	}
}

// @desc: Modify text's color
func (b *Button) SetTextColor(c color.RGBA) {
	b.textColor = c
	b.render()
	if b.pushed {
		b.SSprite.SetImage(b.pressedImg)
	} else {
		b.SSprite.SetImage(b.releasedImg)
	}
}

func (b *Button) Update() error {
	var err error
	var x, y int = ebiten.CursorPosition()

	var flagIn bool = b.SSprite.In(float64(x), float64(y))
	var flagPress bool = inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
	var flagRelease bool = inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft)

	if flagIn {
		if flagPress {
			b.pressed()
		} else if flagRelease {
			b.Trigger()
		}
	}

	if flagRelease {
		b.released()
	}

	err = b.SSprite.Update()
	if err != nil {
		log.Println("[Button.Update] Unable to update button sprite")
		return err
	}

	return nil
}

func (b *Button) Draw(screen *ebiten.Image) {
	b.SSprite.Draw(screen)
}
