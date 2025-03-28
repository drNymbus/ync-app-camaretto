package model

import (
	"log"

	"net"

	"math"
	"time"

	"strconv"

	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"camaretto/model/component"
	"camaretto/model/game"
	"camaretto/view"
	"camaretto/event"
)

const (
	WinWidth int = 1200
	WinHeight int = 900
	// ButtonWidth int = WinWidth / 5
	// ButtonHeight int = WinHeight / 6
	MaxNbPlayers int = 6
)

type AppState int
const (
	MENU AppState = iota
	SCAN
	JOIN
	LOBBY
	GAME
	END
)

type Application struct{
	state AppState

	nbPlayers int
	textCaptureWidth, textCaptureHeight int
	names []*component.TextCapture

	focus int
	cursor *view.Sprite

	xNb, yNb float64
	minusButton, plusButton *component.Button

	start *component.Button
	local, join, host *component.Button

	PlayerInfo *PlayerInfo
	Camaretto *game.Camaretto

	online bool
	hosting bool

	ioMessage chan *Message
	ioError chan error

	server *CamarettoServer
	client *CamarettoClient

	imgBuffer *ebiten.Image
}

func (app *Application) Init(nbPlayers int) {
	app.state = MENU

	app.nbPlayers = 2
	app.names = make([]*component.TextCapture, MaxNbPlayers)

	app.textCaptureWidth, app.textCaptureHeight = WinWidth*3/4, WinHeight/10
	for i := 0; i < MaxNbPlayers; i++ {
		app.names[i] = component.NewTextCapture(55, app.textCaptureWidth, app.textCaptureHeight, 2)
		var diffY float64 = float64((i - MaxNbPlayers/2)*app.textCaptureHeight) + float64(i*10)
		app.names[i].SSprite.SetCenter(float64(WinWidth/2), float64(WinHeight/2) + 50 + diffY, 0)
	}

	app.focus = 0
	app.cursor = view.NewSprite(view.LoadCursorImage(), false, color.RGBA{0, 0, 0, 0}, nil)

	app.xNb, app.yNb = float64(WinWidth)/2, float64(WinHeight)/8
	app.minusButton = component.NewButton("-", color.RGBA{0, 0, 0, 255}, "RED")
	app.minusButton.SSprite.SetCenter(app.xNb - float64(view.ButtonWidth)/2 - 5, app.yNb, 0)
	app.plusButton = component.NewButton("+", color.RGBA{0, 0, 0, 255}, "RED")
	app.plusButton.SSprite.SetCenter(app.xNb + float64(view.ButtonWidth)/2 + 5, app.yNb, 0)

	app.start = component.NewButton("START", color.RGBA{0, 0, 0, 255}, "GREEN")
	app.start.SSprite.SetCenter(app.xNb, float64(WinHeight) - float64(view.ButtonHeight), 0)

	app.local = component.NewButton("Local", color.RGBA{0, 0, 0, 255}, "GREEN")
	app.local.SSprite.SetCenter(app.xNb, float64(WinHeight/2) - float64(view.ButtonHeight) - 5, 0)
	app.host = component.NewButton("Host", color.RGBA{0, 0, 0, 255}, "BLUE")
	app.host.SSprite.SetCenter(app.xNb, float64(WinHeight/2), 0)
	app.join = component.NewButton("Join", color.RGBA{0, 0, 0, 255}, "RED")
	app.join.SSprite.SetCenter(app.xNb, float64(WinHeight/2) + float64(view.ButtonHeight) + 5, 0)

	app.PlayerInfo = &PlayerInfo{}
	app.Camaretto = &game.Camaretto{}

	app.online = false
	app.hosting = false

	app.imgBuffer = ebiten.NewImage(WinWidth, WinHeight)
}

/************ ***************************************************************************** ************/
/************ ********************************** GET/SET ********************************** ************/
/************ ***************************************************************************** ************/

func (app *Application) SetState(s AppState) { app.state = s }
func (app *Application) GetState() AppState { return app.state }

/************ ****************************************************************************** ************/
/************ ********************************** ROUTINE *********************************** ************/
/************ ****************************************************************************** ************/

func (app *Application) startServer() {
	app.server = NewCamarettoServer()
	go app.server.Run()

	log.Println("SERVER LAUNCHED")

	app.joinServer()
}

