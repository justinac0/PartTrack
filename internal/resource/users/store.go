package users

import (
	"PartTrack/internal/db"
	"context"
	"database/sql"
	"fmt"
)

type UserStore struct {
	db *sql.DB
}

func NewStore() *UserStore {
	return &UserStore{
		db: db.GetHandle(),
	}
}

func (s *UserStore) GetAll(ctx context.Context) ([]User, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT id, username, password_hash, role, created, deleted FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		user := User{}
		err := rows.Scan(&user.Id, &user.Username, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.DeletedAt)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (s *UserStore) GetOne(ctx context.Context, id uint64) (*User, error) {
	user := User{}
	row := s.db.QueryRowContext(ctx, "SELECT  id, email, username, password_hash, role, created_at, deleted_at FROM users WHERE id = $1", id)
	err := row.Scan(&user.Id, &user.Email, &user.Username, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.DeletedAt)
	if err != nil {
		return nil, err
	}

	return &user, err
}

func (s *UserStore) Create(ctx context.Context, data User) (*User, error) {
	statement, err := s.db.PrepareContext(ctx, "INSERT INTO users (email, username, password_hash, role, created_at) VALUES ($1, $2, $3, $4, $5)")
	if err != nil {
		return nil, err
	}

	result, err := statement.ExecContext(ctx, &data.Email, &data.Username, &data.PasswordHash, &data.Role, &data.CreatedAt)
	if err != nil {
		return nil, err
	}

	fmt.Println(result)

	return &data, err
}

func (s *UserStore) GetByUsername(ctx context.Context, username string) (*User, error) {
	user := User{}
	row := s.db.QueryRowContext(ctx, "SELECT id, email, username, password_hash, role, created_at, deleted_at FROM users WHERE username = $1", username)
	err := row.Scan(&user.Id, &user.Email, &user.Username, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.DeletedAt)
	if err != nil {
		return nil, err
	}

	return &user, err
}
