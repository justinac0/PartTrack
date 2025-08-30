package users

import (
	"PartTrack/internal/db"
	"PartTrack/internal/resource/sessions"
	"context"
	"fmt"
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

// delete old session in db if exists, creates new one and sets browsers cookie
func recreateSession(c echo.Context, ctx context.Context, userId uint64) error {
	sessionStore := sessions.NewStore()
	err := sessionStore.Delete(ctx, userId)
	if err != nil {
		fmt.Println("[err]: ", err)
	}

	// TODO: write helper for session creation
	// create session
	expiry := time.Now().Add(time.Hour * 24).UTC()
	now := time.Now().UTC()

	session, err := sessionStore.Create(ctx, sessions.Session{
		UserId:    userId,
		SessionId: "default key", // TODO: generate unique
		Expiry:    &expiry,
		Created:   &now,
	})
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	cookie := http.Cookie{
		Name:    "session",
		Value:   session.SessionId,
		Expires: *session.Expiry,
		// TODO: needs to support the following cookie attributes in prod
		// HttpOnly: true,
		// Secure: true,
	}

	c.SetCookie(&cookie)
	return nil
}

func (h *Handler) SignIn(c echo.Context) error {
	if sessions.ValidateSession(c) == nil {
		c.Response().Header().Set("HX-Redirect", "/dashboard")
		return c.NoContent(http.StatusOK)
	}

	username := c.FormValue("username")
	password := c.FormValue("password")

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	user, err := h.store.GetByUsername(ctx, username)
	if err != nil {
		return c.String(http.StatusOK, "incorrect username/password")
	}

	if !verifyPassword(password, user.PasswordHash) {
		return c.String(http.StatusOK, "incorrect username/password")
	}

	err = recreateSession(c, ctx, user.Id)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	c.Response().Header().Set("HX-Redirect", "/dashboard")
	return c.NoContent(http.StatusOK)
}

func (h *Handler) SignOut(c echo.Context) error {
	cookie, err := c.Cookie("session")
	if err != nil {
		c.Response().Header().Set("HX-Redirect", "/")
		return c.NoContent(http.StatusOK)
	}

	cookie.Value = ""
	cookie.Expires = time.Now().UTC()

	c.SetCookie(cookie)

	c.Response().Header().Set("HX-Redirect", "/")
	return c.NoContent(http.StatusOK)
}

func (h *Handler) Register(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	email := c.FormValue("email")
	username := c.FormValue("username")
	password := c.FormValue("password")

	passHash, err := hashPassword(password)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	now := time.Now().UTC()
	user := User{
		Username:     username,
		Email:        email,
		PasswordHash: passHash,
		Role:         RoleGuest,
		CreatedAt:    &now,
	}

	_, err = h.store.Create(ctx, user)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	c.Response().Header().Set("HX-Redirect", "/")
	return c.NoContent(http.StatusOK)
}

func (h *Handler) GetUserById(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	_, err = h.store.GetOne(ctx, id)
	if err != nil {
		panic(err)
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handler) GetUsers(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	_, err := h.store.GetAll(ctx)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