func (app *Application) joinServer() {
	var err error
	app.client = NewCamarettoClient()

	var addr *net.TCPAddr
	addr, err = net.ResolveTCPAddr("tcp", "localhost:5813")
	if err != nil {
		log.Println("[Application.joinServer] Unable to resolve host:", err)
		return
	}

	app.PlayerInfo, err = app.client.Connect(addr, app.PlayerInfo)
	if err != nil {
		log.Println("[ApplicationjoinServer] Connection failed:", err)
		return
	}

	if app.PlayerInfo != nil {
		app.focus = app.PlayerInfo.Index
		app.names[app.PlayerInfo.Index].SetText(app.PlayerInfo.Name)
	}

	app.ioMessage = make(chan *Message)
	app.ioError = make(chan error)

	go app.client.ReceiveMessage(app.ioMessage, app.ioError)
}

func (app *Application) scanServers() {
}

/************ ***************************************************************************** ************/
/************ ********************************** UPDATE *********************************** ************/
/************ ***************************************************************************** ************/

func (app *Application) Hover(x, y float64) {
	if app.state == MENU {
	} else if app.state == SCAN {
	} else if app.state == JOIN {
	} else if app.state == LOBBY {
	} else if app.state == GAME {
		app.Camaretto.Hover(x, y)
	} else if app.state == END {
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
	} else if app.state == SCAN {
	} else if app.state == JOIN {
		if app.hosting && app.host.SSprite.In(x, y) {
			app.host.Pressed()
		} else if app.join.SSprite.In(x, y) {
			app.join.Pressed()
		}
	} else if app.state == LOBBY {
		if app.plusButton.SSprite.In(x, y) {
			app.plusButton.Pressed()
		} else if app.minusButton.SSprite.In(x, y) {
			app.minusButton.Pressed()
		} else if app.start.SSprite.In(x, y) {
			app.start.Pressed()
		} else if !app.online {
			for i, textInput := range app.names {
				if textInput.SSprite.In(x, y) { app.focus = i }
			}
		}
	} else if app.state == GAME {
	} else if app.state == END {
		app.state = MENU
	}
}

func (app *Application) mouseRelease(x, y float64) {
	var err error

	if app.state == MENU {
		app.local.Released()
		app.host.Released()
		app.join.Released()
		if app.local.SSprite.In(x, y) {
			app.online = false
			app.state = LOBBY

		} else if app.host.SSprite.In(x, y) {
			app.online = true
			app.hosting = true

			app.names[0].SSprite.SetCenter(float64(WinWidth)/2, float64(WinHeight)/2 - float64(app.textCaptureHeight)/2, 0)
			app.host.SSprite.Move(float64(WinWidth)/2, float64(WinHeight)/2 + float64(view.ButtonHeight)*2, 0.1)

			app.state = JOIN

		} else if app.join.SSprite.In(x, y) {
			app.online = true
			app.hosting = false

			app.names[0].SSprite.SetCenter(float64(WinWidth)/2, float64(WinHeight)/2 - float64(app.textCaptureHeight)/2, 0)
			app.join.SSprite.Move(float64(WinWidth)/2, float64(WinHeight)/2 + float64(view.ButtonHeight)*2, 0.1)

			app.state = JOIN
		}
	} else if app.state == SCAN {
	} else if app.state == JOIN {
		app.host.Released()
		app.join.Released()
		var playerName string = app.names[0].GetText()
		if len(playerName) > 0 {
			app.PlayerInfo.Name = playerName
			if app.hosting {
				if app.host.SSprite.In(x, y) {
					app.startServer()
					app.state = LOBBY
				}
			} else {
				if app.join.SSprite.In(x, y) {
					app.joinServer()
					app.state = LOBBY
				}
			}
		}
	} else if app.state == LOBBY {
		if !app.online {
			app.plusButton.Released()
			app.minusButton.Released()
			app.start.Released()
			if app.plusButton.SSprite.In(x, y) {
				if app.nbPlayers < 6 { app.nbPlayers++ }
			} else if app.minusButton.SSprite.In(x, y) {
				if app.nbPlayers > 2 { app.nbPlayers-- }
			} else if app.start.SSprite.In(x, y) {
				var playerNames []string = []string{}
				for i := 0; i < app.nbPlayers; i++ {
					playerNames = append(playerNames, app.names[i].GetText())
				}

				app.Camaretto.Init(app.nbPlayers, playerNames, time.Now().UnixNano(), float64(WinWidth), float64(WinHeight))
				app.state = GAME
			}
		} else if app.hosting {
			app.start.Released()
			if app.start.SSprite.In(x, y) {
				var msg *Message = &Message{START, nil, nil}
				for err = app.client.SendMessage(msg); err != nil; {
					err = app.client.SendMessage(msg)
				}
			}
		}
	} else if app.state == GAME {
	} else if app.state == END {
	}
}

