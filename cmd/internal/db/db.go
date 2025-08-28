package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func InitDB() error {
	if db != nil {
		panic("db has already been initialized")
	}

	var err error
	db, err = sql.Open("sqlite3", "test.db")
	if err != nil {
		return err
	}

	// initialize database if it doesn't exist
	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS testing (id INTEGER PRIMARY KEY, name VARCHAR(64))")
	if err != nil {
		return err
	}
	statement.Exec()

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
