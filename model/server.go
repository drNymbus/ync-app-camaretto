package model

import (
	"net"

	"camaretto/model/game"
)

type CamarettoServer struct {
	listener net.Listener
	clients map[string]net.Conn
	game *game.Camaretto
}

// @desc: Create new instance of CamarettoServer then returns it
func NewCamarettoServer() *CamarettoServer {
	var cs *CamarettoServer = &CamarettoServer{}
	return cs
}

// @desc: Start listening on specified port
func (s *CamarettoServer) Start() error {
	return nil
}

// @desc: Accept incoming connections
func (s *CamarettoServer) AcceptConnections() {
}

// @desc: Handle connected client
func (s *CamarettoServer) HandleClient(conn net.Conn) {
}

// @desc: Send game state to every connected client
func (s *CamarettoServer) BroadcastGameState() error {
	return nil
}

// @desc: Stop the server and close all connections with current clients
func (s *CamarettoServer) Stop() error {
	return nil
}
