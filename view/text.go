package view

import (
	"log"

	// "os"
	// "io"
	"bytes"
	// "image"
	"image/color"

	"golang.org/x/text/language"

	"github.com/hajimehoshi/ebiten/v2"
	// "github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var (
	FaceSource *text.GoTextFaceSource
	TextFace *text.GoTextFace
	FontSize float64 = 24
)

func LoadFont() {
	var err error
	// Load font file
	var fontByte []byte = GetFileByte("assets/fonts/NaturalMono_Regular.ttf")
	// var fontByte []byte = getFileByte("assets/fonts/Kenney_Future_Narrow.ttf")
	FaceSource, err = text.NewGoTextFaceSource(bytes.NewReader(fontByte))
	if err != nil { log.Fatal("[parametersInitAssets] Set FaceSource:", err) }

	TextFace = &text.GoTextFace{
		Source: FaceSource,
		Direction: text.DirectionLeftToRight,
		Size: 24, Language: language.English,
	}
}

func TextToImage(s string, clr color.RGBA) (*ebiten.Image, float64, float64) {
	var w, h float64 = text.Measure(s, TextFace, 0.0)

	var op *text.DrawOptions = &text.DrawOptions{}
	op.ColorScale.ScaleWithColor(clr)

	var img *ebiten.Image = ebiten.NewImage(int(w), int(h))
	text.Draw(img, s, &text.GoTextFace{Source: FaceSource, Size: FontSize}, op)
	return img, w, h
}