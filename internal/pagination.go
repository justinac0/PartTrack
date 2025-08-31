package internal

type Page[T any] struct {
	Items     []T
	PageCount int64
	MaxPage   int64
}

const PAGINATION_COUNT = 20
