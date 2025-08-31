package components

import (
	"PartTrack/internal"
	"PartTrack/internal/db"
	"PartTrack/internal/resource/components/model"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type StoreError error

var (
	ErrPageOutOfBounds StoreError = errors.New("page out of bounds")
)

type ComponentStore struct {
	db *sql.DB
}

func NewStore() *ComponentStore {
	return &ComponentStore{
		db: db.GetHandle(),
	}
}

func (s *ComponentStore) GetOne(ctx context.Context, id uint64) (*model.Component, error) {
	row := s.db.QueryRowContext(ctx, "SELECT * FROM components WHERE id = $1;", id)

	var comp model.Component
	err := row.Scan(&comp.Id, &comp.AddedBy, &comp.Name, &comp.Description, &comp.Footprint, &comp.Manufacturer, &comp.Supplier, &comp.Amount, &comp.CreatedAt, &comp.DeletedAt)
	if err != nil {
		return nil, err
	}

	return &comp, nil
}

func (s *ComponentStore) GetPaginated(ctx context.Context, offset int64, search string) (*internal.Page[model.Component], error) {
	searchIncluded := len(search) > 0

	var countRow *sql.Row
	if searchIncluded {
		countRow = s.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM components
			WHERE id::text LIKE '%' || $1 || '%'
			OR name LIKE '%' || $1 || '%'
			OR description LIKE '%' || $1 || '%'
			OR footprint LIKE '%' || $1 || '%'
			OR manufacturer LIKE '%' || $1 || '%'
			OR supplier LIKE '%' || $1 || '%';`, search)
	} else {
		countRow = s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM components;")
	}

	var list internal.Page[model.Component]
	list.Offset = offset

	var rowCount int64
	err := countRow.Scan(&rowCount)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	list.ResultCount = rowCount

	if offset >= list.GetMaxPages() {
		return nil, ErrPageOutOfBounds
	}

	rowOffset := offset * internal.PAGINATION_COUNT
	var rows *sql.Rows
	if searchIncluded {
		rows, err = s.db.QueryContext(ctx,
			`SELECT id, added_by, name, description, footprint, manufacturer, supplier, amount, created_at, deleted_at
			FROM components
			WHERE id::text LIKE '%' || $1 || '%'
				OR name LIKE '%' || $1 || '%'
				OR description LIKE '%' || $1 || '%'
				OR footprint LIKE '%' || $1 || '%'
				OR manufacturer LIKE '%' || $1 || '%'
				OR supplier LIKE '%' || $1 || '%'
			LIMIT $2 OFFSET $3;`,
			search, internal.PAGINATION_COUNT, rowOffset)
	} else {
		rows, err = s.db.QueryContext(ctx, `SELECT * FROM components LIMIT $1 OFFSET $2;`, internal.PAGINATION_COUNT, rowOffset)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var comp model.Component
		err := rows.Scan(&comp.Id, &comp.AddedBy, &comp.Name, &comp.Description, &comp.Footprint, &comp.Manufacturer, &comp.Supplier, &comp.Amount, &comp.CreatedAt, &comp.DeletedAt)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		list.Items = append(list.Items, comp)
	}

	return &list, nil
}
