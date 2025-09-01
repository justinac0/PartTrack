package views

import (
	"PartTrack/internal/db/stores"
	"errors"
)

type SessionHandler struct {
	store *stores.SessionsStore
}

func NewSessionHandler() *SessionHandler {
	return &SessionHandler{
		store: stores.NewSessionsStore(),
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
