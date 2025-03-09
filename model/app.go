package model

import (
	// "log"
	// "math"
	// "strconv"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	// "camaretto/view"
)

const (
	// WinWidth int = 640
	// WinHeight int = 480
	WinWidth int = 1200
	WinHeight int = 900
	ButtonWidth int = WinWidth / 5
	ButtonHeight int = WinHeight / 6
)

type AppState int
const (
	MENU AppState = 0
	GAME AppState = 1
)

type Application struct{
	state AppState
	Camaretto *Camaretto

	Attack *Button
	Shield *Button
	Charge *Button
	Heal *Button

	Info *Button
}

func (app *Application) Init(nbPlayers int) {
	app.state = GAME

	app.Camaretto = NewCamaretto(nbPlayers)

	app.Attack = NewButton(ButtonWidth, ButtonHeight, "ATTACK", color.RGBA{0, 0, 0, 255}, color.RGBA{163, 3, 9, 127})
	app.Shield = NewButton(ButtonWidth, ButtonHeight, "SHIELD", color.RGBA{0, 0, 0, 255}, color.RGBA{2, 42, 201, 127})
	app.Charge = NewButton(ButtonWidth, ButtonHeight, "CHARGE", color.RGBA{0, 0, 0, 255}, color.RGBA{224, 144, 4, 127})
	app.Heal = NewButton(ButtonWidth, ButtonHeight, "HEAL", color.RGBA{0, 0, 0, 255}, color.RGBA{3, 173, 18, 127})

	app.Info = NewButton(WinWidth, WinHeight/16, "This contains information.", color.RGBA{0, 0, 0, 255}, color.RGBA{127, 127, 127, 255})
}

/************ ***************************************************************************** ************/
/************ ********************************** GET/SET ********************************** ************/
/************ ***************************************************************************** ************/

func (app *Application) SetState(s AppState) { app.state = s }
func (app *Application) GetState() AppState { return app.state }

/************ *************************************************************************** ************/
/************ ********************************* UPDATE ********************************** ************/
/************ *************************************************************************** ************/

func (app *Application) Update() {
	var player *Player = app.Camaretto.Players[app.Camaretto.GetPlayerTurn()]
	var playerName string = player.Name

	var msgInfo string = ""
	if app.state == GAME {
		var state GameState = app.Camaretto.GetState()
		if state == SET {
			msgInfo = playerName + " needs to choose an action"
		} else if state == END {
			app.Camaretto.EndTurn()
			app.Camaretto.SetState(SET)
			msgInfo = playerName + " end turn"
		} else {
			var focus FocusState = app.Camaretto.GetFocus()
			if focus == PLAYER {
				msgInfo = playerName + " needs to select a player"
			} else if focus == CARD {
				player = app.Camaretto.Players[app.Camaretto.GetPlayerFocus()]
				playerName = player.Name
				msgInfo = playerName + " needs to select a health card"
			} else if focus == COMPLETE {
				msgInfo = playerName + " reveal card"
			}
		}
		app.Info.SetMessage(msgInfo)
	}
}

/************ *************************************************************************** ************/
/************ ********************************** RENDER ********************************* ************/
/************ *************************************************************************** ************/

func (app *Application) RenderButtons(dst *ebiten.Image) {
	var buttonXPos float64 = 0
	var buttonYPos float64 = float64(WinHeight)*9/10

	buttonXPos = float64(WinWidth)/2
	app.Info.SSprite.SetCenter(buttonXPos, buttonYPos - 120, 0)
	app.Info.SSprite.Display(dst)

	buttonXPos = (float64(WinWidth) * 1/4) + (float64(ButtonWidth)/2)
	app.Attack.SSprite.SetCenter(buttonXPos, buttonYPos, 0)
	app.Attack.SSprite.Display(dst)

	buttonXPos = (float64(WinWidth) * 2/4) + (float64(ButtonWidth)/2)
	app.Shield.SSprite.SetCenter(buttonXPos, buttonYPos, 0)
	app.Shield.SSprite.Display(dst)

	buttonXPos = (float64(WinWidth) * 3/4) + (float64(ButtonWidth)/2)

	var c *Camaretto = app.Camaretto
	if c.Players[c.GetPlayerTurn()].ChargeCard == nil {
		app.Heal.SSprite.SetCenter(0, 0, 0)
		app.Heal.SSprite.SetOffset(0, 0, 0)
		app.Charge.SSprite.SetCenter(buttonXPos, buttonYPos, 0)
		app.Charge.SSprite.Display(dst)
	} else {
		app.Charge.SSprite.SetCenter(0, 0, 0)
		app.Charge.SSprite.SetOffset(0, 0, 0)
		app.Heal.SSprite.SetCenter(buttonXPos, buttonYPos, 0)
		app.Heal.SSprite.Display(dst)
	}
}

func (app *Application) Display(dst *ebiten.Image) {
	if app.state == GAME {
		app.RenderButtons(dst)
		app.Camaretto.Render(dst, float64(WinWidth), float64(WinHeight))
	}
}