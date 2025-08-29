// ref pq: https://pkg.go.dev/github.com/lib/pq#section-readme
package db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type config struct {
	User string `env:"DBUSER"`
	Pass string `env:"DBPASS"`
	Host string `env:"DBHOST"`
	Name string `env:"DBNAME"`
}

var handle *sql.DB

func getConnString() (string, error) {
	err := godotenv.Load()
	if err != nil {
		return "", err
	}

	var cfg config
	err = env.Parse(&cfg)
	if err != nil {
		return "", err
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", cfg.User, cfg.Pass, cfg.Host, cfg.Name)

	return connStr, nil
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
