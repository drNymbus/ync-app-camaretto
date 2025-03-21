package model

import (
	"math"
	"strconv"

	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"camaretto/view"
	"camaretto/event"
)

const (
	WinWidth int = 1200
	WinHeight int = 900
	ButtonWidth int = WinWidth / 5
	ButtonHeight int = WinHeight / 6
)

type AppState int
const (
	MENU AppState = 0
	LOBBY AppState = 1
	GAME AppState = 2
	END AppState = 3
)

type Application struct{
	state AppState

	nbPlayers, maxNbPlayers int
	textCaptureWidth, textCaptureHeight int
	names []*TextCapture

	focus int
	cursor *view.Sprite

	xNb, yNb float64
	minusButton, plusButton *Button

	start *Button
	local, join, host *Button

	Camaretto *Camaretto

	imgBuffer *ebiten.Image
}

func (app *Application) Init(nbPlayers int) {
	app.state = MENU

	app.maxNbPlayers = 6
	app.nbPlayers = 1
	app.names = make([]*TextCapture, 6)

	app.textCaptureWidth, app.textCaptureHeight = WinWidth*3/4, WinHeight/10
	for i := 0; i < 6; i++ {
		app.names[i] = NewTextCapture(55, app.textCaptureWidth, app.textCaptureHeight, 2)
		var diffY float64 = float64((i - app.maxNbPlayers/2)*app.textCaptureHeight) + float64(i*10)
		app.names[i].SSprite.SetCenter(float64(WinWidth/2), float64(WinHeight/2) + 50 + diffY, 0)
	}

	app.focus = 0
	app.cursor = view.NewSprite(view.LoadCursorImage(), false, color.RGBA{0, 0, 0, 0}, nil)

	app.xNb, app.yNb = float64(WinWidth)/2, float64(WinHeight)/8
	app.minusButton = NewButton("-", color.RGBA{0, 0, 0, 255}, "RED")
	app.minusButton.SSprite.SetCenter(app.xNb - float64(view.ButtonWidth)/2 - 5, app.yNb, 0)
	app.plusButton = NewButton("+", color.RGBA{0, 0, 0, 255}, "RED")
	app.plusButton.SSprite.SetCenter(app.xNb + float64(view.ButtonWidth)/2 + 5, app.yNb, 0)

	app.start = NewButton("START", color.RGBA{0, 0, 0, 255}, "GREEN")
	app.start.SSprite.SetCenter(app.xNb, float64(WinHeight) - float64(view.ButtonHeight), 0)

	app.local = NewButton("Local", color.RGBA{0, 0, 0, 255}, "YELLOW")
	app.local.SSprite.SetCenter(app.xNb, float64(WinHeight/2) - float64(view.ButtonHeight) - 5, 0)
	app.host = NewButton("Host", color.RGBA{0, 0, 0, 255}, "YELLOW")
	app.host.SSprite.SetCenter(app.xNb, float64(WinHeight/2), 0)
	app.join = NewButton("Join", color.RGBA{0, 0, 0, 255}, "YELLOW")
	app.join.SSprite.SetCenter(app.xNb, float64(WinHeight/2) + float64(view.ButtonHeight) + 5, 0)

	app.Camaretto = &Camaretto{}
	// app.Camaretto.Init(nbPlayers, float64(WinWidth), float64(WinHeight))

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
	if app.state == MENU {
		if app.local.SSprite.In(x, y) {
			app.local.Pressed()
		} else if app.join.SSprite.In(x, y) {
			app.join.Pressed()
		} else if app.host.SSprite.In(x, y) {
			app.host.Pressed()
		}
	} else if app.state == LOBBY {
		if app.plusButton.SSprite.In(x, y) {
			app.plusButton.Pressed()
		} else if app.minusButton.SSprite.In(x, y) {
			app.minusButton.Pressed()
		} else if app.start.SSprite.In(x, y) {
			app.start.Pressed()
		} else {
			for i, textInput := range app.names {
				if textInput.SSprite.In(x, y) { app.focus = i }
			}
		}
	} else if app.state == END {
		app.state = MENU
	}
}

