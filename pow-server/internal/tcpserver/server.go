package tcpserver

import (
	"encoding/json"
	"log"
	"net"
	"pow-server/config"
	"pow-server/internal/services"
)

const (
	serverType = "tcp"
)

type tcpServer struct {
	listener          net.Listener
	blockchainService services.BlockchainService
	quotesServer      services.QuotesService
	encoder           *json.Encoder
	decoder           *json.Decoder
}

type TcpServer interface {
	StartHandler() error
	CloseListener() error
}

// NewServer creates a new tcp server instance to handle connections from clients
func NewServer(serverAddress string, blockchainService services.BlockchainService, quotesServer services.QuotesService) (
	server TcpServer, err error) {
	listener, err := net.Listen(serverType, serverAddress)
	if err != nil {
		return
	}

	log.Println("start listening")

	server = &tcpServer{
		listener:          listener,
		blockchainService: blockchainService,
		quotesServer:      quotesServer,
	}

	return
}

// StartHandler uses infinite loop to serve multiple client connections
func (s *tcpServer) StartHandler() error {
	log.Println("starting handler at port:", config.Config.ServerPort)
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Println("err accepting connection from client:", err.Error())
		}

		log.Printf("client %s connected to server", conn.RemoteAddr())

		s.decoder = json.NewDecoder(conn)
		s.encoder = json.NewEncoder(conn)

		go func() {
			err = s.router(conn)
		}()

		if err != nil {
			log.Println("handler err:", err.Error())
		}
	}
}

// CloseListener is used to close listener before shutdown
func (s *tcpServer) CloseListener() error {
	err := s.listener.Close()
	if err == nil {
		log.Println("listener closed")
	}

	return err
}
