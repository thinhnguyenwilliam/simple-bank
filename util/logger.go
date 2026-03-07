package util

import (
	"io"
	"log"
	"os"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

func SetupLogger(service string) {
	date := time.Now().Format("2006-01-02")

	logFile := &lumberjack.Logger{
		Filename:   "logs/" + service + "-" + date + ".log",
		MaxSize:    50, // MB
		MaxBackups: 10,
		MaxAge:     30, // days
		Compress:   true,
	}

	multiWriter := io.MultiWriter(os.Stdout, logFile)

	log.SetOutput(multiWriter)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
