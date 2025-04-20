package main

import (
	"log"

	"time"

	"net"

	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"camaretto/model"
	"camaretto/model/component"
	"camaretto/model/netplay"
	// "camaretto/event"
	"camaretto/view"
)

var (
	err error
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

type Application struct {
	// events *event.EventQueue
	imgBuffer *ebiten.Image

	state AppState
	online, hosting bool

	playerInfo *component.PlayerInfo

	menu *model.Menu
	lobby *model.Lobby
	game *model.Game

	ioMessage chan *netplay.Message
	ioError chan error

	server *netplay.CamarettoServer
	client *netplay.CamarettoClient
}

func (app *Application) Init() {
	// app.events = event.NewEventQueue(20)

	app.state = MENU
	app.online, app.hosting = false, false

	app.playerInfo = &component.PlayerInfo{}

	app.menu = &model.Menu{}
	app.menu.Init(WinWidth, WinHeight, app.startLobby, app.startServer, app.joinServer, app.scanServers)

	app.lobby = &model.Lobby{}
	// app.lobby.Init(WinWidth, WinHeight, app.online, app.hosting, app.startGame)

	app.game = &model.Game{}

	app.imgBuffer = ebiten.NewImage(WinWidth, WinHeight)
}

/************ ****************************************************************************** ************/
/************ ********************************** ROUTINE *********************************** ************/
/************ ****************************************************************************** ************/

func (app *Application) startLobby() {
	app.state = LOBBY
	app.menu = &model.Menu{}
	app.lobby.Init(WinWidth, WinHeight, app.menu.Online, app.menu.Hosting, app.startGame)
}

func (app *Application) startGame() {
	app.state = GAME

	var playerNames []string = []string{}
	for i := 0; i < app.lobby.NbPlayers; i++ {
		playerNames = append(playerNames, app.lobby.Names[i].GetText())
	}

	app.lobby = &model.Lobby{}
	var seed int64 = time.Now().UnixNano()
	app.game.Init(seed, playerNames, WinWidth, WinHeight, app.endGame)
}

func (app *Application) endGame() {
	app.state = END

	app.game = &model.Game{}
}

func (app *Application) startServer() {
	app.server = netplay.NewCamarettoServer()
	go app.server.Run()
	log.Println("SERVER LAUNCHED")
}

func (app *Application) joinServer() {
	var err error
	app.client = netplay.NewCamarettoClient()

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
		app.lobby.Names[app.playerInfo.Index].SetText(app.playerInfo.Name)
	}

	app.ioMessage = make(chan *netplay.Message)
	app.ioError = make(chan error)

	go app.client.ReceiveMessage(app.ioMessage, app.ioError)
}

func (app *Application) scanServers() {
}

/************ ***************************************************************************** ************/
/************ ********************************** EBITEN *********************************** ************/
/************ ***************************************************************************** ************/

