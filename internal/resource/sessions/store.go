package sessions

import (
	"PartTrack/internal/db"
	"context"
	"database/sql"
	"errors"
)

type SessionStore struct {
	db *sql.DB
}

func NewStore() *SessionStore {
	return &SessionStore{
		db: db.GetHandle(),
	}
}

func (s *SessionStore) GetBySessionId(ctx context.Context, sessionId string) (*Session, error) {
	session := Session{}
	row := s.db.QueryRowContext(
		ctx,
		`SELECT session_id, user_id, expires_at, created_at
		FROM sessions WHERE session_id = $1;`,
		sessionId)
	err := row.Scan(&session.SessionId, &session.UserId, &session.Expiry, &session.Created)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (s *SessionStore) Create(ctx context.Context, data Session) (*Session, error) {
	statement, err := s.db.PrepareContext(
		ctx,
		`INSERT INTO sessions (session_id, user_id, expires_at, created_at)
		VALUES ($1, $2, $3, $4);`)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	result, err := statement.ExecContext(ctx, &data.SessionId, &data.UserId, &data.Expiry, &data.Created)
	if err != nil {
		return nil, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if affected == 0 {
		return nil, errors.New("no rows affected on session create")
	}

	return &data, nil
}

// func (s *SessionStore) Delete(ctx context.Context, id int64, data Session) (*Session, error) {
// 	return nil, nil
// }

// func (s *SessionStore) Update(ctx context.Context, id int64, data Session) (*Session, error) {
// 	return nil, nil
// }
