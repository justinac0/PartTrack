package components

import (
	"PartTrack/internal"
	"PartTrack/internal/templates"
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

func (h *Handler) ViewOne(c echo.Context) error {
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

	comp, err := h.store.GetOne(ctx, id)
	if err != nil {
		fmt.Println(err)
		return c.NoContent(http.StatusBadRequest)
	}

	c.Response().Header().Add("HX-Redirect", fmt.Sprintf("/protected/components/%d", comp.Id))
	return internal.RenderTempl(c, http.StatusOK, templates.ComponentView(*comp))
}

func (h *Handler) ViewComponents(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	search := c.QueryParam("search")

	path := new(componentPath)
	if err := c.Bind(path); err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	id, err := strconv.ParseUint(path.Id, 10, 64)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	page, err := h.store.GetPaginated(ctx, int64(id), search)
	if err != nil {
		fmt.Println(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	if len(search) > 0 {
		page.SearchQuery = search
	}

	return internal.RenderTempl(c, http.StatusOK, templates.ComponentsPage(page))
}
