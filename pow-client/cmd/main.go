package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"pow-client/config"
	"pow-client/internal/services"
	"pow-client/internal/tcpclient"
	"pow-client/pkg/logging"
	"syscall"
	"time"
)

const (
	requestChallengeHeader = "challenge"
	reconnectInterval = 5
)

func main() {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT)

	//starting logger
	logWriter := logging.GetLogger()
	logger := logWriter.Logger
	//close connection and logfile before shutdown
	defer logWriter.CloseFile()

	logger.Println("start client")

	// loading config from file and env
	err := config.LoadConfig()
	if err != nil {
		logger.Panic(err)
	}

	generateBlockService := services.NewGenerateBlockService(config.Config.MaxSolveIterationLimit)
	challengeService := services.NewChallengeService(generateBlockService)

	// creates client instance
	tcpClient, err := tcpclient.NewClient(fmt.Sprintf("%s:%d", config.Config.ServerHost, config.Config.ServerPort), challengeService)
	if err != nil {
		log.Panicf("tcp client create error: %s", err.Error())
	}

	go func() {
		for {
			func() {
				defer func() {
					err = tcpClient.Disconnect()
					if err != nil {
						log.Println("disconnect tcp client err:", err.Error())
					}

					log.Println("connection to server closed")

					// sleep used for better results visualisation in console
					time.Sleep(time.Second * reconnectInterval)
				}()

				err = tcpClient.Connect()
				if err != nil {
					log.Println("error connecting to server:", err)
					return
				}

				err = tcpClient.SendRequest(requestChallengeHeader, "")
				if err != nil {
					log.Println("error solving the problem:", err)
				}
			}()
		}
	}()

	<-shutdown
}
