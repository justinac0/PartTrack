package user

import (
	"PartTrack/internal/db"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	store *UserStore
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		store: &UserStore{
			db: db.GetHandle(),
		},
	}
}

func (h *UserHandler) GetUser(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	user, err := h.store.GetOne(id)
	if err != nil {
		panic(err)
	}

	log.Println(user)

	return c.NoContent(http.StatusOK)
}

func (h *UserHandler) GetUsers(c echo.Context) error {
	users, err := h.store.GetAll()
	if err != nil {
		panic(err)
	}

	log.Println(users)

	return c.NoContent(http.StatusOK)
}
