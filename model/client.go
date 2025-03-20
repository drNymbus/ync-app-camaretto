package model

import (
	"net"

	// "camaretto/model/game"
)


type CamarettoClient struct {
	conn net.Conn
	updates chan []byte
}

// @desc: Create new instance of CamarettoClient then returns it
func NewCamarettoClient() *CamarettoClient {
	var cc *CamarettoClient = &CamarettoClient{}
	return cc
}

// @desc: Connect to server
func (cc *CamarettoClient) Connect(address string) error {
	return nil
}

// @desc: Send an action to the server
func (cc *CamarettoClient) SendAction(action []byte) error {
	return nil
}

// @desc: Update game state from server's data received
func (cc *CamarettoClient) ReceiveGameState() {
}

// @desc: Disconnect from server
func (cc *CamarettoClient) Disconnect() error {
	return nil
}