func (app *Application) mouseRelease(x, y float64) {
	if app.state == MENU {
		app.local.Released()
		app.host.Released()
		app.join.Released()
		if app.local.SSprite.In(x, y) {
			app.state = LOBBY
		}
	} else if app.state == LOBBY {
		app.plusButton.Released()
		app.minusButton.Released()
		app.start.Released()
		if app.plusButton.SSprite.In(x, y) {
			if app.nbPlayers < 6 { app.nbPlayers++ }
		} else if app.minusButton.SSprite.In(x, y) {
			if app.nbPlayers > 1 { app.nbPlayers-- }
		} else if app.start.SSprite.In(x, y) {
			app.state = GAME
			var playerNames []string = []string{}
			for i := 0; i < app.nbPlayers; i++ {
				playerNames = append(playerNames, app.names[i].GetText())
			}
			app.Camaretto.Init(app.nbPlayers, playerNames, float64(WinWidth), float64(WinHeight))
		}
	}
}

func (app *Application) MouseEventUpdate(e *event.MouseEvent) {
	if app.state == GAME {
		app.Camaretto.EventUpdate(e)
	} else {
		if e.Event == event.PRESSED {
			app.mousePress(e.X, e.Y)
		} else if e.Event == event.RELEASED {
			app.mouseRelease(e.X, e.Y)
		}
	}
}

func (app *Application) KeyEventUpdate(e *event.KeyEvent) {
	if app.state == LOBBY {
		if e.Event == event.PRESSED { app.names[app.focus].HandleEvent(e, nil) }
	}
}

func (app *Application) Update() {
	if app.state == MENU {
	} else if app.state == GAME {
		app.Camaretto.Update()
		if app.Camaretto.IsGameOver() { app.state = END }
	}
}

/************ *************************************************************************** ************/
/************ ********************************** RENDER ********************************* ************/
/************ *************************************************************************** ************/

func (app *Application) Display() *ebiten.Image {
	app.imgBuffer.Clear()
	app.imgBuffer.Fill(color.White)

	if app.state == MENU {
		app.local.SSprite.Display(app.imgBuffer)
		app.join.SSprite.Display(app.imgBuffer)
		app.host.SSprite.Display(app.imgBuffer)
	} else if app.state == LOBBY {
		app.minusButton.SSprite.Display(app.imgBuffer)
		app.plusButton.SSprite.Display(app.imgBuffer)
		app.start.SSprite.Display(app.imgBuffer)

		var x, y float64 = app.xNb, app.yNb - float64(view.ButtonHeight)/2
		var textImg *ebiten.Image
		var tw, th float64
		textImg, tw, th = view.TextToImage(strconv.Itoa(app.nbPlayers), color.RGBA{0, 0, 0, 255})
		op := &ebiten.DrawImageOptions{}; op.GeoM.Translate(x - tw/2, y - th)
		app.imgBuffer.DrawImage(textImg, op)

		for i := 0; i < app.nbPlayers; i++ {
			app.names[i].SSprite.Display(app.imgBuffer)
		}

		x = float64(WinWidth/2 - app.textCaptureWidth/2)
		y = float64(WinHeight/2 + 50) + float64((app.focus - app.maxNbPlayers/2)*app.textCaptureHeight) + float64(app.focus*10)
		app.cursor.Move(x, y, 1)
		app.cursor.Rotate(math.Pi/2, 1)
		app.cursor.Display(app.imgBuffer)

	} else if app.state == GAME {
		app.Camaretto.Render(app.imgBuffer, float64(WinWidth), float64(WinHeight))

	} else if app.state == END {
		img, tw, th := view.TextToImage("C'EST LA FIN!", color.RGBA{0, 0, 0, 255})
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(WinWidth/2) - tw/2, float64(WinHeight/2) - th/2)
		app.imgBuffer.DrawImage(img, op)
	}

	return app.imgBuffer
}