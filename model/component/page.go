package component

import (
	"github.com/hajimehoshi/ebiten/v2"

	"camaretto/event"
)

type PageSignal int
const (
	UPDATE PageSignal = iota // Blank signal; nothing has to be done
	PREVIOUS
	NEXT
)

type Page interface {
	Init(w, h float64)

	MousePress(x, y float64) PageSignal
	MouseRelease(x, y float64) PageSignal
	HandleKeyEvent(e *event.KeyEvent) PageSignal
	Update() PageSignal

	Display(dst *ebiten.Image)
}
