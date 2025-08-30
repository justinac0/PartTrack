package db

import "context"

type Store[T any] interface {
	GetAll(ctx context.Context) ([]T, error)
	GetOne(ctx context.Context, id int64) (*T, error)
	Create(ctx context.Context, data T) (*T, error)
	Delete(ctx context.Context, id int64, data T) (*T, error)
	Update(ctx context.Context, id int64, data T) (*T, error)
}