func (app *Application) Update() error {
/*
	app.events.Update()

	if app.state == GAME {
		if !app.online || app.game.IsMyTurn(app.playerInfo.Index) {
			app.game.Hover(app.events.X, app.events.Y)
		}
	}

	var me *event.MouseEvent = nil
	var ke *event.KeyEvent = nil
	for ;!app.events.IsEmpty(); {
		me = app.events.ReadMouseEvent()
		if me != nil {
			var signal component.PageSignal = component.UPDATE
			if app.state == MENU {
				if me.Event == event.PRESSED {
					signal = app.menu.MousePress(me.X, me.Y)
				} else if me.Event == event.RELEASED {
					signal = app.menu.MouseRelease(me.X, me.Y)
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

					app.menu = &model.Menu{}
					app.state = LOBBY
				}
			} else if app.state == LOBBY {
				if me.Event == event.PRESSED {
					signal = app.lobby.MousePress(me.X, me.Y)
				} else if me.Event == event.RELEASED {
					signal = app.lobby.MouseRelease(me.X, me.Y)
				}

				if signal == component.NEXT {
					if app.online {
						app.client.SendMessage(&netplay.Message{netplay.START, -1, nil, nil, nil})
					} else {
						app.startCamaretto(time.Now().UnixNano())

						app.lobby = &model.Lobby{}
						app.state = GAME
					}
				}
			} else if app.state == GAME {
				if !app.online || app.game.IsMyTurn(app.playerInfo.Index) {
					if me.Event == event.PRESSED {
						app.game.MousePress(me.X, me.Y)
					} else if me.Event == event.RELEASED {
						app.game.MouseRelease(me.X, me.Y)
					}
				}
			} else if app.state == END {
				if me.Event == event.RELEASED {
					app.state = MENU
					app.menu.Init(WinWidth, WinHeight)
				}
			}
		}

		ke = app.events.ReadKeyEvent()
		if ke != nil {
			if app.state == MENU {
				app.menu.HandleKeyEvent(ke)
			} else if app.state == LOBBY {
				app.lobby.HandleKeyEvent(ke)
			}
		}
	}

	if app.state == LOBBY {
		if app.online {
			var message *netplay.Message
			var err error
			select {
				case message = <- app.ioMessage:
					if message.Typ == netplay.PLAYERS { // New player
						app.lobby.NbPlayers = len(message.Players)
						for _, info := range message.Players {
							app.lobby.Names[info.Index].SetText(info.Name)
						}
					} else if message.Typ == netplay.INIT { // Game is starting
						app.lobby.NbPlayers = len(message.Players)
						for _, info := range message.Players {
							app.lobby.Names[info.Index].SetText(info.Name)
						}
						app.startCamaretto(message.Seed)

						app.lobby = &model.Lobby{}
						app.state = GAME
					} else {
						log.Println("[Application.Update] Unparsable message (should not have been sent in the first place)")
					}
					go app.client.ReceiveMessage(app.ioMessage, app.ioError)
				case err = <- app.ioError:
					log.Println("[Application.Update]", err)
				default: // Escape to continue to run program
			}
		}
	} else if app.state == GAME {
		if app.online {
			var message *netplay.Message
			var err error
			select {
				case message = <- app.ioMessage:
					if message.Typ == netplay.ACTION {
						app.game.DeserializeCamaretto(message)
					}
				case err = <- app.ioError:
					log.Println("[Application.Update]", err)
				default: // Escape to continue to run program
			}
		}

		var signal component.PageSignal = app.game.Update()
		if signal == component.NEXT { app.state = END }
	} else if app.state == END {
	}
*/
	var err error

	if app.state == MENU {
		err = app.menu.Update()
		if err != nil {
			log.Println("[Main.Update] Error updating menu:", err)
			return err
		}
	} else if app.state == LOBBY {
		err = app.lobby.Update()
		if err != nil {
			log.Println("[Main.Update] Error updating lobby:", err)
			return err
		}
	} else if app.state == GAME {
		err = app.game.Update()
		if err != nil {
			log.Println("[Main.Update] Error update game:", err)
			return err
		}
	} else if app.state == END {
		app.state = LOBBY
		app.startLobby()
	}

	return nil
}

func (app *Application) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)

	app.imgBuffer.Clear()
	app.imgBuffer.Fill(color.White)

/*
	if app.state == MENU {
		app.menu.Display(app.imgBuffer)
	} else if app.state == LOBBY {
		app.lobby.Display(app.imgBuffer)
	} else if app.state == GAME {
		app.game.Display(app.imgBuffer)
	} else if app.state == END {
		img, tw, th := view.TextToImage("C'EST LA FIN!", color.RGBA{0, 0, 0, 255})
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(WinWidth/2) - tw/2, float64(WinHeight/2) - th/2)
		app.imgBuffer.DrawImage(img, op)
	}
*/

	if app.state == MENU {
		app.menu.Draw(app.imgBuffer)
	} else if app.state == LOBBY {
		app.lobby.Draw(app.imgBuffer)
	} else if app.state == GAME {
		app.game.Draw(app.imgBuffer)
	}

	screen.DrawImage(app.imgBuffer, nil)
}

func (app *Application) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return WinWidth, WinHeight
}

func main() {
	// Loading assets
	view.LoadFont()

	// Init App
	var app *Application = &Application{}
	app.Init()

	// Init Window
	ebiten.SetWindowSize(WinWidth, WinHeight)
	ebiten.SetWindowTitle("Camaretto")

	var icon image.Image
	icon, err = view.InitIcon("assets/amaretto_icon.png")
	if err != nil {
		log.Fatal("[MAIN] InitIcon failed", err)
	}
	ebiten.SetWindowIcon([]image.Image{icon})

	// Game Loop
	if err = ebiten.RunGame(app); err != nil {
		log.Fatal("[MAIN]", err)
	}

	// Free stuff
}
