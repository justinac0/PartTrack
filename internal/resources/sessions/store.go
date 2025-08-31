package sessions

import (
	"PartTrack/internal/db"
	"PartTrack/internal/db/models"
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

func (s *SessionStore) GetBySessionId(ctx context.Context, sessionId string) (*models.Session, error) {
	session := models.Session{}
	row := s.db.QueryRowContext(
		ctx,
		`SELECT session_id, user_id, expires_at, created_at
		FROM sessions WHERE session_id = $1;`,
		sessionId)
	err := row.Scan(&session.SessionId, &session.UserId, &session.ExpiresAt, &session.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (s *SessionStore) GetByUserId(ctx context.Context, userId uint64) (*models.Session, error) {
	session := models.Session{}
	row := s.db.QueryRowContext(
		ctx,
		`SELECT session_id, user_id, expires_at, created_at
		FROM sessions WHERE user_id = $1;`,
		userId)
	err := row.Scan(&session.SessionId, &session.UserId, &session.ExpiresAt, &session.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &session, nil

}

func (s *SessionStore) Create(ctx context.Context, data models.Session) (*models.Session, error) {
	statement, err := s.db.PrepareContext(
		ctx,
		`INSERT INTO sessions (session_id, user_id, expires_at, created_at)
		VALUES ($1, $2, $3, $4);`)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	result, err := statement.ExecContext(ctx, &data.SessionId, &data.UserId, &data.ExpiresAt, &data.CreatedAt)
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

func (s *SessionStore) Delete(ctx context.Context, user_id uint64) error {
	statement, err := s.db.PrepareContext(
		ctx,
		`DELETE FROM sessions WHERE user_id = $1`)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.ExecContext(ctx, user_id)
	if err != nil {
		return err
	}

	return nil
}
