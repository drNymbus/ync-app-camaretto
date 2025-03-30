package model

import (
	"log"

	"time"

	"net"

	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"camaretto/model/game"
	"camaretto/model/component"
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

type Application struct{
	state AppState
	online, hosting bool

	menu *game.Menu
	lobby *game.Lobby

	playerInfo *PlayerInfo
	camaretto *game.Camaretto

	ioMessage chan *Message
	ioError chan error

	server *CamarettoServer
	client *CamarettoClient

	imgBuffer *ebiten.Image
}

func (app *Application) Init(nbPlayers int) {
	app.state = MENU
	app.online, app.hosting = false, false

	app.menu = &game.Menu{}
	app.menu.Init(WinWidth, WinHeight)

	app.lobby = &game.Lobby{}

	app.playerInfo = &PlayerInfo{-1, ""}
	app.camaretto = &game.Camaretto{}

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

	app.playerInfo, err = app.client.Connect(addr, app.playerInfo)
	if err != nil {
		log.Println("[ApplicationjoinServer] Connection failed:", err)
		return
	}

	if app.playerInfo != nil {
		app.lobby.Focus = app.playerInfo.Index
		app.lobby.Names[app.playerInfo.Index].SetText(app.playerInfo.Name)
	}

	app.ioMessage = make(chan *Message)
	app.ioError = make(chan error)

	go app.client.ReceiveMessage(app.ioMessage, app.ioError)
}

func (app *Application) scanServers() {
}

func (app *Application) startCamaretto(seed int64) {
	var playerNames []string = []string{}
	for i := 0; i < app.lobby.NbPlayers; i++ {
		playerNames = append(playerNames, app.lobby.Names[i].GetText())
	}

	app.camaretto.Init(app.lobby.NbPlayers, playerNames, seed, WinWidth, WinHeight)
}

/************ ***************************************************************************** ************/
/************ ********************************** UPDATE *********************************** ************/
/************ ***************************************************************************** ************/

// @desc: true if the player (Application.PlayerInfo) is required to do an action, false otherwise
func (app *Application) isMyTurn() bool {
	var flag bool = false

	if app.camaretto.State == SET {
		flag = (app.camaretto.PlayerTurn == app.PlayerInfo.Index)
	} else {
		if app.camaretto.Focus == PLAYER || app.camaretto.Focus == REVEAL {
			flag = (app.camaretto.PlayerTurn == app.PlayerInfo.Index)
		} else if app.camaretto.Focus == CARD {
			flag = (app.camaretto.PlayerFocus == app.PlayerInfo.Index)
		}
	}

	return flag
}

func (app *Application) Hover(x, y float64) {
	if app.state == GAME {
		if !app.online || app.isMyTurn() {
			app.camaretto.Hover(x, y)
		}
	}
}

func (app *Application) MouseEventUpdate(e *event.MouseEvent) {
	var signal component.PageSignal = component.UPDATE
	if app.state == MENU {
		if e.Event == event.PRESSED {
			signal = app.menu.MousePress(e.X, e.Y)
		} else if e.Event == event.RELEASED {
			signal = app.menu.MouseRelease(e.X, e.Y)
		}

		if signal == component.NEXT {
			app.online, app.hosting = app.menu.Online, app.menu.Hosting
			app.lobby.Init(WinWidth, WinHeight, app.online, app.hosting)

			if app.online {
				app.playerInfo.Name = app.menu.Name.GetText()
				if app.hosting {
					app.startServer()
				} else if app.online {
					app.joinServer()
				}
			}

			app.menu = &game.Menu{}
			app.state = LOBBY
		}
	} else if app.state == LOBBY {
		if e.Event == event.PRESSED {
			signal = app.lobby.MousePress(e.X, e.Y)
		} else if e.Event == event.RELEASED {
			signal = app.lobby.MouseRelease(e.X, e.Y)
		}

		if signal == component.NEXT {
			if app.online {
				app.client.SendMessage(&Message{START, nil, nil})
			} else {
				app.startCamaretto(time.Now().UnixNano())
				app.state = GAME
			}
		}
	} else if app.state == GAME {
		if !app.online || app.isMyTurn() {
			if e.Event == event.PRESSED {
				app.camaretto.MousePress(e.X, e.Y)
			} else if e.Event == event.RELEASED {
				app.camaretto.MouseRelease(e.X, e.Y)
			}
		}
	} else if app.state == END {
		if e.Event == event.RELEASED {
			app.state = MENU
			app.menu.Init(WinWidth, WinHeight)
		}
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
	if app.state == LOBBY {
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
						app.startCamaretto(message.Game.Seed)
						app.state = GAME
					}
					go app.client.ReceiveMessage(app.ioMessage, app.ioError)
				case err = <- app.ioError:
					log.Println("[Application.Update]", err)
				default: // Escape to continue to run program
			}
		}
	} else if app.state == GAME {
		if app.online {
			var message *Message
			var err error
			select {
				case message = <- app.ioMessage:
					if message.Typ == STATE {
						var state *CamarettoState = message.Game
						app.camaretto.State = state.Game
						app.camaretto.Focus = state.Focus
						app.camaretto.PlayerTurn = state.Turn
						app.camaretto.PlayerFocus = state.Player
						app.camaretto.CardFocus = state.Card
						
						for i, revealed := range state.Reveal {
							if revealed { app.camaretto.ToReveal[i].Reveal() }
						}
					}
				case err = <- app.ioError:
					log.Println("[Application.Update]", err)
				default: // Escape to continue to run program
			}
		}

		app.camaretto.Update()
		if app.camaretto.IsGameOver() { app.state = END }
	} else if app.state == END {
		app.client.Disconnect()
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
		app.camaretto.Display(app.imgBuffer)
	} else if app.state == END {
		img, tw, th := view.TextToImage("C'EST LA FIN!", color.RGBA{0, 0, 0, 255})
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(WinWidth/2) - tw/2, float64(WinHeight/2) - th/2)
		app.imgBuffer.DrawImage(img, op)
	}

	return app.imgBuffer
}
