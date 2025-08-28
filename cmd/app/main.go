package main

import (
	"PartTrack/cmd/internal/db"
)

func main() {
	err := db.InitDB()
	if err != nil {
		panic(err)
	}

	defer db.CloseDB()
}
