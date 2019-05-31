package utils

import (
	"database/sql"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	elastic "gopkg.in/olivere/elastic.v5"
)

// PgDB is the client to connet to postgres
var PgDB *sql.DB

// SQLUser is the username to connect to postgres
var SQLUser string

// SQLPassword is password to connect to postgres
var SQLPassword string

// SQLDb is the database name to connect to
var SQLDb string

// SQLHost is the hostname of postgres server
var SQLHost string

// ElasticURL is the URL along with username and password (if applicable) to connect to elasticsearch cluster
var ElasticURL string

// Elasticconn is the client to connect to elasticsearch
var Elasticconn *elastic.Client

// LogLevel to set the log level. allowed levels are INFO, DEBUG, ERROR
var LogLevel string

// Admins holds the email IDs of super admins
var Admins []string

// MovieIndex is the name of elasticsearch index where movie data is saved
var MovieIndex string

// ReadEnvironmentVariables reads and sets the env vars
func ReadEnvironmentVariables() {
	if os.Getenv("IMDB_ENV") != "PRODUCTION" {
		filePath := os.Args[1:]
		err := godotenv.Load(filePath...)
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
	SQLHost = os.Getenv("SQLHost")
	if SQLHost == "" {
		log.Fatalln("SQLHost env var not set")
	}
	SQLUser = os.Getenv("SQLUser")
	if SQLUser == "" {
		log.Fatalln("SQLUser env var not set")
	}
	SQLPassword = os.Getenv("SQLPassword")
	if SQLPassword == "" {
		log.Fatalln("SQLPassword env var not set")
	}
	SQLDb = os.Getenv("SQLDb")
	if SQLDb == "" {
		log.Fatalln("SQLDb env var not set")
	}
	ElasticURL = os.Getenv("ElasticURL")
	if ElasticURL == "" {
		log.Fatalln("ElasticURL env var not set")
	}
	LogLevel = os.Getenv("LogLevel")
	if LogLevel == "" {
		log.Fatalln("LogLevel env var not set. Options: INFO, DEBUG, ERROR")
	}
	adminsEnvVar := os.Getenv("Admins")
	if adminsEnvVar == "" {
		log.Fatalln("Admins env var not set. When specifying multiple use , as delimeter")
	}
	Admins = strings.Split(adminsEnvVar, ",")
	MovieIndex = os.Getenv("MovieIndex")
	if LogLevel == "" {
		log.Fatalln("MovieIndex env var not set")
	}
}
