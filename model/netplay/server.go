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

	seed int64
	action *game.Action
	reveal []bool
}

// @desc: Create new instance of CamarettoServer then returns it
func NewCamarettoServer() *CamarettoServer {
	var err error

	var server *CamarettoServer = &CamarettoServer{}

	var addr *net.TCPAddr
	addr, err = net.ResolveTCPAddr("tcp", "localhost:5813")
	if err != nil { log.Fatal("[NewCamarettoServer] Unable to create ResolveTCPAddr:", err) }

	server.listener, err = net.ListenTCP("tcp", addr)
	if err != nil { log.Fatal("[NewCamarettoServer] Unable to create TCPListener:", err) }

	server.clients = []*ClientConnection{}

	return server
}

// @desc:
func (server *CamarettoServer) handleError(e error, from string, action string) {
	var msg string = "[CamarettoServer." + from + "] " + action + ":"
	log.Println(msg, e)
}

// @desc:
func (server *CamarettoServer) Run() {
	server.lobbyRoutine()
	server.gameRoutine()

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
func (server *CamarettoServer) broadcastMessage(m *Message) {
	log.Println("[CamarettoServer.broadcastMessage] Broadcasting:", m)
	var err error
	for _, conn := range server.clients {
		err = conn.Encoder.Encode(m)
		if err != nil {
			server.handleError(err, "broadcastMessage", "Broadcasting message failed")
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
				server.broadcastMessage(message)
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
	}

	// Send game index position to new player
	playerInfo.Index = len(server.clients)

	err = client.Encoder.Encode(playerInfo)
	if err != nil {
		server.handleError(err, "clientHandshake", "Send player index failed")
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
					c, err = server.listener.AcceptTCP()
					if err != nil {
						server.handleError(err, "acceptConnections", "AcceptTCP failed")
					}

					server.clientHandshake(c)

					var players []*game.PlayerInfo = []*game.PlayerInfo{}
					for _, client := range server.clients {
						players = append(players, client.Info)
					}
					pipe <- &Message{PLAYERS, -1, players, nil, nil}
				}
		}
	}
}

// @desc: Handle new connections to server and update lobby state with all current connections
func (server *CamarettoServer) lobbyRoutine() {
	var stop chan bool = make(chan bool)
	var pipe chan *Message = make(chan *Message)

	go server.broadcastRoutine(pipe, stop)
	go server.acceptConnections(pipe, stop)

	// Wait for first connection
	for ;len(server.clients) < 1; {}
	log.Println(len(server.clients), "CLIENTS")

	var err error
	for { // Wait for host to send START message
		var msg *Message = &Message{}
		err = server.clients[0].Decoder.Decode(msg)
		if err != nil {
			log.Println(msg.Typ, msg.Seed, msg.Players, msg.Action)
			server.handleError(err, "lobbyRoutine", "Receive message from host failed")
		} else if msg.Typ == START {
			stop <- true // Stop background routines
			time.Sleep(time.Second * 2) // Wait for routines to be over

			server.seed = time.Now().UnixNano()

			var players []*game.PlayerInfo = []*game.PlayerInfo{}
			for _, client := range server.clients {
				players = append(players, client.Info)
			}
			server.broadcastMessage(&Message{INIT, server.seed, players, nil, nil})

			return // Exit lobbyRoutine
		} else {
			server.handleError(nil, "lobbyRoutine", "Received a message that should not have been sent")
		}
	}
}

// @desc:
func (server *CamarettoServer) gameRoutine() {
	for {
	}
}
