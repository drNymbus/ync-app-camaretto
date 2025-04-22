package main

import (
	"log"

	"time"

	"net"

	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"camaretto/model/game"
	"camaretto/model/page"
	"camaretto/netplay"
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
	imgBuffer *ebiten.Image

	state AppState
	online, hosting bool

	playerInfo *game.PlayerInfo

	menu *page.Menu
	lobby *page.Lobby

	seed int64
	game *page.Game

	ioMessage chan *netplay.Message
	ioError chan error

	server *netplay.CamarettoServer
	client *netplay.CamarettoClient
}

func (app *Application) Init() {
	app.state = MENU
	app.online, app.hosting = false, false

	app.playerInfo = &game.PlayerInfo{}

	app.menu = &page.Menu{}
	app.menu.Init(WinWidth, WinHeight, app.startLobby, app.startServer, app.joinServer, app.scanServers)
	app.lobby = &page.Lobby{}
	app.game = &page.Game{}

	app.imgBuffer = ebiten.NewImage(WinWidth, WinHeight)
}

/************ ****************************************************************************** ************/
/************ ********************************** ROUTINE *********************************** ************/
/************ ****************************************************************************** ************/

func (app *Application) startLobby() {
	app.state = LOBBY

	app.online = app.menu.Online
	app.hosting = app.menu.Hosting
	app.playerInfo.Name = app.menu.Name.GetText()

	app.menu = &page.Menu{}
	var startFn func() = nil
	if app.online && app.hosting {
		startFn = app.serverStartGame
	} else if !app.online {
		startFn = app.startGame
	}
	app.lobby.Init(WinWidth, WinHeight, app.online, app.hosting, startFn)
}

func (app *Application) startGame() {
	app.state = GAME

	var playerNames []string = []string{}
	for i := 0; i < app.lobby.NbPlayers; i++ {
		playerNames = append(playerNames, app.lobby.Names[i].GetText())
	}

	app.lobby = &page.Lobby{}
	if !app.online { app.seed = time.Now().UnixNano() }
	app.game.Init(app.seed, playerNames, WinWidth, WinHeight, app.online, app.playerInfo, app.endGame)
}

func (app *Application) serverStartGame() {
	app.client.SendMessage(&netplay.Message{netplay.START, -1, nil, nil, nil})
}

func (app *Application) endGame() {
	app.state = END

	app.game = &page.Game{}
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
	addr, err = net.ResolveTCPAddr("tcp", "localhost:58132")
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

		if app.online {
			var message *netplay.Message
			var err error
			select {
				case message = <- app.ioMessage:
					log.Println("[Application.Update] Message Players")
					if message.Typ == netplay.PLAYERS { // New player
						app.lobby.NbPlayers = len(message.Players)
						for _, info := range message.Players {
							app.lobby.Names[info.Index].SetText(info.Name)
						}
					} else if message.Typ == netplay.INIT { // Game is starting
						log.Println("[Application.Update] Message INIT")
						app.lobby.NbPlayers = len(message.Players)
						for _, info := range message.Players {
							app.lobby.Names[info.Index].SetText(info.Name)
						}

						app.seed = message.Seed
						app.startGame()

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
		var old game.Action = *app.game.Camaretto.Current

		err = app.game.Update()
		if err != nil {
			log.Println("[Main.Update] Error update game:", err)
			return err
		}

		if app.online {
			if game.ActionDiff(&old, app.game.Camaretto.Current) {
				var msg *netplay.Message = &netplay.Message{}
				msg.Typ = netplay.ACTION
				msg.Action = app.game.Camaretto.Current
				msg.Reveal = []bool{}
				for _, card := range app.game.Camaretto.ToReveal {
					msg.Reveal = append(msg.Reveal, card.Hidden)
				}
				app.client.SendMessage(msg)
			}

			var message *netplay.Message
			var err error
			select {
				case message = <- app.ioMessage:
					if message.Typ == netplay.ACTION {
						log.Println("[Application.Update] Received new state", app.game.Camaretto.Current, message.Action)
						app.game.Camaretto.ApplyNewState(message.Action, message.Reveal)
						app.game.Camaretto.Update()
					}
					go app.client.ReceiveMessage(app.ioMessage, app.ioError)
				case err = <- app.ioError:
					log.Println("[Application.Update] Error:", err)
					// go app.client.ReceiveMessage(app.ioMessage, app.ioError)
				default: // Escape to continue to run program
			}
		}
	} else if app.state == END {
		app.startLobby()
	}

	return nil
}

func (app *Application) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)

	app.imgBuffer.Clear()
	app.imgBuffer.Fill(color.White)

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
