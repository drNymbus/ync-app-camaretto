package model

import (
	"log"

	"net"
	"encoding/gob"
)

type ClientState int
const (
	SEND ClientState = 0
	RECEIVE ClientState = 1
)

type CamarettoClient struct {
	conn *net.TCPConn
	encoder *gob.Encoder
	decoder *gob.Decoder
}

type ServerState int
const (
	LISTENING ServerState = 0
	RUNNING ServerState = 1
)

type CamarettoServer struct {
	state ServerState
	listener *net.TCPListener
	clients []*CamarettoClient
}

// @desc: Create new instance of CamarettoServer then returns it
func NewCamarettoServer() *CamarettoServer {
	var err error

	var cs *CamarettoServer = &CamarettoServer{}
	cs.state = LISTENING

	var addr *net.TCPAddr
	addr, err = net.ResolveTCPAddr("tcp", "localhost:5813")
	if err != nil { log.Fatal("[NewCamarettoServer] Unable to create ResolveTCPAddr:", err) }

	cs.listener, err = net.ListenTCP("tcp", addr)
	if err != nil { log.Fatal("[NewCamarettoServer] Unable to create TCPListener:", err) }

	cs.clients = make([]*ClientConnection, MaxNbPlayers)
	return cs
}

func (s *CamarettoServer) Run(in, out chan Message) {}

// @desc: Accept incoming connection, adding any new connection to the player pool
func (s *CamarettoServer) AcceptConnection() (*net.TCPConn, error) {
	var err error
	var c *net.TCPConn

	c, err = s.listener.AcceptTCP()

	return c, err
}

// @desc: Handle first client connection, receiving player name then sending index position
func (s *CamarettoServer) HandleFirstConnection(conn *net.TCPConn) string {
	var err error

	s.clients = append(s.clients, conn)

	// Read player name
	var name string
	s.decoder = gob.NewDecoder(conn)
	err = s.decoder.Decode(name)
	if err != nil {
		log.Println("[HandleFirstConnection] Receive player name failed:", err)
	}

	// Send game index position to new player
	var index int = len(s.clients)
	s.encoder = gob.NewEncoder(conn)
	err = s.encoder.Encode(index)
	if err != nil {
		log.Println("[HandleFirstConnection] Send player index failed:", err)
	}

	return name
}

// @desc: Send player names and index position to every connection
func (s *CamarettoServer) BroadcastPlayerPool(names []string) {
	var err error

	var playerPool []PlayerInfo = []PlayerInfo{}
	for i, name := range names {
		playerPool = append(playerPool, PlayerInfo{i, name})
	}

	for _, conn := range s.clients {
		s.encoder = gob.NewEncoder(conn)
		err = s.encoder.Encode(playerPool)
		if err != nil {
			log.Println("[BroadcastPlayerPool] Sending player pool failed:", err)
		}
	}
}
