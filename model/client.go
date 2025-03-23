package model

import (
	"log"

	"net"
	"encoding/gob"

	// "camaretto/model/game"
)

type ClientConnection struct {
	Conn *net.TCPConn
	Encoder *gob.Encoder
	Decoder *gob.Decoder
}

func NewClientConnection(c *net.TCPConn) *ClientConnection {
	var cc *ClientConnection = &ClientConnection{}
	cc.Conn = c
	cc.Encoder = gob.NewEncoder(cc.Conn)
	cc.Decoder = gob.NewDecoder(cc.Conn)
	return cc
}

type CamarettoClient struct {
	Connection *ClientConnection
	msg *Message
}

// @desc: Create new instance of CamarettoClient then returns it
func NewCamarettoClient() *CamarettoClient {
	var cc *CamarettoClient = &CamarettoClient{}
	return cc
}

// @desc: Connect to server
func (cc *CamarettoClient) Connect(addr *net.TCPAddr) {
	var err error
	var c *net.TCPConn
	c, err = net.DialTCP("tcp", nil, addr)
	if err == nil {
		cc.Connection = NewClientConnection(c)
		cc.msg = &Message{}
	} else { log.Println("[Connect] Unable to dial: ", err) }

	err = cc.Connection.Encoder.Encode(&Message{HANDSHAKE, &PlayerInfo{-1, "MARIO"}, nil, nil})
	if err != nil {
		log.Println("[CamarettoClient.Connect] Unable to send player info:", err)
	}

	err = cc.Connection.Decoder.Decode(cc.msg)
	if err != nil {
		log.Println("[CamarettoClient.Connect] Unable to receive player info:", err)
	}
}

func (cc *CamarettoClient) Run(input, output chan Message) {
}

// @desc: Retrieves every address open to connection
func (cc *CamarettoClient) Scan() []*net.TCPAddr {
	return nil
}

// @desc: Disconnect from server
func (cc *CamarettoClient) Disconnect() error {
	return nil
}
