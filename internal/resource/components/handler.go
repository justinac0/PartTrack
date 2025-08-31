package components

import (
	"PartTrack/internal"
	"PartTrack/internal/templates/components"
	"context"
	"fmt"
	"net/http"
	"strconv"
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

type componentPath struct {
	Id string `param:"id"`
}

func (h *Handler) GetPaginated(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	path := new(componentPath)
	if err := c.Bind(path); err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	id, err := strconv.ParseUint(path.Id, 10, 64)
	if err != nil {
		fmt.Println(path.Id)
		return c.NoContent(http.StatusBadRequest)
	}

	page, err := h.store.GetPaginated(ctx, int64(id))
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return internal.RenderTempl(c, http.StatusOK, components.ComponentTable(page))
}
