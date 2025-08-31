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

func (s *ComponentStore) GetPaginated(ctx context.Context, page int64) (*internal.Page[model.Component], error) {
	countRow := s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM components;")

	var list internal.Page[model.Component]
	list.PageCount = page

	var rowCount int64

	err := countRow.Scan(&rowCount)
	if err != nil {
		return nil, err
	}

	list.MaxPage = rowCount / model.PAGINATION_SIZE

	if page > list.MaxPage {
		return nil, errors.New("page out of bounds")
	}

	offset := page * model.PAGINATION_SIZE
	rows, err := s.db.QueryContext(ctx, `SELECT * FROM components LIMIT $1 OFFSET $2;`, model.PAGINATION_SIZE, offset)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var comp model.Component
		err := rows.Scan(&comp.Id, &comp.AddedBy, &comp.Name, &comp.Description, &comp.Footprint, &comp.Manufacturer, &comp.Supplier, &comp.Amount, &comp.CreatedAt, &comp.DeletedAt)
		if err != nil {
			fmt.Println("d")
			return nil, err
		}
		list.Items = append(list.Items, comp)
	}

	return &list, nil
}
