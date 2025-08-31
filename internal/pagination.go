package internal

import (
	"math"
)

type Page[T any] struct {
	Items       []T
	Offset      int64
	ResultCount int64
	SearchQuery string
}

func (p *Page[T]) GetIndex() int64 {
	return p.Offset * PAGINATION_COUNT
}

func (p *Page[T]) GetMaxPages() int64 {
	return int64(math.Ceil(float64(p.ResultCount) / PAGINATION_COUNT))
}

func (p *Page[T]) NextPageIndex() int64 {
	if (p.Offset+1)*PAGINATION_COUNT > p.ResultCount {
		return p.Offset
	}

	return p.Offset + 1
}

func (p *Page[T]) PrevPageIndex() int64 {
	return max(p.Offset-1, 0)
}

const PAGINATION_COUNT = 50
