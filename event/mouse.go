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
	X, Y float64
	events []*MouseEvent
	capacity int
}

func NewMouse(capacity int) *Mouse {
	return &Mouse{0, 0, []*MouseEvent{}, capacity}
}

func (m *Mouse) Update() {
	var x, y int = ebiten.CursorPosition()
	m.X, m.Y = float64(x), float64(y)
	if len(m.events) >= m.capacity { return }

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		m.events = append(m.events, &MouseEvent{m.X, m.Y, ebiten.MouseButtonLeft, PRESSED})
	}
	// if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
	// 	m.events = append(m.events, &MouseEvent{m.X, m.Y, ebiten.MouseButtonRight, PRESSED})
	// }
	// if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonMiddle) {
	// 	m.events = append(m.events, &MouseEvent{m.X, m.Y, ebiten.MouseButtonMiddle, PRESSED})
	// }

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		m.events = append(m.events, &MouseEvent{m.X, m.Y, ebiten.MouseButtonLeft, RELEASED})
	}
	// if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonRight) {
	// 	m.events = append(m.events, &MouseEvent{m.X, m.Y, ebiten.MouseButtonRight, RELEASED})
	// }
	// if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonMiddle) {
	// 	m.events = append(m.events, &MouseEvent{m.X, m.Y, ebiten.MouseButtonMiddle, RELEASED})
	// }
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