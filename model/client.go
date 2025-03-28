package model

import (
	"log"

	"net"
	"encoding/gob"
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
}

// @desc: Create new instance of CamarettoClient then returns it
func NewCamarettoClient() *CamarettoClient {
	var client *CamarettoClient = &CamarettoClient{}
	return client
}

// @desc:
func (client *CamarettoClient) handleError(e error, from string, action string) {
	var msg string = "[CamarettoClient." + from + "] " + action + ":"
	log.Println(msg, e)
}

// @desc: Retrieves every address open to connection
func (client *CamarettoClient) Scan() []*net.TCPAddr {
	return nil
}

// @desc: Connect to server
func (client *CamarettoClient) Connect(addr *net.TCPAddr, info *PlayerInfo) (*PlayerInfo, error) {
	var err error

	var c *net.TCPConn
	c, err = net.DialTCP("tcp", nil, addr)
	if err != nil {
		client.handleError(err, "Connect", "Unable to dial")
		return nil, err
	}

	client.Connection = NewClientConnection(c)
	err = client.Connection.Encoder.Encode(info)
	if err != nil {
		client.handleError(err, "Connect", "Encode player info failed")
		return nil, err
	}

	info = &PlayerInfo{}
	err = client.Connection.Decoder.Decode(info)
	if err != nil {
		client.handleError(err, "Connect", "Decode player info failed")
		return nil, err
	}

	log.Println("[CamarettoClient.Connect] Completed: {", info.Index, ",", info.Name, "}")
	return info, nil
}

// @desc: Disconnect from server
func (client *CamarettoClient) Disconnect() error {
	return nil
}

func (client *CamarettoClient) SendMessage(msg *Message) error {
	var err error

	err = client.Connection.Encoder.Encode(msg)
	if err != nil {
		client.handleError(err, "SendMessage", "Encode message failed")
		return err
	}

	return nil
}

// @desc:
func (client *CamarettoClient) ReceiveMessage(io chan *Message, e chan error) {
	var err error
	var msg *Message = &Message{}

	err = client.Connection.Decoder.Decode(msg)
	if err != nil {
		client.handleError(err, "ReceiveUpdate", "Decode message failed")
		e <- err
	} else {
		io <- msg
	}
}
