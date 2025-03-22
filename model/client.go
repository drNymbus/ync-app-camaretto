package model

import (
	"log"

	"net"
	"encoding/gob"

	// "camaretto/model/game"
)

type CamarettoClient struct {
	conn net.Conn
	decoder *gob.Decoder
	encoder *gob.Encoder
}

// @desc: Create new instance of CamarettoClient then returns it
func NewCamarettoClient() *CamarettoClient {
	var cc *CamarettoClient = &CamarettoClient{}
	return cc
}

func (cc *CamarettoClient) Run(input, output chan Message) {
}

// @desc: Retrieves every address open to connection
func (cc *CamarettoClient) Scan() []*net.TCPAddr {
	return nil
}

// @desc: Connect to server
func (cc *CamarettoClient) Connect(addr *net.TCPAddr) error {
	var err error
	cc.conn, err = net.DialTCP("tcp", nil, addr)
	if err != nil { log.Println("[Connect] Unable to dial: ", err) }
	return err
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
