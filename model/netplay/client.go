package netplay

import (
	"log"

	"net"
	"encoding/gob"

	"camaretto/model/component"
)

type CamarettoClient struct {
	Connection *net.TCPConn
	Encoder *gob.Encoder
	Decoder *gob.Decoder
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
func (client *CamarettoClient) Connect(addr *net.TCPAddr, info *component.PlayerInfo) (*component.PlayerInfo, error) {
	var err error

	var c *net.TCPConn
	c, err = net.DialTCP("tcp", nil, addr)
	if err != nil {
		client.handleError(err, "Connect", "Unable to dial")
		return nil, err
	}

	client.Connection = c
	client.Encoder = gob.NewEncoder(client.Connection)
	client.Decoder = gob.NewDecoder(client.Connection)

	err = client.Encoder.Encode(info)
	if err != nil {
		client.handleError(err, "Connect", "Encode player info failed")
		return nil, err
	}

	info = &component.PlayerInfo{}
	err = client.Decoder.Decode(info)
	if err != nil {
		client.handleError(err, "Connect", "Decode player info failed")
		return nil, err
	}

	log.Println("[CamarettoClient.Connect] Completed: {", info.Index, ",", info.Name, "}")
	return info, nil
}

// @desc: Disconnect from server
func (client *CamarettoClient) Disconnect() error {
	var err error

	err = client.Connection.Close()
	if err != nil {
		client.handleError(err, "Disconnect", "Closing connection failed")
		return err
	}

	client.Connection = nil
	return nil
}

// @desc: Send data "msg" (*Message) to the connection
func (client *CamarettoClient) SendMessage(msg *Message) error {
	var err error

	err = client.Encoder.Encode(msg)
	if err != nil {
		client.handleError(err, "SendMessage", "Encode message failed")
		return err
	}

	return nil
}

// @desc: Receive data from connection storing it into "io" channel
func (client *CamarettoClient) ReceiveMessage(io chan *Message, e chan error) {
	var err error
	var msg *Message = &Message{}

	err = client.Decoder.Decode(msg)
	if err != nil {
		client.handleError(err, "ReceiveUpdate", "Decode message failed")
		e <- err
	} else {
		io <- msg
	}
}
