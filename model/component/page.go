package component

import (
	"github.com/hajimehoshi/ebiten/v2"

	"camaretto/event"
)

type PageSignal int
const (
	PREVIOUS PageSignal = iota
	NEXT
	UPDATE // This is basically a NONE statement, nothing happened and nothing has to be done.
)

type Page interface {
	Init(w, h float64)

	MousePress(x, y float64) PageSignal
	MouseRelease(x, y float64) PageSignal
	HandleKeyEvent(e *event.KeyEvent) PageSignal
	Update() PageSignal

	Display(dst *ebiten.Image)
}
