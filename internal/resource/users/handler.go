package users

import (
	"PartTrack/internal/crypt"
	"PartTrack/internal/resource/sessions"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	store *UserStore
}

func NewHandler() *Handler {
	return &Handler{
		store: NewStore(),
	}
}

// delete old session in db if exists, creates new one and sets browsers cookie
func recreateSession(c echo.Context, ctx context.Context, userId uint64) error {
	sessionStore := sessions.NewStore()
	err := sessionStore.Delete(ctx, userId)
	if err != nil {
		panic(err)
	}

	// TODO: write helper for session creation
	// create session
	expiry := time.Now().Add(time.Hour * 24).UTC()
	now := time.Now().UTC()

	sessionId := uuid.New()

	session, err := sessionStore.Create(ctx, sessions.Session{
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

func (h *Handler) SignIn(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	username := c.FormValue("username")
	password := c.FormValue("password")

	user, err := h.store.GetByUsername(ctx, username)
	if err != nil {
		return c.String(http.StatusOK, "incorrect username/password")
	}

	if !crypt.VerifyPassword(password, user.PasswordHash) {
		return c.String(http.StatusOK, "incorrect username/password")
	}

	err = recreateSession(c, ctx, user.Id)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	c.Response().Header().Set("HX-Redirect", "/protected/dashboard")
	return c.NoContent(http.StatusOK)
}

func (h *Handler) SignOut(c echo.Context) error {
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

func (h *Handler) Register(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	email := c.FormValue("email")
	username := c.FormValue("username")
	password := c.FormValue("password")
	retry_password := c.FormValue("retry_password")

	if password != retry_password {
		return c.String(http.StatusOK, "passwords don't match")
	}

	passHash, err := crypt.HashPassword(password)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	now := time.Now().UTC()
	data := User{
		Username:     username,
		Email:        email,
		PasswordHash: passHash,
		Role:         RoleGuest,
		CreatedAt:    &now,
	}

	_, err = h.store.Create(ctx, data)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	user, err := h.store.GetByUsername(ctx, username)
	if err != nil {
		panic(err)
	}
	err = recreateSession(c, ctx, user.Id)
	if err != nil {
		panic(err)
	}

	c.Response().Header().Set("HX-Redirect", "/protected/dashboard")
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

func (h *Handler) WhoAmI(c echo.Context) error {
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

	sessionStore := sessions.NewStore()
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

func ValidateSession(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	cookie, err := c.Cookie("session")
	if err != nil {
		fmt.Println("a", err)
		return sessions.SessionCookieNotSet
	}

	sessionStore := sessions.NewStore()
	session, err := sessionStore.GetBySessionId(ctx, cookie.Value)
	if err != nil {
		return sessions.SessionNotFound
	}

	userStore := NewStore()
	user, err := userStore.GetOne(ctx, session.UserId)
	if err != nil {
		fmt.Println("b", err)
		return err
	}

	if user.Id != session.UserId {
		fmt.Println("c", err)
		return sessions.SessionNotFound
	}

	if cookie.Value != session.SessionId {
		fmt.Println("d", err)
		return sessions.SessionIdInvalid
	}

	if time.Now().After(*session.ExpiresAt) {
		fmt.Println("e", err)
		return sessions.SessionExpired
	}

	return nil
}
