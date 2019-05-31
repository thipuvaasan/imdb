package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/raazcrzy/imdb/dbConnections"
	"github.com/raazcrzy/imdb/utils"
	"github.com/sirupsen/logrus"
)

// Log is configured with log level for logging
var Log = logrus.New()
var emailKey, categoryKey interface{}

// initializes env vars, Log with log levels, DB connections, and starts server on port 8000
func main() {
	emailKey = "email"
	categoryKey = "category"
	utils.ReadEnvironmentVariables()
	Log.SetLevel(getLogLevel(utils.LogLevel))
	Log.SetOutput(os.Stdout)
	initLogger()
	dbConnections.InitDbs()
	getRoutes()
	fmt.Println("Server starting...")
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}
