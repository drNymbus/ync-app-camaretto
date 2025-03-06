package event

import (
	// "strconv"
	"camaretto/model"
)

// @desc:
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

// @desc:
func mouseOnHealthCard(p *model.Player, x float64, y float64) int {
	if p.HealthCard[0] != nil && p.HealthCard[0].SSprite.In(x, y) {
		return 0
	} else if p.HealthCard[1] != nil && p.HealthCard[1].SSprite.In(x, y) {
		return 1
	}

	return -1
}

// @desc:
func mouseOnDeck(d *model.Deck, x float64, y float64) bool {
	for i := 0; i < d.LenDrawPile; i++ {
		if d.DrawPile[i].SSprite.In(x, y) { return true }
	}
	return false
}

// @desc:
func HandleFocusPlayerRelease(camaretto *model.Camaretto, x float64, y float64) {
	var i int = mouseOnPlayer(camaretto.Players, x, y)
	if i != -1 {
		var state model.GameState = camaretto.GetState()
		if state == model.ATTACK {
			camaretto.SetFocus(model.CARD)
			camaretto.SetPlayerFocus(i)
		} else if state == model.SHIELD {
			camaretto.SetFocus(model.COMPLETE)
			camaretto.SetPlayerFocus(i)
		}
	}
}

// @desc:
func HandleFocusCardRelease(camaretto *model.Camaretto, x float64, y float64) {
	var player *model.Player = camaretto.Players[camaretto.GetPlayerFocus()]
	var i int = mouseOnHealthCard(player, x, y)
	if i != -1 {
		camaretto.SetCardFocus(i)
		camaretto.SetFocus(model.COMPLETE)
	}
}

// @desc:
func HandleFocusCompleteRelease(camaretto *model.Camaretto, x float64, y float64) {
	var state model.GameState = camaretto.GetState()
	if state == model.HEAL {
		var player *model.Player = camaretto.Players[camaretto.GetPlayerTurn()]
		if player.ChargeCard.SSprite.In(x, y) {
			camaretto.Heal()
			camaretto.EndTurn()
		}
	} else {
		if mouseOnDeck(camaretto.DeckPile, x, y) {
			if state == model.ATTACK {
				camaretto.Attack()
			} else if state == model.SHIELD {
				camaretto.Shield()
			} else if state == model.CHARGE {
				camaretto.Charge()
			}
			camaretto.EndTurn()
		}
	}
}

func HandleButtonRelease(app *model.Application, x float64, y float64) {
	if app.Attack.SSprite.In(x, y) {
		app.Camaretto.SetState(model.ATTACK)
		app.Camaretto.SetFocus(model.PLAYER)
	} else if app.Shield.SSprite.In(x, y) {
		app.Camaretto.SetState(model.SHIELD)
		app.Camaretto.SetFocus(model.PLAYER)
	} else if app.Charge.SSprite.In(x, y) {
		app.Camaretto.SetState(model.CHARGE)
		app.Camaretto.SetFocus(model.COMPLETE)
	} else if app.Heal.SSprite.In(x, y) {
		app.Camaretto.SetState(model.HEAL)
		var playerTurn int = app.Camaretto.GetPlayerTurn()
		app.Camaretto.SetPlayerFocus(playerTurn)
		app.Camaretto.SetFocus(model.CARD)
	}
}

// @desc:
func HandleCamarettoMouseRelease(app *model.Application, x float64, y float64) {
	var state model.GameState = app.Camaretto.GetState()
	var focus model.FocusState = app.Camaretto.GetFocus()
 
	if state == model.SET {
		HandleButtonRelease(app, x, y)
	} else {
		if focus == model.PLAYER {
			HandleFocusPlayerRelease(app.Camaretto, x, y)
		} else if focus == model.CARD {
			HandleFocusCardRelease(app.Camaretto, x, y)
		} else if focus == model.COMPLETE {
			HandleFocusCompleteRelease(app.Camaretto, x, y)
		}
	}
}