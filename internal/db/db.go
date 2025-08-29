// ref pq: https://pkg.go.dev/github.com/lib/pq#section-readme
package db

import (
	"database/sql"
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

type config struct {
	User string `env:"DBUSER"`
	Pass string `env:"DBPASS"`
	Host string `env:"DBHOST"`
	Name string `env:"DBNAME"`
}

func InitDB() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}

	var cfg config
	err = env.Parse(&cfg)
	if err != nil {
		return err
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=verify-full", cfg.User, cfg.Pass, cfg.Host, cfg.Name)

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
