package db_test

import (
	"PartTrack/internal/db"
	"testing"
)

func TestDBConnection(t *testing.T) {
	err := db.InitDB()
	if err != nil {
		t.Logf("failed to open connection to db: %v", err)
		t.Fail()
	}

	db.CloseDB()
}
