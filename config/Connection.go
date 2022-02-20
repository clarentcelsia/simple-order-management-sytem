package config

import (
	"database/sql"

	log "github.com/sirupsen/logrus"
)

var dbpostgres *sql.DB

func ConnectDB() *sql.DB {
	// var db *sql.DB
	db, err := ConnectDBPostgres()
	if err != nil {
		log.Fatal(err)
	}
	dbpostgres = db
	return db
}

func GetConnection() (*sql.DB, error) {
	if err := dbpostgres.Ping(); err != nil {
		log.Info("Can't ping to database")
		return nil, err
	}
	return dbpostgres, nil
}
