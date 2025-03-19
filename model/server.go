package model

import (
	"net"

	"camaretto/model/game"
)

type CamarettoServer struct {
	listener net.Listener
	clients map[string]net.Conn
	game *Camaretto
}

// @desc: Create new instance of CamarettoServer then returns it
func NewCamarettoServer() *CamarettoServer {
	var cs *CamarettoServer = &CamarettoServer{}
	return cs
}

// @desc: Start listening on specified port
func (s *Server) Start() error {
}

// @desc: Accept incoming connections
func (s *Server) AcceptConnections() {
}

// @desc: Handle connected client
func (s *Server) HandleClient(conn net.Conn) {
}

// @desc: Send game state to every connected client
func (s *Server) BroadcastGameState() error {
}

// @desc: Stop the server and close all connections with current clients
func (s *Server) Stop() error {
}
