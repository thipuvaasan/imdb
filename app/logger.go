package main

import (
	mylog "log"
	"os"

	"github.com/sirupsen/logrus"
)

/*
Contains functions to initialize Log and set log level
*/

// Logger is the struct to represent Logger
type Logger struct {
	*mylog.Logger
}

var logOK *Logger

func initLogger() {
	logOK = &Logger{
		mylog.New(os.Stdout, "apb_acc", mylog.LstdFlags|mylog.Lshortfile),
	}
}

func getLogLevel(LogLevel string) logrus.Level {
	if LogLevel == "INFO" {
		return logrus.InfoLevel
	} else if LogLevel == "DEBUG" {
		return logrus.DebugLevel
	} else if LogLevel == "ERROR" {
		return logrus.ErrorLevel
	}
	return logrus.ErrorLevel
}
