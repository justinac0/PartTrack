package auth

import (
	"PartTrack/internal/resource/sessions"
	"PartTrack/internal/resource/users"
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := CheckSession(c)
		if err != nil {
			c.Request().Header.Add("HX-Redirect", "/")
			return c.NoContent(http.StatusNetworkAuthenticationRequired)
		}

		return next(c)
	}
}

func CheckSession(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	cookie, err := c.Cookie("session")
	if err != nil {
		return errors.New("session does not exist")
	}

	sessionStore := sessions.NewStore()
	session, err := sessionStore.GetBySessionId(ctx, cookie.Value)
	if err != nil {
		return err
	}

	if cookie.Value != session.SessionId {
		return errors.New("session id's don't match")
	}

	if cookie.Expires.Before(*session.Expiry) {
		return errors.New("session has expired")
	}

	return nil
}

func Setup(e *echo.Echo, userHandler *users.Handler, sessionHandler *sessions.Handler) {
	e.POST("/signin", userHandler.SignIn)
	e.GET("/signout", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})
	e.GET("/register", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})
}
