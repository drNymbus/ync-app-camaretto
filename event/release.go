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
			camaretto.SetFocus(model.REVEAL)
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
		camaretto.SetFocus(model.REVEAL)
	}
}

// @desc:
func HandleFocusRevealRelease(camaretto *model.Camaretto, x float64, y float64) {
	var state model.GameState = camaretto.GetState()
	if state == model.HEAL {
		var player *model.Player = camaretto.Players[camaretto.GetPlayerTurn()]
		if player.ChargeCard.SSprite.In(x, y) {
			camaretto.Heal()
			camaretto.SetFocus(model.COMPLETE)
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
			camaretto.SetFocus(model.COMPLETE)
		}
	}
}

// @desc:
func HandleFocusCompleteRelease(camaretto *model.Camaretto, x float64, y float64) {
	// var state model.GameState = camaretto.GetState()
	// if state == model.HEAL {
	// 	var player *model.Player = camaretto.Players[camaretto.GetPlayerTurn()]
	// 	if player.ChargeCard.SSprite.In(x, y) {
	// 		camaretto.Heal()
	// 		camaretto.SetFocus(model.COMPLETE)
	// 	}
	// } else {
	// 	if mouseOnDeck(camaretto.DeckPile, x, y) {
	// 		if state == model.ATTACK {
	// 			pTurn := camaretto.Players[camaretto.GetPlayerTurn()].Name
	// 			pFocus := camaretto.Players[camaretto.GetPlayerFocus()].Name
	// 			log.Println(pTurn, "attack", pFocus)
	// 			camaretto.Attack()
	// 		} else if state == model.SHIELD {
	// 			pTurn := camaretto.Players[camaretto.GetPlayerTurn()].Name
	// 			pFocus := camaretto.Players[camaretto.GetPlayerFocus()].Name
	// 			log.Println(pTurn, "shield", pFocus)
	// 			camaretto.Shield()
	// 		} else if state == model.CHARGE {
	// 			camaretto.Charge()
	// 		}
	// 		camaretto.SetFocus(model.COMPLETE)
	// 	}
	// }
}

// @desc:
func HandleButtonRelease(app *model.Application, x float64, y float64) {
	if app.Attack.SSprite.In(x, y) {
		log.Println("ATTACK")
		app.Camaretto.SetState(model.ATTACK)
		app.Camaretto.SetFocus(model.PLAYER)
	} else if app.Shield.SSprite.In(x, y) {
		log.Println("SHIELD")
		app.Camaretto.SetState(model.SHIELD)
		app.Camaretto.SetFocus(model.PLAYER)
	} else if app.Charge.SSprite.In(x, y) {
		log.Println("CHARGE")
		app.Camaretto.SetState(model.CHARGE)
		var playerTurn int = app.Camaretto.GetPlayerTurn()
		app.Camaretto.SetPlayerFocus(playerTurn)
		app.Camaretto.SetFocus(model.COMPLETE)
	} else if app.Heal.SSprite.In(x, y) {
		log.Println("HEAL")
		app.Camaretto.SetState(model.HEAL)
		var playerTurn int = app.Camaretto.GetPlayerTurn()
		app.Camaretto.SetPlayerFocus(playerTurn)
		app.Camaretto.SetFocus(model.CARD)
	}
}

// @desc:
func HandleCamarettoMouseRelease(app *model.Application, x float64, y float64) {
	var state model.GameState = app.Camaretto.GetState()
	if state == model.SET {
		HandleButtonRelease(app, x, y)
	} else {
		var focus model.FocusState = app.Camaretto.GetFocus()
		if focus == model.PLAYER {
			HandleFocusPlayerRelease(app.Camaretto, x, y)
		} else if focus == model.CARD {
			HandleFocusCardRelease(app.Camaretto, x, y)
		} else if focus == model.REVEAL {
			HandleFocusCompleteRelease(app.Camaretto, x, y)
		} else if focus == model.COMPLETE {
			HandleFocusCompleteRelease(app.Camaretto, x, y)
		}
	}
}