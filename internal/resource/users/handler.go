package users

import (
	"PartTrack/internal/db"
	"PartTrack/internal/resource/sessions"
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	store *UserStore
}

func NewHandler() *Handler {
	return &Handler{
		store: &UserStore{
			db: db.GetHandle(),
		},
	}
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func verifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (h *Handler) SignIn(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	user, err := h.store.GetByUsername(ctx, username)
	if err != nil {
		log.Println(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	if !verifyPassword(password, user.PasswordHash) {
		c.NoContent(http.StatusInternalServerError)
	}

	// create session
	expiry := time.Now().Add(time.Second * 30)
	now := time.Now()

	sessionStore := sessions.NewStore()
	session, err := sessionStore.Create(ctx, sessions.Session{
		UserId:    user.Id,
		SessionId: "default key",
		Expiry:    &expiry,
		Created:   &now,
	})
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	cookie := http.Cookie{
		Name:  "session",
		Value: session.SessionId,
	}
	c.SetCookie(&cookie)

	c.Response().Header().Add("HX-Redirect", "/dashboard")
	return c.NoContent(http.StatusOK)
}

func (h *Handler) SignOut(c echo.Context) error {
	// passHash, err := hashPassword(password)
	// if err != nil {
	// 	panic(err)
	// }

	return c.NoContent(http.StatusOK)
}

func (h *Handler) GetUserById(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	user, err := h.store.GetOne(ctx, id)
	if err != nil {
		panic(err)
	}

	log.Println(user)

	return c.NoContent(http.StatusOK)
}

func (h *Handler) GetUsers(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	users, err := h.store.GetAll(ctx)
	if err != nil {
		panic(err)
	}

	log.Println(users)

	return c.NoContent(http.StatusOK)
}
