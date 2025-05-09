package netplay

import (
	"log"

	"time"

	"net"
	"encoding/gob"

	"camaretto/model/game"
)

type ClientConnection struct {
	Connection *net.TCPConn
	Encoder *gob.Encoder
	Decoder *gob.Decoder
	Info *game.PlayerInfo
}

func NewClientConnection(c *net.TCPConn) *ClientConnection {
	var client *ClientConnection = &ClientConnection{}
	client.Connection = c
	client.Encoder = gob.NewEncoder(client.Connection)
	client.Decoder = gob.NewDecoder(client.Connection)
	client.Info = nil
	return client
}

type CamarettoServer struct {
	listener *net.TCPListener
	clients []*ClientConnection

	camaretto *game.Camaretto
}

// @desc: Create new instance of CamarettoServer then returns it
func NewCamarettoServer() *CamarettoServer {
	var err error

	var server *CamarettoServer = &CamarettoServer{}

	var addr *net.TCPAddr
	addr, err = net.ResolveTCPAddr("tcp", "localhost:58132")
	if err != nil { log.Fatal("[NewCamarettoServer] Unable to create ResolveTCPAddr:", err) }

	server.listener, err = net.ListenTCP("tcp", addr)
	if err != nil { log.Fatal("[NewCamarettoServer] Unable to create TCPListener:", err) }

	server.clients = []*ClientConnection{}

	server.camaretto = &game.Camaretto{}

	return server
}

// @desc:
func (server *CamarettoServer) handleError(e error, from string, action string) {
	var msg string = "[CamarettoServer." + from + "] " + action + ":"
	log.Println(msg, e)
}

// @desc:
func (server *CamarettoServer) Run() {
	log.Println("[CamarettoServer.Run] Lobby Routine begin")
	server.lobbyRoutine()
	log.Println("[CamarettoServer.Run] Lobby Routine end")
	log.Println("[CamarettoServer.Run] Game Routine begin")
	server.gameRoutine()
	log.Println("[CamarettoServer.Run] Game Routine end")

	var err error
	for _, client := range server.clients {
		err = client.Connection.Close()
		if err != nil {
			server.handleError(err, "Run", "Closing connection failed")
		}
	}

	err = server.listener.Close()
	if err != nil { server.handleError(err, "Run", "Closing listener failed") }
}

// @desc: Send a given message to every current server's connection
func (server *CamarettoServer) broadcastMessage(m *Message, except int) {
	var err error
	for _, conn := range server.clients {
		if conn.Info.Index != except {
			err = conn.Encoder.Encode(m)
			if err != nil {
				server.handleError(err, "broadcastMessage", "Broadcasting message failed")
			}
		}
	}
}

// @desc: Once a message in found is found in pipe channel sends it to all stored connections
// the routine is exited once a value is found in stop channel
func (server *CamarettoServer) broadcastRoutine(pipe chan *Message, stop chan bool) {
	for {
		var message *Message = nil
		select {
			case message = <-pipe:
				log.Println("[CamarettoServer.broadcastMessage] Broadcasting:", message)
				server.broadcastMessage(message, -1)
			case <-stop:
				log.Println("[CamarettoServer.broadcastRoutine] Routine stopped")
				return
			default:
		}
	}
}

// @desc: Handle first client connection, receiving player name then sending back player's index position
func (server *CamarettoServer) clientHandshake(conn *net.TCPConn) {
	var err error

	var client *ClientConnection = NewClientConnection(conn)
	server.clients = append(server.clients, client)

	var playerInfo *game.PlayerInfo = &game.PlayerInfo{}
	// Read player name
	err = client.Decoder.Decode(playerInfo)
	if err != nil {
		server.handleError(err, "clientHandshake", "Receive player name failed")
		return
	}

	// Send game index position to new player
	playerInfo.Index = len(server.clients) - 1

	err = client.Encoder.Encode(playerInfo)
	if err != nil {
		server.handleError(err, "clientHandshake", "Send player index failed")
		return
	}

	client.Info = playerInfo
	log.Println("[CamarettoServer.clientHandshake] Completed: {", playerInfo.Index, ",", playerInfo.Name, "}")
}

// @desc: Wait for a new connection to be opened, handshakes new client
// then trigger a broadcasting message once connection is complete
// the routine stops when a value is found in the stop channel
func (server *CamarettoServer) acceptConnections(pipe chan *Message, stop chan bool) {
	var err error
	var c *net.TCPConn

	for {
		select {
			case <-stop:
				log.Println("[CamarettoServer.acceptConnections] Routine stopped")
				return
			default:
				if len(server.clients) < game.MaxNbPlayers {
					server.listener.SetDeadline(time.Now().Add(5))
					c, err = server.listener.AcceptTCP()
					if err != nil && !err.(net.Error).Timeout() {
						server.handleError(err, "acceptConnections", "AcceptTCP failed")
					} else if err == nil {
						server.clientHandshake(c)

						var players []*game.PlayerInfo = []*game.PlayerInfo{}
						for _, client := range server.clients {
							players = append(players, client.Info)
						}
						pipe <- &Message{PLAYERS, -1, players, nil, -1, game.SET}
					}
				}
		}
	}

	// Timeout is not necessary anymore
	server.listener.SetDeadline(time.Time{})
}

