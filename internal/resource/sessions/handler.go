package sessions

import (
	"errors"
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
	SessionAccessDenied SessionError = errors.New("session access denied")
)
