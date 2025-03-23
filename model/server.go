package model

import (
	"log"

	"net"
	// "encoding/gob"
)

type CamarettoServer struct {
	state AppState
	listener *net.TCPListener
	clients []*ClientConnection
}

// @desc: Create new instance of CamarettoServer then returns it
func NewCamarettoServer() *CamarettoServer {
	var err error

	var cs *CamarettoServer = &CamarettoServer{}
	cs.state = LOBBY

	var addr *net.TCPAddr
	addr, err = net.ResolveTCPAddr("tcp", "localhost:5813")
	if err != nil { log.Fatal("[NewCamarettoServer] Unable to create ResolveTCPAddr:", err) }

	cs.listener, err = net.ListenTCP("tcp", addr)
	if err != nil { log.Fatal("[NewCamarettoServer] Unable to create TCPListener:", err) }

	cs.clients = make([]*ClientConnection, MaxNbPlayers)
	return cs
}

func (s *CamarettoServer) Run(input, output chan *Message) {
	var err error
	if s.state == LOBBY {
		for {
			var c *net.TCPConn
			c, err = s.listener.AcceptTCP()
			if err != nil {
				log.Println("[CamarettoServer.Run - LOBBY] AcceptTCP failed:", err)
			}

			go s.ClientHandshake(c, output)
		}
	} else if s.state == GAME {
	}
}

// @desc: Handle first client connection, receiving player name then sending back player's index position
func (s *CamarettoServer) ClientHandshake(conn *net.TCPConn, output chan *Message) {
	var err error
	var msg *Message

	var client *ClientConnection = NewClientConnection(conn)
	s.clients = append(s.clients, client)

	// Read player name
	// var name string
	// s.decoder = gob.NewDecoder(conn)
	err = client.Decoder.Decode(msg)
	if err == nil {
		output <- msg
	} else {
		log.Println("[HandleFirstConnection] Receive player name failed:", err)
	}

	// Send game index position to new player
	// var index int = len(s.clients)
	msg = &Message{HANDSHAKE, &PlayerInfo{len(s.clients), ""}, nil, nil}
	// s.encoder = gob.NewEncoder(conn)
	err = client.Encoder.Encode(msg)
	if err != nil {
		log.Println("[HandleFirstConnection] Send player index failed:", err)
	}
}

// @desc: Send player names and index position to every connection
func (s *CamarettoServer) BroadcastPlayerPool() {
}
