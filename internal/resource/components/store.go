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
	// countRow := s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM components;")

	// var list ComponentsPaginated
	// err := countRow.Scan(&list.Count)
	// if err != nil {
	// 	return nil, err
	// }

	// fmt.Println(page, list.Count/MAX_PAGE_SIZE)
	// if page  list.Count/MAX_PAGE_SIZE {
	// 	return nil, errors.New("page out of bounds")
	// }

	// list.Components = make([]Component, 0)

	// paginateStmt, err := s.db.PrepareContext(ctx, `SELECT * FROM components LIMIT $1 OFFSET $2;`)
	// if err != nil {
	// 	return nil, err
	// }
	// paginateStmt.ExecContext(ctx)

	return nil, nil
}
