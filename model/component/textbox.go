package component

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

	background *ebiten.Image
	border *ebiten.Image

	image *ebiten.Image

	SSprite *view.Sprite
}

func wrapText(message string, width int) string {
	var step int = width/int(view.FontSize*2/3)
	for i := step; i < len(message); i = i + step {
		message = message[:i] + "\n" + message[i:]
	}
	return message
}

func NewTextBox(w, h float64, msg string, textColor color.RGBA, backgroundColor color.RGBA) *TextBox {
	var tb *TextBox = &TextBox{}

	tb.width, tb.height = w, h

	tb.message = wrapText(msg, int(w))
	tb.textColor = textColor

	tb.length = 0
	tb.count = 0
	tb.speed = 5

	tb.renderBackground(backgroundColor)
	tb.renderBorder()

	tb.SSprite = view.NewSprite(tb.background, nil)

	return tb
}

func (tb *TextBox) renderBackground(c color.RGBA) {
	tb.background = ebiten.NewImage(int(tb.width), int(tb.height))
	tb.background.Fill(c)
}

func (tb *TextBox) renderBorder() {
	tb.border = ebiten.NewImage(int(tb.width), int(tb.height))

	var ii *view.IconImage = view.LoadIconImage()

	var barWidth, barHeight float64 = float64(view.BarWidth), float64(view.BarHeight)
	var barScale float64 = 0.5
	barHeight = barHeight * barScale

	var topMargin, bottomMargin float64 = barHeight + 25, barHeight + 25
	var leftMargin, rightMargin float64 = barHeight + 30, barHeight + 30

	var scaleWidth float64 = (tb.width - (leftMargin+rightMargin) - 50) / barWidth
	var scaleHeight float64 = (tb.height - (topMargin+bottomMargin) - 50) / barWidth

	// Top border
	var op *ebiten.DrawImageOptions = &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scaleWidth, barScale)
	op.GeoM.Translate(-(barWidth*scaleWidth)/2, -barHeight/2)
	op.GeoM.Translate(tb.width/2, topMargin/2)
	tb.border.DrawImage(ii.Bar, op)

	// Bottom border
	op.GeoM.Reset()
	op.GeoM.Scale(scaleWidth, barScale)
	op.GeoM.Translate(-(barWidth*scaleWidth)/2, -barHeight/2)
	op.GeoM.Translate(tb.width/2, tb.height - (bottomMargin/2))
	tb.border.DrawImage(ii.Bar, op)

	// Left border
	op.GeoM.Reset()
	op.GeoM.Scale(scaleHeight, barScale)
	op.GeoM.Translate(-(barWidth*scaleHeight)/2, -barHeight/2)
	op.GeoM.Rotate(math.Pi/2)
	op.GeoM.Translate(leftMargin/2, tb.height/2)
	tb.border.DrawImage(ii.Bar, op)

	// Right border
	op.GeoM.Reset()
	op.GeoM.Scale(scaleHeight, barScale)
	op.GeoM.Translate(-(barWidth*scaleHeight)/2, -barHeight/2)
	op.GeoM.Rotate(math.Pi/2)
	op.GeoM.Translate(tb.width - (rightMargin/2), tb.height/2)
	tb.border.DrawImage(ii.Bar, op)

	var iconScale float64 = 0.1
	var w, h float64 = 0, 0

	w, h = float64(view.CoffeeWidth)*iconScale, float64(view.CoffeeHeight)*iconScale
	// Top left coffee icon
	op.GeoM.Reset()
	op.GeoM.Scale(iconScale, iconScale)
	op.GeoM.Translate(-w/2, -h/2)
	op.GeoM.Translate(leftMargin/2, topMargin/2)
	tb.border.DrawImage(ii.Coffee, op)
	// Bottom right coffee icon
	op.GeoM.Reset()
	op.GeoM.Scale(iconScale, iconScale)
	op.GeoM.Translate(-w/2, -h/2)
	op.GeoM.Translate(tb.width - rightMargin/2, tb.height - bottomMargin/2)
	tb.border.DrawImage(ii.Coffee, op)

	w, h = float64(view.AmarettoWidth)*iconScale, float64(view.AmarettoHeight)*iconScale
	// Top right amaretto icon
	op.GeoM.Reset()
	op.GeoM.Scale(iconScale, iconScale)
	op.GeoM.Translate(-w/2, -h/2)
	op.GeoM.Translate(tb.width - rightMargin/2, topMargin/2)
	tb.border.DrawImage(ii.Amaretto, op)
	// Bottom left amaretto icon
	op.GeoM.Reset()
	op.GeoM.Scale(iconScale, iconScale)
	op.GeoM.Translate(-w/2, -h/2)
	op.GeoM.Translate(leftMargin/2, tb.height - bottomMargin/2)
	tb.border.DrawImage(ii.Amaretto, op)
}

// @desc: true if text is entirely displayed, false otherwise
func (tb *TextBox) Finished() bool { return tb.length == len(tb.message) }

func (tb *TextBox) SetMessage(msg string) {
	tb.message = wrapText(msg, int(tb.width))
	tb.length = 0
	tb.count = 0
}

func (tb *TextBox) SetBackgroundColor(c color.RGBA) {
	tb.renderBackground(c)
}

func (tb *TextBox) Update() error {
	var modify bool = false
	if tb.length < len(tb.message) {
		tb.count++
		if tb.count == tb.speed {
			tb.length++
			tb.count = 0
			modify = true
		}
	}

	if modify {
		var image *ebiten.Image = ebiten.NewImage(int(tb.width), int(tb.height))
		image.DrawImage(tb.background, nil)
		image.DrawImage(tb.border, nil)

		var op *text.DrawOptions = &text.DrawOptions{}
		op.ColorScale.ScaleWithColor(tb.textColor)
		op.LayoutOptions.LineSpacing = view.FontSize + 3
		op.GeoM.Reset()
		op.GeoM.Translate(tb.leftMargin*3/2, tb.topMargin)
		text.Draw(image, tb.message[:tb.length], &text.GoTextFace{Source: view.FaceSource, Size: view.FontSize}, op)
	
		tb.SSprite.SetImage(image)
	}

	tb.SSprite.Update()
	return nil
}

func (tb *TextBox) Draw(screen *ebiten.Image) {
	tb.SSprite.Draw(screen)
}
