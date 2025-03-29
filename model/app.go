package model

import (
	"log"

	"net"

	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"camaretto/model/game"
	"camaretto/view"
	"camaretto/event"
)

const (
	WinWidth int = 1200
	WinHeight int = 900
)

type AppState int
const (
	MENU AppState = iota
	LOBBY
	GAME
	END
)

func (a AppState) String() string {
	var name []string = []string{"MENU", "SCAN", "JOIN", "LOBBY", "GAME", "END"}
	return name[int(a)]
}

type Application struct{
	state AppState

	menu *game.Menu
	lobby *game.Lobby

	online, hosting bool

	PlayerInfo *PlayerInfo
	Camaretto *game.Camaretto

	ioMessage chan *Message
	ioError chan error

	server *CamarettoServer
	client *CamarettoClient

	imgBuffer *ebiten.Image
}

func (app *Application) Init(nbPlayers int) {
	app.state = MENU

	app.menu = &game.Menu{}
	app.menu.Init(WinWidth, WinHeight)

	app.lobby = &game.Lobby{}

	app.PlayerInfo = &PlayerInfo{}
	app.Camaretto = &game.Camaretto{}

	app.imgBuffer = ebiten.NewImage(WinWidth, WinHeight)
}

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
		app.lobby.Focus = app.PlayerInfo.Index
		app.lobby.Names[app.PlayerInfo.Index].SetText(app.PlayerInfo.Name)
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
	if app.state == GAME { app.Camaretto.Hover(x, y) }
}

func (app *Application) MouseEventUpdate(e *event.MouseEvent) {
	if app.state == MENU {
		if e.Event == event.PRESSED {
			app.menu.MousePress(e.X, e.Y)
		} else if e.Event == event.RELEASED {
			app.menu.MouseRelease(e.X, e.Y)
		}

		var playerName string = app.menu.Name.GetText()
		if len(playerName) > 0 {
			app.PlayerInfo.Name = playerName
			if app.menu.Hosting {
			} else {
			}
		}
	} else if app.state == LOBBY {
		if e.Event == event.PRESSED {
			app.lobby.MousePress(e.X, e.Y)
		} else if e.Event == event.RELEASED {
			app.lobby.MouseRelease(e.X, e.Y)
		}
	} else if app.state == GAME {
		app.Camaretto.EventUpdate(e)
	} else if app.state == END {
		app.state = MENU
	}
}

func (app *Application) KeyEventUpdate(e *event.KeyEvent) {
	if app.state == MENU {
		app.menu.HandleKeyEvent(e)
	} else if app.state == LOBBY {
		app.lobby.HandleKeyEvent(e)
	} else if app.state == GAME {
	} else if app.state == END {
	}
}

func (app *Application) Update() {
	if app.state == MENU {
	} else if app.state == LOBBY {
		if app.online {
			var message *Message
			var err error
			select {
				case message = <- app.ioMessage:
					if message.Typ == PLAYERS { // New player
						app.lobby.NbPlayers = len(message.Players)
						for _, info := range message.Players {
							app.lobby.Names[info.Index].SetText(info.Name)
						}
					} else if message.Typ == STATE { // Game is starting
						var playerNames []string = []string{}
						for i := 0; i < app.lobby.NbPlayers; i++ {
							playerNames = append(playerNames, app.lobby.Names[i].GetText())
						}

						app.Camaretto.Init(app.lobby.NbPlayers, playerNames, message.Game.Seed, float64(WinWidth), float64(WinHeight))
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
		app.menu.Display(app.imgBuffer)
	} else if app.state == LOBBY {
		app.lobby.Display(app.imgBuffer)
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