// @desc: Handle new connections to server and update lobby state with all current connections
func (server *CamarettoServer) lobbyRoutine() {
	var stopBroadcast chan bool = make(chan bool)
	var stopAcceptConnections chan bool = make(chan bool)
	var pipe chan *Message = make(chan *Message)

	go server.broadcastRoutine(pipe, stopBroadcast)
	go server.acceptConnections(pipe, stopAcceptConnections)

	// Wait for first connection
	for ;len(server.clients) < 1; {}
	time.Sleep(time.Second * 2) // Wait for handshake to end

	var err error
	for { // Wait for host to send START message
		var msg *Message = &Message{}
		err = server.clients[0].Decoder.Decode(msg)
		if err != nil {
			log.Println(msg.Typ, msg.Seed, msg.Players, msg.Action)
			server.handleError(err, "lobbyRoutine", "Receive message from host failed")
		} else if msg.Typ == START {
			// Stop background routines
			stopBroadcast <- true
			stopAcceptConnections <- true
			time.Sleep(time.Second * 2) // Wait for routines to be over

			var seed int64 = time.Now().UnixNano()

			var names []string = make([]string, len(server.clients))
			var players []*game.PlayerInfo = []*game.PlayerInfo{}
			for _, client := range server.clients {
				names[client.Info.Index] = client.Info.Name
				players = append(players, client.Info)
			}

			server.broadcastMessage(&Message{INIT, seed, players, nil, -1, game.SET}, -1)
			server.camaretto.Init(seed, names, false, 0, 0)

			return // Exit lobbyRoutine
		} else {
			server.handleError(nil, "lobbyRoutine", "Received a message that should not have been sent")
		}
	}
}

// @desc:
func (server *CamarettoServer) getPlayerConnection(index int) *ClientConnection {
	for _, client := range server.clients {
		if client.Info.Index == index {
			return client
		}
	}
	return nil
}

// @desc:
func (server *CamarettoServer) gameRoutine() {
	var err error
	for ;!server.camaretto.IsGameOver(); {
		var msg *Message = &Message{}
		msg.Typ = ACTION
		msg.Action = server.camaretto.Current

		log.Println("[CamarettoServer.gameRoutine] 1.", server.camaretto.Current)
		server.broadcastMessage(msg, -1)

		var player int = -1

		// A player has to choose a card to defend on
		if server.camaretto.Current.Focus == game.CARD {
			player = server.camaretto.Current.PlayerFocus
		} else { // Player's turn
			player = server.camaretto.Current.PlayerTurn
		}

		var clientTurn *ClientConnection = server.getPlayerConnection(player)
		if clientTurn == nil {
			server.handleError(nil, "gameRoutine", "Client connection lost")
		}

		msg = &Message{}
		err = clientTurn.Decoder.Decode(msg)
		if err != nil {
			server.handleError(err, "gameRoutine", "Unable to decode client message")
			continue
		} else {
			log.Println("[CamarettoServer.gameRoutine] Message received:", msg)
		}

		if server.camaretto.Current.State == game.SET {
			server.camaretto.SetState(msg.State)
		} else {
			var index int = msg.Index

			var validIndex bool = true
			if server.camaretto.Current.Focus == game.PLAYER {
				if index < 0 || index > len(server.camaretto.Players)-1 {
					validIndex = false
				} else { server.camaretto.Current.PlayerFocus = index }
			} else if server.camaretto.Current.Focus == game.CARD {
				if index < 0 || index > 1 {
					validIndex = false
				} else { server.camaretto.Current.CardFocus = index }
			} else if server.camaretto.Current.Focus == game.REVEAL {
				if index < 0 || index > len(server.camaretto.Current.Reveal)-1 {
					validIndex = false
				} else { server.camaretto.Current.Reveal[index] = true }
			}

			if !validIndex { server.handleError(nil, "gameRoutine", "Invalid index received from player") }
		}

		log.Println("[CamarettoServer.gameRoutine] 2.", server.camaretto.Current)
		msg = &Message{ACTION, -1, nil, server.camaretto.Current, -1, game.SET};
		server.broadcastMessage(msg, -1)

		server.camaretto.Update(nil)

		if server.camaretto.Current.Focus == game.COMPLETE {
			log.Println("[CamarettoServer.gameRoutine] 3.", server.camaretto.Current)
			msg = &Message{ACTION, -1, nil, server.camaretto.Current, -1, game.SET};
			server.broadcastMessage(msg, -1)
			server.camaretto.Update(nil)
		}
	}
}
