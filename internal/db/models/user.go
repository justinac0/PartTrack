package models

import (
	"PartTrack/internal/db"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	Store *db.Store
}

func (h *UserHandler) GetUser(c echo.Context) error {
	log.Println("not implemented")
	return c.NoContent(http.StatusOK)
}

func (h *UserHandler) GetUsers(c echo.Context) error {
	rows, err := h.Store.DB.Query("SELECT * FROM user;")
	if err != nil {
		panic(err)
	}

	fmt.Println(rows)

	return c.NoContent(http.StatusOK)
}
