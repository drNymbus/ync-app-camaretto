package model

import (
	// "log"
	"math"
	"strconv"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
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
	app.Info.SetMessage("PLAYER" + strconv.Itoa(app.Camaretto.GetPlayerTurn()) + ": Choose an action.")
}

/************ ***************************************************************************** ************/
/************ ********************************** GET/SET ********************************** ************/
/************ ***************************************************************************** ************/

func (app *Application) SetState(s AppState) { app.state = s }
func (app *Application) GetState() AppState { return app.state }

/************ *************************************************************************** ************/
/************ ********************************* UPDATE ********************************** ************/
/************ *************************************************************************** ************/

func (app *Application) Update() {}

/************ *************************************************************************** ************/
/************ *********************************** DRAW ********************************** ************/
/************ *************************************************************************** ************/

func (app *Application) DrawPlayers(dst *ebiten.Image) {
	var nbPlayers int = len(app.Camaretto.Players)
	var angleStep float64 = 2*math.Pi / float64(nbPlayers)
	var radius float64 = 200

	var centerX float64 = float64(WinWidth)/2
	var centerY float64 = (float64(WinHeight) * 6/8)/2

	for i, player := range app.Camaretto.Players {
		var theta float64 = angleStep * float64(i)
		var x float64 = centerX + (radius * math.Cos(theta + math.Pi/2))
		var y float64 = centerY + (radius * math.Sin(theta + math.Pi/2))
		player.Render(dst, x, y, theta)
	}
}

func (app *Application) DrawDeck(dst *ebiten.Image) {
	var centerX float64 = float64(WinWidth)/2
	var centerY float64 = (float64(WinHeight) * 6/8)/2
	
	var deck *Deck = app.Camaretto.DeckPile
	for i, card := range deck.DrawPile[:deck.LenDrawPile] {
		card.SSprite.ResetGeoM()
		card.SSprite.MoveImg(centerX - card.SSprite.Width, centerY - float64(i)*0.2)
		card.SSprite.Display(dst)
	}
	for i, card := range deck.DiscardPile[:deck.LenDiscardPile] {
		card.SSprite.ResetGeoM()
		card.SSprite.MoveImg(centerX, centerY - float64(i)*0.2)
		card.SSprite.Display(dst)
	}
}

func (app *Application) DrawButtons(dst *ebiten.Image) {
	var buttonXPos float64 = 0
	var buttonYPos float64 = float64(WinHeight)*9/10

	buttonXPos = float64(WinWidth)/2
	app.Info.SSprite.ResetGeoM()
	app.Info.SSprite.CenterImg()
	app.Info.SSprite.MoveImg(buttonXPos, buttonYPos - app.Info.SSprite.Height*2)
	app.Info.SSprite.Display(dst)

	buttonXPos = (float64(WinWidth) * 1/4) + (float64(ButtonWidth)/2)
	app.Attack.SSprite.ResetGeoM()
	app.Attack.SSprite.CenterImg()
	app.Attack.SSprite.MoveImg(buttonXPos, buttonYPos)
	app.Attack.SSprite.Display(dst)

	buttonXPos = (float64(WinWidth) * 2/4) + (float64(ButtonWidth)/2)
	app.Shield.SSprite.ResetGeoM()
	app.Shield.SSprite.CenterImg()
	app.Shield.SSprite.MoveImg(buttonXPos, buttonYPos)
	app.Shield.SSprite.Display(dst)

	buttonXPos = (float64(WinWidth) * 3/4) + (float64(ButtonWidth)/2)

	var playerTurn int = app.Camaretto.GetPlayerTurn()
	var p *Player = app.Camaretto.Players[playerTurn]
	if p.ChargeCard == nil {
		app.Heal.SSprite.ResetGeoM()

		app.Charge.SSprite.ResetGeoM()
		app.Charge.SSprite.CenterImg()
		app.Charge.SSprite.MoveImg(buttonXPos, buttonYPos)
		app.Charge.SSprite.Display(dst)
	} else {
		app.Charge.SSprite.ResetGeoM()

		app.Heal.SSprite.ResetGeoM()
		app.Heal.SSprite.CenterImg()
		app.Heal.SSprite.MoveImg(buttonXPos, buttonYPos)
		app.Heal.SSprite.Display(dst)
	}
}