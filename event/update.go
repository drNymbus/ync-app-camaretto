package event

import (
	"camaretto/model"
)

// @desc:
func HandleFocusPlayer(players []*model.Player, e *MouseEvent) int {
	for i, player := range players {
		if !player.Dead {
			var onPlayer bool = false
			if player.HealthCard[0] != nil { onPlayer = onPlayer || player.HealthCard[0].SSprite.In(e.X, e.Y) }
			if player.HealthCard[1] != nil { onPlayer = onPlayer || player.HealthCard[1].SSprite.In(e.X, e.Y) }
			if player.JokerHealth != nil { onPlayer = onPlayer || player.JokerHealth.SSprite.In(e.X, e.Y) }
			if player.ShieldCard != nil { onPlayer = onPlayer || player.ShieldCard.SSprite.In(e.X, e.Y) }
			if player.ChargeCard != nil { onPlayer = onPlayer || player.ChargeCard.SSprite.In(e.X, e.Y) }
	
			if onPlayer { return i }
		}
	}

	return -1
}

// @desc:
func HandleFocusCard(p *model.Player, e *MouseEvent) int {
	if p.HealthCard[0] != nil && p.HealthCard[0].SSprite.In(e.X, e.Y) {
		return 0
	} else if p.HealthCard[1] != nil && p.HealthCard[1].SSprite.In(e.X, e.Y) {
		return 1
	}

	return -1
}

// @desc:
func HandleFocusComplete(d *model.Deck, e *MouseEvent) bool {
	for i := 0; i < d.LenDrawPile; i++ {
		if d.DrawPile[i].SSprite.In(e.X, e.Y) { return true }
	}
	return false
}