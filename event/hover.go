package event

import (
	"image/color"

	"camaretto/model"
)

func hoverCard(c *model.Card, x float64, y float64) {
	if c.SSprite.In(x, y) {
		c.SSprite.EnableBackground()
	} else {
		c.SSprite.DisableBackground()
	}
}

func hoverButton(b *model.Button, x float64, y float64) {
	var color color.RGBA = b.BackgroundColor
	if b.SSprite.In(x, y) {
		color.A = 255
		b.SSprite.SetBackgroundColor(color)
	} else {
		color.A = 127
		b.SSprite.SetBackgroundColor(color)
	}
}

// @desc:
func HandleGameHover(app *model.Application, x float64, y float64) {
	hoverButton(app.Attack, x, y)
	hoverButton(app.Shield, x, y)
	hoverButton(app.Charge, x, y)
	hoverButton(app.Heal, x, y)

	for _, player := range app.Camaretto.Players {
		if player.HealthCard[0] != nil { hoverCard(player.HealthCard[0], x, y) }
		if player.HealthCard[1] != nil { hoverCard(player.HealthCard[1], x, y) }
		if player.ShieldCard != nil { hoverCard(player.ShieldCard, x, y) }
		if player.ChargeCard != nil { hoverCard(player.ChargeCard, x, y) }
	}	
}