package components

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	store *ComponentStore
}

func NewHandler() *Handler {
	return &Handler{
		store: NewStore(),
	}
}

func (h *Handler) GetPaginated(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	page, err := h.store.GetPaginated(ctx, 0)
	if err != nil {
		panic(err)
	}

	fmt.Println(page)

	return c.NoContent(http.StatusOK)
}
