package model

import (
	// "log"
	// "math"
	// "strconv"
	"image/color"

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

	nbPlayers int
	names []string

	minusButton, plusButton *Button

	Camaretto *Camaretto

	imgBuffer *ebiten.Image
}

func (app *Application) Init(nbPlayers int) {
	app.state = GAME

	var x, y float64 = float64(WinWidth)/2, float64(WinHeight)/2 - 200
	app.minusButton = NewButton("-", color.RGBA{0, 0, 0, 255}, "RED")
	app.minusButton.SSprite.SetCenter(x - 100, y, 0)
	app.plusButton = NewButton("+", color.RGBA{0, 0, 0, 255}, "RED")
	app.plusButton.SSprite.SetCenter(x + 100, y, 0)

	app.Camaretto = NewCamaretto(nbPlayers, float64(WinWidth), float64(WinHeight))

	app.imgBuffer = ebiten.NewImage(WinWidth, WinHeight)
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

func (app *Application) mousePress(x, y float64) {
	if app.plusButton.SSprite.In(x, y) {
		app.plusButton.Pressed()
	} else if app.minusButton.SSprite.In(x, y) {
		app.minusButton.Pressed()
	}
}

func (app *Application) mouseRelease(x, y float64) {
	app.plusButton.Released()
	app.minusButton.Released()

	if app.plusButton.SSprite.In(x, y) {
		app.nbPlayers++
	} else if app.minusButton.SSprite.In(x, y) {
		app.nbPlayers--
	}
}

func (app *Application) MouseEventUpdate(e *event.MouseEvent) {
	if app.state == MENU {
		if e.Event == event.PRESSED {
			app.mousePress(e.X, e.Y)
		} else if e.Event == event.RELEASED {
			app.mouseRelease(e.X, e.Y)
		}
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

func (app *Application) Display() *ebiten.Image {
	app.imgBuffer.Clear()
	app.imgBuffer.Fill(color.White)

	if app.state == MENU {
		app.minusButton.SSprite.Display(app.imgBuffer)
		app.plusButton.SSprite.Display(app.imgBuffer)
	} else if app.state == GAME {
		app.Camaretto.Render(app.imgBuffer, float64(WinWidth), float64(WinHeight))
	}

	return app.imgBuffer
}