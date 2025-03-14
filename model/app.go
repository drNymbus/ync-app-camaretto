package model

import (
	// "log"
	// "math"
	// "strconv"
	// "image/color"

	"github.com/hajimehoshi/ebiten/v2"

	// "camaretto/view"
	"camaretto/event"
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
}

func (app *Application) Init(nbPlayers int) {
	app.state = GAME
	app.Camaretto = NewCamaretto(nbPlayers, float64(WinWidth), float64(WinHeight))
}

/************ ***************************************************************************** ************/
/************ ********************************** GET/SET ********************************** ************/
/************ ***************************************************************************** ************/

func (app *Application) SetState(s AppState) { app.state = s }
func (app *Application) GetState() AppState { return app.state }

/************ *************************************************************************** ************/
/************ ********************************* UPDATE ********************************** ************/
/************ *************************************************************************** ************/

func (app *Application) Hover(x, y float64) {
	if app.state == MENU {
	} else if app.state == GAME {
		app.Camaretto.mouseHover(x, y)
	}
}

func (app *Application) EventUpdate(e *event.MouseEvent) {
	if app.state == MENU {
	} else if app.state == GAME {
		app.Camaretto.EventUpdate(e)
	}
}

func (app *Application) Update() {
	if app.state == MENU {
	} else if app.state == GAME {
		app.Camaretto.Update()
	}
}

/************ *************************************************************************** ************/
/************ ********************************** RENDER ********************************* ************/
/************ *************************************************************************** ************/

func (app *Application) Display(dst *ebiten.Image) {
	if app.state == MENU {
	} else if app.state == GAME {
		app.Camaretto.Render(dst, float64(WinWidth), float64(WinHeight))
	}
}