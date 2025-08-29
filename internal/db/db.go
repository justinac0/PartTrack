// ref pq: https://pkg.go.dev/github.com/lib/pq#section-readme
package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB() error {
	// TODO: use env file to load in connection info
	connStr := "postgres://default:password@localhost/PartTrackDB?sslmode=verify-full"

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	return nil
}

func GetHandle() *sql.DB {
	return db
}

func CloseDB() {
	if db != nil {
		db.Close()
	}
}
