package views

import (
	"PartTrack/internal"
	"PartTrack/internal/db/stores"
	"PartTrack/internal/templates"

	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type componentPath struct {
	Id string `param:"id"`
}

type ComponentsHanlder struct {
	Store *stores.ComponentStore
}

func NewComponentsHandler() *ComponentsHanlder {
	return &ComponentsHanlder{
		Store: stores.NewComponentsStore(),
	}
}

func (h *ComponentsHanlder) GetOne(c echo.Context) error {
	path := new(componentPath)
	if err := c.Bind(path); err != nil {
		return internal.OnError(c, http.StatusUnprocessableEntity, err.Error())
	}

	id, err := strconv.ParseUint(path.Id, 10, 64)
	if err != nil {
		return internal.OnError(c, http.StatusUnprocessableEntity, err.Error())
	}

	comp, err := h.Store.GetOne(id)
	if err != nil {
		return internal.OnError(c, http.StatusUnprocessableEntity, err.Error())
	}

	c.Response().Header().Add("HX-Redirect", fmt.Sprintf("/protected/components/%d", comp.Id))
	return internal.RenderTempl(c, http.StatusOK, templates.ComponentView(*comp))
}

func (h *ComponentsHanlder) ComponentsPage(c echo.Context) error {
	search := c.QueryParam("search")

	path := new(componentPath)
	if err := c.Bind(path); err != nil {
		return internal.OnError(c, http.StatusUnprocessableEntity, err.Error())
	}

	id, err := strconv.ParseUint(path.Id, 10, 64)
	if err != nil {
		return internal.OnError(c, http.StatusUnprocessableEntity, err.Error())
	}

	page, err := h.Store.GetPage(int64(id), search)
	if err != nil {
		return internal.OnError(c, http.StatusUnprocessableEntity, err.Error())
	}

	if len(search) > 0 {
		page.SearchQuery = search
	}

	return internal.RenderTempl(c, http.StatusOK, templates.ComponentsPage(page))
}
