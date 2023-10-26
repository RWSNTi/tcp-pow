package tcpclient

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"pow-client/internal/models"
	"pow-client/internal/services"
)

const (
	serverType          = "tcp"
)

type tcpClient struct {
	tcpConn *net.TCPConn
	address *net.TCPAddr
	encoder *json.Encoder
	decoder *json.Decoder

	challengeService services.ChallengeService
}

type TcpClient interface {
	SendRequest(header, payload string) (err error)
	Connect() error
	Disconnect() error
}

// Creating a new tcp client to handle connections to server
func NewClient(address string, challengeService services.ChallengeService) (client TcpClient, err error) {
	serverAddress, err := net.ResolveTCPAddr(serverType, address)
	if err != nil {
		return nil, fmt.Errorf("resolving address failed: %s", err.Error())
	}

	newClient := tcpClient{
		address:          serverAddress,
		challengeService: challengeService,
	}

	return &newClient, err
}

func (c *tcpClient) Connect() (err error) {
	c.tcpConn, err = net.DialTCP(serverType, nil, c.address)
	if err != nil {
		return
	}

	log.Println("connected to server:", c.address.String())

	c.encoder = json.NewEncoder(c.tcpConn)
	c.decoder = json.NewDecoder(c.tcpConn)

	return
}

// SendRequest is used to send initial request to server
func (c *tcpClient) SendRequest(header, payload string) (err error) {
	var tcpRequest models.TCPExchange
	tcpRequest.Header = header
	tcpRequest.Payload = payload

	err = c.encoder.Encode(tcpRequest)
	if err != nil {
		return
	}

	return c.router()
}

func (c *tcpClient) Disconnect() error {
	return c.tcpConn.Close()
}
