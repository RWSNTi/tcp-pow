package tcpclient

import (
	"fmt"
	"log"
	"pow-client/internal/models"
)

const (
	acceptChallengeHeader = "accept-challenge!"
	resultHeader          = "result"
	powAnswerHeader       = "pow"
	quitHeader            = "quit"
)

// router manages all new incoming messages from server untill a message with empty or wrong header will be received,
// or when there will be no response to send from client
func (c *tcpClient) router() (err error) {
	for {
		if !c.decoder.More() {
			return
		}

		var response, message models.TCPExchange
		err = c.decoder.Decode(&message)
		if err != nil {
			return err
		}

		switch message.Header {
		case acceptChallengeHeader:
			response, err = c.solveTheProblemHandler(message)
		case resultHeader:
			response, err = c.printResultHandler(message)
		case quitHeader:
			return nil
		default:
			err = fmt.Errorf("unknown header:%s", message.Header)
		}

		if err != nil {
			response.Header = quitHeader
			log.Println("internal error:", err.Error())
		}

		err = c.encoder.Encode(response)
		if err != nil || response.Header == quitHeader {
			return err
		}
	}
}

func (c *tcpClient) solveTheProblemHandler(message models.TCPExchange) (response models.TCPExchange, err error) {
	log.Println("started to solve the problem")
	responseBlock, err := c.challengeService.SolveTheProblem(message.Payload)
	if err != nil {
		return
	}

	response.Payload = responseBlock
	response.Header = powAnswerHeader

	return
}

func (c *tcpClient) printResultHandler(message models.TCPExchange) (response models.TCPExchange, err error) {
	err = c.challengeService.PrintResult(message.Payload)
	if err != nil {
		return
	}

	response.Header = quitHeader

	return
}