func (app *Application) MouseEventUpdate(e *event.MouseEvent) {
	if app.state == MENU || app.state == LOBBY || app.state == JOIN {
		if e.Event == event.PRESSED {
			app.mousePress(e.X, e.Y)
		} else if e.Event == event.RELEASED {
			app.mouseRelease(e.X, e.Y)
		}
	} else if app.state == SCAN {
	} else if app.state == GAME {
		app.Camaretto.EventUpdate(e)
	} else if app.state == END {
	}
}

func (app *Application) KeyEventUpdate(e *event.KeyEvent) {
	if app.state == MENU {
	} else if app.state == SCAN {
	} else if app.state == JOIN {
		if e.Event == event.PRESSED {
			app.names[0].HandleEvent(e, nil)
		}
	} else if app.state == LOBBY && !app.online {
		if e.Event == event.PRESSED {
			app.names[app.focus].HandleEvent(e, nil)
		}
	} else if app.state == GAME {
	} else if app.state == END {
	}
}

func (app *Application) Update() {
	if app.state == MENU {
	} else if app.state == SCAN {
	} else if app.state == JOIN {
	} else if app.state == LOBBY {
		if app.online {
			var message *Message
			var err error
			select {
				case message = <- app.ioMessage:
					if message.Typ == PLAYERS { // New player
						app.nbPlayers = len(message.Players)
						for _, info := range message.Players {
							app.names[info.Index].SetText(info.Name)
						}
					} else if message.Typ == STATE { // Game is starting
						var playerNames []string = []string{}
						for i := 0; i < app.nbPlayers; i++ {
							playerNames = append(playerNames, app.names[i].GetText())
						}

						app.Camaretto.Init(app.nbPlayers, playerNames, message.Game.Seed, float64(WinWidth), float64(WinHeight))
						app.state = GAME
					}
					go app.client.ReceiveMessage(app.ioMessage, app.ioError)
				case err = <- app.ioError:
					log.Println("[Application.Update]", err)
				default: // Escape to continue to run program
			}
		}
	} else if app.state == GAME {
		app.Camaretto.Update()
		if app.Camaretto.IsGameOver() { app.state = END }
	} else if app.state == END {
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
	} else if app.state == SCAN {
	} else if app.state == JOIN {
		app.names[0].SSprite.Display(app.imgBuffer)
		if app.hosting {
			app.host.SSprite.Display(app.imgBuffer)
		} else {
			app.join.SSprite.Display(app.imgBuffer)
		}
	} else if app.state == LOBBY {
		var x, y float64

		if !app.online {
			app.minusButton.SSprite.Display(app.imgBuffer)
			app.plusButton.SSprite.Display(app.imgBuffer)
		}

		if !app.online || app.hosting {
			app.start.SSprite.Display(app.imgBuffer)
		}

		if !app.online {
			x, y = app.xNb, app.yNb - float64(view.ButtonHeight)/2
			var textImg *ebiten.Image
			var tw, th float64
			textImg, tw, th = view.TextToImage(strconv.Itoa(app.nbPlayers), color.RGBA{0, 0, 0, 255})
			op := &ebiten.DrawImageOptions{}; op.GeoM.Translate(x - tw/2, y - th)
			app.imgBuffer.DrawImage(textImg, op)
		}

		for i := 0; i < app.nbPlayers; i++ {
			var diffY float64 = float64((i - MaxNbPlayers/2)*app.textCaptureHeight) + float64(i*10)
			app.names[i].SSprite.SetCenter(float64(WinWidth/2), float64(WinHeight/2) + 50 + diffY, 0)
			app.names[i].SSprite.Display(app.imgBuffer)
		}

		x = float64(WinWidth/2 - app.textCaptureWidth/2)
		y = float64(WinHeight/2 + 50) + float64((app.focus - MaxNbPlayers/2)*app.textCaptureHeight) + float64(app.focus*10)
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
