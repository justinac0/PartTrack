package db

type Store[T any] interface {
	GetAll() ([]T, error)
	GetOne(id int64) (*T, error)
	Add(data T) (*T, error)
	Delete(id int64, data T) (*T, error)
	Update(id int64, data T) (*T, error)
}
