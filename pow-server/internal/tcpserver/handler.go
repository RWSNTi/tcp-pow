package tcpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"pow-server/config"
	"pow-server/internal/models"
)

const (
	requestChallengeHeader = "challenge"
	acceptChallengeHeader  = "accept-challenge!"
	resultHeader           = "result"
	powAnswerHeader        = "pow"
	quitHeader             = "quit"

	successResultKey = "success"
	failResultKet    = "fail"
)

// router is used to handle requests from client, based on message headers.
func (s *tcpServer) router(conn net.Conn) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.Config.ClientConnectionTimeout)

	defer func() {
		cancel()
		log.Println("closing connection for:", conn.RemoteAddr())
		err = conn.Close()
		if err != nil {
			log.Println("closing connection err", err.Error())
		}
	}()

	e := make(chan struct{})

	go func() {
		for {
			if s.decoder.Buffered() == nil {
				e <- struct{}{}
			}

			var response, message models.TCPExchange
			err = s.decoder.Decode(&message)
			if err != nil {
				e <- struct{}{}
			}

			switch message.Header {
			case requestChallengeHeader:
				response, err = s.sendChallengeToClient()
			case powAnswerHeader:
				response, err = s.processClientWorkResult(message)
			case quitHeader:
				err = nil
				e <- struct{}{}
			default:
				err = fmt.Errorf("unknown header:%s", message.Header)
				e <- struct{}{}
			}

			err = s.encoder.Encode(response)
			if err != nil || response.Header == quitHeader {
				e <- struct{}{}
			}
		}
	}()

	for {
		select{
		case <- ctx.Done():
			log.Println("context timeout")
			return
		case <- e:
			return err
		}
	}
}

func (s *tcpServer) sendChallengeToClient() (response models.TCPExchange, err error) {
	buf, err := json.Marshal(s.blockchainService.GetLastBlock())
	if err != nil {
		log.Println("marshalling challenge err")
		return response, err
	}

	response.Payload = string(buf)
	response.Header = acceptChallengeHeader

	return response, err
}

func (s *tcpServer) processClientWorkResult(message models.TCPExchange) (response models.TCPExchange, err error) {
	var block models.Block
	err = json.Unmarshal([]byte(message.Payload), &block)
	if err != nil {
		return
	}

	var result models.Result
	checkResult := s.blockchainService.AddNewBlock(block)
	if !checkResult {
		result.Key = failResultKet
		return
	}

	quote, err := s.quotesServer.GetQuote()
	if err != nil {
		return
	}

	result.Key = successResultKey
	result.Result = quote

	buf, err := json.Marshal(result)
	if err != nil {
		return
	}

	response.Payload = string(buf)
	response.Header = resultHeader

	return
}
