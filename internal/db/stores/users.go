package stores

import (
	"PartTrack/internal/db"
	"PartTrack/internal/db/models"
	"context"
	"database/sql"
)

type UsersStore struct {
	db *sql.DB
}

func NewUsersStore() *UsersStore {
	return &UsersStore{
		db: db.GetHandle(),
	}
}

func (s *UsersStore) GetAll(ctx context.Context) ([]models.User, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT id, username, password_hash, role, created, deleted FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []models.User{}
	for rows.Next() {
		user := models.User{}
		err := rows.Scan(&user.Id, &user.Username, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.DeletedAt)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (s *UsersStore) GetOne(ctx context.Context, id uint64) (*models.User, error) {
	user := models.User{}
	row := s.db.QueryRowContext(ctx, "SELECT  id, email, username, password_hash, role, created_at, deleted_at FROM users WHERE id = $1", id)
	err := row.Scan(&user.Id, &user.Email, &user.Username, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.DeletedAt)
	if err != nil {
		return nil, err
	}

	return &user, err
}

func (s *UsersStore) Create(ctx context.Context, data models.User) (*models.User, error) {
	statement, err := s.db.PrepareContext(ctx, "INSERT INTO users (email, username, password_hash, role, created_at) VALUES ($1, $2, $3, $4, $5)")
	if err != nil {
		return nil, err
	}

	_, err = statement.ExecContext(ctx, &data.Email, &data.Username, &data.PasswordHash, &data.Role, &data.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &data, err
}

func (s *UsersStore) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	user := models.User{}
	row := s.db.QueryRowContext(ctx, "SELECT id, email, username, password_hash, role, created_at, deleted_at FROM users WHERE username = $1", username)
	err := row.Scan(&user.Id, &user.Email, &user.Username, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.DeletedAt)
	if err != nil {
		return nil, err
	}

	return &user, err
}
