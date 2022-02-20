package config

import (
	"os"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"

	"database/sql"
	"errors"
	"fmt"
)

//Config detail
func Environment() string {
	env := viper.GetString("environment")
	return env
}

func Hostname() string {
	env := viper.GetString("host")
	return env
}

func ListenAndServeServerPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = viper.GetString("port")
	}
	return ":" + port
}

func ProjectName() string {
	env := viper.GetString("token.project-name")
	return env
}

func ProjectSecret() string {
	env := viper.GetString("token.project-secret")
	return env
}

func ConnectDBPostgres() (*sql.DB, error) {
	return SetupDB("database.server", "database.portdb", "database.username", "database.password", "database.dbname")
}

func SetupDB(server string, portdb string, username string, password string, name string) (*sql.DB, error) {
	dbserver := viper.GetString(server)
	if dbserver == "" {
		err := errors.New("undefined config")
		return nil, err
	}
	psql := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbserver,
		viper.GetString(portdb),
		viper.GetString(username),
		viper.GetString(password),
		viper.GetString(name),
	)
	db, errdb := sql.Open("postgres", psql)
	if errdb != nil { //validates if the database exist or not
		return nil, errors.New("can't connect to database")
	}
	if err := db.Ping(); err != nil { //make sure database has been connected
		return nil, errors.New("database not connected")
	}
	return db, nil
}
