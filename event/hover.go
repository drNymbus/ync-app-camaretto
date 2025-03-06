package event

import (
	"camaretto/model"
)

// @desc:
func HandleGameHover(app *model.Application, x float64, y float64) {
	app.Attack.Hover(x, y)
	app.Shield.Hover(x, y)
	app.Charge.Hover(x, y)
	app.Heal.Hover(x, y)

	for _, player := range app.Camaretto.Players {
		if player.HealthCard[0] != nil { player.HealthCard[0].Hover(x, y) }
		if player.HealthCard[1] != nil { player.HealthCard[1].Hover(x, y) }
		if player.ShieldCard != nil { player.ShieldCard.Hover(x, y) }
		if player.ChargeCard != nil { player.ChargeCard.Hover(x, y) }
	}	
}