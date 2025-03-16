package view

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	CardWidth int = 64
	CardHeight int = 64
)

type CardImage struct {
	Card [13]*ebiten.Image
	Joker *ebiten.Image
	Empty *ebiten.Image
	Hidden *ebiten.Image
}

func LoadCardImage() *CardImage {
	var cs *CardImage = &CardImage{}

	var ogSheet *ebiten.Image = GetImage("assets/cards/cardsLarge_tilemap_packed.png")
	// Scale down the original sheet
	var width, height int = ogSheet.Size()
	// var xScale, yScale float64 = 1, 1

	// width, height = int(float64(width) * xScale), int(flloat64(height) * yScale)
	// CardWidth, CardHeight = int(float64(CardWidth) * xScale), int(float64(CardHeight) * yScale)

	// var Sheet *ebiten.Image = ebiten.NewImage(width, height)
	var Sheet *ebiten.Image = ebiten.NewImage(width, height)
	op := &ebiten.DrawImageOptions{};
	// op.GeoM.Scale(xScale, yScale)
	Sheet.DrawImage(ogSheet, op)

	for i := 0; i < 13; i++ { // Init all cards image from Ace to King
		var sx int = i * CardWidth
		var img *ebiten.Image = Sheet.SubImage(image.Rect(sx, 0, sx+CardWidth, CardHeight)).(*ebiten.Image)
		cs.Card[i] = img
	}

	// All other cards are not logically placed in the tilemap sheet
	cs.Joker = Sheet.SubImage(image.Rect((13*CardWidth), (2*CardHeight), (13*CardWidth) + CardWidth, (2*CardHeight) + CardHeight)).(*ebiten.Image)
	cs.Empty = Sheet.SubImage(image.Rect((13*CardWidth), 0, (13*CardWidth) + CardWidth, CardHeight)).(*ebiten.Image)
	cs.Hidden = Sheet.SubImage(image.Rect((13*CardWidth), CardHeight, (13*CardWidth) + CardWidth, CardHeight + CardHeight)).(*ebiten.Image)

	return cs
}