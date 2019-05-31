package dbConnections

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/raazcrzy/imdb/utils"
	"gopkg.in/olivere/elastic.v5"
)

func InitDbs() {
	var dbinfo string
	dbinfo = fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=disable", utils.SQLUser, utils.SQLPassword, utils.SQLDb, utils.SQLHost)
	var err error
	utils.PgDB, err = sql.Open("postgres", dbinfo)
	utils.PgDB.SetMaxOpenConns(100)
	utils.PgDB.SetMaxIdleConns(10)
	utils.PgDB.SetConnMaxLifetime(10 * time.Minute)
	if err != nil {
		log.Fatalln(err)
	}
	utils.Elasticconn, err = elastic.NewClient(elastic.SetURL(utils.ElasticURL), elastic.SetSniff(false))
	if err != nil {
		log.Fatalln(err)
	}
	t, err := utils.PgDB.Begin()
	if err != nil {
		log.Fatalln(err)
	}
	_, err = t.Exec(`
	CREATE SCHEMA IF NOT EXISTS imdb;`)
	if err != nil {
		t.Rollback()
		log.Fatalln(err)
	}
	_, err = t.Exec(`
	CREATE TABLE IF NOT EXISTS imdb.users (
		email VARCHAR(500) NOT NULL PRIMARY KEY,
		created_at integer NOT NULL,
		user_password VARCHAR(32) NOT NULL,
		user_id varchar(32) NOT NULL UNIQUE,
		role varchar(6) NOT NULL
	);`)
	if err != nil {
		t.Rollback()
		log.Fatalln(err)
	}
	_, err = utils.PgDB.Exec(`
	INSERT INTO imdb.users 
	("email", "user_password",
	"user_id", "role",
	"name", "created_at")
	VALUES ($1, 'barx', 'foox', 'admin', 'auto created', 1559217988);`, utils.Admins[0])
	err = t.Commit()
	if err != nil {
		log.Fatalln(err)
	}
}
