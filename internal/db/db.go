// ref pq: https://pkg.go.dev/github.com/lib/pq#section-readme
package db

import (
	"database/sql"
	"errors"
	"os"

	_ "github.com/lib/pq"
)

var handle *sql.DB

func getConnString() (string, error) {
	connectionString := os.Getenv("DATABASE_URL")
	if connectionString == "" {
		return "", errors.New("DATABASE_URL not set")
	}

	return connectionString, nil
}

func Init() error {
	if handle != nil {
		return errors.New("db connection already exists")
	}

	connStr, err := getConnString()
	if err != nil {
		panic(err)
	}

	handle, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	if err = handle.Ping(); err != nil {
		return err
	}

	return nil
}

func GetHandle() *sql.DB {
	return handle
}
