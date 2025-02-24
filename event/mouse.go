package event

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type MouseEventType int
const (
	PRESSED MouseEventType = 0
	RELEASED MouseEventType = 1
)

type MouseEvent struct {
	X, Y float64
	Click ebiten.MouseButton
	MET MouseEventType
}

type Mouse struct {
	events []*MouseEvent
	capacity int
}

func NewMouse(capacity int) *Mouse {
	return &Mouse{[]*MouseEvent{}, capacity}
}

func (m *Mouse) Update() {
	if len(m.events) >= m.capacity { return }

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		var x, y int = ebiten.CursorPosition()
		m.events = append(m.events, &MouseEvent{float64(x), float64(y), ebiten.MouseButtonLeft, PRESSED})
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		var x, y int = ebiten.CursorPosition()
		m.events = append(m.events, &MouseEvent{float64(x), float64(y), ebiten.MouseButtonRight, PRESSED})
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonMiddle) {
		var x, y int = ebiten.CursorPosition()
		m.events = append(m.events, &MouseEvent{float64(x), float64(y), ebiten.MouseButtonMiddle, PRESSED})
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		var x, y int = ebiten.CursorPosition()
		m.events = append(m.events, &MouseEvent{float64(x), float64(y), ebiten.MouseButtonLeft, RELEASED})
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonRight) {
		var x, y int = ebiten.CursorPosition()
		m.events = append(m.events, &MouseEvent{float64(x), float64(y), ebiten.MouseButtonRight, RELEASED})
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonMiddle) {
		var x, y int = ebiten.CursorPosition()
		m.events = append(m.events, &MouseEvent{float64(x), float64(y), ebiten.MouseButtonMiddle, RELEASED})
	}
}

func (m *Mouse) ReadEvent() *MouseEvent {
	if len(m.events) == 0 { return nil }
	var me *MouseEvent = m.events[0]
	m.events = m.events[1:]
	return me
}

func (m *Mouse) EmptyEventQueue() bool {
	return len(m.events) == 0
}