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

func (h *Handler) GetPaginated(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	idStr := c.QueryParam("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		fmt.Println(idStr)
		return c.NoContent(http.StatusBadRequest)
	}

	_, err = h.store.GetPaginated(ctx, id)
	if err != nil {
		panic(err)
	}

	return internal.RenderTempl(c, http.StatusOK, components.ComponentTable())
}
