package event

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"camaretto/model"
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
	var xi, yi int = ebiten.CursorPosition()
	var x, y float64 = float64(xi), float64(yi)
	m.X, m.Y = x, y
	if len(m.events) >= m.capacity { return }

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		m.events = append(m.events, &MouseEvent{x, y, ebiten.MouseButtonLeft, PRESSED})
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		m.events = append(m.events, &MouseEvent{x, y, ebiten.MouseButtonLeft, RELEASED})
	}
}

func (m *Mouse) ReadEvent() *MouseEvent {
	if len(m.events) == 0 { return nil }
	var me *MouseEvent = m.events[0]
	m.events = m.events[1:]
	return me
}

func (m *Mouse) IsEmpty() bool {
	return len(m.events) == 0
}

// @desc: Detect if click is on a player's card, then return the index of the player
func mouseOnPlayer(players []*model.Player, x float64, y float64) int {
	for i, player := range players {
		if !player.Dead {
			var onPlayer bool = false
			if player.HealthCard[0] != nil { onPlayer = onPlayer || player.HealthCard[0].SSprite.In(x, y) }
			if player.HealthCard[1] != nil { onPlayer = onPlayer || player.HealthCard[1].SSprite.In(x, y) }
			if player.JokerHealth != nil { onPlayer = onPlayer || player.JokerHealth.SSprite.In(x, y) }
			if player.ShieldCard != nil { onPlayer = onPlayer || player.ShieldCard.SSprite.In(x, y) }
			if player.ChargeCard != nil { onPlayer = onPlayer || player.ChargeCard.SSprite.In(x, y) }
	
			if onPlayer { return i }
		}
	}

	return -1
}

// @desc: Detect if player's health card have been clicked, then return the index of the card
func mouseOnHealthCard(p *model.Player, x float64, y float64) int {
	if p.HealthCard[0] != nil && p.HealthCard[0].SSprite.In(x, y) {
		return 0
	} else if p.HealthCard[1] != nil && p.HealthCard[1].SSprite.In(x, y) {
		return 1
	}

	return -1
}

// @desc: Detect if the click is on the draw pile of the deck
func mouseOnDeck(d *model.Deck, x float64, y float64) bool {
	for i := 0; i < d.LenDrawPile; i++ {
		if d.DrawPile[i].SSprite.In(x, y) { return true }
	}
	return false
}