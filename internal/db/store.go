package db

import (
	"PartTrack/internal"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type DBErr error

var (
	DBErrNotFound DBErr = errors.New("db could not find requested data")
)

func GetPaginated[T any](db *sql.DB, ctx context.Context) (*internal.Page[T], DBErr) {

	return nil, nil
}

func GetOne[T any](db *sql.DB, resource string, id int64) (*T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	row := db.QueryRowContext(ctx, fmt.Sprintf("SELECT * FROM %s WHERE id = $1;", resource), id)

	var data T
	err := row.Scan(&data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func Create[T any](db *sql.DB, ctx context.Context, data T) (*T, error) {
	return nil, nil
}

func Delete[T any](db *sql.DB, ctx context.Context, id int64, data T) (*T, error) {
	return nil, nil
}

func Update[T any](db *sql.DB, ctx context.Context, id int64, data T) (*T, error) {
	return nil, nil
}
