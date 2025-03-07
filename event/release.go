package event

import (
	"log"
	// "strconv"
	"camaretto/model"
)

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
			camaretto.SetState(model.END)
		}
	} else {
		if mouseOnDeck(camaretto.DeckPile, x, y) {
			if state == model.ATTACK {
				pTurn := camaretto.Players[camaretto.GetPlayerTurn()].Name
				pFocus := camaretto.Players[camaretto.GetPlayerFocus()].Name
				log.Println(pTurn, "attack", pFocus)
				camaretto.Attack()
			} else if state == model.SHIELD {
				pTurn := camaretto.Players[camaretto.GetPlayerTurn()].Name
				pFocus := camaretto.Players[camaretto.GetPlayerFocus()].Name
				log.Println(pTurn, "shield", pFocus)
				camaretto.Shield()
			} else if state == model.CHARGE {
				camaretto.Charge()
			}
			camaretto.SetState(model.END)
		}
	}
}

// @desc:
func HandleButtonRelease(app *model.Application, x float64, y float64) {
	if app.Attack.SSprite.In(x, y) {
		app.Camaretto.SetState(model.ATTACK)
		app.Camaretto.SetFocus(model.PLAYER)
	} else if app.Shield.SSprite.In(x, y) {
		app.Camaretto.SetState(model.SHIELD)
		app.Camaretto.SetFocus(model.PLAYER)
	} else if app.Charge.SSprite.In(x, y) {
		app.Camaretto.SetState(model.CHARGE)
		var playerTurn int = app.Camaretto.GetPlayerTurn()
		app.Camaretto.SetPlayerFocus(playerTurn)
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

	app.Attack.SSprite.Scale(1, 1)
	app.Shield.SSprite.Scale(1, 1)
	app.Charge.SSprite.Scale(1, 1)
	app.Heal.SSprite.Scale(1, 1)

	if state == model.SET {
		HandleButtonRelease(app, x, y)
	} else if state == model.END {
		app.Camaretto.EndTurn()
		app.Camaretto.SetState(model.SET)
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