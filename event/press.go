package event

import (
	"camaretto/model"
)

// @desc:
func HandleButtonPress(app *model.Application, x float64, y float64) {
	// if app.Attack.SSprite.In(x, y) {
	// 	app.Attack.SSprite.Scale(0.95, 0.95)
	// } else if app.Shield.SSprite.In(x, y) {
	// 	app.Shield.SSprite.Scale(0.95, 0.95)
	// } else if app.Charge.SSprite.In(x, y) {
	// 	app.Charge.SSprite.Scale(0.95, 0.95)
	// } else if app.Heal.SSprite.In(x, y) {
	// 	app.Heal.SSprite.Scale(0.95, 0.95)
	// }
}

// @desc:
func HandleCamarettoMousePress(app *model.Application, x float64, y float64) {
	var state model.GameState = app.Camaretto.GetState()

	if state == model.SET {
		HandleButtonPress(app, x, y)
	}
}