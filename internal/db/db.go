// ref pq: https://pkg.go.dev/github.com/lib/pq#section-readme
package db

import (
	"fmt"
	"log"

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

var store *Store

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

func Init() {
	if store != nil {
		log.Fatalln("db store has already been initialized")
		return
	}

	connStr, err := getConnString()
	if err != nil {
		panic(err)
	}

	store, err = NewStore(connStr)
	if err != nil {
		panic(err)
	}

}

func GetStore() *Store {
	return store
}
