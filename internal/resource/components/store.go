package components

import (
	"PartTrack/internal/db"
	"context"
	"database/sql"
)

type ComponentStore struct {
	db *sql.DB
}

func NewStore() *ComponentStore {
	return &ComponentStore{
		db: db.GetHandle(),
	}
}

func (s *ComponentStore) GetPaginated(ctx context.Context, page uint64) (*ComponentsPaginated, error) {

	return nil, nil
}
