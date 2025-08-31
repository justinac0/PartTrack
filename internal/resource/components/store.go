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

func (s *ComponentStore) GetPaginated(ctx context.Context, page int64, search string) (*internal.Page[model.Component], error) {
	searchIncluded := len(search) > 0
	// var filter string
	// var filterCount string

	// // TODO: there has to be a better way to do this
	// if searchIncluded {
	// 	filter = fmt.Sprintf(`SELECT id, added_by, name, description, footprint, manufacturer, supplier, amount, created_at, deleted_at FROM components
	// 		WHERE id::text ~* '%s'
	// 		OR name ~* '%s'
	// 		OR description ~* '%s'
	// 		OR footprint ~* '%s'
	// 		OR manufacturer ~* '%s'
	// 		OR supplier ~* '%s'`, search, search, search, search, search, search)

	// 	filterCount = fmt.Sprintf(`SELECT COUNT(*) FROM components
	// 		WHERE id::text ~* '%s'
	// 		OR name ~* '%s'
	// 		OR description ~* '%s'
	// 		OR footprint ~* '%s'
	// 		OR manufacturer ~* '%s'
	// 		OR supplier ~* '%s';`, search, search, search, search, search, search)
	// }

	// TODO: Ensure SQL injection isn't occuring
	var countRow *sql.Row
	if searchIncluded {
		// countRow = s.db.QueryRowContext(ctx, filterCount)
	} else {
		countRow = s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM components;")
	}

	var list internal.Page[model.Component]
	list.PageCount = page

	var rowCount int64
	err := countRow.Scan(&rowCount)
	if err != nil {
		return nil, err
	}

	list.MaxPage = rowCount / internal.PAGINATION_COUNT

	if page > list.MaxPage {
		return nil, errors.New("page out of bounds")
	}
	var rows *sql.Rows

	offset := page * internal.PAGINATION_COUNT
	if searchIncluded {
		// rows, err = s.db.QueryContext(ctx, fmt.Sprintf(`%s LIMIT $1 OFFSET $2;`, filter), model.PAGINATION_SIZE, offset)
	} else {
		rows, err = s.db.QueryContext(ctx, `SELECT * FROM components LIMIT $1 OFFSET $2;`, internal.PAGINATION_COUNT, offset)
	}
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	for rows.Next() {
		var comp model.Component
		err := rows.Scan(&comp.Id, &comp.AddedBy, &comp.Name, &comp.Description, &comp.Footprint, &comp.Manufacturer, &comp.Supplier, &comp.Amount, &comp.CreatedAt, &comp.DeletedAt)
		if err != nil {
			return nil, err
		}
		list.Items = append(list.Items, comp)
	}

	return &list, nil
}
