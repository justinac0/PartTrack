package sessions

import (
	"PartTrack/internal/db"
	"context"
	"database/sql"
	"time"
)

type Session struct {
	SessionId string     `json:"session_id"`
	UserId    int64      `json:"user_id"`
	Expiry    *time.Time `json:"expiry"`
	Created   *time.Time `json:"created"`
}

type SessionStore struct {
	db *sql.DB
}

func NewStore() *SessionStore {
	return &SessionStore{
		db: db.GetHandle(),
	}
}

func (s *SessionStore) GetBySessionID(ctx context.Context, id string) (*Session, error) {
	return nil, nil
}
