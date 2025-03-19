package model

import (
	"net"

	"camaretto/model/game"
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
func (c *Client) Connect(address string) error {
}

// @desc: Send an action to the server
func (c *Client) SendAction(action []byte) error {
}

// @desc: Update game state from server's data received
func (c *Client) ReceiveGameState() {
}

// @desc: Disconnect from server
func (c *Client) Disconnect() error {
}
