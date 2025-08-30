package sessions

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	store *SessionStore
}

func NewHandler() *Handler {
	return &Handler{
		store: NewStore(),
	}
}

type SessionError error

var (
	SessionExpired      SessionError = errors.New("session has expired")
	SessionIdInvalid    SessionError = errors.New("session id is invalid")
	SessionCookieNotSet SessionError = errors.New("session cookie not found")
	SessionNotFound     SessionError = errors.New("session not found on db")
)

func ValidateSession(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	cookie, err := c.Cookie("session")
	if err != nil {
		return SessionCookieNotSet
	}

	sessionStore := NewStore()
	session, err := sessionStore.GetBySessionId(ctx, cookie.Value)
	if err != nil {
		return SessionNotFound
	}

	if cookie.Value != session.SessionId {
		return SessionIdInvalid
	}

	fmt.Println(session.Expiry, time.Now().UTC())

	if time.Now().After(*session.Expiry) {
		return SessionExpired
	}

	return nil
}
