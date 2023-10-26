package logging

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

type logger struct {
	Logger *log.Logger
	logFile *os.File
}

func GetLogger() *logger {
	var err error
	var logWriter logger
	logWriter.Logger = log.Default()

	// open log file
	logWriter.logFile, err = openLogFile()
	if err != nil {
		log.Panic(err)
	}

	// redirect all the output to file
	wrt := io.MultiWriter(os.Stdout, logWriter.logFile)

	// set log out put
	logWriter.Logger.SetOutput(wrt)

	// optional: log date-time, filename, and line number
	logWriter.Logger.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)

	return &logWriter
}

func(l *logger) CloseFile() {
	err := l.logFile.Close()
	if err != nil {
		log.Println("error closing log file:", err.Error())
	}
	return
}

func openLogFile() (*os.File, error) {
	// path with unix time added to name to save different log files for any new start of program
	path := fmt.Sprintf("./logs/clientlog%d.log", time.Now().Unix)

	logFile, err := os.OpenFile(path, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Panic(err)
	}

	return logFile, nil
}
