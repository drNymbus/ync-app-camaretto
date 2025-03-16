package event

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Keyboard struct {
	Shift bool
	Alt bool
	AltGr bool
}

func NewKeyboard() *Keyboard { return &Keyboard{false, false, false} }

func (k *Keyboard) Update() {
	if inpututil.IsKeyJustPressed(ebiten.KeyShiftLeft) { k.Shift = true }
	if inpututil.IsKeyJustPressed(ebiten.KeyShiftRight) { k.Shift = true }
	if inpututil.IsKeyJustPressed(ebiten.KeyAltLeft) { k.Alt = true }
	if inpututil.IsKeyJustPressed(ebiten.KeyAltRight) { k.AltGr = true }

	if inpututil.IsKeyJustReleased(ebiten.KeyShiftLeft) { k.Shift = false }
	if inpututil.IsKeyJustReleased(ebiten.KeyShiftRight) { k.Shift = false }
	if inpututil.IsKeyJustReleased(ebiten.KeyAltLeft) { k.Alt = false }
	if inpututil.IsKeyJustReleased(ebiten.KeyAltRight) { k.AltGr = false }
}