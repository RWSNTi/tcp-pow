package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"pow-server/config"
	"pow-server/internal/services"
	"pow-server/internal/tcpserver"
	"pow-server/pkg/logging"
	"syscall"
)

func main() {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT)

	//starting logger
	logWriter := logging.GetLogger()
	logger := logWriter.Logger

	logger.Println("starting server")

	// loading config from file and env
	err := config.LoadConfig()
	if err != nil {
		logger.Panic(err)
	}

	// creating a blockchain service and making a blockchain with genesis block
	blockchainService := services.NewBlockchainService(config.Config.Difficulty, config.Config.TargetValue)
	blockchainService.CreateBlockchain()

	// creating quote service, which will get quotes from online API resource
	quoteService := services.NewQuotesService(config.Config.QuotesApiUrl)

	// create server instance
	tcpServer, err := tcpserver.NewServer(fmt.Sprintf("%s:%d", config.Config.ServerHost, config.Config.ServerPort), blockchainService, quoteService)
	if err != nil {
		logger.Panicf("tcp server start error: %s", err.Error())
	}

	//close connection and logfile before shutdown
	defer func() {
		err = tcpServer.CloseListener()
		if err != nil {
			log.Println("closing listener err:", err.Error())
		}

		err = logWriter.CloseFile()
		if err != nil {
			log.Println("closing file err:", err.Error())
		}
	}()

	//start accepting connection
	go func() {
		err = tcpServer.StartHandler()
		if err != nil {
			logger.Panicf("handler starting err: %s", err.Error())
		}
	}()

	<-shutdown
}
