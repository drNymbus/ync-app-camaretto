package event

import (
	// "log"
	// "reflect"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type EventType int
const (
	PRESSED EventType = 0
	RELEASED EventType = 1
)

type MouseEvent struct {
	X, Y float64
	Click ebiten.MouseButton
	Event EventType
}

type KeyEvent struct {
	Key ebiten.Key
	Event EventType
}

type EventQueue struct {
	X, Y float64
	mouse []*MouseEvent
	keyboard []*KeyEvent
	capacity int
}

func NewEventQueue(capacity int) *EventQueue {
	return &EventQueue{0, 0, []*MouseEvent{}, []*KeyEvent{}, capacity}
}

func (q *EventQueue) Update() {
	var xi, yi int = ebiten.CursorPosition()
	q.X, q.Y = float64(xi), float64(yi)

	if len(q.mouse) < q.capacity {
		// Mouse
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			q.mouse = append(q.mouse, &MouseEvent{q.X, q.Y, ebiten.MouseButtonLeft, PRESSED})
		}
		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
			q.mouse = append(q.mouse, &MouseEvent{q.X, q.Y, ebiten.MouseButtonLeft, RELEASED})
		}

		// Keyboard
		var keys []ebiten.Key = make([]ebiten.Key, q.capacity)
		keys = inpututil.AppendJustPressedKeys(keys[:0])
		for _, p := range keys {
			q.keyboard = append(q.keyboard, &KeyEvent{p, PRESSED})
		}
		keys = inpututil.AppendJustReleasedKeys(keys[:0])
		for _, p := range keys {
			q.keyboard = append(q.keyboard, &KeyEvent{p, RELEASED})
		}
	}
}

func (q *EventQueue) ReadMouseEvent() *MouseEvent {
	if len(q.mouse) == 0 { return nil }
	var me *MouseEvent = q.mouse[0]
	q.mouse = q.mouse[1:]
	return me
}

func (q *EventQueue) ReadKeyEvent() *KeyEvent {
	if len(q.keyboard) == 0 { return nil }
	var ke *KeyEvent = q.keyboard[0]
	q.keyboard = q.keyboard[1:]
	return ke
}

func (q *EventQueue) IsEmpty() bool {
	return (len(q.mouse) == 0 && len(q.keyboard) == 0)
}