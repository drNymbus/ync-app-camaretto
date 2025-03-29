package model

import (
	"log"

	"net"

	"camaretto/model/game"
)

type PageSignal int
const (
	PREVIOUS PageSignal = iota
	NEXT
	SETTING // This is to open the settings menu (thoughful future implementation that might be deleted later)
	UPDATE // This is basically a NONE statement, nothing has happened and nothing should be done.
)

type CamarettoServer struct {
	listener *net.TCPListener
	clients []*ClientConnection

	camaretto *CamarettoState
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

	server.camaretto = NewCamarettoState()

	return server
}

// @desc:
func (server *CamarettoServer) handleError(e error, from string, action string) {
	var msg string = "[CamarettoServer." + from + "] " + action + ":"
	log.Println(msg, e)
}

func (server *CamarettoServer) Run() {
	server.lobbyRoutine()
	server.gameRoutine()
}

// @desc: Send a given message to every current server's connection
func (server *CamarettoServer) broadcastMessage(m *Message) {
	var err error
	for _, conn := range server.clients {
		err = conn.Encoder.Encode(m)
		if err != nil {
			server.handleError(err, "broadcastMessage", "Broadcasting message failed")
		}
	}
}

// @desc:
func (server *CamarettoServer) broadcastRoutine(pipe chan *Message) {
	for {
		var message *Message = nil
		select {
			case message = <-pipe:
				server.broadcastMessage(message)
			default:
		}
	}
}

// @desc: Handle first client connection, receiving player name then sending back player's index position
func (server *CamarettoServer) clientHandshake(conn *net.TCPConn) {
	var err error

	var client *ClientConnection = NewClientConnection(conn)
	server.clients = append(server.clients, client)

	var playerInfo *PlayerInfo = &PlayerInfo{}
	// Read player name
	err = client.Decoder.Decode(playerInfo)
	if err != nil {
		server.handleError(err, "clientHandshake", "Receive player name failed")
	}

	// Send game index position to new player
	log.Println(len(server.camaretto.Players))
	playerInfo.Index = len(server.camaretto.Players)
	server.camaretto.Players = append(server.camaretto.Players, playerInfo)

	err = client.Encoder.Encode(playerInfo)
	if err != nil {
		server.handleError(err, "clientHandshake", "Send player index failed")
	}

	log.Println("[CamarettoServer.clientHandshake] Completed: {", playerInfo.Index, ",", playerInfo.Name, "}")
}

// @desc:
func (server *CamarettoServer) acceptConnections(pipe chan *Message) {
	var err error
	var c *net.TCPConn

	for {
		c, err = server.listener.AcceptTCP()
		if err != nil {
			server.handleError(err, "acceptConnections", "AcceptTCP failed")
		}

		server.clientHandshake(c)
		if len(server.clients) == game.MaxNbPlayers { return }
		log.Println("[CamarettoServer.acceptConnections]", server.camaretto.toString())
		pipe <- &Message{PLAYERS, server.camaretto.Players, nil}
	}
}

// @desc:
func (server *CamarettoServer) lobbyRoutine() {
	var pipe chan *Message = make(chan *Message)
	go server.broadcastRoutine(pipe)
	go server.acceptConnections(pipe)

	// Wait for first connection
	for ;len(server.clients) < 1; {}
	log.Println(len(server.clients), "CLIENTS")

	var err error
	for {
		var msg *Message = &Message{}
		err = server.clients[0].Decoder.Decode(msg)
		if err != nil {
			log.Println(msg.Typ, msg.Players, msg.Game)
			server.handleError(err, "lobbyRoutine", "Receive message from host failed")
		} else if msg.Typ == START {
			server.broadcastMessage(&Message{PLAYERS, server.camaretto.Players, nil})
			server.broadcastMessage(&Message{STATE, nil, server.camaretto})
			return
		} else {
			server.handleError(nil, "lobbyRoutine", "Received a message that should not have been sent")
		}
	}
}

// @desc:
func (server *CamarettoServer) gameRoutine() {
}
