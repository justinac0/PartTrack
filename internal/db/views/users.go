package views

import (
	"PartTrack/internal"
	"PartTrack/internal/crypt"
	"PartTrack/internal/db/models"
	"PartTrack/internal/db/stores"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type UsersHandler struct {
	store *stores.UsersStore
}

func NewUsersHandler() *UsersHandler {
	return &UsersHandler{
		store: stores.NewUsersStore(),
	}
}

func (h *UsersHandler) SignIn(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	username := c.FormValue("username")
	password := c.FormValue("password")

	user, err := h.store.GetByUsername(ctx, username)
	if err != nil {
		return internal.OnError(c, http.StatusUnprocessableEntity, err.Error())
	}

	if !crypt.VerifyPassword(password, user.PasswordHash) {
		return internal.OnError(c, http.StatusUnprocessableEntity, "incorrect username/password")
	}

	err = recreateSession(c, ctx, user.Id)
	if err != nil {
		return internal.OnError(c, http.StatusUnprocessableEntity, err.Error())
	}

	c.Response().Header().Set("HX-Redirect", "/protected/dashboard")
	return c.NoContent(http.StatusOK)
}

func (h *UsersHandler) SignOut(c echo.Context) error {
	cookie, err := c.Cookie("session")
	if err != nil {
		c.Response().Header().Add("HX-Redirect", "/")
		return c.NoContent(http.StatusOK)
	}

	cookie.Value = ""
	cookie.Path = "/"
	cookie.Expires = time.Unix(0, 0)

	c.SetCookie(cookie)

	c.Response().Header().Add("HX-Redirect", "/")
	return c.NoContent(http.StatusOK)
}

func (h *UsersHandler) Register(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	email := c.FormValue("email")
	username := c.FormValue("username")
	password := c.FormValue("password")
	retry_password := c.FormValue("retry_password")

	if password != retry_password {
		return internal.OnError(c, http.StatusBadRequest, "passwords don't match")
	}

	passHash, err := crypt.HashPassword(password)
	if err != nil {
		return internal.OnError(c, http.StatusBadRequest, err.Error())
	}

	now := time.Now().UTC()
	data := models.User{
		Username:     username,
		Email:        email,
		PasswordHash: passHash,
		Role:         models.UserRoleGuest,
		CreatedAt:    &now,
	}

	_, err = h.store.Create(ctx, data)
	if err != nil {
		return internal.OnError(c, http.StatusBadRequest, err.Error())
	}

	user, err := h.store.GetByUsername(ctx, username)
	if err != nil {
		return internal.OnError(c, http.StatusBadRequest, err.Error())
	}
	err = recreateSession(c, ctx, user.Id)
	if err != nil {
		return internal.OnError(c, http.StatusBadRequest, err.Error())
	}

	c.Response().Header().Set("HX-Redirect", "/protected/dashboard")
	return c.NoContent(http.StatusOK)
}

func (h *UsersHandler) GetUserById(c echo.Context) error {
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

func (h *UsersHandler) WhoAmI(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	if ValidateSession(c) != nil {
		c.Response().Header().Set("HX-Redirect", "/")
		return c.NoContent(http.StatusUnauthorized)
	}

	cookie, err := c.Cookie("session")
	if err != nil {
		c.Response().Header().Set("HX-Redirect", "/")
		return c.NoContent(http.StatusUnauthorized)
	}

	sessionStore := stores.NewSessionsStore()
	session, err := sessionStore.GetBySessionId(ctx, cookie.Value)
	if err != nil {
		c.Response().Header().Set("HX-Redirect", "/")
		return c.NoContent(http.StatusUnauthorized)
	}

	user, err := h.store.GetOne(ctx, session.UserId)
	if err != nil {
		c.Response().Header().Set("HX-Redirect", "/")
		return c.NoContent(http.StatusUnauthorized)
	}

	return c.String(http.StatusOK, fmt.Sprintf("%s [%s]", user.Username, strings.ToUpper(string(user.Role))))
}

// ----------------------------------------------------------------------------
// TODO: move to somewhere else. not a view
func ValidateSession(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	cookie, err := c.Cookie("session")
	if err != nil {
		fmt.Println("a", err)
		return SessionCookieNotSet
	}

	sessionStore := stores.NewSessionsStore()
	session, err := sessionStore.GetBySessionId(ctx, cookie.Value)
	if err != nil {
		return SessionNotFound
	}

	userStore := stores.NewUsersStore()
	user, err := userStore.GetOne(ctx, session.UserId)
	if err != nil {
		return err
	}

	if user.Id != session.UserId {
		return SessionNotFound
	}

	if cookie.Value != session.SessionId {
		return SessionIdInvalid
	}

	if time.Now().After(*session.ExpiresAt) {
		return SessionExpired
	}

	return nil
}

// delete old session in db if exists, creates new one and sets browsers cookie
func recreateSession(c echo.Context, ctx context.Context, userId uint64) error {
	sessionStore := stores.NewSessionsStore()
	err := sessionStore.Delete(ctx, userId)
	if err != nil {
		panic(err)
	}

	// TODO: write helper for session creation
	// create session
	expiry := time.Now().Add(time.Hour * 24).UTC()
	now := time.Now().UTC()

	sessionId := uuid.New()

	session, err := sessionStore.Create(ctx, models.Session{
		UserId:    userId,
		SessionId: sessionId.String(), // TODO: generate unique
		ExpiresAt: &expiry,
		CreatedAt: &now,
	})
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	cookie := http.Cookie{
		Name:     "session",
		Value:    session.SessionId,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  *session.ExpiresAt,
	}

	c.SetCookie(&cookie)

	return nil
}
